package monitor

import (
	"sync"
	"time"
)

type HistoryEntry struct {
	Timestamp   time.Time         `json:"timestamp"`
	Temperature []TemperatureInfo `json:"temperature"`
	Voltage     VoltageInfo       `json:"voltage"`
}

type HistoryBuffer struct {
	mu      sync.RWMutex
	entries []HistoryEntry
	maxSize int
}

func NewHistoryBuffer(maxSize int) *HistoryBuffer {
	return &HistoryBuffer{
		entries: make([]HistoryEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

func (h *HistoryBuffer) Add(timestamp time.Time, temp []TemperatureInfo, voltage VoltageInfo) {
	h.mu.Lock()
	defer h.mu.Unlock()

	entry := HistoryEntry{
		Timestamp:   timestamp,
		Temperature: temp,
		Voltage:     voltage,
	}

	if len(h.entries) >= h.maxSize {
		// shift left by 1
		copy(h.entries, h.entries[1:])
		h.entries[len(h.entries)-1] = entry
	} else {
		h.entries = append(h.entries, entry)
	}
}

func (h *HistoryBuffer) Get() []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]HistoryEntry, len(h.entries))
	copy(result, h.entries)
	return result
}

func (h *HistoryBuffer) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = h.entries[:0]
}
