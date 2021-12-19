package contactstore

import (
	"context"

	"github.com/platelk/contactgraph/domain/models/contacts"
	"github.com/platelk/contactgraph/domain/models/users"
)

// InMemoryStats store stats about ContactStore in memory
type InMemoryStats struct {
	connection uint
	store      ContactStore
}

// WithStats implements Stats interface for a ContactStore
func WithStats(store ContactStore) *InMemoryStats {
	return &InMemoryStats{
		store: store,
	}
}

// Connect implements ContactStore
func (i *InMemoryStats) Connect(ctx context.Context, from, to users.ID) error {
	err := i.store.Connect(ctx, from, to)
	if err == nil {
		i.connection++
	}
	return err
}

// Lookup implements ContactStore
func (i *InMemoryStats) Lookup(ctx context.Context, userID users.ID) (contacts.List, error) {
	return i.store.Lookup(ctx, userID)
}

// ReverseLookup implements ContactStore
func (i *InMemoryStats) ReverseLookup(ctx context.Context, userID users.ID) (contacts.List, error) {
	return i.store.ReverseLookup(ctx, userID)
}

// Len implements ContactStore
func (i *InMemoryStats) Len() uint {
	return i.store.Len()
}

// Connection implements Stats
func (i *InMemoryStats) Connection() uint {
	return i.connection
}

// Users implements Stats
func (i *InMemoryStats) Users() uint {
	return i.store.Len()
}
