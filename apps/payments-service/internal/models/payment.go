package models

import "time"

// PaymentType represents the type of payment
type PaymentType string

const (
	PaymentTypePremium PaymentType = "premium"
	PaymentTypePayout  PaymentType = "payout"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

// Payment represents a payment or payout in the insurance system
type Payment struct {
	ID            string        `json:"id"`
	Type          PaymentType   `json:"type"`           // premium or payout
	PolicyID      string        `json:"policyId,omitempty"`      // For premium payments
	ClaimID       string        `json:"claimId,omitempty"`       // For claim payouts
	CustomerID    string        `json:"customerId"`
	Amount        float64       `json:"amount"`
	Status        PaymentStatus `json:"status"`
	ProcessedDate *time.Time    `json:"processedDate,omitempty"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

// CreatePaymentRequest represents a request to create a premium payment
type CreatePaymentRequest struct {
	PolicyID   string  `json:"policyId"`
	CustomerID string  `json:"customerId"`
	Amount     float64 `json:"amount"`
}

// CreatePayoutRequest represents a request to create a claim payout
type CreatePayoutRequest struct {
	ClaimID    string  `json:"claimId"`
	CustomerID string  `json:"customerId"`
	Amount     float64 `json:"amount"`
}

// Validate validates a CreatePaymentRequest
func (r *CreatePaymentRequest) Validate() error {
	if r.PolicyID == "" {
		return &ValidationError{Field: "policyId", Message: "policy ID is required"}
	}
	if r.CustomerID == "" {
		return &ValidationError{Field: "customerId", Message: "customer ID is required"}
	}
	if r.Amount <= 0 {
		return &ValidationError{Field: "amount", Message: "amount must be greater than 0"}
	}
	return nil
}

// Validate validates a CreatePayoutRequest
func (r *CreatePayoutRequest) Validate() error {
	if r.ClaimID == "" {
		return &ValidationError{Field: "claimId", Message: "claim ID is required"}
	}
	if r.CustomerID == "" {
		return &ValidationError{Field: "customerId", Message: "customer ID is required"}
	}
	if r.Amount <= 0 {
		return &ValidationError{Field: "amount", Message: "amount must be greater than 0"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
