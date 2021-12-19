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

func TestSetupUpdate_OK(t *testing.T) {
	userStore := userstore.NewInMemory()
	usr, _ := userStore.Add(context.Background(), &users.User{
		PhoneNumber: "+44 020 030 030",
	})
	update := usecases.SetupUpdateUser(logger.Logger{}, userStore)
	res, err := update(context.Background(), &usecases.UserUpdateReq{
		ID:          usr.ID,
		PhoneNumber: "+44 020 030 042 00",
	})
	require.NoError(t, err)
	require.Equal(t, res.User.PhoneNumber, users.PhoneNumber("+44 020 030 042 00"))
}
