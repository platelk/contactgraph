package userstore

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/platelk/contactgraph/domain/models/userquery"
	"github.com/platelk/contactgraph/domain/models/users"
)

type internalUser struct {
	id       users.ID
	nickName string
	phone    uint64
}

// InMemoryImproved is a user repo implementation which will store InMemoryImproved the users.
type InMemoryImproved struct {
	m           sync.RWMutex
	dataByID    map[users.ID]*internalUser
	internalIdx uint64
}

// NewInMemoryImproved will initialise the store
func NewInMemoryImproved() *InMemoryImproved {
	return &InMemoryImproved{dataByID: make(map[users.ID]*internalUser)}
}

// Add implements users.Adder
func (i *InMemoryImproved) Add(ctx context.Context, user *users.User) (*users.User, error) {
	i.m.Lock()
	defer i.m.Unlock()
	user.ID = users.ID(i.internalIdx)
	i.internalIdx++
	phone, _ := strconv.ParseUint(string(user.PhoneNumber), 10, 64)
	i.dataByID[user.ID] = &internalUser{
		id:       user.ID,
		nickName: user.NickName,
		phone:    phone,
	}

	return user, nil
}

// Get will remove the user from the system
func (i *InMemoryImproved) Get(ctx context.Context, userID users.ID) (*users.User, error) {
	i.m.RLock()
	defer i.m.RUnlock()
	usr, ok := i.dataByID[userID]
	if !ok {
		return nil, ErrNotFound
	}
	phone := strconv.FormatUint(usr.phone, 10)
	return &users.User{
		ID:          usr.id,
		NickName:    usr.nickName,
		PhoneNumber: users.PhoneNumber(phone),
	}, nil
}

// Delete will remove the user from the system
func (i *InMemoryImproved) Delete(_ context.Context, user *users.User) (*users.User, error) {
	i.m.Lock()
	defer i.m.Unlock()
	usr, ok := i.dataByID[user.ID]
	if !ok {
		return nil, ErrNotFound
	}

	delete(i.dataByID, usr.id)

	phone := strconv.FormatUint(usr.phone, 10)
	return &users.User{
		ID:          usr.id,
		NickName:    usr.nickName,
		PhoneNumber: users.PhoneNumber(phone),
	}, nil
}

// Update will update user with same ID to the new value
func (i *InMemoryImproved) Update(_ context.Context, user *users.User) (*users.User, error) {
	i.m.Lock()
	defer i.m.Unlock()
	storedUser, ok := i.dataByID[user.ID]
	if !ok {
		return nil, ErrNotFound
	}

	if user.NickName != "" {
		storedUser.nickName = user.NickName
	}
	if user.PhoneNumber != "" {
		phone, _ := strconv.ParseUint(string(user.PhoneNumber), 10, 64)
		storedUser.phone = phone
	}

	phone := strconv.FormatUint(storedUser.phone, 10)
	return &users.User{
		ID:          storedUser.id,
		NickName:    storedUser.nickName,
		PhoneNumber: users.PhoneNumber(phone),
	}, nil
}

// Query will create a queryImproved to search users. implements users.Queryier
func (i *InMemoryImproved) Query() userquery.Queryer {
	return &queryImproved{}
}

// Search will execute the search on the user base
func (i *InMemoryImproved) Search(_ context.Context, q userquery.Queryer) ([]*users.User, error) {
	i.m.RLock()
	defer i.m.RUnlock()
	sQuery, ok := q.(*queryImproved)
	if !ok {
		return nil, ErrQueryNotCompatible
	}
	var res []*users.User
	for _, usr := range i.dataByID {
		if sQuery.match(usr) {
			phone := strconv.FormatUint(usr.phone, 10)
			res = append(res, &users.User{
				ID:          usr.id,
				NickName:    usr.nickName,
				PhoneNumber: users.PhoneNumber(phone),
			})
		}
	}
	return res, nil
}

// -- internal implementation --

type queryImproved struct {
	// TODO: change to string set
	ids          []users.ID
	nickName     []string
	phoneNumbers []uint64
}

func (q *queryImproved) ByID(id users.ID) userquery.Queryer {
	q.ids = append(q.ids, id)
	return q
}

func (q *queryImproved) ByNickName(nickName string) userquery.Queryer {
	q.nickName = append(q.nickName, nickName)
	return q
}

func (q *queryImproved) ByPhoneNumber(phoneNumber users.PhoneNumber) userquery.Queryer {
	phone, _ := strconv.ParseUint(string(phoneNumber), 10, 64)
	fmt.Println(phone)
	q.phoneNumbers = append(q.phoneNumbers, phone)
	return q
}

func (q *queryImproved) match(u *internalUser) bool {
	for _, id := range q.ids {
		if id == u.id {
			return true
		}
	}
	for _, nickName := range q.nickName {
		if nickName == u.nickName {
			return true
		}
	}
	for _, phoneNumber := range q.phoneNumbers {
		if phoneNumber == u.phone {
			return true
		}
	}
	return false
}
