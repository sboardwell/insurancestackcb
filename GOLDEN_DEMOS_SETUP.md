# Golden Demos Setup Progress

This document tracks the progress of adding InsuranceStack and AccountStack to CloudBees Golden Demos (Mimic).

## Overview

**Goal**: Integrate InsuranceStack and AccountStack into CloudBees Golden Demos platform to enable one-command deployment for demos.

**Applications**:
- InsuranceStack (this repo)
- AccountStack (CB-AccountStack/AccountStack)

## Progress Log

### 2025-12-23

#### KUBECONFIG Secret Migration ✅
- **Status**: **SUCCESS** - Deployment working with org-level property
- **Change**: Migrated KUBECONFIG from environment-level secret to organization-level property
- **Rationale**: Allows reuse across multiple environments and projects
- **Key Finding**: CloudBees Unify resolves org-level properties even when Helm charts reference environment-level secrets
  - Helm chart expects secret at environment level
  - Works seamlessly with org-level property instead
  - Simplifies golden demos setup (one KUBECONFIG for all demos)
- **Testing**: ✅ Verified in InsuranceStack deployment
- **Next Steps**:
  - Apply same pattern to AccountStack
  - Consider moving FM_KEY to org-level property as well
  - Document this secret resolution behavior for Mimic integration

## Configuration Requirements

### Secrets & Properties

| Name | Type | Level | Status | Notes |
|------|------|-------|--------|-------|
| KUBECONFIG | Property | Organization | ✅ Working | Org-level property works despite Helm expecting env-level |
| FM_KEY | Secret | Repository | Configured | CloudBees Feature Management API key |

### Deployment Pattern

Both applications use:
- **Architecture**: Full-stack monorepo with multiple components
- **Database**: In-memory (no external database required)
- **Ingress**: Path-based routing for multi-tenancy
- **Feature Management**: Reactive SSE for instant flag updates

## Mimic Integration Plan

### Phase 1: Configuration Setup
- [x] Validate KUBECONFIG as org-level property
- [x] Document required secrets and properties
- [x] Test deployment workflows with new configuration

### Phase 2: Mimic Scenario Creation
- [x] Create `insurancestack-demo` scenario ✅
- [ ] Create `accountstack-demo` scenario (deferred)
- [x] Test end-to-end deployment via Mimic ✅
- [x] Document setup requirements ✅

#### InsuranceStack Scenario Created (2025-12-23)
- **Location**: `/Users/brown/Git_orgs/cb-demos/mimic/scenarios/insurancestack-demo.yaml`
- **Approach**: Minimal scenario with NO file replacements
- **Key Insight**: CloudBees workflow overrides all Helm values at deployment time via `--set` flags
- **Configuration**:
  - Required secrets: `KUBECONFIG` (org-level)
  - Parameters: `project_name`, `target_org`, `environment`
  - Creates: Repository fork, Component, Environment (with FM_TOKEN), Application
  - No file modifications needed - workflow handles all configuration dynamically

#### Successful Test Run (2025-12-23)
- **Run ID**: insurancestack-demo-02a1329e
- **Target CloudBees Org**: Unify Golden Demos (817512df-9fed-417c-9e84-f829d0c33bae)
- **Target GitHub Org**: cb-golden-demo-org
- **Project Name**: insurancestack-stu-test-3
- **Environment**: insurancestack-stu-test-3-dev

**Resources Created:**
- ✅ GitHub Repo: https://github.com/cb-golden-demo-org/insurancestack-stu-test-3
- ✅ CloudBees Component: insurancestack-stu-test-3
- ✅ CloudBees Environment: insurancestack-stu-test-3-dev (with FM_TOKEN)
- ✅ CloudBees Application: insurancestack-stu-test-3-app

**Issues Resolved During Testing:**
1. **Template Repository**: CB-InsuranceStack/InsuranceStack needed to be marked as GitHub template
2. **Application Creation Conflict**: Removed `repository:` field from applications section (was trying to auto-create duplicate component)
3. **GitHub Org Configuration**: Golden demos org only has `cb-golden-demo-org` configured, not `CB-InsuranceStack`

