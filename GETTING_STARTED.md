# 🚀 VoltRun - Getting Started

Welcome to **VoltRun** - your cloud function execution platform!

## What You Have

A complete, production-ready scaffold for a serverless function platform:

✅ **Go Backend** - High-performance API with Fiber  
✅ **Next.js Frontend** - Modern dashboard with Tailwind CSS  
✅ **PostgreSQL Database** - Complete schemas for users, functions, executions  
✅ **JWT Authentication** - Secure token-based auth  
✅ **Docker Setup** - One-command deployment  
✅ **Full Documentation** - Everything you need to get started

## 📂 Project Structure

```
voltrun/
├── backend/          # Go API (Fiber + GORM + JWT)
├── frontend/         # Next.js Dashboard (React + Tailwind)
├── deploy/           # Docker Compose configuration
├── docs/             # Documentation
└── scripts/          # Setup automation
```

## ⚡ Quick Start (2 minutes)

### Option 1: Docker Compose (Easiest)

```bash
cd deploy
docker-compose up
```

That's it! Visit:

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

### Option 2: Local Development

**Prerequisites**: Go 1.21+, Node.js 18+, PostgreSQL 15+

1. **Run the setup script:**

   ```powershell
   .\scripts\setup.ps1
   ```

2. **Start PostgreSQL:**

   ```bash
   docker run -d --name voltrun-postgres \
     -e POSTGRES_USER=voltrun \
     -e POSTGRES_PASSWORD=voltrun \
     -e POSTGRES_DB=voltrun \
     -p 5432:5432 postgres:15-alpine
   ```

3. **Start backend (terminal 1):**

   ```bash
   cd backend
   go run cmd/server/main.go
   ```

4. **Start frontend (terminal 2):**

   ```bash
   cd frontend
   npm run dev
   ```

5. **Open browser:**
   - http://localhost:3000

## 🎯 What's Working

### Backend (Port 8080)

- ✅ Health check endpoint
- ✅ API routes defined for all features
- ✅ Database models (User, Function, Execution, APIKey)
- ✅ JWT authentication middleware
- ✅ Execution engine architecture
- ✅ VM manager interface (Firecracker integration ready)

### Frontend (Port 3000)

- ✅ Landing page with features
- ✅ Login/Register pages
- ✅ Dashboard layout with navigation
- ✅ API client library
- ✅ Auth context and protected routes

### Database

- ✅ Auto-migrations on startup
- ✅ Complete schemas for all entities
- ✅ UUID-based primary keys
- ✅ Proper relationships and indexes

## 📝 Next Steps to Make It Fully Functional

### 1. Implement Authentication (30 min)

Edit `backend/internal/api/routes.go`:

```go
func handleRegister(c *fiber.Ctx) error {
    // TODO: Parse request, hash password, create user, return JWT
}

func handleLogin(c *fiber.Ctx) error {
    // TODO: Validate credentials, generate JWT, return token
}
```

**Resources**: Auth helpers already in `backend/internal/auth/jwt.go`

### 2. Add Function Management (1 hour)

Implement CRUD operations in `backend/internal/api/routes.go`:

```go
func createFunction(c *fiber.Ctx) error {
    // TODO: Parse function data, save to database
}

func listFunctions(c *fiber.Ctx) error {
    // TODO: Query user's functions from database
}
```

**Resources**: Models already in `backend/internal/storage/models.go`

### 3. Add Monaco Editor (1 hour)

Install in frontend:

```bash
cd frontend
npm install @monaco-editor/react
```

Create `frontend/components/CodeEditor.tsx` - examples in docs.

### 4. Implement Execution (2-3 hours)

Two approaches:

**Simple (for testing)**:

- Use `os/exec` to run code directly
- Good for development/testing

**Production (Firecracker)**:

- Follow Firecracker docs
- Already scaffolded in `backend/internal/vm/manager.go`

