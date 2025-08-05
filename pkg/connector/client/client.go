package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"golang.org/x/oauth2"
)

const (
	baseURL    = "https://api.outreach.io/api/v2"
	usersEP    = "users"
	teamsEP    = "teams"
	profilesEP = "profiles"
)

type OutreachClient struct {
	client      *uhttp.BaseHttpClient
	TokenSource oauth2.TokenSource
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

// ConfigOption allows configuration of the client.
type ConfigOption func(client *OutreachClient)

func WithTokenSource(tokenSource oauth2.TokenSource) ConfigOption {
	return func(client *OutreachClient) {
		client.TokenSource = tokenSource
	}
}

func WithAccessToken(accessToken string) ConfigOption {
	return func(client *OutreachClient) {
		client.TokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	}
}

func (c *OutreachClient) ListAllUsers(ctx context.Context, nextPageLink string) ([]*User, string, *v2.RateLimitDescription, error) {
	var (
		requestURL string
		response   UsersResponse
	)

	if nextPageLink != "" {
		requestURL = nextPageLink
	} else {
		usersURL, err := url.JoinPath(baseURL, usersEP)
		if err != nil {
			return nil, "", nil, err
		}

		requestURL = usersURL
	}

	rateLimitDescription := &v2.RateLimitDescription{}
	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		requestURL,
		&response,
		nil,
		rateLimitDescription,
	)
	if err != nil {
		return nil, "", rateLimitDescription, err
	}

	var nextLink string
	if response.Links != nil {
		nextLink = response.Links.Next
	}

	return response.Results, nextLink, rateLimitDescription, nil
}

func (c *OutreachClient) GetUserByID(ctx context.Context, userID string) (*User, *v2.RateLimitDescription, error) {
	var response struct {
		User *User `json:"data"`
	}

	userURL, err := url.JoinPath(baseURL, usersEP, userID)
	if err != nil {
		return nil, nil, err
	}

	rateLimitDescription := &v2.RateLimitDescription{}
	_, err = c.doRequest(
		ctx,
		http.MethodGet,
		userURL,
		&response,
		nil,
		rateLimitDescription,
	)
	if err != nil {
		return nil, rateLimitDescription, err
	}

	return response.User, rateLimitDescription, nil
}

func (c *OutreachClient) UpdateUserProfile(ctx context.Context, userID string, profileID int) (*v2.RateLimitDescription, error) {
	var requestBody struct {
		Data UpdateUsersProfileBody `json:"data"`
	}

	newProfile := DataDetailPair{
		Id:   profileID,
		Type: "profile",
	}

	numericUserID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	requestBody.Data = UpdateUsersProfileBody{
		Id:   numericUserID,
		Type: "user",
		Relationships: UserProfileRelationships{
			Profile: struct {
				Data DataDetailPair `json:"data"`
			}{
				Data: newProfile,
			},
		},
	}

	userURL, err := url.JoinPath(baseURL, usersEP, userID)
	if err != nil {
		return nil, err
	}

	rateLimitDescription := &v2.RateLimitDescription{}
	_, err = c.doRequest(
		ctx,
		http.MethodPatch,
		userURL,
		nil,
		requestBody,
		rateLimitDescription,
	)
	if err != nil {
		return rateLimitDescription, err
	}

	return rateLimitDescription, nil
}

func (c *OutreachClient) DisableUser(ctx context.Context, userID string) (*v2.RateLimitDescription, error) {
	var requestBody struct {
		Data UserLockStatusUpdate `json:"data"`
	}

	numericUserID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	requestBody.Data = UserLockStatusUpdate{
		Id:   numericUserID,
		Type: "user",
		Attributes: struct {
			Locked bool `json:"locked"`
		}{
			Locked: true,
		},
	}

	userURL, err := url.JoinPath(baseURL, usersEP, userID)
	if err != nil {
		return nil, err
	}

	rateLimitDescription := &v2.RateLimitDescription{}
	_, err = c.doRequest(
		ctx,
		http.MethodPatch,
		userURL,
		nil,
		requestBody,
		rateLimitDescription,
	)
	if err != nil {
		return rateLimitDescription, err
	}

	return rateLimitDescription, nil
}

func (c *OutreachClient) CreateUser(ctx context.Context, newUserInfo NewUserBody) (*User, *v2.RateLimitDescription, error) {
	var response struct {
		User *User `json:"data"`
	}

	userURL, err := url.JoinPath(baseURL, usersEP)
	if err != nil {
		return nil, nil, err
	}

	rateLimitDescription := &v2.RateLimitDescription{}
	_, err = c.doRequest(
		ctx,
		http.MethodPost,
		userURL,
		&response,
		newUserInfo,
		rateLimitDescription,
	)
	if err != nil {
		return nil, rateLimitDescription, err
	}

	return response.User, rateLimitDescription, nil
}

