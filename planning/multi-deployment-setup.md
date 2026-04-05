# Multi-Deployment Setup Plan

This guide sets up a flexible deployment system using **Google Cloud only**. Switch between local development, Cloud Run (cheap), and GKE (Kubernetes) with a single command.

## End Goal

```bash
make local           # Run locally with Docker Desktop (local DB)
make local-podman    # Run locally with Podman (local DB)
make cloudrun        # Deploy to Google Cloud Run (shared Cloud SQL)
make k8s             # Deploy to GKE Kubernetes (shared Cloud SQL)
make tunnel          # Expose local dev to internet
make migrate         # Run database migrations
```

All production deployments share the same Cloud SQL database and are accessible via `https://yourdomain.com`

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         PRODUCTION                               │
│                                                                  │
│   Cloud Run ─────┐                                               │
│                  ├──────► Cloud SQL (shared MySQL)               │
│   GKE (K8s) ─────┘                                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                         LOCAL DEV                                │
│                                                                  │
│   Docker/Podman ──────► Local MySQL (Docker container)          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Key principle:** All production deployments (Cloud Run, GKE) connect to the same Cloud SQL database. Data changes in one are visible in all.

---

## Part 1: Buy Domain and Set Up Cloudflare

### 1.1 Purchase Domain via Cloudflare

1. Go to https://dash.cloudflare.com
2. Create account (free)
3. Click **Domain Registration** → **Register Domain**
4. Search for your domain (e.g., `harkhorning.dev`, `harkportfolio.com`)
5. Purchase (~$10/year for `.dev`, `.com`)

If you already own a domain elsewhere:
- At your registrar, change nameservers to Cloudflare's (they'll provide them)

### 1.2 Initial DNS Setup

Once domain is in Cloudflare:

1. Go to **DNS** → **Records**
2. For now, add a placeholder A record:
   ```
   Type: A
   Name: @
   Content: 192.0.2.1 (placeholder)
   Proxy: Yes (orange cloud)
   TTL: Auto
   ```

### 1.3 Enable HTTPS

1. Go to **SSL/TLS** → **Overview**
2. Set mode to **Full (strict)**
3. Go to **SSL/TLS** → **Edge Certificates**
4. Enable **Always Use HTTPS**

---

## Part 2: Set Up Google Cloud

### 2.1 Create GCP Project

1. Go to https://console.cloud.google.com
2. Create new project: `hark-portfolio`
3. Note your **Project ID** (e.g., `hark-portfolio-123456`)

### 2.2 Install gcloud CLI

**Windows:**
```bash
winget install Google.CloudSDK
```

**Mac:**
```bash
brew install google-cloud-sdk
```

Then authenticate:
```bash
gcloud auth login
gcloud config set project YOUR_PROJECT_ID
```

### 2.3 Enable APIs

```bash
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  cloudbuild.googleapis.com \
  container.googleapis.com \
  sqladmin.googleapis.com
```

### 2.4 Create Artifact Registry

```bash
gcloud artifacts repositories create portfolio \
  --repository-format=docker \
  --location=us-central1 \
  --description="Portfolio container images"
```

### 2.5 Authenticate Docker with GCP

```bash
gcloud auth configure-docker us-central1-docker.pkg.dev
```

---

## Part 3: Set Up Cloud SQL (Shared Production Database)

### 3.1 Create Cloud SQL Instance

```bash
gcloud sql instances create portfolio-db \
  --database-version=MYSQL_8_0 \
  --tier=db-f1-micro \
  --region=us-central1 \
  --root-password=YOUR_SECURE_PASSWORD \
  --storage-size=10GB \
  --storage-auto-increase
```

**Cost:** ~$7-10/month for db-f1-micro

### 3.2 Create Database

```bash
gcloud sql databases create portfolio --instance=portfolio-db
```

### 3.3 Get Connection Info

```bash
# Get the connection name (needed for Cloud Run)
gcloud sql instances describe portfolio-db --format='value(connectionName)'
# Output: YOUR_PROJECT_ID:us-central1:portfolio-db

# Get the public IP (needed for local connections)
gcloud sql instances describe portfolio-db --format='value(ipAddresses[0].ipAddress)'
```

