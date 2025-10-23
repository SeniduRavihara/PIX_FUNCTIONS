# VoltRun Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                          VoltRun Platform                            │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────┐
│   End User      │
│   (Browser)     │
└────────┬────────┘
         │ HTTPS
         ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     Frontend Layer (Next.js)                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │  Landing     │  │   Auth UI    │  │  Dashboard   │              │
│  │  Page        │  │ Login/Signup │  │  Functions   │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
│                                                                       │
│  ┌───────────────────────────────────────────────────┐              │
│  │  Components:  Monaco Editor, Function List,       │              │
│  │               Execution Logs, API Key Management  │              │
│  └───────────────────────────────────────────────────┘              │
│                                                                       │
│  Port 3000  •  React 18  •  Tailwind CSS  •  TypeScript             │
└────────┬────────────────────────────────────────────────────────────┘
         │ HTTP/WebSocket
         │ JWT Bearer Token
         ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     Backend Layer (Go)                               │
│                                                                       │
│  ┌─────────────────────────────────────────────────────────┐        │
│  │              Fiber HTTP Server (Port 8080)              │        │
│  │                                                          │        │
│  │  ┌────────────┐  ┌────────────┐  ┌──────────────┐      │        │
│  │  │   Auth     │  │  Function  │  │  Execution   │      │        │
│  │  │  Routes    │  │   Routes   │  │   Routes     │      │        │
│  │  └────────────┘  └────────────┘  └──────────────┘      │        │
│  │                                                          │        │
│  │  ┌─────────────────────────────────────────────────┐   │        │
│  │  │         Middleware Layer                        │   │        │
│  │  │  • CORS  • Logger  • Auth  • Rate Limit        │   │        │
│  │  └─────────────────────────────────────────────────┘   │        │
│  └──────────────────────────┬───────────────────────────┘          │
│                              │                                       │
│  ┌──────────────────────────┴───────────────────────────┐          │
│  │              Business Logic Layer                     │          │
│  │                                                        │          │
│  │  ┌──────────────┐  ┌──────────────┐  ┌────────────┐ │          │
│  │  │  Auth        │  │  Storage     │  │  Execution │ │          │
│  │  │  Service     │  │  Service     │  │  Engine    │ │          │
│  │  │              │  │              │  │            │ │          │
│  │  │  • JWT Gen   │  │  • CRUD      │  │  • Queue   │ │          │
│  │  │  • Validate  │  │  • Query     │  │  • VM Mgmt │ │          │
│  │  │  • Bcrypt    │  │  • Store     │  │  • Execute │ │          │
│  │  └──────────────┘  └──────────────┘  └────────────┘ │          │
│  └────────────────────────────────────────────────────┘           │
│                                                                       │
└────┬──────────────────────────────┬──────────────────────┬──────────┘
     │                              │                       │
     ▼                              ▼                       ▼
┌─────────────┐            ┌─────────────┐      ┌──────────────────┐
│  PostgreSQL │            │   Redis     │      │  VM Orchestrator │
│  Database   │            │   Cache     │      │                  │
│             │            │             │      │  ┌────────────┐  │
│  • Users    │            │  • Session  │      │  │ Firecracker│  │
│  • Functions│            │  • Queue    │      │  │  MicroVM   │  │
│  • Executions            │  • Results  │      │  │            │  │
│  • API Keys │            │             │      │  │ ┌────────┐ │  │
│             │            │             │      │  │ │ Node.js│ │  │
│  Port 5432  │            │  Port 6379  │      │  │ │ Runtime│ │  │
└─────────────┘            └─────────────┘      │  │ └────────┘ │  │
                                                 │  │            │  │
                                                 │  │ ┌────────┐ │  │
                                                 │  │ │ Python │ │  │
                                                 │  │ │ Runtime│ │  │
                                                 │  │ └────────┘ │  │
                                                 │  └────────────┘  │
                                                 │                  │
                                                 │  Isolated VMs    │
                                                 │  Linux Kernel    │
                                                 └──────────────────┘
```

## Request Flow

### 1. User Registration/Login

```
User → Frontend → POST /api/auth/register
                        ↓
                   Validate Input
                        ↓
                   Hash Password (bcrypt)
                        ↓
                   Save to Database
                        ↓
                   Generate JWT Token
                        ↓
                   Return Token + User Data
                        ↓
                   Frontend stores token
                        ↓
                   Redirect to Dashboard
```

### 2. Function Creation

```
User → Frontend → POST /api/functions
                       (JWT Token in Header)
                        ↓
                   Validate JWT
                        ↓
                   Extract User ID
                        ↓
                   Parse Function Data
                    (name, runtime, code)
                        ↓
                   Save to Database
                        ↓
                   Return Function Object
                        ↓
                   Update UI
```

### 3. Function Execution

```
User → Frontend → POST /api/functions/:id/execute
                       (JWT Token + Input Data)
                        ↓
                   Validate JWT
                        ↓
                   Fetch Function from DB
                        ↓
                   Create Execution Record
                    (status: pending)
                        ↓
                   Queue Execution Request
                        ↓
              ┌────────┴─────────┐
              ▼                  ▼
         VM Manager        Update Status
              │              (running)
              ▼
         Create VM
      (Firecracker MicroVM)
              │
              ▼
         Inject User Code
              │
              ▼
         Execute with Limits
       (Memory, CPU, Timeout)
              │
              ▼
         Capture Output
       (stdout, stderr, logs)
              │
              ▼
         Destroy VM
              │
              ▼
         Update Execution Record
        (status: success/failed)
        (output, logs, duration)
              │
              ▼
         Return Results to User
              │
              ▼
         Frontend Displays Output
