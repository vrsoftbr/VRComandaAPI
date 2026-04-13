package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler exposes a lightweight liveness endpoint for monitoring.
// @Summary Health check
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthHandler(c *gin.Context) {
	RespondOK(c, http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
