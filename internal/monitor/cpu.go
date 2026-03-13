package monitor

import (
	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUInfo struct {
	ModelName  string    `json:"model_name"`
	Cores      int       `json:"cores"`
	Threads    int       `json:"threads"`
	UsageTotal float64   `json:"usage_total"`
	UsagePer   []float64 `json:"usage_per"`
}

func collectCPU(snap *SystemSnapshot) {
	info, err := cpu.Info()
	if err == nil && len(info) > 0 {
		snap.CPU.ModelName = info[0].ModelName
		snap.CPU.Cores = int(info[0].Cores)
	}

	counts, err := cpu.Counts(true)
	if err == nil {
		snap.CPU.Threads = counts
	}

	percentTotal, err := cpu.Percent(0, false)
	if err == nil && len(percentTotal) > 0 {
		snap.CPU.UsageTotal = percentTotal[0]
	}

	percentPer, err := cpu.Percent(0, true)
	if err == nil {
		snap.CPU.UsagePer = percentPer
	}
}
