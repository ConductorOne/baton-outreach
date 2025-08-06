package config

import (
	"os"
	"testing"

	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/stretchr/testify/assert"
)

var accessToken = os.Getenv("OUTREACH_ACCESS_TOKEN")

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Outreach
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Outreach{
				AccessToken: accessToken,
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
