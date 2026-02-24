# Signal Copier

A robust Go-based REST API for managing trading signals, backed by PostgreSQL.

## 🏗️ Architecture

- **Backend**: Go (Golang)
- **Database**: PostgreSQL (with GORM Auto-migrations)
- **Containerization**: Docker & Docker Compose

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Make (optional)

### Running the Project

1. **Clone the repository and navigate to the project:**
   ```bash
   git clone <repository-url>
   cd Copier
   ```

2. **Start the services:**
   ```bash
   make start
   ```
   *This will build the API image and start both the API and Postgres containers.*

3. **Verify Installation:**
   - **API Health Check**: [http://localhost:9090/health](http://localhost:9090/health)
   - **API Documentation**: (If available, usually at [http://localhost:9090/api/v1](http://localhost:9090/api/v1))

## 📋 Available Commands

### Docker Lifecycle
- `make build`: Build Docker images
- `make up`: Start all services in detached mode
- `make down`: Stop and remove containers
- `make restart`: Restart all services
- `make logs`: Follow logs for all services

### Development & Database
- `make dev`: Start services in the terminal (not detached)
- `make migrate`: Run database auto-migrations
- `make seed`: Seed the database with initial data
- `make clean`: Deep clean (removes volumes and orphans)

### Testing
- `make test`: Run all tests
- `make test-unit`: Run unit tests
- `make test-integration`: Run integration tests

## 🔧 Configuration

The project uses environment variables. The API looks for a `.env` file in the `api/` directory.

Key configurations include:
- `HTTP_PORT`: Port for the REST API (default: 9090)
- `DB_*`: PostgreSQL connection settings
- `JWT_SECRET`: Secret key for authentication

## 📝 Notes

- **Auto-migrations**: The API automatically runs GORM auto-migrations on startup.
- **Port Mapping**:
  - API: 9090
  - Postgres: 5433 (mapped to 5432 internally)
