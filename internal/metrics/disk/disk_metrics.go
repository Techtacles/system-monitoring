package disk

import (
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

const logtag string = "disk"

type DiskInfo struct {
	PartitionInfo []Partitions
	UsageStat     map[string]UsagePerPath
}

type Partitions struct {
	Device     string
	MountPoint string
	Fstype     string
	Opts       []string
}

type UsagePerPath struct {
	Path        string
	TotalDisk   uint64
	FreeDisk    uint64
	UsedDisk    uint64
	UsedPercent float64
}

func (d *DiskInfo) Collect() error {
	part, err := getPartitions()
	if err != nil {
		logging.Error(logtag, "error instantiating partitions ", err)
		return err
	}

	d.PartitionInfo = part

	usage_path, err := extractDiskInfo()

	if err != nil {
		logging.Error(logtag, "error instantiating usage path metrics", err)
		return err
	}

	d.UsageStat = usage_path

	logging.Info(logtag, "successfully instantiated disk metrics")
	return nil

}

func getPartitions() ([]Partitions, error) {

	partition_stat, err := disk.Partitions(true)

	if err != nil {
		logging.Error(logtag, "error getting disk partitions", err)
		return nil, err
	}

	results := make([]Partitions, 0, len(partition_stat))

	for _, partition := range partition_stat {
		empty_partition := Partitions{}
		empty_partition.Device = partition.Device
		empty_partition.MountPoint = partition.Mountpoint
		empty_partition.Fstype = partition.Fstype
		empty_partition.Opts = partition.Opts
		results = append(results, empty_partition)

	}

	logging.Info(logtag, "partition extraction successful")

	return results, nil

}

func extractDiskInfo() (map[string]UsagePerPath, error) {
	disk_info, err := getPartitions()

	if err != nil {
		logging.Error(logtag, "error extracting disk partitions while extracting disk info", err)
	}

	results := make(map[string]UsagePerPath, len(disk_info))

	for _, v := range disk_info {
		usage_stat, err := disk.Usage(v.MountPoint)

		if err != nil {
			logging.Error(logtag, "error getting usage stat", err)
			return nil, err
		}

		empty_usage_struct := UsagePerPath{}
		empty_usage_struct.FreeDisk = usage_stat.Free
		empty_usage_struct.Path = usage_stat.Path
		empty_usage_struct.TotalDisk = usage_stat.Total
		empty_usage_struct.UsedDisk = usage_stat.Used
		empty_usage_struct.UsedPercent = usage_stat.UsedPercent
		results[usage_stat.Path] = empty_usage_struct

	}

	return results, nil

}
