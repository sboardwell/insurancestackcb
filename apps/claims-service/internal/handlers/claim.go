package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/claims-service/internal/middleware"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/claims-service/internal/models"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/claims-service/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// ClaimHandler handles claim-related HTTP requests
type ClaimHandler struct {
	service *services.ClaimService
	logger  *logrus.Logger
}

// NewClaimHandler creates a new claim handler
func NewClaimHandler(service *services.ClaimService, logger *logrus.Logger) *ClaimHandler {
	return &ClaimHandler{
		service: service,
		logger:  logger,
	}
}

// GetClaims handles GET /claims
// Supports query parameters:
// - policyId: filter by policy ID
// - customerId: filter by customer ID
// - status: filter by status (submitted/under_review/approved/rejected)
// - type: filter by type (accident/theft/damage)
func (h *ClaimHandler) GetClaims(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		h.logger.Warn("User ID not found in context")
		h.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	filters := &models.ClaimFilters{
		PolicyID:   query.Get("policyId"),
		CustomerID: query.Get("customerId"),
		Status:     query.Get("status"),
		Type:       query.Get("type"),
	}

	// Get claims with filters
	claims, err := h.service.GetClaims(filters)
	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve claims")
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve claims")
		return
	}

	h.respondJSON(w, http.StatusOK, claims)
}

// GetClaimByID handles GET /claims/{id}
func (h *ClaimHandler) GetClaimByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	claimID := vars["id"]

	if claimID == "" {
		h.respondError(w, http.StatusBadRequest, "Claim ID is required")
		return
	}

	claim, err := h.service.GetClaimByID(claimID)
	if err != nil {
		h.logger.WithError(err).WithField("claimId", claimID).Warn("Claim not found")
		h.respondError(w, http.StatusNotFound, "Claim not found")
		return
	}

	h.respondJSON(w, http.StatusOK, claim)
}

// CreateClaim handles POST /claims
func (h *ClaimHandler) CreateClaim(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		h.logger.Warn("User ID not found in context")
		h.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid request body")
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.PolicyID == "" {
		h.respondError(w, http.StatusBadRequest, "policyId is required")
		return
	}
	if req.CustomerID == "" {
		h.respondError(w, http.StatusBadRequest, "customerId is required")
		return
	}
	if req.Type == "" {
		h.respondError(w, http.StatusBadRequest, "type is required")
		return
	}
	if req.Amount <= 0 {
		h.respondError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}
	if req.Description == "" {
		h.respondError(w, http.StatusBadRequest, "description is required")
		return
	}

	claim, err := h.service.CreateClaim(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create claim")
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.WithFields(logrus.Fields{
		"claimId":     claim.ID,
		"claimNumber": claim.ClaimNumber,
		"status":      claim.Status,
		"userId":      userID,
	}).Info("Claim created via API")

	// Return a success response with claim details
	response := map[string]interface{}{
		"id":          claim.ID,
		"claimNumber": claim.ClaimNumber,
		"status":      claim.Status,
		"message":     "Claim submitted successfully",
	}

	h.respondJSON(w, http.StatusCreated, response)
}

// UpdateClaim handles PUT /claims/{id}
func (h *ClaimHandler) UpdateClaim(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		h.logger.Warn("User ID not found in context")
		h.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	claimID := vars["id"]

	if claimID == "" {
		h.respondError(w, http.StatusBadRequest, "Claim ID is required")
		return
	}

	var req models.UpdateClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid request body")
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claim, err := h.service.UpdateClaim(claimID, &req)
	if err != nil {
		h.logger.WithError(err).WithField("claimId", claimID).Error("Failed to update claim")
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.WithFields(logrus.Fields{
		"claimId": claim.ID,
		"userId":  userID,
	}).Info("Claim updated via API")

	h.respondJSON(w, http.StatusOK, claim)
}

// UpdateClaimStatus handles PUT /claims/{id}/status
func (h *ClaimHandler) UpdateClaimStatus(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		h.logger.Warn("User ID not found in context")
		h.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	claimID := vars["id"]

	if claimID == "" {
		h.respondError(w, http.StatusBadRequest, "Claim ID is required")
		return
	}

	var req models.UpdateClaimStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid request body")
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Status == "" {
		h.respondError(w, http.StatusBadRequest, "status is required")
		return
	}

	claim, err := h.service.UpdateClaimStatus(claimID, &req)
	if err != nil {
		h.logger.WithError(err).WithField("claimId", claimID).Error("Failed to update claim status")
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.WithFields(logrus.Fields{
		"claimId":   claim.ID,
		"newStatus": claim.Status,
		"userId":    userID,
	}).Info("Claim status updated via API")

	h.respondJSON(w, http.StatusOK, claim)
}

// respondJSON sends a JSON response
func (h *ClaimHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.WithError(err).Error("Failed to encode response")
	}
}

// respondError sends an error response
func (h *ClaimHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
