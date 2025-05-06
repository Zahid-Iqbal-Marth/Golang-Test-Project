package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

func (bp *BatchProcessor) AddPayload(payload Payload) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.payloads = append(bp.payloads, payload)

	if len(bp.payloads) >= bp.batchSize {
		go bp.ProcessBatch()
	}
}

func (bp *BatchProcessor) ProcessBatch() {
	bp.mu.Lock()

	if len(bp.payloads) == 0 {
		bp.restartTimer()
		bp.mu.Unlock()
		return
	}

	payloads := bp.payloads
	bp.payloads = make([]Payload, 0, bp.batchSize)

	bp.restartTimer()
	bp.mu.Unlock()

	batchSize := len(payloads)
	log.Printf("Processing batch: size=%d", batchSize)

	if err := bp.sendBatch(payloads); err != nil {
		log.Printf("ERROR: Failed to process batch after retries: %v", err)
		log.Fatalf("Exiting application due to repeated failures")
	}
}

func (bp *BatchProcessor) sendBatch(payloads []Payload) error {
	jsonData, err := json.Marshal(payloads)
	if err != nil {
		return fmt.Errorf("error marshaling payloads: %w", err)
	}

	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			log.Printf("Retrying batch send: attempt=%d", attempt+1)
			time.Sleep(2 * time.Second)
		}

		start := time.Now()
		resp, err := bp.client.Post(bp.postEndpoint, "application/json", bytes.NewBuffer(jsonData))
		duration := time.Since(start)

		if err != nil {
			log.Printf("WARN: Error sending batch: %v, attempt=%d", err, attempt+1)
			continue
		}

		defer resp.Body.Close()

		log.Printf("Batch sent: size=%d, status_code=%d, duration=%v",
			len(payloads), resp.StatusCode, duration)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		log.Printf("WARN: Bad status code: %d, attempt=%d", resp.StatusCode, attempt+1)
	}

	return fmt.Errorf("failed to send batch after 3 attempts")
}

func NewBatchProcessor(batchSize int, batchInterval time.Duration, postEndpoint string) *BatchProcessor {
	bp := &BatchProcessor{
		payloads:      make([]Payload, 0, batchSize),
		batchSize:     batchSize,
		postEndpoint:  postEndpoint,
		batchInterval: batchInterval,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	bp.restartTimer()

	return bp
}

func (bp *BatchProcessor) restartTimer() {
	if bp.timer != nil {
		bp.timer.Stop()
	}

	log.Printf("Restarting timer with interval: %v", bp.batchInterval)
	bp.timer = time.AfterFunc(bp.batchInterval, func() {
		log.Printf("Timer triggered, processing batch")
		bp.ProcessBatch()
	})
}
