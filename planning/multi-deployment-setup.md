# Multi-Deployment Setup Plan

This guide sets up a flexible deployment system using **Google Cloud only**. Switch between local development, Cloud Run (cheap), and GKE (Kubernetes) with a single command.

**Time estimate:** 2-3 hours for initial setup

**Prerequisites:**
- A credit/debit card (for GCP billing - you won't be charged during free trial)
- Docker Desktop installed
- Git installed
- A terminal (PowerShell on Windows, Terminal on Mac)

---

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

## Glossary (What These Things Are)

| Term | What It Is | Why You Need It |
|------|------------|-----------------|
| **GCP** | Google Cloud Platform - Google's cloud services | Hosts your app in production |
| **Cloud Run** | Serverless container hosting - you upload a container, Google runs it | Cheap way to run your app (~$0-5/mo) |
| **GKE** | Google Kubernetes Engine - managed Kubernetes | For interviews / learning K8s |
| **Cloud SQL** | Managed MySQL database | Stores your data in production |
| **Artifact Registry** | Container image storage (like Docker Hub but private) | Stores your Docker images |
| **Cloudflare** | DNS and CDN provider | Custom domain + free HTTPS |
| **gcloud** | Command-line tool for GCP | How you control GCP from terminal |

---

## Part 1: Buy Domain and Set Up Cloudflare

### 1.1 Create Cloudflare Account

1. Go to https://dash.cloudflare.com
2. Click **Sign Up**
3. Enter email and password
4. Verify your email

### 1.2 Purchase Domain

1. In Cloudflare dashboard, click **Domain Registration** in the left sidebar
2. Click **Register Domain**
3. Search for your domain (e.g., `harkhorning.dev`, `harkportfolio.com`)
4. Add to cart and purchase (~$10/year for `.dev`, `.com`)
5. Complete payment

**Already own a domain elsewhere?** You can transfer it to Cloudflare or just use Cloudflare as DNS:
- At your current registrar (GoDaddy, Namecheap, etc.), find "Nameservers" settings
- Change nameservers to Cloudflare's (they'll show you what to enter)
- Wait 1-24 hours for DNS to propagate

### 1.3 Initial DNS Setup

Once domain is in Cloudflare:

1. Click on your domain to open its dashboard
2. Click **DNS** in the left sidebar
3. Click **Add Record**
4. Fill in:
   ```
   Type: A
   Name: @ (this means the root domain, like example.com)
   Content: 192.0.2.1 (placeholder - we'll change this later)
   Proxy status: Proxied (orange cloud ON)
   TTL: Auto
   ```
5. Click **Save**

### 1.4 Enable HTTPS

1. In left sidebar, click **SSL/TLS** → **Overview**
2. Under "SSL/TLS encryption mode", select **Full (strict)**
3. Click **SSL/TLS** → **Edge Certificates**
4. Find "Always Use HTTPS" and toggle it **ON**

Your domain now has HTTPS. Cloudflare handles certificates automatically.

---

## Part 2: Create Google Cloud Account

### 2.1 Sign Up for GCP

1. Go to https://cloud.google.com
2. Click **Get started for free** (top right)
3. Sign in with your Google account (or create one)
4. You'll need to add a credit card, but:
   - You get $300 free credits for 90 days
   - You won't be charged unless you manually upgrade
   - Even after free trial, this project costs ~$10-15/month

### 2.2 Create a Project

In GCP, everything lives in a "project" - think of it as a folder for your app.

1. Go to https://console.cloud.google.com
2. At the top, click the project dropdown (might say "Select a project")
3. Click **New Project**
4. Enter:
   - Project name: `hark-portfolio` (or whatever you want)
   - Organization: Leave as is (probably "No organization")
5. Click **Create**
6. Wait a few seconds, then select your new project from the dropdown

**Important:** Note your **Project ID** - it's shown under the project name and looks like `hark-portfolio-123456`. You'll need this later.

### 2.3 Enable Billing

GCP requires billing to be enabled, even with free credits.

