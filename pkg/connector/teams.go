package connector

import (
	"context"

	"github.com/conductorone/baton-outreach/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type teamBuilder struct {
	client *client.OutreachClient
}

func (b *teamBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return teamResourceType
}

func (b *teamBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var teamResources []*v2.Resource
	var nextPageToken string

	bag, nextPage, err := client.GetToken(pToken.Token, &v2.ResourceId{ResourceType: teamResourceType.Id})
	if err != nil {
		return nil, "", nil, err
	}

	teams, nextPageLink, err := b.client.ListAllTeams(ctx, nextPage)
	if err != nil {
		return nil, "", nil, err
	}

	for _, team := range teams {
		teamResource, err := parseIntoTeamResource(*team)
		if err != nil {
			return nil, "", nil, err
		}

		teamResources = append(teamResources, teamResource)
	}

	if nextPageLink != "" {
		nextPageToken, err = bag.NextToken(nextPageLink)
		if err != nil {
			return nil, "", nil, err
		}
	}

	return teamResources, nextPageToken, nil, nil
}

func (b *teamBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (b *teamBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func parseIntoTeamResource(team client.Team) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"name":       team.Attributes.Name,
		"created_at": team.Attributes.CreatedAt,
		"updated_at": team.Attributes.UpdatedAt,
	}

	groupTraits := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	ret, err := rs.NewGroupResource(
		team.Attributes.Name,
		teamResourceType,
		team.Id,
		groupTraits,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newTeamBuilder(c *client.OutreachClient) *teamBuilder {
	return &teamBuilder{
		client: c,
	}
}
