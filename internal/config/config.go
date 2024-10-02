package config

import (
	"runtime"
	"time"
)

// Config holds all the configurable parameters
type Config struct {
	FilePath        string        // Path to the file containing IP addresses
	BufferSize      int           // Scanner buffer size
	NumWorkers      int           // Number of workers
	PrintInterval   time.Duration // Interval for printing progress
	LinesChannelCap int           // Capacity of the lines channel buffer
}

// LoadConfig returns the default configuration values
func LoadConfig() Config {
	return Config{
		FilePath:        "ip_addresses_big",
		BufferSize:      8 * 1024 * 1024, // 8 MB buffer size
		NumWorkers:      2 * runtime.NumCPU(),
		PrintInterval:   10 * time.Second,
		LinesChannelCap: 100000,
	}
}
