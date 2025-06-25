# MyExpenses API - Complete Code Explanation for Beginners

## 🎯 Overview

This document explains every concept, pattern, and line of code in the MyExpenses API. It's designed for beginners who want to understand how a real-world Go application is structured.

## 🏗️ Architecture Overview

Our application follows **Clean Architecture**, which separates concerns into distinct layers:

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

## 📁 File Structure Explained

```
MyExpenses/
├── cmd/api/main.go                    # 🚀 Application entry point
├── internal/                          # 🔒 Private application code
│   ├── db/postgres.go                 # 🗄️ Database configuration
│   └── expenses/                      # 💰 Expense-related code
│       ├── domain/                    # 🎯 Core business logic
│       │   ├── expense.go             # 📋 Expense entity definition
│       │   ├── errors.go              # ❌ Domain-specific errors
│       │   └── repository.go          # 📝 Data access interface
│       ├── application/               # 🧠 Business logic orchestration
│       │   └── service.go             # ⚙️ Use cases and services
│       └── infrastructure/            # 🔧 External concerns
│           ├── http/                  # 🌐 HTTP handling
│           │   ├── handlers.go        # 📥 Request/response handling
│           │   └── routes.go          # 🛣️ URL routing
│           └── postgres/              # 🗄️ Database implementation
│               └── repository.go      # 💾 PostgreSQL operations
├── go.mod                             # 📦 Go module definition
├── docker-compose.yml                 # 🐳 Container orchestration
└── README.md                          # 📖 Documentation
```

## 🔑 Key Concepts Explained

### 1. **Packages and Imports**

```go
package domain  // Declares this file belongs to the 'domain' package

import (
    "time"                    // Standard library for time handling
    "github.com/google/uuid"  // External library for UUID generation
)
```

**What it means**: 
- `package` groups related code together
- `import` brings in code from other packages
- Standard library packages (like `time`) don't need a full URL
- External packages (like `github.com/google/uuid`) need the full path

### 2. **Structs and Tags**

```go
type Expense struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    Description string    `json:"description" gorm:"not null"`
    Amount      float64   `json:"amount" gorm:"not null"`
}
```

**What it means**:
- `struct` defines a custom data type with multiple fields
- Backticks contain "tags" that provide metadata
- `json:"id"` tells Go how to serialize this field to JSON
- `gorm:"not null"` tells the database this field is required

### 3. **Pointers and References**

```go
func NewExpense(...) (*Expense, error)  // Returns a pointer to Expense
var expense *Expense                    // Variable that holds a pointer
return &Expense{...}                    // Returns address of new struct
```

**What it means**:
- `*Expense` means "pointer to Expense" (like a reference in other languages)
- `&Expense{...}` creates a struct and returns its memory address
- Pointers are more efficient for large structs (avoid copying)
- `nil` means "no value" or "empty pointer"

### 4. **Interfaces**

```go
type Repository interface {
    Create(ctx context.Context, expense *Expense) error
    GetByID(ctx context.Context, id string) (*Expense, error)
    // ... more methods
}
```

**What it means**:
- `interface` defines a contract (what methods must exist)
- It doesn't specify HOW to implement, just WHAT to implement
- Any struct that has these methods "implements" this interface
- This enables dependency injection and testing

### 5. **Error Handling**

```go
if err != nil {
    return nil, fmt.Errorf("failed to create expense: %w", err)
}
```

**What it means**:
- Go uses explicit error handling (no exceptions)
- `err != nil` checks if an error occurred
- `fmt.Errorf("...: %w", err)` wraps errors with context
- `%w` preserves the original error while adding our message

### 6. **Context**

```go
func Create(ctx context.Context, expense *Expense) error
```

**What it means**:
- `context.Context` carries request-scoped values
- It enables cancellation, timeouts, and request tracing
- `c.Request.Context()` gets the HTTP request context
- This is important for handling long-running operations

## 🎯 Domain Layer Deep Dive

### **expense.go** - The Core Entity

```go
// This is our main business entity - it represents what an expense IS
type Expense struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Description string    `json:"description" gorm:"not null"`
    Amount      float64   `json:"amount" gorm:"not null"`
    Category    string    `json:"category" gorm:"not null"`
    Date        time.Time `json:"date" gorm:"not null"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
