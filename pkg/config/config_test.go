package config

import (
	"testing"

	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Outreach
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Outreach{
				AccessToken: "abc123",
			},
			wantErr: false,
		},
		{
			name: "invalid config - missing required fields",
			config: &Outreach{
				AccessToken: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := field.Validate(Config, tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
