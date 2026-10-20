package postgresrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestMetricRepo(t *testing.T) {
	repo := NewMetricRepo(setupTestDB(t))
	ctx := context.Background()

	key := domain.MetricKey{
		Name:    "name-" + newSuffix(),
		Scope:   domain.MetricScopeUser,
		ScopeID: newUUID(),
		Bucket:  "",
	}
	m := &domain.Metric{
		UUID:     newUUID(),
		Name:     key.Name,
		Scope:    key.Scope,
		ScopeID:  key.ScopeID,
		Type:     domain.MetricTypeCounter,
		Bucket:   key.Bucket,
		Value:    12.3456,
		Version:  1,
		CreateAt: nowUnix(),
	}
	require.NoError(t, repo.Create(ctx, m))

	got, err := repo.GetByUUID(ctx, m.UUID)
	require.NoError(t, err)
	require.Equal(t, 12.3456, got.Value)
	require.Equal(t, key.Name, got.Name)

	byKey, err := repo.GetByKey(ctx, key)
	require.NoError(t, err)
	require.Equal(t, m.UUID, byKey.UUID)

	updated := nowUnix()
	m.Value = 99.5
	m.Version = 2
	m.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, m))

	got, err = repo.GetByUUID(ctx, m.UUID)
	require.NoError(t, err)
	require.Equal(t, 99.5, got.Value)

	list, err := repo.List(ctx, m.UUID+"~", 100)
	require.NoError(t, err)
	require.True(t, hasMetric(list, m.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, m.UUID))
	_, err = repo.GetByUUID(ctx, m.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasMetric(list []*domain.Metric, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}

func TestMetricIncrement(t *testing.T) {
	repo := NewMetricRepo(setupTestDB(t))
	ctx := context.Background()
	key := domain.MetricKey{
		Name:    "counter-" + newSuffix(),
		Scope:   domain.MetricScopeOverall,
	}

	require.NoError(t, repo.Increment(ctx, key, 5))
	got, err := repo.GetByKey(ctx, key)
	require.NoError(t, err)
	require.Equal(t, 5.0, got.Value)

	require.NoError(t, repo.Increment(ctx, key, -2))
	got, err = repo.GetByKey(ctx, key)
	require.NoError(t, err)
	require.Equal(t, 3.0, got.Value)
}

func TestMetricSetGauge(t *testing.T) {
	repo := NewMetricRepo(setupTestDB(t))
	ctx := context.Background()
	key := domain.MetricKey{
		Name:    "gauge-" + newSuffix(),
		Scope:   domain.MetricScopeUser,
		ScopeID: newUUID(),
	}

	require.NoError(t, repo.SetGauge(ctx, key, 100))
	got, err := repo.GetByKey(ctx, key)
	require.NoError(t, err)
	require.Equal(t, 100.0, got.Value)

	// A gauge holds the latest value: writing again replaces it instead of
	// accumulating it.
	require.NoError(t, repo.SetGauge(ctx, key, 200))
	got, err = repo.GetByKey(ctx, key)
	require.NoError(t, err)
	require.Equal(t, 200.0, got.Value)
}