1. In GCP Console, click the hamburger menu (☰) → **Billing**
2. If prompted, link your project to a billing account
3. If you don't have a billing account, create one (this is where you entered your card)

---

## Part 3: Install Google Cloud CLI

The `gcloud` CLI lets you control GCP from your terminal.

### 3.1 Install gcloud

**Windows (PowerShell as Administrator):**
```powershell
winget install Google.CloudSDK
```

After installation, **close and reopen your terminal**.

**Mac:**
```bash
brew install google-cloud-sdk
```

**Linux:**
```bash
# Download and run the installer
curl https://sdk.cloud.google.com | bash
# Restart your shell
exec -l $SHELL
```

### 3.2 Verify Installation

```bash
gcloud --version
```

You should see version info. If you get "command not found", restart your terminal.

### 3.3 Authenticate

This connects your terminal to your Google account:

```bash
gcloud auth login
```

1. A browser window opens
2. Select your Google account
3. Click "Allow"
4. You'll see "You are now authenticated"

### 3.4 Set Your Project

Tell gcloud which project to use (replace with YOUR project ID):

```bash
gcloud config set project YOUR_PROJECT_ID
```

Example:
```bash
gcloud config set project hark-portfolio-123456
```

### 3.5 Enable Required APIs

GCP has many services, but they're disabled by default. Enable the ones you need:

```bash
gcloud services enable run.googleapis.com
gcloud services enable artifactregistry.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable container.googleapis.com
gcloud services enable sqladmin.googleapis.com
```

**What each API does:**
- `run.googleapis.com` - Cloud Run (serverless containers)
- `artifactregistry.googleapis.com` - Store Docker images
- `cloudbuild.googleapis.com` - Build containers in the cloud
- `container.googleapis.com` - GKE (Kubernetes)
- `sqladmin.googleapis.com` - Cloud SQL (database)

Each command takes a few seconds. You might see "Operation finished successfully".

### 3.6 Create Artifact Registry

This is where your Docker images will be stored:

```bash
gcloud artifacts repositories create portfolio --repository-format=docker --location=us-central1 --description="Portfolio container images"
```

### 3.7 Connect Docker to GCP

This lets you push Docker images to Artifact Registry:

```bash
gcloud auth configure-docker us-central1-docker.pkg.dev
```

When prompted, type `Y` and press Enter.

---

## Part 4: Set Up Cloud SQL (Production Database)

Cloud SQL is a managed MySQL database. Google handles backups, updates, and security.

### 4.1 Create the Database Instance

**This command takes 5-10 minutes to complete.** Don't close your terminal.

```bash
gcloud sql instances create portfolio-db --database-version=MYSQL_8_0 --tier=db-f1-micro --region=us-central1 --root-password=CHOOSE_A_SECURE_PASSWORD --storage-size=10GB --storage-auto-increase
```

**Replace `CHOOSE_A_SECURE_PASSWORD` with an actual password!** Write it down - you'll need it later.

**Cost:** ~$7-10/month for `db-f1-micro` (smallest tier)

### 4.2 Create the Database

```bash
gcloud sql databases create portfolio --instance=portfolio-db
```

### 4.3 Get Connection Info

You'll need these values later. Run each command and save the output:

```bash
# Get the connection name (for Cloud Run)
gcloud sql instances describe portfolio-db --format="value(connectionName)"
```

Output looks like: `hark-portfolio-123456:us-central1:portfolio-db`

```bash
# Get the public IP (for connecting from your computer)
gcloud sql instances describe portfolio-db --format="value(ipAddresses[0].ipAddress)"
```

Output looks like: `34.123.45.67`

### 4.4 Allow Connections

By default, Cloud SQL blocks all connections. You need to allow your IP:

```bash
# First, get your current public IP
curl ifconfig.me
```

Note the IP address shown (e.g., `98.76.54.32`).

```bash
# Allow your IP to connect (replace with YOUR IP)
gcloud sql instances patch portfolio-db --authorized-networks=YOUR_IP/32
```

