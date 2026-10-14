package mysqlrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestTemplateRepo(t *testing.T) {
	repo := NewTemplateRepo(setupTestDB(t))
	ctx := context.Background()

	tpl := &domain.Template{
		UUID:        newUUID(),
		Subject:     "subject",
		Body:        "hello",
		Description: "a template",
		Version:     1,
		CreateAt:    nowUnix(),
	}
	require.NoError(t, repo.Create(ctx, tpl))

	got, err := repo.GetByUUID(ctx, tpl.UUID)
	require.NoError(t, err)
	require.Equal(t, tpl.Subject, got.Subject)

	updated := nowUnix()
	tpl.Body = "hello world"
	tpl.Version = 2
	tpl.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, tpl))

	got, err = repo.GetByUUID(ctx, tpl.UUID)
	require.NoError(t, err)
	require.Equal(t, "hello world", got.Body)

	list, err := repo.List(ctx, tpl.UUID+"~", 100)
	require.NoError(t, err)
	require.True(t, hasTemplate(list, tpl.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, tpl.UUID))
	_, err = repo.GetByUUID(ctx, tpl.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasTemplate(list []*domain.Template, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
