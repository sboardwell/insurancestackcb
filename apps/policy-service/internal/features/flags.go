package features

import (
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

// Flags holds all feature flags for the application
type Flags struct {
	maskAmounts bool
	currency    string
	mu          sync.RWMutex
	logger      *logrus.Logger
}

var flags *Flags

// Initialize sets up feature flags
// To integrate with CloudBees Feature Management:
// 1. Install: go get github.com/rollout/rox-go
// 2. Import the SDK
// 3. Replace this implementation with CloudBees Rox SDK initialization
func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
	flags = &Flags{
		logger: logger,
	}

	// Load feature flags from environment variables
	// api.maskAmounts (default: false) - mask dollar amounts in responses
	maskAmountsStr := os.Getenv("FEATURE_MASK_AMOUNTS")
	if maskAmountsStr != "" {
		maskAmounts, err := strconv.ParseBool(maskAmountsStr)
		if err == nil {
			flags.maskAmounts = maskAmounts
		}
	}

	// api.currency (default: "USD") - currency code for amounts
	currency := os.Getenv("FEATURE_CURRENCY")
	if currency == "" {
		currency = "USD" // Default to USD
	}
	flags.currency = currency

	logger.WithFields(logrus.Fields{
		"maskAmounts": flags.maskAmounts,
		"currency":    flags.currency,
	}).Info("Feature flags initialized")

	if apiKey != "" && apiKey != "dev-mode" {
		logger.Warn("CloudBees Feature Management API key provided but SDK not integrated. See flags.go for integration instructions.")
	}

	return flags, nil
}

// GetFlags returns the global flags instance
func GetFlags() *Flags {
	return flags
}

// ShouldMaskAmounts returns whether amounts should be masked in responses
func (f *Flags) ShouldMaskAmounts() bool {
	if f == nil {
		return false
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.maskAmounts
}

// SetMaskAmounts sets the mask amounts flag (for testing/admin purposes)
func (f *Flags) SetMaskAmounts(enabled bool) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.maskAmounts = enabled
	f.logger.WithField("maskAmounts", enabled).Info("Feature flag updated")
}

// GetCurrency returns the currency code for amounts
func (f *Flags) GetCurrency() string {
	if f == nil {
		return "USD"
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.currency
}

// GetCurrencyForUser returns the currency code based on user context (country)
// This demonstrates CloudBees Feature Management targeting by user properties
func (f *Flags) GetCurrencyForUser(userCountry string) string {
	if f == nil {
		return "USD"
	}

	// If FEATURE_CURRENCY is set globally, use that (environment override)
	f.mu.RLock()
	globalCurrency := f.currency
	f.mu.RUnlock()

	// If environment variable explicitly set (not default), use it
	if os.Getenv("FEATURE_CURRENCY") != "" {
		return globalCurrency
	}

	// Otherwise, use country-based targeting (simulates CloudBees targeting rules)
	currency := countryToCurrency(userCountry)

	f.logger.WithFields(logrus.Fields{
		"userCountry": userCountry,
		"currency":    currency,
	}).Debug("Currency determined by user country")

	return currency
}

// countryToCurrency maps country codes to currency codes
// This simulates CloudBees Feature Management targeting rules:
//   IF user.country == "US" THEN currency = "USD"
//   IF user.country == "UK" THEN currency = "GBP"
//   IF user.country == "FR" THEN currency = "EUR"
func countryToCurrency(country string) string {
	countryMap := map[string]string{
		"US": "USD",
		"UK": "GBP",
		"GB": "GBP", // Alternative code for United Kingdom
		"FR": "EUR",
		"DE": "EUR",
		"ES": "EUR",
		"IT": "EUR",
		"NL": "EUR",
		"BE": "EUR",
		"AT": "EUR",
		"PT": "EUR",
		"IE": "EUR",
		"CA": "CAD",
		"AU": "AUD",
		"JP": "JPY",
		"CN": "CNY",
		"IN": "INR",
		"BR": "BRL",
		"MX": "MXN",
	}

	if currency, ok := countryMap[country]; ok {
		return currency
	}

	// Default to USD if country not mapped
	return "USD"
}

// SetCurrency sets the currency code (for testing/admin purposes)
func (f *Flags) SetCurrency(currency string) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.currency = currency
	f.logger.WithField("currency", currency).Info("Feature flag updated")
}

// Shutdown gracefully shuts down the feature management system
func Shutdown() {
	if flags != nil {
		flags.logger.Info("Feature management shutdown complete")
	}
}

/*
CloudBees Feature Management Integration Guide:

To integrate with CloudBees Feature Management (Rox SDK), follow these steps:

1. Install the CloudBees Rox SDK:
   go get github.com/rollout/rox-go/core

2. Update imports:
   import (
       "github.com/rollout/rox-go/core/model"
       "github.com/rollout/rox-go/core/roxx"
   )

3. Replace the Flags struct:
   type Flags struct {
       MaskAmounts model.RoxFlag
       logger      *logrus.Logger
   }

4. Update Initialize function:
   func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
       flags = &Flags{
           logger: logger,
       }

       // Register feature flag: api.maskAmounts (default: false)
       flags.MaskAmounts = model.NewRoxFlag(false)

       // Register with CloudBees
       roxx.Register("api", flags)

       // Setup Rox with API key
       options := roxx.NewRoxOptions(roxx.RoxOptionsBuilder{})
       <-roxx.Setup(apiKey, options)

       logger.Info("CloudBees Feature Management initialized successfully")

       // Fetch latest feature flags
       go func() {
           roxx.Fetch()
           logger.Info("Initial feature flags fetched")
       }()

       return flags, nil
   }

5. Update ShouldMaskAmounts:
   func (f *Flags) ShouldMaskAmounts() bool {
       if f == nil || f.MaskAmounts == nil {
           return false
       }
       return f.MaskAmounts.IsEnabled(nil)
   }

6. Update Shutdown:
   func Shutdown() {
       if flags != nil {
           roxx.Shutdown()
           flags.logger.Info("CloudBees Feature Management shutdown complete")
       }
   }

For more information, see: https://docs.cloudbees.com/docs/cloudbees-feature-management/latest/
*/
