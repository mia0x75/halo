package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

// Thin services share one repo+audit pattern, so the stubs below implement only
// the methods each service calls. Embedding the domain interface satisfies the
// rest of the contract without dead code, and every stub exposes just the knobs
// that drive the three branches under test: the write succeeds, the repo fails,
// or the audit insert fails.

type stubAvatarRepo struct {
	domain.AvatarRepository
	createErr, updateErr, deleteErr error
	lastID                          int64
	created                         *domain.Avatar
	got                             *domain.Avatar
	list                            []*domain.Avatar
}

func (s *stubAvatarRepo) Create(_ context.Context, a *domain.Avatar) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = a
	return s.lastID, nil
}
func (s *stubAvatarRepo) GetByID(context.Context, int64) (*domain.Avatar, error) { return s.got, nil }
func (s *stubAvatarRepo) GetByUUID(context.Context, string) (*domain.Avatar, error) {
	return s.got, nil
}
func (s *stubAvatarRepo) List(context.Context, int64, int) ([]*domain.Avatar, error) {
	return s.list, nil
}
func (s *stubAvatarRepo) Update(context.Context, *domain.Avatar) error { return s.updateErr }
func (s *stubAvatarRepo) Delete(context.Context, int64) error          { return s.deleteErr }

type stubCommentRepo struct {
	domain.CommentRepository
	createErr, updateErr, deleteErr error
	lastID                          int64
	got                             *domain.Comment
	byTicket                        []*domain.Comment
}

func (s *stubCommentRepo) Create(_ context.Context, c *domain.Comment) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	return s.lastID, nil
}
func (s *stubCommentRepo) GetByID(context.Context, int64) (*domain.Comment, error) {
	return s.got, nil
}
func (s *stubCommentRepo) GetByUUID(context.Context, string) (*domain.Comment, error) {
	return s.got, nil
}
func (s *stubCommentRepo) List(context.Context, int64, int) ([]*domain.Comment, error) {
	return nil, nil
}
func (s *stubCommentRepo) ListByTicket(context.Context, uint) ([]*domain.Comment, error) {
	return s.byTicket, nil
}
func (s *stubCommentRepo) Update(context.Context, *domain.Comment) error { return s.updateErr }
func (s *stubCommentRepo) Delete(context.Context, int64) error           { return s.deleteErr }

type stubCronRepo struct {
	domain.CronRepository
	createErr, updateErr, deleteErr error
	lastID                          int64
	got                             *domain.Cron
}

func (s *stubCronRepo) Create(_ context.Context, c *domain.Cron) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	return s.lastID, nil
}
func (s *stubCronRepo) GetByID(context.Context, int64) (*domain.Cron, error) { return s.got, nil }
func (s *stubCronRepo) GetByUUID(context.Context, string) (*domain.Cron, error) {
	return s.got, nil
}
func (s *stubCronRepo) List(context.Context, int64, int) ([]*domain.Cron, error) {
	return nil, nil
}
func (s *stubCronRepo) Update(context.Context, *domain.Cron) error { return s.updateErr }
func (s *stubCronRepo) Delete(context.Context, int64) error        { return s.deleteErr }

type stubEdgeRepo struct {
	domain.EdgeRepository
	createErr, updateErr, deleteErr error
	lastID                          int64
	got                             *domain.Edge
	removed                         []*domain.Edge
	created                         []*domain.Edge
}

func (s *stubCronRepo) ListSchedulable(context.Context) ([]*domain.Cron, error) {
	return nil, nil
}

func (s *stubEdgeRepo) Create(_ context.Context, e *domain.Edge) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.created = append(s.created, e)
	return s.lastID, nil
}
func (s *stubEdgeRepo) GetByID(context.Context, int64) (*domain.Edge, error) { return s.got, nil }
func (s *stubEdgeRepo) GetByUUID(context.Context, string) (*domain.Edge, error) {
	return s.got, nil
}
func (s *stubEdgeRepo) List(context.Context, int64, int) ([]*domain.Edge, error) { return nil, nil }
func (s *stubEdgeRepo) Update(context.Context, *domain.Edge) error               { return s.updateErr }
func (s *stubEdgeRepo) Delete(context.Context, int64) error                      { return s.deleteErr }
func (s *stubEdgeRepo) DeleteByTypeAncestor(_ context.Context, t uint, a uint) error {
	s.removed = append(s.removed, &domain.Edge{Type: domain.EdgeType(t), AncestorID: a})
	return s.deleteErr
}

type stubQueryRepo struct {
	domain.QueryRepository
	createErr, updateErr, deleteErr error
	lastID                          int64
	got                             *domain.Query
	byUser                          []*domain.Query
}

func (s *stubQueryRepo) Create(_ context.Context, q *domain.Query) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	return s.lastID, nil
}
func (s *stubQueryRepo) GetByID(context.Context, int64) (*domain.Query, error) { return s.got, nil }
func (s *stubQueryRepo) GetByUUID(context.Context, string) (*domain.Query, error) {
	return s.got, nil
}
func (s *stubQueryRepo) List(context.Context, int64, int) ([]*domain.Query, error) { return nil, nil }
func (s *stubQueryRepo) ListByUser(context.Context, uint) ([]*domain.Query, error) {
	return s.byUser, nil
}
func (s *stubQueryRepo) Update(context.Context, *domain.Query) error { return s.updateErr }
func (s *stubQueryRepo) Delete(context.Context, int64) error         { return s.deleteErr }

