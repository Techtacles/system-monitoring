package dashboard

import (
	"bytes"
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/techtacles/sysmonitoring/internal/logging"
	"github.com/techtacles/sysmonitoring/internal/metrics/aggregator"
	"github.com/techtacles/sysmonitoring/internal/metrics/cpu"
	"github.com/techtacles/sysmonitoring/internal/metrics/disk"
	"github.com/techtacles/sysmonitoring/internal/metrics/memory"
)

const (
	logtag          = "dashboard"
	refreshInterval = 30 * time.Second // Increased frequency for "live" feel in background
)

var Port string = "8080"

//go:embed images
var imagesDir embed.FS

// Run starts the dashboard server
func Run(enableDocker, enableKubernetes bool, kubeconfigPath string) error {
	ag := aggregator.NewAggregator(enableDocker, enableKubernetes, kubeconfigPath)

	logging.Info(logtag, "performing initial metrics collection")
	ag.CollectAllConcurrent()

	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()
		for range ticker.C {
			logging.Info(logtag, "performing scheduled metrics collection")
			ag.CollectAllConcurrent()
		}
	}()

	// Serve static images
	http.Handle("/images/", http.FileServer(http.FS(imagesDir)))

	// Serve static web assets (CSS, JS)
	http.Handle("/web/", http.FileServer(http.FS(WebAssets)))

	// API Endpoint for raw metrics
	http.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metrics := ag.GetMetrics()
		if err := json.NewEncoder(w).Encode(metrics); err != nil {
			logging.Error(logtag, "error encoding metrics to json", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})

	// API Endpoint for multi-format report export
	http.HandleFunc("/api/report", func(w http.ResponseWriter, r *http.Request) {
		format := r.URL.Query().Get("format")
		if format == "" {
			format = "json"
		}

		metrics := ag.GetMetrics()

		switch format {
		case "csv":
			generateCSVReport(w, metrics)
		case "pdf":
			generatePDFReport(w, metrics)
		default: // default to json
			generateJSONReport(w, metrics)
		}
	})

	// Dashboard UI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFS(WebAssets, "web/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			logging.Error(logtag, "error parsing template", err)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			logging.Error(logtag, "error executing template", err)
		}
	})

	logging.Info(logtag, fmt.Sprintf("starting dashboard server on http://localhost:%s", Port))
	return http.ListenAndServe(":"+Port, nil)
}

func generateJSONReport(w http.ResponseWriter, metrics map[string]interface{}) {
	metricsJSON, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		logging.Error(logtag, "error marshaling metrics for report", err)
		http.Error(w, "error generating report", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("sysmon-report-%s.json", time.Now().Format("2006-01-02-150405"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Write(metricsJSON)
}

func generateCSVReport(w http.ResponseWriter, metrics map[string]interface{}) {
	w.Header().Set("Content-Type", "text/csv")
	filename := fmt.Sprintf("sysmon-report-%s.csv", time.Now().Format("2006-01-02-150405"))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write System Summary
	writer.Write([]string{"Section", "Metric", "Value"})

	// CPU Summary
	if cpuRaw, ok := metrics["cpu"]; ok {
		if c, ok := cpuRaw.(cpu.CpuInfo); ok {
			writer.Write([]string{"CPU", "Average Load", fmt.Sprintf("%.2f%%", c.AveragePercentages)})
			writer.Write([]string{"CPU", "Physical Cores", fmt.Sprintf("%d", c.PhysicalCores)})
			writer.Write([]string{"CPU", "Logical Cores", fmt.Sprintf("%d", c.LogicalCores)})
		}
	}

	// Memory Summary
	if memRaw, ok := metrics["memory"]; ok {
		if m, ok := memRaw.(memory.MemoryInfo); ok {
			writer.Write([]string{"Memory", "Used Percentage", fmt.Sprintf("%.2f%%", m.Vmemory.UsedPercentage)})
			writer.Write([]string{"Memory", "Total", fmt.Sprintf("%.2f GB", float64(m.Vmemory.Total)/1024/1024/1024)})
			writer.Write([]string{"Memory", "Used", fmt.Sprintf("%.2f GB", float64(m.Vmemory.Used)/1024/1024/1024)})
		}
	}

	// Disk Summary
	if diskRaw, ok := metrics["disk"]; ok {
		if d, ok := diskRaw.(disk.DiskInfo); ok {
			for path, usage := range d.UsageStat {
				usedGB := float64(usage.UsedDisk) / 1024 / 1024 / 1024
				totalGB := float64(usage.TotalDisk) / 1024 / 1024 / 1024
				writer.Write([]string{"Disk", path, fmt.Sprintf("%.2f GB / %.2f GB (%.1f%%)", usedGB, totalGB, usage.UsedPercent)})
			}
		}
	}
}

func generatePDFReport(w http.ResponseWriter, metrics map[string]interface{}) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "System Monitoring Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, fmt.Sprintf("Generated on: %s", time.Now().Format(time.RFC1123)))
	pdf.Ln(15)

	// CPU Section
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(190, 8, "CPU Metrics", "1", 0, "L", true, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	if cpuRaw, ok := metrics["cpu"]; ok {
		if c, ok := cpuRaw.(cpu.CpuInfo); ok {
			pdf.Cell(95, 8, fmt.Sprintf("Average Load: %.2f%%", c.AveragePercentages))
			pdf.Ln(6)
			pdf.Cell(95, 8, fmt.Sprintf("Physical Cores: %d", c.PhysicalCores))
			pdf.Ln(6)
			pdf.Cell(95, 8, fmt.Sprintf("Logical Cores: %d", c.LogicalCores))
			pdf.Ln(10)
		}
	}

	// Memory Section
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 8, "Memory Metrics", "1", 0, "L", true, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	if memRaw, ok := metrics["memory"]; ok {
		if m, ok := memRaw.(memory.MemoryInfo); ok {
			usedGB := float64(m.Vmemory.Used) / 1024 / 1024 / 1024
			totalGB := float64(m.Vmemory.Total) / 1024 / 1024 / 1024
			pdf.Cell(95, 8, fmt.Sprintf("Usage: %.2f%% (%.2f GB / %.2f GB)", m.Vmemory.UsedPercentage, usedGB, totalGB))
			pdf.Ln(10)
		}
	}

	// Disk Section
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 8, "Disk Metrics", "1", 0, "L", true, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	if diskRaw, ok := metrics["disk"]; ok {
		if d, ok := diskRaw.(disk.DiskInfo); ok {
			// Sort names for predictable order
			var paths []string
			for p := range d.UsageStat {
				paths = append(paths, p)
			}
			sort.Strings(paths)

			for _, path := range paths {
				usage := d.UsageStat[path]
				usedGB := float64(usage.UsedDisk) / 1024 / 1024 / 1024
				totalGB := float64(usage.TotalDisk) / 1024 / 1024 / 1024
				pdf.Cell(190, 8, fmt.Sprintf("%s: %.2f GB / %.2f GB (%.1f%%)", path, usedGB, totalGB, usage.UsedPercent))
				pdf.Ln(6)
			}
		}
	}

	w.Header().Set("Content-Type", "application/pdf")
	filename := fmt.Sprintf("sysmon-report-%s.pdf", time.Now().Format("2006-01-02-150405"))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		logging.Error(logtag, "error generating pdf", err)
		http.Error(w, "error generating pdf", http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}
