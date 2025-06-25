# MyExpenses Project Structure Guide

This document explains the project structure, architecture patterns, and organization principles used in the MyExpenses application.

## 🏗️ Architecture Overview

The project follows **Clean Architecture** principles, which separates concerns into distinct layers with clear dependencies and responsibilities.

```
┌─────────────────────────────────────┐
│           HTTP Layer                │ ← Handles HTTP requests/responses
├─────────────────────────────────────┤
│        Application Layer            │ ← Business logic and use cases
├─────────────────────────────────────┤
│          Domain Layer               │ ← Core business rules and entities
├─────────────────────────────────────┤
│      Infrastructure Layer           │ ← Database, external services
└─────────────────────────────────────┘
```

### Dependency Flow
- **Outer layers** depend on **inner layers**
- **Inner layers** have **no dependencies** on outer layers
- **Domain layer** is the **most independent** (no external dependencies)

---

## 📁 Project Structure

```
MyExpenses/
├── cmd/                              # 🚀 Application entry points
│   └── api/
│       └── main.go                   # Main application startup
├── internal/                         # 🔒 Private application code
│   ├── db/                          # 🗄️ Database configuration
│   │   └── postgres.go              # PostgreSQL connection setup
│   └── expenses/                    # 💰 Expense domain
│       ├── domain/                  # 🎯 Core business logic
│       │   ├── expense.go           # Expense entity
│       │   ├── errors.go            # Domain-specific errors
│       │   └── repository.go        # Data access interface
│       ├── application/             # 🧠 Business logic orchestration
│       │   └── service.go           # Use cases and services
│       └── infrastructure/          # 🔧 External concerns
│           ├── http/                # 🌐 HTTP handling
│           │   ├── handlers.go      # Request/response handling
│           │   └── routes.go        # URL routing
│           └── postgres/            # 🗄️ Database implementation
│               └── repository.go    # PostgreSQL operations
├── docs/                            # 📚 Documentation
│   ├── diagnostic_commands.md       # Troubleshooting guide
│   ├── api_endpoints_guide.md       # API documentation
│   └── project_structure_guide.md   # This file
├── go.mod                           # 📦 Go module definition
├── go.sum                           # 🔒 Dependency checksums
├── docker-compose.yml               # 🐳 Container orchestration
├── Dockerfile                       # 🐳 Application container
└── README.md                        # 📖 Project overview
```

---

## 🎯 Layer-by-Layer Breakdown

### 1. Domain Layer (`internal/expenses/domain/`)

**Purpose:** Contains the core business logic and entities.

#### Files:
- **`expense.go`** - Core business entity
- **`errors.go`** - Domain-specific error definitions
- **`repository.go`** - Data access interface

#### Key Concepts:
```go
// Business Entity
type Expense struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    Description string    `json:"description" gorm:"not null"`
    Amount      float64   `json:"amount" gorm:"not null"`
    Category    string    `json:"category" gorm:"not null"`
    Date        time.Time `json:"date" gorm:"not null"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Repository Interface (Contract)
type Repository interface {
    Create(ctx context.Context, expense *Expense) error
    GetByID(ctx context.Context, id string) (*Expense, error)
    GetAll(ctx context.Context, filters map[string]interface{}) ([]*Expense, error)
    Update(ctx context.Context, expense *Expense) error
    Delete(ctx context.Context, id string) error
    Exists(ctx context.Context, id string) (bool, error)
}
```

**Design Patterns:**
- **Entity Pattern** - Core business objects
- **Repository Interface** - Data access abstraction
- **Factory Pattern** - Valid object creation
- **Domain Errors** - Business rule violations

---

### 2. Application Layer (`internal/expenses/application/`)

**Purpose:** Orchestrates business logic and coordinates between layers.

#### Files:
- **`service.go`** - Business logic services and use cases

#### Key Concepts:
```go
// Service with dependency injection
type Service struct {
    repo domain.Repository  // Depends on interface, not implementation
}

// Use case methods
func (s *Service) CreateExpense(ctx context.Context, req *CreateExpenseRequest) (*domain.Expense, error)
func (s *Service) GetExpense(ctx context.Context, id string) (*domain.Expense, error)
func (s *Service) GetAllExpenses(ctx context.Context, filters map[string]interface{}) ([]*domain.Expense, error)
func (s *Service) UpdateExpense(ctx context.Context, id string, req *UpdateExpenseRequest) (*domain.Expense, error)
func (s *Service) DeleteExpense(ctx context.Context, id string) error
```

**Design Patterns:**
- **Dependency Injection** - Services receive dependencies
- **Use Case Pattern** - Specific business operations
- **DTO Pattern** - Request/Response data transfer objects
- **Error Wrapping** - Context preservation

---

### 3. Infrastructure Layer (`internal/expenses/infrastructure/`)

**Purpose:** Handles external concerns like databases and HTTP.

#### HTTP Layer (`infrastructure/http/`)
- **`handlers.go`** - HTTP request/response handling
- **`routes.go`** - URL routing configuration

#### Database Layer (`infrastructure/postgres/`)
- **`repository.go`** - PostgreSQL implementation

#### Key Concepts:
```go
// HTTP Handler
type Handler struct {
    service *application.Service  // Depends on application layer
}

// PostgreSQL Repository
type Repository struct {
    db *gorm.DB  // Database connection
}
```

**Design Patterns:**
- **Repository Implementation** - Concrete data access
- **HTTP Handler Pattern** - Request processing
- **Route Grouping** - Organized endpoints
- **Error Mapping** - Database to domain errors

---

### 4. Database Configuration (`internal/db/`)

**Purpose:** Database connection and configuration management.

#### Files:
- **`postgres.go`** - PostgreSQL connection setup

#### Key Concepts:
```go
// Configuration struct
type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

