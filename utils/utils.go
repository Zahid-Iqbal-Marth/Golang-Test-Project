package utils

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Payload struct {
	UserID    int64   `json:"user_id"`
	Total     float64 `json:"total"`
	Title     string  `json:"title"`
	Meta      Meta    `json:"meta"`
	Completed bool    `json:"completed"`
}

type Meta struct {
	Logins       []Login           `json:"logins"`
	PhoneNumbers map[string]string `json:"phone_numbers"`
}

type Login struct {
	Time time.Time `json:"time"`
	IP   string    `json:"ip"`
}

type Config struct {
	BatchSize     int
	BatchInterval time.Duration
	PostEndpoint  string
	ServerPort    string
}

func LoadConfig() (*Config, error) {
	batchSize, err := strconv.Atoi(getEnvWithDefault("BATCH_SIZE", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid BATCH_SIZE: %w", err)
	}

	batchIntervalSec, err := strconv.Atoi(getEnvWithDefault("BATCH_INTERVAL_SEC", "60"))
	if err != nil {
		return nil, fmt.Errorf("invalid BATCH_INTERVAL_SEC: %w", err)
	}

	return &Config{
		BatchSize:     batchSize,
		BatchInterval: time.Duration(batchIntervalSec) * time.Second,
		PostEndpoint:  getEnvWithDefault("POST_ENDPOINT", "https://requestbin.net/r/yourbin"),
		ServerPort:    getEnvWithDefault("SERVER_PORT", "8080"),
	}, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type BatchProcessor struct {
	payloads      []Payload
	mu            sync.Mutex
	batchSize     int
	batchInterval time.Duration
	postEndpoint  string
	timer         *time.Timer
	client        *http.Client
}
