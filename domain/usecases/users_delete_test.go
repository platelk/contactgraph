package usecases_test

import (
	"context"
	"testing"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/domain/usecases"
	"github.com/platelk/contactgraph/infra/logger"
	"github.com/platelk/contactgraph/infra/userstore"

	"github.com/stretchr/testify/require"
)

func TestSetupDelete_OK(t *testing.T) {
	userStore := userstore.NewInMemory()
	usr, err := userStore.Add(context.Background(), &users.User{
		PhoneNumber: "+1 123 456 789 00",
	})
	require.NoError(t, err)
	del := usecases.SetupDeleteUser(logger.Logger{}, userStore)
	res, err := del(context.Background(), &usecases.UserDeleteReq{
		ID: usr.ID,
	})
	require.NoError(t, err)
	require.Equal(t, res.User.PhoneNumber, users.PhoneNumber("+1 123 456 789 00"))
}
