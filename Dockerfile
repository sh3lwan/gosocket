FROM golang:latest AS builder

# Download mysql-client
RUN apt-get update && apt-get install -y default-mysql-client
# Download migrate go library
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install fresh for hot reloading
RUN go install github.com/air-verse/air@latest

# Set the working directory
WORKDIR /app

COPY go.mod .

RUN go mod download && go mod vendor && go mod tidy


# Copy the entrypoint script into the container
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

# Make the entrypoint script executable
RUN chmod +x /usr/local/bin/entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

COPY . .

FROM builder AS development

RUN go build -o ./tmp/main ./cmd/chatapp

#CMD [ "air" ]

# PRODUCTION ##
FROM builder AS production

# Compile Binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/chatapp ./cmd/chatapp

# Set the working directory to the directory containing the binary
WORKDIR /bin

# Run the binary directly
CMD ["./chatapp"]
