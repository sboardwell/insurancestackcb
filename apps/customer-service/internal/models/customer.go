package models

import "time"

// Address represents a customer's address
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
	Country string `json:"country"`
}

// Customer represents an insurance customer in the system
type Customer struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     Address   `json:"address"`
	DateOfBirth string    `json:"dateOfBirth"` // ISO 8601 date format (YYYY-MM-DD)
	RiskScore   int       `json:"riskScore"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateCustomerRequest represents the request body for creating a customer
type CreateCustomerRequest struct {
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Address     Address `json:"address"`
	DateOfBirth string  `json:"dateOfBirth"` // ISO 8601 date format (YYYY-MM-DD)
}

// UpdateCustomerRequest represents the request body for updating a customer
type UpdateCustomerRequest struct {
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Address     Address `json:"address"`
	DateOfBirth string  `json:"dateOfBirth"` // ISO 8601 date format (YYYY-MM-DD)
}
