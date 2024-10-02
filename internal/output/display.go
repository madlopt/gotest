package output

import (
	"fmt"
	"time"
)

// DisplayFinalResults outputs the final results of the IP count
func DisplayFinalResults(uniqueCount int, startTime time.Time) {
	fmt.Printf("\n\nFinal Results:\n")
	fmt.Printf("Total unique IP addresses: %d\n", uniqueCount)
	duration := time.Since(startTime)
	fmt.Printf("Execution Time: %s\n", duration)
}
