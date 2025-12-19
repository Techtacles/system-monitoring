package aggregator

import (
	"sync"

	"github.com/techtacles/sysmonitoring/internal/logging"
	"github.com/techtacles/sysmonitoring/internal/metrics/cpu"
	"github.com/techtacles/sysmonitoring/internal/metrics/disk"
	"github.com/techtacles/sysmonitoring/internal/metrics/host"
	"github.com/techtacles/sysmonitoring/internal/metrics/memory"
	"github.com/techtacles/sysmonitoring/internal/metrics/network"
	"github.com/techtacles/sysmonitoring/internal/metrics/user"
)

const logtag = "aggregator"

type Aggregator struct {
	mu         sync.RWMutex
	allMetrics map[string]interface{}
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		allMetrics: make(map[string]interface{}),
	}
}

// CollectAll collects all metrics sequentially and returns any errors encountered
func (a *Aggregator) CollectAll() map[string]error {
	errors := make(map[string]error)

	if err := a.CollectCPU(); err != nil {
		errors["cpu"] = err
	}
	if err := a.CollectMemory(); err != nil {
		errors["memory"] = err
	}
	if err := a.CollectDisk(); err != nil {
		errors["disk"] = err
	}
	if err := a.CollectNetwork(); err != nil {
		errors["network"] = err
	}
	if err := a.CollectUser(); err != nil {
		errors["user"] = err
	}
	if err := a.CollectHost(); err != nil {
		errors["host"] = err
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// CollectAllConcurrent collects all metrics concurrently for better performance
func (a *Aggregator) CollectAllConcurrent() map[string]error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make(map[string]error)

	collectors := []struct {
		name string
		fn   func() error
	}{
		{"cpu", a.CollectCPU},
		{"memory", a.CollectMemory},
		{"disk", a.CollectDisk},
		{"network", a.CollectNetwork},
		{"user", a.CollectUser},
		{"host", a.CollectHost},
	}

	for _, collector := range collectors {
		wg.Add(1)
		go func(name string, fn func() error) {
			defer wg.Done()
			if err := fn(); err != nil {
				mu.Lock()
				errors[name] = err
				mu.Unlock()
			}
		}(collector.name, collector.fn)
	}

	wg.Wait()

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (a *Aggregator) CollectCPU() error {
	logging.Info(logtag, "collecting cpu metrics")

	c := cpu.CpuInfo{}
	if err := c.Collect(); err != nil {
		logging.Error(logtag, "error collecting cpu metrics", err)
		return err
	}

	a.mu.Lock()
	a.allMetrics["cpu"] = c
	a.mu.Unlock()

	logging.Info(logtag, "successfully collected cpu metrics")
	return nil
}

func (a *Aggregator) CollectMemory() error {
	logging.Info(logtag, "collecting memory metrics")

	m := memory.MemoryInfo{}
	if err := m.Collect(); err != nil {
		logging.Error(logtag, "error collecting memory metrics", err)
		return err
	}

	a.mu.Lock()
	a.allMetrics["memory"] = m
	a.mu.Unlock()

	logging.Info(logtag, "successfully collected memory metrics")
	return nil
}

func (a *Aggregator) CollectDisk() error {
	logging.Info(logtag, "collecting disk metrics")

	d := disk.DiskInfo{}
	if err := d.Collect(); err != nil {
		logging.Error(logtag, "error collecting disk metrics", err)
		return err
	}

	a.mu.Lock()
	a.allMetrics["disk"] = d
	a.mu.Unlock()

	logging.Info(logtag, "successfully collected disk metrics")
	return nil
}

func (a *Aggregator) CollectNetwork() error {
	logging.Info(logtag, "collecting network metrics")

	n := network.NetworkInfo{}
	if err := n.Collect(); err != nil {
		logging.Error(logtag, "error collecting network metrics", err)
		return err
	}

	a.mu.Lock()
	a.allMetrics["network"] = n
	a.mu.Unlock()

	logging.Info(logtag, "successfully collected network metrics")
	return nil
}

func (a *Aggregator) CollectUser() error {
	logging.Info(logtag, "collecting user metrics")

	u := user.UserInfo{}
	if err := u.Collect(); err != nil {
		logging.Error(logtag, "error collecting user metrics", err)
		return err
	}

	a.mu.Lock()
	a.allMetrics["user"] = u
	a.mu.Unlock()

	logging.Info(logtag, "successfully collected user metrics")
	return nil
}

func (a *Aggregator) CollectHost() error {
	logging.Info(logtag, "collecting host metrics")

	h := host.HostInfo{}
	if err := h.Collect(); err != nil {
		logging.Error(logtag, "error collecting host metrics", err)
		return err
	}

	a.mu.Lock()
	a.allMetrics["host"] = h
	a.mu.Unlock()

	logging.Info(logtag, "successfully collected host metrics")
	return nil
}

// GetMetrics returns a copy of all collected metrics (thread-safe)
func (a *Aggregator) GetMetrics() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	metrics := make(map[string]interface{}, len(a.allMetrics))
	for k, v := range a.allMetrics {
		metrics[k] = v
	}
	return metrics
}

// GetMetric returns a specific metric by name (thread-safe)
func (a *Aggregator) GetMetric(name string) (interface{}, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	metric, exists := a.allMetrics[name]
	return metric, exists
}

func (a *Aggregator) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.allMetrics = make(map[string]interface{})
	logging.Info(logtag, "cleared all metrics")
}
