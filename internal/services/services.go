package services

import (
	"os/exec"
	"runtime"
	"strings"
)

type ServiceInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LoadState   string `json:"load_state"`
	ActiveState string `json:"active_state"`
	SubState    string `json:"sub_state"`
}

func isLinux() bool {
	return runtime.GOOS == "linux"
}

func List() ([]ServiceInfo, error) {
	if !isLinux() {
		return nil, nil
	}

	out, err := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend").Output()
	if err != nil {
		return nil, err
	}

	var result []ServiceInfo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		name := strings.TrimSuffix(fields[0], ".service")
		desc := strings.Join(fields[4:], " ")
		result = append(result, ServiceInfo{
			Name:        name,
			LoadState:   fields[1],
			ActiveState: fields[2],
			SubState:    fields[3],
			Description: desc,
		})
	}
	return result, nil
}

func Action(name, action string) error {
	if !isLinux() {
		return nil
	}
	switch action {
	case "start", "stop", "restart", "enable", "disable":
		return exec.Command("systemctl", action, name+".service").Run()
	}
	return nil
}

func Status(name string) (string, error) {
	if !isLinux() {
		return "not available (not linux)", nil
	}
	out, err := exec.Command("systemctl", "status", name+".service", "--no-pager").CombinedOutput()
	if err != nil {
		return string(out), nil
	}
	return string(out), nil
}
