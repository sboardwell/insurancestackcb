package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/models"
	"github.com/CB-InsuranceStack/InsuranceStack/apps/pricing-engine/internal/services"
	"github.com/sirupsen/logrus"
)

// PricingHandler handles pricing-related requests
type PricingHandler struct {
	service *services.PricingService
	logger  *logrus.Logger
}

// NewPricingHandler creates a new pricing handler
func NewPricingHandler(service *services.PricingService, logger *logrus.Logger) *PricingHandler {
	return &PricingHandler{
		service: service,
		logger:  logger,
	}
}

// GetQuote handles POST /quote
func (h *PricingHandler) GetQuote(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req models.QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid request body")
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Calculate quote
	quote, err := h.service.CalculateQuote(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to calculate quote")
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return quote
	respondWithJSON(w, http.StatusOK, quote)
}

// GetRates handles GET /rates
func (h *PricingHandler) GetRates(w http.ResponseWriter, r *http.Request) {
	// Get rates
	rates := h.service.GetRates()

	// Return rates
	respondWithJSON(w, http.StatusOK, rates)
}

// Helper functions

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error encoding response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
