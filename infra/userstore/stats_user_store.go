package userstore

import (
	"context"

	"github.com/platelk/contactgraph/domain/models/userquery"
	"github.com/platelk/contactgraph/domain/models/users"
)

// InMemoryStats store stats about UserStore in memory
type InMemoryStats struct {
	store UserStore
	users uint
}

// WithStats implements Stats interface for a UserStore
func WithStats(store UserStore) *InMemoryStats {
	return &InMemoryStats{
		store: store,
	}
}

// Add implements UserStore
func (i *InMemoryStats) Add(ctx context.Context, user *users.User) (*users.User, error) {
	user, err := i.store.Add(ctx, user)
	if err == nil {
		i.users++
	}
	return user, err
}

// Get implements UserStore
func (i *InMemoryStats) Get(ctx context.Context, userID users.ID) (*users.User, error) {
	return i.store.Get(ctx, userID)
}

// Delete implements UserStore
func (i *InMemoryStats) Delete(ctx context.Context, user *users.User) (*users.User, error) {
	user, err := i.store.Delete(ctx, user)
	if err == nil {
		i.users--
	}
	return user, err
}

// Update implements UserStore
func (i *InMemoryStats) Update(ctx context.Context, user *users.User) (*users.User, error) {
	return i.store.Update(ctx, user)
}

// Query implements UserStore
func (i *InMemoryStats) Query() userquery.Queryer {
	return i.store.Query()
}

// Search implements UserStore
func (i *InMemoryStats) Search(ctx context.Context, q userquery.Queryer) ([]*users.User, error) {
	return i.store.Search(ctx, q)
}

// Users implements Stats
func (i *InMemoryStats) Users() uint {
	return i.users
}
