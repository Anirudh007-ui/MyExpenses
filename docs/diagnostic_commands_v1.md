# MyExpenses Project: Diagnostic & Execution Command Guide

This document records all commands executed to run, diagnose, and fix issues in the MyExpenses Go project. Each command is explained with **why**, **when**, **what** it does, and **expected output**. Use this as a learning and troubleshooting reference for future development.

---

## 1. Initial Setup & Dependency Management

### 1.1. Download Go Dependencies
```bash
go mod tidy
```
- **Why:** Ensures all required Go modules are downloaded and the `go.mod`/`go.sum` files are up to date.
- **When:** After cloning the repo or adding new dependencies.
- **Expected Output:** No errors if all dependencies are available. Warnings if any modules are missing.

---

## 2. Database Setup (PostgreSQL with Docker)

### 2.1. Start Docker (if not running)
```bash
open -a Docker
```
- **Why:** Docker must be running to start containers.
- **When:** Before running any Docker commands on macOS.
- **Expected Output:** No output if Docker starts successfully.

### 2.2. Start PostgreSQL Container
```bash
docker run -d --name postgres-expenses \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=myexpenses \
  -p 5432:5432 \
  postgres:15-alpine
```
- **Why:** Launches a PostgreSQL database for the app to connect to.
- **When:** Before running the Go application.
- **Expected Output:** Container ID if successful. If image is missing, Docker will pull it.

### 2.3. Fixing Database Not Found Errors
If you see:
```
FATAL: database "myexpenses" does not exist (SQLSTATE 3D000)
```
Run:
```bash
docker rm -f postgres-expenses
docker volume prune -f
```
- **Why:** Removes the old container and its data volume to ensure a fresh start.
- **When:** If the database was not created or the container is corrupted.
- **Expected Output:** Confirmation of container and volume removal.

Then, restart the container as in 2.2.

---

## 3. Running the Go Application

### 3.1. Run the Application
```bash
go run cmd/api/main.go
```
- **Why:** Starts the API server.
- **When:** After database is running and dependencies are installed.
- **Expected Output:**
  - Success: Logs indicating server is running (e.g., `Starting server on port 8080`).
  - Failure: Error messages about database connection, missing files, or code issues.

### 3.2. Common Error: Empty or Invalid Go File
If you see:
```
internal/db/db.go:1:1: expected 'package', found 'EOF'
```
- **Why:** An empty or invalid Go file exists (here, `internal/db/db.go`).
- **Fix:** Delete or correct the file.
```bash
rm internal/db/db.go
```
- **Expected Output:** File is removed. Application should now start if no other errors.

---

## 4. Diagnosing Application Status

### 4.1. Check if Application is Running
```bash
ps aux | grep "go run" | grep -v grep
```
- **Why:** Verifies if the Go application is running in the background.
- **When:** If you're unsure whether the server is active.
- **Expected Output:** Process details if running, nothing if not.

### 4.2. Check Health Endpoint
```bash
curl -s http://localhost:8080/health
```
- **Why:** Confirms the API server is up and responding.
- **When:** After starting the server.
- **Expected Output:** JSON response like `{ "status": "ok", "service": "MyExpenses API" }`.

---

## 5. Testing API Endpoints

### 5.1. Create an Expense
```bash
curl -X POST http://localhost:8080/expenses \
  -H "Content-Type: application/json" \
  -d '{"description":"Coffee","amount":4.50,"category":"Food","date":"2024-01-15T08:00:00Z"}'
```
- **Why:** Tests the POST endpoint for creating expenses.
- **Expected Output:** JSON with the created expense and a success message.

### 5.2. Get All Expenses
```bash
curl http://localhost:8080/expenses
```
- **Why:** Tests the GET endpoint for retrieving all expenses.
- **Expected Output:** JSON array of expenses.

### 5.3. Get Expense by ID
```bash
curl http://localhost:8080/expenses/{id}
```
- **Why:** Tests the GET endpoint for a specific expense.
- **Expected Output:** JSON object for the expense with the given ID.

### 5.4. Update an Expense
```bash
curl -X PUT http://localhost:8080/expenses/{id} \
  -H "Content-Type: application/json" \
  -d '{"amount":5.00}'
```
- **Why:** Tests the PUT endpoint for updating an expense.
- **Expected Output:** JSON with the updated expense and a success message.

