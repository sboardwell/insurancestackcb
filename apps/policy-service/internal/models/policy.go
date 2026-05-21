package models

import "time"

// Policy represents an insurance policy in the system
type Policy struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customerId"`
	PolicyNumber string    `json:"policyNumber"`
	Type         string    `json:"type"`   // auto, home, life
	Status       string    `json:"status"` // active, lapsed, cancelled
	Premium      float64   `json:"premium"`
	Coverage     float64   `json:"coverage"`
	Deductible   float64   `json:"deductible"`
	Currency     string    `json:"currency"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	RenewalDate  time.Time `json:"renewalDate,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// PolicyResponse represents a policy in API responses with optional masking
type PolicyResponse struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customerId"`
	PolicyNumber string    `json:"policyNumber"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	Premium      any       `json:"premium"`   // Can be float64 or string (masked)
	Coverage     any       `json:"coverage"`  // Can be float64 or string (masked)
	Deductible   float64   `json:"deductible,omitempty"`
	Currency     string    `json:"currency"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	RenewalDate  time.Time `json:"renewalDate,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ToResponse converts a Policy to PolicyResponse with optional masking and currency override
func (p *Policy) ToResponse(maskAmounts bool, currency string) PolicyResponse {
	resp := PolicyResponse{
		ID:           p.ID,
		CustomerID:   p.CustomerID,
		PolicyNumber: p.PolicyNumber,
		Type:         p.Type,
		Status:       p.Status,
		Currency:     currency, // Use feature flag currency
		Deductible:   p.Deductible,
		StartDate:    p.StartDate,
		EndDate:      p.EndDate,
		RenewalDate:  p.RenewalDate,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}

	if maskAmounts {
		resp.Premium = "***.**"
		resp.Coverage = "***.**"
	} else {
		resp.Premium = p.Premium
		resp.Coverage = p.Coverage
	}

	return resp
}

// CreatePolicyRequest represents the request body for creating a new policy
type CreatePolicyRequest struct {
	CustomerID   string    `json:"customerId"`
	PolicyNumber string    `json:"policyNumber"`
	Type         string    `json:"type"`
	Premium      float64   `json:"premium"`
	Coverage     float64   `json:"coverage"`
	Deductible   float64   `json:"deductible"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
}

// UpdatePolicyRequest represents the request body for updating a policy
type UpdatePolicyRequest struct {
	Status    *string    `json:"status,omitempty"`
	Premium   *float64   `json:"premium,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
}
