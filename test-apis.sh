#!/bin/bash

echo "================================================"
echo "InsuranceStack API Test Suite"
echo "================================================"
echo ""

# Colors
GREEN='\033[0.;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test health endpoints
echo "1. Testing Health Endpoints"
echo "----------------------------"
for port in 8001 8002 8003 8004 8005; do
    response=$(curl -s -w "\n%{http_code}" http://localhost:$port/healthz)
    status=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$status" = "200" ]; then
        echo -e "${GREEN}✓${NC} Port $port: $(echo $body | jq -r '.service // .status')"
    else
        echo -e "${RED}✗${NC} Port $port: Failed (HTTP $status)"
    fi
done

echo ""
echo "2. Testing Data Endpoints"
echo "----------------------------"

# Test policies
policies=$(curl -s http://localhost:8001/policies)
count=$(echo "$policies" | jq '. | length')
echo -e "${YELLOW}Policies:${NC} $count items returned"
if [ "$count" -gt "0" ]; then
    echo "$policies" | jq -r '.[0:2] | .[] | "  - \(.id): \(.policyType) (\(.status))"'
fi

# Test claims
claims=$(curl -s http://localhost:8002/claims)
count=$(echo "$claims" | jq '. | length')
echo -e "${YELLOW}Claims:${NC} $count items returned"
if [ "$count" -gt "0" ]; then
    echo "$claims" | jq -r '.[0:2] | .[] | "  - \(.id): \(.claimType) - \(.status)"'
fi

# Test customers
customers=$(curl -s http://localhost:8004/customers)
count=$(echo "$customers" | jq '. | length')
echo -e "${YELLOW}Customers:${NC} $count items returned"
if [ "$count" -gt "0" ]; then
    echo "$customers" | jq -r '.[0:2] | .[] | "  - \(.id): \(.name) (Risk: \(.riskScore))"'
fi

# Test payments
payments=$(curl -s http://localhost:8005/payments)
count=$(echo "$payments" | jq '. | length')
echo -e "${YELLOW}Payments:${NC} $count items returned"
if [ "$count" -gt "0" ]; then
    echo "$payments" | jq -r '.[0:2] | .[] | "  - \(.id): \(.type) - $\(.amount)"'
fi

echo ""
echo "3. Testing Pricing Engine"
echo "----------------------------"

# Test get rates
rates=$(curl -s http://localhost:8003/rates)
echo -e "${YELLOW}Base Rates:${NC}"
echo "$rates" | jq -r 'to_entries[] | "  - \(.key): $\(.value)"'

# Test quote calculation
quote=$(curl -s -X POST http://localhost:8003/quote \
    -H "Content-Type: application/json" \
    -d '{"policyType":"auto","coverageAmount":50000,"customerId":"cust-001"}')

if [ $? -eq 0 ]; then
    premium=$(echo "$quote" | jq -r '.estimatedPremium // "N/A"')
    echo -e "${GREEN}✓${NC} Quote generated: \$${premium}/year"
else
    echo -e "${RED}✗${NC} Quote generation failed"
fi

echo ""
echo "================================================"
echo "Test Complete!"
echo "================================================"
