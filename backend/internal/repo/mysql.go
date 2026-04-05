package repo

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration

	InitSchema bool // dev field
	SeedData   bool // dev field
}

func DevConfig() Config {
	return Config{
		Host:            "127.0.0.1",
		Port:            3306,
		User:            "root",
		Password:        "Developing",
		Database:        "Development",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

func EnvConfig() Config {
	cfg := DevConfig()

	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Password = password
	}
	if database := os.Getenv("DB_NAME"); database != "" {
		cfg.Database = database
	}
	if maxOpen := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpen != "" {
		if m, err := strconv.Atoi(maxOpen); err == nil {
			cfg.MaxOpenConns = m
		}
	}
	if maxIdle := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdle != "" {
		if m, err := strconv.Atoi(maxIdle); err == nil {
			cfg.MaxIdleConns = m
		}
	}
	if lifetime := os.Getenv("DB_CONN_MAX_LIFETIME"); lifetime != "" {
		if d, err := time.ParseDuration(lifetime); err == nil {
			cfg.ConnMaxLifetime = d
		}
	}

	return cfg
}

func DBConnect(cfg Config) (*sqlx.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("SERVER: failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("SERVER: failed to ping database: %w", err)
	}

	if cfg.InitSchema { // temporary for development
		if err := InitSchema(db); err != nil {
			return nil, fmt.Errorf("SERVER: Failed to init schema: %w", err)
		}
	}
	if cfg.SeedData { // temporary for development
		if err := SeedDevData(db); err != nil {
			return nil, fmt.Errorf("SERVER: Failed to insert data: %w", err)
		}
	}

	return db, nil
}
