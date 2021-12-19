package usecases

import (
	"context"
	"fmt"

	"github.com/platelk/contactgraph/domain/models/contacts"
	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/infra/logger"
)

type contactReverseLookuper interface {
	ReverseLookup(ctx context.Context, user users.ID) (contacts.List, error)
}

type ReverseLookupContactReq struct {
	User users.ID `json:"user_id"`
}

type ReverseLookupContactResp struct {
	Contacts []*users.User
}

type ReverseLookupContact func(ctx context.Context, req *ReverseLookupContactReq) (*ReverseLookupContactResp, error)

func SetupReverseLookupContact(log logger.Logger, lookuper contactReverseLookuper, userStore userGetter) ReverseLookupContact {
	return func(ctx context.Context, req *ReverseLookupContactReq) (*ReverseLookupContactResp, error) {
		contactIDs, err := lookuper.ReverseLookup(ctx, req.User)
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
		return &ReverseLookupContactResp{
			Contacts: userContacts,
		}, nil
	}
}
