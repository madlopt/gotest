
# IP Address Counter

A Go application that efficiently counts unique IPv4 addresses from a potentially very large file using concurrent processing and memory-efficient data structures.

## Table of Contents
- [Project Structure](#project-structure)
- [Features](#features)
- [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Testing](#testing)
- [Performance Tuning](#performance-tuning)

## Project Structure
The directory structure is organized as follows:

```
/cmd
  └── ipcounter
      └── main.go                     # Entry point of the application
/internal
  ├── config                          # Configuration management
  │   └── config.go                   # Configuration settings and constants
  ├── entities                        # Domain entities and utilities
  │   └── ip_converter.go             # IP conversion and related utilities
  ├── file                            # File handling logic
  │   └── file_reader.go              # File reading and size management
  ├── processing                      # Core processing logic
  │   └── counter.go                  # Main counting logic and workers
  ├── bitmap                          # Bitmap-specific operations
  │   └── bitmap_manager.go           # Parallel merging of bitmaps
  └── output                          # Output and reporting
      └── display.go                  # Console output and progress display
/tests                                # Integration and unit tests
  └── ipcounter_test.go               # End-to-end and unit tests for the project
```

## Features
- **Efficient Processing**: Utilizes Go routines and worker pools to process large files concurrently.
- **Roaring Bitmaps**: Uses Roaring Bitmaps to efficiently store and manipulate large sets of IP addresses.
- **Progress Monitoring**: Displays real-time progress, memory usage, and execution time.
- **Configurable**: Allows easy configuration of parameters like buffer size, worker count, and file paths.
- **Unit and Integration Tests**: Includes tests for core functionality.

## Getting Started

### Prerequisites
Ensure you have the following tools installed:
- [Go](https://golang.org/doc/install) (version 1.22 or higher)
- A compatible terminal for running the application

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/ipcounter.git
   cd ipcounter
   ```

2. Initialize the Go module:
   ```bash
   go mod tidy
   ```

## Configuration
The configuration is managed through the `internal/config/config.go` file. You can adjust the following parameters:

- **FilePath**: The path to the file containing the IP addresses.  *Default file path is the "ip_addresses" in the root of the project*.
- **BufferSize**: Buffer size used by the scanner for reading lines.
- **NumWorkers**: Number of concurrent workers to process the file.
- **PrintInterval**: Interval for printing progress updates.
- **LinesChannelCap**: Capacity of the lines channel buffer.

### Example Configuration (`config.go`):
```go
var defaultConfig = Config{
FilePath:        "ip_addresses",           // File path to process
BufferSize:      16 * 1024 * 1024,         // 16 MB buffer size
NumWorkers:      2 * runtime.NumCPU(),     // Double the number of CPU cores
PrintInterval:   10 * time.Second,         // Progress updates every 10 seconds
LinesChannelCap: 100000,                   // Channel capacity
}
```

## Running the Application
1. Build the application:
   ```bash
   go build -o ipcounter ./cmd/ipcounter
   ```

2. Run the application:
   ```bash
   ./ipcounter
   ```

   By default, the application will use the parameters defined in the configuration file. You can modify these parameters as needed.

## Testing
To run the test cases, use the following command:

```bash
go test ./tests
```

## Performance Tuning
You can optimize the performance of the application by tweaking the following parameters in `config.go`:

- **Buffer Size (`BufferSize`)**: Increase the buffer size (e.g., 32MB, 64MB) for larger files if the system has enough memory.
- **Number of Workers (`NumWorkers`)**: Adjust the number of workers based on your system's CPU cores.
- **Channel Capacity (`LinesChannelCap`)**: A larger channel capacity reduces contention but increases memory usage.