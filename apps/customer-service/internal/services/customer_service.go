package services

import (
	"fmt"
	"time"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/customer-service/internal/features"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/customer-service/internal/models"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/customer-service/internal/repository"
	"github.com/sirupsen/logrus"
)

// CustomerService handles business logic for customers
type CustomerService struct {
	repo   *repository.Repository
	flags  *features.Flags
	logger *logrus.Logger
}

// NewCustomerService creates a new customer service
func NewCustomerService(repo *repository.Repository, flags *features.Flags, logger *logrus.Logger) *CustomerService {
	return &CustomerService{
		repo:   repo,
		flags:  flags,
		logger: logger,
	}
}

// GetAllCustomers retrieves all customers
func (s *CustomerService) GetAllCustomers() ([]*models.Customer, error) {
	customers, err := s.repo.GetAllCustomers()
	if err != nil {
		s.logger.Error("Failed to retrieve customers")
		return nil, err
	}

	s.logger.WithField("count", len(customers)).Debug("Retrieved customers")
	return customers, nil
}

// GetCustomerByID retrieves a customer by ID
func (s *CustomerService) GetCustomerByID(customerID string) (*models.Customer, error) {
	customer, err := s.repo.GetCustomerByID(customerID)
	if err != nil {
		s.logger.WithField("customerId", customerID).Warn("Customer not found")
		return nil, err
	}

	s.logger.WithField("customerId", customerID).Debug("Customer retrieved")
	return customer, nil
}

// CreateCustomer creates a new customer
func (s *CustomerService) CreateCustomer(req models.CreateCustomerRequest) (*models.Customer, error) {
	// Generate a new customer ID (in production, this would use UUID or database auto-increment)
	customerID := fmt.Sprintf("cust-%d", time.Now().UnixNano())

	// Default risk score for new customers
	defaultRiskScore := 50

	customer := &models.Customer{
		ID:          customerID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address,
		DateOfBirth: req.DateOfBirth,
		RiskScore:   defaultRiskScore,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateCustomer(customer); err != nil {
		s.logger.WithError(err).Error("Failed to create customer")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"customerId": customerID,
		"email":      req.Email,
	}).Info("Customer created")

	return customer, nil
}

// UpdateCustomer updates an existing customer
func (s *CustomerService) UpdateCustomer(customerID string, req models.UpdateCustomerRequest) (*models.Customer, error) {
	// First, check if customer exists
	existingCustomer, err := s.repo.GetCustomerByID(customerID)
	if err != nil {
		s.logger.WithField("customerId", customerID).Warn("Customer not found for update")
		return nil, err
	}

	// Update the customer fields
	existingCustomer.FirstName = req.FirstName
	existingCustomer.LastName = req.LastName
	existingCustomer.Email = req.Email
	existingCustomer.Phone = req.Phone
	existingCustomer.Address = req.Address
	existingCustomer.DateOfBirth = req.DateOfBirth
	existingCustomer.UpdatedAt = time.Now()

	if err := s.repo.UpdateCustomer(existingCustomer); err != nil {
		s.logger.WithError(err).WithField("customerId", customerID).Error("Failed to update customer")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"customerId": customerID,
		"email":      req.Email,
	}).Info("Customer updated")

	return existingCustomer, nil
}

// DeactivateCustomer deactivates a customer (soft delete)
func (s *CustomerService) DeactivateCustomer(customerID string) error {
	// First, check if customer exists
	_, err := s.repo.GetCustomerByID(customerID)
	if err != nil {
		s.logger.WithField("customerId", customerID).Warn("Customer not found for deactivation")
		return err
	}

	if err := s.repo.DeactivateCustomer(customerID); err != nil {
		s.logger.WithError(err).WithField("customerId", customerID).Error("Failed to deactivate customer")
		return err
	}

	s.logger.WithField("customerId", customerID).Info("Customer deactivated")
	return nil
}