### 3.4 Allow Connections

For Cloud Run (uses Cloud SQL Proxy automatically):
```bash
# No extra config needed - Cloud Run connects via Unix socket
```

For GKE and local development, authorize your IP:
```bash
# Get your current IP
curl ifconfig.me

# Authorize it
gcloud sql instances patch portfolio-db \
  --authorized-networks=YOUR_IP/32
```

Or allow all IPs (less secure, but simpler for dev):
```bash
gcloud sql instances patch portfolio-db \
  --authorized-networks=0.0.0.0/0
```

---

## Part 4: Create Deployment Scripts

### 4.1 File Structure

```
portfolio/
├── Makefile                          # Main entry point
├── deployment/
│   ├── docker/
│   │   └── docker-compose.yml        # Local development
│   ├── kubernetes/
│   │   └── *.yaml                    # GKE manifests
│   ├── cloudrun/
│   │   ├── deploy.sh                 # Cloud Run deployment
│   │   └── env.template              # Environment template
│   ├── terraform/
│   │   └── *.tf                      # GCP infrastructure
│   └── migrations/
│       └── *.sql                     # Database migrations
```

### 4.2 Create `deployment/cloudrun/env.template`

```bash
# Copy this to .env and fill in values
# DO NOT commit .env to git

# Google Cloud
GCP_PROJECT_ID=your-project-id
GCP_REGION=us-central1

# Cloud SQL
CLOUD_SQL_CONNECTION=your-project-id:us-central1:portfolio-db
DB_USER=root
DB_PASSWORD=your-secure-password
DB_NAME=portfolio

# Cloudflare
CLOUDFLARE_API_TOKEN=your-cloudflare-token
CLOUDFLARE_ZONE_ID=your-zone-id
CLOUDFLARE_DNS_RECORD_ID=your-record-id
DOMAIN=yourdomain.com
```

### 4.3 Create `deployment/cloudrun/deploy.sh`

```bash
#!/bin/bash
set -e

# Load environment
if [ -f deployment/cloudrun/.env ]; then
  source deployment/cloudrun/.env
else
  echo "Error: deployment/cloudrun/.env not found"
  echo "Copy deployment/cloudrun/env.template to .env and fill in values"
  exit 1
fi

REGISTRY="us-central1-docker.pkg.dev/${GCP_PROJECT_ID}/portfolio"

echo "=== Building and pushing images ==="

# Build and push frontend
docker build -t ${REGISTRY}/frontend:latest ./frontend
docker push ${REGISTRY}/frontend:latest

# Build and push backend
docker build -t ${REGISTRY}/backend:latest ./backend
docker push ${REGISTRY}/backend:latest

echo "=== Deploying to Cloud Run ==="

# Deploy backend with Cloud SQL connection
gcloud run deploy backend \
  --image ${REGISTRY}/backend:latest \
  --platform managed \
  --region ${GCP_REGION} \
  --allow-unauthenticated \
  --add-cloudsql-instances ${CLOUD_SQL_CONNECTION} \
  --set-env-vars "DB_HOST=/cloudsql/${CLOUD_SQL_CONNECTION},DB_USER=${DB_USER},DB_PASSWORD=${DB_PASSWORD},DB_NAME=${DB_NAME}" \
  --memory 256Mi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 2

# Get backend URL
BACKEND_URL=$(gcloud run services describe backend --region ${GCP_REGION} --format 'value(status.url)')

# Deploy frontend
gcloud run deploy frontend \
  --image ${REGISTRY}/frontend:latest \
  --platform managed \
  --region ${GCP_REGION} \
  --allow-unauthenticated \
  --set-env-vars "BACKEND_URL=${BACKEND_URL}" \
  --memory 256Mi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 2

# Get frontend URL
FRONTEND_URL=$(gcloud run services describe frontend --region ${GCP_REGION} --format 'value(status.url)')

echo "=== Deployment complete ==="
echo "Frontend: ${FRONTEND_URL}"
echo "Backend: ${BACKEND_URL}"

echo ""
echo "=== Updating Cloudflare DNS ==="

CLOUDRUN_HOST=$(echo ${FRONTEND_URL} | sed 's|https://||')

curl -X PUT "https://api.cloudflare.com/client/v4/zones/${CLOUDFLARE_ZONE_ID}/dns_records/${CLOUDFLARE_DNS_RECORD_ID}" \
  -H "Authorization: Bearer ${CLOUDFLARE_API_TOKEN}" \
  -H "Content-Type: application/json" \
  --data "{\"type\":\"CNAME\",\"name\":\"@\",\"content\":\"${CLOUDRUN_HOST}\",\"proxied\":true}"

echo "DNS updated! Site will be live at https://${DOMAIN} shortly."
```

