FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify && go mod tidy

COPY . .

RUN go build -o bin/chatapp

EXPOSE 8080

CMD [ "./bin/chatapp" ]
