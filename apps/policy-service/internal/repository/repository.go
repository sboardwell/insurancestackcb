package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/policy-service/internal/models"
	"github.com/sirupsen/logrus"
)

// Repository provides data access for policies
type Repository struct {
	policies    map[string]*models.Policy
	mu          sync.RWMutex
	logger      *logrus.Logger
	nextID      int
}

// NewRepository creates a new repository and loads data from JSON files
func NewRepository(dataPath string, logger *logrus.Logger) (*Repository, error) {
	repo := &Repository{
		policies: make(map[string]*models.Policy),
		logger:   logger,
		nextID:   1,
	}

	// Load policies
	policiesPath := filepath.Join(dataPath, "policies.json")
	if err := repo.loadPolicies(policiesPath); err != nil {
		// If policies.json doesn't exist, create sample data
		logger.Warnf("Failed to load policies from %s, initializing with sample data: %v", policiesPath, err)
		repo.initializeSamplePolicies()
	}

	logger.Infof("Loaded %d policies from %s", len(repo.policies), dataPath)

	return repo, nil
}

// loadPolicies loads policies from a JSON file
func (r *Repository) loadPolicies(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var policies []*models.Policy
	if err := json.Unmarshal(data, &policies); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, policy := range policies {
		r.policies[policy.ID] = policy
		// Track highest ID for generating new IDs
		var idNum int
		fmt.Sscanf(policy.ID, "pol-%d", &idNum)
		if idNum >= r.nextID {
			r.nextID = idNum + 1
		}
	}

	return nil
}

// initializeSamplePolicies creates sample policies for demo purposes
func (r *Repository) initializeSamplePolicies() {
	samplePolicies := []*models.Policy{
		{
			ID:           "pol-001",
			CustomerID:   "customer-001",
			PolicyNumber: "AUTO-2024-001234",
			Type:         "auto",
			Status:       "active",
			Premium:      1250.00,
			Currency:     "USD",
			StartDate:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:    time.Date(2023, 12, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:           "pol-002",
			CustomerID:   "customer-001",
			PolicyNumber: "HOME-2024-005678",
			Type:         "home",
			Status:       "active",
			Premium:      2100.00,
			Currency:     "USD",
			StartDate:    time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:    time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:    time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:           "pol-003",
			CustomerID:   "customer-002",
			PolicyNumber: "LIFE-2023-009012",
			Type:         "life",
			Status:       "active",
			Premium:      850.00,
			Currency:     "USD",
			StartDate:    time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2033, 6, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:    time.Date(2023, 5, 10, 10, 0, 0, 0, time.UTC),
			UpdatedAt:    time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, policy := range samplePolicies {
		r.policies[policy.ID] = policy
	}
	r.nextID = 4
}

// GetPolicyByID retrieves a policy by ID
func (r *Repository) GetPolicyByID(policyID string) (*models.Policy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	policy, exists := r.policies[policyID]
	if !exists {
		return nil, fmt.Errorf("policy not found")
	}

	return policy, nil
}

// GetPoliciesByCustomerID retrieves all policies for a specific customer
func (r *Repository) GetPoliciesByCustomerID(customerID string) ([]*models.Policy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var customerPolicies []*models.Policy
	for _, policy := range r.policies {
		if policy.CustomerID == customerID {
			customerPolicies = append(customerPolicies, policy)
		}
	}

	return customerPolicies, nil
}

// CreatePolicy creates a new policy
func (r *Repository) CreatePolicy(req models.CreatePolicyRequest) (*models.Policy, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	policy := &models.Policy{
		ID:           fmt.Sprintf("pol-%03d", r.nextID),
		CustomerID:   req.CustomerID,
		PolicyNumber: req.PolicyNumber,
		Type:         req.Type,
		Status:       "active",
		Premium:      req.Premium,
		Coverage:     req.Coverage,
		Deductible:   req.Deductible,
		Currency:     "USD",
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	r.policies[policy.ID] = policy
	r.nextID++

	return policy, nil
}

// UpdatePolicy updates an existing policy
func (r *Repository) UpdatePolicy(policy *models.Policy) (*models.Policy, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.policies[policy.ID]; !exists {
		return nil, fmt.Errorf("policy not found")
	}

	r.policies[policy.ID] = policy
	return policy, nil
}

// GetAllPolicies returns all policies (for testing/admin purposes)
func (r *Repository) GetAllPolicies() []*models.Policy {
	r.mu.RLock()
	defer r.mu.RUnlock()

	policies := make([]*models.Policy, 0, len(r.policies))
	for _, policy := range r.policies {
		policies = append(policies, policy)
	}

	return policies
}
