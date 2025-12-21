package docker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

const logtag = "docker"

var prevCPUStats = make(map[string]container.StatsResponse)

type DockerInfo struct {
	DockerEnv                    string // eg docker-desktop
	ContainersRunning            int
	ContainersPaused             int
	ContainersStopped            int
	NCpu                         int
	MemTotal                     int64
	ContainersDiskUsage          int64
	ImagesDiskUsage              int64
	BuildCacheDiskUsage          int64
	PlatformName                 string
	APIVersion                   string
	OS                           string
	Arch                         string
	TotalImages                  int
	TotalContainers              int
	TotalVolumes                 int
	ContainerCpuMemoryCollection map[string]ContainerMetrics
	ContainerStats               []Containers
	ImageStats                   []Images
	VolumeStats                  []Volumes
}

// example output: {nginx [/sleepy_chaum] [{invalid IP 80 0 tcp}] running 0}
type Containers struct {
	ID                       string
	ImageName                string
	ContainerNames           []string
	ContainerPorts           []container.PortSummary
	ContainerState           container.ContainerState
	ContainerRootSizeInBytes int64
}

// example: all map[51ac:map[cpu:0 disk:[] memory:0 name:/elastic_mirzakhani read_size_in_bytes:0 write_size_in_bytes:0] 8b9:map[cpu:0 disk:[] memory:0 name:/sleepy_chaum read_size_in_bytes:0 write_size_in_bytes:0]]
type ContainerMetrics struct {
	CPUTime       uint64
	CPUPercentage float64
	Memory        uint64
	Name          string
	ReadSize      uint64
	WriteSize     uint64
}

// {09 Dec 2025 23:50:18 UTC 1  243876101}
type Images struct {
	ID                               string
	ImageNames                       []string
	CreatedDate                      string
	NumberOfContainersUsingThisImage int64
	ImageSize                        int64
}

type Volumes struct {
	VolumeName string
	MountPoint string
	Scope      string
	Driver     string
	VolumeSize int64
	Created    string
}

func (d *DockerInfo) Collect() error {
	d.getSystemInfo()
	d.getContainerDiskUsage()
	d.getPlatformInfo()
	container_stats, err := listAllContainers()
	if err != nil {
		return err
	}
	d.ContainerStats = container_stats

	image_stats, err := listAllImages()
	if err != nil {
		return err
	}
	d.ImageStats = image_stats

	volume_stats, err := listAllVolumes()
	if err != nil {
		return err
	}

	container_collection, err := getContainerCollectionMetrics()
	if err != nil {
		return err
	}
	d.ContainerCpuMemoryCollection = container_collection

	d.VolumeStats = volume_stats
	d.TotalContainers = len(d.ContainerStats)
	d.TotalImages = len(d.ImageStats)
	d.TotalVolumes = len(d.VolumeStats)

	return nil
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func getClient() (*client.Client, error) {
	return client.New(client.FromEnv)
}

func listAllContainers() ([]Containers, error) {
	ctx, cancel := getContext()
	defer cancel()

	api_client, err := getClient()
	if err != nil {
		logging.Error(logtag, "failed to create docker client", err)
		return nil, err
	}
	defer api_client.Close()

	result, err := api_client.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		logging.Error(logtag, "failed to list containers", err)
		return nil, err
	}

	list_of_containers := make([]Containers, 0, len(result.Items))

	for _, item := range result.Items {

		list_of_containers = append(list_of_containers, Containers{
			ID:                       item.ID,
			ImageName:                item.Image,
			ContainerNames:           item.Names,
			ContainerPorts:           item.Ports,
			ContainerState:           item.State,
			ContainerRootSizeInBytes: item.SizeRootFs,
		})
	}

	return list_of_containers, nil
}

func listAllImages() ([]Images, error) {
	ctx, cancel := getContext()
	defer cancel()

	api_client, err := getClient()
	defer api_client.Close()

	if err != nil {
		logging.Error(logtag, "failed to create docker client", err)
		return nil, err
	}

	image_list_result, err := api_client.ImageList(ctx, client.ImageListOptions{All: true})
	if err != nil {
		logging.Error(logtag, "error getting image list", err)
		return nil, err
	}
	results := make([]Images, 0, len(image_list_result.Items))
	for _, image := range image_list_result.Items {
		results = append(results, Images{
			ID:                               image.ID,
			ImageNames:                       image.RepoTags,
			CreatedDate:                      time.Unix(image.Created, 0).Format(time.RFC822),
			NumberOfContainersUsingThisImage: image.Containers,
			ImageSize:                        image.Size,
		})
	}
	return results, nil

}

