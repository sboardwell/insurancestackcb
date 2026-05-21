package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/customer-service/internal/models"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/customer-service/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// CustomerHandler handles customer-related requests
type CustomerHandler struct {
	customerService *services.CustomerService
	logger          *logrus.Logger
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(customerService *services.CustomerService, logger *logrus.Logger) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
		logger:          logger,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a success message response
type SuccessResponse struct {
	Message string `json:"message"`
}

// GetCustomers handles GET /customers - returns all customers
func (h *CustomerHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.customerService.GetAllCustomers()
	if err != nil {
		h.logger.WithError(err).Error("Failed to get customers")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve customers",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

// GetCustomerByID handles GET /customers/{id} - returns a specific customer
func (h *CustomerHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["id"]

	if customerID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Customer ID is required",
		})
		return
	}

	customer, err := h.customerService.GetCustomerByID(customerID)
	if err != nil {
		h.logger.WithError(err).WithField("customerId", customerID).Error("Failed to get customer")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "not_found",
			Message: "Customer not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

// CreateCustomer handles POST /customers - creates a new customer
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid request body",
		})
		return
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "FirstName, LastName, and Email are required",
		})
		return
	}

	customer, err := h.customerService.CreateCustomer(req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create customer")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create customer",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// UpdateCustomer handles PUT /customers/{id} - updates an existing customer
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["id"]

	if customerID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Customer ID is required",
		})
		return
	}

	var req models.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid request body",
		})
		return
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "FirstName, LastName, and Email are required",
		})
		return
	}

	customer, err := h.customerService.UpdateCustomer(customerID, req)
	if err != nil {
		h.logger.WithError(err).WithField("customerId", customerID).Error("Failed to update customer")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "not_found",
			Message: "Customer not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

// DeactivateCustomer handles DELETE /customers/{id} - deactivates a customer
func (h *CustomerHandler) DeactivateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["id"]

	if customerID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Customer ID is required",
		})
		return
	}

	err := h.customerService.DeactivateCustomer(customerID)
	if err != nil {
		h.logger.WithError(err).WithField("customerId", customerID).Error("Failed to deactivate customer")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "not_found",
			Message: "Customer not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "Customer deactivated successfully",
	})
}
