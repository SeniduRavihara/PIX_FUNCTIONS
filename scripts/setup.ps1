# VoltRun Quick Setup Script (PowerShell)

Write-Host "⚡ VoltRun Setup" -ForegroundColor Cyan
Write-Host "===============`n" -ForegroundColor Cyan

# Check prerequisites
Write-Host "Checking prerequisites..." -ForegroundColor Yellow

# Check Go
try {
    $goVersion = go version
    Write-Host "✓ Go installed: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ Go not found. Please install Go 1.21+ from https://golang.org/dl/" -ForegroundColor Red
    exit 1
}

# Check Node.js
try {
    $nodeVersion = node --version
    Write-Host "✓ Node.js installed: $nodeVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ Node.js not found. Please install Node.js 18+ from https://nodejs.org/" -ForegroundColor Red
    exit 1
}

# Check Docker
try {
    $dockerVersion = docker --version
    Write-Host "✓ Docker installed: $dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "⚠ Docker not found. Docker is optional but recommended." -ForegroundColor Yellow
}

Write-Host "`nSetting up backend..." -ForegroundColor Yellow
Set-Location backend

# Install Go dependencies
Write-Host "Installing Go dependencies..."
go mod download
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Go dependencies installed" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to install Go dependencies" -ForegroundColor Red
    exit 1
}

# Create .env file
if (-not (Test-Path ".env")) {
    Copy-Item ".env.example" ".env"
    Write-Host "✓ Created backend/.env file" -ForegroundColor Green
} else {
    Write-Host "✓ backend/.env already exists" -ForegroundColor Green
}

Set-Location ..

Write-Host "`nSetting up frontend..." -ForegroundColor Yellow
Set-Location frontend

# Install Node dependencies
Write-Host "Installing Node.js dependencies..."
npm install
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Node.js dependencies installed" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to install Node.js dependencies" -ForegroundColor Red
    exit 1
}

# Create .env.local file
if (-not (Test-Path ".env.local")) {
    Copy-Item ".env.local.example" ".env.local"
    Write-Host "✓ Created frontend/.env.local file" -ForegroundColor Green
} else {
    Write-Host "✓ frontend/.env.local already exists" -ForegroundColor Green
}

Set-Location ..

Write-Host "`n✅ Setup complete!" -ForegroundColor Green
Write-Host "`nNext steps:" -ForegroundColor Cyan
Write-Host "1. Start PostgreSQL:" -ForegroundColor White
Write-Host "   docker run -d --name voltrun-postgres -e POSTGRES_USER=voltrun -e POSTGRES_PASSWORD=voltrun -e POSTGRES_DB=voltrun -p 5432:5432 postgres:15-alpine" -ForegroundColor Gray
Write-Host "`n2. Start backend:" -ForegroundColor White
Write-Host "   cd backend && go run cmd/server/main.go" -ForegroundColor Gray
Write-Host "`n3. Start frontend (in a new terminal):" -ForegroundColor White
Write-Host "   cd frontend && npm run dev" -ForegroundColor Gray
Write-Host "`n4. Open browser:" -ForegroundColor White
Write-Host "   http://localhost:3000" -ForegroundColor Gray
Write-Host "`nOr use Docker Compose:" -ForegroundColor Cyan
Write-Host "   cd deploy && docker-compose up" -ForegroundColor Gray
Write-Host ""
