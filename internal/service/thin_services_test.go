package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

// The thin services all share the same write path: repo write plus audit row in
// one transaction. These tests lock in the shared contract per service: reads
// delegate, writes audit, and both repo and audit failures propagate.

func TestAvatarService(t *testing.T) {
	ctx := context.Background()
	repo := &stubAvatarRepo{lastID: 7, got: &domain.Avatar{UUID: "a1"}}
	logs := &fakeLogRepo{}
	svc := NewAvatarService(&fakeTx{repos: Repos{Avatar: repo, Log: logs}}, repo)

	got, err := svc.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.GetByUUID(ctx, "a1")
	require.NoError(t, err)
	_, err = svc.List(ctx, 0, 10)
	require.NoError(t, err)

	id, err := svc.Create(ctx, &domain.Avatar{UUID: "a1"}, 3)
	require.NoError(t, err)
	require.Equal(t, uint(7), id)
	require.Equal(t, "avatar.create|avatar:a1|", logs.entries[0].Operation)

	require.NoError(t, svc.Update(ctx, &domain.Avatar{UUID: "a1"}, 3))
	require.NoError(t, svc.Delete(ctx, 7, 3))
	require.Len(t, logs.entries, 3)
}

func TestAvatarServiceErrors(t *testing.T) {
	ctx := context.Background()
	boom := errors.New("boom")

	repo := &stubAvatarRepo{createErr: boom}
	svc := NewAvatarService(&fakeTx{repos: Repos{Avatar: repo, Log: &fakeLogRepo{}}}, repo)
	_, err := svc.Create(ctx, &domain.Avatar{UUID: "a"}, 1)
	require.ErrorIs(t, err, boom)

	repo = &stubAvatarRepo{lastID: 1}
	svc = NewAvatarService(&fakeTx{repos: Repos{Avatar: repo, Log: &fakeLogRepo{err: boom}}}, repo)
	_, err = svc.Create(ctx, &domain.Avatar{UUID: "a"}, 1)
	require.ErrorIs(t, err, boom)

	repo = &stubAvatarRepo{updateErr: boom}
	svc = NewAvatarService(&fakeTx{repos: Repos{Avatar: repo, Log: &fakeLogRepo{}}}, repo)
	require.ErrorIs(t, svc.Update(ctx, &domain.Avatar{}, 1), boom)

	repo = &stubAvatarRepo{deleteErr: boom}
	svc = NewAvatarService(&fakeTx{repos: Repos{Avatar: repo, Log: &fakeLogRepo{}}}, repo)
	require.ErrorIs(t, svc.Delete(ctx, 1, 1), boom)
}

func TestClusterService(t *testing.T) {
	ctx := context.Background()
	repo := newFakeClusterRepo()
	_, seedErr := repo.Create(ctx, &domain.Cluster{ClusterID: 1, UUID: "c1"})
	require.NoError(t, seedErr)
	logs := &fakeLogRepo{}
	svc := NewClusterService(&fakeTx{repos: Repos{Cluster: repo, Log: logs}}, repo, &fakeExecutor{})

	got, err := svc.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, "c1", got.UUID)
	byUUID, err := svc.GetByUUID(ctx, "c1")
	require.NoError(t, err)
	require.Equal(t, uint(1), byUUID.ClusterID)
	_, err = svc.List(ctx, 0, 10)
	require.NoError(t, err)

	id, err := svc.Create(ctx, &domain.Cluster{UUID: "c2"}, 3)
	require.NoError(t, err)
	require.NotZero(t, id)
	require.Equal(t, "cluster.create|cluster:c2|", logs.entries[0].Operation)

	require.NoError(t, svc.Update(ctx, &domain.Cluster{ClusterID: 1, UUID: "c1"}, 3))
	require.NoError(t, svc.Delete(ctx, 1, 3))
	require.Len(t, logs.entries, 3)
}

