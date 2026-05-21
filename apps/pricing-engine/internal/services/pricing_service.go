package services

import (
	"fmt"
	"time"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/features"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/models"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PricingService handles pricing calculations
type PricingService struct {
	repo   *repository.Repository
	flags  *features.Flags
	logger *logrus.Logger
}

// NewPricingService creates a new pricing service
func NewPricingService(repo *repository.Repository, flags *features.Flags, logger *logrus.Logger) *PricingService {
	return &PricingService{
		repo:   repo,
		flags:  flags,
		logger: logger,
	}
}

// CalculateQuote calculates an insurance quote based on the request
func (s *PricingService) CalculateQuote(req *models.QuoteRequest) (*models.Quote, error) {
	// Validate request
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	// Get base rate
	baseRate, err := s.repo.GetBaseRateForPolicy(req.PolicyType)
	if err != nil {
		return nil, fmt.Errorf("failed to get base rate: %w", err)
	}

	// Get coverage multiplier
	coverageMultiplier, err := s.repo.GetCoverageMultiplier(req.PolicyType, req.CoverageAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to get coverage multiplier: %w", err)
	}

	// Get age multiplier
	ageMultiplier, err := s.repo.GetAgeMultiplier(req.PolicyType, req.CustomerAge)
	if err != nil {
		return nil, fmt.Errorf("failed to get age multiplier: %w", err)
	}

	// Get risk multiplier
	riskMultiplier, err := s.repo.GetRiskMultiplier(req.PolicyType, req.RiskScore)
	if err != nil {
		return nil, fmt.Errorf("failed to get risk multiplier: %w", err)
	}

	// Calculate base premium
	basePremium := baseRate * coverageMultiplier * ageMultiplier * riskMultiplier

	// Apply dynamic pricing if enabled
	dynamicMultiplier := 1.0
	if s.flags.IsDynamicRatesEnabled() {
		dynamicMultiplier = s.calculateDynamicMultiplier(req)
	}

	adjustedRate := basePremium * dynamicMultiplier

	// Calculate discounts
	discount := s.calculateDiscount(req, adjustedRate)

	// Calculate final premium
	finalPremium := adjustedRate - discount

	// Create quote
	quote := &models.Quote{
		QuoteID:        generateQuoteID(),
		PolicyType:     req.PolicyType,
		CoverageAmount: req.CoverageAmount,
		BaseRate:       basePremium,
		AdjustedRate:   adjustedRate,
		Discount:       discount,
		FinalPremium:   finalPremium,
		ValidUntil:     time.Now().Add(30 * 24 * time.Hour), // Valid for 30 days
		CreatedAt:      time.Now(),
		Factors: &models.Factors{
			BaseMultiplier:     baseRate,
			CoverageMultiplier: coverageMultiplier,
			AgeMultiplier:      ageMultiplier,
			RiskMultiplier:     riskMultiplier,
			DynamicMultiplier:  dynamicMultiplier,
			DiscountAmount:     discount,
		},
	}

	s.logger.WithFields(logrus.Fields{
		"quoteId":      quote.QuoteID,
		"policyType":   req.PolicyType,
		"finalPremium": quote.FinalPremium,
		"dynamicRates": s.flags.IsDynamicRatesEnabled(),
	}).Info("Quote calculated")

	return quote, nil
}

// calculateDynamicMultiplier calculates dynamic pricing adjustments
func (s *PricingService) calculateDynamicMultiplier(req *models.QuoteRequest) float64 {
	dynamicPricing := s.repo.GetDynamicPricing()
	if dynamicPricing == nil || !dynamicPricing.Enabled {
		return 1.0
	}

	multiplier := 1.0

	// Apply seasonality factor
	quarter := getCurrentQuarter()
	if seasonalityFactor, exists := dynamicPricing.Factors.Seasonality[quarter]; exists {
		multiplier *= seasonalityFactor
	}

	// Apply market conditions
	multiplier *= dynamicPricing.Factors.MarketConditions

	// Apply claims history factor
	claimsKey := getClaimsHistoryKey(req.ClaimsHistory)
	if claimsFactor, exists := dynamicPricing.Factors.ClaimsHistory[claimsKey]; exists {
		multiplier *= claimsFactor
	}

	s.logger.WithFields(logrus.Fields{
		"quarter":          quarter,
		"claimsHistory":    req.ClaimsHistory,
		"dynamicMultiplier": multiplier,
	}).Debug("Dynamic multiplier calculated")

	return multiplier
}

// calculateDiscount calculates the total discount based on request parameters
func (s *PricingService) calculateDiscount(req *models.QuoteRequest, adjustedRate float64) float64 {
	discounts := s.repo.GetDiscounts()
	if discounts == nil {
		return 0
	}

	totalDiscount := 0.0

	// Multi-policy discount
	if req.MultiPolicy {
		totalDiscount += adjustedRate * discounts.MultiPolicy
	}

	// Loyalty discount
	if req.LoyaltyYears > 0 {
		loyaltyDiscount := s.getLoyaltyDiscount(req.LoyaltyYears, discounts)
		totalDiscount += adjustedRate * loyaltyDiscount
	}

	// Low risk discount
	if req.RiskScore == 1 {
		totalDiscount += adjustedRate * discounts.LowRisk
	}

	// Paperless billing discount
	if req.PaperlessBill {
		totalDiscount += adjustedRate * discounts.PaperlessBilling
	}

	s.logger.WithFields(logrus.Fields{
		"multiPolicy":   req.MultiPolicy,
		"loyaltyYears":  req.LoyaltyYears,
		"paperlessBill": req.PaperlessBill,
		"riskScore":     req.RiskScore,
		"totalDiscount": totalDiscount,
	}).Debug("Discount calculated")

	return totalDiscount
}

// getLoyaltyDiscount gets the loyalty discount percentage based on years
func (s *PricingService) getLoyaltyDiscount(years int, discounts *models.Discounts) float64 {
	// Find the highest applicable loyalty discount
	applicableDiscount := 0.0

	for yearsStr, discount := range discounts.LoyaltyYears {
		var requiredYears int
		fmt.Sscanf(yearsStr, "%d", &requiredYears)

		if years >= requiredYears && discount > applicableDiscount {
			applicableDiscount = discount
		}
	}

	return applicableDiscount
}

// GetRates returns all available rates
func (s *PricingService) GetRates() *models.RatesResponse {
	rates := s.repo.GetAllRates()

	return &models.RatesResponse{
		Rates:     rates,
		Timestamp: time.Now(),
	}
}

// validateRequest validates the quote request
func (s *PricingService) validateRequest(req *models.QuoteRequest) error {
	if req.PolicyType != "auto" && req.PolicyType != "home" && req.PolicyType != "life" {
		return fmt.Errorf("invalid policy type: %s (must be auto, home, or life)", req.PolicyType)
	}

	if req.CoverageAmount <= 0 {
		return fmt.Errorf("coverage amount must be greater than 0")
	}

	if req.CustomerAge < 18 || req.CustomerAge > 120 {
		return fmt.Errorf("customer age must be between 18 and 120")
	}

	if req.RiskScore < 1 || req.RiskScore > 5 {
		return fmt.Errorf("risk score must be between 1 and 5")
	}

	return nil
}

// Helper functions

func generateQuoteID() string {
	return "Q-" + uuid.New().String()[:8]
}

func getCurrentQuarter() string {
	month := time.Now().Month()
	switch {
	case month >= 1 && month <= 3:
		return "Q1"
	case month >= 4 && month <= 6:
		return "Q2"
	case month >= 7 && month <= 9:
		return "Q3"
	default:
		return "Q4"
	}
}

func getClaimsHistoryKey(claims int) string {
	switch {
	case claims == 0:
		return "0"
	case claims == 1:
		return "1"
	case claims == 2:
		return "2"
	default:
		return "3+"
	}
}
