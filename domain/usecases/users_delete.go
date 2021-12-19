package usecases

import (
	"context"
	"fmt"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/infra/logger"
)

// UserDeleteReq contains the required parameters to delete a new user
type UserDeleteReq struct {
	ID users.ID `json:"id"`
}

// UserDeleteResp contains the field which will be returned on successful user creation
type UserDeleteResp struct {
	User *users.User `json:"user"`
}

// UserDeleter will delete a user in the system
type UserDeleter interface {
	Delete(ctx context.Context, user *users.User) (*users.User, error)
}

// DeleteUser define the function which will delete a user in the system
type DeleteUser func(ctx context.Context, req *UserDeleteReq) (*UserDeleteResp, error)

// SetupDeleteUser will return a configured DeleteUser function which can be used later
func SetupDeleteUser(log logger.Logger, repo UserDeleter) DeleteUser {
	log = log.With().Str("usecase", "user_delete").Logger()
	return deleteUser(repo)
}

func deleteUser(repo UserDeleter) DeleteUser {
	return func(ctx context.Context, req *UserDeleteReq) (*UserDeleteResp, error) {
		newUser, err := repo.Delete(ctx, &users.User{ID: req.ID})
		if err != nil {
			return nil, fmt.Errorf("can't delete user: %w", err)
		}
		return &UserDeleteResp{User: newUser}, nil
	}
}
