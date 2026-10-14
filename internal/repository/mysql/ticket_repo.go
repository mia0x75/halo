package mysqlrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.TicketRepository = (*TicketRepo)(nil)

type TicketRepo struct {
	q *Queries
}

func NewTicketRepo(q *Queries) *TicketRepo {
	return &TicketRepo{q: q}
}

func ticketFromRow(r MmTicket) *domain.Ticket {
	return &domain.Ticket{
		TicketID:   uint(r.TicketID),
		UUID:       r.Uuid,
		ClusterID:  uint(r.ClusterID),
		Database:   r.Database,
		Subject:    r.Subject,
		Content:    r.Content,
		Status:     domain.TicketStatus(r.Status),
		UserID:     uint(r.UserID),
		ReviewerID: uint(r.ReviewerID),
		CronID:     uintFromNull(r.CronID),
		Version:    int(r.Version),
		UpdateAt:   uintFromNull(r.UpdateAt),
		CreateAt:   uint(r.CreateAt),
	}
}

func (r *TicketRepo) Create(ctx context.Context, t *domain.Ticket) (int64, error) {
	id, err := r.q.CreateTicket(ctx, CreateTicketParams{
		Uuid:       t.UUID,
		ClusterID:  uint32(t.ClusterID),
		Database:   t.Database,
		Subject:    t.Subject,
		Content:    t.Content,
		Status:     uint8(t.Status),
		UserID:     uint32(t.UserID),
		ReviewerID: uint32(t.ReviewerID),
		CronID:     nullFromUint(t.CronID),
		Version:    uint32(t.Version),
		UpdateAt:   nullFromUint(t.UpdateAt),
		CreateAt:   uint32(t.CreateAt),
	})
	if err != nil {
		return 0, wrap("ticket.Create", err)
	}
	return id, nil
}

func (r *TicketRepo) GetByID(ctx context.Context, id int64) (*domain.Ticket, error) {
	row, err := r.q.GetTicketByID(ctx, uint32(id))
	if err != nil {
		return nil, wrap("ticket.GetByID", err)
	}
	return ticketFromRow(row), nil
}

func (r *TicketRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Ticket, error) {
	row, err := r.q.GetTicketByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("ticket.GetByUUID", err)
	}
	return ticketFromRow(row), nil
}

func (r *TicketRepo) Update(ctx context.Context, t *domain.Ticket) error {
	err := r.q.UpdateTicket(ctx, UpdateTicketParams{
		Uuid:       t.UUID,
		ClusterID:  uint32(t.ClusterID),
		Database:   t.Database,
		Subject:    t.Subject,
		Content:    t.Content,
		Status:     uint8(t.Status),
		UserID:     uint32(t.UserID),
		ReviewerID: uint32(t.ReviewerID),
		CronID:     nullFromUint(t.CronID),
		Version:    uint32(t.Version),
		UpdateAt:   nullFromUint(t.UpdateAt),
		TicketID:   uint32(t.TicketID),
	})
	return wrap("ticket.Update", err)
}

func (r *TicketRepo) Delete(ctx context.Context, id int64) error {
	return wrap("ticket.Delete", r.q.DeleteTicketByID(ctx, uint32(id)))
}

func (r *TicketRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("ticket.DeleteByUUID", r.q.DeleteTicketByUUID(ctx, uuid))
}

func (r *TicketRepo) List(ctx context.Context, cursor int64, limit int) ([]*domain.Ticket, error) {
	rows, err := r.q.ListTickets(ctx, ListTicketsParams{
		TicketID: uint32(cursor),
		Limit:    int32(limit),
	})
	if err != nil {
		return nil, wrap("ticket.List", err)
	}
	out := make([]*domain.Ticket, 0, len(rows))
	for _, row := range rows {
		out = append(out, ticketFromRow(row))
	}
	return out, nil
}