**Next Steps:**
- Manually trigger `deploy-simple` workflow in the created repo
- Application will deploy to: `https://insurancestack.se-main-demo.sa-demo.beescloud.com/cb-golden-demo-org/insurancestack-stu-test-3-dev`

### Phase 3: Documentation
- [ ] Update Confluence pages with Mimic instructions
- [ ] Create demo scripts and talking points
- [ ] Record demo videos

## Key Findings from Logan Call (2025-12-23)

### Transcript Analysis
Reviewed conversation with Logan (Mimic creator) covering:

1. **Secret Resolution Hierarchy** ✅
   - CloudBees Unify resolves org-level properties even when Helm expects env-level
   - Confirmed: KUBECONFIG at org level works despite Helm chart expecting environment-level
   - Major simplification for golden demos setup

2. **Required Properties/Secrets Feature**
   - Logan built this specifically to eliminate manual post-setup steps
   - Mimic checks for required properties/secrets and prompts user to create if missing
   - All created at organization level for reuse
   - Syntax: `required_properties: [list]` and `required_secrets: [list]`

3. **Naming Convention Change**
   - FM_KEY should be renamed to FM_TOKEN (uppercase) for consistency
   - Update both scenarios and documentation

4. **Local Development Workflow**
   - `local-dev` scenario pack already configured pointing to local mimic/scenarios directory
   - Can test changes without pushing to GitHub
   - Faster iteration cycle

### Current AccountStack Scenario Status
- Location: `/Users/brown/Git_orgs/cb-demos/mimic/scenarios/accountstack-demo.yaml`
- Status: **Incomplete** - still has manual post-setup instructions in details
- Missing: Use of `required_secrets` for KUBECONFIG
- Missing: Automated FM_TOKEN setup (has `create_fm_token_var: true` but mentions manual GitHub secret)

## Recommendations for InsuranceStack Scenario

### 1. Use Required Secrets (Automate Setup)
```yaml
required_secrets:
  - KUBECONFIG  # Org-level, works despite Helm expecting env-level

# Optional: required_properties if needed
required_properties:
  - DOCKER_USERNAME
  - DOCKER_REGISTRY
```

### 2. Scenario Structure
Based on InsuranceStack architecture (6 components vs AccountStack's 4):
- **Repository**: CB-InsuranceStack/InsuranceStack
- **Components**: 6 services (insurance-ui, policy-service, claims-service, pricing-engine, customer-service, payments-service)
- **Replacements**: Similar pattern to AccountStack
  - "insurancestack" → "${project_name}"
  - "insurance-stack-dev" → "${environment}"
  - Base path replacements in Helm values

### 3. Key Differences from AccountStack
| Feature | AccountStack | InsuranceStack |
|---------|--------------|----------------|
| Services | 4 (web + 3 APIs) | 6 (web + 5 APIs) |
| Feature Flags | 9 flags | 7 flags |
| Domain | Financial | Insurance |
| Database | In-memory | In-memory |

### 4. Files to Modify
```yaml
files_to_modify:
  - README.md
  - helm/insurancestack/values.yaml
  - docker-compose.yaml
  - apps/insurance-ui/README.md  # Flag documentation
```

### 5. Naming Convention Updates
- Change `FM_KEY` → `FM_TOKEN` in both scenarios
- Update documentation to reflect uppercase convention
- Align with CloudBees golden demos standards

## Open Questions

1. Should we update AccountStack scenario first to use `required_secrets` as a template?
2. ECR vs Docker Hub - do we need ECR integration for ASPM binary scanner?
3. Should DOCKER_USERNAME/REGISTRY be org-level properties?

## Related Resources

- [InsuranceStack Confluence](https://cloudbees.atlassian.net/wiki/spaces/UDEMO/pages/5809700869)
- [AccountStack Confluence](https://cloudbees.atlassian.net/wiki/spaces/UDEMO/pages/5781618696)
- [Golden Demos Overview](https://cloudbees.atlassian.net/wiki/spaces/UDEMO/pages/5779816519)
- [Mimic GitHub Repository](https://github.com/cb-demos/mimic)
