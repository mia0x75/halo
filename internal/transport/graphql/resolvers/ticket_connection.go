package resolvers

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// ticketConnection builds a cursor page of tickets. The cursor carries the
// ticket id, so the page is a keyset over ticket_id instead of an offset.
func ticketConnection(
	tickets []*domain.Ticket,
	after *string,
	before *string,
	first *int,
	last *int,
) (*gqlapi.TicketConnection, error) {
	if first != nil && last != nil {
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", gqlapi.ReturnCodeInvalidParams)
	}
	if after != nil && before != nil {
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", gqlapi.ReturnCodeInvalidParams)
	}

	from := uint(math.MaxInt64)
	if after != nil {
		id, err := decodeTicketCursor(*after)
		if err != nil {
			return nil, err
		}
		from = id
	}

	page := make([]*domain.Ticket, 0, len(tickets))
	for _, ticket := range tickets {
		if ticket.TicketID >= from {
			continue
		}
		page = append(page, ticket)
	}

	limit := len(page)
	if first != nil && *first < limit {
		limit = *first
	}
	if last != nil && *last < limit {
		limit = *last
	}
	page = page[:limit]

	edges := make([]*gqlapi.TicketEdge, 0, len(page))
	for _, ticket := range page {
		edges = append(edges, &gqlapi.TicketEdge{
			Node:   bridge.TicketFromDomain(ticket),
			Cursor: EncodeCursor(fmt.Sprintf("%d", ticket.TicketID)),
		})
	}

	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: false,
		HasNextPage:     false,
		StartCursor:     "",
		EndCursor:       "",
	}
	if len(edges) > 0 {
		pageInfo.StartCursor = edges[0].Cursor
		pageInfo.EndCursor = edges[len(edges)-1].Cursor
		pageInfo.HasNextPage = len(page) == limit && limit > 0
	}

	return &gqlapi.TicketConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}

// decodeTicketCursor reads the ticket id out of a cursor.
func decodeTicketCursor(cursor string) (uint, error) {
	raw, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("错误代码: %s, 错误信息: 游标(cursor=%s)无效。", gqlapi.ReturnCodeInvalidParams, cursor)
	}
	id, err := strconv.ParseUint(strings.TrimPrefix(string(raw), "cursor"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("错误代码: %s, 错误信息: 游标(cursor=%s)无效。", gqlapi.ReturnCodeInvalidParams, cursor)
	}
	return uint(id), nil
}

// ticketNodes converts the domain tickets into the GraphQL models.
func ticketNodes(tickets []*domain.Ticket) []*models.Ticket {
	out := make([]*models.Ticket, 0, len(tickets))
	for _, ticket := range tickets {
		out = append(out, bridge.TicketFromDomain(ticket))
	}
	return out
}
