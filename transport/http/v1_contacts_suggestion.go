package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/domain/usecases"
)

func (b *Builder) WithV1SuggestionContact(suggestContact usecases.SuggestContact) *Builder {
	b.router.Get("/v1/contact/{userID}/suggest", func(writer http.ResponseWriter, request *http.Request) {
		req, status, err := parseSuggestionContactRequest(request)
		if err != nil {
			writer.WriteHeader(status)
			return
		}
		res, err := suggestContact(request.Context(), req)
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

func parseSuggestionContactRequest(request *http.Request) (*usecases.SuggestContactReq, int, error) {
	userIDStr := chi.URLParam(request, "userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't parse userID from url: %w", err)
	}
	return &usecases.SuggestContactReq{User: users.ID(userID)}, 0, nil
}
