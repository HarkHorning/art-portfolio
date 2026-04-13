package repo

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/config"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(cfg config.DatabaseConfig) error {
	slog.Info("running database migrations")

	// Open a dedicated connection for migrations with multiStatements=true.
	// Uses the standard go-sql-driver DSN format (not a URL).
	var dsn string
	if strings.HasPrefix(cfg.Host, "/cloudsql/") {
		dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true&multiStatements=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Database)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("migration db: %w", err)
	}
	defer db.Close()

	// The Docker healthcheck can pass before MySQL fully accepts client connections.
	// Retry until the connection succeeds.
	for attempt := 1; attempt <= 10; attempt++ {
		if err := db.Ping(); err == nil {
			break
		} else if attempt == 10 {
			return fmt.Errorf("database not ready after 10 attempts: %w", err)
		}
		slog.Info("waiting for database", "attempt", attempt)
		time.Sleep(time.Second)
	}

	dbDriver, err := mysqlmigrate.WithInstance(db, &mysqlmigrate.Config{})
	if err != nil {
		return fmt.Errorf("migration driver: %w", err)
	}

	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "mysql", dbDriver)
	if err != nil {
		return fmt.Errorf("migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	slog.Info("database migrations complete")
	return nil
}
