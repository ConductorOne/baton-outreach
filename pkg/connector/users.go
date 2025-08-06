package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-sdk/pkg/annotations"

	"github.com/conductorone/baton-outreach/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client *client.OutreachClient
}

func (b *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

func (b *userBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var (
		userResources []*v2.Resource
		nextPageToken string
	)
	outAnnotations := annotations.Annotations{}

	bag, nextPage, err := client.GetToken(pToken.Token, &v2.ResourceId{ResourceType: userResourceType.Id})
	if err != nil {
		return nil, "", nil, err
	}

	users, nextPageLink, rateLimitData, err := b.client.ListAllUsers(ctx, nextPage)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return nil, "", outAnnotations, err
	}

	for _, user := range users {
		userResource, err := parseIntoUserResource(*user)
		if err != nil {
			return nil, "", outAnnotations, err
		}

		userResources = append(userResources, userResource)
	}

	if nextPageLink != "" {
		nextPageToken, err = bag.NextToken(nextPageLink)
		if err != nil {
			return nil, "", outAnnotations, err
		}
	}

	return userResources, nextPageToken, outAnnotations, nil
}

// Entitlements always returns an empty slice for users.
func (b *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants implements the Grants function for profiles resource.
func (b *userBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var grantResources []*v2.Grant
	outAnnotations := annotations.Annotations{}

	userID := resource.Id.Resource

	user, rateLimitData, err := b.client.GetUserByID(ctx, userID)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}

		return nil, "", outAnnotations, err
	}

	if user.Relationships == nil || user.Relationships.Profile == nil || user.Relationships.Profile.Data == nil {
		return nil, "", outAnnotations, fmt.Errorf("user {%s} profile is missing", userID)
	}

	userProfile := user.Relationships.Profile
	profileResource := &v2.Resource{
		Id: &v2.ResourceId{
			ResourceType: profileResourceType.Id,
			Resource:     strconv.Itoa(userProfile.Data.Id),
		},
	}

	grantResources = append(grantResources, grant.NewGrant(profileResource, profilePermissionName, resource))

	return grantResources, "", outAnnotations, nil
}

func (b *userBuilder) CreateAccountCapabilityDetails(_ context.Context) (*v2.CredentialDetailsAccountProvisioning, annotations.Annotations, error) {
	return &v2.CredentialDetailsAccountProvisioning{
		SupportedCredentialOptions: []v2.CapabilityDetailCredentialOption{
			v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
		},
		PreferredCredentialOption: v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
	}, nil, nil
}

func (b *userBuilder) CreateAccount(
	ctx context.Context,
	accountInfo *v2.AccountInfo,
	_ *v2.CredentialOptions,
) (connectorbuilder.CreateAccountResponse, []*v2.PlaintextData, annotations.Annotations, error) {
	outAnnotations := annotations.Annotations{}

	newUserInfo, err := createNewUserInfo(accountInfo)
	if err != nil {
		return nil, nil, annotations.Annotations{}, err
	}

	newUser, rateLimitData, err := b.client.CreateUser(ctx, *newUserInfo)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return nil, nil, outAnnotations, err
	}

	userResource, err := parseIntoUserResource(*newUser)
	if err != nil {
		return nil, nil, outAnnotations, err
	}

	caResponse := &v2.CreateAccountResponse_SuccessResult{
		Resource: userResource,
	}

	return caResponse, nil, outAnnotations, nil
}

func createNewUserInfo(accountInfo *v2.AccountInfo) (*client.NewUserBody, error) {
	pMap := accountInfo.Profile.AsMap()

	firstName, ok := pMap["first_name"].(string)
	if !ok || firstName == "" {
		return nil, fmt.Errorf("first_name is required")
	}

	lastName, ok := pMap["last_name"].(string)
	if !ok || lastName == "" {
		return nil, fmt.Errorf("last_name is required")
	}

	email, ok := pMap["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("email is required")
	}

	newUserInfo := &client.NewUserBody{
		Data: struct {
			Type       string                   `json:"type"` // The type should always be 'user'.
			Attributes client.NewUserAttributes `json:"attributes"`
		}{
			Type: "user",
			Attributes: client.NewUserAttributes{
				Email:     email,
				FirstName: firstName,
				LastName:  lastName,
			},
		},
	}

	return newUserInfo, nil
}

func (b *userBuilder) Delete(ctx context.Context, principal *v2.ResourceId) (annotations.Annotations, error) {
	outAnnotations := annotations.Annotations{}

	userID := principal.Resource

	rateLimitData, err := b.client.DisableUser(ctx, userID)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return outAnnotations, err
	}

	disabledUser, rateLimitData, err := b.client.GetUserByID(ctx, userID)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return outAnnotations, fmt.Errorf("error when deleting user. Error: %w", err)
	}

	if isActive(*disabledUser) {
		return outAnnotations, fmt.Errorf("error disabling user. User %s is not locked", userID)
	}

	return outAnnotations, nil
}

func isActive(user client.User) bool {
	return !user.Attributes.Locked
}

func parseIntoUserResource(user client.User) (*v2.Resource, error) {
	var userTraits []rs.UserTraitOption
	var userStatus v2.UserTrait_Status_Status
	primaryEmail := user.Attributes.Email

	profile := map[string]interface{}{
		"user_guid":  user.Attributes.UserGUID,
		"first_name": user.Attributes.FirstName,
		"last_name":  user.Attributes.LastName,
		"email":      primaryEmail,
		"username":   user.Attributes.Username,
		"title":      user.Attributes.Title,
	}

	if user.Attributes.Locked {
		userStatus = v2.UserTrait_Status_STATUS_DISABLED
	} else {
		userStatus = v2.UserTrait_Status_STATUS_ENABLED
	}

	userTraits = append(userTraits,
		rs.WithEmail(primaryEmail, true),
		rs.WithUserLogin(primaryEmail),
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
	)

	if user.Attributes.LastSignInAt != nil {
		userTraits = append(userTraits, rs.WithLastLogin(*user.Attributes.LastSignInAt))
	}

	ret, err := rs.NewUserResource(
		user.Attributes.Name,
		userResourceType,
		user.Id,
		userTraits,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newUserBuilder(c *client.OutreachClient) *userBuilder {
	return &userBuilder{
		client: c,
	}
}