func TestCommentService(t *testing.T) {
	ctx := context.Background()
	repo := &stubCommentRepo{lastID: 9, got: &domain.Comment{UUID: "m1"}, byTicket: []*domain.Comment{{UUID: "m1"}}}
	logs := &fakeLogRepo{}
	svc := NewCommentService(&fakeTx{repos: Repos{Comment: repo, Log: logs}}, repo)

	got, err := svc.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.GetByUUID(ctx, "m1")
	require.NoError(t, err)
	_, err = svc.List(ctx, 0, 10)
	require.NoError(t, err)
	byTicket, err := svc.ListByTicket(ctx, 5)
	require.NoError(t, err)
	require.Equal(t, repo.byTicket, byTicket)

	id, err := svc.Create(ctx, &domain.Comment{UUID: "m1"}, 3)
	require.NoError(t, err)
	require.Equal(t, uint(9), id)
	require.Equal(t, "comment.create|comment:m1|", logs.entries[0].Operation)

	require.NoError(t, svc.Update(ctx, &domain.Comment{UUID: "m1"}, 3))
	require.NoError(t, svc.Delete(ctx, 9, 3))
	require.Len(t, logs.entries, 3)
}

func TestCronService(t *testing.T) {
	ctx := context.Background()
	repo := &stubCronRepo{lastID: 5, got: &domain.Cron{UUID: "r1"}}
	logs := &fakeLogRepo{}
	svc := NewCronService(&fakeTx{repos: Repos{Cron: repo, Log: logs}}, repo)

	got, err := svc.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.GetByUUID(ctx, "r1")
	require.NoError(t, err)
	_, err = svc.List(ctx, 0, 10)
	require.NoError(t, err)

	id, err := svc.Create(ctx, &domain.Cron{UUID: "r1"}, 3)
	require.NoError(t, err)
	require.Equal(t, uint(5), id)
	require.Equal(t, "cron.create|cron:r1|", logs.entries[0].Operation)

	require.NoError(t, svc.Update(ctx, &domain.Cron{UUID: "r1"}, 3))
	require.NoError(t, svc.Delete(ctx, 5, 3))
	require.Len(t, logs.entries, 3)
}

func TestEdgeService(t *testing.T) {
	ctx := context.Background()
	repo := &stubEdgeRepo{lastID: 4, got: &domain.Edge{UUID: "e1"}}
	logs := &fakeLogRepo{}
	svc := NewEdgeService(&fakeTx{repos: Repos{Edge: repo, Log: logs}}, repo)

	got, err := svc.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.GetByUUID(ctx, "e1")
	require.NoError(t, err)
	_, err = svc.List(ctx, 0, 10)
	require.NoError(t, err)

	id, err := svc.Create(ctx, &domain.Edge{UUID: "e1"}, 3)
	require.NoError(t, err)
	require.Equal(t, uint(4), id)
	require.Equal(t, "edge.create|edge:e1|", logs.entries[0].Operation)

	require.NoError(t, svc.Update(ctx, &domain.Edge{UUID: "e1"}, 3))
	require.NoError(t, svc.Delete(ctx, 4, 3))
	require.Len(t, logs.entries, 3)
}

// TestEdgeServiceReplace checks that Replace deletes the old edges, assigns the
// ancestor, fills missing UUIDs and audits, all through the transaction.
func TestEdgeServiceReplace(t *testing.T) {
	ctx := context.Background()
	repo := &stubEdgeRepo{}
	logs := &fakeLogRepo{}
	svc := NewEdgeService(&fakeTx{repos: Repos{Edge: repo, Log: logs}}, repo)

	edges := []*domain.Edge{{DescendantID: 10}, {UUID: "kept", DescendantID: 11}}
	require.NoError(t, svc.Replace(ctx, uint(domain.EdgeTypeUserToRole), 42, edges, 3))

	require.Len(t, repo.removed, 1)
	require.Equal(t, uint(2), uint(repo.removed[0].Type))
	require.Equal(t, uint(42), repo.removed[0].AncestorID)
	require.Len(t, repo.created, 2)
	for _, e := range repo.created {
		require.Equal(t, uint(42), e.AncestorID)
		require.NotEmpty(t, e.UUID)
	}
	require.Equal(t, "kept", repo.created[1].UUID)
	require.Equal(t, "edge.replace|edge:2:42|", logs.entries[0].Operation)
}

