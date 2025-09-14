package logging

import (
    "log/slog"
    "os"
)

// New returns a structured slog.Logger configured for the given environment.
// Env can be: development, production, or test. Defaults to JSON for consistency.
// Set PRETTY_LOGS=true to enable human-readable text format instead of JSON.
func New(env string) *slog.Logger {
    level := new(slog.LevelVar)
    // Default level: Info; Development: Debug
    if env == "development" || env == "dev" {
        level.Set(slog.LevelDebug)
    } else {
        level.Set(slog.LevelInfo)
    }

    // Check for pretty logging
    prettyLogs := os.Getenv("PRETTY_LOGS") == "true"

    var handler slog.Handler
    if prettyLogs {
        // Use custom pretty handler for human-readable output
        handler = NewPrettyHandler(os.Stdout, &slog.HandlerOptions{
            Level:     level,
            AddSource: false,
        })
    } else {
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level:     level,
            AddSource: env == "development",
        })
    }

    return slog.New(handler)
}

