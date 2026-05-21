package models

import "time"

// QuoteRequest represents a request for an insurance quote
type QuoteRequest struct {
	PolicyType     string  `json:"policyType" validate:"required,oneof=auto home life"`
	CoverageAmount int     `json:"coverageAmount" validate:"required,min=1"`
	CustomerAge    int     `json:"customerAge" validate:"required,min=18,max=120"`
	RiskScore      int     `json:"riskScore" validate:"required,min=1,max=5"`
	CustomerID     string  `json:"customerId,omitempty"`
	MultiPolicy    bool    `json:"multiPolicy,omitempty"`
	LoyaltyYears   int     `json:"loyaltyYears,omitempty"`
	PaperlessBill  bool    `json:"paperlessBill,omitempty"`
	ClaimsHistory  int     `json:"claimsHistory,omitempty"`
}

// Quote represents an insurance quote response
type Quote struct {
	QuoteID        string    `json:"quoteId"`
	PolicyType     string    `json:"policyType"`
	CoverageAmount int       `json:"coverageAmount"`
	BaseRate       float64   `json:"baseRate"`
	AdjustedRate   float64   `json:"adjustedRate"`
	Discount       float64   `json:"discount"`
	FinalPremium   float64   `json:"finalPremium"`
	ValidUntil     time.Time `json:"validUntil"`
	CreatedAt      time.Time `json:"createdAt"`
	Factors        *Factors  `json:"factors,omitempty"`
}

// Factors represents the breakdown of pricing factors
type Factors struct {
	BaseMultiplier     float64 `json:"baseMultiplier"`
	CoverageMultiplier float64 `json:"coverageMultiplier"`
	AgeMultiplier      float64 `json:"ageMultiplier"`
	RiskMultiplier     float64 `json:"riskMultiplier"`
	DynamicMultiplier  float64 `json:"dynamicMultiplier,omitempty"`
	DiscountAmount     float64 `json:"discountAmount"`
}

// Rate represents base rates for a policy type
type Rate struct {
	PolicyType string             `json:"policyType"`
	BaseRate   float64            `json:"baseRate"`
	Coverage   map[string]float64 `json:"coverage"`
}

// PricingRules represents the complete pricing rules structure
type PricingRules struct {
	BaseRates      map[string]PolicyRates `json:"baseRates"`
	Discounts      Discounts              `json:"discounts"`
	DynamicPricing DynamicPricing         `json:"dynamicPricing"`
	Metadata       Metadata               `json:"metadata"`
}

// PolicyRates represents rates for a specific policy type
type PolicyRates struct {
	Base           float64            `json:"base"`
	Coverage       map[string]float64 `json:"coverage"`
	AgeMultiplier  map[string]float64 `json:"ageMultiplier"`
	RiskMultiplier map[string]float64 `json:"riskMultiplier"`
}

// Discounts represents available discounts
type Discounts struct {
	MultiPolicy     float64            `json:"multiPolicy"`
	LoyaltyYears    map[string]float64 `json:"loyaltyYears"`
	LowRisk         float64            `json:"lowRisk"`
	PaperlessBilling float64           `json:"paperlessBilling"`
}

// DynamicPricing represents dynamic pricing configuration
type DynamicPricing struct {
	Enabled bool            `json:"enabled"`
	Factors DynamicFactors  `json:"factors"`
}

// DynamicFactors represents dynamic pricing factors
type DynamicFactors struct {
	Seasonality      map[string]float64 `json:"seasonality"`
	MarketConditions float64            `json:"marketConditions"`
	ClaimsHistory    map[string]float64 `json:"claimsHistory"`
}

// Metadata represents pricing rules metadata
type Metadata struct {
	LastUpdated   string `json:"lastUpdated"`
	Version       string `json:"version"`
	EffectiveDate string `json:"effectiveDate"`
}

// RatesResponse represents the response for GET /rates
type RatesResponse struct {
	Rates     []Rate    `json:"rates"`
	Timestamp time.Time `json:"timestamp"`
}
