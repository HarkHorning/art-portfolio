#!/bin/bash
# One-time GCP infrastructure setup.
# Run this once before your first deploy.
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/.env"

echo "=== Setting up GCP infrastructure for project: $GCP_PROJECT_ID ==="

# Enable required APIs
echo "--- Enabling APIs..."
gcloud services enable \
  run.googleapis.com \
  sqladmin.googleapis.com \
  artifactregistry.googleapis.com \
  --project="$GCP_PROJECT_ID"

# Create Artifact Registry repository
echo "--- Creating Artifact Registry repository..."
gcloud artifacts repositories create portfolio \
  --repository-format=docker \
  --location="$GCP_REGION" \
  --description="Portfolio Docker images" \
  --project="$GCP_PROJECT_ID" 2>/dev/null || echo "Repository already exists, skipping."

# Create Cloud SQL instance (takes ~5 minutes)
echo "--- Creating Cloud SQL instance (this takes ~5 minutes)..."
gcloud sql instances create portfolio-db \
  --database-version=MYSQL_8_0 \
  --tier=db-f1-micro \
  --region="$GCP_REGION" \
  --storage-type=SSD \
  --storage-size=10GB \
  --project="$GCP_PROJECT_ID" 2>/dev/null || echo "Instance already exists, skipping."

# Create portfolio database
echo "--- Creating database..."
gcloud sql databases create portfolio \
  --instance=portfolio-db \
  --project="$GCP_PROJECT_ID" 2>/dev/null || echo "Database already exists, skipping."

# Set root password
echo "--- Setting database password..."
gcloud sql users set-password root \
  --host=% \
  --instance=portfolio-db \
  --password="$DB_PASSWORD" \
  --project="$GCP_PROJECT_ID"

# Create service account for the backend
echo "--- Creating service account..."
gcloud iam service-accounts create portfolio-backend \
  --display-name="Portfolio Backend" \
  --project="$GCP_PROJECT_ID" 2>/dev/null || echo "Service account already exists, skipping."

# Grant Cloud SQL access to the service account
echo "--- Granting Cloud SQL access..."
gcloud projects add-iam-policy-binding "$GCP_PROJECT_ID" \
  --member="serviceAccount:portfolio-backend@${GCP_PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"

# Grant GCS access to the service account (for image uploads)
echo "--- Granting GCS access..."
gcloud storage buckets add-iam-policy-binding "gs://${GCS_BUCKET}" \
  --member="serviceAccount:portfolio-backend@${GCP_PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/storage.objectAdmin"

# Print the Cloud SQL connection name to put in .env
echo ""
echo "=== Setup complete ==="
echo ""
echo "Add this to your .env file:"
gcloud sql instances describe portfolio-db \
  --format="value(connectionName)" \
  --project="$GCP_PROJECT_ID" | xargs -I{} echo "CLOUD_SQL_CONNECTION={}"

