package contactstore

import (
	"context"
	"sync"

	"github.com/platelk/contactgraph/domain/models/contacts"
	"github.com/platelk/contactgraph/domain/models/users"
)

// InMemoryMap is a ContactStore implementation based on a HashMap data structure.
type InMemoryMap struct {
	m        sync.RWMutex
	contacts map[users.ID]contacts.List
}

// NewInMemoryMap instantiate a new InMemoryMap
func NewInMemoryMap() *InMemoryMap {
	return &InMemoryMap{contacts: map[users.ID]contacts.List{}}
}

// Connect implements ContactStore.Connect
func (i *InMemoryMap) Connect(_ context.Context, from, to users.ID) error {
	i.m.Lock()
	defer i.m.Unlock()
	contactList, ok := i.contacts[from]
	if !ok {
		contactList = contacts.List{}
		i.contacts[from] = contactList
	}
	contactList[to] = struct{}{}
	return nil
}

// Lookup implements ContactStore.Lookup
func (i *InMemoryMap) Lookup(_ context.Context, userID users.ID) (contacts.List, error) {
	i.m.RLock()
	defer i.m.RUnlock()
	contactList, ok := i.contacts[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return contactList, nil
}

// ReverseLookup implements ContactStore.ReverseLookup
func (i *InMemoryMap) ReverseLookup(_ context.Context, userID users.ID) (contacts.List, error) {
	i.m.RLock()
	defer i.m.RUnlock()
	l := contacts.List{}
	for k, friends := range i.contacts {
		if _, ok := friends[userID]; ok {
			l[k] = struct{}{}
		}
	}
	return l, nil
}

// Len implements ContactStore
func (i *InMemoryMap) Len() uint {
	i.m.RLock()
	defer i.m.RUnlock()
	return uint(len(i.contacts))
}

func (i *InMemoryMap) Iterate(ctx context.Context, it IteratorFunc) error {
	i.m.RLock()
	defer i.m.RUnlock()
	for user, list := range i.contacts {
		it(user, list)
	}
	return nil
}
