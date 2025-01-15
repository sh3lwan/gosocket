FROM golang:latest

# Install dependencies for MySQL client and migrate (via precompiled binary)
RUN apt-get update && \
    apt-get install -y default-mysql-client curl && \
    curl -Ls https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xz -C /usr/local/bin && \
    rm -rf /var/lib/apt/lists/*

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

RUN go mod download

# Copy the entrypoint script into the container
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

# Make the entrypoint script executable
RUN chmod +x /usr/local/bin/entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o /chatapp ./cmd/chatapp

EXPOSE 80

# Run
CMD ["/chatapp"]
