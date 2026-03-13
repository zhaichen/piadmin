package monitor

import (
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

type TemperatureInfo struct {
	SensorKey   string  `json:"sensor_key"`
	Temperature float64 `json:"temperature"`
}

func collectTemperature(snap *SystemSnapshot) {
	temps, err := host.SensorsTemperatures()
	if err == nil && len(temps) > 0 {
		for _, t := range temps {
			if t.Temperature > 0 {
				snap.Temperature = append(snap.Temperature, TemperatureInfo{
					SensorKey:   t.SensorKey,
					Temperature: t.Temperature,
				})
			}
		}
		return
	}

	// fallback: read sysfs directly (Raspberry Pi)
	data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return
	}
	snap.Temperature = append(snap.Temperature, TemperatureInfo{
		SensorKey:   "cpu-thermal",
		Temperature: val / 1000.0,
	})
}
