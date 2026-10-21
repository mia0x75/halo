package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

type fakeMessageRepo struct {
	messages []*domain.Message
	err      error
}

func (f *fakeMessageRepo) Create(_ context.Context, m *domain.Message) error {
	f.messages = append(f.messages, m)
	return f.err
}

func (f *fakeMessageRepo) GetByUUID(_ context.Context, uuid string) (*domain.Message, error) {
	return nil, f.err
}

func (f *fakeMessageRepo) ListByUser(
	_ context.Context,
	_ uint,
	_ uint,
	_ bool,
	_ string,
	_ string,
	_ int,
) ([]*domain.Message, error) {
	return f.messages, f.err
}

func (f *fakeMessageRepo) CountByUser(_ context.Context, _ uint, _ bool, _ string, _ string) (uint, error) {
	return uint(len(f.messages)), f.err
}

func (f *fakeMessageRepo) Delete(_ context.Context, _ string, _ uint) error {
	return f.err
}

func (f *fakeMessageRepo) CountUnread(_ context.Context, _ uint) (uint, error) {
	return 0, f.err
}

func (f *fakeMessageRepo) MarkRead(_ context.Context, _ string, _ uint) error {
	return f.err
}

func (f *fakeMessageRepo) MarkAllRead(_ context.Context, userID uint) (int64, error) {
	if f.err != nil {
		return 0, f.err
	}
	var rows int64
	for _, message := range f.messages {
		if message.UserID == userID && message.ReadAt == nil {
			readAt := uint(1)
			message.ReadAt = &readAt
			rows++
		}
	}
	return rows, nil
}

func TestMessageServiceMarkAllRead(t *testing.T) {
	ctx := context.Background()
	repo := &fakeMessageRepo{messages: []*domain.Message{
		{UserID: 1, Title: "a"},
		{UserID: 1, Title: "b"},
		{UserID: 2, Title: "c"},
	}}
	svc := NewMessageService(repo)

	rows, err := svc.MarkAllRead(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, int64(2), rows)
	require.NotNil(t, repo.messages[0].ReadAt)
	require.NotNil(t, repo.messages[1].ReadAt)
	require.Nil(t, repo.messages[2].ReadAt)
}

func TestMessageServiceMarkAllReadRejectsMissingUser(t *testing.T) {
	svc := NewMessageService(&fakeMessageRepo{})

	_, err := svc.MarkAllRead(context.Background(), 0)
	require.Error(t, err)
}

func TestMessageServiceMarkAllReadPropagatesError(t *testing.T) {
	svc := NewMessageService(&fakeMessageRepo{err: errors.New("boom")})

	_, err := svc.MarkAllRead(context.Background(), 1)
	require.Error(t, err)
}
