package features

import (
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

// Flags holds all feature flags for the application
type Flags struct {
	instantPayouts bool
	mu             sync.RWMutex
	logger         *logrus.Logger
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
	// payments.instantPayouts (default: false) - enable instant payout processing
	instantPayoutsStr := os.Getenv("FEATURE_INSTANT_PAYOUTS")
	if instantPayoutsStr != "" {
		instantPayouts, err := strconv.ParseBool(instantPayoutsStr)
		if err == nil {
			flags.instantPayouts = instantPayouts
		}
	}

	logger.WithFields(logrus.Fields{
		"instantPayouts": flags.instantPayouts,
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

// IsInstantPayoutsEnabled returns whether instant payouts are enabled
func (f *Flags) IsInstantPayoutsEnabled() bool {
	if f == nil {
		return false
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.instantPayouts
}

// SetInstantPayouts sets the instant payouts flag (for testing/admin purposes)
func (f *Flags) SetInstantPayouts(enabled bool) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.instantPayouts = enabled
	f.logger.WithField("instantPayouts", enabled).Info("Feature flag updated")
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
       InstantPayouts model.RoxFlag
       logger         *logrus.Logger
   }

4. Update Initialize function:
   func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
       flags = &Flags{
           logger: logger,
       }

       // Register feature flag: payments.instantPayouts (default: false)
       // This is a HIGH-RISK flag for financial operations
       flags.InstantPayouts = model.NewRoxFlag(false)

       // Register with CloudBees
       roxx.Register("payments", flags)

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

5. Update IsInstantPayoutsEnabled:
   func (f *Flags) IsInstantPayoutsEnabled() bool {
       if f == nil || f.InstantPayouts == nil {
           return false
       }
       return f.InstantPayouts.IsEnabled(nil)
   }

6. Update Shutdown:
   func Shutdown() {
       if flags != nil {
           roxx.Shutdown()
           flags.logger.Info("CloudBees Feature Management shutdown complete")
       }
   }

For more information, see: https://docs.cloudbees.com/docs/cloudbees-feature-management/latest/

IMPORTANT: payments.instantPayouts is a HIGH-RISK feature flag:
- Instant payouts bypass batch reconciliation and fraud detection
- Recommended rollout: 5% -> 10% -> 25% -> 50% -> 100% with monitoring at each stage
- Roll back immediately if fraud rate increases or payment failures spike
- Only enable during business hours when fraud team can monitor
*/
