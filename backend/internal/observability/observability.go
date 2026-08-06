package observability

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"cmd/api/internal/db"
	"cmd/api/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	requestCount        atomic.Int64
	totalRequestNanos   atomic.Int64
	ticketsCreatedCount atomic.Int64
)

type logEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	UserID    string `json:"user_id,omitempty"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Duration  string `json:"duration"`
}

func RequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)

		startedAt := time.Now()
		c.Next()

		duration := time.Since(startedAt)
		requestCount.Add(1)
		totalRequestNanos.Add(duration.Nanoseconds())

		userID := ""
		if value, exists := c.Get("user_id"); exists {
			userID = fmt.Sprint(value)
		}

		entry := logEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Level:     "info",
			Message:   "request completed",
			RequestID: requestID,
			UserID:    userID,
			Method:    c.Request.Method,
			Path:      c.FullPath(),
			Status:    c.Writer.Status(),
			Duration:  duration.String(),
		}

		encoded, err := json.Marshal(entry)
		if err != nil {
			return
		}

		_, _ = os.Stdout.Write(append(encoded, '\n'))
	}
}

func RecordTicketCreated() {
	ticketsCreatedCount.Add(1)
}

func Health(c *gin.Context) {
	status := gin.H{
		"api":     "up",
		"service": "ready",
	}

	if conn := db.GetDB(); conn != nil {
		sqlDB, err := conn.DB()
		if err != nil || sqlDB.Ping() != nil {
			status["database"] = "down"
			c.JSON(http.StatusServiceUnavailable, status)
			return
		}
		status["database"] = "connected"
	} else {
		status["database"] = "down"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}

func Metrics(c *gin.Context) {
	ticketCount, err := repositories.NewTicketsRepository().Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activeUsers, err := repositories.NewRefreshTokenRepository().CountActiveUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	avgDurationSeconds := 0.0
	count := requestCount.Load()
	if count > 0 {
		avgDurationSeconds = float64(totalRequestNanos.Load()) / float64(count) / float64(time.Second)
	}

	metrics := fmt.Sprintf(`
# HELP taskforge_tickets_total Total tickets stored in the system
# TYPE taskforge_tickets_total gauge
taskforge_tickets_total %d
# HELP taskforge_tickets_created_total Tickets created during the current process lifetime
# TYPE taskforge_tickets_created_total counter
taskforge_tickets_created_total %d
# HELP taskforge_api_request_duration_seconds_avg Average API response duration in seconds
# TYPE taskforge_api_request_duration_seconds_avg gauge
taskforge_api_request_duration_seconds_avg %.6f
# HELP taskforge_logged_in_users Current number of active refresh tokens
# TYPE taskforge_logged_in_users gauge
taskforge_logged_in_users %d
`,
		ticketCount,
		ticketsCreatedCount.Load(),
		avgDurationSeconds,
		activeUsers,
	)

	c.Data(http.StatusOK, "text/plain; version=0.0.4", []byte(metrics))
}
