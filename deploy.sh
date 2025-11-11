#!/bin/bash

# CPRoom Production Deployment Script
set -e

echo "🚀 CPRoom Booking System - Production Deployment"
echo "================================================="
echo ""

# Check if .env.prod exists
if [ ! -f .env.prod ]; then
    echo "⚠️  .env.prod file not found!"
    echo "Creating from template..."
    cp .env.prod.example .env.prod
    echo ""
    echo "✅ Created .env.prod file"
    echo "⚠️  IMPORTANT: Edit .env.prod with your production values before continuing!"
    echo ""
    echo "Required changes:"
    echo "  - DB_PASSWORD"
    echo "  - MONGO_ROOT_PASSWORD"
    echo "  - RABBITMQ_PASS"
    echo "  - JWT_SECRET (generate with: openssl rand -base64 64)"
    echo "  - GITHUB_CLIENT_ID"
    echo "  - GITHUB_CLIENT_SECRET"
    echo "  - SMTP configuration"
    echo ""
    read -p "Press Enter after you've configured .env.prod..."
fi

echo "📋 Loading environment variables..."
export $(cat .env.prod | grep -v '^#' | xargs)

echo "🔍 Checking Docker..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

echo "✅ Docker and Docker Compose are installed"
echo ""

# Check if ports are available
echo "🔍 Checking if ports are available..."
if lsof -Pi :80 -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo "⚠️  Port 80 is already in use!"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

if lsof -Pi :8000 -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo "⚠️  Port 8000 is already in use!"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo "✅ Ports are available"
echo ""

# Build and start services
echo "🏗️  Building and starting services..."
echo "This may take several minutes on first run..."
echo ""

docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d --build

echo ""
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Check service status
echo ""
echo "📊 Service Status:"
docker-compose -f docker-compose.prod.yml ps

echo ""
echo "✅ Deployment complete!"
echo ""
echo "📡 Access your application:"
echo "  Frontend:     http://localhost"
echo "  API Gateway:  http://localhost:8000"
echo "  Health Check: http://localhost/health"
echo ""
echo "📝 Useful commands:"
echo "  View logs:        docker-compose -f docker-compose.prod.yml logs -f"
echo "  Stop services:    docker-compose -f docker-compose.prod.yml down"
echo "  Restart services: docker-compose -f docker-compose.prod.yml restart"
echo ""
echo "📖 For more information, see DEPLOYMENT.md"
echo ""
