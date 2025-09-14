package config

import (
    "errors"
    "time"

    env "github.com/caarlos0/env/v10"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
    Env              string        `env:"APP_ENV" envDefault:"development"`
    Port             int           `env:"PORT" envDefault:"3000"`
    RequestTimeout   time.Duration `env:"REQUEST_TIMEOUT" envDefault:"15s"`
    BodyLimitBytes   int64         `env:"BODY_LIMIT_BYTES" envDefault:"10485760"` // 10 MiB

    // CORS
    CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:"," envDefault:"*"`
    CORSAllowedMethods []string `env:"CORS_ALLOWED_METHODS" envSeparator:"," envDefault:"GET,POST,PUT,PATCH,DELETE,OPTIONS"`
    CORSAllowedHeaders []string `env:"CORS_ALLOWED_HEADERS" envSeparator:"," envDefault:"Accept,Authorization,Content-Type,X-Requested-With"`

    // Rate limiting
    RateLimitEnabled bool   `env:"RATE_LIMIT_ENABLED" envDefault:"true"`
    RateLimitPeriod  string `env:"RATE_LIMIT_PERIOD" envDefault:"1m"` // parsed at runtime
    RateLimit        int    `env:"RATE_LIMIT" envDefault:"100"`       // requests per period per IP

    // CORS strict mode: fail startup in production if origins include "*"
    CORSStrict bool `env:"CORS_STRICT" envDefault:"false"`

    // Compression level (1-9)
    CompressionLevel int `env:"COMPRESSION_LEVEL" envDefault:"5"`

    // Observability
    TracingEnabled     bool   `env:"TRACING_ENABLED" envDefault:"false"`
    MetricsEnabled    bool   `env:"METRICS_ENABLED" envDefault:"true"`
    MetricsPath       string `env:"METRICS_PATH" envDefault:"/metrics"`
}

// Load parses environment variables into Config and validates values.
func Load() (*Config, error) {
    var cfg Config
    if err := env.Parse(&cfg); err != nil {
        return nil, err
    }
    if cfg.Port <= 0 || cfg.Port > 65535 {
        return nil, errors.New("invalid PORT")
    }
    if cfg.RequestTimeout <= 0 {
        return nil, errors.New("REQUEST_TIMEOUT must be > 0")
    }
    if cfg.BodyLimitBytes <= 0 || cfg.BodyLimitBytes > 1<<30 { // cap at 1 GiB
        return nil, errors.New("BODY_LIMIT_BYTES must be between 1 and 1073741824 (1GiB)")
    }
    if cfg.RateLimitEnabled && cfg.RateLimit <= 0 {
        return nil, errors.New("RATE_LIMIT must be > 0 when RATE_LIMIT_ENABLED=true")
    }
    if cfg.CompressionLevel < 1 || cfg.CompressionLevel > 9 {
        return nil, errors.New("COMPRESSION_LEVEL must be between 1 and 9")
    }
    return &cfg, nil
}
