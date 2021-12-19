package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/dongri/phonenumber"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/infra/logger"
)

// ErrInvalidUser is returned if one of the provided field isn't validated
var ErrInvalidUser = errors.New("provided user isn't valid")

// UserCreateReq contains the required parameters to create a new user
type UserCreateReq struct {
	NickName    string `json:"nick_name"`
	PhoneNumber string `json:"phone_number"`
}

// UserCreateResp contains the field which will be returned on successful user creation
type UserCreateResp struct {
	User *users.User `json:"user"`
}

// UserAdder will save a new user in the system and generate an id for the User
type UserAdder interface {
	Add(ctx context.Context, user *users.User) (*users.User, error)
}

// CreateUser define the function which will create a user in the system
type CreateUser func(ctx context.Context, req *UserCreateReq) (*UserCreateResp, error)

// SetupCreateUser will return a configured CreateUser function which can be used later
func SetupCreateUser(log logger.Logger, repo UserAdder) CreateUser {
	log = log.With().Str("usecases", "user_create").Logger()
	return validateCreate(createUser(repo))
}

func createUser(repo UserAdder) CreateUser {
	return func(ctx context.Context, req *UserCreateReq) (*UserCreateResp, error) {
		newUser, err := repo.Add(ctx, &users.User{
			NickName:    req.NickName,
			PhoneNumber: users.PhoneNumber(req.PhoneNumber),
		})
		if err != nil {
			return nil, fmt.Errorf("can't save new user: %w", err)
		}
		return &UserCreateResp{User: newUser}, nil
	}
}

func validateCreate(createFunc CreateUser) CreateUser {
	return func(ctx context.Context, req *UserCreateReq) (*UserCreateResp, error) {
		err := validateUser(&users.User{
			NickName: req.NickName,
		})
		if err != nil {
			return nil, fmt.Errorf("can't validate user: %s: %w", err.Error(), ErrInvalidUser)
		}

		parsedPhone := phonenumber.ParseWithLandLine(req.PhoneNumber, "")
		if parsedPhone == "" {
			return nil, fmt.Errorf("can't parse user phone number %v for country: %v", req.PhoneNumber, "unknown")
		}
		req.PhoneNumber = parsedPhone
		return createFunc(ctx, req)
	}
}
