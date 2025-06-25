# MyExpenses API

A scalable REST API for managing personal expenses built with Go using Clean Architecture principles.

## Architecture

This project follows **Clean Architecture** with the following layers:

- **Domain Layer**: Business entities and interfaces
- **Application Layer**: Use cases and business logic
- **Infrastructure Layer**: Database, HTTP handlers, external services
- **Presentation Layer**: HTTP routes and middleware

## Features

- ✅ CRUD operations for expenses
- ✅ Advanced filtering and search
- ✅ PostgreSQL database with GORM
- ✅ RESTful API design
- ✅ Input validation
- ✅ Error handling
- ✅ Docker support
- ✅ Environment configuration

## API Endpoints

### POST /expenses
Create a new expense.

**Request Body:**
```json
{
  "description": "Grocery shopping",
  "amount": 45.50,
  "category": "Food",
  "date": "2024-01-15T10:30:00Z"
}
```

**Response:**
```json
{
  "message": "Expense created successfully",
  "data": {
    "id": "uuid-here",
    "description": "Grocery shopping",
    "amount": 45.50,
    "category": "Food",
    "date": "2024-01-15T10:30:00Z",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### GET /expenses
Get all expenses with optional filtering.

**Query Parameters:**
- `category` - Filter by category (partial match)
- `date_from` - Filter expenses from this date (YYYY-MM-DD)
- `date_to` - Filter expenses until this date (YYYY-MM-DD)
- `min_amount` - Minimum amount filter
- `max_amount` - Maximum amount filter
- `description` - Filter by description (partial match)

**Example:**
```
GET /expenses?category=Food&min_amount=10&date_from=2024-01-01
```

### GET /expenses/{id}
Get a specific expense by ID.

### PUT /expenses/{id}
Update an existing expense.

**Request Body:**
```json
{
  "description": "Updated grocery shopping",
  "amount": 50.00,
  "category": "Food & Dining",
  "date": "2024-01-15T10:30:00Z"
}
```

### DELETE /expenses/{id}
Delete an expense.

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker (optional)

### Local Development

1. **Clone the repository:**
```bash
git clone <your-repo-url>
cd MyExpenses
```

2. **Install dependencies:**
```bash
go mod tidy
```

3. **Set up environment variables:**
Create a `.env` file:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=myexpenses
DB_SSLMODE=disable
PORT=8080
```

4. **Start PostgreSQL:**
```bash
# Using Docker
docker run -d \
  --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=myexpenses \
  -p 5432:5432 \
  postgres:15-alpine
```

5. **Run the application:**
```bash
go run cmd/api/main.go
```

### Using Docker Compose

1. **Start all services:**
```bash
docker-compose up -d
```

2. **View logs:**
```bash
docker-compose logs -f app
```

3. **Stop services:**
```bash
docker-compose down
```

## Testing the API

### Create an expense:
```bash
curl -X POST http://localhost:8080/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Coffee",
    "amount": 4.50,
    "category": "Food",
    "date": "2024-01-15T08:00:00Z"
  }'
```

### Get all expenses:
```bash
curl http://localhost:8080/expenses
```

### Get expenses with filters:
```bash
curl "http://localhost:8080/expenses?category=Food&min_amount=5"
```

### Update an expense:
```bash
curl -X PUT http://localhost:8080/expenses/{expense-id} \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5.00
  }'
```

### Delete an expense:
```bash
curl -X DELETE http://localhost:8080/expenses/{expense-id}
```

## Project Structure

```
MyExpenses/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── db/
│   │   └── postgres.go            # Database configuration
│   └── expenses/
│       ├── domain/                # Domain layer
│       │   ├── expense.go         # Expense entity
│       │   ├── errors.go          # Domain errors
│       │   └── repository.go      # Repository interface
│       ├── application/           # Application layer
│       │   └── service.go         # Business logic
│       └── infrastructure/        # Infrastructure layer
│           ├── http/
│           │   ├── handlers.go    # HTTP handlers
│           │   └── routes.go      # Route configuration
│           └── postgres/
│               └── repository.go  # PostgreSQL implementation
├── docker-compose.yml             # Docker services
├── Dockerfile                     # Application container
├── go.mod                         # Go modules
└── README.md                      # Documentation
```

## Scalability Features

This implementation includes several scalability features used by big tech companies:

1. **Clean Architecture**: Separation of concerns for easy maintenance
2. **Repository Pattern**: Database abstraction for easy testing and switching
3. **Dependency Injection**: Loose coupling between components
4. **Context Support**: Proper request cancellation and timeouts
5. **Error Handling**: Structured error responses
6. **Input Validation**: Request validation at multiple layers
7. **Database Indexing**: GORM handles indexing automatically
8. **Connection Pooling**: GORM manages database connections

## Future Enhancements

- [ ] Authentication and authorization
- [ ] Rate limiting
- [ ] Caching layer (Redis)
- [ ] Event-driven architecture
- [ ] API versioning
- [ ] Swagger documentation
- [ ] Unit and integration tests
- [ ] Monitoring and logging
- [ ] GraphQL support
- [ ] Real-time notifications

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License. 