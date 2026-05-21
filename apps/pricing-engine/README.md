# Pricing Engine Service

A FAST-PATH, production-ready Go API service for insurance pricing calculations with CloudBees Feature Management integration.

## Overview

The Pricing Engine service provides real-time insurance quote calculations for auto, home, and life insurance policies. Built for speed and minimal governance, this service emphasizes rapid deployment and iteration while maintaining production quality.

## FAST-PATH Service

This is a **FAST-PATH** service, designed for:
- **Rapid Development**: Quick iterations without extensive approval processes
- **Speed over Perfection**: Focus on getting features to market quickly
- **Minimal Governance**: Streamlined reviews and faster deployment cycles
- **Production Quality**: Fast doesn't mean sloppy - maintains high standards
- **Feature Flag Control**: Dynamic pricing adjustments without redeployment

## Features

- RESTful API for insurance quote calculations
- Support for auto, home, and life insurance policies
- Multi-factor pricing based on coverage, age, and risk score
- Discount calculations (multi-policy, loyalty, paperless billing)
- Real-time feature flag system using CloudBees Feature Management
- Feature flag: `pricing.dynamicRates` - enable/disable real-time rate adjustments based on seasonality, market conditions, and claims history
- Proper error handling and structured logging
- CORS support
- Graceful shutdown
- Health check endpoint
- JSON structured logging

## Project Structure

```
pricing-engine/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── health.go           # Health check handler
│   │   └── pricing.go          # Pricing endpoints
│   ├── services/                # Business logic
│   │   └── pricing_service.go  # Pricing calculations
│   ├── repository/              # Data access layer
│   │   └── repository.go       # Pricing rules loader
│   ├── features/                # Feature flags
│   │   └── flags.go            # CloudBees FM/Rox integration
│   ├── models/                  # Data models
│   │   └── pricing.go          # Quote, Rate, and pricing models
│   ├── middleware/              # HTTP middleware
│   │   ├── logging.go          # Request logging
│   │   ├── cors.go             # CORS configuration
│   │   └── auth.go             # Authentication
│   └── auth/                    # Authentication utilities
│       ├── jwt.go              # JWT token management
│       └── password.go         # Password hashing
├── go.mod                       # Go module definition
├── Dockerfile                   # Docker configuration
├── Makefile                     # Build automation
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
  "service": "pricing-engine"
}
```

### Calculate Quote

**POST /quote**

Calculates an insurance quote based on policy parameters.

**Request Body:**
```json
{
  "policyType": "auto",
  "coverageAmount": 500000,
  "customerAge": 35,
  "riskScore": 2,
  "customerId": "CUST-12345",
  "multiPolicy": true,
  "loyaltyYears": 5,
  "paperlessBill": true,
  "claimsHistory": 0
}
```

**Response:**
```json
{
  "quoteId": "Q-a3b4c5d6",
  "policyType": "auto",
  "coverageAmount": 500000,
  "baseRate": 1240.0,
  "adjustedRate": 1178.0,
  "discount": 200.88,
  "finalPremium": 977.12,
  "validUntil": "2025-01-20T10:30:00Z",
  "createdAt": "2024-12-21T10:30:00Z",
  "factors": {
    "baseMultiplier": 800.0,
    "coverageMultiplier": 1.55,
    "ageMultiplier": 1.0,
    "riskMultiplier": 1.0,
    "dynamicMultiplier": 0.95,
    "discountAmount": 200.88
  }
}
```

**Policy Types:**
- `auto`: Auto insurance
- `home`: Home insurance
- `life`: Life insurance

**Risk Score:** 1-5 (1 = lowest risk, 5 = highest risk)

**Coverage Amounts:**
- Auto: 250000, 300000, 400000, 500000
- Home: 500000, 650000, 750000, 1000000, 1200000
- Life: 250000, 500000, 750000, 1000000

### Get Base Rates

**GET /rates**

Returns current base rates for all policy types.

**Response:**
```json
{
  "rates": [
    {
      "policyType": "auto",
      "baseRate": 800,
      "coverage": {
        "250000": 1.0,
        "300000": 1.15,
        "400000": 1.35,
        "500000": 1.55
      }
    },
    {
      "policyType": "home",
      "baseRate": 1200,
      "coverage": {
        "500000": 1.0,
        "650000": 1.2,
        "750000": 1.4,
        "1000000": 1.8,
        "1200000": 2.2
      }
    },
    {
      "policyType": "life",
      "baseRate": 500,
      "coverage": {
        "250000": 1.0,
        "500000": 1.8,
        "750000": 2.5,
        "1000000": 3.2
      }
    }
  ],
  "timestamp": "2024-12-21T10:30:00Z"
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8003` |
| `DATA_PATH` | Path to seed data directory | `../../data/seed` |
| `CLOUDBEES_FM_API_KEY` | CloudBees Feature Management API key | `dev-mode` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `JWT_SECRET` | JWT signing secret | `dev-secret-key-change-in-production` |
| `FEATURE_DYNAMIC_RATES` | Enable dynamic rates in dev mode (true/false) | `false` |

