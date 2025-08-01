package connector

import (
	"context"
	"io"

	"github.com/conductorone/baton-outreach/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"golang.org/x/oauth2"
)

type Connector struct {
	client *client.OutreachClient
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(_ context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(d.client),
		newTeamBuilder(d.client),
		//newRoleBuilder(d.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Connector) Asset(_ context.Context, _ *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(_ context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Outreach",
		Description: "Baton connector to sync users, teams and roles from Outreach",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(_ context.Context) (annotations.Annotations, error) {
	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, accessToken string) (*Connector, error) {
	c, err := client.New(ctx, client.WithAccessToken(accessToken))
	if err != nil {
		return nil, err
	}

	return &Connector{
		client: c,
	}, nil
}

// NewWithTokenSource returns a new instance of the connector using a provided Token Source.
func NewWithTokenSource(ctx context.Context, tokenSource oauth2.TokenSource) (*Connector, error) {
	clientOptions := []client.ConfigOption{
		client.WithTokenSource(tokenSource),
	}

	c, err := client.New(ctx, clientOptions...)
	if err != nil {
		return nil, err
	}

	return &Connector{
		client: c,
	}, nil
}
