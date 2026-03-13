package monitor

import (
	"runtime"

	"github.com/shirou/gopsutil/v3/host"
)

func collectHostInfo(snap *SystemSnapshot) {
	info, err := host.Info()
	if err != nil {
		return
	}
	snap.Hostname = info.Hostname
	snap.OS = info.OS
	snap.Platform = info.Platform + " " + info.PlatformVersion
	snap.Arch = runtime.GOARCH
	snap.KernelVer = info.KernelVersion
	snap.Uptime = info.Uptime
}
