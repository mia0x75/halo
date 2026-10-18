package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/analyzer/ir"
	"github.com/mia0x75/halo/internal/domain"
)

func TestUserServiceReadsAndUpdate(t *testing.T) {
	ctx := context.Background()
	users := newFakeUserRepo()
	_, err := users.Create(ctx, &domain.User{UserID: 1, UUID: "u1", Email: "a@b.test"})
	require.NoError(t, err)
	logs := &fakeLogRepo{}
	svc := NewUserService(&fakeTx{repos: Repos{User: users, Log: logs}}, users)

	got, err := svc.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, "a@b.test", got.Email)
	byUUID, err := svc.GetByUUID(ctx, "u1")
	require.NoError(t, err)
	require.Equal(t, uint(1), byUUID.UserID)
	_, err = svc.List(ctx, 0, 10)
	require.NoError(t, err)

	require.NoError(t, svc.Update(ctx, &domain.User{UserID: 1, UUID: "u1", Email: "c@d.test"}, 7))
	require.Equal(t, "user.update|user:u1|", logs.entries[0].Operation)
}

func TestStatementTypeFromIR(t *testing.T) {
	cases := []struct {
		name string
		st   *ir.Statement
		want domain.StatementType
	}{
		{name: "nil", st: nil, want: domain.StatementTypeSelect},
		{name: "select", st: &ir.Statement{Kind: ir.KindSelect}, want: domain.StatementTypeSelect},
		{name: "insert", st: &ir.Statement{Kind: ir.KindInsert}, want: domain.StatementTypeInsert},
		{name: "update", st: &ir.Statement{Kind: ir.KindUpdate}, want: domain.StatementTypeUpdate},
		{name: "delete", st: &ir.Statement{Kind: ir.KindDelete}, want: domain.StatementTypeDelete},
		{name: "other", st: &ir.Statement{Kind: ir.KindOther}, want: domain.StatementTypeSelect},
		{name: "ddl-nil-info", st: &ir.Statement{Kind: ir.KindDDL}, want: domain.StatementTypeAlterTable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, statementTypeFromIR(tc.st))
		})
	}
}

func TestStatementTypeFromDDL(t *testing.T) {
	cases := map[string]domain.StatementType{
		"create_table":    domain.StatementTypeCreateTable,
		"drop_table":      domain.StatementTypeDropTable,
		"alter_table":     domain.StatementTypeAlterTable,
		"create_index":    domain.StatementTypeCreateIndex,
		"drop_index":      domain.StatementTypeDropIndex,
		"create_database": domain.StatementTypeCreateDatabase,
		"drop_database":   domain.StatementTypeDropDatabase,
		"truncate_table":  domain.StatementTypeTruncateTable,
		"create_view":     domain.StatementTypeCreateView,
		"CREATE_TABLE":    domain.StatementTypeCreateTable,
		"unknown_action":  domain.StatementTypeAlterTable,
	}
	for action, want := range cases {
		t.Run(action, func(t *testing.T) {
			require.Equal(t, want, statementTypeFromDDL(&ir.Statement{DDL: &ir.DDLInfo{Action: action}}))
		})
	}
}

func TestStatementStatusFromViolations(t *testing.T) {
	warn := domain.Violation{Level: uint8(domain.RuleLevelWarning)}
	fail := domain.Violation{Level: uint8(domain.RuleLevelError)}
	msg := domain.Violation{Level: uint8(domain.RuleLevelMessage)}

	require.Equal(t, domain.StatementStatusWaitingForMrv, statementStatusFromViolations(nil))
	require.Equal(t, domain.StatementStatusWaitingForMrv, statementStatusFromViolations([]domain.Violation{msg}))
	require.Equal(t, domain.StatementStatusVldWarning, statementStatusFromViolations([]domain.Violation{warn}))
	require.Equal(t, domain.StatementStatusVldFailure, statementStatusFromViolations([]domain.Violation{warn, fail}))
	require.Equal(t, domain.StatementStatusVldFailure, statementStatusFromViolations([]domain.Violation{fail, warn}))
}

func TestJoinViolations(t *testing.T) {
	require.Equal(t, "", joinViolations(nil))
	require.Equal(t, "a; b", joinViolations([]domain.Violation{{Message: "a"}, {Message: "b"}}))
}

func TestStatementServiceGetByUUID(t *testing.T) {
	ctx := context.Background()
	statements := newFakeStatementRepo()
	require.NoError(t, statements.Create(ctx, &domain.Statement{TicketID: 1, Sequence: 1, UUID: "s1"}))
	svc := NewStatementService(statements, nil, nil)

	got, err := svc.GetByUUID(ctx, "s1")
	require.NoError(t, err)
	require.Equal(t, uint(1), got.TicketID)

	_, err = svc.GetByUUID(ctx, "missing")
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestStatementServiceUpdateResultMissing(t *testing.T) {
	svc := NewStatementService(newFakeStatementRepo(), nil, nil)
	err := svc.UpdateResult(context.Background(), 1, 1, StatementResult{Results: "x"})
	require.ErrorIs(t, err, domain.ErrNotFound)
}
