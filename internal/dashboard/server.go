package dashboard

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/techtacles/sysmonitoring/internal/logging"
	"github.com/techtacles/sysmonitoring/internal/metrics/aggregator"
)

const (
	logtag          = "dashboard"
	refreshInterval = 30 * time.Second // Increased frequency for "live" feel in background
)

var Port string = "8080"

// Run starts the dashboard server
func Run() error {
	ag := aggregator.NewAggregator()

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

	// API Endpoint for raw metrics
	http.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metrics := ag.GetMetrics()
		if err := json.NewEncoder(w).Encode(metrics); err != nil {
			logging.Error(logtag, "error encoding metrics to json", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})

	// Dashboard UI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("dashboard").Parse(tmpl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			logging.Error(logtag, "error parsing template", err)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			logging.Error(logtag, "error executing template", err)
		}
	})

	logging.Info(logtag, fmt.Sprintf("starting dashboard server on %s", Port))
	return http.ListenAndServe(":"+Port, nil)
}
