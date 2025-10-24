# Quick Start Guide

## Easiest Way to Run (Development Mode)

### Step 1: Start Database Services Only

```powershell
cd deploy
docker-compose -f docker-compose.dev.yml up -d
```

This starts:

- PostgreSQL on port 5432
- Redis on port 6379

### Step 2: Run Backend (New Terminal)

```powershell
cd backend

# First time only: create .env file
Copy-Item .env.example .env

# Install dependencies and run
go mod download
go run cmd/server/main.go
```

Backend will start on **http://localhost:8080**

### Step 3: Run Frontend (New Terminal)

```powershell
cd frontend

# First time only: create .env file
Copy-Item .env.local.example .env.local

# Install dependencies and run
npm install
npm run dev
```

Frontend will start on **http://localhost:3000**

### Step 4: Open Browser

Visit **http://localhost:3000**

---

## Using Full Docker Compose (Production-like)

If you want everything in Docker:

```powershell
cd deploy
docker-compose up --build
```

**Note:** First build takes 5-10 minutes. Subsequent builds are faster.

---

## Quick Test

```powershell
# Test backend health
curl http://localhost:8080/health

# Should return:
# {"status":"healthy","service":"voltrun-backend","version":"1.0.0"}
```

---

## Stop Services

### Development mode:

```powershell
# Stop databases
cd deploy
docker-compose -f docker-compose.dev.yml down

# Stop backend: Ctrl+C in terminal
# Stop frontend: Ctrl+C in terminal
```

### Full Docker mode:

```powershell
cd deploy
docker-compose down
```

---

## Troubleshooting

**"go: command not found"**

- Install Go from https://go.dev/dl/ (need version 1.21+)

**"npm: command not found"**

- Install Node.js from https://nodejs.org/ (need version 18+)

**Port already in use:**

```powershell
# Check what's using the port
netstat -ano | findstr :8080

# Kill the process (replace <PID>)
taskkill /PID <PID> /F
```

**Database connection error:**

```powershell
# Check if PostgreSQL is running
docker ps

# View logs
docker logs voltrun-postgres
```

---

## What's Next?

1. **Implement Authentication** - Edit `backend/internal/api/routes.go`
2. **Test API Endpoints** - Use the login/register pages
3. **Add Features** - See `GETTING_STARTED.md` for detailed guide

---

**Recommended:** Use development mode (docker-compose.dev.yml) for faster iteration!
