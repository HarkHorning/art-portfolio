package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Environment string

const (
	EnvLocal    Environment = "local"
	EnvCloudRun Environment = "cloudrun"
	EnvK8s      Environment = "k8s"
)

type Config struct {
	Environment Environment
	Server      ServerConfig
	Database    DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	SeedData bool
}

func Load() Config {
	env := detectEnvironment()

	cfg := Config{
		Environment: env,
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "127.0.0.1"),
			Port:            getEnvInt("DB_PORT", 3306),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", "devpassword"),
			Database:        getEnv("DB_NAME", "portfolio"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			SeedData: getEnvBool("DB_SEED_DATA", false),
		},
	}

	return cfg
}

func detectEnvironment() Environment {
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		switch env {
		case "cloudrun":
			return EnvCloudRun
		case "k8s", "kubernetes":
			return EnvK8s
		default:
			return EnvLocal
		}
	}

	// Auto-detect based on GCP markers
	if os.Getenv("K_SERVICE") != "" {
		return EnvCloudRun
	}
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return EnvK8s
	}

	return EnvLocal
}

func SetupLogging(env Environment) {
	var handler slog.Handler

	if env == EnvLocal {
		// Human-readable text for local dev
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		// JSON for cloud environments (Cloud Logging, etc.)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	slog.SetDefault(slog.New(handler))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