func (c *OutreachClient) ListAllTeams(ctx context.Context, nextPageLink string) ([]*Team, string, *v2.RateLimitDescription, error) {
	var (
		requestURL string
		response   TeamsResponse
	)

	if nextPageLink != "" {
		requestURL = nextPageLink
	} else {
		teamsURL, err := url.JoinPath(baseURL, teamsEP)
		if err != nil {
			return nil, "", nil, err
		}

		requestURL = teamsURL
	}

	rateLimitDescription := &v2.RateLimitDescription{}
	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		requestURL,
		&response,
		nil,
		rateLimitDescription,
	)
	if err != nil {
		return nil, "", rateLimitDescription, err
	}

	var nextLink string
	if response.Links != nil {
		nextLink = response.Links.Next
	}

	return response.Results, nextLink, rateLimitDescription, nil
}

func (c *OutreachClient) GetTeamByID(ctx context.Context, teamID string) (*Team, *v2.RateLimitDescription, error) {
	var response struct {
		Team *Team `json:"data"`
	}

	teamURL, err := url.JoinPath(baseURL, teamsEP, teamID)
	if err != nil {
		return nil, nil, err
	}

	rateLimitDescription := &v2.RateLimitDescription{}
	_, err = c.doRequest(
		ctx,
		http.MethodGet,
		teamURL,
		&response,
		nil,
		rateLimitDescription,
	)
	if err != nil {
		return nil, rateLimitDescription, err
	}

	return response.Team, rateLimitDescription, nil
}

func (c *OutreachClient) ListAllProfiles(ctx context.Context, nextPageLink string) ([]*Profile, string, *v2.RateLimitDescription, error) {
	var (
		requestURL string
		response   ProfilesResponse
	)

	if nextPageLink != "" {
		requestURL = nextPageLink
	} else {
		rolesURL, err := url.JoinPath(baseURL, profilesEP)
		if err != nil {
			return nil, "", nil, err
		}

		requestURL = rolesURL
	}

	rateLimitDescription := &v2.RateLimitDescription{}
	_, err := c.doRequest(
		ctx,
		http.MethodGet,
		requestURL,
		&response,
		nil,
		rateLimitDescription,
	)
	if err != nil {
		return nil, "", rateLimitDescription, err
	}

	var nextLink string
	if response.Links != nil {
		nextLink = response.Links.Next
	}

	return response.Results, nextLink, rateLimitDescription, nil
}

func (c *OutreachClient) UpdateTeamMembers(ctx context.Context, teamID string, teamMembers []DataDetailPair) (*v2.RateLimitDescription, error) {
	var requestBody struct {
		Data UpdateTeamBody `json:"data"`
	}

	numericTeamID, err := strconv.Atoi(teamID)
	if err != nil {
		return nil, err
	}

	requestBody.Data = UpdateTeamBody{
		Id:   numericTeamID,
		Type: "team",
		Relationships: UpdateTeamRelationships{
			Users: struct {
				Data []DataDetailPair `json:"data"`
			}{
				Data: teamMembers,
			},
		},
	}

	teamURL, err := url.JoinPath(baseURL, teamsEP, teamID)
	if err != nil {
		return nil, err
	}

	rateLimitDescription := &v2.RateLimitDescription{}
	_, err = c.doRequest(
		ctx,
		http.MethodPatch,
		teamURL,
		nil,
		requestBody,
		rateLimitDescription,
	)
	if err != nil {
		return rateLimitDescription, err
	}

	return rateLimitDescription, nil
}

func (c *OutreachClient) doRequest(
	ctx context.Context,
	method string,
	endpointUrl string,
	res interface{},
	body interface{},
	rateLimitDescription *v2.RateLimitDescription,
	reqOpts ...ReqOpt,
) (http.Header, error) {
	var (
		resp           *http.Response
		responseHeader http.Header
		err            error
		errResponse    ErrorResponse
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

	opts := []uhttp.RequestOption{uhttp.WithBearerToken(token)}
	if body != nil {
		opts = append(opts, uhttp.WithJSONBody(body), uhttp.WithContentType("application/vnd.api+json"))
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
		doOptions := []uhttp.DoOption{uhttp.WithErrorResponse(&errResponse)}
		if res != nil {
			doOptions = append(doOptions, uhttp.WithResponse(&res))
		}
		if rateLimitDescription != nil {
			doOptions = append(doOptions, uhttp.WithRatelimitData(rateLimitDescription))
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

	if resp != nil {
		responseHeader = resp.Header
	}

	return responseHeader, nil
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
