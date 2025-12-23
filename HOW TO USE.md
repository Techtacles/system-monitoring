# How to Use System Monitoring Tool

This guide provides instructions on how to effectively use the **System Monitoring Tool** to monitor your system's performance via the web dashboard or directly through the terminal.

## Quick Start
1. **Clone and Setup:**
   ```bash
   git clone <repository-url>
   cd system-monitoring
   go mod download
   ```
2. **Run the Dashboard:**
   ```bash
   go run main.go start
   ```
3. **View in Browser:** Open [http://localhost:8080](http://localhost:8080)

---

##  टर्मिनल (CLI) Usage
You can retrieve real-time metrics directly in your terminal using the `get_metrics` command. These are displayed in a clean, human-readable table format.

### 1. Basic Commands
To get a single snapshot of specific metrics:
```bash
go run main.go get_metrics cpu
go run main.go get_metrics memory
go run main.go get_metrics disk
go run main.go get_metrics network
go run main.go get_metrics docker
go run main.go get_metrics host
go run main.go get_metrics user
```

### 2. Monitoring All Metrics
To see all available system metrics at once:
```bash
go run main.go get_metrics all
```

### 3. Real-time Auto-Refresh
You can turn the CLI into a live monitor by adding the `-a` (auto) and `-r` (refresh interval) flags:
```bash
# Refresh CPU metrics every 2 seconds
go run main.go get_metrics cpu -a -r 2

# Monitor all metrics every 10 seconds
go run main.go get_metrics all -a -r 10
```

---

## Web Dashboard Usage
The dashboard provides a visual overview of your system health with live-updating charts.

### 1. Start the Server
```bash
go run main.go start
```

### 2. Custom Port
If you need to run the dashboard on a different port:
```bash
go run main.go start -p 9090
```

### 3. Docker Telemetry
To include Docker metrics (containers, images, volumes) in your dashboard, make sure Docker is running and use the `-d` flag:
```bash
go run main.go start -d
```

---

## Metric Breakdown
| Metric | Description |
| :--- | :--- |
| **CPU** | Core counts, usage percentages, and top CPU-consuming processes. |
| **Memory** | Virtual and Swap memory usage, and processes with high memory footprint. |
| **Disk** | Disk usage per path/partition and device information. |
| **Network** | Established vs Total connections and detailed Interface I/O stats. |
| **Host** | Uptime, Kernel version, and Load Averages (1m, 5m, 15m). |
| **User** | Current logged-in user details and system architecture. |
| **Docker** | Container status, Image sizes, and Docker-specific resource usage. |

---

## Best Practices
- **Run as Admin/Sudo:** Some metrics (like certain process details or disk info) may require elevated permissions for full visibility.
- **Refresh Rates:** While the tool is lightweight, very fast refresh rates (e.g., < 1s) might increase the tool's own CPU usage slightly on very busy systems.
- **Docker Tooling:** Always ensure the Docker Daemon is running before using the `-d` or `docker` commands to avoid connection errors.
