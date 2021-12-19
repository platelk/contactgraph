package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/platelk/contactgraph/domain/usecases"
)

// WithV1CreateUser will add http endpoint to create new user
func (b *Builder) WithV1CreateUser(createUser usecases.CreateUser) *Builder {
	b.router.Post("/v1/user", func(writer http.ResponseWriter, request *http.Request) {
		req, status, err := parseRequest(request)
		if err != nil {
			writer.WriteHeader(status)
			return
		}
		res, err := createUser(request.Context(), req)
		switch {
		case errors.Is(err, usecases.ErrInvalidUser):
			writer.WriteHeader(http.StatusBadRequest)
			_, err = writer.Write([]byte(err.Error()))
			if err != nil {
				b.log.Error().Err(err).Msg("can't write response.")
			}
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

func parseRequest(request *http.Request) (*usecases.UserCreateReq, int, error) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't read body: %w", err)
	}
	var req usecases.UserCreateReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't parse json body: %w", err)
	}

	return &req, 0, nil
}