### 4.4 Create the Makefile

Create `Makefile` in project root:

```makefile
.PHONY: local local-podman cloudrun k8s tunnel migrate clean help

# Load environment if exists
-include deployment/cloudrun/.env
export

# Default target
help:
	@echo "Portfolio Deployment Targets:"
	@echo ""
	@echo "  make local          - Run with Docker Compose (local DB)"
	@echo "  make local-podman   - Run with Podman Compose (local DB)"
	@echo "  make cloudrun       - Deploy to Google Cloud Run"
	@echo "  make k8s            - Deploy to GKE Kubernetes"
	@echo "  make tunnel         - Expose local dev via Cloudflare Tunnel"
	@echo "  make migrate        - Run migrations on production DB"
	@echo "  make clean          - Stop all local containers"
	@echo ""

# ------------------------------------------------------------------------------
# Local Development (uses local MySQL container)
# ------------------------------------------------------------------------------

local:
	docker compose -f deployment/docker/docker-compose.yml up --build

local-down:
	docker compose -f deployment/docker/docker-compose.yml down

local-podman:
	podman-compose -f deployment/docker/docker-compose.yml up --build

local-podman-down:
	podman-compose -f deployment/docker/docker-compose.yml down

# ------------------------------------------------------------------------------
# Database Migrations
# ------------------------------------------------------------------------------

migrate:
	@echo "Running migrations on Cloud SQL..."
	gcloud sql connect portfolio-db --user=root < deployment/migrations/all.sql

migrate-local:
	@echo "Running migrations on local database..."
	docker exec -i $$(docker ps -qf "name=mysql") mysql -uroot -pdevpassword portfolio < deployment/migrations/all.sql

# ------------------------------------------------------------------------------
# Google Cloud Run
# ------------------------------------------------------------------------------

cloudrun:
	@chmod +x deployment/cloudrun/deploy.sh
	@./deployment/cloudrun/deploy.sh

cloudrun-logs:
	gcloud run logs read backend --region us-central1 --limit 50

# ------------------------------------------------------------------------------
# GKE Kubernetes
# ------------------------------------------------------------------------------

k8s-create-cluster:
	gcloud container clusters create-auto portfolio-cluster \
		--region=us-central1

k8s-credentials:
	gcloud container clusters get-credentials portfolio-cluster \
		--region=us-central1

k8s-create-secret:
	kubectl create namespace portfolio --dry-run=client -o yaml | kubectl apply -f -
	kubectl create secret generic db-secret \
		--namespace=portfolio \
		--from-literal=host=$$(gcloud sql instances describe portfolio-db --format='value(ipAddresses[0].ipAddress)') \
		--from-literal=user=${DB_USER} \
		--from-literal=password=${DB_PASSWORD} \
		--dry-run=client -o yaml | kubectl apply -f -

k8s:
	kubectl apply -f deployment/kubernetes/

k8s-status:
	kubectl get all -n portfolio

k8s-down:
	kubectl delete namespace portfolio

k8s-ip:
	@kubectl get svc frontend -n portfolio -o jsonpath='{.status.loadBalancer.ingress[0].ip}'

# ------------------------------------------------------------------------------
# Cloudflare Tunnel
# ------------------------------------------------------------------------------

tunnel:
	@echo "Starting Cloudflare Tunnel..."
	cloudflared tunnel --url http://localhost:3000

# ------------------------------------------------------------------------------
# Cleanup
# ------------------------------------------------------------------------------

clean:
	-docker compose -f deployment/docker/docker-compose.yml down -v
	-podman-compose -f deployment/docker/docker-compose.yml down -v 2>/dev/null || true
```

