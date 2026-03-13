package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Auth    AuthConfig    `yaml:"auth"`
	Monitor MonitorConfig `yaml:"monitor"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type AuthConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Password string `yaml:"password"`
	TokenTTL int    `yaml:"token_ttl"`
}

type MonitorConfig struct {
	Interval int `yaml:"interval"`
}

func (c *Config) MonitorInterval() time.Duration {
	return time.Duration(c.Monitor.Interval) * time.Second
}

func (c *Config) TokenTTLDuration() time.Duration {
	return time.Duration(c.Auth.TokenTTL) * time.Second
}

func defaults() *Config {
	return &Config{
		Server: ServerConfig{
			Addr: ":8080",
		},
		Auth: AuthConfig{
			Enabled:  true,
			Password: "",
			TokenTTL: 86400,
		},
		Monitor: MonitorConfig{
			Interval: 2,
		},
	}
}

func Load(path string) *Config {
	cfg := defaults()

	if path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			_ = yaml.Unmarshal(data, cfg)
		}
	} else {
		for _, p := range []string{"./piadmin.yaml", "/etc/piadmin/config.yaml"} {
			data, err := os.ReadFile(p)
			if err == nil {
				_ = yaml.Unmarshal(data, cfg)
				break
			}
		}
	}

	applyEnv(cfg)

	if cfg.Auth.Enabled && cfg.Auth.Password == "" {
		cfg.Auth.Password = generatePassword()
		fmt.Printf("Generated password: %s\n", cfg.Auth.Password)
	}

	return cfg
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("PIADMIN_ADDR"); v != "" {
		cfg.Server.Addr = v
	}
	if v := os.Getenv("PIADMIN_PASSWORD"); v != "" {
		cfg.Auth.Password = v
	}
	if v := os.Getenv("PIADMIN_AUTH_ENABLED"); v != "" {
		cfg.Auth.Enabled = v == "true" || v == "1"
	}
	if v := os.Getenv("PIADMIN_MONITOR_INTERVAL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Monitor.Interval = n
		}
	}
}

func generatePassword() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
