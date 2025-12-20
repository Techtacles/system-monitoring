package docker

import (
	"context"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

const logtag = "docker"

type DockerInfo struct {
	DockerEnv           string // eg docker-desktop
	ContainersRunning   int
	ContainersPaused    int
	ContainersStopped   int
	NCpu                int
	MemTotal            int64
	ContainersDiskUsage int64
	ImagesDiskUsage     int64
	BuildCacheDiskUsage int64
	PlatformName        string
	APIVersion          string
	OS                  string
	Arch                string
	TotalImages         int
	TotalContainers     int
	TotalVolumes        int
	ContainerStats      []Containers
	ImageStats          []Images
	VolumeStats         []Volumes
}

// example output: {nginx [/sleepy_chaum] [{invalid IP 80 0 tcp}] running 0}
type Containers struct {
	ImageName                string
	ContainerNames           []string
	ContainerPorts           []container.PortSummary
	ContainerState           container.ContainerState
	ContainerRootSizeInBytes int64
}

// {09 Dec 2025 23:50:18 UTC 1  243876101}
type Images struct {
	ImageNames                       []string
	CreatedDate                      string
	NumberOfContainersUsingThisImage int64
	ImageAnnotation                  string
	ImageSize                        int64
}

type Volumes struct {
	VolumeName string
	MountPoint string
	Scope      string
	Driver     string
	VolumeSize int64
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
			ImageNames:                       image.RepoTags,
			CreatedDate:                      time.Unix(image.Created, 0).Format("02 Jan 2006 15:04:05 UTC"),
			NumberOfContainersUsingThisImage: image.Containers,
			ImageSize:                        image.Size,
		})
	}
	return results, nil

}

func listAllVolumes() ([]Volumes, error) {
	ctx, cancel := getContext()
	defer cancel()

	api_client, err := getClient()

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
		results = append(results, Volumes{
			VolumeName: i.Name,
			MountPoint: i.Mountpoint,
			Scope:      i.Scope,
			Driver:     i.Driver,
			VolumeSize: i.UsageData.Size,
		})
	}
	return results, nil

}

func (d *DockerInfo) getContainerDiskUsage() error {
	ctx, cancel := getContext()
	defer cancel()

	api_client, err := getClient()

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
