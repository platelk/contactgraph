package usecases

import (
	"context"
	"fmt"

	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/infra/logger"
)

type contactConnector interface {
	Connect(ctx context.Context, user1, user2 users.ID) error
}

type ConnectContactReq struct {
	From users.ID `json:"from"`
	To   users.ID `json:"to"`
}
type ConnectContactResp struct{}

type ConnectContact func(ctx context.Context, req *ConnectContactReq) (*ConnectContactResp, error)

func SetupConnectContact(log logger.Logger, connector contactConnector) ConnectContact {
	return func(ctx context.Context, req *ConnectContactReq) (*ConnectContactResp, error) {
		if err := connector.Connect(ctx, req.From, req.To); err != nil {
			return nil, fmt.Errorf("can't connect user[%v] and user[%v]: %w", req.From, req.To, err)
		}
		return &ConnectContactResp{}, nil
	}
}
