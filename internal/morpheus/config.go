package morpheus

import "github.com/hoeg/bluepill/internal/values"

type EnforcementConfig struct {
	Whitelist values.Whitelist
	Enforce   bool
}

func NewEnforcementConfig() *EnforcementConfig {
	return nil
}
