package usecases

import (
	"context"
	"fmt"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/infra/logger"
)

// UserUpdateReq contains the required parameters to UpdateUser a new user
type UserUpdateReq struct {
	ID          users.ID `json:"id"`
	NickName    string   `json:"nick_name"`
	PhoneNumber string   `json:"phone_number"`
}

// UserUpdateResp contains the field which will be returned on successful user update
type UserUpdateResp struct {
	User *users.User `json:"user"`
}

// UserUpdater will update a user in the system
type UserUpdater interface {
	Update(ctx context.Context, user *users.User) (*users.User, error)
}

// UpdateUser define the function which will UpdateUser a user in the system
type UpdateUser func(ctx context.Context, req *UserUpdateReq) (*UserUpdateResp, error)

// SetupUpdateUser will return a configured UpdateUser function which can be used later
func SetupUpdateUser(log logger.Logger, repo UserUpdater) UpdateUser {
	log = log.With().Str("usecase", "user_update").Logger()
	return validateUpdate(log, updateUser(repo))
}

func updateUser(repo UserUpdater) UpdateUser {
	return func(ctx context.Context, req *UserUpdateReq) (*UserUpdateResp, error) {
		newUser, err := repo.Update(ctx, &users.User{
			ID:          req.ID,
			NickName:    req.NickName,
			PhoneNumber: users.ParsePhoneNumber(req.PhoneNumber),
		})
		if err != nil {
			return nil, fmt.Errorf("can't save new user: %w", err)
		}
		return &UserUpdateResp{User: newUser}, nil
	}
}

func validateUpdate(log logger.Logger, updateFunc UpdateUser) UpdateUser {
	return func(ctx context.Context, req *UserUpdateReq) (*UserUpdateResp, error) {
		log.Debug().Interface("req", req).Msg("receive update")
		if err := validateNickName(req.NickName); req.NickName != "" && err != nil {
			return nil, fmt.Errorf("can't validate user: %s: %w", err.Error(), ErrInvalidUser)
		}
		if err := validatePhoneNumber(req.PhoneNumber); req.PhoneNumber != "" && err != nil {
			return nil, fmt.Errorf("can't validate user: %s: %w", err.Error(), ErrInvalidUser)
		}
		return updateFunc(ctx, req)
	}
}
