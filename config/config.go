package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config stores runtime settings for the API server
type Config struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	//APIKeyHeader is the header that clients must send
	APIKeyHeader string

	//APIKeys are the allowed keys
	APIKeys []string
}

// Addr returns the host:port string for http.Server
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Load reads the environment variables and falls back to safe defaults
func Load() (*Config, error) {

	port, err := envInt("PORT", 8080)
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:         envStr("HOST", "0.0.0.0"),
		Port:         port,
		ReadTimeout:  envDur("READ_TIMEOUT", 5*time.Second),
		WriteTimeout: envDur("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  envDur("IDLE_TIMEOUT", 60*time.Second),
		APIKeyHeader: envStr("API_KEY_HEADER", "X-Api-Key"),
		APIKeys:      envList("API_KEYS"),
	}, nil
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) (int, error) {
	v := os.Getenv(key)

	if v == "" {
		return def, nil
	}

	n, err := strconv.Atoi(v)

	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", key)
	}
	return n, nil
}

func envDur(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)

	if v == "" {
		return def
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}

	return d
}

func envList(key string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, " ")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
