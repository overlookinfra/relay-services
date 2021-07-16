package opt

import (
	"fmt"
	"regexp"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Debug determines whether this server starts with debugging enabled.
	Debug bool

	// Name is the name of this deployment, if known.
	Name string

	// Namespace is the Kubernetes namespace this deployment is running in, if
	// known.
	Namespace string

	// ControllerMaxReconcileBackoffDuration is the amount of time the
	// controller may wait to reprocess an object that has encountered an error.
	ControllerMaxReconcileBackoffDuration time.Duration

	// ControllerConditionTypeInclusionPattern is a regular expression used to
	// match node condition types to determine whether they should be
	// propagated.
	ControllerConditionTypeInclusionPattern *regexp.Regexp
}

func NewConfig(defaultName string) (*Config, error) {
	viper.SetEnvPrefix("relay_node_tainter")
	viper.AutomaticEnv()

	viper.SetDefault("name", defaultName)
	viper.SetDefault("controller_max_reconcile_backoff_duration", 1*time.Minute)

	cfg := &Config{
		Debug:                                 viper.GetBool("debug"),
		Name:                                  viper.GetString("name"),
		Namespace:                             viper.GetString("namespace"),
		ControllerMaxReconcileBackoffDuration: viper.GetDuration("controller_max_reconcile_backoff_duration"),
	}

	if pattern, err := regexp.Compile(viper.GetString("controller_condition_type_inclusion_pattern")); err != nil {
		return nil, fmt.Errorf("failed to parse regular expression for controller condition type inclusion pattern: %w", err)
	} else {
		cfg.ControllerConditionTypeInclusionPattern = pattern
	}

	return cfg, nil
}
