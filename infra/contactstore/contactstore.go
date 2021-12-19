package contactstore

import (
	"context"
	"errors"

	"github.com/platelk/contactgraph/domain/models/contacts"
	"github.com/platelk/contactgraph/domain/models/users"
)

// ErrNotFound is an error returned if a connection or a user is not found.
var ErrNotFound = errors.New("contact not found")

// ContactStore define the common functionality expected from a store which hold the connection between users.
type ContactStore interface {
	Connect(ctx context.Context, from, to users.ID) error
	Lookup(ctx context.Context, userID users.ID) (contacts.List, error)
	ReverseLookup(ctx context.Context, userID users.ID) (contacts.List, error)
	Len() uint
}

type IteratorFunc func(userID users.ID, friends contacts.List)

type IterableStore interface {
	ContactStore
	Iterate(ctx context.Context, it IteratorFunc) error
}

// Stats keep statistic about ContactStore
type Stats interface {
	Connection() uint
	Users() uint
}
