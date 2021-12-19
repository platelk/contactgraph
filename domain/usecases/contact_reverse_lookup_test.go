package usecases_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/domain/usecases"
	"github.com/platelk/contactgraph/infra/contactstore"
	"github.com/platelk/contactgraph/infra/logger"
	"github.com/platelk/contactgraph/infra/userstore"
)

var reverseLookupResult *usecases.ReverseLookupContactResp

func benchSetupReverseLookupContact(b *testing.B, pop, conn uint, contactStore contactstore.ContactStore, userStore userstore.UserStore) {
	uc := usecases.SetupReverseLookupContact(logger.Logger{}, contactStore, userStore)
	r := rand.New(rand.NewSource(99))
	b.Run("random_user", func(b *testing.B) {
		var res *usecases.ReverseLookupContactResp
		for n := 0; n < b.N; n++ {
			res, _ = uc(context.Background(), &usecases.ReverseLookupContactReq{User: users.ID(r.Int63n(int64(pop)))})
		}
		reverseLookupResult = res
	})
}
