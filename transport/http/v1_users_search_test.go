package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/domain/usecases"
	"github.com/platelk/contactgraph/infra/logger"

	"github.com/stretchr/testify/require"
)

func TestBuilder_WithV1SearchUser(t *testing.T) {
	router := NewBuilder(logger.Logger{}, Config{}).
		WithV1SearchUser(func(ctx context.Context, req *usecases.UserSearchReq) (*usecases.UserSearchResp, error) {
			require.NotEmpty(t, req.IDs)
			require.Contains(t, req.IDs, users.ID(41))
			require.Contains(t, req.IDs, users.ID(42))
			return &usecases.UserSearchResp{}, nil
		}).router

	req := httptest.NewRequest("GET", "http://localhost/v1/users?id=41&id=42", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