type stubRoleRepo struct {
	domain.RoleRepository
	createErr, updateErr, deleteErr error
	lastID                          int64
	got                             *domain.Role
}

func (s *stubRoleRepo) Create(_ context.Context, r *domain.Role) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	return s.lastID, nil
}
func (s *stubRoleRepo) GetByID(context.Context, int64) (*domain.Role, error) { return s.got, nil }
func (s *stubRoleRepo) GetByUUID(context.Context, string) (*domain.Role, error) {
	return s.got, nil
}
func (s *stubRoleRepo) List(context.Context, int64, int) ([]*domain.Role, error) { return nil, nil }
func (s *stubRoleRepo) Update(context.Context, *domain.Role) error               { return s.updateErr }
func (s *stubRoleRepo) Delete(context.Context, int64) error                      { return s.deleteErr }

// Error-only repos: Create returns no id.

type stubGlossaryRepo struct {
	domain.GlossaryRepository
	createErr, updateErr, deleteErr error
	got                             *domain.Glossary
}

func (s *stubGlossaryRepo) Create(context.Context, *domain.Glossary) error { return s.createErr }
func (s *stubGlossaryRepo) GetByUUID(context.Context, string) (*domain.Glossary, error) {
	return s.got, nil
}
func (s *stubGlossaryRepo) GetByKey(context.Context, string, uint8) (*domain.Glossary, error) {
	return s.got, nil
}
func (s *stubGlossaryRepo) List(context.Context, string, int) ([]*domain.Glossary, error) {
	return nil, nil
}
func (s *stubGlossaryRepo) Update(context.Context, *domain.Glossary) error { return s.updateErr }
func (s *stubGlossaryRepo) DeleteByUUID(context.Context, string) error     { return s.deleteErr }
func (s *stubGlossaryRepo) DeleteByKey(context.Context, string, uint8) error {
	return s.deleteErr
}

type stubOptionRepo struct {
	domain.OptionRepository
	createErr, updateErr, deleteErr error
	got                             *domain.Option
}

func (s *stubOptionRepo) Create(context.Context, *domain.Option) error { return s.createErr }
func (s *stubOptionRepo) GetByUUID(context.Context, string) (*domain.Option, error) {
	return s.got, nil
}
func (s *stubOptionRepo) GetByName(context.Context, string) (*domain.Option, error) {
	return s.got, nil
}
func (s *stubOptionRepo) List(context.Context, string, int) ([]*domain.Option, error) {
	return nil, nil
}
func (s *stubOptionRepo) Update(context.Context, *domain.Option) error { return s.updateErr }
func (s *stubOptionRepo) DeleteByUUID(context.Context, string) error   { return s.deleteErr }

type stubRuleRepo struct {
	domain.RuleRepository
	createErr, updateErr, deleteErr error
	got                             *domain.Rule
}

func (s *stubRuleRepo) Create(context.Context, *domain.Rule) error { return s.createErr }
func (s *stubRuleRepo) GetByUUID(context.Context, string) (*domain.Rule, error) {
	return s.got, nil
}
func (s *stubRuleRepo) GetByName(context.Context, string) (*domain.Rule, error) {
	return s.got, nil
}
func (s *stubRuleRepo) List(context.Context, string, int) ([]*domain.Rule, error) { return nil, nil }
func (s *stubRuleRepo) Update(context.Context, *domain.Rule) error                { return s.updateErr }
func (s *stubRuleRepo) DeleteByUUID(context.Context, string) error                { return s.deleteErr }

type stubMetricRepo struct {
	domain.MetricRepository
	createErr, updateErr, deleteErr error
	got                             *domain.Metric
}

func (s *stubMetricRepo) Create(context.Context, *domain.Metric) error { return s.createErr }
func (s *stubMetricRepo) GetByUUID(context.Context, string) (*domain.Metric, error) {
	return s.got, nil
}
func (s *stubMetricRepo) GetByKey(context.Context, domain.MetricKey) (*domain.Metric, error) {
	return s.got, nil
}
func (s *stubMetricRepo) List(context.Context, string, int) ([]*domain.Metric, error) {
	return nil, nil
}
func (s *stubMetricRepo) Update(context.Context, *domain.Metric) error { return s.updateErr }
func (s *stubMetricRepo) DeleteByUUID(context.Context, string) error  { return s.deleteErr }

type stubTemplateRepo struct {
	domain.TemplateRepository
	createErr, updateErr, deleteErr error
	got                             *domain.Template
}

func (s *stubTemplateRepo) Create(context.Context, *domain.Template) error { return s.createErr }
func (s *stubTemplateRepo) GetByUUID(context.Context, string) (*domain.Template, error) {
	return s.got, nil
}
func (s *stubTemplateRepo) List(context.Context, string, int) ([]*domain.Template, error) {
	return nil, nil
}
func (s *stubTemplateRepo) Update(context.Context, *domain.Template) error { return s.updateErr }
func (s *stubTemplateRepo) DeleteByUUID(context.Context, string) error     { return s.deleteErr }
