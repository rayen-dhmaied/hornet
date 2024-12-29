# HorNet: Horizon Social Network

HorNet is a microservices-based social network application designed to facilitate user posts, reactions, and connectionns management. Each service is built in Go and deployed with Docker, making it scalable and flexible.

## Project Structure

```plaintext
.
├── api
│   ├── connections # Connections ms implementation
│   ├── openapi     # Ms API specifications
│   │   ├── connections.json
│   │   ├── posts.json
│   │   └── reactions.json
│   ├── posts       # Posts ms implementation
│   └── reactions   # Reactions ms implementation
├── cmd
│   ├── connections
│   │   └── main.go
│   ├── posts
│   │   └── main.go
│   └── reactions
│       └── main.go
├── common
│   └── logger
│       └── logger.go
├── config
│   ├── connections
│   │   └── config.go
│   ├── posts
│   │   └── config.go
│   └── reactions
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

This command builds the binary for the specified service. Replace `<service_name>` with `posts`, `reactions`, or `connections`. By default, `SERVICE` is set to `posts`.

### Build the Docker Container

```sh
make build-container SERVICE=<service_name>
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