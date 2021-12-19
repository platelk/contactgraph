package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/domain/usecases"
)

// WithV1SearchUser will add http endpoint to search users
// Note: here pagination is not implemented, so too many users can break the response
func (b *Builder) WithV1SearchUser(searchUser usecases.SearchUser) *Builder {
	b.router.Get("/v1/users", func(writer http.ResponseWriter, request *http.Request) {
		req, status, err := parseSearchRequest(request)
		if err != nil {
			writer.WriteHeader(status)
			return
		}
		res, err := searchUser(request.Context(), req)
		switch {
		case err != nil:
			b.log.Error().Err(err).Send()
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = writer.Write([]byte(err.Error()))
			if err != nil {
				b.log.Error().Err(err).Msg("can't write response.")
			}
		default:
			data, _ := json.Marshal(res)
			_, _ = writer.Write(data)
		}
	})
	return b
}

func parseSearchRequest(request *http.Request) (*usecases.UserSearchReq, int, error) {
	var ids []users.ID
	for _, id := range request.URL.Query()["id"] {
		parsedID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("can't parse id: %w", err)
		}
		ids = append(ids, users.ID(parsedID))
	}
	return &usecases.UserSearchReq{
		IDs:       ids,
		Emails:    request.URL.Query()["email"],
		FirstName: request.URL.Query()["first_name"],
		LastName:  request.URL.Query()["last_name"],
		NickName:  request.URL.Query()["nick_name"],
	}, 0, nil
}
