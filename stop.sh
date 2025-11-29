#!/bin/bash

# VoltRun Stop Script
# Stops all VoltRun services

set -e

echo "🛑 Stopping VoltRun services..."

# Stop frontend
if [ -f .voltrun/frontend.pid ]; then
    FRONTEND_PID=$(cat .voltrun/frontend.pid)
    if ps -p $FRONTEND_PID > /dev/null 2>&1; then
        echo "Stopping Frontend (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID 2>/dev/null || true
    fi
fi

# Stop backend
if [ -f .voltrun/backend.pid ]; then
    BACKEND_PID=$(cat .voltrun/backend.pid)
    if ps -p $BACKEND_PID > /dev/null 2>&1; then
        echo "Stopping Backend (PID: $BACKEND_PID)..."
        kill $BACKEND_PID 2>/dev/null || true
    fi
fi

# Stop Docker containers
echo "Stopping Docker containers..."
cd deploy
docker-compose -f docker-compose.dev.yml down
cd ..

# Clean up PID files
rm -rf .voltrun

echo "✅ All VoltRun services stopped"
