package features

import (
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

// Flags holds all feature flags for the application
type Flags struct {
	autoApproval bool
	mu           sync.RWMutex
	logger       *logrus.Logger
}

var flags *Flags

// Initialize sets up feature flags
// To integrate with CloudBees Feature Management:
// See the integration guide at the bottom of this file
func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
	flags = &Flags{
		logger: logger,
	}

	// Load feature flags from environment variables
	// claims.autoApproval (default: false) - enable automatic approval of low-value claims
	autoApprovalStr := os.Getenv("FEATURE_AUTO_APPROVAL")
	if autoApprovalStr != "" {
		autoApproval, err := strconv.ParseBool(autoApprovalStr)
		if err == nil {
			flags.autoApproval = autoApproval
		}
	}

	logger.WithFields(logrus.Fields{
		"autoApproval": flags.autoApproval,
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

// IsAutoApprovalEnabled returns whether automatic approval for low-value claims is enabled
func (f *Flags) IsAutoApprovalEnabled() bool {
	if f == nil {
		return false
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.autoApproval
}

// SetAutoApproval sets the auto approval flag (for testing/admin purposes)
func (f *Flags) SetAutoApproval(enabled bool) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.autoApproval = enabled
	f.logger.WithField("autoApproval", enabled).Info("Feature flag updated")
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
       AutoApproval model.RoxFlag
       logger       *logrus.Logger
   }

4. Update Initialize function:
   func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
       flags = &Flags{
           logger: logger,
       }

       // Register feature flag: claims.autoApproval (default: false)
       flags.AutoApproval = model.NewRoxFlag(false)

       // Register with CloudBees
       roxx.Register("claims", flags)

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

5. Update IsAutoApprovalEnabled:
   func (f *Flags) IsAutoApprovalEnabled() bool {
       if f == nil || f.AutoApproval == nil {
           return false
       }
       return f.AutoApproval.IsEnabled(nil)
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