// Connection function
func Connect(config *Config) (*gorm.DB, error)
```

**Design Patterns:**
- **Configuration Pattern** - Centralized settings
- **Environment Variables** - Flexible configuration
- **Connection Factory** - Database connection creation

---

### 5. Application Entry Point (`cmd/api/`)

**Purpose:** Wires together all layers and starts the application.

#### Files:
- **`main.go`** - Application startup and dependency injection

#### Key Concepts:
```go
func main() {
    // 1. Load configuration
    dbConfig := db.NewConfig()
    
    // 2. Connect to database
    database, err := db.Connect(dbConfig)
    
    // 3. Initialize layers
    repo := postgres.NewRepository(database)
    service := application.NewService(repo)
    
    // 4. Setup HTTP server
    router := gin.Default()
    http.SetupRoutes(router, service)
    
    // 5. Start server
    router.Run(":8080")
}
```

**Design Patterns:**
- **Dependency Injection** - Wiring components
- **Startup Sequence** - Ordered initialization
- **Error Handling** - Graceful startup failures

---

## 🔄 Data Flow Examples

### Creating an Expense
```
1. HTTP Request → POST /expenses
   ↓
2. Router → CreateExpense handler
   ↓
3. Handler → Parse JSON, validate request
   ↓
4. Service → Call domain.NewExpense() (business rules)
   ↓
5. Domain → Validate and create expense entity
   ↓
6. Service → Call repo.Create() (persistence)
   ↓
7. Repository → Execute SQL INSERT
   ↓
8. Database → Store the data
   ↓
9. Response → Return 201 Created with expense data
```

### Getting Expenses
```
1. HTTP Request → GET /expenses?category=Food
   ↓
2. Router → GetAllExpenses handler
   ↓
3. Handler → Parse query parameters
   ↓
4. Service → Call repo.GetAll() with filters
   ↓
5. Repository → Execute SQL SELECT with WHERE clauses
   ↓
6. Database → Return filtered results
   ↓
7. Response → Return 200 OK with expenses array
```

---

## 🎨 Design Patterns Used

### 1. Clean Architecture
- **Separation of Concerns** - Each layer has specific responsibilities
- **Dependency Inversion** - High-level modules don't depend on low-level modules
- **Interface Segregation** - Small, focused interfaces

### 2. Repository Pattern
- **Abstraction** - Data access is abstracted through interfaces
- **Testability** - Easy to mock for testing
- **Flexibility** - Can swap database implementations

### 3. Dependency Injection
- **Loose Coupling** - Components don't create their dependencies
- **Testability** - Dependencies can be injected for testing
- **Flexibility** - Easy to change implementations

### 4. Factory Pattern
- **Valid Object Creation** - Ensures business rules are enforced
- **Centralized Logic** - Creation logic is in one place
- **Validation** - Objects are validated at creation time

### 5. Error Wrapping
- **Context Preservation** - Original errors are preserved
- **Debugging** - Stack traces show the full error chain
- **User-Friendly** - Errors can be translated for users

---

## 🧪 Testing Strategy

### Unit Testing
- **Domain Layer** - Test business rules and validation
- **Application Layer** - Test use cases with mocked repositories
- **Infrastructure Layer** - Test HTTP handlers and database operations

### Integration Testing
- **Database Integration** - Test with real PostgreSQL
- **API Integration** - Test full HTTP request/response cycle

### Test Structure
```
tests/
├── unit/
│   ├── domain/
│   ├── application/
│   └── infrastructure/
├── integration/
│   ├── api/
│   └── database/
└── fixtures/
    └── test_data.json
```

---

## 🚀 Deployment Structure

### Docker
- **Multi-stage builds** - Separate build and runtime environments
- **Environment variables** - Configuration through environment
- **Health checks** - Container health monitoring

### Docker Compose
- **Service orchestration** - Multiple services (app + database)
- **Network isolation** - Services communicate through networks
- **Volume persistence** - Database data persistence

---

## 📈 Scalability Considerations

### Current Architecture Benefits
- **Horizontal Scaling** - Stateless application layer
- **Database Scaling** - Repository pattern allows database changes
- **Caching** - Can add caching layer without changing business logic
- **Microservices** - Each domain can become a separate service

### Future Enhancements
- **Event-Driven Architecture** - Add event publishing/subscribing
- **CQRS** - Separate read and write models
- **GraphQL** - Add GraphQL layer for flexible queries
- **Real-time** - WebSocket support for real-time updates

---

## 🔧 Development Workflow

### Adding New Features
1. **Domain Layer** - Define entities and business rules
2. **Application Layer** - Implement use cases
3. **Infrastructure Layer** - Add HTTP handlers and database operations
4. **Testing** - Write unit and integration tests
5. **Documentation** - Update API documentation

### Code Organization Principles
- **Single Responsibility** - Each file has one clear purpose
- **Open/Closed Principle** - Open for extension, closed for modification
- **Interface Segregation** - Small, focused interfaces
- **Dependency Inversion** - Depend on abstractions, not concretions

---

## 📚 Learning Resources

### Clean Architecture
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)

### Go Best Practices
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Design Patterns
- [Design Patterns: Elements of Reusable Object-Oriented Software](https://en.wikipedia.org/wiki/Design_Patterns)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)

This structure provides a solid foundation for building scalable, maintainable applications that can evolve with business needs. 