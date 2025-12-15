package cpu

import (
	"os/exec"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

const logtag string = "cpu"

type CpuInfo struct {
	PhysicalCores      int
	LogicalCores       int
	Percentages        []float64
	AveragePercentages float64
	Runtime            string
	User               string
	Processes          []ProcessInfo
}

type ProcessInfo struct {
	Pid          int
	CpuPercent   float64
	ChildrenPids []int32
	ParentPid    int32
	IsChild      bool
	IsRunning    bool
	ProcessName  string
	NumThreads   int32
	Username     string
}

func (c *CpuInfo) collectCores() error {
	logging.Info(logtag, "collecting cpu cores")
	physical_core_count, err := cpu.Counts(false)

	if err != nil {
		logging.Error(logtag, "error retrieving cpu cores", err)
		return err
	}

	logical_core_count, err := cpu.Counts(true)

	if err != nil {

		return err
	}

	c.PhysicalCores = physical_core_count
	c.LogicalCores = logical_core_count

	return nil

}

func (c *CpuInfo) collectPercentages(interval time.Duration) error {
	percentages, err := cpu.Percent(interval, true)

	if err != nil {
		logging.Error(logtag, "error collecting percentages", err)
		return err
	}

	c.Percentages = percentages

	var sum float64

	for _, p := range c.Percentages {
		sum += p
	}

	average := sum / float64(len(c.Percentages))

	c.AveragePercentages = average

	return nil
}

func (c *CpuInfo) getUser() error {
	var cmd *exec.Cmd
	c.getRuntime()

	switch c.Runtime {
	case "windows":
		cmd = exec.Command("cmd", "/C", "whoami")
	default:
		cmd = exec.Command("whoami")

	}
	stdout, err := cmd.Output()
	if err != nil {
		logging.Error(logtag, "error executing command:", err)
		return nil
	}

	c.User = string(stdout)
	return nil
}

func (c *CpuInfo) getRuntime() error {
	c.Runtime = runtime.GOOS
	return nil
}

func getProcesses() ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	time.Sleep(500 * time.Millisecond)

	results := make([]ProcessInfo, 0, len(procs))

	for _, proc := range procs {
		running, err := proc.IsRunning()
		if err != nil || !running {
			continue
		}

		cpuusage, err := proc.CPUPercent()
		if err != nil || cpuusage <= 0.9 {
			continue
		}

		p := ProcessInfo{
			Pid:        int(proc.Pid),
			CpuPercent: cpuusage,
			IsRunning:  true,
		}

		if threads, err := proc.NumThreads(); err == nil {
			p.NumThreads = threads
		}

		if name, err := proc.Name(); err == nil {
			p.ProcessName = name
		}

		if user, err := proc.Username(); err == nil {
			p.Username = user
		}

		if ppid, err := proc.Ppid(); err == nil {
			p.ParentPid = ppid
			p.IsChild = ppid != 0
		}

		if children, err := proc.Children(); err == nil {
			p.ChildrenPids = make([]int32, 0, len(children))
			for _, child := range children {
				p.ChildrenPids = append(p.ChildrenPids, child.Pid)
			}
		}

		results = append(results, p)
	}

	return results, nil
}
