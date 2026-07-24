package connector

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/conductorone/baton-outreach/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	eopt "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sdk/pkg/types/sessions"
	"github.com/stretchr/testify/require"
)

// fakeSessionStore is a minimal in-memory implementation of
// sessions.SessionStore, sufficient for exercising session.GetJSON /
// session.SetManyJSON without needing a real dotc1z-backed store.
type fakeSessionStore struct {
	mu   sync.Mutex
	data map[string][]byte
}

func newFakeSessionStore() *fakeSessionStore {
	return &fakeSessionStore{data: make(map[string][]byte)}
}

func (f *fakeSessionStore) Get(_ context.Context, key string, _ ...sessions.SessionStoreOption) ([]byte, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.data[key]
	return v, ok, nil
}

func (f *fakeSessionStore) GetMany(_ context.Context, keys []string, _ ...sessions.SessionStoreOption) (map[string][]byte, []string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	found := make(map[string][]byte)
	var missing []string
	for _, k := range keys {
		if v, ok := f.data[k]; ok {
			found[k] = v
		} else {
			missing = append(missing, k)
		}
	}
	return found, missing, nil
}

func (f *fakeSessionStore) Set(_ context.Context, key string, value []byte, _ ...sessions.SessionStoreOption) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.data[key] = value
	return nil
}

func (f *fakeSessionStore) SetMany(_ context.Context, values map[string][]byte, _ ...sessions.SessionStoreOption) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for k, v := range values {
		f.data[k] = v
	}
	return nil
}

func (f *fakeSessionStore) Delete(_ context.Context, key string, _ ...sessions.SessionStoreOption) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.data, key)
	return nil
}

func (f *fakeSessionStore) Clear(_ context.Context, _ ...sessions.SessionStoreOption) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.data = make(map[string][]byte)
	return nil
}

func (f *fakeSessionStore) GetAll(_ context.Context, _ string, _ ...sessions.SessionStoreOption) (map[string][]byte, string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make(map[string][]byte, len(f.data))
	for k, v := range f.data {
		out[k] = v
	}
	return out, "", nil
}

func TestUserBuilder_Grants_ProfileSynced(t *testing.T) {
	ctx := context.Background()

	const userID = 42
	const profileID = 7

	user := &client.User{
		Id:   userID,
		Type: "user",
		Relationships: &client.UserRelationships{
			Profile: &struct {
				Data *client.DataDetailPair `json:"data,omitempty"`
			}{
				Data: &client.DataDetailPair{Id: profileID, Type: "profile"},
			},
		},
	}

	userBytes, err := json.Marshal(user)
	require.NoError(t, err)

	store := newFakeSessionStore()
	require.NoError(t, store.Set(ctx, "42", userBytes))

	b := newUserBuilder(nil, true)

	resource := &v2.Resource{
		Id: &v2.ResourceId{ResourceType: userResourceType.Id, Resource: "42"},
	}

	grants, results, err := b.Grants(ctx, resource, rs.SyncOpAttrs{Session: store})
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, grants, 1)

	g := grants[0]
	require.Equal(t, profileResourceType.Id, g.Entitlement.Resource.Id.ResourceType)

	expectedProfileResource := &v2.Resource{
		Id: &v2.ResourceId{
			ResourceType: profileResourceType.Id,
			Resource:     "7",
		},
	}
	require.Equal(t, eopt.NewEntitlementID(expectedProfileResource, profilePermissionName), g.Entitlement.Id)
}

// TestUserBuilder_Grants_ProfileNotSynced verifies that Grants() behaves
// identically regardless of syncProfiles - the gate now lives entirely at
// the resource-type annotation level (see ResourceType tests below), not
// inside Grants itself. With syncProfiles false, this still exercises the
// real emission path and expects the same single profile grant as the
// syncProfiles-true case.
func TestUserBuilder_Grants_ProfileNotSynced(t *testing.T) {
	ctx := context.Background()

	const userID = 42
	const profileID = 7

	user := &client.User{
		Id:   userID,
		Type: "user",
		Relationships: &client.UserRelationships{
			Profile: &struct {
				Data *client.DataDetailPair `json:"data,omitempty"`
			}{
				Data: &client.DataDetailPair{Id: profileID, Type: "profile"},
			},
		},
	}

	userBytes, err := json.Marshal(user)
	require.NoError(t, err)

	store := newFakeSessionStore()
	require.NoError(t, store.Set(ctx, "42", userBytes))

	b := newUserBuilder(nil, false)

	resource := &v2.Resource{
		Id: &v2.ResourceId{ResourceType: userResourceType.Id, Resource: "42"},
	}

	grants, results, err := b.Grants(ctx, resource, rs.SyncOpAttrs{Session: store})
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, grants, 1)

	g := grants[0]
	require.Equal(t, profileResourceType.Id, g.Entitlement.Resource.Id.ResourceType)

	expectedProfileResource := &v2.Resource{
		Id: &v2.ResourceId{
			ResourceType: profileResourceType.Id,
			Resource:     "7",
		},
	}
	require.Equal(t, eopt.NewEntitlementID(expectedProfileResource, profilePermissionName), g.Entitlement.Id)
}

func TestUserBuilder_ResourceType_ProfileSynced(t *testing.T) {
	preCount := len(userResourceType.GetAnnotations())

	b := newUserBuilder(nil, true)
	rt := b.ResourceType(context.Background())

	require.Equal(t, userResourceType.Id, rt.Id)
	rtAnnotations := annotations.Annotations(rt.GetAnnotations())
	require.True(t, rtAnnotations.Contains(&v2.SkipEntitlements{}))
	require.False(t, rtAnnotations.Contains(&v2.SkipEntitlementsAndGrants{}))

	// The package-level var must never be mutated by ResourceType().
	require.Equal(t, preCount, len(userResourceType.GetAnnotations()))
	require.Zero(t, len(userResourceType.GetAnnotations()))
}

func TestUserBuilder_ResourceType_ProfileNotSynced(t *testing.T) {
	preCount := len(userResourceType.GetAnnotations())

	b := newUserBuilder(nil, false)
	rt := b.ResourceType(context.Background())

	require.Equal(t, userResourceType.Id, rt.Id)
	rtAnnotations := annotations.Annotations(rt.GetAnnotations())
	require.True(t, rtAnnotations.Contains(&v2.SkipEntitlementsAndGrants{}))

	// The package-level var must never be mutated by ResourceType().
	require.Equal(t, preCount, len(userResourceType.GetAnnotations()))
	require.Zero(t, len(userResourceType.GetAnnotations()))
}
