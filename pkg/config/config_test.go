package config

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	// There is no point on validating config since it's an OAuth connector.
	// Providing and access token is optional for a CLI execution and doing tests.
}
