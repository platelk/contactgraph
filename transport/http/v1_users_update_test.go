package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/domain/usecases"
	"github.com/platelk/contactgraph/infra/logger"

	"github.com/stretchr/testify/require"
)

func TestBuilder_WithV1UpdateUser(t *testing.T) {
	router := NewBuilder(logger.Logger{}, Config{}).WithV1UpdateUser(func(ctx context.Context, req *usecases.UserUpdateReq) (*usecases.UserUpdateResp, error) {
		require.NotEmpty(t, req.ID)
		return &usecases.UserUpdateResp{User: &users.User{
			ID:          req.ID,
			NickName:    req.NickName,
			PhoneNumber: users.PhoneNumber(req.PhoneNumber),
		}}, nil
	}).router

	req := httptest.NewRequest("PUT", "http://localhost/v1/user", strings.NewReader(`
	{"id": 42, "phone_number": "+1 123 456 789 00", "nick_name": "test"}
	`))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
