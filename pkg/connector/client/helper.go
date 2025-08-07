package client

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
)

func GetToken(token string, resourceID *v2.ResourceId) (*pagination.Bag, string, error) {
	bag := &pagination.Bag{}
	err := bag.Unmarshal(token)
	if err != nil {
		return nil, "", err
	}

	if bag.Current() == nil {
		bag.Push(pagination.PageState{
			ResourceTypeID: resourceID.ResourceType,
			ResourceID:     resourceID.Resource,
		})
	}

	return bag, bag.PageToken(), nil
}
