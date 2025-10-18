# Mock Backend Jobs

A lightweight Golang service using Echo framework that simulates a job scheduling system with configurable execution time and failure scenarios.

## Features

- **3 Hardcoded Jobs**:
  - `user_annonymization`
  - `zendesk_import`
  - `blueshift_export`

- **Job States**: PENDING, RUNNING, SUCCESS, FAILED
- **In-Memory Storage**: Job status stored in memory with UUID tracking
- **Configurable Jobs**: Configure execution time and failure behavior via environment variables
- **Background Execution**: Jobs run asynchronously
- **Health Check Endpoint**: Monitor service health

## Endpoints

### Health Check
```
GET /health
```
Returns service health and available jobs.

**Response:**
```json
{
  "status": "healthy",
  "service": "mock-backend-jobs",
  "available_jobs": ["user_annonymization", "zendesk_import", "blueshift_export"]
}
```

### Trigger Job
```
POST /api/v1/system/scheduler/trigger
```

**Request Body:**
```json
{
  "job_name": "user_annonymization"
}
```

**Response:**
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Job triggered successfully",
  "status": "PENDING"
}
```

### Get Job Status
```
GET /api/v1/system/scheduler/status?job_id=<job_id>
```

**Response:**
```json
{
  "job": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "user_annonymization",
    "status": "SUCCESS",
    "started_at": "2025-10-18T10:00:00Z",
    "finished_at": "2025-10-18T10:00:05Z"
  }
}
```

**Get All Jobs** (omit job_id parameter):
```
GET /api/v1/system/scheduler/status
```

## CI/CD with GitHub Actions

This project includes automated GitHub Actions workflows for building and deploying Docker images to GitHub Container Registry (GHCR).

### Workflows

1. **Test Workflow** (`.github/workflows/test.yml`)
   - Runs on every push and pull request
   - Executes `go vet`, `go fmt`, and tests
   - Verifies build completes successfully

2. **Docker Build & Push** (`.github/workflows/docker-build-push.yml`)
   - Builds multi-arch Docker images (amd64, arm64)
   - Pushes to `ghcr.io` with automatic tagging
   - **Triggered only on version tags** (e.g., `v1.0.0`)

### Image Tags

Images are automatically tagged when you push a version tag:
- **Semver tags**: `ghcr.io/your-username/mock-backend-jobs:v1.0.0`, `v1.0`, `v1`
- **Latest**: `ghcr.io/your-username/mock-backend-jobs:latest`

### Using GHCR Images

**Pull the image:**
```bash
docker pull ghcr.io/your-username/mock-backend-jobs:latest
```

**Run the image:**
```bash
docker run -p 8080:8080 ghcr.io/your-username/mock-backend-jobs:latest
```

**Deploy with Helm:**
```bash
helm install mock-backend-jobs ./helm/mock-backend-jobs \
  --set image.repository=ghcr.io/your-username/mock-backend-jobs \
  --set image.tag=latest
```

### Creating a Release

To create a new release with semantic versioning:

```bash
# Tag the release
git tag -a v1.0.0 -m "Release version 1.0.0"

# Push the tag
git push origin v1.0.0
```

GitHub Actions will automatically build and push images with tags:
- `ghcr.io/your-username/mock-backend-jobs:v1.0.0`
- `ghcr.io/your-username/mock-backend-jobs:v1.0`
- `ghcr.io/your-username/mock-backend-jobs:v1`
- `ghcr.io/your-username/mock-backend-jobs:latest`

### GHCR Authentication

For private repositories, authenticate with GitHub Container Registry:

```bash
# Create a GitHub Personal Access Token with read:packages scope
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

For Kubernetes deployments, create an image pull secret:

```bash
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=USERNAME \
  --docker-password=$GITHUB_TOKEN \
  --docker-email=your-email@example.com
```

Then reference it in `values.yaml`:
```yaml
imagePullSecrets:
  - name: ghcr-secret
```

## Quick Start

### Local Development

1. **Clone and navigate to the project:**
```bash
cd /path/to/mock-backend-jobs
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Run the service:**
```bash
go run cmd/server/main.go
```

The service will start on `http://localhost:8080`

### Using Docker

1. **Build the image:**
```bash
docker build -t mock-backend-jobs:latest .
```

2. **Run the container:**
```bash
docker run -p 8080:8080 mock-backend-jobs:latest
```

### Using Docker with Custom Configuration

