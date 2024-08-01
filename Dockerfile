FROM golang:latest AS builder

# Download mysql-client & migrate go library
RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash && apt-get update && apt-get install -y default-mysql-client && apt-get install -y migrate && rm -rf /var/lib/apt/lists/* && curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

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

#RUN go build -buildvsc=false -o /app/tmp/main .
# Install fresh for hot reloading
#RUN go install github.com/pilu/fresh@latest
#RUN go install github.com/air-verse/air@latest
# Run the binary

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/chatapp .
# Run the binary
#CMD [ "./bin/chatapp" ]
CMD [ "air" ]

## PRODUCTION ##
FROM builder AS production

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/chatapp .
# Run the binary
CMD [ "./bin/chatapp" ]
