package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestOptionServiceUpdateValidatesValue(t *testing.T) {
	ctx := context.Background()
	repo := &stubOptionRepo{}
	logs := &fakeLogRepo{}
	svc := NewOptionService(&fakeTx{repos: Repos{Option: repo, Log: logs}}, repo)

	err := svc.Update(ctx, &domain.Option{UUID: "o1", Name: "smtp.enabled", Parameter: "yes"}, 7)
	require.ErrorIs(t, err, domain.ErrInvalidInput)
	require.Empty(t, logs.entries)
}

func TestOptionServiceUpdateRecordsOldAndNewValue(t *testing.T) {
	ctx := context.Background()
	repo := &stubOptionRepo{existing: &domain.Option{UUID: "o1", Name: "smtp.host", Parameter: "old"}}
	logs := &fakeLogRepo{}
	svc := NewOptionService(&fakeTx{repos: Repos{Option: repo, Log: logs}}, repo)

	require.NoError(t, svc.Update(ctx, &domain.Option{UUID: "o1", Name: "smtp.host", Parameter: "new"}, 7))
	require.Len(t, logs.entries, 1)
	require.Equal(t, "option.update|option:o1|old->new", logs.entries[0].Operation)
}

func TestOptionServiceUpdateMasksSecretInAudit(t *testing.T) {
	ctx := context.Background()
	repo := &stubOptionRepo{existing: &domain.Option{UUID: "o1", Name: "smtp.password", Parameter: "old"}}
	logs := &fakeLogRepo{}
	svc := NewOptionService(&fakeTx{repos: Repos{Option: repo, Log: logs}}, repo)

	require.NoError(t, svc.Update(ctx, &domain.Option{UUID: "o1", Name: "smtp.password", Parameter: "new"}, 7))
	require.Len(t, logs.entries, 1)
	require.Equal(t, "option.update|option:o1|********->********", logs.entries[0].Operation)
}
