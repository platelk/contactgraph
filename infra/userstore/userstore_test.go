package userstore

import (
	"context"
	"errors"
	"testing"

	"github.com/platelk/contactgraph/domain/models/userquery"
	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/domain/usecases"

	"github.com/stretchr/testify/require"
)

type wrongQuery struct{}

func (w *wrongQuery) ByLastName(lastName string) userquery.Queryer {
	panic("implement me")
}

func (w *wrongQuery) ByNickName(nickName string) userquery.Queryer {
	panic("implement me")
}

func (w *wrongQuery) ByPhoneNumber(phoneNumber users.PhoneNumber) userquery.Queryer {
	panic("implement me")
}

func (w *wrongQuery) ByID(id users.ID) userquery.Queryer {
	panic("implement me")
}

func (w *wrongQuery) ByEmail(email string) userquery.Queryer {
	panic("implement me")
}

func (w *wrongQuery) ByFirstName(firstName string) userquery.Queryer {
	panic("implement me")
}

type userStore interface {
	usecases.UserAdder
	usecases.UserUpdater
	usecases.UserDeleter
	usecases.UserSearcher
}

func runTestSuite(t *testing.T, store userStore) {
	runTestAdd(t, store)
	runTestDelete(t, store)
	runTestUpdate(t, store)
	runTestSearch(t, store)
}

func runTestSearch(t *testing.T, store userStore) {
	t.Run("search by ID", func(t *testing.T) {
		usr, _ := store.Add(context.Background(), &users.User{
			PhoneNumber: "44020030042",
		})
		res, _ := store.Search(context.Background(), store.Query().ByID(usr.ID))
		require.NotEmpty(t, res)
		require.Equal(t, usr.ID, res[0].ID)
	})
	t.Run("search by phone number", func(t *testing.T) {
		usr, err := store.Add(context.Background(), &users.User{
			PhoneNumber: "33939348888",
		})
		require.NoError(t, err)
		res, err := store.Search(context.Background(), store.Query().ByPhoneNumber(usr.PhoneNumber))
		require.NoError(t, err)
		require.NotEmpty(t, res)
		require.Equal(t, usr.ID, res[0].ID)
	})
	t.Run("error on no compatible queryImproved", func(t *testing.T) {
		_, err := store.Search(context.Background(), &wrongQuery{})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrQueryNotCompatible))
	})
}

func runTestUpdate(t *testing.T, store userStore) {
	t.Run("update not found user", func(t *testing.T) {
		_, err := store.Update(context.Background(), &users.User{
			ID: 404,
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrNotFound))
	})
	t.Run("update email", func(t *testing.T) {
		usr, err := store.Add(context.Background(), &users.User{
			NickName:    "test",
			PhoneNumber: "33939348893",
		})
		require.NoError(t, err)
		_, err = store.Update(context.Background(), &users.User{
			ID:          usr.ID,
			NickName:    "updated",
			PhoneNumber: "33939348893",
		})
		require.NoError(t, err)

		res, _ := store.Search(context.Background(), store.Query().ByNickName("updated"))
		require.NotEmpty(t, res)
		usrUpdated := res[0]

		require.Equal(t, users.PhoneNumber("33939348893"), usrUpdated.PhoneNumber)
		require.Equal(t, "updated", usrUpdated.NickName)
	})
}

func runTestDelete(t *testing.T, store userStore) {
	t.Run("delete after add", func(t *testing.T) {
		usr, err := store.Add(context.Background(), &users.User{
			NickName:    "test-delete-1",
			PhoneNumber: "33939348893",
		})
		require.NoError(t, err)

		res, _ := store.Search(context.Background(), store.Query().ByNickName(usr.NickName))
		require.NotEmpty(t, res)

		_, err = store.Delete(context.Background(), usr)
		require.NoError(t, err)

		res, _ = store.Search(context.Background(), store.Query().ByNickName(usr.NickName))
		require.Empty(t, res)
	})
	t.Run("delete unknown user", func(t *testing.T) {
		_, err := store.Delete(context.Background(), &users.User{
			ID:          users.ID(42),
			NickName:    "test-delete-2",
			PhoneNumber: "33939348893",
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrNotFound))
	})
}

func runTestAdd(t *testing.T, store userStore) {
	t.Run("add basic user", func(t *testing.T) {
		usr, err := store.Add(context.Background(), &users.User{
			NickName:    "test-add-1",
			PhoneNumber: "33939348893",
		})
		require.NoError(t, err)
		require.NotEmpty(t, usr)
	})
	t.Run("add multiple user", func(t *testing.T) {
		usr, err := store.Add(context.Background(), &users.User{
			NickName:    "test-add-multiple-1",
			PhoneNumber: "33939348893",
		})
		require.NoError(t, err)
		require.NotEmpty(t, usr)
		usr2, err := store.Add(context.Background(), &users.User{
			NickName:    "test-add-multiple-2",
			PhoneNumber: "33939348893",
		})
		require.NoError(t, err)
		require.NotEmpty(t, usr2)
		res, _ := store.Search(context.Background(), store.Query().ByNickName("test-add-multiple-2"))
		require.NotEmpty(t, res)
		require.Equal(t, usr2.ID, res[0].ID, "user: %v != found: %v", usr2.ID, res[0].ID)
	})
}
