# Payments Service

A production-ready Go API service for managing insurance payments and claim payouts with CloudBees Feature Management integration.

## HIGHEST-RISK SERVICE - RESTRICTED DEPLOYMENT

This is the HIGHEST-RISK service in the InsuranceStack platform due to its handling of financial transactions. Special governance measures are enforced:

- **Restricted Deployment Windows**: Deployments only allowed during business hours (9 AM - 5 PM EST, weekdays)
- **Extended Approval Process**: Requires approval from both Finance and Engineering leads
- **Mandatory Pre-Production Validation**: All changes must pass extended testing in staging
- **Enhanced Monitoring**: Real-time alerts for all failed transactions and anomalies
- **Audit Logging**: All payment operations are logged for compliance
- **Rate Limiting**: Strict rate limits to prevent abuse
- **Zero-Downtime Deployments**: Blue-green deployment strategy required

## Features

- RESTful API for payment and payout management
- Premium payment processing for insurance policies
- Claim payout processing
- Feature flag system ready for CloudBees Feature Management integration
- Feature flag: `payments.instantPayouts` - toggle between instant and batch payout processing
- Environment-based feature flags (with CloudBees integration guide included)
- Proper error handling and logging
- CORS support
- Graceful shutdown
- Health check endpoint
- JSON structured logging

## Project Structure

```
payments-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── health.go           # Health check handler
│   │   └── payment.go          # Payment endpoints
│   ├── services/                # Business logic
│   │   └── payment_service.go  # Payment business logic
│   ├── repository/              # Data access layer
│   │   └── repository.go       # Repository implementation
│   ├── features/                # Feature flags
│   │   └── flags.go            # CloudBees FM/Rox integration
│   ├── models/                  # Data models
│   │   └── payment.go          # Payment model
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
  "service": "payments-service"
}
```

### List Payments

**GET /payments**

Returns all payments for the authenticated user.

**Headers:**
- `X-User-ID` (optional): User ID for authentication

**Response:**
```json
[
  {
    "id": "pay-001",
    "type": "premium",
    "policyId": "pol-001",
    "customerId": "cust-001",
    "amount": 150.00,
    "status": "completed",
    "processedDate": "2024-12-20T10:30:00Z",
    "createdAt": "2024-12-20T10:00:00Z",
    "updatedAt": "2024-12-20T10:30:00Z"
  }
]
```

### Get Payment by ID

**GET /payments/{id}**

Returns a specific payment by ID.

**Parameters:**
- `id` (path): Payment ID

**Response:**
```json
{
  "id": "pay-001",
  "type": "premium",
  "policyId": "pol-001",
  "customerId": "cust-001",
  "amount": 150.00,
  "status": "completed",
  "processedDate": "2024-12-20T10:30:00Z",
  "createdAt": "2024-12-20T10:00:00Z",
  "updatedAt": "2024-12-20T10:30:00Z"
}
```

### Create Premium Payment

**POST /payments**

Creates a new premium payment for an insurance policy.

**Request Body:**
```json
{
  "policyId": "pol-001",
  "customerId": "cust-001",
  "amount": 150.00
}
```

**Response:**
```json
{
  "id": "pay-002",
  "type": "premium",
  "policyId": "pol-001",
  "customerId": "cust-001",
  "amount": 150.00,
  "status": "pending",
  "createdAt": "2024-12-21T10:00:00Z",
  "updatedAt": "2024-12-21T10:00:00Z"
}
```

### Create Claim Payout

**POST /payouts**

Creates a new payout for an insurance claim.

**Request Body:**
```json
{
  "claimId": "claim-001",
  "customerId": "cust-001",
  "amount": 5000.00
}
```

**Response:**
```json
{
  "id": "pay-003",
  "type": "payout",
  "claimId": "claim-001",
  "customerId": "cust-001",
  "amount": 5000.00,
  "status": "pending",
  "createdAt": "2024-12-21T10:00:00Z",
  "updatedAt": "2024-12-21T10:00:00Z"
}
```

### Process Payment

**PUT /payments/{id}/process**

Processes a pending payment or payout.

**Parameters:**
- `id` (path): Payment ID

