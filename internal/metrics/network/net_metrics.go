package network

import (
	"runtime"

	gnet "github.com/shirou/gopsutil/v4/net"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

const logtag string = "network"

type NetworkInfo struct {
	NumEstablishedConnections int // this gets the number of established
	NumTotalConnections       int
	Runtime                   string
	IOStats                   []IOInfo
	Connections               []ConnStatInfo
}

type ConnStatInfo struct {
	Family     uint32
	LocalAddr  gnet.Addr
	RemoteAddr gnet.Addr
	Status     string
	Pid        int32
	Type       uint32
}

type IOInfo struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
	Errin       uint64 `json:"errin"`
	Errout      uint64 `json:"errout"`
	Dropin      uint64 `json:"dropin"`
	Dropout     uint64 `json:"dropout"`
}

func (n *NetworkInfo) Collect() error {
	n.Runtime = runtime.GOOS

	iostat, err := collectIOStats()
	if err != nil {
		return err
	}

	constat, established_conn, total_conn, err := collectConnections()
	if err != nil {
		return err
	}
	n.IOStats = iostat
	n.Connections = constat
	n.NumEstablishedConnections = established_conn
	n.NumTotalConnections = total_conn

	return nil
}

func collectIOStats() ([]IOInfo, error) {
	logging.Info(logtag, "collecting network IO stats")

	stats, err := gnet.IOCounters(true)
	if err != nil {
		logging.Error(logtag, "failed to get IO counters", err)
		return nil, err
	}

	results := make([]IOInfo, 0, len(stats))

	for _, s := range stats {
		results = append(results, IOInfo{
			Name:        s.Name,
			BytesSent:   s.BytesSent,
			BytesRecv:   s.BytesRecv,
			PacketsSent: s.PacketsSent,
			PacketsRecv: s.PacketsRecv,
			Errin:       s.Errin,
			Errout:      s.Errout,
			Dropin:      s.Dropin,
			Dropout:     s.Dropout,
		})
	}

	return results, nil
}

func collectConnections() ([]ConnStatInfo, int, int, error) {
	logging.Info(logtag, "collecting network connections")

	conns, err := gnet.Connections("all")
	if err != nil {
		logging.Error(logtag, "failed to get connections", err)
		return nil, 0, 0, err
	}
	num_established_connection := 0
	num_total_connection := 0

	results := make([]ConnStatInfo, 0, len(conns))

	for _, c := range conns {

		// Skip LISTEN-only or empty connections. Focus on only established connections
		if c.Status == "LISTEN" || c.Status == "" {
			num_total_connection += 1
			continue

		}
		if c.Status == "ESTABLISHED" {
			num_established_connection += 1
		}

		results = append(results, ConnStatInfo{
			Family:     c.Family,
			Type:       c.Type,
			LocalAddr:  c.Laddr,
			RemoteAddr: c.Raddr,
			Status:     c.Status,
			Pid:        c.Pid,
		})
		num_total_connection = num_total_connection + num_established_connection

	}

	return results, num_established_connection, num_total_connection, nil
}
