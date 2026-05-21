# InsuranceStack

A reference insurance platform for demonstrating governed CI/CD at scale in regulated enterprise environments.

## Overview

InsuranceStack is a realistic **insurance management platform** demonstrating how insurance companies manage policies, claims, pricing, and payments — designed for enablement, workshops, and executive demonstrations of modern CI/CD governance.

> "A reference insurance platform for demonstrating governed CI/CD at scale"

**Key Features**:
- Multi-service microservices architecture (6 services + UI)
- Realistic insurance domain model (policies, claims, payments, pricing)
- Differentiated governance workflows (fast-path vs. gated deployments)
- Local-first development (works offline)
- Built-in feature management (CloudBees FM integration)
- Comprehensive test coverage for SmartTests demonstrations
- CloudBees Unify CI/CD workflows with promotion and evidence collection

## Why Insurance?

Insurance is the perfect domain for demonstrating enterprise CI/CD because:
- **Universally understood** - Everyone knows policies, claims, and payments
- **Regulation-heavy** - Perfect for governance and compliance demos
- **Clear service boundaries** - Natural microservice decomposition
- **Risk differentiation** - Some services need more controls than others
- **Real-world complexity** - Feels like actual enterprise systems

## Quick Start

### Local Development with Docker Compose

1. **Set up environment** (optional - for CloudBees FM integration):

```bash
# Copy example environment file
cp .env.example .env

# Edit .env and add your CloudBees FM API key
# CLOUDBEES_FM_API_KEY=your-actual-fm-key-here
```

2. **Start all services**:

```bash
docker compose up --build
```

This starts:
- Insurance UI (port 3000)
- Policy Service (port 8001)
- Claims Service (port 8002)
- Pricing Engine (port 8003)
- Customer Service (port 8004)
- Payments Service (port 8005)

3. **Access the application**: http://localhost:3000

**Note**: The `.env` file is gitignored for security. Without an FM key, the app works perfectly with hardcoded default flag values.

## Service Architecture

### Core Services

- **apps/policy-service** (port 8001)
  - Manages insurance policies
  - Depends on pricing-engine
  - Exposes `/policies`, `/policies/{id}`

- **apps/claims-service** (port 8002)
  - Handles insurance claims
  - Depends on policy-service
  - **Higher governance requirements** - Manual approvals, extra testing
  - Exposes `/claims`, `/claims/{id}`

- **apps/pricing-engine** (port 8003)
  - Calculates insurance premiums
  - Stateless service
  - **Fast-path deployment** - Frequently changed, minimal gates
  - Exposes `/calculate`, `/rates`

- **apps/customer-service** (port 8004)
  - Customer profiles and management
  - Referenced by policy and claims services
  - Exposes `/customers`, `/customers/{id}`

- **apps/payments-service** (port 8005)
  - Simulates premium payments and claim payouts
  - **Highest risk classification** - Restricted deployment windows
  - Exposes `/payments`, `/payouts`

### Frontend

- **apps/insurance-ui** (port 3000)
  - React frontend
  - Policy management, claims submission, customer portal

## Domain Model

### Policy States
- **Active** - Policy is in force
- **Lapsed** - Missed payments
- **Cancelled** - Policy terminated

### Claim States
- **Submitted** - Initial claim filing
- **Under Review** - Being assessed
- **Approved** - Claim accepted, payment authorized
- **Rejected** - Claim denied

### Governance Tiers

**Fast Path** (pricing-engine):
- Auto-deploy to dev/test
- Minimal approval gates
- Unit tests only
- Changes affect new policies only

**Standard Path** (policy-service, customer-service):
- Auto-deploy to dev
- Manual approval for production
- Full test suite required

**Governed Path** (claims-service, payments-service):
- Manual approval at every stage
- Security scanning required
- Evidence collection mandatory
- Restricted deployment windows (payments only)
- Extended testing requirements

## Feature Management

InsuranceStack uses **CloudBees Feature Management** for all feature flags with **fully reactive, real-time updates**.

### Key Insurance Flags

- **`claims.autoApproval`** - Toggle automatic claim approval for low-value claims
- **`pricing.dynamicRates`** - Switch between fixed and dynamic pricing algorithms
- **`ui.claimsFastTrack`** - Show/hide expedited claims processing
- **`payments.instantPayouts`** - Enable instant vs. batch payouts
- **`policies.cancelSelfService`** - Allow customers to cancel policies online

### Live Demo Flow
1. Open InsuranceStack in browser
2. Open CloudBees FM dashboard
3. Toggle `claims.autoApproval` flag
4. Watch claims approval workflow change instantly

See [Feature Flags Reference](config/README.md) for complete flag list.

## Testing

Run all tests locally:

```bash
make test
```

Or run specific test suites:

```bash
make test-unit          # ~150-200 unit tests
make test-integration   # ~50-80 integration tests
make test-e2e          # ~20-30 end-to-end tests
```

High test volume demonstrates CloudBees SmartTests impact analysis and test subsetting.

## CI/CD & Governance

### CloudBees Unify Workflows

InsuranceStack demonstrates:
- **Reusable workflows** - Shared pipeline templates
- **Environment promotion** - dev → test → prod with evidence
- **Manual approval gates** - Required for governed services
- **Evidence collection** - Test results, scan reports, approvals
- **Component dependencies** - Claims depends on policy, policy depends on pricing

### Deployment Examples

**Fast Path** (pricing-engine):
```
Commit → Build → Test → Deploy Dev (auto) → Deploy Test (auto) → Deploy Prod (auto)
Time: ~10 minutes
```

**Governed Path** (claims-service):
```
Commit → Build → Test → Security Scan → Deploy Dev (auto)
  → Approval → Deploy Test (manual)
  → Approval + Evidence → Deploy Prod (manual, restricted window)
Time: ~2-4 hours (with approvals)
```

## Documentation

- [Architecture Details](docs/ARCHITECTURE.md) - Component design and tech stack
- [Deployment Guide](DEPLOYMENT.md) - Kubernetes and Helm setup
- [Changelog](CHANGELOG.md) - Version history and migration guides
- [Demo Flow](docs/DEMO.md) - 20-minute demonstration script

## Technology Stack

- **Frontend**: React, TypeScript, Vite
- **Backend**: Go (all microservices)
- **Infrastructure**: Docker, Docker Compose, Kubernetes, Helm
- **CI/CD**: CloudBees Unify
- **Feature Management**: CloudBees Feature Management (optional)
- **Testing**: Jest, Playwright, Go testing

## Positioning

**InsuranceStack** is for:
- Governance and compliance demonstrations
- CI/CD pipeline differentiation (fast vs. gated)
- Evidence collection and audit trail demos
- Promotion workflow demonstrations
- Executive demonstrations of regulated CI/CD

**Based on**: InsuranceStack reference architecture
**Related**: SquidStack (for deeper platform demonstrations)

## Non-Goals

- No real payment gateways or PCI compliance
- No real identity providers or SSO
- No production-grade security hardening
- No external SaaS dependencies required
- No actual insurance calculations or underwriting

This is about **demonstrating CI/CD process, not production insurance systems**.

## Success Criteria

This project succeeds if:
- A full governance flow can be demoed in 20-30 minutes
- Governance differences between services are obvious
- Promotion and evidence collection are tangible
- It feels realistic to enterprise engineers
- Fast-path vs. gated deployments are clearly differentiated

## License

MIT

## Support

For issues or questions, please open a GitHub issue or contact the CloudBees SE team.
