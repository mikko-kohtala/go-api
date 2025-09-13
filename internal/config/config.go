package config

import (
    "errors"
    "time"

    env "github.com/caarlos0/env/v10"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
    Env              string        `env:"APP_ENV" envDefault:"development"`
    Port             int           `env:"PORT" envDefault:"8080"`
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
    return &cfg, nil
}
