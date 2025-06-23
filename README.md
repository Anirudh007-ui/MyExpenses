# MyExpenses

A cloud-ready, containerized Go CRUD API for managing personal expenses.

## System Flow Diagram

```mermaid
flowchart TD
    "Client" -->|"REST API"| "GoAPI"
    "GoAPI" -->|"CRUD"| "PostgreSQL"
    "GoAPI" -->|"Cache"| "Redis"
    "GoAPI" -->|"Queue"| "SQS/PubSub"
    "GoAPI" -->|"Storage"| "S3/GCS"
    "GoAPI" -->|"Scheduled"| "CronJobs"
```

## Features

- Add, Get, Update, Delete Expenses
- Import Expenses from CSV, Excel, Google Sheets
- Cloud-agnostic (AWS/GCP)
- Dockerized
- Redis Caching
- SQS/PubSub for async tasks
- S3/GCS for file storage
- Health checks
- Swagger API docs

## Tech Stack

- Go
- PostgreSQL
- Redis
- Docker
- SQS/PubSub
- S3/GCS

## Getting Started

1. Clone the repo
2. Follow the setup instructions (to be added) 