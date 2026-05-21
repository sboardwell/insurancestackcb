package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/models"
	"github.com/sirupsen/logrus"
)

// Repository provides data access for pricing rules
type Repository struct {
	pricingRules *models.PricingRules
	mu           sync.RWMutex
	logger       *logrus.Logger
}

// NewRepository creates a new repository and loads pricing rules from JSON file
func NewRepository(dataPath string, logger *logrus.Logger) (*Repository, error) {
	repo := &Repository{
		logger: logger,
	}

	// Load pricing rules
	if err := repo.loadPricingRules(filepath.Join(dataPath, "pricing-rules.json")); err != nil {
		return nil, fmt.Errorf("failed to load pricing rules: %w", err)
	}

	logger.Infof("Loaded pricing rules from %s (version: %s)", dataPath, repo.pricingRules.Metadata.Version)

	return repo, nil
}

// loadPricingRules loads pricing rules from a JSON file
func (r *Repository) loadPricingRules(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var rules models.PricingRules
	if err := json.Unmarshal(data, &rules); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.pricingRules = &rules

	return nil
}

// GetPricingRules returns the pricing rules
func (r *Repository) GetPricingRules() *models.PricingRules {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.pricingRules
}

// GetBaseRateForPolicy returns the base rate for a given policy type
func (r *Repository) GetBaseRateForPolicy(policyType string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.pricingRules == nil {
		return 0, fmt.Errorf("pricing rules not loaded")
	}

	rates, exists := r.pricingRules.BaseRates[policyType]
	if !exists {
		return 0, fmt.Errorf("policy type %s not found", policyType)
	}

	return rates.Base, nil
}

// GetCoverageMultiplier returns the coverage multiplier for a given policy type and coverage amount
func (r *Repository) GetCoverageMultiplier(policyType string, coverageAmount int) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.pricingRules == nil {
		return 0, fmt.Errorf("pricing rules not loaded")
	}

	rates, exists := r.pricingRules.BaseRates[policyType]
	if !exists {
		return 0, fmt.Errorf("policy type %s not found", policyType)
	}

	// Convert coverage amount to string for lookup
	coverageStr := fmt.Sprintf("%d", coverageAmount)
	multiplier, exists := rates.Coverage[coverageStr]
	if !exists {
		// Find the closest coverage amount
		return 1.0, nil // Default multiplier
	}

	return multiplier, nil
}

// GetAgeMultiplier returns the age multiplier for a given policy type and age
func (r *Repository) GetAgeMultiplier(policyType string, age int) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.pricingRules == nil {
		return 0, fmt.Errorf("pricing rules not loaded")
	}

	rates, exists := r.pricingRules.BaseRates[policyType]
	if !exists {
		return 0, fmt.Errorf("policy type %s not found", policyType)
	}

	// Determine age range
	var ageRange string
	switch {
	case age >= 18 && age <= 24:
		ageRange = "18-24"
	case age >= 25 && age <= 34:
		ageRange = "25-34"
	case age >= 35 && age <= 49:
		ageRange = "35-49"
	case age >= 50 && age <= 64:
		ageRange = "50-64"
	case age >= 65:
		ageRange = "65+"
	default:
		return 1.0, nil // Default multiplier
	}

	// Try exact match first
	if multiplier, exists := rates.AgeMultiplier[ageRange]; exists {
		return multiplier, nil
	}

	// Try alternate range format for home insurance
	if policyType == "home" {
		altRanges := map[string]string{
			"18-24": "18-34",
			"25-34": "18-34",
		}
		if altRange, exists := altRanges[ageRange]; exists {
			if multiplier, exists := rates.AgeMultiplier[altRange]; exists {
				return multiplier, nil
			}
		}
	}

	return 1.0, nil // Default multiplier
}

// GetRiskMultiplier returns the risk multiplier for a given policy type and risk score
func (r *Repository) GetRiskMultiplier(policyType string, riskScore int) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.pricingRules == nil {
		return 0, fmt.Errorf("pricing rules not loaded")
	}

	rates, exists := r.pricingRules.BaseRates[policyType]
	if !exists {
		return 0, fmt.Errorf("policy type %s not found", policyType)
	}

	riskStr := fmt.Sprintf("%d", riskScore)
	multiplier, exists := rates.RiskMultiplier[riskStr]
	if !exists {
		return 1.0, nil // Default multiplier
	}

	return multiplier, nil
}

// GetDiscounts returns the discounts configuration
func (r *Repository) GetDiscounts() *models.Discounts {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.pricingRules == nil {
		return nil
	}

	return &r.pricingRules.Discounts
}

// GetDynamicPricing returns the dynamic pricing configuration
func (r *Repository) GetDynamicPricing() *models.DynamicPricing {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.pricingRules == nil {
		return nil
	}

	return &r.pricingRules.DynamicPricing
}

// GetAllRates returns all base rates (for GET /rates endpoint)
func (r *Repository) GetAllRates() []models.Rate {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.pricingRules == nil {
		return []models.Rate{}
	}

	rates := make([]models.Rate, 0, len(r.pricingRules.BaseRates))
	for policyType, policyRates := range r.pricingRules.BaseRates {
		rates = append(rates, models.Rate{
			PolicyType: policyType,
			BaseRate:   policyRates.Base,
			Coverage:   policyRates.Coverage,
		})
	}

	return rates
}