func TestGlossaryService(t *testing.T) {
	ctx := context.Background()
	repo := &stubGlossaryRepo{got: &domain.Glossary{UUID: "g1"}}
	logs := &fakeLogRepo{}
	svc := NewGlossaryService(&fakeTx{repos: Repos{Glossary: repo, Log: logs}}, repo)

	got, err := svc.GetByUUID(ctx, "g1")
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.GetByKey(ctx, "grp", 1)
	require.NoError(t, err)
	_, err = svc.List(ctx, "", 10)
	require.NoError(t, err)

	require.NoError(t, svc.Create(ctx, &domain.Glossary{UUID: "g1"}, 3))
	require.Equal(t, "glossary.create|glossary:g1|", logs.entries[0].Operation)
	require.NoError(t, svc.Update(ctx, &domain.Glossary{UUID: "g1"}, 3))
	require.NoError(t, svc.DeleteByUUID(ctx, "g1", 3))
	require.NoError(t, svc.DeleteByKey(ctx, "grp", 1, 3))
	require.Len(t, logs.entries, 4)
}

func TestOptionService(t *testing.T) {
	ctx := context.Background()
	repo := &stubOptionRepo{got: &domain.Option{UUID: "o1"}}
	logs := &fakeLogRepo{}
	svc := NewOptionService(&fakeTx{repos: Repos{Option: repo, Log: logs}}, repo)

	got, err := svc.GetByUUID(ctx, "o1")
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.GetByName(ctx, "n")
	require.NoError(t, err)
	_, err = svc.List(ctx, "", 10)
	require.NoError(t, err)

	require.NoError(t, svc.Create(ctx, &domain.Option{UUID: "o1"}, 3))
	require.Equal(t, "option.create|option:o1|", logs.entries[0].Operation)
	require.NoError(t, svc.Update(ctx, &domain.Option{UUID: "o1"}, 3))
	require.NoError(t, svc.DeleteByUUID(ctx, "o1", 3))
	require.Len(t, logs.entries, 3)
}

func TestQueryService(t *testing.T) {
	ctx := context.Background()
	repo := &stubQueryRepo{lastID: 2, got: &domain.Query{UUID: "q1"}, byUser: []*domain.Query{{UUID: "q1"}}}
	logs := &fakeLogRepo{}
	svc := NewQueryService(&fakeTx{repos: Repos{Query: repo, Log: logs}}, repo)

	got, err := svc.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.GetByUUID(ctx, "q1")
	require.NoError(t, err)
	_, err = svc.List(ctx, 0, 10)
	require.NoError(t, err)
	byUser, err := svc.ListByUser(ctx, 5)
	require.NoError(t, err)
	require.Equal(t, repo.byUser, byUser)

	id, err := svc.Create(ctx, &domain.Query{UUID: "q1"}, 3)
	require.NoError(t, err)
	require.Equal(t, uint(2), id)
	require.Equal(t, "query.create|query:q1|", logs.entries[0].Operation)
	require.NoError(t, svc.Update(ctx, &domain.Query{UUID: "q1"}, 3))
	require.NoError(t, svc.Delete(ctx, 2, 3))
	require.Len(t, logs.entries, 3)
}