### 5.5. Delete an Expense
```bash
curl -X DELETE http://localhost:8080/expenses/{id}
```
- **Why:** Tests the DELETE endpoint for removing an expense.
- **Expected Output:** JSON with a success message.

---

## 6. General Troubleshooting

- **Check Docker is running:**
  - If you see `Cannot connect to the Docker daemon`, start Docker Desktop.
- **Check for empty or invalid Go files:**
  - Remove or fix files that cause `expected 'package', found 'EOF'` errors.
- **Check logs for errors:**
  - Application logs will indicate missing dependencies, database issues, or code errors.
- **Remove and recreate containers if needed:**
  - Use `docker rm -f <container>` and `docker volume prune -f` to reset state.

---

## 7. Future Diagnostics

- **Document every new command and its purpose.**
- **Add new troubleshooting steps as you encounter new issues.**
- **Keep this file updated for team learning and onboarding.**

---

# 📚 Summary

This guide is a living document. Update it with every new command, error, and fix you encounter. It will help you and your team diagnose, fix, and understand the MyExpenses project efficiently. 

# Diagnostic Commands Guide

This guide provides comprehensive diagnostic commands for troubleshooting the MyExpenses API service. Each command includes detailed explanations of when to use it, what it does, and what outcomes to expect.

## Table of Contents

