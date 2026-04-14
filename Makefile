.PHONY: local local-down local-logs local-detach db db-down dev clean cloudrun cloudrun-logs help

# Default target
help:
	@echo ""
	@echo "Portfolio Development Commands"
	@echo "=============================="
	@echo ""
	@echo "Local Development:"
	@echo "  make local        - Build and run all containers"
	@echo "  make local-down   - Stop all containers"
	@echo "  make local-logs   - Tail container logs"
	@echo "  make local-detach - Run all containers in background"
	@echo ""
	@echo "Native Development (backend runs on your machine):"
	@echo "  make dev          - Start DB container, run backend natively"
	@echo "  make db           - Start MySQL container only"
	@echo "  make db-down      - Stop MySQL container"
	@echo ""
	@echo "Production (requires GCP setup):"
	@echo "  make cloudrun     - Deploy to Google Cloud Run"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean        - Stop all containers and wipe data volumes"
	@echo ""

# ==============================================================================
# Local Development (fully containerized)
# ==============================================================================

local:
	powershell -NoProfile -Command "docker compose -f deployment/docker/docker-compose.yml up --build"

local-down:
	powershell -NoProfile -Command "docker compose -f deployment/docker/docker-compose.yml down"

local-logs:
	powershell -NoProfile -Command "docker compose -f deployment/docker/docker-compose.yml logs -f"

local-detach:
	powershell -NoProfile -Command "docker compose -f deployment/docker/docker-compose.yml up --build -d"

# ==============================================================================
# Native Development (DB in Docker, backend runs on your machine)
# ==============================================================================

dev:
	powershell -NoProfile -Command "docker compose -f deployment/docker/docker-compose.yml up -d --wait mysql"
	powershell -NoProfile -Command "$$env:DB_PORT='3307'; $$env:DB_SEED_DATA='true'; $$env:ADMIN_USERNAME='admin'; $$env:ADMIN_PASSWORD='devpassword'; $$env:ADMIN_SESSION_SECRET='dev-secret'; $$env:GCS_BUCKET='hark-portfolio-images'; Set-Location backend; go run ./cmd/server"

db:
	powershell -NoProfile -Command "docker compose -f deployment/docker/docker-compose.yml up -d mysql"

db-down:
	powershell -NoProfile -Command "docker compose -f deployment/docker/docker-compose.yml stop mysql"

# ==============================================================================
# Google Cloud Run
# ==============================================================================

cloudrun:
	@echo "Deploying to Cloud Run..."
	powershell -NoProfile -Command "bash deployment/cloudrun/deploy.sh"

cloudrun-logs:
	powershell -NoProfile -Command "gcloud run logs read backend --region us-central1 --limit 50"

# ==============================================================================
# Cleanup
# ==============================================================================

clean:
	powershell -NoProfile -Command "docker compose -f deployment/docker/docker-compose.yml down -v"
	@echo "Cleaned up containers and volumes"
