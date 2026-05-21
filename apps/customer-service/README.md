# Customer Service

A production-ready Go API service for managing insurance customer profiles with CloudBees Feature Management integration.

## Features

- RESTful API for customer profile management
- Feature flag system ready for CloudBees Feature Management integration
- Customer risk score tracking
- Proper error handling and logging
- CORS support
- Graceful shutdown
- Health check endpoint
- JSON structured logging

## Project Structure

```
customer-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── health.go           # Health check handler
│   │   └── customer.go         # Customer endpoints
│   ├── services/                # Business logic
│   │   └── customer_service.go # Customer business logic
│   ├── repository/              # Data access layer
│   │   └── repository.go       # Repository implementation
│   ├── features/                # Feature flags
│   │   └── flags.go            # CloudBees FM/Rox integration
│   ├── models/                  # Data models
│   │   └── customer.go         # Customer model
│   └── middleware/              # HTTP middleware
│       ├── logging.go          # Request logging
│       ├── cors.go             # CORS configuration
│       └── auth.go             # Authentication
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## API Endpoints

### Health Check

**GET /healthz**

Returns the health status of the service.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-12-21T10:30:00Z",
  "service": "customer-service"
}
```

### List All Customers

**GET /customers**

Returns all customer profiles in the system.

**Response:**
```json
[
  {
    "id": "cust-001",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "phone": "+1-555-0123",
    "address": "123 Main St, Anytown, ST 12345",
    "dateOfBirth": "1985-06-15T00:00:00Z",
    "riskScore": 75,
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-12-12T15:30:00Z"
  }
]
```

### Get Customer by ID

**GET /customers/{id}**

Returns a specific customer by ID.

**Parameters:**
- `id` (path): Customer ID

**Response:**
```json
{
  "id": "cust-001",
  "firstName": "John",
  "lastName": "Doe",
  "email": "john.doe@example.com",
  "phone": "+1-555-0123",
  "address": "123 Main St, Anytown, ST 12345",
  "dateOfBirth": "1985-06-15T00:00:00Z",
  "riskScore": 75,
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-12-12T15:30:00Z"
}
```

**Error Responses:**
- `404 Not Found` - Customer does not exist

### Create New Customer

**POST /customers**

Creates a new customer profile.

**Request Body:**
```json
{
  "firstName": "Jane",
  "lastName": "Smith",
  "email": "jane.smith@example.com",
  "phone": "+1-555-0456",
  "address": "456 Oak Ave, Somewhere, ST 67890",
  "dateOfBirth": "1990-03-22T00:00:00Z"
}
```

**Response:**
```json
{
  "id": "cust-002",
  "firstName": "Jane",
  "lastName": "Smith",
  "email": "jane.smith@example.com",
  "phone": "+1-555-0456",
  "address": "456 Oak Ave, Somewhere, ST 67890",
  "dateOfBirth": "1990-03-22T00:00:00Z",
  "riskScore": 50,
  "createdAt": "2024-12-21T10:30:00Z",
  "updatedAt": "2024-12-21T10:30:00Z"
}
```

### Update Customer

**PUT /customers/{id}**

Updates an existing customer profile.

**Parameters:**
- `id` (path): Customer ID

**Request Body:**
```json
{
  "firstName": "Jane",
  "lastName": "Smith-Jones",
  "email": "jane.smith@example.com",
  "phone": "+1-555-9999",
  "address": "789 New St, Elsewhere, ST 11111",
  "dateOfBirth": "1990-03-22T00:00:00Z"
}
```

**Response:**
```json
{
  "id": "cust-002",
  "firstName": "Jane",
  "lastName": "Smith-Jones",
  "email": "jane.smith@example.com",
  "phone": "+1-555-9999",
  "address": "789 New St, Elsewhere, ST 11111",
  "dateOfBirth": "1990-03-22T00:00:00Z",
  "riskScore": 50,
  "createdAt": "2024-12-21T10:30:00Z",
  "updatedAt": "2024-12-21T11:45:00Z"
}
```

**Error Responses:**
- `404 Not Found` - Customer does not exist

### Deactivate Customer

**DELETE /customers/{id}**

Deactivates a customer (soft delete).

**Parameters:**
- `id` (path): Customer ID

