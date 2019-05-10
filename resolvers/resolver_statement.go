package resolvers

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
)

type statementResolver struct{ *Resolver }

// Ticket 语句所属工单信息
func (r *statementResolver) Ticket(ctx context.Context, obj *models.Statement) (ticket *models.Ticket, err error) {
	rc := gqlapi.ReturnCodeOK
	found := false
	ticket = &models.Ticket{
		TicketID: obj.TicketID,
	}
	if found, err = g.Engine.Get(ticket); err != nil {
		ticket = nil
		rc = gqlapi.ReturnCodeUnknowError
		err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
	} else if !found {
		ticket = nil
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 语句(uuid=%s)的工单不存在。", rc, obj.UUID)
	}

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
