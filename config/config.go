package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// JwtSecret is used to sign JWT access tokens.
	JwtSecret string `envconfig:"JWT_SECRET"`

	// JaegerHost is the OpenTelemetry Jaeger collector endpoint.
	JaegerHost string `envconfig:"JAEGER_HOST"`

	// ServiceName is used for logging, tracing and metrics labels.
	ServiceName string `envconfig:"SERVICE_NAME"`

	// PublicHTTPAddr is the public HTTP server port.
	PublicHTTPAddr string `envconfig:"PUBLIC_HTTP_ADDR"`

	Environment string `envconfig:"ENVIRONMENT"`

	DBConfig    DBConfig
	RedisConfig RedisConfig
}

type DBConfig struct {
	// Host is the database host.
	Host string `envconfig:"PG_HOST"`
	// Port is the database port.
	Port string `envconfig:"PG_PORT"`
	// User is the database user.
	User string `envconfig:"PG_USER"`
	// Password is the database password.
	Password string `envconfig:"PG_PASSWORD"`
}

// GetDSN returns a postgres DSN for gorm/pgx.
func (cfg *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.User)
}

type RedisConfig struct {
	// Host is the redis host.
	Host string `envconfig:"REDIS_HOST"`
	// Port is the redis port.
	Port string `envconfig:"REDIS_PORT"`
	// Username is the redis username (ACL). If empty, the "default" user is used.
	Username string `envconfig:"REDIS_USERNAME"`
	// Password is the redis password.
	Password string `envconfig:"REDIS_PASSWORD"`
}

// GetRedisDSN returns a redis URI.
func (cfg *Config) GetRedisDSN() string {
	return fmt.Sprintf("redis://:%s@%s:%s/0", cfg.RedisConfig.Password, cfg.RedisConfig.Host, cfg.RedisConfig.Port)
}

// GetConfig reads environment variables and returns the parsed Config.
func GetConfig() (*Config, error) {
	config := &Config{}

	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// InitENV loads `.env.local` (optional) and `.env` (required) from the given directory.
func InitENV(dir string) error {
	if err := godotenv.Load(filepath.Join(dir, ".env.local")); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("godotenv.Load: %w", err)
		}
	}

	if err := godotenv.Load(filepath.Join(dir, ".env")); err != nil {
		return fmt.Errorf("godotenv.Load: %w", err)
	}
	return nil
}
