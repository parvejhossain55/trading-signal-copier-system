# Signal Copier Project

A full-stack application with Go API backend and Next.js frontend, containerized with Docker Compose.

## 🏗️ Architecture

- **API**: Go-based REST API with gRPC support
- **Web**: Next.js frontend application
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Brokers**: Kafka and NATS
- **Containerization**: Docker & Docker Compose

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Make (optional, for convenience commands)

### Running the Project

1. **Clone and navigate to the project:**
   ```bash
   cd /Users/morshedulmunna/Developer/signal_copier
   ```

2. **Start all services:**
   ```bash
   make start
   ```
   Or manually:
   ```bash
   docker-compose up -d
   ```

3. **Access the applications:**
   - **Web Application**: http://localhost:3000
   - **API**: http://localhost:9090
   - **API Health Check**: http://localhost:9090/health

## 📋 Available Commands

### Docker Compose Commands

```bash
# Build all images
make build
# or
docker-compose build

# Start all services
make up
# or
docker-compose up -d

# Stop all services
make down
# or
docker-compose down

# View logs
make logs
# or
docker-compose logs -f

# Restart services
make restart
# or
docker-compose restart
```

### Development Commands

```bash
# Start in development mode (with live reload)
make dev

# Start only API with dependencies
make dev-api

# Start only Web
make dev-web
```

### Database Commands

```bash
# Run migrations
make migrate
# or
docker-compose exec api /app/main migrate
```

### Testing Commands

```bash
# Run all tests
make test

# Run unit tests
make test-unit

# Run integration tests
make test-integration

# Run end-to-end tests
make test-e2e
```

### Health & Status

```bash
# Check service status
make status

# Check health of all services
make health
```

### Cleanup Commands

```bash
# Clean up containers and volumes
make clean

# Remove all project images
make clean-images
```

## 🗄️ Services

| Service | Port | Description |
|---------|------|-------------|
| Web (Next.js) | 3000 | Frontend application |
| API (Go) | 9090 | REST API server |
| API (gRPC) | 50001 | gRPC server |
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache |
| Kafka | 9092 | Message broker |
| NATS | 4222 | Message broker |
| NATS HTTP | 8222 | NATS monitoring |

## 🔧 Configuration

### Environment Variables

The project uses environment variables for configuration:

- **API**: `/api/.env` (automatically created from `example.env`)
- **Web**: `/web/.env` (automatically created from `example.env`)

Key configurations:
- Database connection settings
- Redis cache settings
- Message broker configurations
- JWT secrets
- Service ports

### Docker Services

All services are configured to work together:
- Services communicate using Docker network
- Health checks ensure proper startup order
- Volumes persist data across restarts

## 🏃‍♂️ Development

### Project Structure

```
signal_copier/
├── api/                 # Go API backend
│   ├── cmd/            # CLI commands
│   ├── internal/       # Internal packages
│   ├── config/         # Configuration
│   ├── database/       # Database models and migrations
│   └── proto/          # gRPC protocol buffers
├── web/                # Next.js frontend
│   ├── src/            # Source code
│   ├── components/     # React components
│   └── app/            # Next.js app directory
├── docker-compose.yml  # Docker Compose configuration
└── Makefile           # Convenience commands
```

### Adding New Services

1. Add service definition to `docker-compose.yml`
2. Update environment variables if needed
3. Add health checks for proper startup order
4. Update Makefile with new commands if needed

### Debugging

- View logs: `make logs` or `docker-compose logs -f [service]`
- Check health: `make health`
- Access service shell: `docker-compose exec [service] sh`

## 🐛 Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 3000, 9090, 5432, 6379, 9092, 4222 are available
2. **Permission issues**: Ensure Docker has proper permissions
3. **Build failures**: Check Dockerfile syntax and dependencies
4. **Service not starting**: Check logs with `make logs [service]`

### Reset Everything

```bash
make clean
make start
```

## 📝 Notes

- The project uses Go 1.24.2 and Node.js 20
- All services include health checks for reliable startup
- Data persistence is handled through Docker volumes
- The API supports both REST and gRPC protocols
- Message brokers (Kafka and NATS) are configured for event-driven architecture
