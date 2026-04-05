# Portfolio

My personal portfolio site, built to learn cloud infrastructure and deployment pipelines.

## What's in here

- **Frontend** - SvelteKit app with an art gallery and visit tracker
- **Backend** - Go API (Gin) that serves art data and logs visits
- **MySQL** - Stores everything

The whole thing runs on Azure Kubernetes Service. Infrastructure is managed with Terraform, and GitHub Actions handles CI/CD.

## Running locally

```bash
cd deployment/docker
docker compose up --build
```

Frontend runs on `localhost:3000`, backend on `localhost:8080`.

## Deploying

Infrastructure first:
```bash
cd deployment/terraform
terraform init && terraform apply
```

Then push to Kubernetes:
```bash
az aks get-credentials --resource-group Portfolio --name hark-portfolio-aks
kubectl apply -f deployment/kubernetes/
```

Pushes to `main` auto-deploy via GitHub Actions.

## Project layout

```
frontend/              # SvelteKit
backend/               # Go API
deployment/
  ├── docker/          # Local dev with Docker Compose
  ├── kubernetes/      # K8s manifests
  └── terraform/       # Azure infrastructure
```