## Feature Flags

### pricing.dynamicRates

**Default:** `false`

When enabled, dynamic pricing adjustments are applied based on:
- **Seasonality**: Q1 (0.95x), Q2 (1.0x), Q3 (1.05x), Q4 (1.02x)
- **Market Conditions**: Global multiplier for market volatility
- **Claims History**: Adjustments based on customer claims (0 claims: 0.9x, 1: 1.0x, 2: 1.15x, 3+: 1.3x)

**Use Cases:**
- A/B testing dynamic pricing strategies
- Seasonal rate adjustments
- Market-responsive pricing
- Instant rollback if pricing issues detected

### CloudBees Integration

The service uses the CloudBees Rox SDK (`github.com/rollout/rox-go/v5/core`) for real-time feature flag management. Flags can be toggled instantly without redeploying the service.

**Development Mode:** When no CloudBees API key is provided, the service falls back to environment variables (`FEATURE_DYNAMIC_RATES`).

**Production Mode:** Provide `CLOUDBEES_FM_API_KEY` to use CloudBees Feature Management for centralized control and real-time updates.

## Pricing Calculation

The pricing engine calculates quotes using the following formula:

```
Base Premium = Base Rate × Coverage Multiplier × Age Multiplier × Risk Multiplier

Dynamic Adjustment = Base Premium × Dynamic Multiplier (if enabled)

Total Discounts = (Multi-Policy + Loyalty + Low Risk + Paperless) discounts

Final Premium = Dynamic Adjusted Premium - Total Discounts
```

### Discounts

- **Multi-Policy**: 15% discount for customers with multiple policies
- **Loyalty Years**: 0% (1yr), 5% (2yr), 8% (3yr), 12% (5yr), 18% (10yr+)
- **Low Risk**: 10% for risk score of 1
- **Paperless Billing**: 3% for opting into paperless billing

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Access to seed data files at `/data/seed/pricing-rules.json`
- CloudBees Feature Management account (optional, for production)

### Installation

1. Install dependencies:

```bash
cd /Users/brown/git_orgs/CB-InsuranceStack/InsuranceStack/apps/pricing-engine
go mod download
```

2. Set up environment variables:

```bash
export PORT=8003
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
go build -o bin/pricing-engine cmd/server/main.go
./bin/pricing-engine
```

### Using Make

```bash
# Build the application
make build

# Run in development mode
make run-dev

# Run tests
make test

# Format code
make fmt

# Build Docker image
make docker-build

# Run in Docker
make docker-run

# Show all available commands
make help
```

## Testing the API

### Health Check
```bash
curl http://localhost:8003/healthz
```

### Get Base Rates
```bash
curl http://localhost:8003/rates
```

### Calculate Auto Insurance Quote
```bash
curl -X POST http://localhost:8003/quote \
  -H "Content-Type: application/json" \
  -d '{
    "policyType": "auto",
    "coverageAmount": 500000,
    "customerAge": 35,
    "riskScore": 2,
    "multiPolicy": true,
    "loyaltyYears": 5,
    "paperlessBill": true,
    "claimsHistory": 0
  }'
```

### Calculate Home Insurance Quote
```bash
curl -X POST http://localhost:8003/quote \
  -H "Content-Type: application/json" \
  -d '{
    "policyType": "home",
    "coverageAmount": 750000,
    "customerAge": 42,
    "riskScore": 1,
    "multiPolicy": false,
    "loyaltyYears": 3,
    "paperlessBill": true,
    "claimsHistory": 1
  }'
```

### Test Dynamic Rates (Enabled)
```bash
# Start server with dynamic rates enabled
export FEATURE_DYNAMIC_RATES=true
go run cmd/server/main.go &

# Calculate quote - rates will be adjusted by seasonality and claims history
curl -X POST http://localhost:8003/quote \
  -H "Content-Type: application/json" \
  -d '{
    "policyType": "life",
    "coverageAmount": 1000000,
    "customerAge": 45,
    "riskScore": 3,
    "claimsHistory": 2
  }'
```

## Development

### FAST-PATH Development Principles

1. **Speed First**: Prioritize rapid development and deployment
2. **Iterate Quickly**: Ship features fast, improve based on feedback
3. **Feature Flags**: Use flags for gradual rollouts and quick rollbacks
4. **Minimal Process**: Streamlined reviews, trust the team
5. **Production Quality**: Fast doesn't mean broken - maintain standards

