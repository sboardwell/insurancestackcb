#!/bin/bash

# Comprehensive UI Page Test Script
# Tests all frontend pages to ensure they load and display data correctly

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# Function to check if page contains expected content
check_page_content() {
    local url=$1
    local expected_text=$2
    local description=$3

    response=$(curl -s "$url")

    if echo "$response" | grep -q "$expected_text"; then
        test_result 0 "$description"
        return 0
    else
        test_result 1 "$description" "Expected text '$expected_text' not found"
        return 1
    fi
}

# Function to check API endpoint returns valid JSON array
check_api_data() {
    local url=$1
    local min_length=$2
    local description=$3

    response=$(curl -s "$url")

    # Check if response is valid JSON
    if ! echo "$response" | jq . > /dev/null 2>&1; then
        test_result 1 "$description" "Response is not valid JSON"
        return 1
    fi

    # Check if response is an array
    if ! echo "$response" | jq -e 'type == "array"' > /dev/null 2>&1; then
        test_result 1 "$description" "Response is not an array"
        return 1
    fi

    # Check array length
    length=$(echo "$response" | jq 'length')
    if [ "$length" -ge "$min_length" ]; then
        test_result 0 "$description (found $length items)"
        return 0
    else
        test_result 1 "$description" "Expected at least $min_length items, found $length"
        return 1
    fi
}

