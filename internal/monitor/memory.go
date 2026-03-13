package monitor

import (
	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	UsedPercent float64 `json:"used_percent"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
}

func collectMemory(snap *SystemSnapshot) {
	v, err := mem.VirtualMemory()
	if err == nil {
		snap.Memory.Total = v.Total
		snap.Memory.Used = v.Used
		snap.Memory.Available = v.Available
		snap.Memory.UsedPercent = v.UsedPercent
	}

	s, err := mem.SwapMemory()
	if err == nil {
		snap.Memory.SwapTotal = s.Total
		snap.Memory.SwapUsed = s.Used
	}
}