```

**Key Points**:
- This is the "source of truth" for what an expense looks like
- GORM tags tell the database how to store this data
- JSON tags tell Go how to convert this to/from JSON
- `autoCreateTime` and `autoUpdateTime` are automatically managed

### **errors.go** - Domain-Specific Errors

```go
var (
    ErrInvalidDescription = errors.New("invalid description: cannot be empty")
    ErrInvalidAmount      = errors.New("invalid amount: must be greater than 0")
    // ...
)
```

**Key Points**:
- These are business rule violations
- They're specific to our domain (expenses)
- They provide clear, actionable error messages
- They're used throughout the application for consistency

### **repository.go** - Data Access Contract

```go
type Repository interface {
    Create(ctx context.Context, expense *Expense) error
    GetByID(ctx context.Context, id string) (*Expense, error)
    GetAll(ctx context.Context, filters map[string]interface{}) ([]*Expense, error)
    Update(ctx context.Context, expense *Expense) error
    Delete(ctx context.Context, id string) error
    Exists(ctx context.Context, id string) (bool, error)
}
```

**Key Points**:
- This is an interface (contract) for data access
- It doesn't specify HOW to store data, just WHAT operations are available
- We can implement this with PostgreSQL, MySQL, in-memory storage, etc.
- This enables easy testing and swapping implementations

## 🧠 Application Layer Deep Dive

### **service.go** - Business Logic Orchestration

```go
type Service struct {
    repo domain.Repository  // Dependency injection
}

func (s *Service) CreateExpense(ctx context.Context, req *CreateExpenseRequest) (*domain.Expense, error) {
    // Step 1: Create domain object (enforces business rules)
    expense, err := domain.NewExpense(req.Description, req.Amount, req.Category, req.Date)
    if err != nil {
        return nil, fmt.Errorf("failed to create expense: %w", err)
    }

    // Step 2: Save to repository
    if err := s.repo.Create(ctx, expense); err != nil {
        return nil, fmt.Errorf("failed to save expense: %w", err)
    }

    return expense, nil
}
```

**Key Points**:
- This layer orchestrates business operations
- It coordinates between domain logic and data persistence
- It handles error wrapping and context propagation
- It implements "use cases" (specific business operations)

### **Request/Response DTOs**

```go
type CreateExpenseRequest struct {
    Description string    `json:"description" binding:"required"`
    Amount      float64   `json:"amount" binding:"required,gt=0"`
    Category    string    `json:"category" binding:"required"`
    Date        time.Time `json:"date" binding:"required"`
}
```

**Key Points**:
- DTO = Data Transfer Object
- These are separate from domain models for API flexibility
- `binding:"required"` provides validation
- They can evolve independently of domain models

## 🔧 Infrastructure Layer Deep Dive

### **postgres/repository.go** - Database Implementation

```go
type Repository struct {
    db *gorm.DB  // GORM database connection
}

