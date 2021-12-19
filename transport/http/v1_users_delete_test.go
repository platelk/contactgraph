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

func TestBuilder_WithV1DeleteUser(t *testing.T) {
	router := NewBuilder(logger.Logger{}, Config{}).WithV1DeleteUser(func(ctx context.Context, req *usecases.UserDeleteReq) (*usecases.UserDeleteResp, error) {
		require.NotEmpty(t, req.ID)
		return &usecases.UserDeleteResp{User: &users.User{
			ID: req.ID,
		}}, nil
	}).router

	req := httptest.NewRequest("DELETE", "http://localhost/v1/user", strings.NewReader(`
	{"id": 42}
	`))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
