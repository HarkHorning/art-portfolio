package main

import (
	"log/slog"
	"os"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/api"
	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/config"
	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/repo"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()
	config.SetupLogging(cfg.Environment)

	slog.Info("starting server",
		"environment", cfg.Environment,
		"port", cfg.Server.Port,
	)

	if err := repo.RunMigrations(cfg.Database); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}

	db, err := repo.DBConnect(cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("connected to database", "host", cfg.Database.Host)

	router := api.Routes(db)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
