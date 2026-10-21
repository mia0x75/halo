package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestGroupServiceCreate(t *testing.T) {
	ctx := context.Background()
	repo := &stubGroupRepo{}
	logs := &fakeLogRepo{}
	svc := NewGroupService(&fakeTx{repos: Repos{Group: repo, Log: logs}}, repo, newFakeClusterRepo(), NewLogService(logs))

	group := &domain.Group{UUID: "g1", Name: "机房 A"}
	require.NoError(t, svc.Create(ctx, group, 7))
	require.Equal(t, domain.GroupStatusNormal, group.Status)
	require.Equal(t, 1, group.Version)
	require.Len(t, logs.entries, 1)
	require.Equal(t, "group.create|group:g1|机房 A", logs.entries[0].Operation)
}

func TestGroupServiceCreateRejectsDuplicateName(t *testing.T) {
	ctx := context.Background()
	repo := &stubGroupRepo{existing: &domain.Group{UUID: "g1", Name: "机房 A"}}
	logs := &fakeLogRepo{}
	svc := NewGroupService(&fakeTx{repos: Repos{Group: repo, Log: logs}}, repo, newFakeClusterRepo(), NewLogService(logs))

	err := svc.Create(ctx, &domain.Group{UUID: "g2", Name: "机房 A"}, 7)
	require.ErrorIs(t, err, domain.ErrConflict)
	require.Empty(t, logs.entries)
}

func TestGroupServiceCreateRejectsEmptyName(t *testing.T) {
	ctx := context.Background()
	repo := &stubGroupRepo{}
	logs := &fakeLogRepo{}
	svc := NewGroupService(&fakeTx{repos: Repos{Group: repo, Log: logs}}, repo, newFakeClusterRepo(), NewLogService(logs))

	err := svc.Create(ctx, &domain.Group{UUID: "g2", Name: "  "}, 7)
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestGroupServiceUpdate(t *testing.T) {
	ctx := context.Background()
	repo := &stubGroupRepo{existing: &domain.Group{UUID: "g1", Name: "机房 A", Version: 3}}
	logs := &fakeLogRepo{}
	svc := NewGroupService(&fakeTx{repos: Repos{Group: repo, Log: logs}}, repo, newFakeClusterRepo(), NewLogService(logs))

	group := &domain.Group{UUID: "g1", Name: "机房 B", Status: domain.GroupStatusBlocked}
	require.NoError(t, svc.Update(ctx, group, 7))
	require.Equal(t, 4, group.Version)
	require.Len(t, logs.entries, 1)
	require.Equal(t, "group.update|group:g1|机房 B", logs.entries[0].Operation)
}

func TestGroupServiceDelete(t *testing.T) {
	ctx := context.Background()
	repo := &stubGroupRepo{existing: &domain.Group{UUID: "g1", Name: "机房 A"}}
	logs := &fakeLogRepo{}
	svc := NewGroupService(&fakeTx{repos: Repos{Group: repo, Log: logs}}, repo, newFakeClusterRepo(), NewLogService(logs))

	require.NoError(t, svc.DeleteByUUID(ctx, "g1", 7))
	require.Len(t, logs.entries, 1)
	require.Equal(t, "group.delete|group:g1|", logs.entries[0].Operation)
}

func TestGroupServiceDeleteRejectsGroupWithClusters(t *testing.T) {
	ctx := context.Background()
	repo := &stubGroupRepo{existing: &domain.Group{UUID: "g1", Name: "机房 A", GroupID: 42}}
	logs := &fakeLogRepo{}
	clusters := newFakeClusterRepo()
	_, err := clusters.Create(ctx, &domain.Cluster{UUID: "c1", GroupID: 42})
	require.NoError(t, err)
	svc := NewGroupService(&fakeTx{repos: Repos{Group: repo, Log: logs}}, repo, clusters, NewLogService(logs))

	err = svc.DeleteByUUID(ctx, "g1", 7)
	require.ErrorIs(t, err, domain.ErrConflict)
	require.Empty(t, repo.deleted)
	require.Empty(t, logs.entries)
}
