package resolvers

import (
	"fmt"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/domain"
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
	page, hasNextPage, hasPreviousPage, err := pageInMemory(
		tickets,
		func(ticket *domain.Ticket) int { return int(ticket.TicketID) },
		after, before, first, last,
	)
	if err != nil {
		return nil, err
	}

	edges := make([]*gqlapi.TicketEdge, 0, len(page))
	for _, ticket := range page {
		edges = append(edges, &gqlapi.TicketEdge{
			Node:   bridge.TicketFromDomain(ticket),
			Cursor: EncodeCursor(fmt.Sprintf("%d", ticket.TicketID)),
		})
	}

	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     "",
		EndCursor:       "",
	}
	if len(edges) > 0 {
		pageInfo.StartCursor = edges[0].Cursor
		pageInfo.EndCursor = edges[len(edges)-1].Cursor
	}

	// The store handed back the complete list, so the total spans every match
	// and not just this page; the side nav badge reads it directly.
	return &gqlapi.TicketConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(tickets),
	}, nil
}