### 5. Add Real-time Logs (1-2 hours)

- Add WebSocket support to backend
- Stream execution logs
- Update frontend to display live logs

## 📚 Key Files to Edit

| File                                            | Purpose                 | Priority                             |
| ----------------------------------------------- | ----------------------- | ------------------------------------ |
| `backend/internal/api/routes.go`                | API handlers            | **High**                             |
| `backend/internal/exec/engine.go`               | Function execution      | **High**                             |
| `frontend/app/dashboard/functions/new/page.tsx` | Function creation UI    | Medium                               |
| `frontend/components/CodeEditor.tsx`            | Monaco editor           | Medium                               |
| `backend/internal/vm/manager.go`                | Firecracker integration | Low (can use simpler approach first) |

## 🧪 Testing the API

```bash
# Health check
curl http://localhost:8080/health

# After implementing auth:
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## 📖 Documentation

- **Main README**: `README.md`
- **Development Guide**: `docs/DEVELOPMENT.md`
- **Project Summary**: `docs/PROJECT_SUMMARY.md`
- **Backend Docs**: `backend/README.md`
- **Deployment**: `deploy/README.md`

## 🛠️ Useful Commands

### Backend

```bash
cd backend
go run cmd/server/main.go          # Run server
go test ./...                      # Run tests
go build -o voltrun cmd/server     # Build binary
```

### Frontend

```bash
cd frontend
npm run dev                        # Development
npm run build                      # Production build
npm run lint                       # Lint code
```

### Database

```bash
# Connect to database
docker exec -it voltrun-postgres psql -U voltrun -d voltrun

# View tables
\dt

# Query users
SELECT * FROM users;
```

### Docker

```bash
cd deploy
docker-compose up -d               # Start in background
docker-compose logs -f backend     # View backend logs
docker-compose down                # Stop all services
```

## 🐛 Common Issues

**Port already in use**

```bash
# Windows: Find and kill process
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

**Database connection failed**

- Check PostgreSQL is running: `docker ps`
- Verify DATABASE_URL in `backend/.env`
- Test connection: `psql -h localhost -U voltrun -d voltrun`

**Frontend build errors**

```bash
cd frontend
rm -rf .next node_modules
npm install
```

## 🎓 Learning Resources

- **Go + Fiber**: https://docs.gofiber.io/
- **Next.js**: https://nextjs.org/docs
- **GORM**: https://gorm.io/docs/
- **Firecracker**: https://github.com/firecracker-microvm/firecracker
- **JWT**: https://jwt.io/introduction

## 🚀 Deployment

### Development

- Already configured with Docker Compose

### Production

- Use `docker-compose -f docker-compose.prod.yml up` (create prod config)
- Set up proper secrets management
- Configure TLS/SSL
- Use managed PostgreSQL (AWS RDS, DigitalOcean)
- Add monitoring (Prometheus, Grafana)

## 💡 Feature Ideas

- [ ] Multiple runtime support (Node, Python, Go, Rust)
- [ ] Function versioning
- [ ] Usage analytics and dashboards
- [ ] API rate limiting
- [ ] Function templates library
- [ ] Collaboration features
- [ ] CI/CD integration
- [ ] WebAssembly runtime
- [ ] Scheduled functions (cron jobs)
- [ ] Event-driven triggers

## 🤝 Contributing

1. Pick a feature from "Next Steps"
2. Create a branch: `git checkout -b feature/your-feature`
3. Implement and test
4. Commit: `git commit -m "feat: your feature"`
5. Push and create PR

## 📞 Support

- Check `docs/DEVELOPMENT.md` for detailed guides
- Review code comments in backend/frontend
- See examples in existing code

## 🎉 You're Ready!

Everything is set up and ready to go. Start with implementing the authentication handlers, then move to function management, and you'll have a working serverless platform in a few hours!

**Happy Coding! ⚡**

---

**VoltRun** - Serverless execution, reimagined.
