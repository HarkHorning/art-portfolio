# Build, push, and deploy to Cloud Run.
# Run this every time you want to deploy a new version.
# Usage: .\deploy.ps1

$ErrorActionPreference = "Stop"

$ScriptDir  = Split-Path -Parent $MyInvocation.MyCommand.Path
$RepoRoot   = Split-Path -Parent (Split-Path -Parent $ScriptDir)
$EnvFile    = Join-Path $ScriptDir ".env"

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

$PROJECT    = $env:GCP_PROJECT_ID
$REGION     = $env:GCP_REGION
$REGISTRY   = "${REGION}-docker.pkg.dev/${PROJECT}/portfolio"
$BACKEND_IMAGE  = "${REGISTRY}/backend:latest"
$FRONTEND_IMAGE = "${REGISTRY}/frontend:latest"

$DB_HOST    = "/cloudsql/$($env:CLOUD_SQL_CONNECTION)"
$CONN       = $env:CLOUD_SQL_CONNECTION

if (-not $CONN) {
    Write-Error "CLOUD_SQL_CONNECTION is empty in .env file."
    exit 1
}

Write-Host "=== Deploying portfolio to Cloud Run ===" -ForegroundColor Cyan

Write-Host "--- Configuring Docker auth..."
gcloud auth configure-docker "${REGION}-docker.pkg.dev" --quiet

Write-Host "--- Building backend image..."
docker build -t $BACKEND_IMAGE "$RepoRoot\backend"

Write-Host "--- Pushing backend image..."
docker push $BACKEND_IMAGE

Write-Host "--- Building frontend image..."
docker build -t $FRONTEND_IMAGE "$RepoRoot\frontend"

Write-Host "--- Pushing frontend image..."
docker push $FRONTEND_IMAGE

Write-Host "--- Deploying backend..."
gcloud run deploy portfolio-backend `
    --image=$BACKEND_IMAGE `
    --platform=managed `
    --region=$REGION `
    --service-account="portfolio-backend@${PROJECT}.iam.gserviceaccount.com" `
    --add-cloudsql-instances=$CONN `
    --set-env-vars="ENVIRONMENT=cloudrun,DB_HOST=$DB_HOST,DB_USER=$($env:DB_USER),DB_PASSWORD=$($env:DB_PASSWORD),DB_NAME=$($env:DB_NAME),DB_SEED_DATA=false,ADMIN_USERNAME=$($env:ADMIN_USERNAME),ADMIN_PASSWORD=$($env:ADMIN_PASSWORD),ADMIN_SESSION_SECRET=$($env:ADMIN_SESSION_SECRET),GCS_BUCKET=$($env:GCS_BUCKET)" `
    --allow-unauthenticated `
    --port=8080 `
    --project=$PROJECT

$BACKEND_URL = gcloud run services describe portfolio-backend `
    --platform=managed `
    --region=$REGION `
    --format="value(status.url)" `
    --project=$PROJECT

if (-not $BACKEND_URL) {
    Write-Error "Backend URL is empty - backend deploy failed. Check logs before deploying frontend."
    exit 1
}
Write-Host "Backend deployed at: $BACKEND_URL"

Write-Host "--- Deploying frontend..."
gcloud run deploy portfolio-frontend `
    --image=$FRONTEND_IMAGE `
    --platform=managed `
    --region=$REGION `
    --set-env-vars="BACKEND_URL=$BACKEND_URL" `
    --allow-unauthenticated `
    --port=3000 `
    --project=$PROJECT

$FRONTEND_URL = gcloud run services describe portfolio-frontend `
    --platform=managed `
    --region=$REGION `
    --format="value(status.url)" `
    --project=$PROJECT
Write-Host "Frontend deployed at: $FRONTEND_URL"

Write-Host "--- Updating backend CORS with frontend URL..."
gcloud run services update portfolio-backend `
    --platform=managed `
    --region=$REGION `
    --update-env-vars="CORS_ORIGIN=$FRONTEND_URL" `
    --project=$PROJECT

Write-Host ""
Write-Host "=== Deploy complete ===" -ForegroundColor Green
Write-Host ""
Write-Host "Frontend: $FRONTEND_URL" -ForegroundColor Yellow
Write-Host "Backend:  $BACKEND_URL"  -ForegroundColor Yellow
Write-Host ""
