package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// ProfileResourceTypeID is the resource type ID for profiles, exported so it
// can be referenced when checking the sync filter (WillSyncResourceType).
const ProfileResourceTypeID = "profile"

// The user resource type is for all user objects from the database.
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

var profileResourceType = &v2.ResourceType{
	Id:          ProfileResourceTypeID,
	DisplayName: "Profile",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_ROLE},
}