func TestRoleService(t *testing.T) {
	ctx := context.Background()
	repo := &stubRoleRepo{lastID: 6, got: &domain.Role{UUID: "r1"}}
	logs := &fakeLogRepo{}
	svc := NewRoleService(&fakeTx{repos: Repos{Role: repo, Log: logs}}, repo)

	got, err := svc.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.GetByUUID(ctx, "r1")
	require.NoError(t, err)
	_, err = svc.List(ctx, 0, 10)
	require.NoError(t, err)

	id, err := svc.Create(ctx, &domain.Role{UUID: "r1"}, 3)
	require.NoError(t, err)
	require.Equal(t, uint(6), id)
	require.Equal(t, "role.create|role:r1|", logs.entries[0].Operation)
	require.NoError(t, svc.Update(ctx, &domain.Role{UUID: "r1"}, 3))
	require.NoError(t, svc.Delete(ctx, 2, 3))
	require.Len(t, logs.entries, 3)
}

func TestRuleService(t *testing.T) {
	ctx := context.Background()
	repo := &stubRuleRepo{got: &domain.Rule{UUID: "r1"}}
	logs := &fakeLogRepo{}
	svc := NewRuleService(&fakeTx{repos: Repos{Rule: repo, Log: logs}}, repo)

	got, err := svc.GetByUUID(ctx, "r1")
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.GetByName(ctx, "n")
	require.NoError(t, err)
	_, err = svc.List(ctx, "", 10)
	require.NoError(t, err)

	require.NoError(t, svc.Create(ctx, &domain.Rule{UUID: "r1"}, 3))
	require.Equal(t, "rule.create|rule:r1|", logs.entries[0].Operation)
	require.NoError(t, svc.Update(ctx, &domain.Rule{UUID: "r1"}, 3))
	require.NoError(t, svc.DeleteByUUID(ctx, "r1", 3))
	require.Len(t, logs.entries, 3)
}

func TestRoleServiceErrors(t *testing.T) {
	ctx := context.Background()
	boom := errors.New("boom")

	repo := &stubRoleRepo{createErr: boom}
	svc := NewRoleService(&fakeTx{repos: Repos{Role: repo, Log: &fakeLogRepo{}}}, repo)
	_, err := svc.Create(ctx, &domain.Role{}, 1)
	require.ErrorIs(t, err, boom)

	repo = &stubRoleRepo{lastID: 1}
	svc = NewRoleService(&fakeTx{repos: Repos{Role: repo, Log: &fakeLogRepo{err: boom}}}, repo)
	_, err = svc.Create(ctx, &domain.Role{}, 1)
	require.ErrorIs(t, err, boom)

	repo = &stubRoleRepo{updateErr: boom}
	svc = NewRoleService(&fakeTx{repos: Repos{Role: repo, Log: &fakeLogRepo{}}}, repo)
	require.ErrorIs(t, svc.Update(ctx, &domain.Role{}, 1), boom)

	repo = &stubRoleRepo{deleteErr: boom}
	svc = NewRoleService(&fakeTx{repos: Repos{Role: repo, Log: &fakeLogRepo{}}}, repo)
	require.ErrorIs(t, svc.Delete(ctx, 1, 1), boom)
}

func TestMetricService(t *testing.T) {
	ctx := context.Background()
	repo := &stubMetricRepo{got: &domain.Metric{UUID: "m1"}}
	logs := &fakeLogRepo{}
	svc := NewMetricService(&fakeTx{repos: Repos{Metric: repo, Log: logs}}, repo, &fakeTicketRepo{}, &stubRuleHitRepo{})

	got, err := svc.GetByUUID(ctx, "m1")
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.GetByKey(ctx, domain.MetricKey{Name: "n"})
	require.NoError(t, err)
	_, err = svc.List(ctx, "", 10)
	require.NoError(t, err)

	require.NoError(t, svc.Create(ctx, &domain.Metric{UUID: "m1"}, 3))
	require.Equal(t, "metric.create|metric:m1|", logs.entries[0].Operation)
	require.NoError(t, svc.Update(ctx, &domain.Metric{UUID: "m1"}, 3))
	require.NoError(t, svc.DeleteByUUID(ctx, "m1", 3))
	require.Len(t, logs.entries, 3)
}

