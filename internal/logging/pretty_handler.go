package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorWhite  = "\033[97m"
)

// PrettyHandler implements a custom slog.Handler for pretty local logging
type PrettyHandler struct {
	opts  slog.HandlerOptions
	mu    *sync.Mutex
	out   io.Writer
	attrs []slog.Attr
}

// NewPrettyHandler creates a new pretty handler for local development
func NewPrettyHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{
		out:  out,
		opts: *opts,
		mu:   &sync.Mutex{},
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Time with milliseconds
	timeStr := r.Time.Format("15:04:05.000")

	// Level with color
	var levelStr string
	var levelColor string
	switch r.Level {
	case slog.LevelDebug:
		levelStr = "DEBUG"
		levelColor = colorGray
	case slog.LevelInfo:
		levelStr = "INFO"
		levelColor = colorGreen
	case slog.LevelWarn:
		levelStr = "WARN"
		levelColor = colorYellow
	case slog.LevelError:
		levelStr = "ERROR"
		levelColor = colorRed
	default:
		levelStr = r.Level.String()
		levelColor = colorWhite
	}

	// Extract request ID and direction from attributes
	requestID := ""
	direction := ""
	fields := make([]string, 0)

	// Process handler's attributes
	for _, attr := range h.attrs {
		if attr.Key == "request_id" {
			requestID = attr.Value.String()
		} else if attr.Key == "direction" {
			direction = attr.Value.String()
		} else if attr.Key != "component" {
			fields = append(fields, formatAttr(attr))
		}
	}

	// Process record's attributes
	r.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "request_id" {
			requestID = attr.Value.String()
		} else if attr.Key == "direction" {
			direction = attr.Value.String()
		} else if attr.Key != "component" {
			fields = append(fields, formatAttr(attr))
		}
		return true
	})

	// Build the log line
	var logLine strings.Builder

	// Time
	logLine.WriteString(colorGray)
	logLine.WriteString(timeStr)
	logLine.WriteString(colorReset)
	logLine.WriteString(" ")

	// Level
	logLine.WriteString(levelColor)
	logLine.WriteString(levelStr)
	logLine.WriteString(colorReset)
	logLine.WriteString(" ")

	// Check if this is an HTTP method log (GET, POST, etc.)
	isHTTPLog := false
	httpMethods := []string{"GET ", "POST ", "PUT ", "DELETE ", "PATCH ", "HEAD ", "OPTIONS "}
	for _, method := range httpMethods {
		if strings.HasPrefix(r.Message, method) {
			isHTTPLog = true
			break
		}
	}

	// Add arrow for HTTP requests based on direction
	if isHTTPLog && direction != "" {
		if direction == "incoming" {
			logLine.WriteString(colorCyan)
			logLine.WriteString("→ ")
			logLine.WriteString(colorReset)
		} else if direction == "outgoing" {
			logLine.WriteString(colorGreen)
			logLine.WriteString("← ")
			logLine.WriteString(colorReset)
		}
	}

	// Message
	logLine.WriteString(r.Message)

	// Add request ID if present
	if requestID != "" {
		logLine.WriteString(" ")
		logLine.WriteString(colorGray)
		logLine.WriteString("{id:\"")
		logLine.WriteString(requestID)
		logLine.WriteString("\"")
		if len(fields) > 0 {
			logLine.WriteString(", ")
			logLine.WriteString(strings.Join(fields, ", "))
		}
		logLine.WriteString("}")
		logLine.WriteString(colorReset)
	} else if len(fields) > 0 {
		// Add fields if any (no request ID)
		logLine.WriteString(" ")
		logLine.WriteString(colorGray)
		logLine.WriteString("{")
		logLine.WriteString(strings.Join(fields, ", "))
		logLine.WriteString("}")
		logLine.WriteString(colorReset)
	}

	logLine.WriteString("\n")

	_, err := h.out.Write([]byte(logLine.String()))
	return err
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := *h
	newHandler.attrs = append(newHandler.attrs, attrs...)
	return &newHandler
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we'll ignore groups in pretty mode
	return h
}

// formatAttr formats an attribute as key:value
func formatAttr(attr slog.Attr) string {
	switch attr.Value.Kind() {
	case slog.KindString:
		// Quote strings if they contain spaces
		str := attr.Value.String()
		if strings.Contains(str, " ") {
			return fmt.Sprintf("%s:\"%s\"", attr.Key, str)
		}
		return fmt.Sprintf("%s:%s", attr.Key, str)
	case slog.KindTime:
		t := attr.Value.Time()
		return fmt.Sprintf("%s:%s", attr.Key, t.Format("15:04:05"))
	case slog.KindDuration:
		d := attr.Value.Duration()
		// Format duration nicely
		if d < time.Millisecond {
			return fmt.Sprintf("%s:%.2fµs", attr.Key, float64(d.Nanoseconds())/1000)
		} else if d < time.Second {
			return fmt.Sprintf("%s:%.2fms", attr.Key, float64(d.Nanoseconds())/1e6)
		}
		return fmt.Sprintf("%s:%s", attr.Key, d.String())
	default:
		return fmt.Sprintf("%s:%v", attr.Key, attr.Value.Any())
	}
}