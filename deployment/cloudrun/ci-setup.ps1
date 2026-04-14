# One-time setup: Workload Identity Federation for GitHub Actions.
# Run this once from PowerShell (gcloud must be authenticated).
# Usage: .\ci-setup.ps1

$ErrorActionPreference = "Stop"

$PROJECT  = "hark-portfolio"
$REGION   = "us-central1"
$REPO     = "HarkHorning/art-portfolio"
$SA_NAME  = "github-actions"
$SA_EMAIL = "${SA_NAME}@${PROJECT}.iam.gserviceaccount.com"
$POOL     = "github-pool"
$PROVIDER = "github-provider"

Write-Host "=== CI/CD: Workload Identity Federation setup ===" -ForegroundColor Cyan

# --- Service account ---
Write-Host "--- Creating service account..."
gcloud iam service-accounts create $SA_NAME `
    --display-name="GitHub Actions" `
    --project=$PROJECT

# --- IAM roles ---
Write-Host "--- Granting IAM roles..."

# Push images to Artifact Registry
gcloud projects add-iam-policy-binding $PROJECT `
    --member="serviceAccount:$SA_EMAIL" `
    --role="roles/artifactregistry.writer" `
    --quiet

# Deploy Cloud Run services
gcloud projects add-iam-policy-binding $PROJECT `
    --member="serviceAccount:$SA_EMAIL" `
    --role="roles/run.developer" `
    --quiet

# Read Cloud Run service URLs (needed to update CORS)
gcloud projects add-iam-policy-binding $PROJECT `
    --member="serviceAccount:$SA_EMAIL" `
    --role="roles/run.viewer" `
    --quiet

# Act as the backend's runtime service account during deploy
gcloud iam service-accounts add-iam-policy-binding `
    "portfolio-backend@${PROJECT}.iam.gserviceaccount.com" `
    --member="serviceAccount:$SA_EMAIL" `
    --role="roles/iam.serviceAccountUser" `
    --project=$PROJECT

# --- Workload Identity Pool ---
Write-Host "--- Creating Workload Identity Pool..."
gcloud iam workload-identity-pools create $POOL `
    --location="global" `
    --display-name="GitHub Actions Pool" `
    --project=$PROJECT

# --- OIDC Provider ---
Write-Host "--- Creating OIDC provider..."
gcloud iam workload-identity-pools providers create-oidc $PROVIDER `
    --location="global" `
    --workload-identity-pool=$POOL `
    --display-name="GitHub Provider" `
    --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" `
    --attribute-condition="assertion.repository=='$REPO'" `
    --issuer-uri="https://token.actions.githubusercontent.com" `
    --project=$PROJECT

# --- Allow GitHub repo to impersonate the SA ---
Write-Host "--- Binding GitHub repo to service account..."
$POOL_RESOURCE = gcloud iam workload-identity-pools describe $POOL `
    --location="global" `
    --project=$PROJECT `
    --format="value(name)"

gcloud iam service-accounts add-iam-policy-binding $SA_EMAIL `
    --role="roles/iam.workloadIdentityUser" `
    --member="principalSet://iam.googleapis.com/${POOL_RESOURCE}/attribute.repository/${REPO}" `
    --project=$PROJECT

# --- Print values needed as GitHub secrets ---
$PROVIDER_RESOURCE = gcloud iam workload-identity-pools providers describe $PROVIDER `
    --location="global" `
    --workload-identity-pool=$POOL `
    --project=$PROJECT `
    --format="value(name)"

Write-Host ""
Write-Host "=== Done. Add these as GitHub Actions secrets ===" -ForegroundColor Green
Write-Host ""
Write-Host "GCP_WIF_PROVIDER     = $PROVIDER_RESOURCE" -ForegroundColor Yellow
Write-Host "GCP_SERVICE_ACCOUNT  = $SA_EMAIL" -ForegroundColor Yellow
Write-Host ""
Write-Host "Also add these secrets (same values as your .env file):" -ForegroundColor Cyan
Write-Host "  GCP_PROJECT_ID, GCP_REGION, CLOUD_SQL_CONNECTION"
Write-Host "  DB_USER, DB_PASSWORD, DB_NAME"
Write-Host "  ADMIN_USERNAME, ADMIN_PASSWORD, ADMIN_SESSION_SECRET"
Write-Host "  GCS_BUCKET"
Write-Host ""
