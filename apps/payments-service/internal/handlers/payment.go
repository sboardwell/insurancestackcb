package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/payments-service/internal/models"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/payments-service/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	service *services.PaymentService
	logger  *logrus.Logger
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(service *services.PaymentService, logger *logrus.Logger) *PaymentHandler {
	return &PaymentHandler{
		service: service,
		logger:  logger,
	}
}

// GetPayments handles GET /payments
func (h *PaymentHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := h.service.GetAllPayments()
	if err != nil {
		h.logger.WithError(err).Error("Failed to get payments")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

// GetPaymentByID handles GET /payments/{id}
func (h *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID := vars["id"]

	payment, err := h.service.GetPaymentByID(paymentID)
	if err != nil {
		h.logger.WithError(err).WithField("paymentId", paymentID).Error("Failed to get payment")
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

// CreatePayment handles POST /payments
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode payment request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.WithError(err).Error("Payment request validation failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := h.service.CreatePayment(req.PolicyID, req.CustomerID, req.Amount)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create payment")
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

// CreatePayout handles POST /payouts
func (h *PaymentHandler) CreatePayout(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePayoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode payout request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.WithError(err).Error("Payout request validation failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := h.service.CreatePayout(req.ClaimID, req.CustomerID, req.Amount)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create payout")
		http.Error(w, "Failed to create payout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

// ProcessPayment handles PUT /payments/{id}/process
func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID := vars["id"]

	payment, err := h.service.ProcessPayment(paymentID)
	if err != nil {
		h.logger.WithError(err).WithField("paymentId", paymentID).Error("Failed to process payment")

		// Check if it's a not found error
		if err.Error() == "payment not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Check if it's a validation error (already processed)
		if err.Error() == "payment already processed" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}
