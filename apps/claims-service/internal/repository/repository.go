package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/claims-service/internal/models"
	"github.com/sirupsen/logrus"
)

// Policy represents an insurance policy (minimal structure needed for filtering)
type Policy struct {
	ID         string `json:"id"`
	CustomerID string `json:"customerId"`
}

// Repository provides data access for claims
type Repository struct {
	claims   map[string]*models.Claim
	policies map[string]*Policy // policyID -> Policy
	mu       sync.RWMutex
	logger   *logrus.Logger
}

// NewRepository creates a new repository and loads data from JSON files
func NewRepository(dataPath string, logger *logrus.Logger) (*Repository, error) {
	repo := &Repository{
		claims:   make(map[string]*models.Claim),
		policies: make(map[string]*Policy),
		logger:   logger,
	}

	// Load policies first (needed for customer filtering)
	if err := repo.loadPolicies(filepath.Join(dataPath, "policies.json")); err != nil {
		// Log warning but continue - policies file may not exist yet
		logger.Warnf("Failed to load policies: %v", err)
	}

	// Load claims
	if err := repo.loadClaims(filepath.Join(dataPath, "claims.json")); err != nil {
		// Log warning but continue - claims file may not exist yet
		logger.Warnf("Failed to load claims: %v", err)
	}

	logger.Infof("Loaded %d policies and %d claims from %s", len(repo.policies), len(repo.claims), dataPath)

	return repo, nil
}

// loadPolicies loads policies from a JSON file
func (r *Repository) loadPolicies(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var policies []*Policy
	if err := json.Unmarshal(data, &policies); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, pol := range policies {
		r.policies[pol.ID] = pol
	}

	return nil
}

// loadClaims loads claims from a JSON file
func (r *Repository) loadClaims(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var claims []*models.Claim
	if err := json.Unmarshal(data, &claims); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, claim := range claims {
		r.claims[claim.ID] = claim
	}

	return nil
}

// GetClaimByID retrieves a claim by ID
func (r *Repository) GetClaimByID(claimID string) (*models.Claim, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	claim, exists := r.claims[claimID]
	if !exists {
		return nil, fmt.Errorf("claim not found")
	}

	return claim, nil
}

// GetAllClaims returns all claims
func (r *Repository) GetAllClaims() []*models.Claim {
	r.mu.RLock()
	defer r.mu.RUnlock()

	claims := make([]*models.Claim, 0, len(r.claims))
	for _, claim := range r.claims {
		claims = append(claims, claim)
	}

	return claims
}

// GetPolicyIDsByCustomerID retrieves all policy IDs for a given customer
func (r *Repository) GetPolicyIDsByCustomerID(customerID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var policyIDs []string
	for _, pol := range r.policies {
		if pol.CustomerID == customerID {
			policyIDs = append(policyIDs, pol.ID)
		}
	}

	return policyIDs
}

// GetClaimsByFilter retrieves claims matching the given filters
func (r *Repository) GetClaimsByFilter(filters *models.ClaimFilters) []*models.Claim {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*models.Claim
	for _, claim := range r.claims {
		if claim.Matches(filters) {
			filtered = append(filtered, claim)
		}
	}

	return filtered
}

// CreateClaim creates a new claim
func (r *Repository) CreateClaim(claim *models.Claim) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.claims[claim.ID] = claim
	return nil
}

// UpdateClaim updates an existing claim
func (r *Repository) UpdateClaim(claim *models.Claim) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.claims[claim.ID]; !exists {
		return fmt.Errorf("claim not found")
	}

	r.claims[claim.ID] = claim
	return nil
}
