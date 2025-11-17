//go:build !generate

package main

import (
	"context"

	cfg "github.com/conductorone/baton-outreach/pkg/config"
	"github.com/conductorone/baton-outreach/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
)

var version = "dev"

func main() {
	ctx := context.Background()
	config.RunConnector(
		ctx,
		"baton-outreach",
		version,
		cfg.Config,
		connector.New,
		connectorrunner.WithSessionStoreEnabled(),
	)
}
