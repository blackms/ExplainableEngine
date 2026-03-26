package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// --- statusRecorder tests ---

func TestStatusRecorderCapturesCode(t *testing.T) {
	rec := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: rec, status: http.StatusOK}

	sr.WriteHeader(http.StatusNotFound)

	if sr.status != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", sr.status)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("underlying recorder should also have 404, got %d", rec.Code)
	}
}

func TestStatusRecorderDefaultsTo200(t *testing.T) {
	rec := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: rec, status: http.StatusOK}

	// Write body without calling WriteHeader explicitly.
	_, _ = sr.Write([]byte("ok"))

	if sr.status != http.StatusOK {
		t.Fatalf("expected default status 200, got %d", sr.status)
	}
}

// --- LoggingMiddleware tests ---

func TestLoggingMiddlewareOutputsJSON(t *testing.T) {
	// Capture stdout.
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	t.Setenv("LOG_LEVEL", "INFO")

	handler := LoggingMiddleware(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("X-Request-Id", "test-id-123")
		rw.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/explain/1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = origStdout

	var entry logEntry
	if err := json.Unmarshal(out, &entry); err != nil {
		t.Fatalf("expected valid JSON log, got error: %v\nraw: %s", err, string(out))
	}

	if entry.Method != "GET" {
		t.Errorf("expected method GET, got %s", entry.Method)
	}
	if entry.Path != "/api/v1/explain/1" {
		t.Errorf("expected path /api/v1/explain/1, got %s", entry.Path)
	}
	if entry.Status != 200 {
		t.Errorf("expected status 200, got %d", entry.Status)
	}
	if entry.RequestID != "test-id-123" {
		t.Errorf("expected request_id test-id-123, got %s", entry.RequestID)
	}
	if entry.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestLoggingMiddlewareWarningLevelSkipsSuccess(t *testing.T) {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	t.Setenv("LOG_LEVEL", "WARNING")

	handler := LoggingMiddleware(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = origStdout

	if strings.TrimSpace(string(out)) != "" {
		t.Fatalf("expected no output at WARNING level for 200, got: %s", string(out))
	}
}

func TestLoggingMiddlewareWarningLevelLogsErrors(t *testing.T) {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	t.Setenv("LOG_LEVEL", "WARNING")

	handler := LoggingMiddleware(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))

	req := httptest.NewRequest(http.MethodGet, "/fail", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = origStdout

	var entry logEntry
	if err := json.Unmarshal(out, &entry); err != nil {
		t.Fatalf("expected JSON log for error at WARNING level, got: %s", string(out))
	}
	if entry.Status != 500 {
		t.Errorf("expected status 500, got %d", entry.Status)
	}
}

// --- CORS middleware tests ---

func TestCORSMiddlewareDefaultAllowsAll(t *testing.T) {
	handler := CORSMiddleware("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected *, got %s", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got != "GET, POST, OPTIONS" {
		t.Errorf("unexpected Allow-Methods: %s", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); got != "Content-Type, X-Request-Id" {
		t.Errorf("unexpected Allow-Headers: %s", got)
	}
}

func TestCORSMiddlewareSpecificOrigin(t *testing.T) {
	handler := CORSMiddleware("http://a.com, http://b.com")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Allowed origin.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://a.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://a.com" {
		t.Errorf("expected http://a.com, got %s", got)
	}

	// Disallowed origin.
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.Header.Set("Origin", "http://evil.com")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if got := rec2.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected empty Allow-Origin for disallowed, got %s", got)
	}
}

func TestCORSMiddlewarePreflightReturns204(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("should not reach"))
	})

	handler := CORSMiddleware("*")(inner)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/explain", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS preflight, got %d", rec.Code)
	}
	if body := rec.Body.String(); body != "" {
		t.Errorf("expected empty body for preflight, got %s", body)
	}
}
