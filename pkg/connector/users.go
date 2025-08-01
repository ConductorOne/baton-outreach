package connector

import (
	"context"

	"github.com/conductorone/baton-outreach/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client *client.OutreachClient
}

func (b *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

func (b *userBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var userResources []*v2.Resource
	var nextPageToken string

	bag, nextPage, err := client.GetToken(pToken.Token, &v2.ResourceId{ResourceType: userResourceType.Id})
	if err != nil {
		return nil, "", nil, err
	}

	users, nextPageLink, err := b.client.ListAllUsers(ctx, nextPage)
	if err != nil {
		return nil, "", nil, err
	}

	for _, user := range users {
		userResource, err := parseIntoUserResource(*user)
		if err != nil {
			return nil, "", nil, err
		}

		userResources = append(userResources, userResource)
	}

	if nextPageLink != "" {
		nextPageToken, err = bag.NextToken(nextPageLink)
		if err != nil {
			return nil, "", nil, err
		}
	}

	return userResources, nextPageToken, nil, nil
}

// Entitlements always returns an empty slice for users.
func (b *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants implements the Grants function for roles resource.
func (b *userBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var grantResources []*v2.Grant
	//userID := resource.Id.Resource
	//
	//user, err := b.client.GetUserByID(ctx, userID)
	//if err != nil {
	//	return nil, "", nil, err
	//}
	//
	//if user.Relationships == nil || user.Relationships.Role == nil {
	//	return nil, "", nil, fmt.Errorf("user {%s} role is missing", userID)
	//}
	//
	//userRole := user.Relationships.Role
	//roleResource := &v2.Resource{
	//	Id: &v2.ResourceId{
	//		ResourceType: roleResourceType.Id,
	//		Resource:     userRole.Data.Id,
	//	},
	//}
	//
	//grantResources = append(grantResources, grant.NewGrant(roleResource, rolePermissionName, resource))

	return grantResources, "", nil, nil
}

func parseIntoUserResource(user client.User) (*v2.Resource, error) {
	var userTraits []rs.UserTraitOption
	userStatus := v2.UserTrait_Status_STATUS_DISABLED
	primaryEmail := user.Attributes.Email

	profile := map[string]interface{}{
		"user_guid":  user.Attributes.UserGuid,
		"first_name": user.Attributes.FirstName,
		"last_name":  user.Attributes.LastName,
		"email":      primaryEmail,
		"username":   user.Attributes.Username,
		"title":      user.Attributes.Title,
	}

	userTraits = append(userTraits,
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
		rs.WithEmail(primaryEmail, true),
		rs.WithUserLogin(primaryEmail),
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
