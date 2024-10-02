package processing

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/shirou/gopsutil/mem"
	"ipcounter/internal/bitmap"
	"ipcounter/internal/config"
	"ipcounter/internal/entities"
)

// CountUniqueIPs initializes the IP counting process using workers and tracks progress
func CountUniqueIPs(cfg config.Config, wg *sync.WaitGroup, uniqueCount *int, mu *sync.Mutex, startTime time.Time, fileSize int64) error {
	defer wg.Done()

	// Create a Roaring Bitmap for each worker
	bitmaps := make([]*roaring.Bitmap, cfg.NumWorkers)
	localCounts := make([]int, cfg.NumWorkers) // Track counts locally for each worker

	for i := range bitmaps {
		bitmaps[i] = roaring.New()
	}

	// Create a channel for lines and start the worker pool
	lines := make(chan string, cfg.LinesChannelCap)
	var workerWG sync.WaitGroup

	for i := 0; i < cfg.NumWorkers; i++ {
		workerWG.Add(1)
		go func(lines <-chan string, localBitmap *roaring.Bitmap, localCount *int, wg *sync.WaitGroup) {
			defer wg.Done()
			for line := range lines {
				ipInt, valid := entities.ConvertIPToUint32(line)
				if valid {
					localBitmap.Add(ipInt)
					*localCount++
				}
			}
		}(lines, bitmaps[i], &localCounts[i], &workerWG)
	}

	// Open the file and read it line-by-line with a larger buffer
	file, err := os.Open(cfg.FilePath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Printf("Error closing file: %v\n", cerr)
		}
	}()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, cfg.BufferSize)
	scanner.Buffer(buf, cfg.BufferSize)

	// Start a ticker to print progress every specified interval
	ticker := time.NewTicker(cfg.PrintInterval)
	defer ticker.Stop()

	var linesProcessed int64 = 0

	// Progress monitoring goroutine
	go func() {
		for range ticker.C {
			// Calculate the progress percentage based on lines processed and file size
			progressPercentage := int((linesProcessed * 100) / fileSize)
			// Calculate running time
			runningTime := time.Since(startTime).Round(time.Second)

			// Calculate the number of unique IP addresses found so far
			var currentUniqueCount int
			for _, count := range localCounts {
				currentUniqueCount += count
			}

			// Get current memory usage
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			vMem, _ := mem.VirtualMemory()
			totalMem := vMem.Total
			memoryUsage := float64(memStats.Alloc) / float64(totalMem) * 100
			usedMemoryMB := float64(memStats.Alloc) / 1024 / 1024

			// Print progress update
			fmt.Printf(
				"\rProcessing IP addresses: %d%% | Unique IP addresses found: %d | Running time: %s | Current memory usage: %.2f%% (%.2f MB)",
				progressPercentage, currentUniqueCount, runningTime.String(), memoryUsage, usedMemoryMB)
		}
	}()

	// Process each line in the file
	for scanner.Scan() {
		lines <- scanner.Text()
		linesProcessed++
	}

	// Properly handle scanner errors
	if err := scanner.Err(); err != nil {
		return err
	}

	close(lines) // Close the lines channel to signal workers to finish
	workerWG.Wait()

	// Merge bitmaps in parallel
	finalBitmap := bitmap.MergeBitmapsParallel(bitmaps)

	// Get total unique IP count
	mu.Lock()
	*uniqueCount = int(finalBitmap.GetCardinality())
	mu.Unlock()

	return nil
}
