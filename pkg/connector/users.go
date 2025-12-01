package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-outreach/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/session"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userBuilder struct {
	client *client.OutreachClient
}

func (b *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

func (b *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, attr rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	var (
		userResources []*v2.Resource
		nextPageToken string
	)
	outAnnotations := annotations.Annotations{}
	token := attr.PageToken.Token

	bag, nextPage, err := client.GetToken(token, &v2.ResourceId{ResourceType: userResourceType.Id})
	if err != nil {
		return nil, nil, err
	}

	users, nextPageLink, rateLimitData, err := b.client.ListAllUsers(ctx, nextPage)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return nil, &rs.SyncOpResults{
			Annotations: outAnnotations,
		}, fmt.Errorf("error listing users: %w", err)
	}

	err = session.SetManyJSON(ctx, attr.Session, parseJSONCache(users))
	if err != nil {
		return nil, &rs.SyncOpResults{
			Annotations: outAnnotations,
		}, fmt.Errorf("error caching users in session: %w", err)
	}

	for _, user := range users {
		userResource, err := parseIntoUserResource(*user)
		if err != nil {
			return nil, &rs.SyncOpResults{
				Annotations: outAnnotations,
			}, fmt.Errorf("error getting user resource: %w", err)
		}

		userResources = append(userResources, userResource)
	}

	if nextPageLink != "" {
		nextPageToken, err = bag.NextToken(nextPageLink)
		if err != nil {
			return nil, &rs.SyncOpResults{
				Annotations: outAnnotations,
			}, fmt.Errorf("error parsing next page token for users: %w", err)
		}
	}

	return userResources, &rs.SyncOpResults{
		NextPageToken: nextPageToken,
		Annotations:   outAnnotations,
	}, nil
}

// Entitlements always returns an empty slice for users.
func (b *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

// Grants implements the Grants function for profiles resource.
func (b *userBuilder) Grants(ctx context.Context, resource *v2.Resource, attr rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	var grantResources []*v2.Grant
	outAnnotations := annotations.Annotations{}

	userID := resource.Id.Resource
	var user *client.User

	cachedUser, found, err := session.GetJSON[*client.User](ctx, attr.Session, userID)
	if err != nil {
		return nil, &rs.SyncOpResults{Annotations: outAnnotations}, err
	}

	if found {
		user = cachedUser
	} else {
		// user not found in cache
		u, rateLimitData, err := b.client.GetUserByID(ctx, userID)
		if err != nil {
			if rateLimitData != nil {
				outAnnotations.WithRateLimiting(rateLimitData)
			}

			return nil, &rs.SyncOpResults{Annotations: outAnnotations}, err
		}
		user = u
	}

	if user.Relationships == nil || user.Relationships.Profile == nil || user.Relationships.Profile.Data == nil {
		return nil, &rs.SyncOpResults{
			Annotations: outAnnotations,
		}, status.Errorf(codes.NotFound, "user {%s} profile is missing", userID)
	}

	userProfile := user.Relationships.Profile
	profileResource := &v2.Resource{
		Id: &v2.ResourceId{
			ResourceType: profileResourceType.Id,
			Resource:     strconv.Itoa(userProfile.Data.Id),
		},
	}

	grantResources = append(grantResources, grant.NewGrant(profileResource, profilePermissionName, resource))

	return grantResources, &rs.SyncOpResults{
		Annotations: outAnnotations,
	}, nil
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
	_ *v2.LocalCredentialOptions,
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

	// Login field contains the user email directly taken from C1.
	email := accountInfo.Login

	firstName, ok := pMap["first_name"].(string)
	if !ok || firstName == "" {
		return nil, fmt.Errorf("first_name is required")
	}

	lastName, ok := pMap["last_name"].(string)
	if !ok || lastName == "" {
		return nil, fmt.Errorf("last_name is required")
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
		rs.WithLastLogin(user.Attributes.LastSignInAt),
		rs.WithEmail(primaryEmail, true),
		rs.WithUserLogin(primaryEmail),
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
	)

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

func parseJSONCache(users []*client.User) map[string]*client.User {
	usersMap := make(map[string]*client.User)
	for _, user := range users {
		userIDStr := strconv.Itoa(user.Id)
		usersMap[userIDStr] = user
	}
	return usersMap
}
