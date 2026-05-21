package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/customer-service/internal/models"
	"github.com/sirupsen/logrus"
)

// Repository provides data access for customers
type Repository struct {
	customers map[string]*models.Customer
	mu        sync.RWMutex
	logger    *logrus.Logger
}

// NewRepository creates a new repository and loads data from JSON files
func NewRepository(dataPath string, logger *logrus.Logger) (*Repository, error) {
	repo := &Repository{
		customers: make(map[string]*models.Customer),
		logger:    logger,
	}

	// Load customers
	customersPath := filepath.Join(dataPath, "customers.json")
	if err := repo.loadCustomers(customersPath); err != nil {
		// If customers.json doesn't exist, start with empty repository
		logger.Warnf("Could not load customers from %s: %v. Starting with empty repository.", customersPath, err)
	} else {
		logger.Infof("Loaded %d customers from %s", len(repo.customers), dataPath)
	}

	return repo, nil
}

// loadCustomers loads customers from a JSON file
func (r *Repository) loadCustomers(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var customers []*models.Customer
	if err := json.Unmarshal(data, &customers); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, customer := range customers {
		r.customers[customer.ID] = customer
	}

	return nil
}

// GetAllCustomers returns all customers
func (r *Repository) GetAllCustomers() ([]*models.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customers := make([]*models.Customer, 0, len(r.customers))
	for _, customer := range r.customers {
		customers = append(customers, customer)
	}

	return customers, nil
}

// GetCustomerByID retrieves a customer by ID
func (r *Repository) GetCustomerByID(customerID string) (*models.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customer, exists := r.customers[customerID]
	if !exists {
		return nil, fmt.Errorf("customer not found")
	}

	return customer, nil
}

// GetCustomerByEmail retrieves a customer by email address
func (r *Repository) GetCustomerByEmail(email string) (*models.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, customer := range r.customers {
		if customer.Email == email {
			return customer, nil
		}
	}

	return nil, fmt.Errorf("customer not found")
}

// CreateCustomer creates a new customer
func (r *Repository) CreateCustomer(customer *models.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if customer with same ID already exists
	if _, exists := r.customers[customer.ID]; exists {
		return fmt.Errorf("customer with ID %s already exists", customer.ID)
	}

	r.customers[customer.ID] = customer
	return nil
}

// UpdateCustomer updates an existing customer
func (r *Repository) UpdateCustomer(customer *models.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if customer exists
	if _, exists := r.customers[customer.ID]; !exists {
		return fmt.Errorf("customer not found")
	}

	r.customers[customer.ID] = customer
	return nil
}

// DeactivateCustomer deactivates a customer (soft delete)
func (r *Repository) DeactivateCustomer(customerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if customer exists
	if _, exists := r.customers[customerID]; !exists {
		return fmt.Errorf("customer not found")
	}

	// In a real implementation, we would set a status field or deletion timestamp
	// For now, we'll just remove it from the map (hard delete for simplicity)
	delete(r.customers, customerID)
	return nil
}
