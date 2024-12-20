ARG SERVICE
RUN if [ -z "${SERVICE}" ] || { [ "${SERVICE}" != "reactions" ] && [ "${SERVICE}" != "posts" ] && [ "${SERVICE}" != "users" ] } ; then \
      echo "Error: SERVICE must be 'reactions', 'posts' or 'users'"; exit 1; \
    fi

FROM golang:1.23.2 AS builder

# Set the working directory
WORKDIR /app

# Copy all files
COPY . .

# Set environment variable to disable CGO (C Go) support
ENV CGO_ENABLED=0

# Download dependencies specified in the go.mod file
RUN go get -d -v .

# Build the Go application
RUN go build -a -installsuffix cgo -o app ./cmd/${SERVICE}/main.go

FROM scratch AS runtime

# Copy the compiled binary from the builder stage to the runtime image
COPY --from=builder ./app ./

# Expose port 8080
EXPOSE 8080/tcp

# Set environment variable for Gin
ENV GIN_MODE=release

# Set the command to run the application
ENTRYPOINT ["./app"]
