package models

import "time"

// Claim represents an insurance claim
type Claim struct {
	ID            string     `json:"id"`
	PolicyID      string     `json:"policyId"`
	CustomerID    string     `json:"customerId"`
	ClaimNumber   string     `json:"claimNumber"`
	Type          string     `json:"type"`          // accident, theft, damage
	Status        string     `json:"status"`        // submitted, under_review, approved, rejected
	Amount        float64    `json:"amount"`
	Description   string     `json:"description"`
	SubmittedDate time.Time  `json:"submittedDate"`
	ReviewedDate  *time.Time `json:"reviewedDate"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// ClaimFilters represents filters for claim queries
type ClaimFilters struct {
	PolicyID   string
	CustomerID string
	Status     string
	Type       string
}

// Matches checks if a claim matches the given filters
func (c *Claim) Matches(filters *ClaimFilters) bool {
	// Policy ID filter
	if filters.PolicyID != "" && c.PolicyID != filters.PolicyID {
		return false
	}

	// Customer ID filter
	if filters.CustomerID != "" && c.CustomerID != filters.CustomerID {
		return false
	}

	// Status filter
	if filters.Status != "" && c.Status != filters.Status {
		return false
	}

	// Type filter
	if filters.Type != "" && c.Type != filters.Type {
		return false
	}

	return true
}

// CreateClaimRequest represents a request to create a new claim
type CreateClaimRequest struct {
	PolicyID    string  `json:"policyId"`
	CustomerID  string  `json:"customerId"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// UpdateClaimRequest represents a request to update a claim
type UpdateClaimRequest struct {
	Amount      *float64 `json:"amount,omitempty"`
	Description *string  `json:"description,omitempty"`
}

// UpdateClaimStatusRequest represents a request to update claim status
type UpdateClaimStatusRequest struct {
	Status string `json:"status"`
	Notes  string `json:"notes,omitempty"`
}

// ValidateClaimType checks if the claim type is valid
func ValidateClaimType(claimType string) bool {
	validTypes := map[string]bool{
		"accident": true,
		"theft":    true,
		"damage":   true,
	}
	return validTypes[claimType]
}

// ValidateClaimStatus checks if the claim status is valid
func ValidateClaimStatus(status string) bool {
	validStatuses := map[string]bool{
		"submitted":    true,
		"under_review": true,
		"approved":     true,
		"rejected":     true,
	}
	return validStatuses[status]
}
