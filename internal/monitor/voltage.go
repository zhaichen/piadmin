package monitor

import (
	"os/exec"
	"strconv"
	"strings"
)

type VoltageInfo struct {
	Core     float64      `json:"core"`
	SdramC   float64      `json:"sdram_c"`
	SdramI   float64      `json:"sdram_i"`
	SdramP   float64      `json:"sdram_p"`
	Throttle ThrottleInfo `json:"throttle"`
}

type ThrottleInfo struct {
	Raw                  uint64 `json:"raw"`
	UnderVoltage         bool   `json:"under_voltage"`
	FreqCapped           bool   `json:"freq_capped"`
	Throttled            bool   `json:"throttled"`
	SoftTempLimit        bool   `json:"soft_temp_limit"`
	UnderVoltageOccurred bool   `json:"under_voltage_occurred"`
	FreqCappedOccurred   bool   `json:"freq_capped_occurred"`
	ThrottledOccurred    bool   `json:"throttled_occurred"`
	SoftTempLimitOccurred bool  `json:"soft_temp_limit_occurred"`
}

func collectVoltage(snap *SystemSnapshot) {
	if _, err := exec.LookPath("vcgencmd"); err != nil {
		return
	}

	snap.Voltage = VoltageInfo{
		Core:   readVoltage("core"),
		SdramC: readVoltage("sdram_c"),
		SdramI: readVoltage("sdram_i"),
		SdramP: readVoltage("sdram_p"),
	}

	snap.Voltage.Throttle = readThrottle()
}

func readVoltage(component string) float64 {
	out, err := exec.Command("vcgencmd", "measure_volts", component).Output()
	if err != nil {
		return 0
	}
	// output: "volt=1.2000V\n"
	s := strings.TrimSpace(string(out))
	s = strings.TrimPrefix(s, "volt=")
	s = strings.TrimSuffix(s, "V")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func readThrottle() ThrottleInfo {
	out, err := exec.Command("vcgencmd", "get_throttled").Output()
	if err != nil {
		return ThrottleInfo{}
	}
	// output: "throttled=0x50000\n"
	s := strings.TrimSpace(string(out))
	s = strings.TrimPrefix(s, "throttled=")
	raw, err := strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 64)
	if err != nil {
		return ThrottleInfo{}
	}
	return ThrottleInfo{
		Raw:                   raw,
		UnderVoltage:          raw&(1<<0) != 0,
		FreqCapped:            raw&(1<<1) != 0,
		Throttled:             raw&(1<<2) != 0,
		SoftTempLimit:         raw&(1<<3) != 0,
		UnderVoltageOccurred:  raw&(1<<16) != 0,
		FreqCappedOccurred:    raw&(1<<17) != 0,
		ThrottledOccurred:     raw&(1<<18) != 0,
		SoftTempLimitOccurred: raw&(1<<19) != 0,
	}
}
