# VoltRun ⚡

**Cloud Function Execution Platform**

VoltRun is a cloud function platform similar to Firebase Cloud Functions or AWS Lambda, built from scratch. It allows users to upload code in multiple languages (JavaScript, Python, etc.), execute it on-demand in isolated Firecracker MicroVMs, and get real-time results via API or dashboard.

## 🚀 Features

- **Multi-language Support**: Execute JavaScript, Python, and more
- **Secure Isolation**: Each function runs in its own Firecracker MicroVM
- **Real-time Execution**: Live logs and results via WebSocket/REST API
- **Code Management**: Built-in Monaco editor for writing and deploying functions
- **Authentication**: JWT-based auth with API key management
- **Monitoring**: Execution logs, metrics, and dashboard analytics

## 🧩 Technology Stack

**Backend (Go)**

- Language: Go (Golang)
- VM Management: Firecracker MicroVM
- Database: PostgreSQL
- Message Queue: NATS/Redis (optional)
- API Framework: Fiber
- Logging: Zap

**Frontend (Next.js)**

- Framework: Next.js 14+ (App Router)
- Styling: Tailwind CSS + shadcn/ui
- Code Editor: Monaco Editor
- Auth: JWT / Firebase Auth

## 📁 Project Structure

```
voltrun/
├── backend/                      # Go backend service
│   ├── cmd/
│   │   └── server/               # Main server entry point
│   ├── internal/
│   │   ├── api/                  # HTTP routes & controllers
│   │   ├── vm/                   # Firecracker VM management
│   │   ├── exec/                 # Function execution engine
│   │   ├── storage/              # Database & S3 interfaces
│   │   ├── auth/                 # JWT authentication
│   │   └── utils/                # Config, logging, helpers
│   ├── go.mod
│   └── Dockerfile
│
├── frontend/                     # Next.js web dashboard
│   ├── app/                      # App Router pages
│   ├── components/               # UI components
│   ├── lib/                      # API client & utilities
│   └── package.json
│
├── deploy/                       # Deployment configs
│   ├── docker-compose.yml
│   ├── firecracker-template.json
│   └── config/
│
├── docs/                         # Documentation
└── scripts/                      # DevOps scripts
```

## 🛠️ Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- Firecracker (for local VM execution)
- PostgreSQL 15+

### Backend Setup

```bash
cd backend
go mod download
go run cmd/server/main.go
```

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

### Docker Setup

```bash
docker-compose up -d
```

## 📊 API Endpoints

- `POST /api/functions` - Create new function
- `GET /api/functions` - List user functions
- `POST /api/functions/:id/execute` - Execute function
- `GET /api/executions/:id/logs` - Get execution logs
- `POST /api/auth/login` - User authentication
- `POST /api/keys` - Generate API key

## 🔐 Security

- JWT-based authentication
- API key authorization for function invocation
- Complete VM isolation per execution
- Resource limits (CPU, memory, timeout)
- Code sandboxing with Firecracker

## 🚧 Roadmap

- [ ] Core execution engine with Firecracker
- [ ] REST API with Fiber
- [ ] PostgreSQL schemas and migrations
- [ ] Frontend dashboard with Monaco editor
- [ ] Real-time logs via WebSocket
- [ ] Multi-runtime support (Node.js, Python, Go)
- [ ] Usage tracking and analytics
- [ ] WebAssembly runtime support
- [ ] AI-driven optimization suggestions

## 📝 License

MIT

---

**VoltRun** - Serverless execution, reimagined.