func (r *Repository) GetByID(ctx context.Context, id string) (*domain.Expense, error) {
    // Parse UUID
    uuid, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid UUID format: %w", err)
    }

    // Query database
    var expense domain.Expense
    if err := r.db.WithContext(ctx).Where("id = ?", uuid).First(&expense).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrExpenseNotFound
        }
        return nil, fmt.Errorf("failed to get expense: %w", err)
    }

    return &expense, nil
}
```

**Key Points**:
- This implements the `domain.Repository` interface
- It translates domain operations into SQL queries
- It handles database-specific error mapping
- It uses GORM for convenient database operations

### **http/handlers.go** - HTTP Request Handling

```go
func (h *Handler) CreateExpense(c *gin.Context) {
    // Parse request
    var req application.CreateExpenseRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid request body",
            "details": err.Error(),
        })
        return
    }

    // Call business logic
    expense, err := h.service.CreateExpense(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create expense",
        })
        return
    }

    // Return response
    c.JSON(http.StatusCreated, gin.H{
        "message": "Expense created successfully",
        "data":    expense,
    })
}
```

**Key Points**:
- This handles HTTP-specific concerns
- It parses JSON requests and validates them
- It calls business logic and handles errors
- It formats responses with appropriate HTTP status codes

### **http/routes.go** - URL Routing

```go
func SetupRoutes(router *gin.Engine, service *application.Service) {
    handler := NewHandler(service)

    expenses := router.Group("/expenses")
    {
        expenses.POST("", handler.CreateExpense)           // POST /expenses
        expenses.GET("", handler.GetAllExpenses)          // GET /expenses
        expenses.GET("/:id", handler.GetExpense)          // GET /expenses/{id}
        expenses.PUT("/:id", handler.UpdateExpense)       // PUT /expenses/{id}
        expenses.DELETE("/:id", handler.DeleteExpense)    // DELETE /expenses/{id}
    }
}
```

**Key Points**:
- This maps URLs to handler functions
- It follows REST conventions
- Route groups organize related endpoints
- URL parameters (like `:id`) get passed to handlers

## 🚀 Application Entry Point

### **main.go** - Wiring Everything Together

```go
func main() {
    // 1. Load environment variables
    godotenv.Load()

    // 2. Connect to database
    dbConfig := db.NewConfig()
    database, err := db.Connect(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

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

**Key Points**:
- This is where dependency injection happens
- It wires together all the layers
- It handles startup errors gracefully
- It follows a clear initialization sequence

## 🔄 Data Flow Example

Let's trace a request through the entire system:

```
1. HTTP Request: POST /expenses
   ↓
2. Router: Routes to CreateExpense handler
   ↓
3. Handler: Parses JSON, validates request
   ↓
4. Service: Calls domain.NewExpense() (business rules)
   ↓
5. Domain: Validates and creates expense entity
   ↓
6. Service: Calls repo.Create() (persistence)
   ↓
7. Repository: Executes SQL INSERT
   ↓
8. Database: Stores the data
   ↓
9. Response: Returns 201 Created with expense data
```

## 🎯 Key Design Patterns

### 1. **Dependency Injection**
```go
// Instead of creating dependencies inside:
func NewService() *Service {
    repo := postgres.NewRepository()  // ❌ Tight coupling
    return &Service{repo: repo}
}

// We inject them from outside:
func NewService(repo domain.Repository) *Service {
    return &Service{repo: repo}  // ✅ Loose coupling
}
```

### 2. **Interface Segregation**
```go
// Small, focused interfaces:
type Repository interface {
    Create(ctx context.Context, expense *Expense) error
    GetByID(ctx context.Context, id string) (*Expense, error)
    // ... only what we need
}
```

### 3. **Error Wrapping**
```go
// Preserve original error while adding context:
if err := repo.Create(ctx, expense); err != nil {
    return fmt.Errorf("failed to save expense: %w", err)
}
```

### 4. **Factory Pattern**
```go
// Ensure valid object creation:
func NewExpense(description string, amount float64, ...) (*Expense, error) {
    // Validation logic here
    return &Expense{...}, nil
}
```

## 🧪 Testing Concepts

The architecture makes testing easy:

```go
// Mock repository for testing:
type MockRepository struct {
    expenses map[string]*domain.Expense
}

func (m *MockRepository) Create(ctx context.Context, expense *domain.Expense) error {
    m.expenses[expense.ID.String()] = expense
    return nil
}

// Test the service with mock:
func TestCreateExpense(t *testing.T) {
    mockRepo := &MockRepository{expenses: make(map[string]*domain.Expense)}
    service := application.NewService(mockRepo)
    
    // Test business logic without database
    expense, err := service.CreateExpense(ctx, req)
    // Assertions...
}
```

## 🚀 Next Steps

1. **Run the application**: `go run cmd/api/main.go`
2. **Test the API**: Use curl or Postman to test endpoints
3. **Add features**: Implement authentication, caching, etc.
4. **Write tests**: Add unit and integration tests
5. **Deploy**: Containerize and deploy to cloud

## 📚 Additional Resources

- [Go Language Tour](https://tour.golang.org/)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Gin Web Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)

This architecture follows industry best practices used by companies like Google, Netflix, and Uber for building scalable, maintainable microservices. 