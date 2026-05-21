package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/payments-service/internal/models"
	"github.com/sirupsen/logrus"
)

// Repository provides data access for payments
type Repository struct {
	payments map[string]*models.Payment
	mu       sync.RWMutex
	logger   *logrus.Logger
}

// NewRepository creates a new repository and loads data from JSON files
func NewRepository(dataPath string, logger *logrus.Logger) (*Repository, error) {
	repo := &Repository{
		payments: make(map[string]*models.Payment),
		logger:   logger,
	}

	// Try to load payments if the file exists
	paymentsFile := filepath.Join(dataPath, "payments.json")
	if _, err := os.Stat(paymentsFile); err == nil {
		if err := repo.loadPayments(paymentsFile); err != nil {
			logger.WithError(err).Warn("Failed to load payments, starting with empty data")
		}
	} else {
		logger.Info("No payments data file found, starting with empty data")
	}

	logger.Infof("Loaded %d payments from %s", len(repo.payments), dataPath)

	return repo, nil
}

// loadPayments loads payments from a JSON file
func (r *Repository) loadPayments(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var payments []*models.Payment
	if err := json.Unmarshal(data, &payments); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, payment := range payments {
		r.payments[payment.ID] = payment
	}

	return nil
}

// GetAllPayments returns all payments
func (r *Repository) GetAllPayments() ([]*models.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	payments := make([]*models.Payment, 0, len(r.payments))
	for _, payment := range r.payments {
		payments = append(payments, payment)
	}

	return payments, nil
}

// GetPaymentByID retrieves a payment by ID
func (r *Repository) GetPaymentByID(paymentID string) (*models.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	payment, exists := r.payments[paymentID]
	if !exists {
		return nil, fmt.Errorf("payment not found")
	}

	return payment, nil
}

// CreatePayment creates a new payment
func (r *Repository) CreatePayment(payment *models.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.payments[payment.ID] = payment
	return nil
}

// UpdatePayment updates an existing payment
func (r *Repository) UpdatePayment(payment *models.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.payments[payment.ID]; !exists {
		return fmt.Errorf("payment not found")
	}

	r.payments[payment.ID] = payment
	return nil
}
