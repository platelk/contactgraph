package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/domain/usecases"
	"github.com/platelk/contactgraph/infra/logger"

	"github.com/stretchr/testify/require"
)

func TestBuilder_WithV1CreateUser(t *testing.T) {
	router := NewBuilder(logger.Logger{}, Config{}).
		WithV1CreateUser(func(ctx context.Context, req *usecases.UserCreateReq) (*usecases.UserCreateResp, error) {
			require.NotEmpty(t, req.NickName)
			require.NotEmpty(t, req.PhoneNumber)
			return &usecases.UserCreateResp{User: &users.User{
				ID:          42,
				NickName:    req.NickName,
				PhoneNumber: users.PhoneNumber(req.PhoneNumber),
			}}, nil
		}).router

	req := httptest.NewRequest("POST", "http://localhost/v1/user", strings.NewReader(`
	{"phone_number": "+1 123 456 789 00", "nick_name": "test"}
	`))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, string(body), "42")
}
