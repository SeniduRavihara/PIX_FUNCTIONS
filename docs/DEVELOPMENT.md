# VoltRun Development Guide

## Quick Start

### Backend

```bash
cd backend
go mod download
go run cmd/server/main.go
```

Backend will start on http://localhost:8080

### Frontend

```bash
cd frontend
npm install
cp .env.local.example .env.local
npm run dev
```

Frontend will start on http://localhost:3000

### Docker

```bash
cd deploy
docker-compose up -d
```

## Project Status

### ✅ Completed

- Backend folder structure with Go modules
- API routes with Fiber framework
- PostgreSQL database models (User, Function, Execution, APIKey)
- JWT authentication and middleware
- VM manager stub for Firecracker integration
- Execution engine with placeholder implementations
- Structured logging with Zap
- Configuration management
- Docker setup (docker-compose, Dockerfile)
- Frontend Next.js structure
- Authentication pages (login/register)
- Dashboard layout with navigation
- API client library
- Auth context provider

### 🚧 In Progress

- Monaco editor integration
- Function creation UI
- Execution interface

### 📋 Todo

- Actual Firecracker VM implementation
- Node.js and Python runners
- WebSocket for real-time logs
- API key generation and management
- Function versioning
- Usage tracking and billing hooks
- CI/CD pipelines
- Production deployment scripts

## Architecture Overview

```
┌─────────────┐
│   Frontend  │  Next.js Dashboard
│   (Port 3000)│  - Monaco Editor
└──────┬──────┘  - Function Management
       │
       │ HTTP/WebSocket
       ▼
┌─────────────┐
│   Backend   │  Go API Server
│   (Port 8080)│  - Routes & Auth
└──────┬──────┘  - Storage Layer
       │
       ├────────┐
       │        │
       ▼        ▼
┌──────────┐ ┌─────────┐
│PostgreSQL│ │  VMs    │  Firecracker MicroVMs
│          │ │         │  - Isolated Execution
└──────────┘ └─────────┘  - Multiple Runtimes
```

## Development Workflow

1. **Create Feature Branch**

   ```bash
   git checkout -b feature/your-feature
   ```

2. **Make Changes**

   - Backend changes in `backend/`
   - Frontend changes in `frontend/`

3. **Test Locally**

   - Run backend: `cd backend && go run cmd/server/main.go`
   - Run frontend: `cd frontend && npm run dev`

4. **Commit and Push**
   ```bash
   git add .
   git commit -m "feat: your feature description"
   git push origin feature/your-feature
   ```

## API Testing

Using curl:

```bash
# Health check
curl http://localhost:8080/health

# Login (placeholder - implement in backend first)
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# List functions (requires JWT token)
curl http://localhost:8080/api/functions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Database Setup

1. Start PostgreSQL:

```bash
docker run -d \
  --name voltrun-postgres \
  -e POSTGRES_USER=voltrun \
  -e POSTGRES_PASSWORD=voltrun \
  -e POSTGRES_DB=voltrun \
  -p 5432:5432 \
  postgres:15-alpine
```

2. Migrations run automatically on backend startup

3. Connect to database:

```bash
psql -h localhost -U voltrun -d voltrun
# Password: voltrun
```

## Firecracker Setup (Linux Only)

VoltRun uses Firecracker for VM isolation. This is currently stubbed out for development.

To enable Firecracker:

1. Download Firecracker:

```bash
curl -Lo firecracker https://github.com/firecracker-microvm/firecracker/releases/download/v1.5.0/firecracker-v1.5.0-x86_64
chmod +x firecracker
sudo mv firecracker /usr/bin/
```

2. Configure kernel and rootfs paths in `.env`

3. Implement actual VM creation in `backend/internal/vm/manager.go`

## Next Steps

1. **Implement Authentication Handlers**

   - Add user registration logic in `backend/internal/api/routes.go`
   - Hash passwords with bcrypt
   - Generate JWT tokens

2. **Build Function Management**

   - Create function CRUD operations
   - Store code in database
   - Implement function listing

3. **Add Monaco Editor**

   - Install `@monaco-editor/react`
   - Create code editor component
   - Add syntax highlighting for JS/Python

4. **Implement Execution**

   - Finish Firecracker VM integration
   - Create Node.js and Python runners
   - Capture stdout/stderr

5. **Add Real-time Logs**
   - Set up WebSocket server
   - Stream execution logs
   - Update frontend to display live logs

## Troubleshooting

### Backend won't start

- Check PostgreSQL is running: `docker ps | grep postgres`
- Verify DATABASE_URL in backend/.env
- Check port 8080 is not in use: `netstat -an | findstr :8080`

### Frontend build errors

- Delete `.next` folder: `rm -rf .next`
- Clear node_modules: `rm -rf node_modules && npm install`
- Check Node version: `node --version` (should be 18+)

### Database connection issues

- Test connection: `psql -h localhost -U voltrun -d voltrun`
- Check PostgreSQL logs: `docker logs voltrun-postgres`
- Restart database: `docker restart voltrun-postgres`

## Resources

- [Go Fiber Documentation](https://docs.gofiber.io/)
- [Next.js Documentation](https://nextjs.org/docs)
- [Firecracker Documentation](https://github.com/firecracker-microvm/firecracker)
- [GORM Documentation](https://gorm.io/docs/)
- [Tailwind CSS](https://tailwindcss.com/docs)