**Response:**
```json
{
  "message": "Customer deactivated successfully"
}
```

**Error Responses:**
- `404 Not Found` - Customer does not exist

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8004` |
| `DATA_PATH` | Path to seed data directory | `../../data/seed` |
| `CLOUDBEES_FM_API_KEY` | CloudBees Feature Management API key (optional) | `dev-mode` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |

## Feature Flags

The codebase is ready for CloudBees Feature Management integration. See `internal/features/flags.go` for detailed integration instructions. Once integrated, flags can be toggled in real-time without redeploying the service.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Access to seed data files (included in repository)
- CloudBees Feature Management account (optional, for production feature flag management)

### Installation

1. Install dependencies:

```bash
go mod download
```

2. Set up environment variables:

```bash
export PORT=8004
export DATA_PATH=../../data/seed
export CLOUDBEES_FM_API_KEY=your-api-key-here
export LOG_LEVEL=info
```

3. Run the server:

```bash
go run cmd/server/main.go
```

Or build and run:

```bash
go build -o bin/customer-service cmd/server/main.go
./bin/customer-service
```

### Testing the API

#### Health Check
```bash
curl http://localhost:8004/healthz
```

#### Get All Customers
```bash
curl http://localhost:8004/customers
```

#### Get Specific Customer
```bash
curl http://localhost:8004/customers/cust-001
```

#### Create New Customer
```bash
curl -X POST http://localhost:8004/customers \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Jane",
    "lastName": "Smith",
    "email": "jane.smith@example.com",
    "phone": "+1-555-0456",
    "address": "456 Oak Ave, Somewhere, ST 67890",
    "dateOfBirth": "1990-03-22T00:00:00Z"
  }'
```

#### Update Customer
```bash
curl -X PUT http://localhost:8004/customers/cust-001 \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe-Smith",
    "email": "john.doe@example.com",
    "phone": "+1-555-9999",
    "address": "789 New St, Elsewhere, ST 11111",
    "dateOfBirth": "1985-06-15T00:00:00Z"
  }'
```

#### Deactivate Customer
```bash
curl -X DELETE http://localhost:8004/customers/cust-001
```

## Development

### Build

```bash
go build -o bin/customer-service cmd/server/main.go
```

### Run Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

## Production Considerations

1. **Authentication**: The current implementation uses middleware for authentication. In production, implement proper JWT token validation or OAuth2.

2. **CORS**: The CORS middleware currently allows all origins (`*`). In production, specify exact allowed origins.

3. **Database**: Data is loaded from JSON files. In production, integrate with a proper database (PostgreSQL, MySQL, etc.).

4. **Monitoring**: Add metrics collection (Prometheus), distributed tracing (OpenTelemetry), and error tracking (Sentry).

5. **Rate Limiting**: Implement rate limiting to prevent abuse.

6. **TLS**: Enable HTTPS with proper certificates.

7. **Secrets Management**: Use a secrets manager (AWS Secrets Manager, HashiCorp Vault) for sensitive configuration.

8. **Container Deployment**: Dockerfile included for containerized deployment.

## Architecture

### Layered Architecture

1. **Handlers Layer** (`internal/handlers/`): HTTP request/response handling
2. **Services Layer** (`internal/services/`): Business logic and feature flag application
3. **Repository Layer** (`internal/repository/`): Data access abstraction
4. **Models Layer** (`internal/models/`): Domain models and data structures

### Middleware

- **Logging**: Logs all HTTP requests with method, path, status, and duration
- **CORS**: Handles cross-origin resource sharing
- **Auth**: Extracts and validates authentication

### Feature Management

CloudBees Feature Management (Rox SDK) integration is ready for runtime feature toggling without deployments. Feature flags can be fetched on startup and updated in real-time.

## License

Copyright 2024 CB-InsuranceStack

## Quick Start

The fastest way to get started:

```bash
# From the customer-service directory
export DATA_PATH=../../data/seed
go run cmd/server/main.go
```

Then test the API:

```bash
# Health check
curl http://localhost:8004/healthz

# Get all customers
curl http://localhost:8004/customers

# Get specific customer
curl http://localhost:8004/customers/cust-001
```

## Support

For issues or questions, please contact the development team.
