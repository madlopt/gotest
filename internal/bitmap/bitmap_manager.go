package bitmap

import (
	"github.com/RoaringBitmap/roaring"
	"runtime"
	"sync"
)

// MergeBitmapsParallel merges multiple Roaring Bitmaps in parallel
func MergeBitmapsParallel(bitmaps []*roaring.Bitmap) *roaring.Bitmap {
	numWorkers := runtime.NumCPU()
	groupSize := (len(bitmaps) + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	mergedBitmaps := make([]*roaring.Bitmap, numWorkers)

	// Optimize each individual bitmap before merging
	for _, bitmap := range bitmaps {
		if bitmap != nil {
			bitmap.RunOptimize() // Apply run-length encoding compression to reduce memory usage
		}
	}

	for i := 0; i < numWorkers; i++ {
		start := i * groupSize
		end := start + groupSize
		if end > len(bitmaps) {
			end = len(bitmaps)
		}

		if start >= len(bitmaps) {
			break
		}

		wg.Add(1)
		go func(idx, start, end int) {
			defer wg.Done()
			merged := roaring.New()
			for _, bitmap := range bitmaps[start:end] {
				merged.Or(bitmap)
			}
			mergedBitmaps[idx] = merged
		}(i, start, end)
	}

	wg.Wait()

	// Combine all intermediate merged bitmaps into the final one
	finalBitmap := roaring.New()
	for _, bitmap := range mergedBitmaps {
		if bitmap != nil {
			finalBitmap.Or(bitmap)
		}
	}

	return finalBitmap
}
