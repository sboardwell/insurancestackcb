# Claims Service

A production-ready Go API service for managing insurance claims with CloudBees Feature Management integration. **This is the governance showcase service** demonstrating approval workflows and higher approval requirements for insurance claims processing.

## Features

- RESTful API for insurance claims management
- Claims submission, review, and approval workflow
- CloudBees Feature Management integration for dynamic feature control
- Automatic approval for low-value claims (feature flag controlled)
- CORS support for cross-origin requests
- Request logging and authentication middleware
- Docker support for containerized deployment
- Graceful shutdown handling
- Health check endpoint
- Governance and compliance workflow showcase

## API Endpoints

### Health Check
```
GET /healthz
```
Returns the health status of the service.

**Response:**
```json
{
  "status": "healthy",
  "service": "claims-service"
}
```

### List Claims
```
GET /claims
```
Retrieves a list of insurance claims with optional filtering.

**Query Parameters:**
- `policyId` (string) - Filter by policy ID
- `customerId` (string) - Filter by customer ID
- `status` (string) - Filter by status (submitted/under_review/approved/rejected)
- `type` (string) - Filter by type (accident/theft/damage)

**Example Requests:**
```bash
# List all claims
curl "http://localhost:8002/claims"

# Filter by policy
curl "http://localhost:8002/claims?policyId=pol-001"

# Filter by status
curl "http://localhost:8002/claims?status=under_review"
```

**Response:**
```json
[
  {
    "id": "claim-001",
    "policyId": "pol-001",
    "customerId": "cust-001",
    "claimNumber": "CLM-2024-001",
    "type": "accident",
    "status": "under_review",
    "amount": 5000.00,
    "description": "Vehicle collision on highway",
    "submittedDate": "2024-12-13T10:00:00Z",
    "reviewedDate": null,
    "createdAt": "2024-12-13T10:00:00Z",
    "updatedAt": "2024-12-13T10:00:00Z"
  }
]
```

### Get Claim by ID
```
GET /claims/{id}
```
Retrieves a specific claim by ID.

**Example:**
```bash
curl "http://localhost:8002/claims/claim-001"
```

**Response:**
```json
{
  "id": "claim-001",
  "policyId": "pol-001",
  "customerId": "cust-001",
  "claimNumber": "CLM-2024-001",
  "type": "accident",
  "status": "under_review",
  "amount": 5000.00,
  "description": "Vehicle collision on highway",
  "submittedDate": "2024-12-13T10:00:00Z",
  "reviewedDate": null,
  "createdAt": "2024-12-13T10:00:00Z",
  "updatedAt": "2024-12-13T10:00:00Z"
}
```

### Submit New Claim
```
POST /claims
```
Submits a new insurance claim.

**Request Body:**
```json
{
  "policyId": "pol-001",
  "customerId": "cust-001",
  "type": "accident",
  "amount": 5000.00,
  "description": "Vehicle collision on highway"
}
```

**Response:**
```json
{
  "id": "claim-001",
  "claimNumber": "CLM-2024-001",
  "status": "submitted",
  "message": "Claim submitted successfully"
}
```

### Update Claim
```
PUT /claims/{id}
```
Updates an existing claim.

**Request Body:**
```json
{
  "amount": 5500.00,
  "description": "Updated description with additional details"
}
```

### Change Claim Status
```
PUT /claims/{id}/status
```
Changes the status of a claim (for approval workflows).

**Request Body:**
```json
{
  "status": "approved",
  "notes": "Claim approved after review"
}
```

## Feature Flags

### `claims.autoApproval` (default: false)
Controls whether low-value claims are automatically approved without manual review.

**When disabled (false):**
- All claims require manual review and approval
- Claims move to "under_review" status upon submission
- Claims must be explicitly approved via the status change endpoint

**When enabled (true):**
- Claims under a certain threshold (e.g., $1000) are automatically approved
- Higher-value claims still require manual review
- Demonstrates governance workflows with conditional automation

**Configuration:**
Set up this feature flag in CloudBees Feature Management dashboard with the key `claims.autoApproval`.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8002` |
| `CLOUDBEES_FM_API_KEY` | CloudBees Feature Management API key | (required) |
| `DATA_PATH` | Path to seed data directory | `/data/seed` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `FEATURE_AUTO_APPROVAL` | Enable auto-approval for low-value claims | `false` |

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Docker (optional, for containerized deployment)
- CloudBees Feature Management account

### Installation

1. Clone the repository:
```bash
cd /Users/brown/git_orgs/CB-InsuranceStack/InsuranceStack/apps/claims-service
```

2. Install dependencies:
```bash
make deps
```

3. Set up environment variables:
```bash
export CLOUDBEES_FM_API_KEY="your-api-key-here"
export DATA_PATH="../../data/seed"
export LOG_LEVEL="info"
export FEATURE_AUTO_APPROVAL="false"
```

### Running Locally

#### Development Mode
```bash
make run-dev
```

#### Production Mode
```bash
make build
make run
```

#### With Custom Configuration
```bash
export PORT=8002
export LOG_LEVEL=debug
export DATA_PATH=/path/to/data
go run cmd/server/main.go
```

### Running with Docker

#### Build Docker Image
```bash
make docker-build
```

#### Run Docker Container
```bash
make docker-run
```

Or manually:
```bash
docker build -t insurancestack/claims-service:latest .

