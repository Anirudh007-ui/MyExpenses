# MyExpenses Troubleshooting Guide

This document provides solutions to common issues encountered while developing, running, and testing the MyExpenses application.

## 🚨 Common Issues & Solutions

### 1. Go Compilation Errors

#### Issue: Empty or Invalid Go File
```
internal/db/db.go:1:1: expected 'package', found 'EOF'
```

**Cause:** An empty or malformed Go file exists.

**Solution:**
```bash
# Remove the empty file
rm internal/db/db.go

# Or check if the file has proper Go syntax
cat internal/db/db.go
```

**Prevention:** Always ensure Go files have proper package declarations and syntax.

---

#### Issue: Missing Dependencies
```
cannot find package "github.com/gin-gonic/gin"
```

**Cause:** Dependencies not downloaded or `go.mod` file corrupted.

**Solution:**
```bash
# Download dependencies
go mod tidy

# If that fails, try cleaning and re-downloading
go clean -modcache
go mod download
```

**Prevention:** Run `go mod tidy` after adding new imports.

---

### 2. Database Connection Issues

#### Issue: Database Not Found
```
FATAL: database "myexpenses" does not exist (SQLSTATE 3D000)
```

**Cause:** PostgreSQL container was started without the correct environment variables or database wasn't created.

**Solution:**
```bash
# Stop and remove the old container
docker rm -f postgres-expenses

# Remove any old volumes
docker volume prune -f

# Start a fresh container with correct environment
docker run -d --name postgres-expenses \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=myexpenses \
  -p 5432:5432 \
  postgres:15-alpine

# Wait for database to be ready
sleep 5
```

**Prevention:** Always use the correct environment variables when starting PostgreSQL.

---

#### Issue: Connection Refused
```
failed to connect to database: dial tcp 127.0.0.1:5432: connect: connection refused
```

**Cause:** PostgreSQL container is not running or port is not exposed.

**Solution:**
```bash
# Check if container is running
docker ps | grep postgres

# If not running, start it
docker start postgres-expenses

# If container doesn't exist, create it
docker run -d --name postgres-expenses \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=myexpenses \
  -p 5432:5432 \
  postgres:15-alpine
```

**Prevention:** Use Docker Compose for easier container management.

---

#### Issue: Authentication Failed
```
FATAL: password authentication failed for user "postgres"
```

**Cause:** Wrong username/password combination.

**Solution:**
```bash
# Check environment variables in your application
echo $DB_USER
echo $DB_PASSWORD

# Ensure they match the PostgreSQL container settings
# Default: postgres/password
```

**Prevention:** Use consistent environment variables across development and production.

---

### 3. Docker Issues

#### Issue: Docker Daemon Not Running
```
Cannot connect to the Docker daemon at unix:///Users/ani/.docker/run/docker.sock
```

**Cause:** Docker Desktop is not running.

**Solution:**
```bash
# Start Docker Desktop (macOS)
open -a Docker

# Wait for Docker to start
sleep 10

# Verify Docker is running
docker ps
```

**Prevention:** Start Docker Desktop before running any Docker commands.

---

#### Issue: Port Already in Use
```
Error response from daemon: driver failed programming external connectivity on endpoint postgres-expenses: Error starting userland proxy: listen tcp 0.0.0.0:5432: bind: address already in use
```

**Cause:** Port 5432 is already occupied by another PostgreSQL instance.

**Solution:**
```bash
# Check what's using port 5432
lsof -i :5432

# Stop the conflicting service
sudo brew services stop postgresql  # If using Homebrew PostgreSQL
# OR
sudo systemctl stop postgresql      # If using system PostgreSQL

# Or use a different port
docker run -d --name postgres-expenses \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=myexpenses \
  -p 5433:5432 \  # Changed to port 5433
  postgres:15-alpine
```

**Prevention:** Use Docker Compose to manage port conflicts.

---

### 4. Application Startup Issues

#### Issue: Application Won't Start
```
Failed to start server: listen tcp :8080: bind: address already in use
```

**Cause:** Port 8080 is already in use by another application.

**Solution:**
```bash
# Check what's using port 8080
lsof -i :8080

# Kill the process using port 8080
kill -9 <PID>

# Or use a different port
export PORT=8081
go run cmd/api/main.go
```

**Prevention:** Use environment variables for port configuration.

---

#### Issue: Environment Variables Not Loaded
```
No .env file found, using system environment variables
```

**Cause:** `.env` file is missing or not in the correct location.

**Solution:**
```bash
# Create .env file in project root
cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=myexpenses
DB_SSLMODE=disable
PORT=8080
EOF
```

**Prevention:** Always include a `.env.example` file in your repository.

---

### 5. API Testing Issues

#### Issue: Connection Refused on API Calls
```
curl: (7) Failed to connect to localhost port 8080: Connection refused
```

**Cause:** Application is not running or not listening on the expected port.

**Solution:**
```bash
# Check if application is running
ps aux | grep "go run" | grep -v grep

# If not running, start it
go run cmd/api/main.go &

# Wait for startup
sleep 3

# Test health endpoint
curl http://localhost:8080/health
```

