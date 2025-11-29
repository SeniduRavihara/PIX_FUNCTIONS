# 🚀 VoltRun - Quick Start Guide

## How to Run VoltRun (3 Options)

### ⚡ Option 1: Automated Startup (EASIEST)

```bash
./start.sh
```

This single command starts everything:

- PostgreSQL database
- Redis cache
- Go backend (port 8080)
- Next.js frontend (port 3000)

Then open: **http://localhost:3000**

To stop everything:

```bash
./stop.sh
```

---

### 🔧 Option 2: Manual Startup (Step by Step)

**Terminal 1 - Database:**

```bash
cd deploy
docker-compose -f docker-compose.dev.yml up
```

**Terminal 2 - Backend:**

```bash
cd backend
go run cmd/server/main.go
```

**Terminal 3 - Frontend:**

```bash
cd frontend
npm run dev
```

**Terminal 4 - Stop Everything:**

```bash
cd deploy
docker-compose -f docker-compose.dev.yml down
```

---

### 🐳 Option 3: Full Docker (Everything Containerized)

```bash
cd deploy
docker-compose up --build
```

This builds and runs everything in Docker containers.

---

## 🎯 First Time Usage

1. **Register Account**

   - Go to http://localhost:3000
   - Click "Sign Up"
   - Create your account

2. **Create a Function**

   - Click "New Function"
   - Choose runtime (Node.js or Python)
   - Paste your code
   - Click "Create"

3. **Execute Function**
   - Click on your function
   - Enter JSON input: `{"name": "World"}`
   - Click "Execute Function"
   - See results!

---

## 📝 Example Functions

### Node.js Example

```javascript
exports.handler = async (event) => {
  console.log("Received event:", JSON.stringify(event));

  return {
    statusCode: 200,
    body: JSON.stringify({
      message: `Hello, ${event.name || "World"}!`,
      timestamp: new Date().toISOString(),
      input: event,
    }),
  };
};
```

**Test with:**

```json
{
  "name": "VoltRun",
  "action": "test"
}
```

---

### Python Example

```python
import json
from datetime import datetime

def handler(event):
    print(f'Received event: {json.dumps(event)}')

    name = event.get('name', 'World')

    return {
        'statusCode': 200,
        'body': json.dumps({
            'message': f'Hello, {name}!',
            'timestamp': datetime.now().isoformat(),
            'input': event
        })
    }
```

**Test with:**

```json
{
  "name": "Python User",
  "version": "3.11"
}
```

---

## 🔍 Verify Installation

### Check Backend Health

```bash
curl http://localhost:8080/health
```

Should return:

```json
{
  "status": "healthy",
  "service": "voltrun-backend",
  "version": "1.0.0"
}
```

### Check Database Connection

```bash
docker exec -it voltrun-postgres psql -U voltrun -d voltrun -c "SELECT COUNT(*) FROM users;"
```

---

## 🐛 Troubleshooting

### Backend won't start

```bash
# Check if port 8080 is available
lsof -i :8080

# Check if PostgreSQL is running
docker ps | grep postgres

# Check logs
cd backend
go run cmd/server/main.go
```

### Frontend won't start

```bash
# Reinstall dependencies
cd frontend
rm -rf node_modules .next
npm install
npm run dev
```

### Database connection failed

```bash
# Restart PostgreSQL
cd deploy
docker-compose -f docker-compose.dev.yml restart postgres

# Check connection
docker logs voltrun-postgres
```

### "Module not found" errors in Go

```bash
cd backend
go mod tidy
go mod download
```

---

## 📊 Services Overview

| Service    | Port | URL                   | Status Check               |
| ---------- | ---- | --------------------- | -------------------------- |
| Frontend   | 3000 | http://localhost:3000 | Open in browser            |
| Backend    | 8080 | http://localhost:8080 | curl localhost:8080/health |
| PostgreSQL | 5432 | localhost:5432        | docker ps \| grep postgres |
| Redis      | 6379 | localhost:6379        | docker ps \| grep redis    |

---

## 🔐 Default Credentials

**Database:**

- Host: localhost
- Port: 5432
- User: voltrun
- Password: voltrun
- Database: voltrun

**Admin Account:**
Create via registration form (no default admin)

---

## 📁 Important Files

```
PIX_FUNCTIONS/
├── start.sh                    # ⚡ One-command startup
├── stop.sh                     # 🛑 Stop all services
├── backend/
│   ├── cmd/server/main.go      # Backend entry point
│   └── internal/
│       ├── api/routes.go       # All API endpoints
│       ├── runners/            # Code execution
│       └── storage/models.go   # Database models
├── frontend/
│   ├── app/                    # Next.js pages
│   └── lib/api.ts              # API client
└── deploy/
    └── docker-compose.dev.yml  # Development setup
```

---

## 🎓 What You Can Do

✅ **User Management**

- Register/Login
- JWT authentication
- Secure password storage

✅ **Function Management**

- Create functions (Node.js, Python)
- Edit code
- Delete functions
- View function details

✅ **Execution**

- Execute functions with JSON input
- Real-time results
- View execution logs
- See execution history

✅ **API Keys**

- Generate API keys
- Use for programmatic access
- Secure key storage

✅ **Monitoring**

- Execution history
- Status tracking
- Performance metrics
- Error logs

---

## 🚀 Next Steps

1. **Test the Platform:**

   - Create a function
   - Execute it
   - View results and logs

2. **Explore Features:**

   - Try different runtimes
   - Test with various inputs
   - Check execution history

3. **API Integration:**

   - Create an API key
   - Use it to call functions programmatically

4. **Production Deployment:**
   - Set up proper secrets
   - Configure TLS/SSL
   - Use managed database

---

## 📞 Quick Commands Reference

```bash
# Start everything
./start.sh

# Stop everything
./stop.sh

# Backend only
cd backend && go run cmd/server/main.go

# Frontend only
cd frontend && npm run dev

# Database only
cd deploy && docker-compose -f docker-compose.dev.yml up

# View logs
docker logs voltrun-postgres
docker logs voltrun-redis

# Clean restart
./stop.sh && ./start.sh

# Check running processes
ps aux | grep "go run"
ps aux | grep "next-server"
```

---

**You're all set! 🎉**

Run `./start.sh` and start building with VoltRun!