docker run -p 8002:8002 \
  -e CLOUDBEES_FM_API_KEY="your-api-key" \
  -e DATA_PATH=/data/seed \
  -v $(pwd)/../../data/seed:/data/seed \
  insurancestack/claims-service:latest
```

## Project Structure

```
claims-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── features/
│   │   └── flags.go             # Feature flag management
│   ├── handlers/
│   │   ├── health.go            # Health check handler
│   │   └── claim.go             # Claims handlers
│   ├── middleware/
│   │   ├── auth.go              # Authentication middleware
│   │   ├── cors.go              # CORS middleware
│   │   └── logging.go           # Logging middleware
│   ├── models/
│   │   └── claim.go             # Claim data models
│   ├── repository/
│   │   └── repository.go        # Data access layer
│   └── services/
│       └── claim_service.go     # Business logic
├── Dockerfile                    # Docker configuration
├── Makefile                      # Build automation
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
└── README.md                     # This file
```

## Development

### Building
```bash
make build
```

### Running Tests
```bash
make test
```

### Code Formatting
```bash
make fmt
```

### Linting
```bash
make lint
```

### Cleaning Build Artifacts
```bash
make clean
```

## Authentication

The service uses a simple header-based authentication for demo purposes:

```bash
curl -H "X-User-ID: user-001" "http://localhost:8002/claims?policyId=pol-001"
```

If no `X-User-ID` header is provided, it defaults to `user-001`.

**Note:** In production, this should be replaced with proper JWT token validation or session-based authentication.

## Claim Types

Available claim types:
- `accident` - Vehicle accidents and collisions
- `theft` - Theft or vandalism claims
- `damage` - Property or vehicle damage

## Claim Status

- `submitted` - Claim has been submitted and awaiting review
- `under_review` - Claim is being reviewed by an adjuster
- `approved` - Claim has been approved for payment
- `rejected` - Claim has been denied

## Governance Workflow

This service demonstrates governance and approval workflows:

1. **Claim Submission**: Customer submits a claim with details and amount
2. **Automatic Triage**: Low-value claims (< $1000) can be auto-approved if feature flag is enabled
3. **Manual Review**: High-value claims require explicit approval via status change endpoint
4. **Audit Trail**: All claim status changes are logged with timestamps
5. **Approval Requirements**: Different thresholds can be enforced for different claim amounts

## Logging

The service uses structured JSON logging with the following levels:
- `debug` - Detailed debug information
- `info` - General informational messages
- `warn` - Warning messages
- `error` - Error messages

Each request is logged with:
- HTTP method
- Path
- Status code
- Duration
- Remote address
- User agent

Example log entry:
```json
{
  "level": "info",
  "method": "GET",
  "path": "/transactions",
  "status": 200,
  "duration": "2.345ms",
  "remote": "127.0.0.1",
  "msg": "HTTP request"
}
```

## Health Checks

The `/healthz` endpoint can be used for:
- Kubernetes liveness/readiness probes
- Load balancer health checks
- Monitoring systems

The Docker image includes a built-in health check that runs every 30 seconds.

## CloudBees Feature Management Integration

This service integrates with CloudBees Feature Management using the Rox SDK.

### Setup

1. Create a CloudBees Feature Management account
2. Create a new application in the dashboard
3. Copy your API key
4. Set the `CLOUDBEES_FM_API_KEY` environment variable

### Feature Flag Configuration

In the CloudBees dashboard:
1. Navigate to your application
2. Create a new flag named `claims.autoApproval`
3. Set the default value to `false`
4. Configure targeting rules as needed
5. Deploy the configuration

### Testing Feature Flags

With flag disabled (default):
```bash
# All claims require manual review
curl -X POST "http://localhost:8002/claims" \
  -H "Content-Type: application/json" \
  -d '{"policyId":"pol-001","customerId":"cust-001","type":"accident","amount":500,"description":"Minor accident"}'
# Response will show status: "under_review"
```

Enable the flag in CloudBees dashboard, then:
```bash
# Low-value claims are auto-approved
curl -X POST "http://localhost:8002/claims" \
  -H "Content-Type: application/json" \
  -d '{"policyId":"pol-001","customerId":"cust-001","type":"accident","amount":500,"description":"Minor accident"}'
# Response will show status: "approved" (if amount < $1000)
```

## Error Handling

The API returns standard HTTP status codes:
- `200 OK` - Successful request
- `400 Bad Request` - Invalid parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

Error responses include a message:
```json
{
  "error": "Invalid startDate format. Use ISO 8601 (YYYY-MM-DD or RFC3339)"
}
```

## CORS Configuration

The service is configured to accept requests from any origin with:
- Methods: GET, POST, PUT, DELETE, OPTIONS
- Headers: Accept, Authorization, Content-Type, X-CSRF-Token, X-User-ID
- Credentials: Enabled

**Note:** In production, configure `AllowedOrigins` to specific domains.

## Performance Considerations

- Claims are loaded into memory from JSON files on startup
- Read operations are protected with RWMutex for thread safety
- Results are sorted by submission date (most recent first)
- No database required for this demo service

## License

Copyright (c) 2024 CB-InsuranceStack

## Support

For issues and questions, please contact the development team.
