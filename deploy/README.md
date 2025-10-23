# VoltRun Deployment Configuration

This directory contains deployment and infrastructure configuration files.

## Files

- `docker-compose.yml` - Docker Compose configuration for local development
- `firecracker-template.json` - Firecracker VM configuration template
- `config/` - Additional configuration files

## Quick Start

### Development Environment

```bash
docker-compose up -d
```

This will start:

- PostgreSQL database (port 5432)
- Redis cache (port 6379)
- Backend API (port 8080)
- Frontend dashboard (port 3000)

### Firecracker Setup

For VM-based function execution, you'll need Firecracker installed:

```bash
# Download Firecracker (Linux only)
curl -Lo firecracker https://github.com/firecracker-microvm/firecracker/releases/download/v1.5.0/firecracker-v1.5.0-x86_64
chmod +x firecracker
sudo mv firecracker /usr/bin/

# Download kernel and rootfs
mkdir -p /var/lib/voltrun
# Download vmlinux.bin and rootfs.ext4 (see Firecracker docs)
```

### Environment Variables

Backend environment variables:

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret for JWT tokens
- `FIRECRACKER_BIN` - Path to Firecracker binary
- `KERNEL_PATH` - Path to VM kernel
- `ROOTFS_PATH` - Path to VM root filesystem

## Production Deployment

For production, consider:

- Using managed PostgreSQL (AWS RDS, DigitalOcean)
- Setting up proper secrets management
- Configuring TLS/SSL certificates
- Setting up monitoring with Prometheus/Grafana
- Using container orchestration (Kubernetes, ECS)
