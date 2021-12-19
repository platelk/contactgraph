package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/platelk/contactgraph/domain/usecases"
)

// WithV1DeleteUser will add http endpoint to delete new user
func (b *Builder) WithV1DeleteUser(deleteUser usecases.DeleteUser) *Builder {
	b.router.Delete("/v1/user", func(writer http.ResponseWriter, request *http.Request) {
		req, status, err := parseDeleteRequest(request)
		if err != nil {
			writer.WriteHeader(status)
			return
		}
		res, err := deleteUser(request.Context(), req)
		switch {
		case err != nil:
			b.log.Error().Err(err).Send()
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = writer.Write([]byte(err.Error()))
			if err != nil {
				b.log.Error().Err(err).Msg("can't write response.")
			}
		default:
			data, err := json.Marshal(res)
			if err != nil {
				b.log.Error().Err(err).Msg("can't marshal response.")
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = writer.Write(data)
			if err != nil {
				b.log.Error().Err(err).Msg("can't write response.")
			}
		}
	})
	return b
}

func parseDeleteRequest(request *http.Request) (*usecases.UserDeleteReq, int, error) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't read body: %w", err)
	}
	var req usecases.UserDeleteReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't parse json body: %w", err)
	}
	return &req, 0, nil
}