```bash
docker run -p 8080:8080 \
  -e USER_ANNONYMIZATION_DURATION=10s \
  -e USER_ANNONYMIZATION_SHOULD_FAIL=false \
  -e ZENDESK_IMPORT_DURATION=15s \
  -e ZENDESK_IMPORT_SHOULD_FAIL=false \
  -e BLUESHIFT_EXPORT_DURATION=20s \
  -e BLUESHIFT_EXPORT_SHOULD_FAIL=true \
  mock-backend-jobs:latest
```

## Configuration

Configure job behavior using environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `USER_ANNONYMIZATION_DURATION` | Execution time for user_annonymization | `10s` |
| `USER_ANNONYMIZATION_SHOULD_FAIL` | Force failure for user_annonymization | `true` or `false` |
| `ZENDESK_IMPORT_DURATION` | Execution time for zendesk_import | `15s` |
| `ZENDESK_IMPORT_SHOULD_FAIL` | Force failure for zendesk_import | `true` or `false` |
| `BLUESHIFT_EXPORT_DURATION` | Execution time for blueshift_export | `20s` |
| `BLUESHIFT_EXPORT_SHOULD_FAIL` | Force failure for blueshift_export | `true` or `false` |

### Default Configuration

- `user_annonymization`: 5s execution, success
- `zendesk_import`: 8s execution, success
- `blueshift_export`: 10s execution, **fails by default**

Note: Jobs have a 10% random failure chance even when not configured to fail.

## Kubernetes Deployment

### Using Helm

1. **Install the chart:**
```bash
helm install mock-backend-jobs ./helm/mock-backend-jobs
```

2. **Install with custom values:**
```bash
helm install mock-backend-jobs ./helm/mock-backend-jobs \
  --set image.repository=your-registry/mock-backend-jobs \
  --set image.tag=v1.0.0
```

3. **Install with custom job configuration:**
```bash
helm install mock-backend-jobs ./helm/mock-backend-jobs \
  --set-string 'env[0].name=USER_ANNONYMIZATION_DURATION' \
  --set-string 'env[0].value=10s' \
  --set-string 'env[1].name=USER_ANNONYMIZATION_SHOULD_FAIL' \
  --set-string 'env[1].value=false'
```

4. **Upgrade the deployment:**
```bash
helm upgrade mock-backend-jobs ./helm/mock-backend-jobs
```

5. **Uninstall:**
```bash
helm uninstall mock-backend-jobs
```

### Enable Ingress

Edit `helm/mock-backend-jobs/values.yaml`:

```yaml
ingress:
  enabled: true
  className: "nginx"
  hosts:
    - host: mock-backend-jobs.yourdomain.com
      paths:
        - path: /
          pathType: Prefix
```

## Testing the API

### Trigger a Job

```bash
curl -X POST http://localhost:8080/api/v1/system/scheduler/trigger \
  -H "Content-Type: application/json" \
  -d '{"job_name": "user_annonymization"}'
```

### Check Job Status

```bash
# Get specific job
curl http://localhost:8080/api/v1/system/scheduler/status?job_id=<job_id>

# Get all jobs
curl http://localhost:8080/api/v1/system/scheduler/status
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Example Workflow

```bash
# 1. Start the service
go run cmd/server/main.go

# 2. Trigger a job
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/system/scheduler/trigger \
  -H "Content-Type: application/json" \
  -d '{"job_name": "user_annonymization"}')

echo $RESPONSE
# {"job_id":"123e4567-e89b-12d3-a456-426614174000","message":"Job triggered successfully","status":"PENDING"}

# 3. Extract job ID and check status
JOB_ID=$(echo $RESPONSE | jq -r '.job_id')
curl http://localhost:8080/api/v1/system/scheduler/status?job_id=$JOB_ID

# 4. Wait and check again
sleep 6
curl http://localhost:8080/api/v1/system/scheduler/status?job_id=$JOB_ID
```

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── handlers/
│   │   └── scheduler.go      # HTTP handlers
│   ├── models/
│   │   └── job.go            # Data models
│   └── scheduler/
│       └── manager.go        # Job management logic
├── helm/
│   └── mock-backend-jobs/    # Helm chart
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
├── Dockerfile                # Container image definition
├── go.mod                    # Go module dependencies
└── README.md
```

## Development

### Build

```bash
go build -o bin/server cmd/server/main.go
```

### Run Tests (when implemented)

```bash
go test ./...
```

### Format Code

```bash
go fmt ./...
```

## License

MIT