1. [Database Diagnostics](#database-diagnostics)
2. [Application Diagnostics](#application-diagnostics)
3. [Network Diagnostics](#network-diagnostics)
4. [Docker Diagnostics](#docker-diagnostics)
5. [Go Application Diagnostics](#go-application-diagnostics)
6. [Error Resolution Workflows](#error-resolution-workflows)

---

## Database Diagnostics

### 1. Check PostgreSQL Container Status

**Command:**
```bash
docker ps | grep postgres
```

**When to use:**
- When the application fails to connect to the database
- After starting or restarting the PostgreSQL container
- To verify the container is running and healthy

**What it does:**
- Lists all running Docker containers
- Filters to show only PostgreSQL-related containers
- Shows container ID, image, status, and port mappings

**Expected outcome:**
```
8531d45f49c0   postgres:15-alpine   "docker-entrypoint.s…"   5 minutes ago   Up 5 minutes   0.0.0.0:5432->5432/tcp   postgres-expenses
```

**Troubleshooting:**
- If no output: Container is not running
- If status shows "Exited": Container crashed or stopped
- If port mapping is missing: Container not properly configured

---

### 2. List All PostgreSQL Containers (Including Stopped)

**Command:**
```bash
docker ps -a | grep postgres
```

**When to use:**
- When `docker ps` shows no containers
- To see if containers exist but are stopped
- Before cleaning up old containers

**What it does:**
- Lists all Docker containers (running and stopped)
- Shows container status, exit codes, and creation time

**Expected outcome:**
```
8531d45f49c0   postgres:15-alpine   "docker-entrypoint.s…"   7 minutes ago   Up 7 minutes   0.0.0.0:5432->5432/tcp   postgres-expenses
```

---

### 3. Check Database Existence (From Container)

**Command:**
```bash
docker exec postgres-expenses psql -U postgres -l
```

**When to use:**
- When application reports "database does not exist"
- After container startup to verify database creation
- To see all available databases

**What it does:**
- Connects to PostgreSQL from inside the container
- Lists all databases with their owners and encodings
- Shows database privileges

**Expected outcome:**
```
                                                 List of databases
    Name    |  Owner   | Encoding |  Collate   |   Ctype    | ICU Locale | Locale Provider |   Access privileges
------------+----------+----------+------------+------------+------------+-----------------+----------------
 myexpenses | postgres | UTF8     | en_US.utf8 | en_US.utf8 |            | libc            | 
 postgres   | postgres | UTF8     | en_US.utf8 | en_US.utf8 |            | libc            | 
 template0  | postgres | UTF8     | en_US.utf8 | en_US.utf8 |            | libc            | =c/postgres
 template1  | postgres | UTF8     | en_US.utf8 | en_US.utf8 |            | libc            | =c/postgres
```

**Troubleshooting:**
- If `myexpenses` database is missing: Container initialization failed
- If no databases listed: PostgreSQL not properly started

---

### 4. Test Database Connection (From Container)

**Command:**
```bash
docker exec postgres-expenses psql -U postgres -d myexpenses -c "SELECT 1;"
```

**When to use:**
- To verify database connectivity from within container
- After database creation to test basic functionality
- To check if PostgreSQL is responding to queries

**What it does:**
- Connects to the `myexpenses` database
- Executes a simple SELECT query
- Tests basic database functionality

**Expected outcome:**
```
 ?column? 
----------
        1
(1 row)
```

**Troubleshooting:**
- If connection fails: Database doesn't exist or permissions issue
- If query fails: PostgreSQL service issues

---

### 5. Check Database Tables

**Command:**
```bash
docker exec postgres-expenses psql -U postgres -d myexpenses -c "\dt"
```

**When to use:**
- After application startup to verify table creation
- When AutoMigrate fails or tables are missing
- To see all tables in the database

**What it does:**
- Lists all tables in the `myexpenses` database
- Shows table names, types, and owners

**Expected outcome:**
```
          List of relations
 Schema |   Name   | Type  |  Owner   
--------+----------+-------+----------
 public | expenses | table | postgres
```

**Troubleshooting:**
- If no tables listed: AutoMigrate failed or not run
- If `expenses` table missing: Application initialization issue

---

### 6. Check Table Schema

**Command:**
```bash
docker exec postgres-expenses psql -U postgres -d myexpenses -c "\d expenses"
```

**When to use:**
- To verify table structure is correct
- When experiencing data type issues
- To check if all required columns exist

**What it does:**
- Shows detailed table schema
- Lists all columns with data types and constraints
- Shows indexes and primary keys

**Expected outcome:**
```
                              Table "public.expenses"
   Column    |           Type           | Collation | Nullable |      Default      
-------------+--------------------------+-----------+----------+-------------------
 id          | uuid                     |           | not null | gen_random_uuid()
 description | text                     |           | not null | 
 amount      | numeric                  |           | not null | 
 category    | text                     |           | not null | 
 date        | timestamp with time zone |           | not null | 
 created_at  | timestamp with time zone |           |          | 
 updated_at  | timestamp with time zone |           |          | 
Indexes:
    "expenses_pkey" PRIMARY KEY, btree (id)
```

---

### 7. Test Database Connection (From Host)

**Command:**
```bash
PGPASSWORD=password psql -h localhost -p 5432 -U postgres -l
```

**When to use:**
- To test connectivity from the host machine
- When application can't connect to database
- To verify port forwarding is working

**What it does:**
- Connects to PostgreSQL from the host machine
- Uses environment variable for password
- Lists all databases

**Expected outcome:**
```
                                 List of databases
    Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
------------+----------+----------+------------+------------+-----------------------
 myexpenses | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 postgres   | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres
 template1  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres
```

**Troubleshooting:**
- If connection fails: Port forwarding issue or wrong PostgreSQL instance
- If `myexpenses` missing: Database not created properly

---

### 8. Test Specific Database Connection (From Host)

**Command:**
```bash
PGPASSWORD=password psql -h localhost -p 5432 -U postgres -d myexpenses -c "SELECT 1;"
```

**When to use:**
- To test specific database connectivity
- When application reports database connection errors
- To verify the exact database the application should connect to

**What it does:**
- Connects to the specific `myexpenses` database
- Executes a test query
- Verifies database is accessible

**Expected outcome:**
```
 ?column? 
----------
        1
(1 row)
```

**Troubleshooting:**
- If fails with "database does not exist": Wrong database or container issue
- If connection refused: Port or network issue

---

## Application Diagnostics

### 9. Check Application Process Status

**Command:**
```bash
ps aux | grep "go run" | grep -v grep
```

**When to use:**
- When application seems to have stopped
- After starting application in background
- To verify application is still running

**What it does:**
- Lists all running processes
- Filters for Go application processes
- Excludes the grep command itself

**Expected outcome:**
```
ani  6081  0.0  0.0  1234567  1234 pts/0    S    00:10   0:00 go run cmd/api/main.go
```

**Troubleshooting:**
- If no output: Application is not running
- If multiple processes: Multiple instances running

---

### 10. Check Background Jobs

**Command:**
```bash
jobs
```

**When to use:**
- When application was started in background
- To see status of background processes
- To manage background jobs

**What it does:**
- Lists all background jobs in current shell
- Shows job status (running, stopped, etc.)

**Expected outcome:**
```
[1]  + running    go run cmd/api/main.go
```

**Troubleshooting:**
- If no output: No background jobs in current shell
- If job shows "stopped": Application crashed or was stopped

---

### 11. Test Health Endpoint

**Command:**
```bash
curl -s http://localhost:8080/health
```

**When to use:**
- To verify application is running and responding
- After application startup
- To test basic HTTP connectivity

**What it does:**
- Makes HTTP GET request to health endpoint
- Returns application status
- Tests basic functionality

**Expected outcome:**
```json
{"service":"MyExpenses API","status":"ok"}
```

**Troubleshooting:**
- If connection refused: Application not running or wrong port
- If 404: Health endpoint not configured
- If 500: Application error

---

### 12. Test Health Endpoint with Pretty Output

**Command:**
```bash
curl -s http://localhost:8080/health | jq .
```

**When to use:**
- When you want formatted JSON output
- For better readability of API responses
- When debugging API responses

**What it does:**
- Makes HTTP request and formats JSON response
- Uses `jq` for JSON pretty printing

**Expected outcome:**
```json
{
  "service": "MyExpenses API",
  "status": "ok"
}
```

---

## Network Diagnostics

### 13. Check Port Usage

**Command:**
```bash
lsof -i :5432
```

**When to use:**
- When database connection fails
- To see what's using port 5432
- To identify port conflicts

**What it does:**
- Lists all processes using port 5432
- Shows process details and connection types
- Identifies potential conflicts

**Expected outcome:**
```
COMMAND    PID USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
com.docke 4990  ani   73u  IPv6 0x5da1bfce2596147a      0t0  TCP *:postgresql (LISTEN)
```

**Troubleshooting:**
- If multiple processes: Port conflict
- If no Docker process: Container not running
- If local PostgreSQL: Potential conflict

---

### 14. Check Application Port Usage

**Command:**
```bash
lsof -i :8080
```

**When to use:**
- When application won't start
- To see if port 8080 is already in use
- To identify port conflicts

**What it does:**
- Lists all processes using port 8080
- Shows which application is using the port

**Expected outcome:**
```
COMMAND    PID USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
go        6081  ani    3u  IPv6 0x9c75ad060923e8f5      0t0  TCP localhost:8080 (LISTEN)
```

---

## Docker Diagnostics

### 15. Check Container Logs

**Command:**
```bash
docker logs postgres-expenses
```

**When to use:**
- When container fails to start
- To see database initialization process
- To debug container startup issues

**What it does:**
- Shows all logs from the PostgreSQL container
- Displays startup sequence and errors
- Shows database creation process

**Expected outcome:**
```
The files belonging to this database system will be owned by user "postgres".
This user must also own the server process.

The database cluster will be initialized with locale "en_US.utf8".
...
CREATE DATABASE
...
PostgreSQL init process complete; ready for start up.
...
database system is ready to accept connections
```

**Troubleshooting:**
- If no "CREATE DATABASE": Database not created
- If "init process complete" missing: Container failed to initialize
- If errors: Configuration or resource issues

---

### 16. Check Container Resource Usage

**Command:**
```bash
docker stats postgres-expenses
```

**When to use:**
- When container seems slow or unresponsive
- To monitor resource consumption
- To identify performance issues

**What it does:**
- Shows real-time resource usage
- Displays CPU, memory, network, and disk usage

**Expected outcome:**
```
CONTAINER ID   NAME              CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O         PIDS
22f369be2342   postgres-expenses  0.00%     15.23MiB / 2GiB      0.74%     1.23kB / 2.45kB   0B / 0B           7
```

---

## Go Application Diagnostics

### 17. Run Application with Error Output

**Command:**
```bash
go run cmd/api/main.go
```

**When to use:**
- To start the application
- To see detailed error messages
- To debug startup issues

**What it does:**
- Compiles and runs the Go application
- Shows all log output and errors
- Displays startup sequence

**Expected outcome:**
```
2025/06/26 00:16:34 No .env file found, using system environment variables
2025/06/26 00:16:34 Successfully connected to PostgreSQL database
2025/06/26 00:16:34 Starting server on port 8080
```

**Troubleshooting:**
- If database connection fails: Check database setup
- If port binding fails: Check port availability
- If compilation errors: Check Go code

---

### 18. Run Application in Background

**Command:**
```bash
go run cmd/api/main.go &
```

**When to use:**
- To start application without blocking terminal
- When you need to run other commands
- For testing scenarios

**What it does:**
- Starts application in background
- Returns control to terminal
- Shows process ID

**Expected outcome:**
```
[1] 6081
```

---

### 19. Check Go Module Dependencies

**Command:**
```bash
go mod tidy
```

**When to use:**
- When dependencies are missing
- After adding new imports
- To clean up unused dependencies

**What it does:**
- Downloads missing dependencies
- Removes unused dependencies
- Updates go.mod and go.sum files

**Expected outcome:**
```
go: downloading github.com/gin-gonic/gin v1.9.1
go: downloading gorm.io/gorm v1.25.5
...
```

---

## Error Resolution Workflows

### Workflow 1: Database Connection Issues

**Symptoms:**
- Application fails with "database does not exist"
- Connection refused errors
- Port conflicts

**Diagnostic Steps:**
1. Check container status: `docker ps | grep postgres`
2. Check port conflicts: `lsof -i :5432`
3. Test container connectivity: `docker exec postgres-expenses psql -U postgres -l`
4. Test host connectivity: `PGPASSWORD=password psql -h localhost -p 5432 -U postgres -l`
5. Check container logs: `docker logs postgres-expenses`

**Common Solutions:**
- Stop conflicting local PostgreSQL: `brew services stop postgresql`
- Restart container: `docker restart postgres-expenses`
- Recreate container with proper environment variables

---

### Workflow 2: Application Startup Issues

**Symptoms:**
- Application won't start
- Port already in use
- Compilation errors

**Diagnostic Steps:**
1. Check port usage: `lsof -i :8080`
2. Check Go dependencies: `go mod tidy`
3. Run application directly: `go run cmd/api/main.go`
4. Check for background processes: `ps aux | grep "go run"`

**Common Solutions:**
- Kill conflicting processes
- Fix compilation errors
- Update dependencies

---

### Workflow 3: API Endpoint Issues

**Symptoms:**
- Health endpoint not responding
- CRUD operations failing
- Database table issues

**Diagnostic Steps:**
1. Test health endpoint: `curl -s http://localhost:8080/health`
2. Check database tables: `docker exec postgres-expenses psql -U postgres -d myexpenses -c "\dt"`
3. Test database connectivity: `docker exec postgres-expenses psql -U postgres -d myexpenses -c "SELECT 1;"`
4. Check application logs for errors

**Common Solutions:**
- Restart application
- Recreate database tables
- Fix database schema issues

---

## Quick Reference Commands

### Essential Commands for Daily Use

```bash
# Start PostgreSQL container
docker run -d --name postgres-expenses -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=myexpenses -p 5432:5432 postgres:15-alpine

# Start application
go run cmd/api/main.go

# Test health
curl -s http://localhost:8080/health | jq .

# Check database
docker exec postgres-expenses psql -U postgres -d myexpenses -c "\dt"

# Stop and clean up
docker stop postgres-expenses && docker rm postgres-expenses
```

### Emergency Commands

```bash
# Force stop all containers
docker stop $(docker ps -q)

# Remove all containers
docker rm $(docker ps -aq)

# Clean up volumes
docker volume prune -f

# Kill all Go processes
pkill -f "go run"
```

---

## Best Practices

1. **Always check container status first** before troubleshooting application issues
2. **Use specific database names** to avoid confusion with other PostgreSQL instances
3. **Test connectivity from both container and host** to isolate network issues
4. **Check logs immediately** when containers or applications fail to start
5. **Use pretty-printed JSON** (`jq`) for better debugging of API responses
6. **Keep diagnostic commands handy** for quick troubleshooting
7. **Document any custom configurations** that might affect the setup

---

## Common Error Messages and Solutions

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `database "myexpenses" does not exist` | Database not created or wrong PostgreSQL instance | Check container logs, recreate container |
| `connection refused` | Port conflict or service not running | Check port usage, restart service |
| `port already in use` | Another service using the port | Kill conflicting process or change port |
| `compilation failed` | Missing dependencies or syntax errors | Run `go mod tidy`, fix code |
| `table does not exist` | AutoMigrate failed | Check application logs, restart application |

This diagnostic guide should help you quickly identify and resolve most issues with the MyExpenses API service. 