package userstore

import (
	"context"
	"sync"

	"github.com/platelk/contactgraph/domain/models/userquery"
	"github.com/platelk/contactgraph/domain/models/users"
)

// InMemory is a user repo implementation which will store inmemory the users.
type InMemory struct {
	m           sync.RWMutex
	dataByID    map[users.ID]*users.User
	internalIdx uint64
}

// NewInMemory will initialise the store
func NewInMemory() *InMemory {
	return &InMemory{dataByID: make(map[users.ID]*users.User)}
}

// Add implements users.Adder
func (i *InMemory) Add(ctx context.Context, user *users.User) (*users.User, error) {
	i.m.Lock()
	defer i.m.Unlock()
	user.ID = users.ID(i.internalIdx)
	i.internalIdx++
	i.dataByID[user.ID] = &(*user) // nolint

	return user, nil
}

// Get will remove the user from the system
func (i *InMemory) Get(ctx context.Context, userID users.ID) (*users.User, error) {
	i.m.RLock()
	defer i.m.RUnlock()
	usr, ok := i.dataByID[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return usr, nil
}

// Delete will remove the user from the system
func (i *InMemory) Delete(ctx context.Context, user *users.User) (*users.User, error) {
	i.m.Lock()
	defer i.m.Unlock()
	usr, ok := i.dataByID[user.ID]
	if !ok {
		return nil, ErrNotFound
	}

	delete(i.dataByID, usr.ID)

	return usr, nil
}

// Update will update user with same ID to the new value
func (i *InMemory) Update(ctx context.Context, user *users.User) (*users.User, error) {
	i.m.Lock()
	defer i.m.Unlock()
	storedUser, ok := i.dataByID[user.ID]
	if !ok {
		return nil, ErrNotFound
	}

	if user.NickName != "" {
		storedUser.NickName = user.NickName
	}
	if user.PhoneNumber != "" {
		storedUser.PhoneNumber = user.PhoneNumber
	}
	return storedUser, nil
}

// Query will create a queryImproved to search users. implements users.Queryier
func (i *InMemory) Query() userquery.Queryer {
	return &query{}
}

// Search will execute the search on the user base
func (i *InMemory) Search(ctx context.Context, q userquery.Queryer) ([]*users.User, error) {
	i.m.RLock()
	defer i.m.RUnlock()
	sQuery, ok := q.(*query)
	if !ok {
		return nil, ErrQueryNotCompatible
	}
	var res []*users.User
	for _, usr := range i.dataByID {
		if sQuery.match(usr) {
			res = append(res, usr)
		}
	}
	return res, nil
}

// -- internal implementation --

type query struct {
	// TODO: change to string set
	ids          []users.ID
	nickName     []string
	phoneNumbers []users.PhoneNumber
}

func (q *query) ByID(id users.ID) userquery.Queryer {
	q.ids = append(q.ids, id)
	return q
}

func (q *query) ByNickName(nickName string) userquery.Queryer {
	q.nickName = append(q.nickName, nickName)
	return q
}

func (q *query) ByPhoneNumber(phoneNumber users.PhoneNumber) userquery.Queryer {
	q.phoneNumbers = append(q.phoneNumbers, phoneNumber)
	return q
}

func (q *query) match(u *users.User) bool {
	for _, id := range q.ids {
		if id == u.ID {
			return true
		}
	}
	for _, nickName := range q.nickName {
		if nickName == u.NickName {
			return true
		}
	}
	for _, phoneNumber := range q.phoneNumbers {
		if phoneNumber == u.PhoneNumber {
			return true
		}
	}
	return false
}
