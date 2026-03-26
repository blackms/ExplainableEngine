package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// logLevel represents the minimum severity level for log output.
type logLevel int

const (
	logLevelDebug   logLevel = iota
	logLevelInfo
	logLevelWarning
)

// logEntry is the structured JSON log record emitted per request.
type logEntry struct {
	Timestamp  string `json:"timestamp"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Status     int    `json:"status"`
	DurationMs int64  `json:"duration_ms"`
	RequestID  string `json:"request_id"`
}

// statusRecorder wraps http.ResponseWriter to capture the response status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code and delegates to the wrapped writer.
func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// parseLogLevel converts a LOG_LEVEL string to a logLevel constant.
// Defaults to logLevelInfo for unrecognised values.
func parseLogLevel(s string) logLevel {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG":
		return logLevelDebug
	case "WARNING":
		return logLevelWarning
	default:
		return logLevelInfo
	}
}

// LoggingMiddleware emits a structured JSON log line for every request.
// The log level is read from the LOG_LEVEL environment variable on each
// request so that it can be changed at runtime without a restart.
//
// At INFO (default) and DEBUG levels every request is logged.
// At WARNING level only responses with status >= 400 are logged.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		level := parseLogLevel(os.Getenv("LOG_LEVEL"))
		if level == logLevelWarning && rec.status < 400 {
			return
		}

		entry := logEntry{
			Timestamp:  start.UTC().Format(time.RFC3339),
			Method:     r.Method,
			Path:       r.URL.Path,
			Status:     rec.status,
			DurationMs: time.Since(start).Milliseconds(),
			RequestID:  w.Header().Get("X-Request-Id"),
		}

		data, err := json.Marshal(entry)
		if err != nil {
			return
		}
		fmt.Fprintln(os.Stdout, string(data))
	})
}
