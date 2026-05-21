package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/policy-service/internal/middleware"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/policy-service/internal/models"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/policy-service/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// PolicyHandler handles policy-related requests
type PolicyHandler struct {
	policyService *services.PolicyService
	logger        *logrus.Logger
}

// NewPolicyHandler creates a new policy handler
func NewPolicyHandler(policyService *services.PolicyService, logger *logrus.Logger) *PolicyHandler {
	return &PolicyHandler{
		policyService: policyService,
		logger:        logger,
	}
}

// GetPolicies handles GET /policies - returns all policies for current customer
func (h *PolicyHandler) GetPolicies(w http.ResponseWriter, r *http.Request) {
	customerID := middleware.GetUserID(r)

	policies, err := h.policyService.GetPoliciesByCustomerID(customerID)
	if err != nil {
		h.logger.WithError(err).WithField("customerId", customerID).Error("Failed to get policies")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve policies",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(policies)
}

// GetPolicyByID handles GET /policies/{id} - returns a specific policy
func (h *PolicyHandler) GetPolicyByID(w http.ResponseWriter, r *http.Request) {
	customerID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	policyID := vars["id"]

	if policyID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Policy ID is required",
		})
		return
	}

	policy, err := h.policyService.GetPolicyByID(policyID, customerID)
	if err != nil {
		if err.Error() == "unauthorized" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error:   "forbidden",
				Message: "You do not have access to this policy",
			})
			return
		}

		h.logger.WithError(err).WithFields(logrus.Fields{
			"customerId": customerID,
			"policyId":   policyID,
		}).Error("Failed to get policy")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "not_found",
			Message: "Policy not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(policy)
}

// CreatePolicy handles POST /policies - creates a new policy
func (h *PolicyHandler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	customerID := middleware.GetUserID(r)

	var req models.CreatePolicyRequest
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
	if req.PolicyNumber == "" || req.Type == "" || req.Premium <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Missing required fields: policyNumber, type, and premium are required",
		})
		return
	}

	// Validate policy type
	if req.Type != "auto" && req.Type != "home" && req.Type != "life" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid policy type. Must be one of: auto, home, life",
		})
		return
	}

	policy, err := h.policyService.CreatePolicy(customerID, req)
	if err != nil {
		h.logger.WithError(err).WithField("customerId", customerID).Error("Failed to create policy")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create policy",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(policy)
}

// UpdatePolicy handles PUT /policies/{id} - updates a policy
func (h *PolicyHandler) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	customerID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	policyID := vars["id"]

	if policyID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Policy ID is required",
		})
		return
	}

	var req models.UpdatePolicyRequest
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

	// Validate status if provided
	if req.Status != nil {
		status := *req.Status
		if status != "active" && status != "lapsed" && status != "cancelled" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error:   "bad_request",
				Message: "Invalid status. Must be one of: active, lapsed, cancelled",
			})
			return
		}
	}

	policy, err := h.policyService.UpdatePolicy(policyID, customerID, req)
	if err != nil {
		if err.Error() == "unauthorized" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error:   "forbidden",
				Message: "You do not have access to this policy",
			})
			return
		}

		h.logger.WithError(err).WithFields(logrus.Fields{
			"customerId": customerID,
			"policyId":   policyID,
		}).Error("Failed to update policy")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "not_found",
			Message: "Policy not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(policy)
}

// DeletePolicy handles DELETE /policies/{id} - cancels a policy
func (h *PolicyHandler) DeletePolicy(w http.ResponseWriter, r *http.Request) {
	customerID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	policyID := vars["id"]

	if policyID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "bad_request",
			Message: "Policy ID is required",
		})
		return
	}

	// Create an update request to set status to cancelled
	cancelled := "cancelled"
	req := models.UpdatePolicyRequest{
		Status: &cancelled,
	}

	policy, err := h.policyService.UpdatePolicy(policyID, customerID, req)
	if err != nil {
		if err.Error() == "unauthorized" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error:   "forbidden",
				Message: "You do not have access to this policy",
			})
			return
		}

		h.logger.WithError(err).WithFields(logrus.Fields{
			"customerId": customerID,
			"policyId":   policyID,
		}).Error("Failed to cancel policy")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "not_found",
			Message: "Policy not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(policy)
}
