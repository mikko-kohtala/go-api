package logging

import (
    "log/slog"
    "os"
)

// New returns a structured slog.Logger configured for the given environment.
// Env can be: development, production, or test. Defaults to JSON for consistency.
func New(env string) *slog.Logger {
    level := new(slog.LevelVar)
    // Default level: Info; Development: Debug
    if env == "development" || env == "dev" {
        level.Set(slog.LevelDebug)
    } else {
        level.Set(slog.LevelInfo)
    }
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: env == "development"})
    return slog.New(handler)
}

