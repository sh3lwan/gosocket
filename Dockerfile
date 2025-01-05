FROM golang:latest

# Download mysql-client & migrate go library
RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash && apt-get update && apt-get install -y default-mysql-client && apt-get install -y migrate && rm -rf /var/lib/apt/lists/*

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o /chatapp ./cmd/chatapp

EXPOSE 8000

# Run
CMD ["/chatapp"]