Example:
```bash
gcloud sql instances patch portfolio-db --authorized-networks=98.76.54.32/32
```

**Note:** If your IP changes (common with home internet), you'll need to run this again.

**Simpler but less secure option** (allows any IP to connect):
```bash
gcloud sql instances patch portfolio-db --authorized-networks=0.0.0.0/0
```

---

## Part 5: Create Deployment Files

### 5.1 Folder Structure

Create these folders and files in your project:

```
portfolio/
├── Makefile
├── deployment/
│   ├── docker/
│   │   └── docker-compose.yml    (already exists)
│   ├── kubernetes/
│   │   └── *.yaml                (already exists)
│   └── cloudrun/
│       ├── deploy.sh
│       └── env.template
```

### 5.2 Create Environment Template

Create file `deployment/cloudrun/env.template`:

```bash
# Google Cloud Run Environment Configuration
#
# SETUP INSTRUCTIONS:
# 1. Copy this file to `.env` in the same folder:
#    cp env.template .env
# 2. Fill in all the values below
# 3. NEVER commit .env to git (it's in .gitignore)

# ============================================================================
# GOOGLE CLOUD
# ============================================================================

# Your GCP Project ID (find it at https://console.cloud.google.com)
# Example: hark-portfolio-123456
GCP_PROJECT_ID=

# GCP Region (us-central1 is cheapest and has good availability)
GCP_REGION=us-central1

# ============================================================================
# CLOUD SQL DATABASE
# ============================================================================

# Connection name from: gcloud sql instances describe portfolio-db --format="value(connectionName)"
# Example: hark-portfolio-123456:us-central1:portfolio-db
CLOUD_SQL_CONNECTION=

# Database user (root is fine for this project)
DB_USER=root

# The password you set when creating Cloud SQL
DB_PASSWORD=

# Database name
DB_NAME=portfolio

# ============================================================================
# CLOUDFLARE
# ============================================================================

# API Token - create at: https://dash.cloudflare.com/profile/api-tokens
# Use template "Edit zone DNS", select your zone
CLOUDFLARE_API_TOKEN=

# Zone ID - find in Cloudflare dashboard, right sidebar of your domain
CLOUDFLARE_ZONE_ID=

# DNS Record ID - see instructions in multi-deployment-setup.md Part 9
CLOUDFLARE_DNS_RECORD_ID=

# Your domain name
DOMAIN=
```

### 5.3 Create Your .env File

```bash
cd deployment/cloudrun
cp env.template .env
```

Now edit `.env` and fill in all the values.

### 5.4 Create Deploy Script

Create file `deployment/cloudrun/deploy.sh`:

