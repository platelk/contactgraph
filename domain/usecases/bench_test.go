package usecases_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/infra/contactstore"
	"github.com/platelk/contactgraph/infra/userstore"
)

var population = [...]uint{100, 1_000, 10_000, 100_000, 1_000_000, 10_000_000}
var connection = [...]uint{1, 5, 10, 50}

func generateUser(n uint, usersStore userstore.UserStore) {
	for i := uint(0); i < n; i++ {
		usersStore.Add(context.Background(), &users.User{
			NickName:    fmt.Sprintf("test-%d", i),
			PhoneNumber: users.PhoneNumber(fmt.Sprintf("%012d", i)),
		})
	}
}

func generateConnection(population uint, avgConn uint, store contactstore.ContactStore) {
	r := rand.New(rand.NewSource(99))
	for i := uint(0); i < population; i++ {
		for j := uint(0); j < avgConn; j++ {
			to := r.Int63n(int64(population))
			if to == int64(i) {
				continue
			}
			_ = store.Connect(context.Background(), users.ID(i), users.ID(to))
		}
	}
}

func BenchmarkContact(b *testing.B) {
	for _, pop := range population {
		usrStore := userstore.NewInMemory()
		generateUser(pop, usrStore)
		for _, conn := range connection {
			connStore := contactstore.NewShardStore(
				contactstore.NewInMemoryMap(),
				contactstore.NewInMemoryMap(),
				contactstore.NewInMemoryMap(),
				contactstore.NewInMemoryMap(),
				contactstore.NewInMemoryMap(),
			)
			generateConnection(pop, conn, connStore)
			b.Run(fmt.Sprintf("lookup_%v_%v", pop, conn), func(b *testing.B) {
				benchSetupLookupContact(b, pop, conn, connStore, usrStore)
			})
			b.Run(fmt.Sprintf("reverse_lookup_%v_%v", pop, conn), func(b *testing.B) {
				benchSetupReverseLookupContact(b, pop, conn, connStore, usrStore)
			})
			b.Run(fmt.Sprintf("suggestion_%v_%v", pop, conn), func(b *testing.B) {
				benchSetupSuggestContact(b, pop, conn, connStore, usrStore)
			})
		}
	}
}
