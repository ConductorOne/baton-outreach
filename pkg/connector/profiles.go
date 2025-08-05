package connector

import (
	"context"
	"strconv"

	"github.com/conductorone/baton-outreach/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const profilePermissionName = "assigned"
const defaultProfileID = 2 // This is the ID for the 'Default' profile, a system-provided profile.

type profileBuilder struct {
	client *client.OutreachClient
}

func (b *profileBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return profileResourceType
}

func (b *profileBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var (
		profileResources []*v2.Resource
		nextPageToken    string
	)
	outAnnotations := annotations.Annotations{}

	bag, nextPage, err := client.GetToken(pToken.Token, &v2.ResourceId{ResourceType: profileResourceType.Id})
	if err != nil {
		return nil, "", nil, err
	}

	profiles, nextPageLink, rateLimitData, err := b.client.ListAllProfiles(ctx, nextPage)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return nil, "", outAnnotations, err
	}

	for _, profile := range profiles {
		profileResource, err := parseIntoProfileResource(*profile)
		if err != nil {
			return nil, "", outAnnotations, err
		}

		profileResources = append(profileResources, profileResource)
	}

	if nextPageLink != "" {
		nextPageToken, err = bag.NextToken(nextPageLink)
		if err != nil {
			return nil, "", outAnnotations, err
		}
	}

	return profileResources, nextPageToken, outAnnotations, nil
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

func (b *profileBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	outAnnotations := annotations.Annotations{}
	profileID, err := strconv.Atoi(entitlement.Resource.Id.Resource)
	if err != nil {
		return outAnnotations, err
	}
	userID := principal.Id.Resource

	rateLimitData, err := b.client.UpdateUserProfile(ctx, userID, profileID)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return outAnnotations, err
	}

	return outAnnotations, nil
}

func (b *profileBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	outAnnotations := annotations.Annotations{}
	profileID := defaultProfileID
	userID := grant.Principal.Id.Resource

	rateLimitData, err := b.client.UpdateUserProfile(ctx, userID, profileID)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return outAnnotations, err
	}

	return outAnnotations, nil
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
