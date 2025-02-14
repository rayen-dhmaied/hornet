# HorNet: Horizon Social Network

HorNet is a microservices-based social network application designed to facilitate user posts, and followers management. Each service is built in Go and deployed with Docker, making it scalable and flexible.

## Project Structure

```plaintext
.
├── api
│   ├── followers
│   │   ├── handler
│   │   │   └── handler.go
│   │   ├── model
│   │   │   └── model.go
│   │   ├── repository
│   │   │   └── repository.go
│   │   ├── router.go
│   │   └── service
│   │       └── service.go
│   ├── openapi
│   │   ├── followers.json
│   │   └── posts.json
│   └── posts
│       ├── handler
│       │   └── handler.go
│       ├── model
│       │   └── model.go
│       ├── repository
│       │   └── repository.go
│       ├── router.go
│       └── service
│           └── service.go
├── cmd
│   ├── followers
│   │   └── main.go
│   └── posts
│       └── main.go
├── common
│   └── logger
│       └── logger.go
├── config
│   ├── followers
│   │   └── config.go
│   └── posts
│       └── config.go
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Prerequisites

- **Go**: Make sure you have Go installed (version 1.23.2 or later).
- **Docker**: Required to build and run Docker containers.
- **GolangCI-Lint**: Recommended for linting the Go code.

## Usage

### Build the Binary

```sh
make build SERVICE=<service_name>
```

This command builds the binary for the specified service. Replace `<service_name>` with `posts`, or `followers`. By default, `SERVICE` is set to `posts`.

### Build the Docker Container

```sh
make build-container SERVICE=<service_name> PORT=<service_port>
```

Builds a Docker image for the specified service. Make sure Docker is running.

### Run the Binary Directly

```sh
make run SERVICE=<service_name>
```

Runs the binary directly, which is useful for local testing outside of Docker.

### Lint the Code

```sh
make lint
```

Checks the code for linting errors. This requires `golangci-lint` to be installed.

### Format the Code

```sh
make fmt
```

Formats the Go code according to Go standards.

### Clean up

```sh
make clean
```

Removes any binaries generated during the build process.

### Help

```sh
make help
```

Displays help information for all available `make` commands and variables.


## Environment Variables


| Variable Name        | Description                          | Default Value | Required |
|----------------------|--------------------------------------|---------------|----------|
| `POSTS_PORT`               | The port on which the posts server runs    | `8080`        | No       |
| `FOLLOWERS_PORT`               | The port on which the followers server runs    | `8080`        | No       |
| `MONGO_URI`       | URI for the MongoDB database      | -             | Yes      |
| `MONGO_DB`       | The name the MongoDB database      | -             | Yes      |
| `NEO4J_URI`           | URI for the Neo4j database | - | Yes       |
| `NEO4J_DB`           | The name of Neo4j database | `neo4j` | Yes       |
| `NEO4J_USER`           | User for the Neo4j database | - | Yes       |
| `NEO4J_PASSWORD`           | The password of Neo4j database | - | Yes       |