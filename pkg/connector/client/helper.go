package client

import (
	"context"
	"fmt"
	"time"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"golang.org/x/oauth2"
)

func GetToken(token string, resourceID *v2.ResourceId) (*pagination.Bag, string, error) {
	bag := &pagination.Bag{}
	err := bag.Unmarshal(token)
	if err != nil {
		return nil, "", err
	}

	if bag.Current() == nil {
		bag.Push(pagination.PageState{
			ResourceTypeID: resourceID.ResourceType,
			ResourceID:     resourceID.Resource,
		})
	}

	return bag, bag.PageToken(), nil
}

// ConfigOption allows configuration of the client.
type ConfigOption func(client *OutreachClient)

func WithTokenSource(tokenSource oauth2.TokenSource) ConfigOption {
	return func(client *OutreachClient) {
		client.TokenSource = tokenSource
	}
}

// WithRefreshToken receives a Refresh Token, Client ID and Client Secret from the platform to be able to renew the token when expired.
// This ConfigOption is intended for CLI executions.
func WithRefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) ConfigOption {
	return func(client *OutreachClient) {
		token := &oauth2.Token{
			AccessToken:  "",
			RefreshToken: refreshToken,
			Expiry:       time.Now().Add(-1 * time.Second),
		}

		config := oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				TokenURL: authURL,
			},
		}
		tokenSource := oauth2.ReuseTokenSource(token, config.TokenSource(ctx, token))

		client.TokenSource = tokenSource
	}
}

func WithAccessToken(accessToken string) ConfigOption {
	return func(client *OutreachClient) {
		client.TokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	}
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

func (er *ErrorResponse) Message() string {
	if er.Error == "" && er.Description == "" {
		return "Error response empty"
	}

	return fmt.Sprintf("API error response: %s | Error description: %s", er.Error, er.Description)
}
