# 📂 Mini Projects Collection

Welcome to this curated collection of small-scale projects and coding exercises! This repository serves as a personal playground for exploring algorithms, design patterns, concurrency, data structures, and language-specific features. It's ideal for practicing coding skills, experimenting with ideas, and building reusable components.

Currently, the focus is on Go (Golang) projects, with plans to expand to other languages like Python in the future. Each project is self-contained, well-tested, and demonstrates practical concepts in a concise manner.

## 🚀 Why This Repository?
- **Learning-Oriented**: Projects are designed to highlight key programming concepts without unnecessary complexity.
- **Modular Structure**: Easy to navigate and extend.
- **Real-World Inspired**: Implementations draw from common scenarios like batch processing, in-memory databases, and data pipelines.
- **Open for Contributions**: Feel free to suggest improvements or add your own mini-projects!

## 📁 Repository Structure
The projects are organized by programming language for clarity:

- **`go_projects/`**: Contains all Go-based implementations.
  - Each subfolder (e.g., `batcher-queue/`) includes source code, tests, and a main executable where applicable.
- **`python_projects/`**: (Coming soon) Python implementations.

## 🐹 Go Projects
These projects emphasize Go's strengths in concurrency, generics, and efficient data handling. Below is a summary table for quick reference:

| Project Name        | Description                                                                                                                        | Key Concepts Covered                                                    | Directory Path                |
| ------------------- | ---------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- | ----------------------------- |
| **Batcher Queue**   | A generic batching library that groups items based on size limits or time intervals, with support for graceful shutdown.           | Generics, Concurrency (Channels, Goroutines), Timeouts, Sync Primitives | `./go_projects/batcher-queue` |
| **DbSim (TypeBox)** | An in-memory key-value store simulator supporting scalars, lists, and nested objects, with a simple SQL-like command parser.       | In-Memory Storage, Interfaces, Parsing, Recursion                       | `./go_projects/db-sim`        |
| **Async Pipeline**  | A multi-stage asynchronous data processing pipeline (inspired by Unix pipes) for tasks like user data handling and spam detection. | Channels, Worker Pools, Fan-Out/Fan-In Patterns, Synchronization        | `./go_projects/pipeline`      |

### Detailed Project Overviews

#### Batcher Queue
- **Purpose**: Efficiently batches events or items to reduce overhead in scenarios like logging, API calls, or database writes.
- **Features**:
  - Configurable batch size and flush interval.
  - Thread-safe with mutexes and channels.
  - Handles partial batches on timeout or shutdown.
- **Example Usage** (from `main.go`):
  ```go
  package main

  import (
      "fmt"
      "time"
      "github.com/moguchev/stepik/4/4.6/HW/batcher" // Adjust import path as needed
  )

  func main() {
      handler := func(batch []string) {
          fmt.Println(">>> Flushed:", batch)
          time.Sleep(500 * time.Millisecond)
      }
      b := batcher.NewBatcher(5, 2*time.Second, handler)
      for i := 1; i <= 12; i++ {
          b.Add(fmt.Sprintf("event-%d", i))
          time.Sleep(300 * time.Millisecond)
      }
      b.Close()
      fmt.Println("Batcher gracefully shut down.")
  }
  ```
- **Testing**: Comprehensive unit tests in `batcher_queue_test.go` cover flushing by size, timeout, shutdown, and edge cases.

#### DbSim (TypeBox)
- **Purpose**: Simulates a lightweight database for storing and querying typed data, useful for prototyping or educational purposes.
- **Features**:
  - Supports scalar values, lists, and nested objects.
  - Command-based interface (e.g., SET, OBJECT) with parsing.
  - In-memory storage for fast operations.
- **Example Usage** (from `main.go`):
  ```go
  package main

  import (
      "bufio"
      "os"
      "strconv"
      "strings"
      "github.com/dim4d/DbSim/core"
      "github.com/dim4d/DbSim/storage"
  )

  func main() {
      // Simplified scanner and command processing
      tb := storage.NewTypeBox()
      // Process commands like "SET key type value" or "OBJECT key n" with fields
  }
  ```
- **Note**: The full implementation includes recursive parsing for nested structures.

#### Async Pipeline
- **Purpose**: Demonstrates building scalable, concurrent data pipelines for processing streams of data.
- **Features**:
  - Multi-stage processing with fan-out/fan-in for parallelism.
  - Error handling and synchronization using Go's concurrency primitives.
- **Usage**: Navigate to the directory and run with `go run .` or build the executable.

## 🏃 Quick Start
To get started with any Go project:

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/mini-projects.git
   cd mini-projects
   ```

2. **Navigate to a Project**:
   ```bash
   cd go_projects/batcher-queue/app
   ```

3. **Run Tests**:
   ```bash
   go test -v ./...
   ```

4. **Build and Run**:
   ```bash
   go build
   ./app  # Or the generated executable
   ```

Ensure you have Go installed (version 1.18+ for generics support). No external dependencies are required beyond the standard library.

## 🔧 Contributing
Contributions are welcome! If you'd like to add a new project, fix a bug, or improve documentation:
- Fork the repo.
- Create a feature branch.
- Submit a pull request with a clear description.

Please follow Go's coding conventions and include tests for new features.

---
