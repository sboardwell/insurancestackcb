# InsuranceStack - Docker Quick Start

## Prerequisites
- Docker and Docker Compose installed
- Ports 8001-8005 available

## Quick Start

### 1. Build all services
```bash
docker-compose build
```

### 2. Start backend services
```bash
docker-compose up -d policy-service claims-service pricing-engine customer-service payments-service
```

### 3. Verify all services are running
```bash
docker-compose ps
```

### 4. Test health endpoints
```bash
for port in 8001 8002 8003 8004 8005; do
  echo "Testing port $port..."
  curl -s http://localhost:$port/healthz | jq '.'
done
```

## Service Ports

| Service | Port | Description |
|---------|------|-------------|
| policy-service | 8001 | Insurance policies management |
| claims-service | 8002 | Claims processing and approval |
| pricing-engine | 8003 | Quote calculation (FAST-PATH) |
| customer-service | 8004 | Customer profile management |
| payments-service | 8005 | Payment processing (HIGHEST RISK) |

## API Examples

### Get all policies
```bash
curl http://localhost:8001/policies
```

### Get all claims
```bash
curl http://localhost:8002/claims
```

### Calculate insurance quote
```bash
curl -X POST http://localhost:8003/quote \
  -H "Content-Type: application/json" \
  -d '{"policyType":"auto","coverageAmount":50000,"customerId":"cust-001"}'
```

### Get customers
```bash
curl http://localhost:8004/customers
```

### Get payments
```bash
curl http://localhost:8005/payments
```

## Environment Variables

All services support these environment variables (set in docker-compose.yaml):

- `PORT` - Service port
- `LOG_LEVEL` - Logging level (debug, info, warn, error)
- `CLOUDBEES_FM_API_KEY` - CloudBees Feature Management API key
- `DATA_PATH` - Path to seed data files
- `JWT_SECRET` - JWT signing secret

## Feature Flags

Each service includes CloudBees Feature Management integration:

- **policy-service**: `policy.maskPremiums` - Controls premium visibility
- **claims-service**: `claims.autoApproval` - Auto-approve claims under $1000
- **pricing-engine**: `pricing.dynamicRates` - Real-time rate adjustments  
- **payments-service**: `payments.instantPayouts` - Instant vs batch payouts

## Governance Tiers

Services are classified by governance requirements:

- **FAST-PATH**: pricing-engine (minimal checks, quick deployment)
- **STANDARD**: policy-service, customer-service (normal governance)
- **GOVERNED**: claims-service, payments-service (enhanced checks, requires approval)
- **HIGHEST RISK**: payments-service (critical, restricted deployment windows)

## Stopping Services

```bash
docker-compose down
```

## Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f policy-service
```

## Environment-Agnostic Builds

All Docker images are built without environment-specific configuration:
- No BASE_PATH or environment variables baked into images
- Same image deploys to dev, staging, production
- Environment-specific config injected at runtime

This ensures:
- Build once, deploy anywhere
- Consistent behavior across environments
- Simplified CI/CD pipelines
