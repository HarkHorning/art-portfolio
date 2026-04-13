package repo

import (
	"fmt"
	"strings"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/config"
	"github.com/jmoiron/sqlx"
)

func DBConnect(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	var dsn string
	if strings.HasPrefix(cfg.Host, "/cloudsql/") {
		// Cloud Run uses Unix socket for Cloud SQL
		dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Database,
		)
	} else {
		// Standard TCP connection (local dev, GKE, etc.)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if cfg.SeedData {
		if err := SeedDevData(db); err != nil {
			return nil, fmt.Errorf("failed to seed data: %w", err)
		}
	}

	return db, nil
}
