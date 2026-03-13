package monitor

import (
	"github.com/shirou/gopsutil/v3/net"
)

type NetworkInfo struct {
	Name      string `json:"name"`
	BytesSent uint64 `json:"bytes_sent"`
	BytesRecv uint64 `json:"bytes_recv"`
	PacksSent uint64 `json:"packs_sent"`
	PacksRecv uint64 `json:"packs_recv"`
}

func collectNetwork(snap *SystemSnapshot) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return
	}
	for _, c := range counters {
		if c.Name == "lo" || c.Name == "lo0" {
			continue
		}
		snap.Network = append(snap.Network, NetworkInfo{
			Name:      c.Name,
			BytesSent: c.BytesSent,
			BytesRecv: c.BytesRecv,
			PacksSent: c.PacketsSent,
			PacksRecv: c.PacketsRecv,
		})
	}
}