func TestTemplateService(t *testing.T) {
	ctx := context.Background()
	repo := &stubTemplateRepo{got: &domain.Template{UUID: "t1"}}
	logs := &fakeLogRepo{}
	svc := NewTemplateService(&fakeTx{repos: Repos{Template: repo, Log: logs}}, repo)

	got, err := svc.GetByUUID(ctx, "t1")
	require.NoError(t, err)
	require.Equal(t, repo.got, got)
	_, err = svc.List(ctx, "", 10)
	require.NoError(t, err)

	require.NoError(t, svc.Create(ctx, &domain.Template{UUID: "t1"}, 3))
	require.Equal(t, "template.create|template:t1|", logs.entries[0].Operation)
	require.NoError(t, svc.Update(ctx, &domain.Template{UUID: "t1"}, 3))
	require.NoError(t, svc.DeleteByUUID(ctx, "t1", 3))
	require.Len(t, logs.entries, 3)
}

// TestErrorOnlyServices checks the Create(error) shape shared by the entity
// services whose repository does not return a generated id.
func TestErrorOnlyServices(t *testing.T) {
	ctx := context.Background()
	boom := errors.New("boom")

	require.ErrorIs(t, NewGlossaryService(&fakeTx{repos: Repos{Glossary: &stubGlossaryRepo{createErr: boom}, Log: &fakeLogRepo{}}}, &stubGlossaryRepo{}).Create(ctx, &domain.Glossary{}, 1), boom)
	require.ErrorIs(t, NewOptionService(&fakeTx{repos: Repos{Option: &stubOptionRepo{createErr: boom}, Log: &fakeLogRepo{}}}, &stubOptionRepo{}).Create(ctx, &domain.Option{}, 1), boom)
	require.ErrorIs(t, NewRuleService(&fakeTx{repos: Repos{Rule: &stubRuleRepo{createErr: boom}, Log: &fakeLogRepo{}}}, &stubRuleRepo{}).Create(ctx, &domain.Rule{}, 1), boom)
	require.ErrorIs(t, NewMetricService(&fakeTx{repos: Repos{Metric: &stubMetricRepo{createErr: boom}, Log: &fakeLogRepo{}}}, &stubMetricRepo{}, &fakeTicketRepo{}, &stubRuleHitRepo{}).Create(ctx, &domain.Metric{}, 1), boom)
	require.ErrorIs(t, NewTemplateService(&fakeTx{repos: Repos{Template: &stubTemplateRepo{createErr: boom}, Log: &fakeLogRepo{}}}, &stubTemplateRepo{}).Create(ctx, &domain.Template{}, 1), boom)

	require.ErrorIs(t, NewGlossaryService(&fakeTx{repos: Repos{Glossary: &stubGlossaryRepo{}, Log: &fakeLogRepo{err: boom}}}, &stubGlossaryRepo{}).Create(ctx, &domain.Glossary{}, 1), boom)
	require.ErrorIs(t, NewOptionService(&fakeTx{repos: Repos{Option: &stubOptionRepo{}, Log: &fakeLogRepo{err: boom}}}, &stubOptionRepo{}).Create(ctx, &domain.Option{}, 1), boom)
	require.ErrorIs(t, NewRuleService(&fakeTx{repos: Repos{Rule: &stubRuleRepo{}, Log: &fakeLogRepo{err: boom}}}, &stubRuleRepo{}).Create(ctx, &domain.Rule{}, 1), boom)
	require.ErrorIs(t, NewMetricService(&fakeTx{repos: Repos{Metric: &stubMetricRepo{}, Log: &fakeLogRepo{err: boom}}}, &stubMetricRepo{}, &fakeTicketRepo{}, &stubRuleHitRepo{}).Create(ctx, &domain.Metric{}, 1), boom)
	require.ErrorIs(t, NewTemplateService(&fakeTx{repos: Repos{Template: &stubTemplateRepo{}, Log: &fakeLogRepo{err: boom}}}, &stubTemplateRepo{}).Create(ctx, &domain.Template{}, 1), boom)
}