### Project Layout

The service follows Go best practices with a clean architecture:

1. **Handlers Layer** (`internal/handlers/`): HTTP request/response handling
2. **Services Layer** (`internal/services/`): Business logic and pricing calculations
3. **Repository Layer** (`internal/repository/`): Data access and rule loading
4. **Models Layer** (`internal/models/`): Domain models and data structures
5. **Features Layer** (`internal/features/`): Feature flag management

### Middleware

- **Logging**: Logs all HTTP requests with method, path, status, and duration
- **CORS**: Handles cross-origin resource sharing
- **Auth**: JWT token validation (bypassed for health check)

### Feature Flag Architecture

Feature flags are initialized on startup and can be updated in real-time via CloudBees. The service checks flag status on each request, allowing instant behavior changes without downtime.

```go
// Check if dynamic rates should be applied
if flags.IsDynamicRatesEnabled() {
    dynamicMultiplier = calculateDynamicMultiplier(req)
}
```

## Docker

### Build Image
```bash
docker build -t insurancestack/pricing-engine:latest .
```

### Run Container
```bash
docker run -p 8003:8003 \
  -e PORT=8003 \
  -e LOG_LEVEL=info \
  -e CLOUDBEES_FM_API_KEY=your-key \
  -v /path/to/data/seed:/data/seed \
  -e DATA_PATH=/data/seed \
  insurancestack/pricing-engine:latest
```

## Production Considerations

1. **Authentication**: Implement proper JWT token validation or OAuth2.

2. **CORS**: The CORS middleware currently allows all origins (`*`). In production, specify exact allowed origins.

3. **Database**: Pricing rules are loaded from JSON files. In production, consider storing rules in a database for easier updates.

4. **Monitoring**: Add metrics collection (Prometheus), distributed tracing (OpenTelemetry), and error tracking (Sentry).

5. **Rate Limiting**: Implement rate limiting to prevent abuse.

6. **TLS**: Enable HTTPS with proper certificates.

7. **Secrets Management**: Use a secrets manager (AWS Secrets Manager, HashiCorp Vault) for sensitive configuration.

8. **Feature Flag Management**: Use CloudBees Feature Management dashboard to control flags across environments.

9. **Caching**: Consider caching pricing rules for improved performance.

## Data Model

### QuoteRequest
- `policyType`: Policy type (auto, home, life)
- `coverageAmount`: Coverage amount in dollars
- `customerAge`: Customer age (18-120)
- `riskScore`: Risk score (1-5)
- `customerId`: Customer identifier (optional)
- `multiPolicy`: Multi-policy discount flag
- `loyaltyYears`: Years of customer loyalty
- `paperlessBill`: Paperless billing flag
- `claimsHistory`: Number of previous claims

### Quote
- `quoteId`: Unique quote identifier
- `policyType`: Policy type
- `coverageAmount`: Coverage amount
- `baseRate`: Base calculated rate
- `adjustedRate`: Rate after dynamic adjustments
- `discount`: Total discount amount
- `finalPremium`: Final premium after discounts
- `validUntil`: Quote expiration date
- `createdAt`: Quote creation timestamp
- `factors`: Breakdown of pricing factors

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests manually
go test -v -race ./...
```

## License

Copyright 2024 CB-InsuranceStack

## Quick Start

The fastest way to get started:

```bash
# From the pricing-engine directory
export DATA_PATH=../../data/seed
go run cmd/server/main.go
```

Then test the API:

```bash
# Health check
curl http://localhost:8003/healthz

# Get base rates
curl http://localhost:8003/rates

# Calculate a quote
curl -X POST http://localhost:8003/quote \
  -H "Content-Type: application/json" \
  -d '{
    "policyType": "auto",
    "coverageAmount": 500000,
    "customerAge": 35,
    "riskScore": 2,
    "multiPolicy": true,
    "loyaltyYears": 5,
    "paperlessBill": true
  }'
```

## CloudBees Feature Management Setup

1. Sign up for CloudBees Feature Management
2. Create a new application
3. Get your API key
4. Set the environment variable:
   ```bash
   export CLOUDBEES_FM_API_KEY=your-actual-api-key
   ```
5. Run the service - it will automatically register flags with CloudBees
6. Toggle `pricing.dynamicRates` in the CloudBees dashboard - changes apply instantly!

## FAST-PATH Success Metrics

- **Time to Market**: Days, not weeks or months
- **Deployment Frequency**: Multiple times per week
- **Feature Flag Usage**: Active feature flagging for safe, fast deployments
- **Customer Feedback Loop**: Rapid iteration based on real-world usage
- **Team Velocity**: High throughput with maintained quality

## Support

For issues or questions, please contact the development team.
