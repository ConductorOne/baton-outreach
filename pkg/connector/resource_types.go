package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// The user resource type is for all user objects from the database.
// newUserBuilder clones this and adds SkipEntitlements, or
// SkipEntitlementsAndGrants when profile isn't synced.
var userResourceType = &v2.ResourceType{
	Id:          "user",
	DisplayName: "User",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_USER},
}

var teamResourceType = &v2.ResourceType{
	Id:          "team",
	DisplayName: "Team",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_GROUP},
}

// ProfileResourceTypeID is referenced when gating cross-type profile grants.
const ProfileResourceTypeID = "profile"

var profileResourceType = &v2.ResourceType{
	Id:          ProfileResourceTypeID,
	DisplayName: "Profile",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_ROLE},
}
