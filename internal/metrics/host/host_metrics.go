package host

import (
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

const logtag string = "host"

type HostInfo struct {
	Uptime        uint64
	BootTime      uint64
	LoadAvg1      float64
	LoadAvg5      float64
	LoadAvg15     float64
	OS            string
	Platform      string
	PlatformVer   string
	KernelVersion string
}

func (h *HostInfo) Collect() error {
	hostStat, err := host.Info()
	if err != nil {
		logging.Error(logtag, "error getting host info", err)
		return err
	}

	loadStat, err := load.Avg()
	if err != nil {
		logging.Error(logtag, "error getting load average", err)
		return err
	}

	h.Uptime = hostStat.Uptime
	h.BootTime = hostStat.BootTime
	h.OS = hostStat.OS
	h.Platform = hostStat.Platform
	h.PlatformVer = hostStat.PlatformVersion
	h.KernelVersion = hostStat.KernelVersion

	h.LoadAvg1 = loadStat.Load1
	h.LoadAvg5 = loadStat.Load5
	h.LoadAvg15 = loadStat.Load15

	logging.Info(logtag, "successfully collected host info")
	return nil
}
