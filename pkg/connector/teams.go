package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-outreach/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const teamPermissionName = "member"

type teamBuilder struct {
	client *client.OutreachClient
}

func (b *teamBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return teamResourceType
}

func (b *teamBuilder) List(ctx context.Context, _ *v2.ResourceId, attr rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	var (
		teamResources []*v2.Resource
		nextPageToken string
	)
	outAnnotations := annotations.Annotations{}
	token := attr.PageToken.Token

	bag, nextPage, err := client.GetToken(token, &v2.ResourceId{ResourceType: teamResourceType.Id})
	if err != nil {
		return nil, nil, err
	}

	teams, nextPageLink, rateLimitData, err := b.client.ListAllTeams(ctx, nextPage)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return nil, &rs.SyncOpResults{
			Annotations: outAnnotations,
		}, err
	}

	for _, team := range teams {
		teamResource, err := parseIntoTeamResource(*team)
		if err != nil {
			return nil, &rs.SyncOpResults{
				Annotations: outAnnotations,
			}, err
		}

		teamResources = append(teamResources, teamResource)
	}

	if nextPageLink != "" {
		nextPageToken, err = bag.NextToken(nextPageLink)
		if err != nil {
			return nil, &rs.SyncOpResults{
				Annotations: outAnnotations,
			}, err
		}
	}

	return teamResources, &rs.SyncOpResults{
		Annotations:   outAnnotations,
		NextPageToken: nextPageToken,
	}, nil
}

func (b *teamBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	var outAnnotations annotations.Annotations

	displayName := fmt.Sprintf("Member of %s", resource.DisplayName)
	description := fmt.Sprintf("Member of the %s team.", resource.DisplayName)

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDisplayName(displayName),
		entitlement.WithDescription(description),
	}

	return []*v2.Entitlement{entitlement.NewAssignmentEntitlement(resource, teamPermissionName, assigmentOptions...)},
		&rs.SyncOpResults{
			Annotations: outAnnotations,
		}, nil
}

func (b *teamBuilder) Grants(ctx context.Context, resource *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	var grantResources []*v2.Grant
	outAnnotations := annotations.Annotations{}
	logger := ctxzap.Extract(ctx)

	teamID := resource.Id.Resource

	teamDetails, rateLimitData, err := b.client.GetTeamByID(ctx, teamID)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return nil, &rs.SyncOpResults{
			Annotations: outAnnotations,
		}, err
	}

	if teamDetails.Relationships == nil || teamDetails.Relationships.Users == nil || teamDetails.Relationships.Users.Data == nil {
		logger.Warn(fmt.Sprintf("the team {%s} does not have any members", teamID))
		return nil, &rs.SyncOpResults{
			Annotations: outAnnotations,
		}, nil
	}

	teamMembers := *teamDetails.Relationships.Users.Data
	for _, member := range teamMembers {
		userResource := &v2.Resource{
			Id: &v2.ResourceId{
				ResourceType: userResourceType.Id,
				Resource:     strconv.Itoa(member.Id),
			},
		}

		grantResources = append(grantResources, grant.NewGrant(resource, teamPermissionName, userResource))
	}

	return grantResources, &rs.SyncOpResults{
		Annotations: outAnnotations,
	}, nil
}

func (b *teamBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	var teamMembers []client.DataDetailPair
	outAnnotations := annotations.Annotations{}

	teamID := entitlement.Resource.Id.Resource
	userID, err := strconv.Atoi(principal.Id.Resource)
	if err != nil {
		return outAnnotations, err
	}

	teamDetails, rateLimitData, err := b.client.GetTeamByID(ctx, teamID)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return outAnnotations, err
	}

	if teamDetails.Relationships != nil && teamDetails.Relationships.Users != nil && teamDetails.Relationships.Users.Data != nil {
		teamMembers = *teamDetails.Relationships.Users.Data
	}

	for _, member := range teamMembers {
		if member.Id == userID {
			// It doesn't fail when "re-adding" an existing user to a Team, but to avoid the unnecessary request, I added this validation and the annotation.
			outAnnotations.Update(&v2.GrantAlreadyExists{})
			return outAnnotations, nil
		}
	}
	teamMembers = append(teamMembers, client.DataDetailPair{
		Id:   userID,
		Type: "user",
	})

	rateLimitData, err = b.client.UpdateTeamMembers(ctx, teamID, teamMembers)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return outAnnotations, err
	}

	return outAnnotations, nil
}

func (b *teamBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	outAnnotations := annotations.Annotations{}
	updatedTeamMembers := make([]client.DataDetailPair, 0)

	teamID := grant.Entitlement.Resource.Id.Resource
	userID, err := strconv.Atoi(grant.Principal.Id.Resource)
	if err != nil {
		return nil, err
	}

	teamDetails, rateLimitData, err := b.client.GetTeamByID(ctx, teamID)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return outAnnotations, err
	}

	if teamDetails.Relationships == nil || teamDetails.Relationships.Users == nil || teamDetails.Relationships.Users.Data == nil {
		return nil, fmt.Errorf("revoke tried on the team {%s} but the members list was not accessible", teamID)
	}

	teamMembers := *teamDetails.Relationships.Users.Data
	for _, member := range teamMembers {
		if member.Id == userID {
			continue
		}

		updatedTeamMembers = append(updatedTeamMembers, member)
	}

	rateLimitData, err = b.client.UpdateTeamMembers(ctx, teamID, updatedTeamMembers)
	if err != nil {
		if rateLimitData != nil {
			outAnnotations.WithRateLimiting(rateLimitData)
		}
		return outAnnotations, err
	}

	return outAnnotations, nil
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
