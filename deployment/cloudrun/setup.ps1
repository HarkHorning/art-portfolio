# One-time GCP infrastructure setup.
# Run this once before your first deploy.
# Usage: .\setup.ps1

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$EnvFile = Join-Path $ScriptDir ".env"

if (-not (Test-Path $EnvFile)) {
    Write-Error ".env file not found. Copy env.template to .env and fill it in."
    exit 1
}

# Load .env
Get-Content $EnvFile | ForEach-Object {
    if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
        [System.Environment]::SetEnvironmentVariable($matches[1].Trim(), $matches[2].Trim())
    }
}

$PROJECT  = $env:GCP_PROJECT_ID
$REGION   = $env:GCP_REGION
$DB_PASS  = $env:DB_PASSWORD
$BUCKET   = $env:GCS_BUCKET
$SA       = "portfolio-backend@${PROJECT}.iam.gserviceaccount.com"

Write-Host "=== Setting up GCP infrastructure for project: $PROJECT ===" -ForegroundColor Cyan

Write-Host "--- Enabling APIs..."
gcloud services enable run.googleapis.com sqladmin.googleapis.com artifactregistry.googleapis.com --project=$PROJECT

Write-Host "--- Creating Artifact Registry repository..."
gcloud artifacts repositories create portfolio `
    --repository-format=docker `
    --location=$REGION `
    --description="Portfolio Docker images" `
    --project=$PROJECT 2>$null
if ($LASTEXITCODE -ne 0) { Write-Host "Repository already exists, skipping." }

Write-Host "--- Creating Cloud SQL instance (this takes ~5 minutes)..."
gcloud sql instances create portfolio-db `
    --database-version=MYSQL_8_0 `
    --tier=db-f1-micro `
    --region=$REGION `
    --storage-type=SSD `
    --storage-size=10GB `
    --project=$PROJECT 2>$null
if ($LASTEXITCODE -ne 0) { Write-Host "Instance already exists, skipping." }

Write-Host "--- Creating database..."
gcloud sql databases create portfolio `
    --instance=portfolio-db `
    --project=$PROJECT 2>$null
if ($LASTEXITCODE -ne 0) { Write-Host "Database already exists, skipping." }

Write-Host "--- Setting database password..."
gcloud sql users set-password root `
    --host=% `
    --instance=portfolio-db `
    --password=$DB_PASS `
    --project=$PROJECT

Write-Host "--- Creating service account..."
gcloud iam service-accounts create portfolio-backend `
    --display-name="Portfolio Backend" `
    --project=$PROJECT 2>$null
if ($LASTEXITCODE -ne 0) { Write-Host "Service account already exists, skipping." }

Write-Host "--- Granting Cloud SQL access..."
gcloud projects add-iam-policy-binding $PROJECT `
    --member="serviceAccount:$SA" `
    --role="roles/cloudsql.client"

Write-Host "--- Granting GCS access..."
gcloud storage buckets add-iam-policy-binding "gs://$BUCKET" `
    --member="serviceAccount:$SA" `
    --role="roles/storage.objectAdmin"

Write-Host ""
Write-Host "=== Setup complete ===" -ForegroundColor Green
Write-Host ""
Write-Host "Add this to your .env file:" -ForegroundColor Yellow
$CONNECTION = gcloud sql instances describe portfolio-db --format="value(connectionName)" --project=$PROJECT
Write-Host "CLOUD_SQL_CONNECTION=$CONNECTION" -ForegroundColor Yellow
Write-Host ""
