package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"golang.org/x/oauth2"
)

const (
	baseURL = "https://api.outreach.io/api/v2"
	usersEP = "users"
	teamsEP = "teams"
	rolesEP = "roles"
)

type OutreachClient struct {
	client      *uhttp.BaseHttpClient
	TokenSource oauth2.TokenSource
}

// ConfigOption allows configuration of the client.
type ConfigOption func(client *OutreachClient)

func WithTokenSource(tokenSource oauth2.TokenSource) ConfigOption {
	return func(client *OutreachClient) {
		client.TokenSource = tokenSource
	}
}

func (c *OutreachClient) ListAllUsers(ctx context.Context, nextPageLink string) ([]*User, string, error) {
	var (
		requestURL string
		response   UsersResponse
	)

	if nextPageLink != "" {
		requestURL = nextPageLink
	} else {
		usersURL, err := url.JoinPath(baseURL, usersEP)
		if err != nil {
			return nil, "", err
		}

		requestURL = usersURL
	}

	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		requestURL,
		&response,
		nil,
	)
	if err != nil {
		return nil, "", err
	}

	var nextLink string
	if response.Links != nil {
		nextLink = response.Links.Next
	}

	return response.Results, nextLink, nil
}

func (c *OutreachClient) ListAllTeams(ctx context.Context, nextPageLink string) ([]*Team, string, error) {
	var (
		requestURL string
		response   TeamsResponse
	)

	if nextPageLink != "" {
		requestURL = nextPageLink
	} else {
		teamsURL, err := url.JoinPath(baseURL, teamsEP)
		if err != nil {
			return nil, "", err
		}

		requestURL = teamsURL
	}

	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		requestURL,
		&response,
		nil,
	)
	if err != nil {
		return nil, "", err
	}

	var nextLink string
	if response.Links != nil {
		nextLink = response.Links.Next
	}

	return response.Results, nextLink, nil
}

func (c *OutreachClient) ListAllRoles(ctx context.Context, nextPageLink string) ([]*Role, string, error) {
	var (
		requestURL string
		response   RolesResponse
	)

	if nextPageLink != "" {
		requestURL = nextPageLink
	} else {
		rolesURL, err := url.JoinPath(baseURL, rolesEP)
		if err != nil {
			return nil, "", err
		}

		requestURL = rolesURL
	}

	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		requestURL,
		&response,
		nil,
	)
	if err != nil {
		return nil, "", err
	}

	var nextLink string
	if response.Links != nil {
		nextLink = response.Links.Next
	}

	return response.Results, nextLink, nil
}

func (c *OutreachClient) doRequest(
	ctx context.Context,
	method string,
	endpointUrl string,
	res interface{},
	body interface{},
	reqOpts ...ReqOpt,
) (http.Header, error) {
	var (
		resp *http.Response
		err  error
	)

	urlAddress, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, err
	}

	for _, o := range reqOpts {
		o(urlAddress)
	}

	accessToken, err := c.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	var token string
	if accessToken != nil {
		token = accessToken.AccessToken
	}

	opts := []uhttp.RequestOption{uhttp.WithBearerToken(token), uhttp.WithAcceptJSONHeader(), uhttp.WithContentTypeJSONHeader()}
	if body != nil {
		opts = append(opts, uhttp.WithJSONBody(body))
	}

	req, err := c.client.NewRequest(
		ctx,
		method,
		urlAddress,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	switch method {
	case http.MethodGet, http.MethodPut, http.MethodPost, http.MethodPatch:
		var doOptions []uhttp.DoOption
		if res != nil {
			doOptions = append(doOptions, uhttp.WithResponse(&res))
		}
		resp, err = c.client.Do(req, doOptions...)
		if resp != nil {
			defer resp.Body.Close()
		}

	case http.MethodDelete:
		resp, err = c.client.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
	}
	if err != nil {
		return nil, err
	}

	return resp.Header, nil
}

func New(ctx context.Context, cOpts ...ConfigOption) (*OutreachClient, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	cli, err := uhttp.NewBaseHttpClientWithContext(ctx, httpClient)
	if err != nil {
		return nil, err
	}

	icClient := OutreachClient{
		client: cli,
	}
	for _, option := range cOpts {
		option(&icClient)
	}

	return &icClient, nil
}
