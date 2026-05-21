# Policy Service

A production-ready Go API service for managing insurance policies with CloudBees Feature Management integration.

## Features

- RESTful API for insurance policy management
- Feature flag system ready for CloudBees Feature Management integration
- Feature flag: `api.maskAmounts` - dynamically mask premium amounts in responses
- Environment-based feature flags (with CloudBees integration guide included)
- Proper error handling and logging
- CORS support
- Graceful shutdown
- Health check endpoint
- JSON structured logging

## Project Structure

```
policy-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── health.go           # Health check handler
│   │   └── policy.go           # Policy endpoints
│   ├── services/                # Business logic
│   │   └── policy_service.go   # Policy business logic
│   ├── repository/              # Data access layer
│   │   └── repository.go       # Repository implementation
│   ├── features/                # Feature flags
│   │   └── flags.go            # CloudBees FM/Rox integration
│   ├── models/                  # Data models
│   │   └── policy.go           # Policy model
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
  "timestamp": "2024-12-13T10:30:00Z",
  "service": "policy-service"
}
```

### List All Policies

**GET /policies**

Returns all insurance policies for the authenticated customer.

**Headers:**
- `X-User-ID` (optional): Customer ID, defaults to `customer-001` if not provided

**Response (when maskAmounts = false):**
```json
[
  {
    "id": "pol-001",
    "customerId": "customer-001",
    "policyNumber": "AUTO-2024-001234",
    "type": "auto",
    "status": "active",
    "premium": 1250.00,
    "currency": "USD",
    "startDate": "2024-01-01T00:00:00Z",
    "endDate": "2025-01-01T00:00:00Z",
    "createdAt": "2023-12-15T10:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

**Response (when maskAmounts = true):**
```json
[
  {
    "id": "pol-001",
    "customerId": "customer-001",
    "policyNumber": "AUTO-2024-001234",
    "type": "auto",
    "status": "active",
    "premium": "***.**",
    "currency": "USD",
    "startDate": "2024-01-01T00:00:00Z",
    "endDate": "2025-01-01T00:00:00Z",
    "createdAt": "2023-12-15T10:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

### Get Policy by ID

**GET /policies/{id}**

Returns a specific policy by ID. The policy must belong to the authenticated customer.

**Headers:**
- `X-User-ID` (optional): Customer ID, defaults to `customer-001` if not provided

**Parameters:**
- `id` (path): Policy ID

**Response:**
```json
{
  "id": "pol-001",
  "customerId": "customer-001",
  "policyNumber": "AUTO-2024-001234",
  "type": "auto",
  "status": "active",
  "premium": 1250.00,
  "currency": "USD",
  "startDate": "2024-01-01T00:00:00Z",
  "endDate": "2025-01-01T00:00:00Z",
  "createdAt": "2023-12-15T10:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Error Responses:**

- `404 Not Found` - Policy does not exist
- `403 Forbidden` - Policy does not belong to the customer

### Create New Policy

**POST /policies**

Creates a new insurance policy for a customer.

**Headers:**
- `X-User-ID` (optional): Customer ID, defaults to `customer-001` if not provided

**Request Body:**
```json
{
  "customerId": "customer-001",
  "policyNumber": "AUTO-2024-001235",
  "type": "auto",
  "premium": 1500.00,
  "startDate": "2025-01-01T00:00:00Z",
  "endDate": "2026-01-01T00:00:00Z"
}
```

**Response:** `201 Created` with the created policy object

### Update Policy

**PUT /policies/{id}**

Updates an existing policy.

**Headers:**
- `X-User-ID` (optional): Customer ID, defaults to `customer-001` if not provided

**Request Body:**
```json
{
  "status": "active",
  "premium": 1350.00,
  "endDate": "2026-06-01T00:00:00Z"
}
```

**Response:** `200 OK` with the updated policy object

### Cancel Policy

**DELETE /policies/{id}**

Cancels a policy by setting its status to "cancelled".

**Headers:**
- `X-User-ID` (optional): Customer ID, defaults to `customer-001` if not provided

**Response:** `200 OK` with the cancelled policy object

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8001` |
| `DATA_PATH` | Path to seed data directory | `../../data/seed` |
| `CLOUDBEES_FM_API_KEY` | CloudBees Feature Management API key (optional) | `dev-mode` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `FEATURE_MASK_AMOUNTS` | Enable premium masking (true/false) | `false` |

## Feature Flags

### api.maskAmounts

**Default:** `false`

When enabled, all premium amounts in policy responses are masked as `"***.**"` for privacy and security.

**Current Implementation:** This flag is controlled via the `FEATURE_MASK_AMOUNTS` environment variable.

**CloudBees Integration:** The codebase is ready for CloudBees Feature Management integration. See `internal/features/flags.go` for detailed integration instructions. Once integrated, flags can be toggled in real-time without redeploying the service.

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
export PORT=8001
export DATA_PATH=/Users/brown/git_orgs/CB-InsuranceStack/InsuranceStack/data/seed
export CLOUDBEES_FM_API_KEY=your-api-key-here
export LOG_LEVEL=info
```

3. Run the server:

```bash
go run cmd/server/main.go
```

Or build and run:

```bash
go build -o bin/policy-service cmd/server/main.go
./bin/policy-service
```

### Testing the API

#### Health Check
```bash
curl http://localhost:8001/healthz
```

#### Get All Policies
```bash
curl http://localhost:8001/policies
```

#### Get Specific Policy
```bash
curl http://localhost:8001/policies/pol-001
```

#### Create New Policy
```bash
curl -X POST http://localhost:8001/policies \
  -H "Content-Type: application/json" \
  -d '{"customerId":"customer-001","policyNumber":"AUTO-2025-001","type":"auto","premium":1500.00,"startDate":"2025-01-01T00:00:00Z","endDate":"2026-01-01T00:00:00Z"}'
```

#### Test with Different Customer
```bash
curl -H "X-User-ID: customer-002" http://localhost:8001/policies
```

## Development

### Build

```bash
go build -o bin/policy-service cmd/server/main.go
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

1. **Authentication**: The current implementation uses a simple `X-User-ID` header for demo purposes. In production, implement proper JWT token validation or OAuth2.

2. **CORS**: The CORS middleware currently allows all origins (`*`). In production, specify exact allowed origins.

3. **Database**: Data is loaded from JSON files. In production, integrate with a proper database (PostgreSQL, MySQL, etc.).

4. **Monitoring**: Add metrics collection (Prometheus), distributed tracing (OpenTelemetry), and error tracking (Sentry).

5. **Rate Limiting**: Implement rate limiting to prevent abuse.

6. **TLS**: Enable HTTPS with proper certificates.

7. **Secrets Management**: Use a secrets manager (AWS Secrets Manager, HashiCorp Vault) for sensitive configuration.

8. **Container Deployment**: Create a Dockerfile for containerized deployment.

## Architecture

### Layered Architecture

1. **Handlers Layer** (`internal/handlers/`): HTTP request/response handling
2. **Services Layer** (`internal/services/`): Business logic and feature flag application
3. **Repository Layer** (`internal/repository/`): Data access abstraction
4. **Models Layer** (`internal/models/`): Domain models and data structures

### Middleware

- **Logging**: Logs all HTTP requests with method, path, status, and duration
- **CORS**: Handles cross-origin resource sharing
- **Auth**: Extracts and validates customer authentication

### Feature Management

CloudBees Feature Management (Rox SDK) is integrated for runtime feature toggling without deployments. Feature flags are fetched on startup and can be updated in real-time.

## License

Copyright 2024 CB-InsuranceStack

## Quick Start

The fastest way to get started:

```bash
# From the policy-service directory
export DATA_PATH=../../data/seed
go run cmd/server/main.go
```

Then test the API:

```bash
# Health check
curl http://localhost:8001/healthz

# Get policies
curl http://localhost:8001/policies

# Get specific policy
curl http://localhost:8001/policies/pol-001

# Test with masked amounts
export FEATURE_MASK_AMOUNTS=true
go run cmd/server/main.go
# In another terminal:
curl http://localhost:8001/policies
# You should see "***.**" instead of actual premium amounts
```

## Support

For issues or questions, please contact the development team.