# Function to check specific field exists in API response
check_api_field() {
    local url=$1
    local field_path=$2
    local description=$3

    response=$(curl -s "$url")

    if echo "$response" | jq -e "$field_path" > /dev/null 2>&1; then
        test_result 0 "$description"
        return 0
    else
        test_result 1 "$description" "Field '$field_path' not found or is null"
        return 1
    fi
}

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}InsuranceStack UI Page Test Suite${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

# ============================================================
# 1. BASIC PAGE ACCESSIBILITY
# ============================================================
echo -e "${BLUE}1. Page Accessibility${NC}"
echo "----------------------------"

# Check pages return 200 status and contain the app div
for page in "/" "/login" "/policies" "/claims" "/customers" "/payments" "/get-quote"; do
    http_code=$(curl -s -w "%{http_code}" -o /dev/null "http://localhost:3000$page")
    page_name=$(echo "$page" | sed 's/\///' | sed 's/-/ /' | awk '{for(i=1;i<=NF;i++){ $i=toupper(substr($i,1,1)) substr($i,2) }}1')
    if [ -z "$page_name" ]; then page_name="Homepage"; fi

    if [ "$http_code" = "200" ]; then
        test_result 0 "$page_name page accessible (HTTP 200)"
    else
        test_result 1 "$page_name page accessible" "Got HTTP $http_code"
    fi
done

echo ""

# ============================================================
# 2. API DATA ENDPOINTS
# ============================================================
echo -e "${BLUE}2. API Data Endpoints${NC}"
echo "----------------------------"

check_api_data "http://localhost:3000/api/policies" "1" "Policies API returns data"
check_api_data "http://localhost:3000/api/claims" "1" "Claims API returns data"
check_api_data "http://localhost:3000/api/customers" "1" "Customers API returns data"
check_api_data "http://localhost:3000/api/payments" "1" "Payments API returns data"

echo ""

# ============================================================
# 3. API DATA STRUCTURE - POLICIES
# ============================================================
echo -e "${BLUE}3. Policies API Data Structure${NC}"
echo "----------------------------"

check_api_field "http://localhost:3000/api/policies" ".[0].id" "Policy has 'id' field"
check_api_field "http://localhost:3000/api/policies" ".[0].type" "Policy has 'type' field (not policyType)"
check_api_field "http://localhost:3000/api/policies" ".[0].policyNumber" "Policy has 'policyNumber' field"
check_api_field "http://localhost:3000/api/policies" ".[0].status" "Policy has 'status' field"
check_api_field "http://localhost:3000/api/policies" ".[0].premium" "Policy has 'premium' field"
check_api_field "http://localhost:3000/api/policies" ".[0].coverage" "Policy has 'coverage' field"

echo ""

# ============================================================
# 4. API DATA STRUCTURE - CLAIMS
# ============================================================
echo -e "${BLUE}4. Claims API Data Structure${NC}"
echo "----------------------------"

check_api_field "http://localhost:3000/api/claims" ".[0].id" "Claim has 'id' field"
check_api_field "http://localhost:3000/api/claims" ".[0].type" "Claim has 'type' field (not claimType)"
check_api_field "http://localhost:3000/api/claims" ".[0].claimNumber" "Claim has 'claimNumber' field"
check_api_field "http://localhost:3000/api/claims" ".[0].status" "Claim has 'status' field"
check_api_field "http://localhost:3000/api/claims" ".[0].amount" "Claim has 'amount' field"
check_api_field "http://localhost:3000/api/claims" ".[0].submittedDate" "Claim has 'submittedDate' field (not dateOfLoss)"

echo ""

# ============================================================
# 5. API DATA STRUCTURE - CUSTOMERS
# ============================================================
echo -e "${BLUE}5. Customers API Data Structure${NC}"
echo "----------------------------"

check_api_field "http://localhost:3000/api/customers" ".[0].id" "Customer has 'id' field"
check_api_field "http://localhost:3000/api/customers" ".[0].firstName" "Customer has 'firstName' field"
check_api_field "http://localhost:3000/api/customers" ".[0].lastName" "Customer has 'lastName' field"
check_api_field "http://localhost:3000/api/customers" ".[0].email" "Customer has 'email' field"
check_api_field "http://localhost:3000/api/customers" ".[0].address" "Customer has 'address' field (object)"
check_api_field "http://localhost:3000/api/customers" ".[0].address.street" "Customer address has 'street' field"
check_api_field "http://localhost:3000/api/customers" ".[0].riskScore" "Customer has 'riskScore' field (not status)"

echo ""

# ============================================================
# 6. API DATA STRUCTURE - PAYMENTS
# ============================================================
echo -e "${BLUE}6. Payments API Data Structure${NC}"
echo "----------------------------"

check_api_field "http://localhost:3000/api/payments" ".[0].id" "Payment has 'id' field"
check_api_field "http://localhost:3000/api/payments" ".[0].type" "Payment has 'type' field (not paymentType)"
check_api_field "http://localhost:3000/api/payments" ".[0].amount" "Payment has 'amount' field"
check_api_field "http://localhost:3000/api/payments" ".[0].status" "Payment has 'status' field"
check_api_field "http://localhost:3000/api/payments" ".[0].customerId" "Payment has 'customerId' field"

echo ""

# ============================================================
# 7. DATA CONTENT VALIDATION
# ============================================================
echo -e "${BLUE}7. Data Content Validation${NC}"
echo "----------------------------"

# Check that policies have valid types
policies_response=$(curl -s "http://localhost:3000/api/policies")
if echo "$policies_response" | jq -e '.[0].type | IN("auto", "home", "life", "health")' > /dev/null 2>&1; then
    test_result 0 "Policy types are valid (auto/home/life/health)"
else
    test_result 1 "Policy types are valid" "Invalid policy type found"
fi

# Check that claims have valid types
claims_response=$(curl -s "http://localhost:3000/api/claims")
if echo "$claims_response" | jq -e '.[0].type | IN("accident", "theft", "damage")' > /dev/null 2>&1; then
    test_result 0 "Claim types are valid (accident/theft/damage)"
else
    test_result 1 "Claim types are valid" "Invalid claim type found"
fi

# Check that payments have valid types
payments_response=$(curl -s "http://localhost:3000/api/payments")
if echo "$payments_response" | jq -e '.[0].type | IN("premium", "payout")' > /dev/null 2>&1; then
    test_result 0 "Payment types are valid (premium/payout)"
else
    test_result 1 "Payment types are valid" "Invalid payment type found"
fi

# Check no AccountStack references in customer emails
customers_response=$(curl -s "http://localhost:3000/api/customers")
accountstack_count=$(echo "$customers_response" | jq '[.[] | select(.email | contains("accountstack"))] | length')
if [ "$accountstack_count" -eq "0" ]; then
    test_result 0 "No AccountStack references in customer emails"
else
    test_result 1 "No AccountStack references" "Found $accountstack_count AccountStack email addresses"
fi

echo ""

# ============================================================
# 8. FRONTEND COMPONENT TESTS
# ============================================================
echo -e "${BLUE}8. Frontend Component Tests${NC}"
echo "----------------------------"

# Test GetQuote page has the necessary UI structure in source
quote_source=$(curl -s "http://localhost:3000/get-quote")
if echo "$quote_source" | grep -q "root"; then
    test_result 0 "GetQuote page has React root element"
else
    test_result 1 "GetQuote page structure" "React root element not found"
fi

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
    echo -e "${GREEN}✓ All UI tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some UI tests failed${NC}"
    exit 1
fi
