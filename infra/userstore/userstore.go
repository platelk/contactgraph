package userstore

import (
	"context"
	"errors"

	"github.com/platelk/contactgraph/domain/models/userquery"
	"github.com/platelk/contactgraph/domain/models/users"
)

// ErrAlreadyExist is returned if the email is already present in the store
var ErrAlreadyExist = errors.New("user already exist")

// ErrNotFound is returned if the id of the user isn't found in the store
var ErrNotFound = errors.New("user not found")

// ErrQueryNotCompatible is returned if the email is already present in the store
var ErrQueryNotCompatible = errors.New("the provided queryImproved is not compatible")

// UserStore define the interface to store and retrieve users of the systems.
type UserStore interface {
	Add(ctx context.Context, user *users.User) (*users.User, error)
	Get(ctx context.Context, userID users.ID) (*users.User, error)
	Delete(ctx context.Context, user *users.User) (*users.User, error)
	Update(ctx context.Context, user *users.User) (*users.User, error)
	Query() userquery.Queryer
	Search(ctx context.Context, q userquery.Queryer) ([]*users.User, error)
}

// Stats keep statistic about UserStore
type Stats interface {
	Users() uint
}
