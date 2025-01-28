# GoSocket Backend Application

## Overview
**GoSocket** is a backend application built with Go for managing real-time communication and backend services. This application serves as the core for handling API requests, WebSocket connections, and business logic for your project.

## Features
- **WebSocket Support**: Efficient real-time communication using WebSockets.
- **RESTful APIs**: Endpoints for managing resources.
- **Dockerized Deployment**: Easily deployable with Docker.
- **Scalable Architecture**: Designed for scalability and high performance.

---

## Getting Started

### Prerequisites
Ensure you have the following installed:
- [Go](https://golang.org/doc/install) (v1.20 or later)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/)

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/sh3lwan/gosocket.git
   cd gosocket
   ```

2. Build and run application using Docker:
   ```bash
   go build -o gosocket ./cmd/chatapp/main.go
   ```

3. Run application locally:
   ```bash
   ./gosocket
   ```
---

## Configuration

### Environment Variables
Create a `.env` file in the project root with the following variables:

```env
APP_PORT=8080
DB_DATABASE=chat
DB_USERNAME=root
DB_PASSWORD=password
DB_HOST=mysql_db
DB_PORT=3306
DB_DRIVER=mysql
```

### Docker Setup
1. Build and start the services:
   ```bash
   docker-compose up --build
   ```

2. Access the application:
   - REST API: `http://localhost:8000/api`
   - WebSocket: `ws://localhost:8000/ws`

---

## API Endpoints

### RESTful APIs
| Method | Endpoint         | Description                |
|--------|------------------|----------------------------|
| GET    | `/api/health`    | Health check endpoint.     |
| POST   | `/api/login`     | Authenticate a user.       |

### WebSocket
- **Endpoint**: `/ws`
- Handles real-time communication for clients.

---

## Testing
Run unit tests using:
```bash
go test ./...
```

---

## Contributing
We welcome contributions! To contribute:
1. Fork the repository.
2. Create a new branch:
   ```bash
   git checkout -b feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add your message"
   ```
4. Push to the branch:
   ```bash
   git push origin feature-name
   ```
5. Create a pull request.

---

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Acknowledgments
- Built with [Go](https://golang.org/).
- Inspired by modern backend development best practices.