```bash
#!/bin/bash
set -e

# ============================================================================
# Cloud Run Deployment Script
# ============================================================================
# This script:
# 1. Builds Docker images for frontend and backend
# 2. Pushes them to Google Artifact Registry
# 3. Deploys them to Cloud Run
# 4. Updates Cloudflare DNS to point to your app
# ============================================================================

# Load environment variables
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
if [ -f "$SCRIPT_DIR/.env" ]; then
  source "$SCRIPT_DIR/.env"
else
  echo "Error: .env file not found!"
  echo "Copy env.template to .env and fill in your values"
  exit 1
fi

# Verify required variables are set
REQUIRED_VARS="GCP_PROJECT_ID GCP_REGION CLOUD_SQL_CONNECTION DB_USER DB_PASSWORD DB_NAME"
for var in $REQUIRED_VARS; do
  if [ -z "${!var}" ]; then
    echo "Error: $var is not set in .env"
    exit 1
  fi
done

REGISTRY="us-central1-docker.pkg.dev/${GCP_PROJECT_ID}/portfolio"

echo ""
echo "=========================================="
echo "  Building and Pushing Docker Images"
echo "=========================================="
echo ""

# Navigate to project root
cd "$SCRIPT_DIR/../.."

# Build and push frontend
echo "Building frontend..."
docker build -t ${REGISTRY}/frontend:latest ./frontend
echo "Pushing frontend..."
docker push ${REGISTRY}/frontend:latest

# Build and push backend
echo "Building backend..."
docker build -t ${REGISTRY}/backend:latest ./backend
echo "Pushing backend..."
docker push ${REGISTRY}/backend:latest

echo ""
echo "=========================================="
echo "  Deploying to Cloud Run"
echo "=========================================="
echo ""

# Deploy backend
echo "Deploying backend..."
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
  --max-instances 2 \
  --quiet

# Get backend URL
BACKEND_URL=$(gcloud run services describe backend --region ${GCP_REGION} --format "value(status.url)")
echo "Backend deployed: ${BACKEND_URL}"

# Deploy frontend
echo "Deploying frontend..."
gcloud run deploy frontend \
  --image ${REGISTRY}/frontend:latest \
  --platform managed \
  --region ${GCP_REGION} \
  --allow-unauthenticated \
  --set-env-vars "BACKEND_URL=${BACKEND_URL}" \
  --memory 256Mi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 2 \
  --quiet

# Get frontend URL
FRONTEND_URL=$(gcloud run services describe frontend --region ${GCP_REGION} --format "value(status.url)")
echo "Frontend deployed: ${FRONTEND_URL}"

echo ""
echo "=========================================="
echo "  Deployment Complete!"
echo "=========================================="
echo ""
echo "Frontend: ${FRONTEND_URL}"
echo "Backend:  ${BACKEND_URL}"
echo ""

# Update Cloudflare DNS if configured
if [ -n "$CLOUDFLARE_API_TOKEN" ] && [ -n "$CLOUDFLARE_ZONE_ID" ] && [ -n "$CLOUDFLARE_DNS_RECORD_ID" ]; then
  echo "Updating Cloudflare DNS..."

  # Extract hostname from URL (remove https://)
  CLOUDRUN_HOST=$(echo ${FRONTEND_URL} | sed 's|https://||')

  curl -s -X PUT "https://api.cloudflare.com/client/v4/zones/${CLOUDFLARE_ZONE_ID}/dns_records/${CLOUDFLARE_DNS_RECORD_ID}" \
    -H "Authorization: Bearer ${CLOUDFLARE_API_TOKEN}" \
    -H "Content-Type: application/json" \
    --data "{\"type\":\"CNAME\",\"name\":\"@\",\"content\":\"${CLOUDRUN_HOST}\",\"proxied\":true}" \
    > /dev/null

  echo "DNS updated! Your site will be live at https://${DOMAIN} shortly."
else
  echo "Cloudflare not configured. Set CLOUDFLARE_* variables in .env to enable auto DNS updates."
  echo "For now, manually point your domain to: ${FRONTEND_URL}"
fi
```

Make the script executable (Mac/Linux):
```bash
chmod +x deployment/cloudrun/deploy.sh
```

### 5.5 Create Makefile

Create `Makefile` in your project root:

```makefile
.PHONY: local local-down cloudrun clean help

help:
	@echo ""
	@echo "Portfolio Deployment Commands"
	@echo "=============================="
	@echo ""
	@echo "  make local        - Run locally with Docker Compose"
	@echo "  make local-down   - Stop local containers"
	@echo "  make cloudrun     - Deploy to Google Cloud Run"
	@echo "  make clean        - Stop containers and remove data"
	@echo ""

# Local Development
local:
	docker compose -f deployment/docker/docker-compose.yml up --build

local-down:
	docker compose -f deployment/docker/docker-compose.yml down

# Google Cloud Run
cloudrun:
ifeq ($(OS),Windows_NT)
	@powershell -ExecutionPolicy Bypass -File deployment/cloudrun/deploy.ps1
else
	@bash deployment/cloudrun/deploy.sh
endif

# Cleanup
clean:
	docker compose -f deployment/docker/docker-compose.yml down -v
```

**Windows Users:** Create `deployment/cloudrun/deploy.ps1` with the PowerShell version of the deploy script (I can create this if needed).

---

