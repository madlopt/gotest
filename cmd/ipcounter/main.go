package main

import (
	"fmt"
	"sync"
	"time"

	"ipcounter/internal/config"
	"ipcounter/internal/file"
	"ipcounter/internal/output"
	"ipcounter/internal/processing"
)

func main() {
	startTime := time.Now()
	cfg := config.LoadConfig()

	// Get file size
	fileSize, err := file.GetFileSize(cfg.FilePath)
	if err != nil {
		fmt.Println("Failed to get file size:", err)
		return
	}

	var uniqueCount int
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	err = processing.CountUniqueIPs(cfg, &wg, &uniqueCount, &mu, startTime, fileSize)
	if err != nil {
		fmt.Println("Failed to count unique IPs:", err)
		return
	}

	wg.Wait()
	output.DisplayFinalResults(uniqueCount, startTime)
}
