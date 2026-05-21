# InsuranceStack Helm Chart

A Helm chart for deploying the InsuranceStack application with CloudBees Feature Management integration.

## Overview

This Helm chart deploys a complete insurance platform with 6 microservices, each configured with appropriate governance tiers:

### Services

| Service | Port | Governance Tier | Description |
|---------|------|----------------|-------------|
| **insurance-ui** | 3000 | standard | React-based web interface |
| **policy-service** | 8001 | standard | Policy management service |
| **claims-service** | 8002 | governed (high) | Claims processing with high governance |
| **pricing-engine** | 8003 | fast-path (minimal) | Real-time pricing calculations |
| **customer-service** | 8004 | standard | Customer management |
| **payments-service** | 8005 | governed (highest-risk) | Payment processing with highest security |

## Governance Tiers

The chart implements three governance tiers through labels and annotations:

- **fast-path**: Minimal governance for high-performance services (pricing-engine)
- **standard**: Standard governance for regular services (policy, customer, UI)
- **governed**: High governance for sensitive operations (claims, payments)

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- cert-manager (for TLS certificates)
- nginx-ingress-controller

## Installation

### Basic Installation

```bash
helm install insurancestack ./insurancestack
```

### Production Installation

```bash
helm install insurancestack ./insurancestack \
  --set deployment.baseDomain=insurancestack.example.com \
  --set deployment.organizationName=CB-InsuranceStack \
  --set deployment.environmentName=production \
  --set auth.jwtSecret="your-secure-jwt-secret" \
  --set cloudbees.fmKey="your-cloudbees-fm-api-key" \
  --set featureFlags.fmKey="your-cloudbees-fm-api-key"
```

### Development Installation

```bash
helm install insurancestack ./insurancestack \
  --set deployment.baseDomain=localhost \
  --set insuranceUI.showDemoCredentials=true
```

## Configuration

### Global Settings

| Parameter | Description | Default |
|-----------|-------------|---------|
| `global.imageRegistry` | Global Docker registry | `""` |
| `global.imagePullSecrets` | Image pull secrets | `[]` |

### Deployment Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `deployment.baseDomain` | Base domain for ingress | `insurancestack.se-main-demo.sa-demo.beescloud.com` |
| `deployment.organizationName` | Organization name for URL path | `""` |
| `deployment.environmentName` | Environment name for URL path | `""` |

### Authentication

| Parameter | Description | Default |
|-----------|-------------|---------|
| `auth.jwtSecret` | JWT secret for token signing | `dev-secret-key-change-in-production-use-a-long-random-string` |
| `auth.username` | Demo username | `demo@insurancestack.com` |
| `auth.password` | Demo password | `demo123` |

### CloudBees Feature Management

| Parameter | Description | Default |
|-----------|-------------|---------|
| `cloudbees.fmKey` | CloudBees FM API key | `local-mode` |
| `cloudbees.environment` | CloudBees FM environment | `production` |
| `featureFlags.enabled` | Enable feature flags | `true` |

### Service-Specific Configuration

Each service has the following configurable parameters:

- `enabled`: Enable/disable the service (default: `true`)
- `replicaCount`: Number of replicas (default: `2`)
- `governanceTier`: Governance tier label
- `image.repository`: Docker image repository
- `image.tag`: Image tag (default: `latest`)
- `image.pullPolicy`: Image pull policy (default: `IfNotPresent`)
- `service.type`: Kubernetes service type (default: `ClusterIP`)
- `service.port`: Service port
- `service.targetPort`: Container port
- `resources`: Resource requests and limits
- `env`: Service-specific environment variables

### Environment Variables by Service

#### Policy Service (8001)
- `DATA_PATH`: `/data/policies`
- `FEATURE_APPROVAL_WORKFLOW`: `true`
- `FEATURE_RISK_ASSESSMENT`: `true`
- `FEATURE_MULTI_PRODUCT`: `false`

#### Claims Service (8002)
- `DATA_PATH`: `/data/claims`
- `FEATURE_AUTO_APPROVAL`: `false`
- `FEATURE_FRAUD_DETECTION`: `true`
- `FEATURE_ADVANCED_ANALYTICS`: `true`

#### Pricing Engine (8003)
- `DATA_PATH`: `/data/pricing`
- `FEATURE_DYNAMIC_PRICING`: `true`
- `FEATURE_REALTIME_QUOTES`: `true`
- `FEATURE_BULK_DISCOUNT`: `false`

#### Customer Service (8004)
- `DATA_PATH`: `/data/customers`
- `FEATURE_KYC`: `true`
- `FEATURE_DOCUMENT_UPLOAD`: `true`
- `FEATURE_COMMUNICATION_PREFS`: `false`

#### Payments Service (8005)
- `DATA_PATH`: `/data/payments`
- `FEATURE_RECURRING_PAYMENTS`: `true`
- `FEATURE_MULTI_CURRENCY`: `false`
- `FEATURE_REFUND_PROCESSING`: `true`

## Labels and Annotations

### Governance Labels

All services include the following governance labels:

```yaml
governance.insurancestack.io/tier: [fast-path|standard|governed]
governance.insurancestack.io/level: [minimal|high|highest-risk]  # For governed services
```

### Annotations

Services include descriptive annotations:

```yaml
governance.insurancestack.io/description: "Service description with governance information"
```

## Upgrading

```bash
helm upgrade insurancestack ./insurancestack \
  --set cloudbees.fmKey="new-api-key"
```

## Uninstalling

```bash
helm uninstall insurancestack
```

## Advanced Configuration

### Custom Image Tags

```bash
helm install insurancestack ./insurancestack \
  --set policyService.image.tag=v1.2.3 \
  --set claimsService.image.tag=v1.2.3 \
  --set pricingEngine.image.tag=v1.2.3 \
  --set customerService.image.tag=v1.2.3 \
  --set paymentsService.image.tag=v1.2.3 \
  --set insuranceUI.image.tag=v1.2.3
```

### Resource Limits

```bash
helm install insurancestack ./insurancestack \
  --set pricingEngine.resources.limits.cpu=2000m \
  --set pricingEngine.resources.limits.memory=2Gi
```

### Disable Specific Services

```bash
helm install insurancestack ./insurancestack \
  --set pricingEngine.enabled=false
```

## Health Checks

All backend services include health check endpoints:

- Liveness probe: `GET /healthz` (initial delay: 10s)
- Readiness probe: `GET /healthz` (initial delay: 5s)

## Security Considerations

1. **Change default credentials** in production:
   - Set `auth.jwtSecret` to a secure random string
   - Update `auth.username` and `auth.password`

2. **Governance tiers** are configured based on risk:
   - Payments and Claims use "governed" tier
   - Pricing uses "fast-path" for performance
   - Others use "standard" tier

3. **Secrets** are managed through Kubernetes secrets:
   - JWT secrets
   - CloudBees FM API keys
   - Authentication credentials

## Troubleshooting

### Check pod status
```bash
kubectl get pods -l app.kubernetes.io/name=insurancestack
```

### View logs
```bash
kubectl logs -l app.kubernetes.io/component=policy-service
```

### Check governance labels
```bash
kubectl get pods -l governance.insurancestack.io/tier=governed
```

## License

Copyright (c) 2024 InsuranceStack Team
