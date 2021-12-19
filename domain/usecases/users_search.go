package usecases

import (
	"context"
	"fmt"

	"github.com/platelk/contactgraph/domain/models/userquery"
	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/infra/logger"
)

// UserSearchReq contains the required parameters to search users
type UserSearchReq struct {
	IDs          []users.ID
	Emails       []string
	FirstName    []string
	LastName     []string
	NickName     []string
	PhoneNumbers []string
}

// UserSearchResp contains the field which will be returned on successful user search
type UserSearchResp struct {
	Users []*users.User `json:"users"`
}

// UserSearcher will allow searching users based on different criteria
type UserSearcher interface {
	Query() userquery.Queryer
	Search(ctx context.Context, query userquery.Queryer) ([]*users.User, error)
}

// SearchUser define the function which will search for users in the system
type SearchUser func(ctx context.Context, req *UserSearchReq) (*UserSearchResp, error)

// SetupSearchUser will return a configured CreateUser function which can be used later
func SetupSearchUser(log logger.Logger, repo UserSearcher) SearchUser {
	log = log.With().Str("usecase", "user_search").Logger()
	return searchUser(repo)
}

func searchUser(repo UserSearcher) SearchUser {
	return func(ctx context.Context, req *UserSearchReq) (*UserSearchResp, error) {
		qBuilder := repo.Query()

		for _, id := range req.IDs {
			qBuilder = qBuilder.ByID(id)
		}
		for _, nickName := range req.NickName {
			qBuilder = qBuilder.ByNickName(nickName)
		}
		for _, phoneNumber := range req.PhoneNumbers {
			qBuilder = qBuilder.ByPhoneNumber(users.ParsePhoneNumber(phoneNumber))
		}

		users, err := repo.Search(ctx, qBuilder)
		if err != nil {
			return nil, fmt.Errorf("can't perform search: %w", err)
		}

		return &UserSearchResp{Users: users}, nil
	}
}
