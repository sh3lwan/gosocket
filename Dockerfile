FROM golang:latest

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
