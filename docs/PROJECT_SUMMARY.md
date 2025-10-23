# VoltRun - Project Summary

## 🎯 Project Overview

**VoltRun** is a cloud function execution platform similar to AWS Lambda or Firebase Cloud Functions, built from scratch using Go (backend) and Next.js (frontend). It allows users to upload code in multiple languages, execute it in isolated Firecracker MicroVMs, and retrieve results via REST API or web dashboard.

## 📦 What's Been Built

### ✅ Backend (Go)

**Location**: `backend/`

#### Structure

```
backend/
├── cmd/server/main.go         # Entry point with Fiber app
├── internal/
│   ├── api/routes.go          # All API endpoints defined
│   ├── auth/
│   │   ├── jwt.go             # JWT token generation/validation
│   │   └── middleware.go       # Auth middleware for protected routes
│   ├── storage/
│   │   ├── models.go           # Database schemas (User, Function, Execution, APIKey)
│   │   └── db.go               # GORM database setup & migrations
│   ├── vm/manager.go          # Firecracker VM lifecycle management (stub)
│   ├── exec/engine.go         # Function execution engine
│   └── utils/
│       ├── logger.go           # Structured logging with Zap
│       └── config.go           # Environment configuration
├── go.mod                     # Go dependencies (Fiber, GORM, JWT, Zap)
├── Dockerfile                 # Multi-stage Docker build
└── .env.example               # Environment variables template
```

#### Key Features

- **API Framework**: Fiber (Express-like for Go)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens with bcrypt password hashing
- **Logging**: Structured logging with Zap
- **API Endpoints**:
  - Auth: `/api/auth/register`, `/api/auth/login`, `/api/auth/refresh`
  - Functions: `/api/functions` (CRUD + execute)
  - Executions: `/api/executions` (history + logs)
  - API Keys: `/api/keys` (management)
  - Health: `/health`

#### Database Models

1. **User**: email, password, name
2. **Function**: name, description, runtime (nodejs/python), code, memory, timeout
3. **Execution**: function_id, status, input/output, logs, duration
4. **APIKey**: user_id, key (hashed), prefix, expiry

### ✅ Frontend (Next.js)

**Location**: `frontend/`

#### Structure

```
frontend/
├── app/
│   ├── layout.tsx             # Root layout with AuthProvider
│   ├── page.tsx               # Landing page with features
│   ├── login/page.tsx         # Login form
│   ├── register/page.tsx      # Registration form
│   └── dashboard/
│       ├── layout.tsx         # Dashboard navigation & header
│       ├── page.tsx           # Functions list
│       ├── executions/        # (stub for execution history)
│       └── keys/              # (stub for API key management)
├── lib/
│   ├── api.ts                 # API client with all endpoints
│   └── auth-context.tsx       # React Context for authentication
└── .env.local.example         # Frontend environment variables
```

#### Key Features

- **Framework**: Next.js 14+ with App Router
- **Styling**: Tailwind CSS
- **Authentication**: JWT-based with React Context
- **Pages**:
  - Landing page with hero section
  - Login/Register forms with validation
  - Dashboard with protected routes
  - Functions list (placeholder - connects to API)
- **API Client**: Centralized fetch wrapper with token management

### ✅ Deployment

**Location**: `deploy/`

- **docker-compose.yml**: Full stack setup (PostgreSQL, Redis, Backend, Frontend)
- **firecracker-template.json**: VM configuration template
- **README.md**: Deployment instructions

### ✅ Documentation

**Location**: `docs/` and `scripts/`

- **README.md**: Main project documentation
- **docs/DEVELOPMENT.md**: Comprehensive development guide
- **backend/README.md**: Backend-specific docs
- **deploy/README.md**: Deployment guide
- **scripts/setup.ps1**: Automated PowerShell setup script

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

```bash
cd deploy
docker-compose up
```

- Backend: http://localhost:8080
- Frontend: http://localhost:3000
- PostgreSQL: localhost:5432

### Option 2: Manual Setup

```powershell
# Run setup script
.\scripts\setup.ps1

# Start PostgreSQL
docker run -d --name voltrun-postgres \
  -e POSTGRES_USER=voltrun \
  -e POSTGRES_PASSWORD=voltrun \
  -e POSTGRES_DB=voltrun \
  -p 5432:5432 postgres:15-alpine

# Start backend (terminal 1)
cd backend
go run cmd/server/main.go

# Start frontend (terminal 2)
cd frontend
npm run dev
```

## 🏗️ Architecture

```
┌─────────────────┐
│   Next.js UI    │  Port 3000
│  - Login/Signup │  - Monaco Editor (TODO)
│  - Dashboard    │  - Function Management
└────────┬────────┘
         │ HTTP (JWT)
         ▼
┌─────────────────┐
│   Go Backend    │  Port 8080
│  - Fiber API    │  - JWT Auth
│  - GORM ORM     │  - Execution Engine
└────┬───────┬────┘
     │       │
     │       └──────────┐
     ▼                  ▼
┌─────────┐      ┌──────────────┐
│Postgres │      │ Firecracker  │
│  - Users│      │  MicroVMs    │
│  - Funcs│      │  (Stubbed)   │
└─────────┘      └──────────────┘
```