## Part 6: Set Up GitHub Actions (CI/CD)

This automatically deploys your app when you push to GitHub.

### 6.1 Create a Service Account

A service account is like a robot user that GitHub Actions uses to deploy.

1. Go to https://console.cloud.google.com/iam-admin/serviceaccounts
2. Make sure your project is selected at the top
3. Click **Create Service Account**
4. Enter:
   - Name: `github-actions`
   - Description: `Used by GitHub Actions for deployment`
5. Click **Create and Continue**
6. Add these roles (click **Add Another Role** between each):
   - `Cloud Run Admin`
   - `Artifact Registry Writer`
   - `Cloud SQL Client`
   - `Service Account User`
7. Click **Continue** → **Done**

### 6.2 Create a Key for the Service Account

1. Click on the service account you just created (`github-actions@...`)
2. Click **Keys** tab
3. Click **Add Key** → **Create new key**
4. Select **JSON**
5. Click **Create**
6. A JSON file downloads - **keep this safe, don't share it!**

### 6.3 Add Secrets to GitHub

1. Go to your GitHub repo → **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret** for each of these:

| Secret Name | Value |
|-------------|-------|
| `GCP_PROJECT_ID` | Your project ID (e.g., `hark-portfolio-123456`) |
| `GCP_SA_KEY` | Entire contents of the JSON file you downloaded |
| `CLOUD_SQL_CONNECTION` | Connection name (e.g., `hark-portfolio-123456:us-central1:portfolio-db`) |
| `DB_USER` | `root` |
| `DB_PASSWORD` | Your Cloud SQL password |

### 6.4 Create GitHub Actions Workflow

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
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Authenticate to Google Cloud
        uses: google-github-actions/auth@v2
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY }}

      - name: Set up Cloud SDK
        uses: google-github-actions/setup-gcloud@v2

      - name: Configure Docker for Artifact Registry
        run: gcloud auth configure-docker ${GCP_REGION}-docker.pkg.dev --quiet

      - name: Build and push frontend
        run: |
          docker build -t ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/frontend:${{ github.sha }} ./frontend
          docker push ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/frontend:${{ github.sha }}

      - name: Build and push backend
        run: |
          docker build -t ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/backend:${{ github.sha }} ./backend
          docker push ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/backend:${{ github.sha }}

      - name: Deploy backend to Cloud Run
        run: |
          gcloud run deploy backend \
            --image ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/backend:${{ github.sha }} \
            --region ${GCP_REGION} \
            --platform managed \
            --allow-unauthenticated \
            --add-cloudsql-instances ${{ secrets.CLOUD_SQL_CONNECTION }} \
            --set-env-vars "DB_HOST=/cloudsql/${{ secrets.CLOUD_SQL_CONNECTION }},DB_USER=${{ secrets.DB_USER }},DB_PASSWORD=${{ secrets.DB_PASSWORD }},DB_NAME=portfolio"

      - name: Deploy frontend to Cloud Run
        run: |
          gcloud run deploy frontend \
            --image ${REGISTRY}/${GCP_PROJECT_ID}/portfolio/frontend:${{ github.sha }} \
            --region ${GCP_REGION} \
            --platform managed \
            --allow-unauthenticated
