package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Zahid-Iqbal-Marth/Golang-Test-Project/utils"
	"github.com/gin-gonic/gin"
)

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