---

## Part 5: Set Up GKE (Google Kubernetes Engine)

### 5.1 Create GKE Autopilot Cluster

```bash
gcloud container clusters create-auto portfolio-cluster \
  --region=us-central1
```

### 5.2 Get Cluster Credentials

```bash
gcloud container clusters get-credentials portfolio-cluster \
  --region=us-central1
```

### 5.3 Update Kubernetes Manifests for GCP

Update image references in all deployment files:
```yaml
# Change from Azure ACR:
image: harkportfolioacr.azurecr.io/backend:v1

# To GCP Artifact Registry:
image: us-central1-docker.pkg.dev/YOUR_PROJECT_ID/portfolio/backend:latest
```

Update `backend-deployment.yaml` to connect to Cloud SQL:
```yaml
env:
  - name: DB_HOST
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: host
  - name: DB_USER
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: user
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: password
  - name: DB_PORT
    value: "3306"
  - name: DB_NAME
    value: "portfolio"
```

### 5.4 Create Kubernetes Secret

```bash
make k8s-create-secret
```

---

## Part 6: Terraform for GCP Infrastructure

### 6.1 Create `deployment/terraform/main.tf`

```hcl
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Artifact Registry
resource "google_artifact_registry_repository" "portfolio" {
  location      = var.region
  repository_id = "portfolio"
  format        = "DOCKER"
}

# Cloud SQL
resource "google_sql_database_instance" "portfolio" {
  name             = "portfolio-db"
  database_version = "MYSQL_8_0"
  region           = var.region

  settings {
    tier = "db-f1-micro"

    ip_configuration {
      ipv4_enabled = true
      authorized_networks {
        name  = "allow-all"
        value = "0.0.0.0/0"
      }
    }
  }

  deletion_protection = false
}

resource "google_sql_database" "portfolio" {
  name     = "portfolio"
  instance = google_sql_database_instance.portfolio.name
}

resource "google_sql_user" "root" {
  name     = "root"
  instance = google_sql_database_instance.portfolio.name
  password = var.db_password
}

# GKE Autopilot Cluster
resource "google_container_cluster" "portfolio" {
  name     = "portfolio-cluster"
  location = var.region

  enable_autopilot = true

  network    = "default"
  subnetwork = "default"
}
```

### 6.2 Create `deployment/terraform/variables.tf`

```hcl
variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "db_password" {
  description = "Cloud SQL root password"
  type        = string
  sensitive   = true
}
```

### 6.3 Create `deployment/terraform/outputs.tf`

```hcl
output "artifact_registry_url" {
  value = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.portfolio.repository_id}"
}

output "cloud_sql_ip" {
  value = google_sql_database_instance.portfolio.public_ip_address
}

output "cloud_sql_connection" {
  value = google_sql_database_instance.portfolio.connection_name
}

output "gke_cluster_name" {
  value = google_container_cluster.portfolio.name
}
```

---

## Part 7: Cloudflare Tunnel Setup

### 7.1 Install Cloudflared

**Windows:**
```bash
winget install Cloudflare.cloudflared
```

**Mac:**
```bash
brew install cloudflared
```

### 7.2 Authenticate

```bash
cloudflared tunnel login
```

### 7.3 Use It

```bash
make tunnel
```

---

## Part 8: Get Cloudflare API Credentials

### 8.1 Get Zone ID

1. Cloudflare Dashboard → your domain
2. Right sidebar → **Zone ID**

### 8.2 Create API Token

1. **My Profile** → **API Tokens**
2. **Create Token** → Use template: **Edit zone DNS**
3. Zone Resources: Include → Specific Zone → your domain

### 8.3 Get DNS Record ID

```bash
curl -X GET "https://api.cloudflare.com/client/v4/zones/YOUR_ZONE_ID/dns_records" \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

---

## Part 9: GitHub Actions for GCP

Create `.github/workflows/deploy.yaml`:

```yaml
name: Build and Deploy