**Prevention:** Always check application status before testing endpoints.

---

#### Issue: Invalid JSON Response
```
curl: (3) URL using bad/illegal format or missing URL
```

**Cause:** Malformed curl command or JSON.

**Solution:**
```bash
# Use proper JSON formatting
curl -X POST http://localhost:8080/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Coffee",
    "amount": 4.50,
    "category": "Food",
    "date": "2024-01-15T08:00:00Z"
  }'

# Or use a JSON file
echo '{"description":"Coffee","amount":4.50,"category":"Food","date":"2024-01-15T08:00:00Z"}' > test.json
curl -X POST http://localhost:8080/expenses \
  -H "Content-Type: application/json" \
  -d @test.json
```

**Prevention:** Use proper JSON formatting and validate with tools like `jq`.

---

### 6. Database Migration Issues

#### Issue: Table Already Exists
```
ERROR: relation "expenses" already exists
```

**Cause:** Database table was already created in a previous run.

**Solution:**
```bash
# Drop and recreate the database
docker exec -it postgres-expenses psql -U postgres -c "DROP DATABASE myexpenses;"
docker exec -it postgres-expenses psql -U postgres -c "CREATE DATABASE myexpenses;"

# Or restart the container
docker restart postgres-expenses
```

**Prevention:** Use proper migration tools for production environments.

---

#### Issue: UUID Extension Not Available
```
ERROR: function gen_random_uuid() does not exist
```

**Cause:** PostgreSQL UUID extension not enabled.

**Solution:**
```bash
# Enable UUID extension
docker exec -it postgres-expenses psql -U postgres -d myexpenses -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
```

**Prevention:** Include extension creation in your migration scripts.

---

## 🔧 Debugging Techniques

### 1. Application Logs
```bash
# Run application with verbose logging
GIN_MODE=debug go run cmd/api/main.go

# Check application logs
tail -f /var/log/myexpenses.log  # If logging to file
```

### 2. Database Debugging
```bash
# Connect to PostgreSQL container
docker exec -it postgres-expenses psql -U postgres -d myexpenses

# Check tables
\dt

# Check table structure
\d expenses

# Query data
SELECT * FROM expenses;
```

### 3. Network Debugging
```bash
# Check if ports are listening
netstat -tulpn | grep :8080
netstat -tulpn | grep :5432

# Test database connection
telnet localhost 5432

# Test API connection
telnet localhost 8080
```

### 4. Docker Debugging
```bash
# Check container status
docker ps -a

# Check container logs
docker logs postgres-expenses

# Check container resources
docker stats postgres-expenses

# Inspect container configuration
docker inspect postgres-expenses
```

---

## 🛠️ Maintenance Commands

### Database Maintenance
```bash
# Backup database
docker exec postgres-expenses pg_dump -U postgres myexpenses > backup.sql

# Restore database
docker exec -i postgres-expenses psql -U postgres myexpenses < backup.sql

# Reset database
docker exec postgres-expenses psql -U postgres -c "DROP DATABASE myexpenses; CREATE DATABASE myexpenses;"
```

### Application Maintenance
```bash
# Clean Go cache
go clean -cache -modcache

# Update dependencies
go get -u ./...

# Run tests
go test ./...

# Build application
go build -o myexpenses cmd/api/main.go
```

### Docker Maintenance
```bash
# Clean up unused containers
docker container prune

# Clean up unused images
docker image prune

# Clean up unused volumes
docker volume prune

# Clean up everything
docker system prune -a
```

---

## 📋 Pre-flight Checklist

Before running the application, ensure:

- [ ] Docker Desktop is running
- [ ] PostgreSQL container is running and healthy
- [ ] Database `myexpenses` exists
- [ ] Environment variables are set correctly
- [ ] Port 8080 is available
- [ ] All Go dependencies are downloaded
- [ ] No syntax errors in Go files

### Quick Health Check
```bash
# 1. Check Docker
docker ps

# 2. Check database
docker exec postgres-expenses psql -U postgres -d myexpenses -c "SELECT 1;"

# 3. Check application
curl -s http://localhost:8080/health

# 4. Check dependencies
go mod verify
```

---

## 🆘 Getting Help

### Logs to Collect
When reporting issues, include:

1. **Application logs** - Full startup and error logs
2. **Database logs** - PostgreSQL container logs
3. **System information** - OS, Go version, Docker version
4. **Environment** - Environment variables (without sensitive data)
5. **Steps to reproduce** - Exact commands and sequence

### Useful Commands for Debugging
```bash
# System information
go version
docker --version
uname -a

# Application status
ps aux | grep go
netstat -tulpn | grep :8080

# Database status
docker ps | grep postgres
docker logs postgres-expenses

# Network connectivity
curl -v http://localhost:8080/health
telnet localhost 5432
```

This troubleshooting guide should help you resolve most common issues. If you encounter a new problem, document it here for future reference. 