package resolvers

import (
	"context"

	"github.com/mia0x75/halo/gqlapi"
)

// TicketStatusChanged 工单状态变化通知，TODO: 目前调不通，主要是WebSocket连接会主动断开
func (r subscriptionRootResolver) TicketStatusChanged(ctx context.Context) (<-chan *gqlapi.TicketStatusChangePayload, error) {
	requires := []gqlapi.RoleEnum{gqlapi.RoleEnumDeveloper, gqlapi.RoleEnumReviewer, gqlapi.RoleEnumAdmin}
	user, err := SubscriptionAuth(ctx, requires)
	if err != nil {
		return nil, err
	}

	// TODO: 解決不同鉴权的问题
	events := make(chan *gqlapi.TicketStatusChangePayload, 1)

	go func() {
		<-ctx.Done()
		TicketSub.Lock()
		delete(TicketSub.Subscribers, user.UUID)
		TicketSub.Unlock()
	}()

	TicketSub.Lock()
	TicketSub.Subscribers[user.UUID] = events
	TicketSub.Unlock()

	return events, nil
}
