package process

import (
	"sort"

	"github.com/shirou/gopsutil/v3/process"
)

type Info struct {
	PID     int32   `json:"pid"`
	Name    string  `json:"name"`
	Status  string  `json:"status"`
	CPU     float64 `json:"cpu"`
	Memory  float32 `json:"memory"`
	User    string  `json:"user"`
	Command string  `json:"command"`
}

func List() ([]Info, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var result []Info
	for _, p := range procs {
		name, _ := p.Name()
		status, _ := p.Status()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		user, _ := p.Username()
		cmdline, _ := p.Cmdline()

		statusStr := ""
		if len(status) > 0 {
			statusStr = status[0]
		}

		result = append(result, Info{
			PID:     p.Pid,
			Name:    name,
			Status:  statusStr,
			CPU:     cpuPercent,
			Memory:  memPercent,
			User:    user,
			Command: cmdline,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CPU > result[j].CPU
	})

	return result, nil
}

func Kill(pid int32, force bool) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	if force {
		return p.Kill()
	}
	return p.Terminate()
}
