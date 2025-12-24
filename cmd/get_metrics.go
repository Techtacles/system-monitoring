package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"strings"

	"github.com/spf13/cobra"
	"github.com/techtacles/sysmonitoring/internal/logging"
	"github.com/techtacles/sysmonitoring/internal/metrics/aggregator"
	"github.com/techtacles/sysmonitoring/internal/metrics/cpu"
	"github.com/techtacles/sysmonitoring/internal/metrics/disk"
	"github.com/techtacles/sysmonitoring/internal/metrics/docker"
	"github.com/techtacles/sysmonitoring/internal/metrics/host"
	"github.com/techtacles/sysmonitoring/internal/metrics/kubernetes"
	"github.com/techtacles/sysmonitoring/internal/metrics/memory"
	"github.com/techtacles/sysmonitoring/internal/metrics/network"
	"github.com/techtacles/sysmonitoring/internal/metrics/user"
)

const metricsLogTag = "get_metrics"

var collectAutoRefresh bool
var refreshInterval int
var getKubeconfigPath string

var GetMetricCmd = &cobra.Command{
	Use:   "get_metrics",
	Short: "Get a particular metric",
	Long:  `Get a particular metric. Can take in args like cpu, disk, host, memory, network, user, docker, all, kubernetes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			logging.Info(metricsLogTag, "No metrics passed in. Please ensure metrics is either: cpu, disk, host, memory, network, user, docker, all, kubernetes")
			return nil
		}

		newAgg := aggregator.NewAggregator(false, false, getKubeconfigPath)
		metricList := args
		if len(args) == 1 && args[0] == "all" {
			metricList = []string{"cpu", "memory", "disk", "host", "network", "user", "docker", "kubernetes"}
		}

		collectAndPrint := func() {
			for _, metric := range metricList {
				var err error
				switch metric {
				case "cpu":
					err = newAgg.CollectCPU()
				case "memory":
					err = newAgg.CollectMemory()
				case "disk":
					err = newAgg.CollectDisk()
				case "host":
					err = newAgg.CollectHost()
				case "network":
					err = newAgg.CollectNetwork()
				case "user":
					err = newAgg.CollectUser()
				case "docker":
					err = newAgg.CollectDocker()
				case "kubernetes":
					err = newAgg.CollectKubernetes()
				}

				if err != nil {
					logging.Error(metricsLogTag, fmt.Sprintf("failed to collect %s metrics", metric), err)
					continue
				}

				if val, ok := newAgg.GetMetric(metric); ok {
					printMetric(metric, val)
				}
			}
		}

		collectAndPrint()

		if collectAutoRefresh {
			ticker := time.NewTicker(time.Duration(refreshInterval) * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				logging.Info(metricsLogTag, fmt.Sprintf("refreshing every %d seconds", refreshInterval))
				collectAndPrint()
			}
		}

		return nil
	},
}

func printMetric(name string, result interface{}) {
	fmt.Printf("\n--- %s Metrics ---\n", strings.ToUpper(name))
	switch name {
	case "cpu":
		if info, ok := result.(cpu.CpuInfo); ok {
			printCPUTable(info)
		} else {
			fmt.Printf("%+v\n", result)
		}
	case "memory":
		if info, ok := result.(memory.MemoryInfo); ok {
			printMemoryTable(info)
		} else {
			fmt.Printf("%+v\n", result)
		}
	case "disk":
		if info, ok := result.(disk.DiskInfo); ok {
			printDiskTable(info)
		} else {
			fmt.Printf("%+v\n", result)
		}
	case "network":
		if info, ok := result.(network.NetworkInfo); ok {
			printNetworkTable(info)
		} else {
			fmt.Printf("%+v\n", result)
		}
	case "host":
		if info, ok := result.(host.HostInfo); ok {
			printHostTable(info)
		} else {
			fmt.Printf("%+v\n", result)
		}
	case "user":
		if info, ok := result.(user.UserInfo); ok {
			printUserTable(info)
		} else {
			fmt.Printf("%+v\n", result)
		}
	case "docker":
		if info, ok := result.(docker.DockerInfo); ok {
			printDockerTable(info)
		} else {
			fmt.Printf("%+v\n", result)
		}
	case "kubernetes":
		if info, ok := result.(kubernetes.KubeInfo); ok {
			printKubernetesTable(info)
		} else {
			fmt.Printf("%+v\n", result)
		}
	default:
		fmt.Printf("%+v\n", result)
	}
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func printCPUTable(info cpu.CpuInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "Physical Cores:\t%d\n", info.PhysicalCores)
	fmt.Fprintf(w, "Logical Cores:\t%d\n", info.LogicalCores)
	fmt.Fprintf(w, "Average Usage:\t%.2f%%\n", info.AveragePercentages)
	w.Flush()

	fmt.Println("\nCore Percentages:")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	header := ""
	for i := range info.Percentages {
		header += fmt.Sprintf("Core %d\t", i)
	}
	fmt.Fprintln(w, header)

	row := ""
	for _, p := range info.Percentages {
		row += fmt.Sprintf("%.2f%%\t", p)
	}
	fmt.Fprintln(w, row)
	w.Flush()

	if len(info.Processes) > 0 {
		fmt.Println("\nTop Processes (CPU):")
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "PID\tName\tCPU%\tThreads\tUser")
		for _, p := range info.Processes {
			fmt.Fprintf(w, "%d\t%s\t%.2f%%\t%d\t%s\n", p.Pid, p.ProcessName, p.CpuPercent, p.NumThreads, p.Username)
		}
		w.Flush()
	}
}

func printKubernetesTable(info kubernetes.KubeInfo) {
	fmt.Printf("Summary: %d Nodes, %d Pods, %d Services, %d Deployments\n",
		len(info.NodeStats), len(info.PodStats), len(info.ServiceStats), len(info.DeploymentStats))

	if len(info.NodeStats) > 0 {
		fmt.Println("\nNodes:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Name\tStatus\tAge")
		for _, n := range info.NodeStats {
			status := "Ready"
			if n.Unschedulable {
				status = "Unschedulable"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", n.Name, status, time.Since(n.CreationTimestamp).Round(time.Second))
		}
		w.Flush()
	}

	if len(info.PodStats) > 0 {
		fmt.Println("\nPods (Top 10):")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Namespace\tName\tPhase\tIP\tNode")
		displayPods := info.PodStats
		if len(displayPods) > 10 {
			displayPods = displayPods[:10]
		}
		for _, p := range displayPods {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", p.Namespace, p.Name, p.Phase, p.PodIP, p.NodeName)
		}
		w.Flush()
	}

	if len(info.PersistentVolumeStats) > 0 || len(info.PersistentVolumeClaimStats) > 0 {
		fmt.Printf("\nStorage: %d PVs, %d PVCs\n", len(info.PersistentVolumeStats), len(info.PersistentVolumeClaimStats))
	}
}

func printMemoryTable(info memory.MemoryInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Resource\tTotal\tAvailable/Free\tUsed\tUsed %")
	fmt.Fprintf(w, "Virtual\t%s\t%s\t%s\t%.2f%%\n",
		formatBytes(info.Vmemory.Total),
		formatBytes(info.Vmemory.Available),
		formatBytes(info.Vmemory.Used),
		info.Vmemory.UsedPercentage)
	fmt.Fprintf(w, "Swap\t%s\t%s\t%s\t%.2f%%\n",
		formatBytes(info.SwapMemoryTotal),
		formatBytes(info.SwapMemoryFree),
		formatBytes(info.SwapMemoryUsed),
		info.SwapMemoryUsedPercent)
	w.Flush()

	if len(info.ProcessInfo) > 0 {
		fmt.Println("\nTop Processes (Memory):")
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "PID\tName\tMem%\tRSS\tVMS\tUser")
		for _, p := range info.ProcessInfo {
			fmt.Fprintf(w, "%d\t%s\t%.2f%%\t%s\t%s\t%s\n",
				p.Pid, p.ProcessName, p.MemPercent,
				formatBytes(p.PhysicalMemorySize),
				formatBytes(p.VirtualMemorySize),
				p.Username)
		}
		w.Flush()
	}
}

func printDiskTable(info disk.DiskInfo) {
	fmt.Println("Disk Usage:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Path\tTotal\tUsed\tFree\tUsed %")
	for _, u := range info.UsageStat {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%.2f%%\n",
			u.Path,
			formatBytes(u.TotalDisk),
			formatBytes(u.UsedDisk),
			formatBytes(u.FreeDisk),
			u.UsedPercent)
	}
	w.Flush()

	fmt.Println("\nPartitions:")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Device\tMount\tType\tOpts")
	for _, p := range info.PartitionInfo {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Device, p.MountPoint, p.Fstype, strings.Join(p.Opts, ","))
	}
	w.Flush()
}

func printNetworkTable(info network.NetworkInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "OS:\t%s\n", info.Runtime)
	fmt.Fprintf(w, "Established Conns:\t%d\n", info.NumEstablishedConnections)
	fmt.Fprintf(w, "Total Conns:\t%d\n", info.NumTotalConnections)
	w.Flush()

	fmt.Println("\nInterface IO Stats:")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Name\tSent\tRecv\tErr In/Out\tDrop In/Out")
	for _, io := range info.IOStats {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d/%d\t%d/%d\n",
			io.Name,
			formatBytes(io.BytesSent),
			formatBytes(io.BytesRecv),
			io.Errin, io.Errout,
			io.Dropin, io.Dropout)
	}
	w.Flush()
}

func printHostTable(info host.HostInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "OS:\t%s\n", info.OS)
	fmt.Fprintf(w, "Platform:\t%s (%s)\n", info.Platform, info.PlatformVer)
	fmt.Fprintf(w, "Kernel:\t%s\n", info.KernelVersion)
	fmt.Fprintf(w, "Uptime:\t%d s\n", info.Uptime)
	fmt.Fprintf(w, "Load Avg:\t%.2f, %.2f, %.2f (1m, 5m, 15m)\n", info.LoadAvg1, info.LoadAvg5, info.LoadAvg15)
	w.Flush()
}

func printUserTable(info user.UserInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "Username:\t%s\n", info.Username)
	fmt.Fprintf(w, "Full Name:\t%s\n", info.FullName)
	fmt.Fprintf(w, "Home Dir:\t%s\n", info.HomeDir)
	fmt.Fprintf(w, "Config:\t%s (%s)\n", info.Runtime, info.Arch)
	w.Flush()
}

func printDockerTable(info docker.DockerInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "Docker Env:\t%s\n", info.DockerEnv)
	fmt.Fprintf(w, "Containers:\t%d running, %d paused, %d stopped\n",
		info.ContainersRunning, info.ContainersPaused, info.ContainersStopped)
	fmt.Fprintf(w, "Resources:\t%d CPUs, %s Total Mem\n", info.NCpu, formatBytes(uint64(info.MemTotal)))
	fmt.Fprintf(w, "Storage Usage:\tContainers: %s, Images: %s, Cache: %s\n",
		formatBytes(uint64(info.ContainersDiskUsage)),
		formatBytes(uint64(info.ImagesDiskUsage)),
		formatBytes(uint64(info.BuildCacheDiskUsage)))
	w.Flush()

	if len(info.ContainerStats) > 0 {
		fmt.Println("\nContainers:")
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "ID\tImage\tNames\tState")
		for _, c := range info.ContainerStats {
			id := c.ID
			if len(id) > 12 {
				id = id[:12]
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n",
				id, c.ImageName, strings.Join(c.ContainerNames, ","), c.ContainerState)
		}
		w.Flush()
	}

	if len(info.ImageStats) > 0 {
		fmt.Println("\nImages:")
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "ID\tTags\tSize\tCreated")
		for _, img := range info.ImageStats {
			id := img.ID
			if strings.HasPrefix(id, "sha256:") {
				id = id[7:19]
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				id, strings.Join(img.ImageNames, ","), formatBytes(uint64(img.ImageSize)), img.CreatedDate)
		}
		w.Flush()
	}
}

func init() {
	rootCmd.AddCommand(GetMetricCmd)
	GetMetricCmd.Flags().BoolVarP(&collectAutoRefresh, "auto", "a", false, "Whether to autorefresh every 30 seconds")
	GetMetricCmd.Flags().IntVarP(&refreshInterval, "refresh", "r", 30, "Number of seconds to autorefresh")
	GetMetricCmd.Flags().StringVarP(&getKubeconfigPath, "kubeconfig", "", "", "absolute path to the kubeconfig file (optional)")
	GetMetricCmd.Flags().BoolVarP(&collectDocker, "docker", "d", false, "Whether to collect docker metrics. Make sure docker is running when passing this flag")
}
