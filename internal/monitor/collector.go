package monitor

import (
	"sync"
	"time"
)

type SystemSnapshot struct {
	Timestamp   time.Time         `json:"timestamp"`
	Hostname    string            `json:"hostname"`
	OS          string            `json:"os"`
	Platform    string            `json:"platform"`
	Arch        string            `json:"arch"`
	KernelVer   string            `json:"kernel_version"`
	Uptime      uint64            `json:"uptime"`
	CPU         CPUInfo           `json:"cpu"`
	Memory      MemoryInfo        `json:"memory"`
	Disks       []DiskInfo        `json:"disks"`
	Network     []NetworkInfo     `json:"network"`
	Temperature []TemperatureInfo `json:"temperature"`
}

type Collector struct {
	mu       sync.RWMutex
	snapshot *SystemSnapshot
	interval time.Duration
	subs     map[chan *SystemSnapshot]struct{}
	subsMu   sync.Mutex
	stopCh   chan struct{}
}

func NewCollector(interval time.Duration) *Collector {
	return &Collector{
		interval: interval,
		subs:     make(map[chan *SystemSnapshot]struct{}),
		stopCh:   make(chan struct{}),
	}
}

func (c *Collector) Start() {
	// collect immediately
	c.collect()

	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.collect()
			case <-c.stopCh:
				return
			}
		}
	}()
}

func (c *Collector) Stop() {
	close(c.stopCh)
}

func (c *Collector) collect() {
	snap := &SystemSnapshot{
		Timestamp: time.Now(),
	}

	collectHostInfo(snap)
	collectCPU(snap)
	collectMemory(snap)
	collectDisk(snap)
	collectNetwork(snap)
	collectTemperature(snap)

	c.mu.Lock()
	c.snapshot = snap
	c.mu.Unlock()

	c.publish(snap)
}

func (c *Collector) GetSnapshot() *SystemSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.snapshot
}

func (c *Collector) Subscribe() chan *SystemSnapshot {
	ch := make(chan *SystemSnapshot, 1)
	c.subsMu.Lock()
	c.subs[ch] = struct{}{}
	c.subsMu.Unlock()
	return ch
}

func (c *Collector) Unsubscribe(ch chan *SystemSnapshot) {
	c.subsMu.Lock()
	delete(c.subs, ch)
	c.subsMu.Unlock()
	close(ch)
}

func (c *Collector) publish(snap *SystemSnapshot) {
	c.subsMu.Lock()
	defer c.subsMu.Unlock()
	for ch := range c.subs {
		select {
		case ch <- snap:
		default:
			// subscriber is slow, skip
		}
	}
}
