# MyExpenses API Endpoints Guide

This document provides a complete guide to all API endpoints in the MyExpenses application, including request/response examples, testing commands, and error handling.

## 🚀 Base URL
```
http://localhost:8080
```

## 📋 API Endpoints Overview

| Method | Endpoint | Description | Status Code |
|--------|----------|-------------|-------------|
| GET | `/health` | Health check | 200 |
| POST | `/expenses` | Create expense | 201 |
| GET | `/expenses` | Get all expenses | 200 |
| GET | `/expenses/{id}` | Get expense by ID | 200 |
| PUT | `/expenses/{id}` | Update expense | 200 |
| DELETE | `/expenses/{id}` | Delete expense | 200 |

---

## 🔍 Health Check

### GET /health
**Purpose:** Verify the API server is running and healthy.

**Request:**
```bash
curl -X GET http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "service": "MyExpenses API"
}
```

**Status Code:** 200 OK

---

## 💰 Expense Management

### POST /expenses
**Purpose:** Create a new expense.

**Request Body:**
```json
{
  "description": "Coffee",
  "amount": 4.50,
  "category": "Food",
  "date": "2024-01-15T08:00:00Z"
}
```

**Testing Command:**
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

**Expected Response:**
```json
{
  "message": "Expense created successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "description": "Coffee",
    "amount": 4.50,
    "category": "Food",
    "date": "2024-01-15T08:00:00Z",
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-01-15T08:00:00Z"
  }
}
```

**Status Code:** 201 Created

**Validation Rules:**
- `description`: Required, non-empty string
- `amount`: Required, greater than 0
- `category`: Required, non-empty string
- `date`: Required, valid ISO 8601 date

---

### GET /expenses
**Purpose:** Retrieve all expenses with optional filtering.

**Query Parameters:**
- `category` - Filter by category (partial match, case-insensitive)
- `date_from` - Filter expenses from this date (YYYY-MM-DD)
- `date_to` - Filter expenses until this date (YYYY-MM-DD)
- `min_amount` - Minimum amount filter
- `max_amount` - Maximum amount filter
- `description` - Filter by description (partial match, case-insensitive)

**Testing Commands:**

**Get all expenses:**
```bash
curl -X GET http://localhost:8080/expenses
```

**Get expenses with filters:**
```bash
# Filter by category
curl -X GET "http://localhost:8080/expenses?category=Food"

# Filter by amount range
curl -X GET "http://localhost:8080/expenses?min_amount=10&max_amount=50"

# Filter by date range
curl -X GET "http://localhost:8080/expenses?date_from=2024-01-01&date_to=2024-01-31"

# Filter by description
curl -X GET "http://localhost:8080/expenses?description=coffee"

# Multiple filters
curl -X GET "http://localhost:8080/expenses?category=Food&min_amount=5&date_from=2024-01-01"
```

**Expected Response:**
```json
{
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "description": "Coffee",
      "amount": 4.50,
      "category": "Food",
      "date": "2024-01-15T08:00:00Z",
      "created_at": "2024-01-15T08:00:00Z",
      "updated_at": "2024-01-15T08:00:00Z"
    }
  ],
  "count": 1
}
```

**Status Code:** 200 OK

---

### GET /expenses/{id}
**Purpose:** Retrieve a specific expense by ID.

**Path Parameter:**
- `id` - UUID of the expense

**Testing Command:**
```bash
curl -X GET http://localhost:8080/expenses/123e4567-e89b-12d3-a456-426614174000
```

**Expected Response:**
```json
{
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "description": "Coffee",
    "amount": 4.50,
    "category": "Food",
    "date": "2024-01-15T08:00:00Z",
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-01-15T08:00:00Z"
  }
}
```

**Status Code:** 200 OK

**Error Response (404 Not Found):**
```json
{
  "error": "Expense not found"
}
```

---

### PUT /expenses/{id}
**Purpose:** Update an existing expense (partial updates supported).

**Path Parameter:**
- `id` - UUID of the expense

