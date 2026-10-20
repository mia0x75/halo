package resolvers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
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

// statementCursorUpperBound 是语句游标的首屏上界：UUID 都小于它。
const statementCursorUpperBound = "ffffffff-ffff-ffff-ffff-ffffffffffff"

// Statements 全量语句列表，按 UUID 倒序分页。
func (r *queryRootResolver) Statements(ctx context.Context, after *string, first *int) (*gqlapi.StatementConnection, error) {
	limit := 20
	if first != nil && *first > 0 {
		limit = *first
	}

	// 仓储层按 uuid 倒序取「小于游标」的一段，首屏用上界哨兵。
	cursor := statementCursorUpperBound
	if after != nil && *after != "" {
		decoded, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, fmt.Errorf("错误代码: %s, 错误信息: 游标无效。", gqlapi.ReturnCodeInvalidParams)
		}
		cursor = strings.TrimPrefix(string(decoded), "cursor")
	}

	domainStmts, err := r.app.Services.Statement.List(ctx, cursor, limit)
	if err != nil {
		return nil, err
	}

	edges := []*gqlapi.StatementEdge{}
	for _, stmt := range domainStmts {
		edges = append(edges, &gqlapi.StatementEdge{
			Node:   bridge.StatementFromDomain(stmt),
			Cursor: EncodeCursor(stmt.UUID),
		})
	}

	pageInfo := &gqlapi.PageInfo{}
	if len(edges) > 0 {
		pageInfo = &gqlapi.PageInfo{
			HasPreviousPage: after != nil && *after != "",
			HasNextPage:     len(edges) == limit,
			StartCursor:     edges[0].Cursor,
			EndCursor:       edges[len(edges)-1].Cursor,
		}
	}

	return &gqlapi.StatementConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}
