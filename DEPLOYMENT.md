# CPRoom Production Deployment Guide

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose installed
- Ports 80 and 8000 available
- GitHub OAuth app created (for OAuth login)

### 1. Setup Environment Variables

```bash
# Copy the example environment file
cp .env.prod.example .env.prod

# Edit the file with your production values
nano .env.prod
```

**Required Changes:**
- `DB_PASSWORD` - Strong database password
- `MONGO_ROOT_PASSWORD` - MongoDB admin password
- `RABBITMQ_PASS` - RabbitMQ password
- `JWT_SECRET` - Generate with: `openssl rand -base64 64`
- `GITHUB_CLIENT_ID` - From GitHub OAuth app
- `GITHUB_CLIENT_SECRET` - From GitHub OAuth app
- `SMTP_*` - Email server configuration

### 2. Configure GitHub OAuth

1. Go to: https://github.com/settings/developers
2. Click "OAuth Apps" → "New OAuth App"
3. Fill in:
   - **Application name**: CPRoom Booking System
   - **Homepage URL**: `http://your-domain.com` (or `http://localhost` for testing)
   - **Authorization callback URL**: `http://your-domain.com:8000/auth/github/callback`
4. Copy Client ID and Client Secret to `.env.prod`

### 3. Build and Start Services

```bash
# Load environment variables
export $(cat .env.prod | grep -v '^#' | xargs)

# Build and start all services
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d --build

# Check status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f
```

### 4. Access the Application

- **Frontend**: http://localhost (port 80)
- **API Gateway**: http://localhost:8000
- **Health Check**: http://localhost/health

## 📋 Service Architecture

```
┌─────────────────────────────────────────┐
│         Frontend (Port 80)              │
│    React + TypeScript + Nginx           │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      API Gateway - Kong (Port 8000)     │
│        JWT Authentication               │
└─────────────────┬───────────────────────┘
                  │
    ┌─────────────┼─────────────┐
    │             │             │
┌───▼──┐   ┌─────▼──┐   ┌─────▼──────┐
│Auth  │   │Booking │   │Notification│
│8081  │   │8083    │   │8084        │
└──────┘   └────────┘   └────────────┘
    │             │             │
┌───▼──┐   ┌─────▼──┐   ┌─────▼──────┐
│Room  │   │Approval│   │  RabbitMQ  │
│8082  │   │8085    │   │            │
└──────┘   └────────┘   └────────────┘
    │             │             │
┌───▼─────────────▼─────────────▼──────┐
│       PostgreSQL & MongoDB            │
└───────────────────────────────────────┘
```

## 🔒 Security Features

### Exposed Ports (Minimal)
Only 2 ports are exposed to the host:
- **Port 80**: Frontend (Nginx)
- **Port 8000**: API Gateway (Kong)

All backend services are isolated in the internal network.

### Built-in Security
- JWT authentication on all API routes
- HTTPS ready (add SSL certificates to Kong)
- Security headers in Nginx
- Database credentials isolated
- No direct database access from outside

## 🛠️ Management Commands

### View Logs
```bash
# All services
docker-compose -f docker-compose.prod.yml logs -f

# Specific service
docker-compose -f docker-compose.prod.yml logs -f frontend
docker-compose -f docker-compose.prod.yml logs -f api-gateway
docker-compose -f docker-compose.prod.yml logs -f auth-service
```

### Stop Services
```bash
# Stop all services
docker-compose -f docker-compose.prod.yml down

# Stop and remove volumes (WARNING: Deletes all data)
docker-compose -f docker-compose.prod.yml down -v
```

### Restart Services
```bash
# Restart all
docker-compose -f docker-compose.prod.yml restart

# Restart specific service
docker-compose -f docker-compose.prod.yml restart frontend
```

### Scale Services (if needed)
```bash
# Scale booking service to 3 instances
docker-compose -f docker-compose.prod.yml up -d --scale booking-service=3
```

### Update Services
```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose -f docker-compose.prod.yml up -d --build
```

## 📊 Health Checks

All services have health checks configured:
- **Frontend**: `/health` endpoint
- **API Gateway**: Kong health check
- **Databases**: Connection checks
- **RabbitMQ**: Diagnostics ping

Check health status:
```bash
docker-compose -f docker-compose.prod.yml ps
```

## 🔍 Troubleshooting

### Service won't start
```bash
# Check logs
docker-compose -f docker-compose.prod.yml logs [service-name]

# Check if port is already in use
sudo lsof -i :80
sudo lsof -i :8000
```

### Database connection issues
```bash
# Check if postgres is healthy
docker-compose -f docker-compose.prod.yml exec postgres pg_isready

# Check MongoDB
docker-compose -f docker-compose.prod.yml exec mongo mongosh --eval "db.adminCommand('ping')"
```

### OAuth not working
1. Verify GitHub OAuth callback URL matches exactly
2. Check `GITHUB_REDIRECT_URL` in `.env.prod` uses port 8000
3. Check `FRONTEND_URL` points to correct domain

### Clear and restart
```bash
# Stop everything
docker-compose -f docker-compose.prod.yml down

# Remove old images
docker-compose -f docker-compose.prod.yml down --rmi all

# Rebuild from scratch
docker-compose -f docker-compose.prod.yml up -d --build --force-recreate
```

## 🌐 Production Deployment (Custom Domain)

### 1. Update Environment Variables
```bash
# In .env.prod
FRONTEND_URL=https://your-domain.com
GITHUB_REDIRECT_URL=https://your-domain.com:8000/auth/github/callback
VITE_API_BASE_URL=https://your-domain.com:8000
```

### 2. Configure SSL/TLS
Add SSL certificates to Kong or use a reverse proxy (Nginx/Traefik) in front.

### 3. Update GitHub OAuth
Update callback URL in GitHub OAuth app settings.

### 4. Use Production SMTP
Configure with your email service (SendGrid, AWS SES, etc.)

## 📦 Data Persistence

Volumes are created for persistent storage:
- `postgres_data` - PostgreSQL database
- `mongo_data` - MongoDB database

### Backup Data
```bash
# Backup PostgreSQL
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U postgres cproom_db > backup.sql

# Backup MongoDB
docker-compose -f docker-compose.prod.yml exec mongo mongodump --out /tmp/backup
docker cp cproom-mongo:/tmp/backup ./mongo-backup
```

### Restore Data
```bash
# Restore PostgreSQL
docker-compose -f docker-compose.prod.yml exec -T postgres psql -U postgres cproom_db < backup.sql

# Restore MongoDB
docker cp ./mongo-backup cproom-mongo:/tmp/backup
docker-compose -f docker-compose.prod.yml exec mongo mongorestore /tmp/backup
```

## 🎯 Performance Tips

1. **Use Docker volumes** for better I/O performance
2. **Monitor resource usage**: `docker stats`
3. **Set resource limits** in compose file if needed
4. **Use logging drivers** for production logging
5. **Enable Docker BuildKit** for faster builds:
   ```bash
   export DOCKER_BUILDKIT=1
   ```

## 📝 Notes

- The original `backend/docker-compose.yml` is preserved for development
- This production compose includes only essential exposed ports
- All services use health checks for reliability
- Automatic restarts enabled with `unless-stopped`
- Frontend is optimized with Nginx and gzip compression

---

**Need Help?** Check the logs first, most issues are configuration-related.
