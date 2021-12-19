package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/platelk/contactgraph/domain/usecases"
)

func (b *Builder) WithV1ConnectContact(connectContact usecases.ConnectContact) *Builder {
	b.router.Post("/v1/contact", func(writer http.ResponseWriter, request *http.Request) {
		req, status, err := parseConnectContactRequest(request)
		if err != nil {
			writer.WriteHeader(status)
			_, err = writer.Write([]byte(err.Error()))
			if err != nil {
				b.log.Error().Err(err).Msg("can't write response.")
			}
			return
		}
		res, err := connectContact(request.Context(), req)
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

func parseConnectContactRequest(request *http.Request) (*usecases.ConnectContactReq, int, error) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't read body: %w", err)
	}
	var req usecases.ConnectContactReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't parse json body: %w", err)
	}

	return &req, 0, nil
}
