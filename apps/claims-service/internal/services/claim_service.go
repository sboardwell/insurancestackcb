package services

import (
	"fmt"
	"sort"
	"time"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/claims-service/internal/features"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/claims-service/internal/models"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/claims-service/internal/repository"
	"github.com/sirupsen/logrus"
)

const (
	// AutoApprovalThreshold is the maximum amount for automatic approval (in dollars)
	AutoApprovalThreshold = 1000.0
)

// ClaimService handles business logic for claims
type ClaimService struct {
	repo   *repository.Repository
	flags  *features.Flags
	logger *logrus.Logger
}

// NewClaimService creates a new claim service
func NewClaimService(repo *repository.Repository, flags *features.Flags, logger *logrus.Logger) *ClaimService {
	return &ClaimService{
		repo:   repo,
		flags:  flags,
		logger: logger,
	}
}

// GetClaimByID retrieves a claim by ID
func (s *ClaimService) GetClaimByID(claimID string) (*models.Claim, error) {
	return s.repo.GetClaimByID(claimID)
}

// GetClaims retrieves claims with optional filters
func (s *ClaimService) GetClaims(filters *models.ClaimFilters) ([]*models.Claim, error) {
	var claims []*models.Claim

	if filters == nil || (filters.PolicyID == "" && filters.CustomerID == "" && filters.Status == "" && filters.Type == "") {
		// No filters - return all claims
		claims = s.repo.GetAllClaims()
	} else {
		// Apply filters
		claims = s.repo.GetClaimsByFilter(filters)
	}

	// Sort by submission date descending (most recent first)
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].SubmittedDate.After(claims[j].SubmittedDate)
	})

	s.logger.WithFields(logrus.Fields{
		"count":   len(claims),
		"filters": filters,
	}).Info("Retrieved claims")

	return claims, nil
}

// CreateClaim creates a new claim with governance rules applied
func (s *ClaimService) CreateClaim(req *models.CreateClaimRequest) (*models.Claim, error) {
	// Validate claim type
	if !models.ValidateClaimType(req.Type) {
		return nil, fmt.Errorf("invalid claim type: %s (must be accident, theft, or damage)", req.Type)
	}

	// Validate amount
	if req.Amount <= 0 {
		return nil, fmt.Errorf("claim amount must be greater than 0")
	}

	// Generate claim number
	claimNumber := s.generateClaimNumber()

	// Determine initial status based on auto-approval feature flag
	status := "under_review"
	autoApprovalEnabled := s.flags.IsAutoApprovalEnabled()

	// Apply governance rule: auto-approve low-value claims if feature flag is enabled
	if autoApprovalEnabled && req.Amount < AutoApprovalThreshold {
		status = "approved"
		s.logger.WithFields(logrus.Fields{
			"claimNumber": claimNumber,
			"amount":      req.Amount,
			"threshold":   AutoApprovalThreshold,
		}).Info("Auto-approved low-value claim")
	} else {
		s.logger.WithFields(logrus.Fields{
			"claimNumber":       claimNumber,
			"amount":            req.Amount,
			"autoApprovalFlag":  autoApprovalEnabled,
			"threshold":         AutoApprovalThreshold,
		}).Info("Claim requires manual review")
	}

	now := time.Now()
	claim := &models.Claim{
		ID:            s.generateClaimID(),
		PolicyID:      req.PolicyID,
		CustomerID:    req.CustomerID,
		ClaimNumber:   claimNumber,
		Type:          req.Type,
		Status:        status,
		Amount:        req.Amount,
		Description:   req.Description,
		SubmittedDate: now,
		ReviewedDate:  nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// If auto-approved, set reviewed date
	if status == "approved" {
		claim.ReviewedDate = &now
	}

	if err := s.repo.CreateClaim(claim); err != nil {
		return nil, fmt.Errorf("failed to create claim: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"claimId":     claim.ID,
		"claimNumber": claim.ClaimNumber,
		"status":      claim.Status,
		"amount":      claim.Amount,
	}).Info("Claim created successfully")

	return claim, nil
}

// UpdateClaim updates an existing claim
func (s *ClaimService) UpdateClaim(claimID string, req *models.UpdateClaimRequest) (*models.Claim, error) {
	// Get existing claim
	claim, err := s.repo.GetClaimByID(claimID)
	if err != nil {
		return nil, err
	}

	// Only allow updates to claims that are in submitted or under_review status
	if claim.Status == "approved" || claim.Status == "rejected" {
		return nil, fmt.Errorf("cannot update claim with status: %s", claim.Status)
	}

	// Update fields if provided
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return nil, fmt.Errorf("claim amount must be greater than 0")
		}
		claim.Amount = *req.Amount
	}

	if req.Description != nil {
		claim.Description = *req.Description
	}

	claim.UpdatedAt = time.Now()

	if err := s.repo.UpdateClaim(claim); err != nil {
		return nil, fmt.Errorf("failed to update claim: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"claimId":     claim.ID,
		"claimNumber": claim.ClaimNumber,
	}).Info("Claim updated successfully")

	return claim, nil
}

// UpdateClaimStatus updates the status of a claim (for approval workflows)
func (s *ClaimService) UpdateClaimStatus(claimID string, req *models.UpdateClaimStatusRequest) (*models.Claim, error) {
	// Get existing claim
	claim, err := s.repo.GetClaimByID(claimID)
	if err != nil {
		return nil, err
	}

	// Validate new status
	if !models.ValidateClaimStatus(req.Status) {
		return nil, fmt.Errorf("invalid claim status: %s", req.Status)
	}

	// Prevent status changes on already finalized claims
	if claim.Status == "approved" || claim.Status == "rejected" {
		return nil, fmt.Errorf("cannot change status of finalized claim (current status: %s)", claim.Status)
	}

	oldStatus := claim.Status
	claim.Status = req.Status
	claim.UpdatedAt = time.Now()

	// If approving or rejecting, set reviewed date
	if req.Status == "approved" || req.Status == "rejected" {
		now := time.Now()
		claim.ReviewedDate = &now
	}

	if err := s.repo.UpdateClaim(claim); err != nil {
		return nil, fmt.Errorf("failed to update claim status: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"claimId":     claim.ID,
		"claimNumber": claim.ClaimNumber,
		"oldStatus":   oldStatus,
		"newStatus":   req.Status,
		"notes":       req.Notes,
	}).Info("Claim status updated")

	return claim, nil
}

// generateClaimID generates a unique claim ID
func (s *ClaimService) generateClaimID() string {
	return fmt.Sprintf("claim-%d", time.Now().UnixNano())
}

// generateClaimNumber generates a human-readable claim number
func (s *ClaimService) generateClaimNumber() string {
	now := time.Now()
	return fmt.Sprintf("CLM-%d-%06d", now.Year(), now.Unix()%1000000)
}
