# System Monitoring Tool

A lightweight, Go-based system monitoring agent designed to collect and analyze system performance metrics including CPU, memory, disk, network, and OS information.

## Overview

This project provides a modular framework for gathering system telemetry. It leverages `gopsutil` for cross-platform metric collection and includes a structured logging system. The agent is designed to be extensible, with a planned dashboard for real-time visualization.

## Features

- **CPU Monitoring**:
  - Collects physical and logical core counts.
  - Monitors CPU usage percentages (per core and average).
  - Identifies resource-intensive processes (CPU usage > 0.9%).
  - Tracks process details: PID, name, user, thread count, parent/child relationships.
  
- **Modular Architecture**:
  - Separation of concerns with dedicated packages for `cpu`, `memory`, `disk`, `network`, `os`, and `trace`.
  - Built-in `aggregator` for metric consolidation.
  - Structured logging via `zerolog`.

## Project Structure

```
├── cmd/                # Application entry points
│   ├── cli.go          # CLI root command definition
│   └── run.go          # 'start' command implementation
├── internal/           # Private application layout code
│   ├── dashboard/      # Web dashboard implementation and embedded assets
│   ├── logging/        # Logging configuration and helpers
│   └── metrics/        # Metric collection modules
│       ├── aggregator/ # Aggregates collected metrics
│       ├── cpu/        # CPU stats and process monitoring
│       ├── disk/       # Disk usage and IO stats
│       ├── memory/     # RAM and Swap usage
│       ├── network/    # Network interface stats
│       ├── host/       # Host stats
│       └── user/       # User stats
├── main.go             # Main application entry point
├── experiment.go       # Experimental code snippets
└── go.mod              # Go module definition
```

## Prerequisites

- **Go**: Version 1.24.0 or higher
- **Dependencies**: Managed via `go.mod`
  - `github.com/shirou/gopsutil/v4`
  - `github.com/rs/zerolog`
  - `github.com/spf13/cobra`

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd system-monitoring
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

## Running the Application

You can run the application directly using `go run` or by building a binary.

### Using `go run`

```bash
# Start the dashboard on the default port (8080)
go run main.go start

# Start on a specific port
go run main.go start --port 9090
```

### Building the Binary

```bash
# Build the binary
go build -o sysmon main.go

# Run the binary
./sysmon start

# Run with custom port
./sysmon start --port 5555
```


If your OS is Apple, and you download the binary from the latest artifact, it might sometimes get flagged. To resolve this, run the command: 
```bash
chmod +x sysmon-darwin
xattr -d com.apple.quarantine sysmon-darwin
```
Once running, open your browser and navigate to `http://localhost:8080` (or your custom port) to view the dashboard.


## Future enhancements

- Add integrations such as prometheus and grafana
- Add more features to the dashboard
- Add alerting
- Connect to slack/teams/webhook for notifications
- Add time series analysis
