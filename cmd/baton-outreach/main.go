//go:build !generate

package main

import (
	"context"
	"fmt"
	"os"

	cfg "github.com/conductorone/baton-outreach/pkg/config"
	"github.com/conductorone/baton-outreach/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-outreach",
		getConnector,
		cfg.Config,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, config *cfg.Outreach) (types.ConnectorServer, error) {
	var cb *connector.Connector
	l := ctxzap.Extract(ctx)
	if err := field.Validate(cfg.Config, config); err != nil {
		return nil, err
	}

	accessToken := config.AccessToken
	if accessToken != "" {
		cbWithAccessToken, err := connector.NewWithAccessToken(ctx, accessToken)
		if err != nil {
			l.Error("error creating connector with access token", zap.Error(err))
			return nil, err
		}

		cb = cbWithAccessToken
	}

	refreshToken := config.RefreshToken
	outreachClientID := config.OutreachClientId
	outreachClientSecret := config.OutreachClientSecret

	if outreachClientID != "" && outreachClientSecret != "" && refreshToken != "" {
		cbWithRefreshToken, err := connector.NewWithRefreshToken(ctx, outreachClientID, outreachClientSecret, refreshToken)
		if err != nil {
			l.Error("error creating connector with refresh token", zap.Error(err))
			return nil, err
		}

		cb = cbWithRefreshToken
	}

	if cb == nil {
		return nil, fmt.Errorf("connector initialization failed")
	}

	conn, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}
	return conn, nil
}
