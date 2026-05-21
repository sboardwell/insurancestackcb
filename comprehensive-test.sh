#!/bin/bash

# Comprehensive Test Script for InsuranceStack
# Tests all services, APIs, and validates data integrity

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_TOTAL=0

# Function to print test result
test_result() {
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC}: $2"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: $2"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        if [ -n "$3" ]; then
            echo -e "  ${YELLOW}Details: $3${NC}"
        fi
    fi
}

# Function to test HTTP endpoint
test_endpoint() {
    local url=$1
    local expected_code=$2
    local description=$3

    response=$(curl -s -w "\n%{http_code}" "$url")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "$expected_code" ]; then
        test_result 0 "$description"
        return 0
    else
        test_result 1 "$description" "Expected $expected_code, got $http_code"
        return 1
    fi
}

# Function to test JSON array length
test_array_length() {
    local url=$1
    local min_length=$2
    local description=$3

    response=$(curl -s "$url")
    length=$(echo "$response" | jq 'length' 2>/dev/null || echo "0")

    if [ "$length" -ge "$min_length" ]; then
        test_result 0 "$description (found $length items)"
        return 0
    else
        test_result 1 "$description" "Expected at least $min_length items, found $length"
        return 1
    fi
}

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}InsuranceStack Comprehensive Test Suite${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

# ============================================================
# 1. SERVICE HEALTH CHECKS
# ============================================================
echo -e "${BLUE}1. Service Health Checks${NC}"
echo "----------------------------"

test_endpoint "http://localhost:8001/healthz" "200" "Policy Service health check"
test_endpoint "http://localhost:8002/healthz" "200" "Claims Service health check"
test_endpoint "http://localhost:8003/healthz" "200" "Pricing Engine health check"
test_endpoint "http://localhost:8004/healthz" "200" "Customer Service health check"
test_endpoint "http://localhost:8005/healthz" "200" "Payments Service health check"
test_endpoint "http://localhost:3000" "200" "Frontend UI accessible"

echo ""

# ============================================================
# 2. DATA ENDPOINTS - BASIC RETRIEVAL
# ============================================================
echo -e "${BLUE}2. Data Endpoints - Basic Retrieval${NC}"
echo "----------------------------"

test_array_length "http://localhost:8001/policies" "1" "Policies endpoint returns data"
test_array_length "http://localhost:8002/claims" "1" "Claims endpoint returns data"
test_array_length "http://localhost:8004/customers" "1" "Customers endpoint returns data"
test_array_length "http://localhost:8005/payments" "1" "Payments endpoint returns data"

echo ""

# ============================================================
# 3. SPECIFIC DATA VALIDATION
# ============================================================
echo -e "${BLUE}3. Data Validation${NC}"
echo "----------------------------"

# Test specific policy retrieval
response=$(curl -s "http://localhost:8001/policies")
if echo "$response" | jq -e '.[] | select(.id == "pol-001")' > /dev/null 2>&1; then
    test_result 0 "Policy pol-001 exists for default customer"
else
    test_result 1 "Policy pol-001 exists for default customer" "Policy not found"
fi

# Test specific claim retrieval
response=$(curl -s "http://localhost:8002/claims")
if echo "$response" | jq -e '.[] | select(.id == "claim-001")' > /dev/null 2>&1; then
    test_result 0 "Claim claim-001 exists"
else
    test_result 1 "Claim claim-001 exists" "Claim not found"
fi

# Test specific customer retrieval
response=$(curl -s "http://localhost:8004/customers")
if echo "$response" | jq -e '.[] | select(.id == "cust-001")' > /dev/null 2>&1; then
    test_result 0 "Customer cust-001 exists"
else
    test_result 1 "Customer cust-001 exists" "Customer not found"
fi

# Test specific payment retrieval
response=$(curl -s "http://localhost:8005/payments")
if echo "$response" | jq -e '.[] | select(.type == "premium")' > /dev/null 2>&1; then
    test_result 0 "Premium payments exist"
else
    test_result 1 "Premium payments exist" "No premium payments found"
fi

if echo "$response" | jq -e '.[] | select(.type == "payout")' > /dev/null 2>&1; then
    test_result 0 "Payout payments exist"
else
    test_result 1 "Payout payments exist" "No payout payments found"
fi

echo ""

# ============================================================
# 4. PRICING ENGINE TESTS
# ============================================================
echo -e "${BLUE}4. Pricing Engine${NC}"
echo "----------------------------"

# Test base rates endpoint
test_endpoint "http://localhost:8003/rates" "200" "Pricing engine rates endpoint"

# Test quote calculation
quote_response=$(curl -s -X POST http://localhost:8003/quote \
  -H "Content-Type: application/json" \
  -d '{"policyType":"auto","coverage":500000,"deductible":1000,"age":30,"zipCode":"94102"}')

if echo "$quote_response" | jq -e '.premium' > /dev/null 2>&1; then
    test_result 0 "Quote calculation returns premium"
else
    test_result 1 "Quote calculation returns premium" "No premium in response"
fi

echo ""

# ============================================================
# 5. AUTHENTICATION TESTS
# ============================================================
echo -e "${BLUE}5. Authentication & Authorization${NC}"
echo "----------------------------"

# Test that health endpoints work without auth
test_endpoint "http://localhost:8001/healthz" "200" "Health endpoint accessible without auth"

# Test that data endpoints work with default customer
test_endpoint "http://localhost:8001/policies" "200" "Policies accessible with default auth"
test_endpoint "http://localhost:8002/claims" "200" "Claims accessible with default auth"
test_endpoint "http://localhost:8004/customers" "200" "Customers accessible with default auth"
test_endpoint "http://localhost:8005/payments" "200" "Payments accessible with default auth"

echo ""

# ============================================================
# 6. DATA INTEGRITY CHECKS
# ============================================================
echo -e "${BLUE}6. Data Integrity${NC}"
echo "----------------------------"

# Check that customer emails are InsuranceStack
customers_response=$(curl -s "http://localhost:8004/customers")
accountstack_count=$(echo "$customers_response" | jq '[.[] | select(.email | contains("accountstack"))] | length' 2>/dev/null || echo "0")

if [ "$accountstack_count" -eq "0" ]; then
    test_result 0 "No AccountStack references in customer emails"
else
    test_result 1 "No AccountStack references in customer emails" "Found $accountstack_count AccountStack emails"
fi

# Check that all policies have required fields
policies_response=$(curl -s "http://localhost:8001/policies")
policies_with_required=$(echo "$policies_response" | jq '[.[] | select(.id and .policyNumber and .customerId and .type and .status)] | length' 2>/dev/null || echo "0")
total_policies=$(echo "$policies_response" | jq 'length' 2>/dev/null || echo "0")

if [ "$policies_with_required" -eq "$total_policies" ]; then
    test_result 0 "All policies have required fields"
else
    test_result 1 "All policies have required fields" "Only $policies_with_required of $total_policies have required fields"
fi

# Check that all claims have required fields
claims_response=$(curl -s "http://localhost:8002/claims")
claims_with_required=$(echo "$claims_response" | jq '[.[] | select(.id and .policyId and .customerId and .status)] | length' 2>/dev/null || echo "0")
total_claims=$(echo "$claims_response" | jq 'length' 2>/dev/null || echo "0")

if [ "$claims_with_required" -eq "$total_claims" ]; then
    test_result 0 "All claims have required fields"
else
    test_result 1 "All claims have required fields" "Only $claims_with_required of $total_claims have required fields"
fi

echo ""

# ============================================================
# 7. FRONTEND UI TESTS
# ============================================================
echo -e "${BLUE}7. Frontend UI${NC}"
echo "----------------------------"

# Test that main pages are accessible
test_endpoint "http://localhost:3000/" "200" "Homepage accessible"
test_endpoint "http://localhost:3000/login" "200" "Login page accessible"
test_endpoint "http://localhost:3000/policies" "200" "Policies page accessible"
test_endpoint "http://localhost:3000/claims" "200" "Claims page accessible"
test_endpoint "http://localhost:3000/customers" "200" "Customers page accessible"
test_endpoint "http://localhost:3000/payments" "200" "Payments page accessible"
test_endpoint "http://localhost:3000/get-quote" "200" "Get Quote page accessible"

# Test that frontend API proxying works
test_endpoint "http://localhost:3000/api/policies" "200" "Frontend proxies policies API"
test_endpoint "http://localhost:3000/api/claims" "200" "Frontend proxies claims API"
test_endpoint "http://localhost:3000/api/customers" "200" "Frontend proxies customers API"
test_endpoint "http://localhost:3000/api/payments" "200" "Frontend proxies payments API"

echo ""

# ============================================================
# 8. DOCKER CONTAINER STATUS
# ============================================================
echo -e "${BLUE}8. Docker Container Status${NC}"
echo "----------------------------"

containers=("insurancestack-policy-service" "insurancestack-claims-service" "insurancestack-pricing-engine" "insurancestack-customer-service" "insurancestack-payments-service" "insurancestack-ui")

for container in "${containers[@]}"; do
    if docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
        status=$(docker inspect -f '{{.State.Health.Status}}' "$container" 2>/dev/null || echo "running")
        if [ "$status" = "healthy" ] || [ "$status" = "running" ]; then
            test_result 0 "Container $container is running"
        else
            test_result 1 "Container $container is healthy" "Status: $status"
        fi
    else
        test_result 1 "Container $container is running" "Container not found"
    fi
done

echo ""

# ============================================================
# SUMMARY
# ============================================================
echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""
echo -e "Total Tests: ${BLUE}$TESTS_TOTAL${NC}"
echo -e "Passed:      ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed:      ${RED}$TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi
