package connector

import (
	"context"
	"strconv"

	"github.com/conductorone/baton-outreach/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleBuilder struct {
	client *client.OutreachClient
}

func (b *roleBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return roleResourceType
}

func (b *roleBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var roleResources []*v2.Resource
	var nextPageToken string

	bag, nextPage, err := client.GetToken(pToken.Token, &v2.ResourceId{ResourceType: roleResourceType.Id})
	if err != nil {
		return nil, "", nil, err
	}

	roles, nextPageLink, err := b.client.ListAllRoles(ctx, nextPage)
	if err != nil {
		return nil, "", nil, err
	}

	for _, role := range roles {
		roleResource, err := parseIntoRoleResource(*role)
		if err != nil {
			return nil, "", nil, err
		}

		roleResources = append(roleResources, roleResource)
	}

	if nextPageLink != "" {
		nextPageToken, err = bag.NextToken(nextPageLink)
		if err != nil {
			return nil, "", nil, err
		}
	}

	return roleResources, nextPageToken, nil, nil
}

func (b *roleBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (b *roleBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func parseIntoRoleResource(role client.Role) (*v2.Resource, error) {
	var resourceOptions []rs.ResourceOption
	profile := map[string]interface{}{
		"name":       role.Attributes.Name,
		"created_at": role.Attributes.CreatedAt,
		"updated_at": role.Attributes.UpdatedAt,
	}

	roleTraits := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
	}

	if role.Relationships != nil && role.Relationships.ParentRole != nil {
		parentRoleData := role.Relationships.ParentRole
		if parentRoleData != nil {
			parentRoleID := &v2.ResourceId{
				ResourceType: roleResourceType.Id,
				Resource:     strconv.Itoa(role.Relationships.ParentRole.Data.Id),
			}

			resourceOptions = append(resourceOptions, rs.WithParentResourceID(parentRoleID))
		}
	}
	
	ret, err := rs.NewRoleResource(
		role.Attributes.Name,
		roleResourceType,
		role.Id,
		roleTraits,
		resourceOptions...,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newRoleBuilder(c *client.OutreachClient) *roleBuilder {
	return &roleBuilder{
		client: c,
	}
}