```

---

## Part 7: Get Cloudflare API Credentials

To automatically update DNS when you deploy, you need Cloudflare API access.

### 7.1 Get Your Zone ID

1. Go to https://dash.cloudflare.com
2. Click on your domain
3. Scroll down on the right sidebar
4. Find **Zone ID** and copy it

### 7.2 Create an API Token

1. Click your profile icon (top right) → **My Profile**
2. Click **API Tokens** in the left sidebar
3. Click **Create Token**
4. Click **Use template** next to "Edit zone DNS"
5. Under "Zone Resources":
   - Select **Include**
   - Select **Specific zone**
   - Select your domain
6. Click **Continue to summary**
7. Click **Create Token**
8. **Copy the token immediately** - you can't see it again!

### 7.3 Get the DNS Record ID

You need the ID of the DNS record you created earlier. Run this command (replace with your values):

```bash
curl -X GET "https://api.cloudflare.com/client/v4/zones/YOUR_ZONE_ID/dns_records" -H "Authorization: Bearer YOUR_API_TOKEN" -H "Content-Type: application/json"
```

This returns JSON. Look for the record where `"name"` matches your domain. Copy the `"id"` field.

Example response (simplified):
```json
{
  "result": [
    {
      "id": "abc123def456",    <-- THIS IS YOUR DNS RECORD ID
      "name": "yourdomain.com",
      "type": "A",
      "content": "192.0.2.1"
    }
  ]
}
```

### 7.4 Add to Your .env

Add these to your `deployment/cloudrun/.env`:

```
CLOUDFLARE_API_TOKEN=your-token-here
CLOUDFLARE_ZONE_ID=your-zone-id-here
CLOUDFLARE_DNS_RECORD_ID=your-record-id-here
DOMAIN=yourdomain.com
```

---

## Part 8: Update .gitignore

Make sure sensitive files aren't committed to git. Add these to `.gitignore`:

```
# Cloud Run secrets
deployment/cloudrun/.env

# Cloudflare
.cloudflared/

# Terraform
*.tfvars
*.tfstate
*.tfstate.*
.terraform/
.terraform.lock.hcl

# Service account keys (NEVER commit these)
*-key.json
*-credentials.json
```

---

## Part 9: Test Your Setup

### 9.1 Test Locally

```bash
make local
```

Open http://localhost:3000 - you should see your app.

Press `Ctrl+C` to stop.

### 9.2 Deploy to Cloud Run

```bash
make cloudrun
```

This takes a few minutes. When done, you'll see URLs for your frontend and backend.

### 9.3 Test Your Domain

If you configured Cloudflare, go to `https://yourdomain.com`. It might take a few minutes for DNS to propagate.

---

## Checklist

### One-Time Setup
- [ ] Created Cloudflare account
- [ ] Purchased/configured domain
- [ ] Enabled HTTPS in Cloudflare
- [ ] Created GCP account
- [ ] Created GCP project
- [ ] Enabled billing
- [ ] Installed gcloud CLI
- [ ] Authenticated with `gcloud auth login`
- [ ] Enabled required APIs
- [ ] Created Artifact Registry
- [ ] Created Cloud SQL instance
- [ ] Created database
- [ ] Authorized your IP for Cloud SQL
- [ ] Created `deployment/cloudrun/.env`
- [ ] Created GitHub service account
- [ ] Added secrets to GitHub
- [ ] Got Cloudflare API credentials

### Testing
- [ ] `make local` works
- [ ] `make cloudrun` deploys successfully
- [ ] Domain resolves with HTTPS

---

## Cost Summary

| Resource | Monthly Cost |
|----------|--------------|
| Cloud SQL (db-f1-micro) | ~$7-10 |
| Cloud Run (low traffic) | ~$0-5 |
| Artifact Registry | ~$0.10/GB |
| Cloudflare | Free |
| **Total** | **~$10-15/mo** |

---

## Troubleshooting

### "Permission denied" when running gcloud commands
- Run `gcloud auth login` again
- Make sure your project is set: `gcloud config set project YOUR_PROJECT_ID`

### Cloud Run: "Container failed to start"
- Check logs: `gcloud run logs read backend --region us-central1 --limit 50`
- Common causes: wrong environment variables, database connection issues

### Can't connect to Cloud SQL
- Check your IP is authorized: Go to GCP Console → SQL → your instance → Connections
- Make sure the password is correct in your .env

### Cloudflare: "Too many redirects"
- Make sure SSL mode is set to **Full (strict)**, not just "Full" or "Flexible"

### Docker build fails
- Make sure Docker Desktop is running
- Try `docker system prune` to clear old data

### DNS not working
- DNS changes can take up to 24 hours (usually minutes)
- Check Cloudflare dashboard to verify records are correct
- Try `nslookup yourdomain.com` to see current DNS
