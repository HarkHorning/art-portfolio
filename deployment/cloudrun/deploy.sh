#!/bin/bash
# Build, push, and deploy to Cloud Run.
# Run this every time you want to deploy a new version.
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
source "$SCRIPT_DIR/.env"

REGISTRY="${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT_ID}/portfolio"
BACKEND_IMAGE="${REGISTRY}/backend:latest"
FRONTEND_IMAGE="${REGISTRY}/frontend:latest"

echo "=== Deploying portfolio to Cloud Run ==="

# Authenticate Docker with Artifact Registry
echo "--- Configuring Docker auth..."
gcloud auth configure-docker "${GCP_REGION}-docker.pkg.dev" --quiet

# Build and push backend
echo "--- Building backend image..."
docker build -t "$BACKEND_IMAGE" "$REPO_ROOT/backend"
echo "--- Pushing backend image..."
docker push "$BACKEND_IMAGE"

# Build and push frontend
echo "--- Building frontend image..."
docker build -t "$FRONTEND_IMAGE" "$REPO_ROOT/frontend"
echo "--- Pushing frontend image..."
docker push "$FRONTEND_IMAGE"

# Deploy backend Cloud Run service
echo "--- Deploying backend..."
gcloud run deploy portfolio-backend \
  --image="$BACKEND_IMAGE" \
  --platform=managed \
  --region="$GCP_REGION" \
  --service-account="portfolio-backend@${GCP_PROJECT_ID}.iam.gserviceaccount.com" \
  --add-cloudsql-instances="$CLOUD_SQL_CONNECTION" \
  --set-env-vars="\
ENVIRONMENT=cloudrun,\
DB_HOST=/cloudsql/${CLOUD_SQL_CONNECTION},\
DB_USER=${DB_USER},\
DB_PASSWORD=${DB_PASSWORD},\
DB_NAME=${DB_NAME},\
DB_SEED_DATA=false,\
ADMIN_USERNAME=${ADMIN_USERNAME},\
ADMIN_PASSWORD=${ADMIN_PASSWORD},\
ADMIN_SESSION_SECRET=${ADMIN_SESSION_SECRET},\
GCS_BUCKET=${GCS_BUCKET}" \
  --allow-unauthenticated \
  --port=8080 \
  --project="$GCP_PROJECT_ID"

# Get backend URL
BACKEND_URL=$(gcloud run services describe portfolio-backend \
  --platform=managed \
  --region="$GCP_REGION" \
  --format="value(status.url)" \
  --project="$GCP_PROJECT_ID")
echo "Backend deployed at: $BACKEND_URL"

# Deploy frontend Cloud Run service
echo "--- Deploying frontend..."
gcloud run deploy portfolio-frontend \
  --image="$FRONTEND_IMAGE" \
  --platform=managed \
  --region="$GCP_REGION" \
  --set-env-vars="BACKEND_URL=${BACKEND_URL}" \
  --allow-unauthenticated \
  --port=3000 \
  --project="$GCP_PROJECT_ID"

# Get frontend URL
FRONTEND_URL=$(gcloud run services describe portfolio-frontend \
  --platform=managed \
  --region="$GCP_REGION" \
  --format="value(status.url)" \
  --project="$GCP_PROJECT_ID")
echo "Frontend deployed at: $FRONTEND_URL"

# Update backend CORS with the actual frontend URL
echo "--- Updating backend CORS with frontend URL..."
gcloud run services update portfolio-backend \
  --platform=managed \
  --region="$GCP_REGION" \
  --update-env-vars="CORS_ORIGIN=${FRONTEND_URL}" \
  --project="$GCP_PROJECT_ID"

echo ""
echo "=== Deploy complete ==="
echo ""
echo "Frontend: $FRONTEND_URL"
echo "Backend:  $BACKEND_URL"
echo ""
