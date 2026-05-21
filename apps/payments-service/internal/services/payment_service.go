package services

import (
	"fmt"
	"time"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/payments-service/internal/features"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/payments-service/internal/models"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/payments-service/internal/repository"
	"github.com/sirupsen/logrus"
)

// PaymentService handles payment business logic
type PaymentService struct {
	repo   *repository.Repository
	flags  *features.Flags
	logger *logrus.Logger
}

// NewPaymentService creates a new payment service
func NewPaymentService(repo *repository.Repository, flags *features.Flags, logger *logrus.Logger) *PaymentService {
	return &PaymentService{
		repo:   repo,
		flags:  flags,
		logger: logger,
	}
}

// GetAllPayments returns all payments
func (s *PaymentService) GetAllPayments() ([]*models.Payment, error) {
	return s.repo.GetAllPayments()
}

// GetPaymentByID returns a payment by ID
func (s *PaymentService) GetPaymentByID(paymentID string) (*models.Payment, error) {
	return s.repo.GetPaymentByID(paymentID)
}

// CreatePayment creates a new premium payment
func (s *PaymentService) CreatePayment(policyID, customerID string, amount float64) (*models.Payment, error) {
	payment := &models.Payment{
		ID:         fmt.Sprintf("pay-%d", time.Now().UnixNano()),
		Type:       models.PaymentTypePremium,
		PolicyID:   policyID,
		CustomerID: customerID,
		Amount:     amount,
		Status:     models.PaymentStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	s.logger.WithFields(logrus.Fields{
		"paymentId":  payment.ID,
		"policyId":   policyID,
		"customerId": customerID,
		"amount":     amount,
	}).Info("Creating premium payment")

	if err := s.repo.CreatePayment(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// CreatePayout creates a new claim payout
func (s *PaymentService) CreatePayout(claimID, customerID string, amount float64) (*models.Payment, error) {
	payment := &models.Payment{
		ID:         fmt.Sprintf("pay-%d", time.Now().UnixNano()),
		Type:       models.PaymentTypePayout,
		ClaimID:    claimID,
		CustomerID: customerID,
		Amount:     amount,
		Status:     models.PaymentStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	s.logger.WithFields(logrus.Fields{
		"paymentId":  payment.ID,
		"claimId":    claimID,
		"customerId": customerID,
		"amount":     amount,
	}).Info("Creating claim payout")

	if err := s.repo.CreatePayment(payment); err != nil {
		return nil, err
	}

	// Check if instant payouts are enabled
	if s.flags.IsInstantPayoutsEnabled() && payment.Type == models.PaymentTypePayout {
		s.logger.WithField("paymentId", payment.ID).Info("Instant payouts enabled - processing immediately")
		return s.ProcessPayment(payment.ID)
	}

	s.logger.WithField("paymentId", payment.ID).Info("Instant payouts disabled - payout queued for batch processing")
	return payment, nil
}

// ProcessPayment processes a pending payment
func (s *PaymentService) ProcessPayment(paymentID string) (*models.Payment, error) {
	payment, err := s.repo.GetPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	// Validate that payment is in pending state
	if payment.Status != models.PaymentStatusPending {
		return nil, fmt.Errorf("payment already processed")
	}

	// Simulate payment processing
	s.logger.WithFields(logrus.Fields{
		"paymentId": payment.ID,
		"type":      payment.Type,
		"amount":    payment.Amount,
	}).Info("Processing payment")

	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)

	// Update payment status
	now := time.Now()
	payment.Status = models.PaymentStatusCompleted
	payment.ProcessedDate = &now
	payment.UpdatedAt = now

	if err := s.repo.UpdatePayment(payment); err != nil {
		return nil, err
	}

	s.logger.WithField("paymentId", payment.ID).Info("Payment processed successfully")

	return payment, nil
}