```

## Component Responsibilities

### Frontend (Next.js)

| Component        | Responsibility                       |
| ---------------- | ------------------------------------ |
| **Pages**        | Route handling, data fetching        |
| **Components**   | Reusable UI (Monaco, forms, tables)  |
| **API Client**   | HTTP requests, token management      |
| **Auth Context** | Global auth state                    |
| **Hooks**        | Custom logic (useAuth, useFunctions) |

### Backend (Go)

| Package              | Responsibility               |
| -------------------- | ---------------------------- |
| **cmd/server**       | Application entry point      |
| **internal/api**     | HTTP routes and handlers     |
| **internal/auth**    | JWT, bcrypt, middleware      |
| **internal/storage** | Database models and queries  |
| **internal/vm**      | Firecracker VM lifecycle     |
| **internal/exec**    | Execution engine and runners |
| **internal/utils**   | Config, logging, helpers     |

### Database (PostgreSQL)

| Table          | Purpose                              |
| -------------- | ------------------------------------ |
| **users**      | User accounts (email, password hash) |
| **functions**  | Function code and metadata           |
| **executions** | Execution history and logs           |
| **api_keys**   | API keys for programmatic access     |

## Security Layers

```
┌────────────────────────────────────────┐
│  1. HTTPS/TLS Encryption               │
└────────────────┬───────────────────────┘
                 ▼
┌────────────────────────────────────────┐
│  2. JWT Authentication                 │
│     • Token validation                 │
│     • User identification              │
└────────────────┬───────────────────────┘
                 ▼
┌────────────────────────────────────────┐
│  3. Authorization Checks               │
│     • User owns resource?              │
│     • Permissions valid?               │
└────────────────┬───────────────────────┘
                 ▼
┌────────────────────────────────────────┐
│  4. Input Validation                   │
│     • Sanitize input                   │
│     • Type checking                    │
└────────────────┬───────────────────────┘
                 ▼
┌────────────────────────────────────────┐
│  5. VM Isolation (Firecracker)         │
│     • Complete process isolation       │
│     • Resource limits                  │
│     • Network isolation                │
└────────────────────────────────────────┘
```

## Scalability Considerations

### Horizontal Scaling

```
┌────────────┐
│ Load       │
│ Balancer   │
└─────┬──────┘
      │
      ├──────────┬──────────┬──────────┐
      ▼          ▼          ▼          ▼
  Backend 1  Backend 2  Backend 3  Backend N
      │          │          │          │
      └──────────┴──────────┴──────────┘
                     │
                     ▼
            ┌────────────────┐
            │   PostgreSQL   │
            │   (Primary)    │
            └────────┬───────┘
                     │
                     ├──────────┬
                     ▼          ▼
              Read Replica  Read Replica
```

### VM Pool Management

```
┌──────────────────────────────────┐
│      VM Pool Manager             │
│                                  │
│  ┌────────┐  ┌────────┐         │
│  │ Warm   │  │ Cold   │         │
│  │ VMs    │  │ VMs    │         │
│  │ (Ready)│  │(Create)│         │
│  └────────┘  └────────┘         │
│                                  │
│  • Pre-warm common runtimes     │
│  • Auto-scale based on demand   │
│  • Health checks & recycling    │
└──────────────────────────────────┘
```

## Data Flow

```
┌──────┐     ┌─────────┐     ┌──────────┐     ┌────┐
│Client│────▶│Frontend │────▶│ Backend  │────▶│ DB │
└──────┘     └─────────┘     └──────────┘     └────┘
   ▲             │                 │              │
   │             │                 ▼              │
   │             │            ┌─────────┐         │
   │             │            │   VM    │         │
   │             │            │ Engine  │         │
   │             │            └─────────┘         │
   │             │                 │              │
   │             │                 ▼              ▼
   │             │            Execute Code    Store Results
   │             │                 │              │
   │             │◀────────────────┴──────────────┘
   │             │          Return Output
   │◀────────────┘
   │        Update UI
```

## Technology Stack Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                      Client Side                             │
│  React 18 • Next.js 14 • TypeScript • Tailwind CSS          │
│  Monaco Editor • SWR • Axios/Fetch                          │
└─────────────────────────────────────────────────────────────┘
                           ▲│
                            ││ HTTP/WebSocket
                           │▼
┌─────────────────────────────────────────────────────────────┐
│                      Server Side                             │
│  Go 1.21+ • Fiber v2 • GORM • JWT • Zap                     │
│  Firecracker SDK • WebSocket                                │
└─────────────────────────────────────────────────────────────┘
                           ▲│
                            ││
                           │▼
┌─────────────────────────────────────────────────────────────┐
│                    Data & Infrastructure                     │
│  PostgreSQL 15 • Redis • Docker • Linux Kernel              │
│  Firecracker MicroVMs • Node.js • Python                    │
└─────────────────────────────────────────────────────────────┘
```

## Development vs Production

### Development

- Docker Compose for all services
- Hot reload for frontend (npm run dev)
- Development database with test data
- Simplified VM execution (direct exec)

### Production

- Kubernetes/Docker Swarm orchestration
- Managed PostgreSQL (RDS, DigitalOcean)
- Redis cluster for caching
- Dedicated VM worker nodes
- CDN for frontend assets
- Monitoring (Prometheus, Grafana)
- Centralized logging (ELK stack)

---

This architecture is designed to be:

- **Scalable**: Horizontally scale backend and VMs
- **Secure**: Multiple layers of isolation and auth
- **Fast**: Efficient Go backend, microVMs for isolation
- **Maintainable**: Clean separation of concerns
- **Observable**: Structured logging and metrics
