package contactstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/platelk/contactgraph/domain/models/users"
)

func runTestSuiteContactStore(t *testing.T, storeCreator func() ContactStore) {
	t.Run("connect", func(t *testing.T) {
		runTestConnectContactStore(t, storeCreator)
	})
	t.Run("lookup", func(t *testing.T) {
		runTestLookupContactStore(t, storeCreator)
	})
	t.Run("reverse_lookup", func(t *testing.T) {
		runTestReverseLookupContactStore(t, storeCreator)
	})
}

func runTestConnectContactStore(t *testing.T, storeCreator func() ContactStore) {
	t.Run("connect 2 user", func(t *testing.T) {
		from, to := 42, 43
		store := storeCreator()
		err := store.Connect(context.Background(), users.ID(from), users.ID(to))
		require.NoError(t, err)
	})
}

func runTestLookupContactStore(t *testing.T, storeCreator func() ContactStore) {
	t.Run("lookup with 2 contact", func(t *testing.T) {
		store := storeCreator()
		err := store.Connect(context.Background(), users.ID(42), users.ID(43))
		require.NoError(t, err)
		err = store.Connect(context.Background(), users.ID(42), users.ID(44))
		require.NoError(t, err)
		friends, err := store.Lookup(context.Background(), users.ID(42))
		require.NoError(t, err)
		require.Len(t, friends, 2)
		require.Contains(t, friends, users.ID(43))
		require.Contains(t, friends, users.ID(44))
	})
	t.Run("not found", func(t *testing.T) {
		store := storeCreator()
		_, err := store.Lookup(context.Background(), users.ID(42))
		require.ErrorIs(t, err, ErrNotFound)
	})
}

func runTestReverseLookupContactStore(t *testing.T, storeCreate func() ContactStore) {
	t.Run("2 contact have users in there contacts", func(t *testing.T) {
		store := storeCreate()

		err := store.Connect(context.Background(), users.ID(43), users.ID(42))
		require.NoError(t, err)
		err = store.Connect(context.Background(), users.ID(44), users.ID(42))
		require.NoError(t, err)

		friends, err := store.ReverseLookup(context.Background(), users.ID(42))
		require.NoError(t, err)
		require.Len(t, friends, 2)
		require.Contains(t, friends, users.ID(43))
		require.Contains(t, friends, users.ID(44))
	})
}