## 📋 Implementation Status

### ✅ Completed

- [x] Project structure and folder organization
- [x] Go backend with Fiber framework
- [x] PostgreSQL models and GORM setup
- [x] JWT authentication system
- [x] API route definitions
- [x] Execution engine architecture
- [x] Next.js frontend setup
- [x] Authentication pages (login/register)
- [x] Dashboard layout and navigation
- [x] API client library
- [x] Docker and docker-compose setup
- [x] Documentation (README, dev guide)
- [x] Setup automation script

### 🚧 Partially Implemented

- [ ] Auth handlers (routes defined, logic TODO)
- [ ] Function CRUD (routes defined, logic TODO)
- [ ] VM manager (interface ready, Firecracker integration TODO)
- [ ] Execution engine (structure ready, actual execution TODO)

### 📝 Next Steps

1. **Implement Auth Handlers**

   - Add registration logic with password hashing
   - Implement login with JWT generation
   - Test with frontend forms

2. **Build Function Management**

   - Create function storage/retrieval
   - Implement function list page
   - Add function creation form

3. **Add Monaco Editor**

   - Install `@monaco-editor/react`
   - Create code editor component
   - Add syntax highlighting

4. **Implement Execution**

   - Integrate Firecracker SDK (Linux only)
   - Build Node.js runner
   - Build Python runner
   - Capture output and logs

5. **Real-time Features**
   - Add WebSocket support
   - Stream execution logs
   - Live status updates

## 🔧 Technology Stack

| Component         | Technology    | Purpose                     |
| ----------------- | ------------- | --------------------------- |
| **Backend**       | Go 1.21+      | High-performance API server |
| **API Framework** | Fiber v2      | Fast HTTP routing           |
| **Database**      | PostgreSQL 15 | Relational data storage     |
| **ORM**           | GORM          | Database abstraction        |
| **Auth**          | JWT + bcrypt  | Secure authentication       |
| **Logging**       | Zap           | Structured logging          |
| **VM**            | Firecracker   | Isolated execution          |
| **Frontend**      | Next.js 14+   | React framework             |
| **Styling**       | Tailwind CSS  | Utility-first CSS           |
| **State**         | React Context | Client state management     |
| **Deployment**    | Docker        | Containerization            |

## 📚 Available Commands

### Backend

```bash
cd backend
go run cmd/server/main.go      # Run server
go build -o voltrun cmd/server  # Build binary
go test ./...                   # Run tests
```

### Frontend

```bash
cd frontend
npm run dev                     # Development server
npm run build                   # Production build
npm start                       # Start production server
```

### Docker

```bash
cd deploy
docker-compose up               # Start all services
docker-compose down             # Stop all services
docker-compose logs -f backend  # View backend logs
```

## 🔐 Environment Variables

### Backend (`.env`)

```
PORT=8080
DATABASE_URL=host=localhost user=voltrun password=voltrun dbname=voltrun port=5432 sslmode=disable
JWT_SECRET=change-in-production
FIRECRACKER_BIN=/usr/bin/firecracker
KERNEL_PATH=/var/lib/voltrun/vmlinux.bin
ROOTFS_PATH=/var/lib/voltrun/rootfs.ext4
```

### Frontend (`.env.local`)

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

## 🎯 API Endpoints

| Method | Endpoint                     | Description      | Auth Required |
| ------ | ---------------------------- | ---------------- | ------------- |
| GET    | `/health`                    | Health check     | No            |
| POST   | `/api/auth/register`         | Register user    | No            |
| POST   | `/api/auth/login`            | Login user       | No            |
| POST   | `/api/auth/refresh`          | Refresh token    | Yes           |
| GET    | `/api/functions`             | List functions   | Yes           |
| POST   | `/api/functions`             | Create function  | Yes           |
| GET    | `/api/functions/:id`         | Get function     | Yes           |
| PUT    | `/api/functions/:id`         | Update function  | Yes           |
| DELETE | `/api/functions/:id`         | Delete function  | Yes           |
| POST   | `/api/functions/:id/execute` | Execute function | Yes           |
| GET    | `/api/executions`            | List executions  | Yes           |
| GET    | `/api/executions/:id`        | Get execution    | Yes           |
| GET    | `/api/executions/:id/logs`   | Get logs         | Yes           |
| GET    | `/api/keys`                  | List API keys    | Yes           |
| POST   | `/api/keys`                  | Create API key   | Yes           |
| DELETE | `/api/keys/:id`              | Delete API key   | Yes           |

## 🐛 Troubleshooting

**Backend won't start**

- Ensure PostgreSQL is running
- Check `DATABASE_URL` in `.env`
- Verify port 8080 is available

**Frontend errors**

- Delete `.next` folder
- Run `npm install` again
- Check `NEXT_PUBLIC_API_URL`

**Database issues**

- Test connection: `psql -h localhost -U voltrun -d voltrun`
- Check Docker logs: `docker logs voltrun-postgres`

## 📄 License

MIT License - See LICENSE file for details

---

**Built with ⚡ by VoltRun Team**
