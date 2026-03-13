package gpio

import (
	"fmt"
	"os"
	"path/filepath"
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

const gpioBase = "/sys/class/gpio"

func isRaspberryPi() bool {
	return runtime.GOOS == "linux"
}

func Available() bool {
	if !isRaspberryPi() {
		return false
	}
	_, err := os.Stat(gpioBase)
	return err == nil
}

func ListPins() ([]PinInfo, error) {
	if !Available() {
		return nil, fmt.Errorf("GPIO not available on this platform")
	}

	// Common Raspberry Pi GPIO pins (BCM numbering)
	bcmPins := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}

	var result []PinInfo
	for _, pin := range bcmPins {
		info := PinInfo{Number: pin, Available: true}

		dir := filepath.Join(gpioBase, fmt.Sprintf("gpio%d", pin))
		if _, err := os.Stat(dir); err != nil {
			// pin not exported
			info.Mode = "unexported"
			info.Available = true
			result = append(result, info)
			continue
		}

		// read direction
		dirData, err := os.ReadFile(filepath.Join(dir, "direction"))
		if err == nil {
			info.Mode = strings.TrimSpace(string(dirData))
		}

		// read value
		valData, err := os.ReadFile(filepath.Join(dir, "value"))
		if err == nil {
			info.Value, _ = strconv.Atoi(strings.TrimSpace(string(valData)))
		}

		result = append(result, info)
	}
	return result, nil
}

func ExportPin(pin int) error {
	dir := filepath.Join(gpioBase, fmt.Sprintf("gpio%d", pin))
	if _, err := os.Stat(dir); err == nil {
		return nil // already exported
	}
	return os.WriteFile(filepath.Join(gpioBase, "export"), []byte(strconv.Itoa(pin)), 0644)
}

func UnexportPin(pin int) error {
	return os.WriteFile(filepath.Join(gpioBase, "unexport"), []byte(strconv.Itoa(pin)), 0644)
}

func SetDirection(pin int, direction string) error {
	if direction != "in" && direction != "out" {
		return fmt.Errorf("direction must be 'in' or 'out'")
	}
	if err := ExportPin(pin); err != nil {
		return err
	}
	return os.WriteFile(
		filepath.Join(gpioBase, fmt.Sprintf("gpio%d", pin), "direction"),
		[]byte(direction), 0644,
	)
}

func SetValue(pin int, value int) error {
	if value != 0 && value != 1 {
		return fmt.Errorf("value must be 0 or 1")
	}
	return os.WriteFile(
		filepath.Join(gpioBase, fmt.Sprintf("gpio%d", pin), "value"),
		[]byte(strconv.Itoa(value)), 0644,
	)
}

func ReadValue(pin int) (int, error) {
	data, err := os.ReadFile(filepath.Join(gpioBase, fmt.Sprintf("gpio%d", pin), "value"))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}
