# VoltRun Backend

Go-based backend service for VoltRun cloud function platform.

## Features

- RESTful API built with Fiber
- JWT authentication
- PostgreSQL database with GORM
- Firecracker VM management
- Function execution engine
- Structured logging with Zap

## Project Structure

```
backend/
├── cmd/
│   └── server/           # Main application entry point
├── internal/
│   ├── api/              # HTTP handlers and routes
│   ├── auth/             # JWT authentication & middleware
│   ├── exec/             # Function execution engine
│   ├── storage/          # Database models and queries
│   ├── vm/               # Firecracker VM management
│   └── utils/            # Configuration and logging
├── go.mod
├── go.sum
├── Dockerfile
└── .env.example
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- (Optional) Firecracker for VM execution

### Installation

1. Clone the repository and navigate to backend:

```bash
cd backend
```

2. Install dependencies:

```bash
go mod download
```

3. Copy environment file:

```bash
cp .env.example .env
```

4. Start PostgreSQL (or use Docker):

```bash
docker run -d \
  --name voltrun-postgres \
  -e POSTGRES_USER=voltrun \
  -e POSTGRES_PASSWORD=voltrun \
  -e POSTGRES_DB=voltrun \
  -p 5432:5432 \
  postgres:15-alpine
```

5. Run the server:

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/refresh` - Refresh JWT token

### Functions

- `GET /api/functions` - List user's functions
- `POST /api/functions` - Create new function
- `GET /api/functions/:id` - Get function details
- `PUT /api/functions/:id` - Update function
- `DELETE /api/functions/:id` - Delete function
- `POST /api/functions/:id/execute` - Execute function

### Executions

- `GET /api/executions` - List executions
- `GET /api/executions/:id` - Get execution details
- `GET /api/executions/:id/logs` - Get execution logs

### API Keys

- `GET /api/keys` - List API keys
- `POST /api/keys` - Create API key
- `DELETE /api/keys/:id` - Delete API key

### Health Check

- `GET /health` - Service health status

## Development

### Run with hot reload:

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with air
air
```

### Run tests:

```bash
go test ./...
```

### Build for production:

```bash
go build -o voltrun cmd/server/main.go
```

### Docker build:

```bash
docker build -t voltrun-backend .
docker run -p 8080:8080 voltrun-backend
```

## Environment Variables

See `.env.example` for all available configuration options.

Key variables:

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret for signing JWT tokens
- `ENVIRONMENT` - Environment (development, production)

## Database Migrations

Migrations are handled automatically by GORM on startup. Models are defined in `internal/storage/models.go`.

## Architecture

### Request Flow

1. HTTP request → Fiber router
2. Authentication middleware validates JWT
3. Controller handles business logic
4. Storage layer interacts with database
5. Execution engine manages VM lifecycle
6. Response sent back to client

### Function Execution

1. User uploads function code
2. On execution request, create Firecracker VM
3. Inject code into VM
4. Execute with resource limits
5. Capture output and logs
6. Destroy VM
7. Return results to user

## Contributing

1. Create feature branch
2. Make changes
3. Run tests
4. Submit pull request

## License

MIT