**Response:**
```json
{
  "id": "pay-002",
  "type": "premium",
  "policyId": "pol-001",
  "customerId": "cust-001",
  "amount": 150.00,
  "status": "completed",
  "processedDate": "2024-12-21T10:05:00Z",
  "createdAt": "2024-12-21T10:00:00Z",
  "updatedAt": "2024-12-21T10:05:00Z"
}
```

**Error Responses:**

- `404 Not Found` - Payment does not exist
- `400 Bad Request` - Invalid payment data or payment already processed

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8005` |
| `DATA_PATH` | Path to seed data directory | `../../data/seed` |
| `CLOUDBEES_FM_API_KEY` | CloudBees Feature Management API key (optional) | `dev-mode` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `FEATURE_INSTANT_PAYOUTS` | Enable instant payouts vs batch processing (true/false) | `false` |

## Feature Flags

### payments.instantPayouts

**Default:** `false`

When enabled, claim payouts are processed instantly. When disabled, payouts are queued for batch processing at end of business day. This demonstrates controlled rollout of high-risk financial features.

**Risk Level:** HIGH - Instant payouts bypass fraud detection checks and batch reconciliation
**Recommended Strategy:** Enable for 5% of claims initially, monitor for fraud, gradually increase

**Current Implementation:** This flag is controlled via the `FEATURE_INSTANT_PAYOUTS` environment variable.

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
export PORT=8005
export DATA_PATH=../../data/seed
export CLOUDBEES_FM_API_KEY=your-api-key-here
export LOG_LEVEL=info
export FEATURE_INSTANT_PAYOUTS=false
```

3. Run the server:

```bash
go run cmd/server/main.go
```

Or build and run:

```bash
go build -o bin/payments-service cmd/server/main.go
./bin/payments-service
```

### Testing the API

#### Health Check
```bash
curl http://localhost:8005/healthz
```

#### List All Payments
```bash
curl http://localhost:8005/payments
```

#### Get Specific Payment
```bash
curl http://localhost:8005/payments/pay-001
```

#### Create Premium Payment
```bash
curl -X POST http://localhost:8005/payments \
  -H "Content-Type: application/json" \
  -d '{"policyId":"pol-001","customerId":"cust-001","amount":150.00}'
```

#### Create Claim Payout
```bash
curl -X POST http://localhost:8005/payouts \
  -H "Content-Type: application/json" \
  -d '{"claimId":"claim-001","customerId":"cust-001","amount":5000.00}'
```

#### Process Payment
```bash
curl -X PUT http://localhost:8005/payments/pay-001/process
```

## Development

### Build

```bash
go build -o bin/payments-service cmd/server/main.go
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
- **Auth**: Extracts and validates user authentication

### Feature Management

CloudBees Feature Management (Rox SDK) is integrated for runtime feature toggling without deployments. Feature flags are fetched on startup and can be updated in real-time.

## License

Copyright 2024 CB-InsuranceStack

## Quick Start

The fastest way to get started:

```bash
# From the payments-service directory
export DATA_PATH=../../data/seed
go run cmd/server/main.go
```

Then test the API:

```bash
# Health check
curl http://localhost:8005/healthz

# List payments
curl http://localhost:8005/payments

# Create a premium payment
curl -X POST http://localhost:8005/payments \
  -H "Content-Type: application/json" \
  -d '{"policyId":"pol-001","customerId":"cust-001","amount":150.00}'

# Test with instant payouts enabled
export FEATURE_INSTANT_PAYOUTS=true
go run cmd/server/main.go
# In another terminal:
curl -X POST http://localhost:8005/payouts \
  -H "Content-Type: application/json" \
  -d '{"claimId":"claim-001","customerId":"cust-001","amount":5000.00}'
# Payout will be processed instantly instead of queued for batch processing
```

## Support

For issues or questions, please contact the development team.

## Governance and Compliance

This service is subject to enhanced governance requirements:

1. **Change Management**: All changes require CAB (Change Advisory Board) approval
2. **Audit Trail**: Complete audit logs maintained for 7 years for regulatory compliance
3. **PCI DSS**: Payment card data handling must comply with PCI DSS standards
4. **SOX Compliance**: Financial controls must meet SOX requirements
5. **Fraud Detection**: All transactions monitored for fraudulent patterns
6. **Disaster Recovery**: RPO of 15 minutes, RTO of 1 hour
