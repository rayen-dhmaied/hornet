# HorNet: Horizon Social Network

Microservices-based social network built with Go, featuring posts management and follower relationships. Deployed with Docker and Kubernetes using Istio service mesh.

## Architecture

- **Posts Service** - User posts management (MongoDB)
- **Followers Service** - Follower relationships and social graph (Neo4j)
- **Clean Architecture** - Handler → Service → Repository pattern
- **Service Mesh** - Istio with mTLS and Keycloak authentication

## Project Structure

```
.
├── api/                    # API layer
│   ├── followers/          # Followers service endpoints
│   │   ├── handler/        # HTTP handlers
│   │   ├── model/          # Data models
│   │   ├── repository/     # Database layer
│   │   ├── service/        # Business logic
│   │   └── router.go       # Route definitions
│   ├── posts/              # Posts service endpoints
│   │   └── (same structure)
│   └── openapi/            # OpenAPI specifications
├── cmd/                    # Service entry points
│   ├── followers/main.go   # Followers service
│   └── posts/main.go       # Posts service
├── common/                 # Shared utilities
│   └── logger/             # Logging package
├── config/                 # Configuration per service
├── Dockerfile              # Multi-stage build
└── Makefile                # Build automation
```

## Quick Start

### Prerequisites

- Go 1.23.2+
- Docker
- MongoDB (for Posts Service)
- Neo4j (for Followers Service)

### Build & Run

```bash
# Build service binary
make build SERVICE=posts

# Build Docker image
make build-container SERVICE=posts PORT=8080

# Run locally
make run SERVICE=posts

# Code quality
make lint
make fmt

# Clean build artifacts
make clean
```

### Docker Build

The Dockerfile uses multi-stage builds with distroless base images for minimal size (<10MB per service):

```bash
docker build --build-arg SERVICE=posts --build-arg PORT=8080 -t hornet-posts .
docker build --build-arg SERVICE=followers --build-arg PORT=8081 -t hornet-followers .
```

## Configuration

### Posts Service

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `POSTS_PORT` | HTTP server port | `8080` | No |
| `MONGO_URI` | MongoDB connection string | - | Yes |
| `MONGO_DB` | MongoDB database name | - | Yes |
| `FOLLOWERS_SERVICE_URL` | Followers service endpoint | - | Yes |

### Followers Service

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `FOLLOWERS_PORT` | HTTP server port | `8080` | No |
| `NEO4J_URI` | Neo4j connection string | - | Yes |
| `NEO4J_DB` | Neo4j database name | `neo4j` | Yes |
| `NEO4J_USER` | Neo4j username | - | Yes |
| `NEO4J_PASSWORD` | Neo4j password | - | Yes |
| `POSTS_SERVICE_URL` | Posts service endpoint | - | Yes |

### Example `.env`

```bash
# Posts Service
POSTS_PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DB=hornet
FOLLOWERS_SERVICE_URL=http://followers-service:8081

# Followers Service
FOLLOWERS_PORT=8081
NEO4J_URI=bolt://localhost:7687
NEO4J_DB=neo4j
NEO4J_USER=neo4j
NEO4J_PASSWORD=password
POSTS_SERVICE_URL=http://posts-service:8080
```

## Development

### Makefile Targets

```bash
make build SERVICE=<service>           # Build binary
make build-container SERVICE=<service> # Build Docker image
make run SERVICE=<service>             # Run service locally
make lint                              # Run golangci-lint
make fmt                               # Format code
make clean                             # Remove binaries
make help                              # Show all commands
```

### API Documentation

OpenAPI specifications available in `api/openapi/`:
- `posts.json` - Posts Service API spec
- `followers.json` - Followers Service API spec

## Deployment

### Kubernetes + Istio

Services are designed to run in Kubernetes with Istio service mesh:

```yaml
# Example deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: posts-service
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: posts
        image: hornet-posts:latest
        env:
        - name: MONGO_URI
          valueFrom:
            secretKeyRef:
              name: mongodb-secret
              key: uri
```

Istio handles:
- mTLS between services
- JWT authentication via Keycloak
- Load balancing and retries
- Traffic management
