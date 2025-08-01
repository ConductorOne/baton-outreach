package connector

import (
	"context"

	"github.com/conductorone/baton-outreach/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const profilePermissionName = "assigned"

type profileBuilder struct {
	client *client.OutreachClient
}

func (b *profileBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return profileResourceType
}

func (b *profileBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var profileResources []*v2.Resource
	var nextPageToken string

	bag, nextPage, err := client.GetToken(pToken.Token, &v2.ResourceId{ResourceType: profileResourceType.Id})
	if err != nil {
		return nil, "", nil, err
	}

	profiles, nextPageLink, err := b.client.ListAllProfiles(ctx, nextPage)
	if err != nil {
		return nil, "", nil, err
	}

	for _, profile := range profiles {
		profileResource, err := parseIntoProfileResource(*profile)
		if err != nil {
			return nil, "", nil, err
		}

		profileResources = append(profileResources, profileResource)
	}

	if nextPageLink != "" {
		nextPageToken, err = bag.NextToken(nextPageLink)
		if err != nil {
			return nil, "", nil, err
		}
	}

	return profileResources, nextPageToken, nil, nil
}

func (b *profileBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var profileEntitlements []*v2.Entitlement

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDisplayName(resource.DisplayName),
		entitlement.WithDescription(resource.DisplayName),
	}

	profileEntitlements = append(profileEntitlements, entitlement.NewPermissionEntitlement(resource, profilePermissionName, assigmentOptions...))

	return profileEntitlements, "", nil, nil
}

// Grants function gets implemented on the users resource, since the users records have that data.
func (b *profileBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func parseIntoProfileResource(prof client.Profile) (*v2.Resource, error) {
	var resourceOptions []rs.ResourceOption

	profile := map[string]interface{}{
		"name":       prof.Attributes.Name,
		"created_at": prof.Attributes.CreatedAt,
		"updated_at": prof.Attributes.UpdatedAt,
		"is_admin":   prof.Attributes.IsAdmin,
	}

	profileTraits := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
	}

	ret, err := rs.NewRoleResource(
		prof.Attributes.Name,
		profileResourceType,
		prof.Id,
		profileTraits,
		resourceOptions...,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newProfileBuilder(c *client.OutreachClient) *profileBuilder {
	return &profileBuilder{
		client: c,
	}
}
