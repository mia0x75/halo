package resolvers

import (
	"context"

	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// TicketStatusChanged 工单状态变化通知。
func (r *subscriptionRootResolver) TicketStatusChanged(ctx context.Context) (<-chan *gqlapi.TicketStatusChangePayload, error) {
	requires := []gqlapi.RoleEnum{gqlapi.RoleEnumDeveloper, gqlapi.RoleEnumReviewer, gqlapi.RoleEnumAdmin}
	user, err := SubscriptionAuth(ctx, requires)
	if err != nil {
		return nil, err
	}

	events := make(chan *gqlapi.TicketStatusChangePayload, 1)

	go func() {
		<-ctx.Done()
		TicketSub.Lock()
		if current, ok := TicketSub.Subscribers[user.UUID]; ok && current == events {
			delete(TicketSub.Subscribers, user.UUID)
			// 关闭通道，订阅方知道连接已结束。
			close(events)
		}
		TicketSub.Unlock()
	}()

	TicketSub.Lock()
	TicketSub.Subscribers[user.UUID] = events
	TicketSub.Unlock()

	return events, nil
}
