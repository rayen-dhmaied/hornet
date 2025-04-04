FROM golang:1.23.2 AS builder
ARG SERVICE

# Set the working directory
WORKDIR /app

# Copy all files
COPY . .

# Set environment variable to disable CGO (C Go) support
ENV CGO_ENABLED=0

# Download dependencies specified in the go.mod file
RUN go mod download

# Build the Go application
RUN go build -o app -ldflags="-s -w" ./cmd/${SERVICE}/main.go

FROM gcr.io/distroless/static-debian12 AS runtime
ARG PORT

# Copy the compiled binary from the builder stage to the runtime image
COPY --from=builder ./app ./

# Expose port
EXPOSE ${PORT}/tcp

# Set environment variable for Gin
ENV GIN_MODE=release

# Set the command to run the application
ENTRYPOINT ["./app"]
