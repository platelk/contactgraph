package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/infra/contactstore"
	"github.com/platelk/contactgraph/infra/userstore"
)

// DevGenerateReq define the fields required to generate a population with their connections.
type DevGenerateReq struct {
	Population uint `json:"population"`
	Connection uint `json:"connection"`
}

// WithV1DevGenerateData will add default endpoints to check its status
func (b *Builder) WithV1DevGenerateData(userStore userstore.UserStore, contactStore contactstore.ContactStore) *Builder {
	b.router.Post("/v1/dev/generate", func(writer http.ResponseWriter, request *http.Request) {
		req, status, err := parseDevGenerateRequest(request)
		if err != nil {
			writer.WriteHeader(status)
			return
		}
		go func() {
			generateUser(req.Population, userStore)
			generateConnection(req.Population, req.Connection, contactStore)
		}()
		switch {
		case err != nil:
			b.log.Error().Err(err).Send()
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = writer.Write([]byte(err.Error()))
			if err != nil {
				b.log.Error().Err(err).Msg("can't write response.")
			}
		default:
			writer.WriteHeader(http.StatusAccepted)
		}
	})
	return b
}

func parseDevGenerateRequest(request *http.Request) (*DevGenerateReq, int, error) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't read body: %w", err)
	}
	var req DevGenerateReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("can't parse json body: %w", err)
	}

	return &req, 0, nil
}

func generateUser(n uint, usersStore userstore.UserStore) {
	for i := uint(0); i < n; i++ {
		_, _ = usersStore.Add(context.Background(), &users.User{
			NickName:    fmt.Sprintf("usr%d", i),
			PhoneNumber: users.PhoneNumber(fmt.Sprintf("33%09d", i)),
		})
	}
}

func generateConnection(population uint, avgConn uint, store contactstore.ContactStore) {
	r := rand.New(rand.NewSource(99))
	for i := uint(0); i < population; i++ {
		for j := uint(0); j < avgConn; j++ {
			to := r.Int63n(int64(population))
			if to == int64(i) {
				continue
			}
			_ = store.Connect(context.Background(), users.ID(i), users.ID(to))
		}
	}
}