**Request Body (all fields optional):**
```json
{
  "description": "Updated Coffee",
  "amount": 5.00,
  "category": "Food & Beverages",
  "date": "2024-01-15T09:00:00Z"
}
```

**Testing Commands:**

**Update all fields:**
```bash
curl -X PUT http://localhost:8080/expenses/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated Coffee",
    "amount": 5.00,
    "category": "Food & Beverages",
    "date": "2024-01-15T09:00:00Z"
  }'
```

**Partial update (only amount):**
```bash
curl -X PUT http://localhost:8080/expenses/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5.00
  }'
```

**Expected Response:**
```json
{
  "message": "Expense updated successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "description": "Updated Coffee",
    "amount": 5.00,
    "category": "Food & Beverages",
    "date": "2024-01-15T09:00:00Z",
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**Status Code:** 200 OK

**Error Response (404 Not Found):**
```json
{
  "error": "Expense not found"
}
```

---

### DELETE /expenses/{id}
**Purpose:** Delete an expense.

**Path Parameter:**
- `id` - UUID of the expense

**Testing Command:**
```bash
curl -X DELETE http://localhost:8080/expenses/123e4567-e89b-12d3-a456-426614174000
```

**Expected Response:**
```json
{
  "message": "Expense deleted successfully"
}
```

**Status Code:** 200 OK

**Error Response (404 Not Found):**
```json
{
  "error": "Expense not found"
}
```

---

## 🚨 Error Handling

### Common Error Responses

**400 Bad Request (Validation Error):**
```json
{
  "error": "Invalid request body",
  "details": "Key: 'CreateExpenseRequest.Amount' Error:Field validation for 'Amount' failed on the 'gt' tag"
}
```

**404 Not Found:**
```json
{
  "error": "Expense not found"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Failed to create expense"
}
```

---

## 🧪 Complete Testing Workflow

Here's a complete testing sequence to verify all endpoints:

```bash
# 1. Health check
curl -X GET http://localhost:8080/health

# 2. Create an expense
curl -X POST http://localhost:8080/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Lunch",
    "amount": 15.50,
    "category": "Food",
    "date": "2024-01-15T12:00:00Z"
  }'

# 3. Get all expenses
curl -X GET http://localhost:8080/expenses

# 4. Get specific expense (replace {id} with actual ID from step 2)
curl -X GET http://localhost:8080/expenses/{id}

# 5. Update expense (replace {id} with actual ID)
curl -X PUT http://localhost:8080/expenses/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 16.00
  }'

# 6. Delete expense (replace {id} with actual ID)
curl -X DELETE http://localhost:8080/expenses/{id}

# 7. Verify deletion
curl -X GET http://localhost:8080/expenses/{id}
```

---

## 📊 Response Format Standards

### Success Responses
- **Single Resource:** Wrapped in `data` field
- **Multiple Resources:** Wrapped in `data` array with `count` field
- **Actions:** Include `message` field with success description

### Error Responses
- **User-friendly:** `error` field with readable message
- **Debugging:** `details` field with technical information (when applicable)
- **Consistent:** Same structure across all endpoints

---

## 🔧 Advanced Usage

### Filtering Examples
```bash
# Get all food expenses over $10
curl "http://localhost:8080/expenses?category=Food&min_amount=10"

# Get expenses from last week
curl "http://localhost:8080/expenses?date_from=2024-01-08&date_to=2024-01-14"

# Search for coffee-related expenses
curl "http://localhost:8080/expenses?description=coffee"
```

### Using with jq (JSON processor)
```bash
# Get only expense IDs
curl -s http://localhost:8080/expenses | jq '.data[].id'

# Get total amount of all expenses
curl -s http://localhost:8080/expenses | jq '[.data[].amount] | add'

# Filter by category and get count
curl -s "http://localhost:8080/expenses?category=Food" | jq '.count'
```

---

## 📝 Notes

- All dates should be in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)
- UUIDs are automatically generated for new expenses
- The API supports partial updates (only send fields you want to change)
- All responses are in JSON format
- The API follows REST conventions
- Filtering is case-insensitive for text fields 