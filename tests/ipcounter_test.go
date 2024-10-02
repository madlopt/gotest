package tests

import (
	"os"
	"sync"
	"testing"
	"time"

	"ipcounter/internal/config"
	"ipcounter/internal/file"
	"ipcounter/internal/processing"
)

// createTestFile creates a temporary file with the given content and returns the file path.
func createTestFile(t *testing.T, content string) string {
	// Use os.CreateTemp instead of ioutil.TempFile for Go 1.17+
	tempFile, err := os.CreateTemp("", "ip_addresses_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	// Write the content to the temporary file
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	// Close the file and handle the error properly
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	return tempFile.Name()
}

// deleteTestFile deletes the specified temporary file.
func deleteTestFile(t *testing.T, filePath string) {
	if err := os.Remove(filePath); err != nil {
		t.Fatalf("Failed to delete temp file: %v", err)
	}
}

// TestCountUniqueIPs is an integration test that verifies the entire counting flow.
func TestCountUniqueIPs(t *testing.T) {
	// Mock data: Multiple lines with a few duplicate IP addresses
	testContent := `192.168.1.1
192.168.1.2
192.168.1.1
192.168.1.3
192.168.1.4
192.168.1.2`

	// Create a temporary test file with mock IP addresses
	filePath := createTestFile(t, testContent)
	defer deleteTestFile(t, filePath) // Ensure cleanup after test

	// Initialize the config with the test file path
	cfg := config.Config{
		FilePath:        filePath,
		BufferSize:      16 * 1024, // Use a smaller buffer for testing
		NumWorkers:      2,         // Use a couple of workers for testing
		PrintInterval:   2 * time.Second,
		LinesChannelCap: 10,
	}

	// Use a wait group and mutex for synchronization
	var uniqueCount int
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Get the size of the test file
	fileSize, err := file.GetFileSize(filePath)
	if err != nil {
		t.Fatalf("Failed to get test file size: %v", err)
	}

	startTime := time.Now()
	wg.Add(1)

	// Run the main counting function as part of the test
	err = processing.CountUniqueIPs(cfg, &wg, &uniqueCount, &mu, startTime, fileSize)
	if err != nil {
		t.Fatalf("CountUniqueIPs failed: %v", err)
	}

	// Wait for counting to complete
	wg.Wait()

	// Verify the final count of unique IP addresses
	expectedUniqueCount := 4
	if uniqueCount != expectedUniqueCount {
		t.Errorf("Expected %d unique IP addresses, but got %d", expectedUniqueCount, uniqueCount)
	}
}
