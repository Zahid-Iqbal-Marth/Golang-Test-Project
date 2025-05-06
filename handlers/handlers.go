package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Zahid-Iqbal-Marth/Golang-Test-Project/utils"
	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		log.Printf("Request started: method=%s, path=%s, remote_addr=%s",
			c.Request.Method, c.Request.URL.Path, c.ClientIP())

		c.Next()

		log.Printf("Request completed: method=%s, path=%s, duration=%v, status=%d",
			c.Request.Method, c.Request.URL.Path, time.Since(start), c.Writer.Status())
	}
}

// HealthzHandler handles health check requests
func HealthzHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// LogHandler handles payload logging requests
func LogHandler(processor *utils.BatchProcessor) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("ERROR: Failed to read request body: %v", err)
			c.String(http.StatusBadRequest, "Failed to read request body")
			return
		}

		var payload utils.Payload
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("ERROR: Failed to parse JSON payload: %v", err)
			c.String(http.StatusBadRequest, "Invalid JSON payload")
			return
		}

		processor.AddPayload(payload)

		c.Status(http.StatusAccepted)
	}
}