func calculateCPUPercent(
	id string,
	curr container.StatsResponse,
) float64 {

	prev, ok := prevCPUStats[id]
	if !ok {
		prevCPUStats[id] = curr
		return 0.0 // first sample
	}

	cpuDelta := float64(
		curr.CPUStats.CPUUsage.TotalUsage -
			prev.CPUStats.CPUUsage.TotalUsage,
	)

	systemDelta := float64(
		curr.CPUStats.SystemUsage -
			prev.CPUStats.SystemUsage,
	)

	onlineCPUs := float64(curr.CPUStats.OnlineCPUs)
	if onlineCPUs == 0 {
		onlineCPUs = float64(len(curr.CPUStats.CPUUsage.PercpuUsage))
	}

	prevCPUStats[id] = curr

	if systemDelta > 0 && cpuDelta > 0 {
		return (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}

	return 0.0
}

func getContainerCollectionMetrics() (map[string]ContainerMetrics, error) {
	ctx, cancel := getContext()
	defer cancel()

	api_client, err := getClient()
	defer api_client.Close()

	if err != nil {
		logging.Error(logtag, "failed to create docker client", err)
		return nil, err
	}
	all_containers, err := listAllContainers()
	if err != nil {
		return nil, err
	}

	container_collection_metrics := make(map[string]ContainerMetrics, len(all_containers))
	all_container_ids := make([]string, 0, len(all_containers))

	for _, container := range all_containers {
		all_container_ids = append(all_container_ids, container.ID)
	}

	for _, id := range all_container_ids {
		container_stats_result, err := api_client.ContainerStats(ctx, id, client.ContainerStatsOptions{})
		if err != nil {
			logging.Error(logtag, "failed to get container stats", err)
			return nil, err
		}
		var stats container.StatsResponse
		err = json.NewDecoder(container_stats_result.Body).Decode(&stats)

		if err != nil {
			logging.Error(logtag, "failed to decode container stats", err)
			return nil, err
		}
		defer container_stats_result.Body.Close()
		cpuPercent := calculateCPUPercent(id, stats)

		// define the container metrics nested dict
		container_collection_metrics[id] = ContainerMetrics{
			CPUTime:       stats.CPUStats.CPUUsage.TotalUsage,
			CPUPercentage: cpuPercent,
			Memory:        stats.MemoryStats.Usage,
			Name:          stats.Name,
			ReadSize:      stats.StorageStats.ReadSizeBytes,
			WriteSize:     stats.StorageStats.WriteSizeBytes,
		}

	}

	return container_collection_metrics, nil

}

func listAllVolumes() ([]Volumes, error) {
	ctx, cancel := getContext()
	defer cancel()

	api_client, err := getClient()
	defer api_client.Close()

	if err != nil {
		logging.Error(logtag, "failed to create docker client", err)
		return nil, err
	}

	volume_list_result, err := api_client.VolumeList(ctx, client.VolumeListOptions{})
	if err != nil {
		logging.Error(logtag, "failed to get volumes", err)
		return nil, err
	}
	results := make([]Volumes, 0, len(volume_list_result.Items))
	for _, i := range volume_list_result.Items {
		vol := Volumes{
			VolumeName: i.Name,
			MountPoint: i.Mountpoint,
			Scope:      i.Scope,
			Driver:     i.Driver,
			Created:    i.CreatedAt,
		}
		if i.UsageData != nil {
			vol.VolumeSize = i.UsageData.Size
		}
		results = append(results, vol)
	}
	return results, nil

}

func (d *DockerInfo) getContainerDiskUsage() error {
	ctx, cancel := getContext()
	defer cancel()

	api_client, err := getClient()
	defer api_client.Close()

	if err != nil {
		logging.Error(logtag, "failed to create docker client", err)
		return err
	}

	disk_usage_result, err := api_client.DiskUsage(ctx, client.DiskUsageOptions{
		Containers: true,
		Images:     true,
		BuildCache: true,
		Volumes:    true,
	})

	d.ContainersDiskUsage = disk_usage_result.Containers.TotalSize
	d.ImagesDiskUsage = disk_usage_result.Images.TotalSize
	d.BuildCacheDiskUsage = disk_usage_result.BuildCache.TotalSize

	if err != nil {
		logging.Error(logtag, "unable to calculate disk usage", err)
		return err
	}
	return nil
}

func (d *DockerInfo) getSystemInfo() error {
	ctx, cancel := getContext()
	defer cancel()

	api_client, err := getClient()
	defer api_client.Close()

	if err != nil {
		logging.Error(logtag, "failed to create docker client", err)
		return err
	}

	system_info, err := api_client.Info(ctx, client.InfoOptions{})
	if err != nil {
		logging.Error(logtag, "failed to get system info", err)
		return err
	}
	d.DockerEnv = system_info.Info.Name
	d.ContainersPaused = system_info.Info.ContainersPaused
	d.ContainersRunning = system_info.Info.ContainersRunning
	d.ContainersStopped = system_info.Info.ContainersStopped
	d.NCpu = system_info.Info.NCPU
	d.MemTotal = system_info.Info.MemTotal

	return nil
}

func (d *DockerInfo) getPlatformInfo() error {
	ctx, cancel := getContext()
	defer cancel()

	api_client, err := getClient()
	defer api_client.Close()

	if err != nil {
		logging.Error(logtag, "failed to create docker client", err)
		return err
	}

	server_version, err := api_client.ServerVersion(ctx, client.ServerVersionOptions{})
	if err != nil {
		logging.Error(logtag, "failed to get server version ", err)
		return err
	}
	d.PlatformName = server_version.Platform.Name
	d.Arch = server_version.Arch
	d.OS = server_version.Os
	d.APIVersion = server_version.APIVersion
	return nil

}
