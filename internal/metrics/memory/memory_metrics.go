package memory

import (
	"strings"

	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

const logtag string = "memory"

type MemoryInfo struct {
	Vmemory               VirtualMemoryInfo
	SwapMemoryTotal       uint64
	SwapMemoryUsed        uint64
	SwapMemoryFree        uint64
	SwapMemoryUsedPercent float64
	ProcessInfo           []ProcessInfo
}

type ProcessInfo struct {
	Pid                   int
	MemPercent            float32
	ChildrenPids          []int32
	ParentPid             int32
	IsChild               bool
	ProcessName           string
	NumThreads            int32
	Username              string
	VirtualMemorySize     uint64
	PhysicalMemorySize    uint64
	MaxPhysicalMemorySize uint64
	MemoryUsedByHeap      uint64
	MemoryUsedByStack     uint64
	LockedMemory          uint64
}

type VirtualMemoryInfo struct {
	Total          uint64
	Available      uint64
	Used           uint64
	UsedPercentage float64
	Free           uint64
	Shared         uint64
	SReclaimable   uint64
	SUnreclaimable uint64
	Active         uint64
	Inactive       uint64
}

func (m *MemoryInfo) Collect() error {
	processes, err := getProcesses()

	if err != nil {
		logging.Error(logtag, "error instantiating memory processes", err)
		return err
	}

	m.ProcessInfo = processes

	m.getSwapMemoryInfo()

	vm, err := getVirtualMemoryInfo()
	if err != nil {
		logging.Error(logtag, "error instantiating virtual mem", err)
		return err
	}

	m.Vmemory = vm

	return nil
}

func (m *MemoryInfo) getSwapMemoryInfo() error {
	swap_memory_info, err := mem.SwapMemory()
	if err != nil {
		logging.Error(logtag, "unable to get swap memory info", err)
		return err
	}

	m.SwapMemoryTotal = swap_memory_info.Total
	m.SwapMemoryUsed = swap_memory_info.Used
	m.SwapMemoryFree = swap_memory_info.Free
	m.SwapMemoryUsedPercent = swap_memory_info.UsedPercent
	logging.Info(logtag, "successfully instantiated swap memory")
	return nil

}

func getVirtualMemoryInfo() (VirtualMemoryInfo, error) {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return VirtualMemoryInfo{}, err
	}

	return VirtualMemoryInfo{
		Total:          vmem.Total,
		Available:      vmem.Available,
		Used:           vmem.Used,
		UsedPercentage: vmem.UsedPercent,
		Free:           vmem.Free,
		Shared:         vmem.Shared,
		SReclaimable:   vmem.Sreclaimable,
		SUnreclaimable: vmem.Sunreclaim,
		Active:         vmem.Active,
		Inactive:       vmem.Inactive,
	}, nil
}

func getProcesses() ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		logging.Error(logtag, "error retrieving processes", err)
		return nil, err
	}

	results := make([]ProcessInfo, 0, len(procs))

	for _, proc := range procs {
		username, err := proc.Username()

		if err != nil {
			continue
		}

		mempercent, err := proc.MemoryPercent()
		if err != nil {
			continue
		}

		if username == "root" || strings.HasPrefix(username, "_") {
			continue
		}

		// we dont want to add noise by getting memory <1
		if mempercent <= 1 {
			continue
		}

		ppid, err := proc.Ppid()
		if err != nil || ppid == 0 {
			continue
		}

		meminfo, err := proc.MemoryInfo()
		if err != nil {
			continue
		}
		p := ProcessInfo{
			Pid:                   int(proc.Pid),
			ParentPid:             ppid,
			IsChild:               true,
			MemPercent:            mempercent,
			Username:              username,
			VirtualMemorySize:     meminfo.VMS,
			PhysicalMemorySize:    meminfo.RSS,
			MaxPhysicalMemorySize: meminfo.HWM,
			MemoryUsedByHeap:      meminfo.Data,
			MemoryUsedByStack:     meminfo.Stack,
			LockedMemory:          meminfo.Locked,
		}

		p.ProcessName, _ = proc.Name()
		p.NumThreads, _ = proc.NumThreads()
		if children, err := proc.Children(); err == nil {
			p.ChildrenPids = make([]int32, 0, len(children))
			for _, c := range children {
				p.ChildrenPids = append(p.ChildrenPids, c.Pid)
			}
		}

		results = append(results, p)

	}
	logging.Info(logtag, "finished instantiating GetProcesses()")
	return results, nil

}
