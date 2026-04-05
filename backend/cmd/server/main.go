package main

import (
	"log"

	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/api"
	"github.com/HarkHorning/portfolio-go-svelte-azure-k8/internal/repo"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	cfg := repo.EnvConfig()
	cfg.InitSchema = true
	cfg.SeedData = true
	db, err := repo.DBConnect(cfg)
	if err != nil {
		log.Fatalf("SERVER: Failed to connect to database: %v", err)
	}
	defer db.Close()
	router := api.Routes(db)

	log.Println("SERVER: Starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("SERVER: Failed to start server:", err)
	}
}
