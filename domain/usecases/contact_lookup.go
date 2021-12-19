package usecases

import (
	"context"
	"fmt"

	"github.com/platelk/contactgraph/domain/models/contacts"
	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/infra/logger"
)

type contactLookuper interface {
	Lookup(ctx context.Context, user users.ID) (contacts.List, error)
}

type userGetter interface {
	Get(ctx context.Context, userID users.ID) (*users.User, error)
}

type LookupContactReq struct {
	User users.ID `json:"user_id"`
}

type LookupContactResp struct {
	Contacts []*users.User
}

type LookupContact func(ctx context.Context, req *LookupContactReq) (*LookupContactResp, error)

func SetupLookupContact(log logger.Logger, lookuper contactLookuper, userStore userGetter) LookupContact {
	return func(ctx context.Context, req *LookupContactReq) (*LookupContactResp, error) {
		contactIDs, err := lookuper.Lookup(ctx, req.User)
		if err != nil {
			return nil, fmt.Errorf("can't lookup user [%v]: %w", req.User, err)
		}
		userContacts := make([]*users.User, len(contactIDs))
		var i int
		for contactID := range contactIDs {
			usr, err := userStore.Get(ctx, contactID)
			if err != nil {
				return nil, fmt.Errorf("can't retrieve user [%v] information: %w", contactID, err)
			}
			userContacts[i] = usr
			i++
		}
		return &LookupContactResp{
			Contacts: userContacts,
		}, nil
	}
}
