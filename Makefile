.PHONY: local local-down local-logs cloudrun clean help

# Default target
help:
	@echo ""
	@echo "Portfolio Development Commands"
	@echo "=============================="
	@echo ""
	@echo "Local Development:"
	@echo "  make local        - Run with Docker Compose"
	@echo "  make local-down   - Stop containers"
	@echo "  make local-logs   - View container logs"
	@echo "  make clean        - Stop and remove all data"
	@echo ""
	@echo "Production (requires GCP setup):"
	@echo "  make cloudrun     - Deploy to Google Cloud Run"
	@echo ""

# ==============================================================================
# Local Development
# ==============================================================================

local:
	docker compose -f deployment/docker/docker-compose.yml up --build

local-down:
	docker compose -f deployment/docker/docker-compose.yml down

local-logs:
	docker compose -f deployment/docker/docker-compose.yml logs -f

local-detach:
	docker compose -f deployment/docker/docker-compose.yml up --build -d

# ==============================================================================
# Google Cloud Run
# ==============================================================================

cloudrun:
	@echo "Deploying to Cloud Run..."
	@bash deployment/cloudrun/deploy.sh

cloudrun-logs:
	gcloud run logs read backend --region us-central1 --limit 50

# ==============================================================================
# Cleanup
# ==============================================================================

clean:
	docker compose -f deployment/docker/docker-compose.yml down -v
	@echo "Cleaned up containers and volumes"
