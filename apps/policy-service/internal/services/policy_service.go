package services

import (
	"fmt"
	"time"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/policy-service/internal/features"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/policy-service/internal/models"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/policy-service/internal/repository"
	"github.com/sirupsen/logrus"
)

// PolicyService handles business logic for policies
type PolicyService struct {
	repo   *repository.Repository
	flags  *features.Flags
	logger *logrus.Logger
}

// NewPolicyService creates a new policy service
func NewPolicyService(repo *repository.Repository, flags *features.Flags, logger *logrus.Logger) *PolicyService {
	return &PolicyService{
		repo:   repo,
		flags:  flags,
		logger: logger,
	}
}

// GetPolicyByID retrieves a policy by ID and applies masking if needed
func (s *PolicyService) GetPolicyByID(policyID string, customerID string) (*models.PolicyResponse, error) {
	policy, err := s.repo.GetPolicyByID(policyID)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"policyId":   policyID,
			"customerId": customerID,
		}).Warn("Policy not found")
		return nil, err
	}

	// Verify the policy belongs to the requesting customer
	if policy.CustomerID != customerID {
		s.logger.WithFields(logrus.Fields{
			"policyId":   policyID,
			"customerId": customerID,
			"ownerId":    policy.CustomerID,
		}).Warn("Unauthorized access attempt")
		return nil, fmt.Errorf("unauthorized")
	}

	// Apply masking and currency based on feature flags
	maskAmounts := s.flags.ShouldMaskAmounts()
	currency := s.flags.GetCurrency()
	s.logger.WithFields(logrus.Fields{
		"policyId":    policyID,
		"customerId":  customerID,
		"maskAmounts": maskAmounts,
		"currency":    currency,
	}).Debug("Retrieving policy")

	response := policy.ToResponse(maskAmounts, currency)
	return &response, nil
}

// GetPoliciesByCustomerID retrieves all policies for a customer with optional masking
func (s *PolicyService) GetPoliciesByCustomerID(customerID string) ([]models.PolicyResponse, error) {
	policies, err := s.repo.GetPoliciesByCustomerID(customerID)
	if err != nil {
		s.logger.WithField("customerId", customerID).Error("Failed to retrieve policies")
		return nil, err
	}

	// Apply masking and currency based on feature flags
	maskAmounts := s.flags.ShouldMaskAmounts()
	currency := s.flags.GetCurrency()
	s.logger.WithFields(logrus.Fields{
		"customerId":  customerID,
		"count":       len(policies),
		"maskAmounts": maskAmounts,
		"currency":    currency,
	}).Debug("Retrieving policies")

	responses := make([]models.PolicyResponse, len(policies))
	for i, policy := range policies {
		responses[i] = policy.ToResponse(maskAmounts, currency)
	}

	return responses, nil
}

// CreatePolicy creates a new policy for a customer
func (s *PolicyService) CreatePolicy(customerID string, req models.CreatePolicyRequest) (*models.PolicyResponse, error) {
	// Use the customerID from the authenticated request
	if req.CustomerID == "" {
		req.CustomerID = customerID
	}

	// Verify customer owns this policy
	if req.CustomerID != customerID {
		s.logger.WithFields(logrus.Fields{
			"customerId":        customerID,
			"requestCustomerId": req.CustomerID,
		}).Warn("Unauthorized policy creation attempt")
		return nil, fmt.Errorf("unauthorized")
	}

	policy, err := s.repo.CreatePolicy(req)
	if err != nil {
		s.logger.WithField("customerId", customerID).Error("Failed to create policy")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"policyId":   policy.ID,
		"customerId": customerID,
		"type":       policy.Type,
	}).Info("Policy created successfully")

	// Apply masking and currency based on feature flags
	maskAmounts := s.flags.ShouldMaskAmounts()
	currency := s.flags.GetCurrency()
	response := policy.ToResponse(maskAmounts, currency)
	return &response, nil
}

// UpdatePolicy updates an existing policy
func (s *PolicyService) UpdatePolicy(policyID string, customerID string, req models.UpdatePolicyRequest) (*models.PolicyResponse, error) {
	// Get existing policy to verify ownership
	policy, err := s.repo.GetPolicyByID(policyID)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"policyId":   policyID,
			"customerId": customerID,
		}).Warn("Policy not found")
		return nil, err
	}

	// Verify the policy belongs to the requesting customer
	if policy.CustomerID != customerID {
		s.logger.WithFields(logrus.Fields{
			"policyId":   policyID,
			"customerId": customerID,
			"ownerId":    policy.CustomerID,
		}).Warn("Unauthorized update attempt")
		return nil, fmt.Errorf("unauthorized")
	}

	// Apply updates
	now := time.Now()
	if req.Status != nil {
		policy.Status = *req.Status
	}
	if req.Premium != nil {
		policy.Premium = *req.Premium
	}
	if req.EndDate != nil {
		policy.EndDate = *req.EndDate
	}
	policy.UpdatedAt = now

	// Update in repository
	updatedPolicy, err := s.repo.UpdatePolicy(policy)
	if err != nil {
		s.logger.WithField("policyId", policyID).Error("Failed to update policy")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"policyId":   policyID,
		"customerId": customerID,
	}).Info("Policy updated successfully")

	// Apply masking and currency based on feature flags
	maskAmounts := s.flags.ShouldMaskAmounts()
	currency := s.flags.GetCurrency()
	response := updatedPolicy.ToResponse(maskAmounts, currency)
	return &response, nil
}
