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

func TestSetupCreate_OK(t *testing.T) {
	create := usecases.SetupCreateUser(logger.Logger{}, userstore.NewInMemory())
	res, err := create(context.Background(), &usecases.UserCreateReq{
		NickName:    "test",
		PhoneNumber: "+33 93 934 8893",
	})
	require.NoError(t, err)
	require.Equal(t, res.User.PhoneNumber, users.PhoneNumber("33939348893"))
}
