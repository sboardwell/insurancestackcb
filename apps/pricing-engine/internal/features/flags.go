package features

import (
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

// Flags holds all feature flags for the application
type Flags struct {
	dynamicRates bool
	mu           sync.RWMutex
	logger       *logrus.Logger
}

var flags *Flags

// Initialize sets up feature flags
// To integrate with CloudBees Feature Management:
// 1. Install: go get github.com/rollout/rox-go/v5/core
// 2. Import the SDK
// 3. Replace this implementation with CloudBees Rox SDK initialization
func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
	flags = &Flags{
		logger: logger,
	}

	// Load feature flags from environment variables
	// pricing.dynamicRates (default: false) - enable real-time rate adjustments
	dynamicRatesStr := os.Getenv("FEATURE_DYNAMIC_RATES")
	if dynamicRatesStr != "" {
		dynamicRates, err := strconv.ParseBool(dynamicRatesStr)
		if err == nil {
			flags.dynamicRates = dynamicRates
		}
	}

	logger.WithFields(logrus.Fields{
		"dynamicRates": flags.dynamicRates,
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

// IsDynamicRatesEnabled returns whether dynamic rates are enabled
func (f *Flags) IsDynamicRatesEnabled() bool {
	if f == nil {
		return false
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.dynamicRates
}

// SetDynamicRates sets the dynamic rates flag (for testing/admin purposes)
func (f *Flags) SetDynamicRates(enabled bool) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.dynamicRates = enabled
	f.logger.WithField("dynamicRates", enabled).Info("Feature flag updated")
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
   go get github.com/rollout/rox-go/v5/core

2. Update imports:
   import (
       "github.com/rollout/rox-go/v5/core"
   )

3. Replace the Flags struct:
   type Flags struct {
       DynamicRates *core.RoxFlag
       logger       *logrus.Logger
   }

4. Update Initialize function:
   func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
       flags = &Flags{
           logger: logger,
       }

       // Register feature flag: pricing.dynamicRates (default: false)
       flags.DynamicRates = core.NewRoxFlag(false)

       // Register with CloudBees
       core.Register("pricing", flags)

       // Setup Rox with API key
       options := core.NewRoxOptions(core.RoxOptionsBuilder{})
       <-core.Setup(apiKey, options)

       logger.Info("CloudBees Feature Management initialized successfully")

       // Fetch latest feature flags
       go func() {
           core.Fetch()
           logger.Info("Initial feature flags fetched")
       }()

       return flags, nil
   }

5. Update IsDynamicRatesEnabled:
   func (f *Flags) IsDynamicRatesEnabled() bool {
       if f == nil || f.DynamicRates == nil {
           return false
       }
       return f.DynamicRates.IsEnabled(nil)
   }

6. Update Shutdown:
   func Shutdown() {
       if flags != nil {
           core.Shutdown()
           flags.logger.Info("CloudBees Feature Management shutdown complete")
       }
   }

Feature Flags:
- pricing.dynamicRates (default: false) - enable real-time rate adjustments

Environment Variables (Current Implementation):
- FEATURE_DYNAMIC_RATES: Set to "true" to enable dynamic pricing

For more information, see: https://docs.cloudbees.com/docs/cloudbees-feature-management/latest/
*/
