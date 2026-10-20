package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

type fakeMessageRepo struct {
	domain.MessageRepository
	created []*domain.Message
}

func (f *fakeMessageRepo) Create(_ context.Context, m *domain.Message) error {
	f.created = append(f.created, m)
	return nil
}

func TestMessageWriterNotifiesSubmitterOnStatusChange(t *testing.T) {
	repo := &fakeMessageRepo{}
	handler := GlobalMessageWriter(repo)

	err := handler(context.Background(), &TicketStatusPatched{
		User:   models.User{UserID: 5},
		Ticket: models.Ticket{UUID: "t1", Subject: "订单归档", UserID: 7, ReviewerID: 5},
	})
	require.NoError(t, err)
	require.Len(t, repo.created, 1)
	require.Equal(t, uint(7), repo.created[0].UserID)
	require.Equal(t, "工单状态更新", repo.created[0].Title)
}

func TestMessageWriterSkipsActor(t *testing.T) {
	repo := &fakeMessageRepo{}
	handler := GlobalMessageWriter(repo)

	err := handler(context.Background(), &TicketStatusPatched{
		User:   models.User{UserID: 7},
		Ticket: models.Ticket{UUID: "t1", Subject: "订单归档", UserID: 7, ReviewerID: 5},
	})
	require.NoError(t, err)
	require.Empty(t, repo.created)
}

func TestMessageWriterNotifiesReviewerOnTicketCreated(t *testing.T) {
	repo := &fakeMessageRepo{}
	handler := GlobalMessageWriter(repo)

	err := handler(context.Background(), &TicketCreated{
		User:   models.User{UserID: 7},
		Ticket: models.Ticket{UUID: "t1", Subject: "订单归档", UserID: 7, ReviewerID: 5},
	})
	require.NoError(t, err)
	require.Len(t, repo.created, 1)
	require.Equal(t, uint(5), repo.created[0].UserID)
	require.Equal(t, "待审核工单", repo.created[0].Title)
}

func TestMessageWriterNotifiesSubmitterOnComment(t *testing.T) {
	repo := &fakeMessageRepo{}
	handler := GlobalMessageWriter(repo)

	err := handler(context.Background(), &CommentCreated{
		User:   models.User{UserID: 5},
		Ticket: models.Ticket{UUID: "t1", Subject: "订单归档", UserID: 7, ReviewerID: 5},
	})
	require.NoError(t, err)
	require.Len(t, repo.created, 1)
	require.Equal(t, uint(7), repo.created[0].UserID)
	require.Equal(t, "工单有新评论", repo.created[0].Title)
}

func TestMessageWriterIgnoresOtherEvents(t *testing.T) {
	repo := &fakeMessageRepo{}
	handler := GlobalMessageWriter(repo)

	err := handler(context.Background(), &UserRegistered{User: models.User{UserID: 7}})
	require.NoError(t, err)
	require.Empty(t, repo.created)
}
