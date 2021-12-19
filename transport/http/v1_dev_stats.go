package http

import (
	"encoding/json"
	"net/http"

	"github.com/platelk/contactgraph/infra/contactstore"
	"github.com/platelk/contactgraph/infra/userstore"
)

type statsResp struct {
	Users         uint `json:"users"`
	Connections   uint `json:"connections"`
	UserConnected uint `json:"user_connected"`
}

// WithV1DevStats will add default endpoints to check its status
func (b *Builder) WithV1DevStats(contactStats contactstore.Stats, userStats userstore.Stats) *Builder {
	b.router.Get("/v1/dev/stats", func(writer http.ResponseWriter, request *http.Request) {
		resp := &statsResp{
			Users:         userStats.Users(),
			Connections:   contactStats.Connection(),
			UserConnected: contactStats.Users(),
		}
		data, err := json.Marshal(resp)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(data)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	return b
}
