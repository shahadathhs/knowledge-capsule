package logger

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"knowledge-capsule/pkg/contextkeys"
)

// Event types for structured logging
const (
	EventRequest   = "request"
	EventAuth      = "auth"
	EventUser      = "user"
	EventCapsule   = "capsule"
	EventTopic     = "topic"
	EventAdmin     = "admin"
	EventSearch    = "search"
	EventUpload    = "upload"
	EventChat      = "chat"
	EventSeed      = "seed"
	EventError     = "error"
	EventPanic    = "panic"
)

// FromRequest extracts common request context for logging.
func FromRequest(r *http.Request) []slog.Attr {
	attrs := []slog.Attr{
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	}
	if r.URL.RawQuery != "" {
		attrs = append(attrs, slog.String("query", r.URL.RawQuery))
	}
	if userID, ok := r.Context().Value(contextkeys.UserContextKey).(string); ok && userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}
	if role, ok := r.Context().Value(contextkeys.RoleContextKey).(string); ok && role != "" {
		attrs = append(attrs, slog.String("role", role))
	}
	return attrs
}

// Info logs an info-level event with attrs.
func Info(event string, attrs ...slog.Attr) {
	slog.Default().Info(event, toAny(attrs)...)
}

// InfoRequest logs a request completion with status and duration.
func InfoRequest(r *http.Request, status int, duration time.Duration, attrs ...slog.Attr) {
	all := append(FromRequest(r),
		slog.Int("status", status),
		slog.Duration("duration_ms", duration),
	)
	all = append(all, attrs...)
	level := slog.LevelInfo
	if status >= 400 {
		level = slog.LevelWarn
	}
	slog.Default().Log(r.Context(), level, EventRequest, toAny(all)...)
}

// Error logs an error event.
func Error(event string, err error, attrs ...slog.Attr) {
	all := append([]slog.Attr{slog.String("error", err.Error())}, attrs...)
	slog.Default().Error(event, toAny(all)...)
}

// ErrorRequest logs an error in request context.
func ErrorRequest(r *http.Request, event string, err error, attrs ...slog.Attr) {
	all := append(FromRequest(r), slog.String("error", err.Error()))
	all = append(all, attrs...)
	slog.Default().Error(event, toAny(all)...)
}

// Debug logs a debug-level event.
func Debug(event string, attrs ...slog.Attr) {
	slog.Default().Debug(event, toAny(attrs)...)
}

// LogEvent logs a domain event (auth, user, capsule, etc.) with context.
func LogEvent(event string, r *http.Request, attrs ...slog.Attr) {
	all := append(FromRequest(r), attrs...)
	slog.Default().Info(event, toAny(all)...)
}

// Attr creates a slog.Attr (convenience).
func Attr(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

func toAny(attrs []slog.Attr) []any {
	out := make([]any, 0, len(attrs)*2)
	for _, a := range attrs {
		out = append(out, a.Key, a.Value)
	}
	return out
}

// WithContext returns a logger with request context attached.
func WithContext(ctx context.Context, attrs ...slog.Attr) *slog.Logger {
	return slog.Default().With(toAny(attrs)...)
}