on:
  push:
    branches:
      - main

env:
  GCP_PROJECT_ID: ${{ secrets.GCP_PROJECT_ID }}
  GCP_REGION: us-central1
  REGISTRY: us-central1-docker.pkg.dev

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Auth to Google Cloud
        uses: google-github-actions/auth@v2
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY }}

      - name: Set up Cloud SDK
        uses: google-github-actions/setup-gcloud@v2

      - name: Configure Docker
        run: gcloud auth configure-docker ${GCP_REGION}-docker.pkg.dev

      - name: Build and push frontend
        run: |
          docker build -t ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/frontend:${{ github.sha }} ./frontend
          docker push ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/frontend:${{ github.sha }}

      - name: Build and push backend
        run: |
          docker build -t ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/backend:${{ github.sha }} ./backend
          docker push ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/backend:${{ github.sha }}

      - name: Deploy to Cloud Run
        run: |
          gcloud run deploy backend \
            --image ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/backend:${{ github.sha }} \
            --region ${GCP_REGION} \
            --allow-unauthenticated \
            --add-cloudsql-instances ${{ secrets.CLOUD_SQL_CONNECTION }} \
            --set-env-vars "DB_HOST=/cloudsql/${{ secrets.CLOUD_SQL_CONNECTION }},DB_USER=${{ secrets.DB_USER }},DB_PASSWORD=${{ secrets.DB_PASSWORD }},DB_NAME=portfolio"

          gcloud run deploy frontend \
            --image ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/frontend:${{ github.sha }} \
            --region ${GCP_REGION} \
            --allow-unauthenticated
```

---

## Part 10: Update .gitignore

```
# Cloud Run
deployment/cloudrun/.env

# Cloudflare
.cloudflared/

# Terraform
*.tfvars
*.tfstate
*.tfstate.*
.terraform/
.terraform.lock.hcl
```

---

## Checklist

### One-time Setup
- [ ] Purchase domain on Cloudflare
- [ ] Set SSL mode to Full (strict)
- [ ] Create GCP project
- [ ] Enable required APIs
- [ ] Create Artifact Registry
- [ ] Create Cloud SQL instance
- [ ] Create `deployment/cloudrun/.env`
- [ ] Create Cloudflare API token
- [ ] Install cloudflared

### For GKE (optional)
- [ ] Create GKE Autopilot cluster
- [ ] Get cluster credentials
- [ ] Create db-secret in cluster
- [ ] Update K8s manifests with GCP registry URLs

### Testing
- [ ] `make local` works
- [ ] `make cloudrun` deploys successfully
- [ ] `make k8s` deploys successfully
- [ ] Domain resolves with HTTPS

---

## Cost Summary

| Resource | Monthly Cost |
|----------|--------------|
| Cloud SQL (db-f1-micro) | ~$7-10 |
| Cloud Run (low traffic) | ~$0-5 |
| GKE Autopilot (when used) | ~$30+ |
| Artifact Registry | ~$0.10/GB |
| Cloudflare | Free |
| **Total (Cloud Run only)** | **~$10-15/mo** |
| **Total (with GKE)** | **~$40-50/mo** |

---

## Quick Reference

| Command | Database | Cost |
|---------|----------|------|
| `make local` | Local MySQL | Free |
| `make cloudrun` | Cloud SQL (shared) | ~$10/mo |
| `make k8s` | Cloud SQL (shared) | ~$40/mo |
| `make tunnel` | Local MySQL | Free |

---

## Troubleshooting

### Cloud Run: Can't connect to Cloud SQL
- Verify `--add-cloudsql-instances` flag is set
- Check connection name format: `project:region:instance`
- DB_HOST should be `/cloudsql/CONNECTION_NAME` for Cloud Run

### GKE: Can't connect to Cloud SQL
- Authorize GKE node IPs or use Cloud SQL Proxy sidecar
- Simpler: authorize 0.0.0.0/0 for dev (not recommended for production)

### Cloudflare: Too many redirects
- Set SSL mode to **Full (strict)**
