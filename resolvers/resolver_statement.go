package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/models"
)

type statementResolver struct{ *Resolver }

// Ticket 语句所属工单信息
func (r *statementResolver) Ticket(ctx context.Context, obj *models.Statement) (ticket *models.Ticket, err error) {
	rc := gqlapi.ReturnCodeOK
	t, getErr := r.app.Services.Ticket.Get(ctx, obj.TicketID)
	if getErr != nil {
		ticket = nil
		if errors.Is(getErr, domain.ErrNotFound) {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 语句(uuid=%s)的工单不存在。", rc, obj.UUID)
		} else {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, getErr.Error())
		}
		return
	}
	ticket = bridge.TicketFromDomain(t)

	return
}

// TypeDesc 语句类型
func (r *statementResolver) TypeDesc(ctx context.Context, obj *models.Statement) (string, error) {
	for k, v := range gqlapi.StatementTypeEnumMap {
		if v == obj.Type {
			return string(k), nil
		}
	}
	return "OTHER", nil
}
