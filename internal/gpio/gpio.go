package gpio

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type PinInfo struct {
	Number    int    `json:"number"`
	Mode      string `json:"mode"`
	Value     int    `json:"value"`
	Available bool   `json:"available"`
}

type PinRequest struct {
	Pin       int    `json:"pin"`
	Direction string `json:"direction"` // "in" or "out"
	Value     int    `json:"value"`     // 0 or 1
}

const gpioChip = "gpiochip0"

func Available() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	_, err := exec.LookPath("gpioget")
	return err == nil
}

func ListPins() ([]PinInfo, error) {
	if !Available() {
		return nil, fmt.Errorf("GPIO not available (gpioget/gpioset not found)")
	}

	// Common Raspberry Pi GPIO pins (BCM numbering)
	bcmPins := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}

	// Get all line info at once
	lineInfo := getLineInfo()

	var result []PinInfo
	for _, pin := range bcmPins {
		info := PinInfo{Number: pin, Available: true}

		if li, ok := lineInfo[pin]; ok {
			info.Mode = li.direction
			info.Value = li.value
		} else {
			info.Mode = "input"
			// read current value
			val, err := readPin(pin)
			if err == nil {
				info.Value = val
			}
		}

		result = append(result, info)
	}
	return result, nil
}

type lineInfoEntry struct {
	direction string
	value     int
}

func getLineInfo() map[int]lineInfoEntry {
	result := make(map[int]lineInfoEntry)

	out, err := exec.Command("gpioinfo", gpioChip).Output()
	if err != nil {
		return result
	}

	for _, line := range strings.Split(string(out), "\n") {
		// format: "	line   2:      unnamed       input  active-high"
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "line") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		numStr := strings.TrimSuffix(parts[1], ":")
		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}

		direction := "input"
		for _, p := range parts {
			if p == "output" {
				direction = "output"
				break
			}
		}

		entry := lineInfoEntry{direction: direction}

		// try to read current value
		val, err := readPin(num)
		if err == nil {
			entry.value = val
		}

		result[num] = entry
	}

	return result
}

func readPin(pin int) (int, error) {
	out, err := exec.Command("gpioget", gpioChip, strconv.Itoa(pin)).Output()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(out)))
}

func ExportPin(pin int) error {
	// With libgpiod, pins don't need explicit export
	// Just verify we can read it
	_, err := readPin(pin)
	return err
}

func UnexportPin(pin int) error {
	// No-op with libgpiod (no persistent export concept)
	return nil
}

func SetDirection(pin int, direction string) error {
	if direction != "in" && direction != "out" {
		return fmt.Errorf("direction must be 'in' or 'out'")
	}
	if direction == "out" {
		// set output with default value 0
		return exec.Command("gpioset", gpioChip, fmt.Sprintf("%d=0", pin)).Run()
	}
	// for input, just read to verify
	_, err := readPin(pin)
	return err
}

func SetValue(pin int, value int) error {
	if value != 0 && value != 1 {
		return fmt.Errorf("value must be 0 or 1")
	}
	return exec.Command("gpioset", gpioChip, fmt.Sprintf("%d=%d", pin, value)).Run()
}

func ReadValue(pin int) (int, error) {
	return readPin(pin)
}
