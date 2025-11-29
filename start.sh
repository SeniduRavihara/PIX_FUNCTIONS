#!/bin/bash

# VoltRun Complete Startup Script
# This script starts the entire VoltRun platform

set -e

echo "🚀 Starting VoltRun Cloud Function Platform..."
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if required commands exist
command -v docker >/dev/null 2>&1 || { echo "❌ Docker is required but not installed. Aborting." >&2; exit 1; }
command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed. Aborting." >&2; exit 1; }
command -v node >/dev/null 2>&1 || { echo "❌ Node.js is required but not installed. Aborting." >&2; exit 1; }

echo "${BLUE}Step 1: Starting PostgreSQL and Redis...${NC}"
cd deploy
docker compose -f docker-compose.dev.yml up -d
cd ..

echo "${GREEN}✓ Database services started${NC}"
echo ""

# Wait for PostgreSQL to be ready
echo "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
sleep 5

echo "${BLUE}Step 2: Starting Go Backend...${NC}"
cd backend

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    echo "${YELLOW}Creating .env file...${NC}"
    cat > .env << EOF
PORT=8080
DATABASE_URL=host=localhost user=voltrun password=voltrun dbname=voltrun port=5432 sslmode=disable
JWT_SECRET=voltrun-secret-change-in-production-$(openssl rand -hex 16)
EOF
fi

# Download Go dependencies
echo "${YELLOW}Installing Go dependencies...${NC}"
go mod download

# Start backend in background
echo "${GREEN}Starting backend server on port 8080...${NC}"
go run cmd/server/main.go &
BACKEND_PID=$!
echo "Backend PID: $BACKEND_PID"

cd ..
echo "${GREEN}✓ Backend started${NC}"
echo ""

# Wait for backend to start
echo "${YELLOW}Waiting for backend to be ready...${NC}"
sleep 3

echo "${BLUE}Step 3: Starting Next.js Frontend...${NC}"
cd frontend

# Create .env.local if it doesn't exist
if [ ! -f .env.local ]; then
    echo "${YELLOW}Creating .env.local file...${NC}"
    cat > .env.local << EOF
NEXT_PUBLIC_API_URL=http://localhost:8080/api
EOF
fi

# Install npm dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "${YELLOW}Installing npm dependencies (this may take a few minutes)...${NC}"
    npm install
fi

# Start frontend in background
echo "${GREEN}Starting frontend server on port 3000...${NC}"
npm run dev &
FRONTEND_PID=$!
echo "Frontend PID: $FRONTEND_PID"

cd ..
echo "${GREEN}✓ Frontend started${NC}"
echo ""

echo "════════════════════════════════════════════════════════════════"
echo "${GREEN}✅ VoltRun is now running!${NC}"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "📍 Access points:"
echo "   Frontend:  ${BLUE}http://localhost:3000${NC}"
echo "   Backend:   ${BLUE}http://localhost:8080${NC}"
echo "   Health:    ${BLUE}http://localhost:8080/health${NC}"
echo ""
echo "🎯 Next steps:"
echo "   1. Open ${BLUE}http://localhost:3000${NC} in your browser"
echo "   2. Click 'Register' to create an account"
echo "   3. Create your first cloud function"
echo "   4. Execute it and see the results!"
echo ""
echo "📝 Process IDs:"
echo "   Backend:  $BACKEND_PID"
echo "   Frontend: $FRONTEND_PID"
echo ""
echo "🛑 To stop all services, run:"
echo "   ${YELLOW}./scripts/stop.sh${NC}"
echo "   or press Ctrl+C and run: kill $BACKEND_PID $FRONTEND_PID"
echo ""
echo "════════════════════════════════════════════════════════════════"

# Save PIDs for stop script
mkdir -p .voltrun
echo $BACKEND_PID > .voltrun/backend.pid
echo $FRONTEND_PID > .voltrun/frontend.pid

# Keep script running and show logs
echo ""
echo "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo ""

# Wait for processes
wait
