FROM golang:latest

# Download mysql-client & migrate go library
RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash && apt-get update && apt-get install -y default-mysql-client && apt-get install -y migrate && rm -rf /var/lib/apt/lists/*

# Copy the entrypoint script into the container
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

# Make the entrypoint script executable
RUN chmod +x /usr/local/bin/entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify && go mod tidy

COPY . .

RUN go build -o bin/chatapp

EXPOSE 8080

CMD [ "./bin/chatapp" ]
