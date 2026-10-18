package domain

import (
	"errors"
	"fmt"
	"testing"
)

// The integer enum values below are persisted in the metadata schema, so they
// are part of the storage contract (ADR 0002 keeps the mapping in code). These
// tests fail loudly if a constant is ever renumbered.

func TestTicketStatusValues(t *testing.T) {
	want := map[TicketStatus]int{
		TicketStatusWaitingForVld: 1,
		TicketStatusVldFailure:    2,
		TicketStatusVldWarning:    3,
		TicketStatusWaitingForMrv: 4,
		TicketStatusMrvFailure:    5,
		TicketStatusLgtm:          6,
		TicketStatusDone:          7,
		TicketStatusExecFailure:   8,
		TicketStatusClosed:        9,
	}
	seen := make(map[int]struct{}, len(want))
	for got, value := range want {
		if int(got) != value {
			t.Errorf("ticket status %d: got %d, want %d", value, int(got), value)
		}
		if _, dup := seen[int(got)]; dup {
			t.Errorf("ticket status %d is not unique", got)
		}
		seen[int(got)] = struct{}{}
	}
}

func TestStatementTypeValues(t *testing.T) {
	want := map[StatementType]int{
		StatementTypeAlterTable:     1,
		StatementTypeAnalyzeTable:   2,
		StatementTypeBegin:          3,
		StatementTypeCommit:         4,
		StatementTypeCreateDatabase: 5,
		StatementTypeCreateIndex:    6,
		StatementTypeCreateTable:    7,
		StatementTypeCreateView:     8,
		StatementTypeCreateUser:     9,
		StatementTypeDelete:         10,
		StatementTypeDropDatabase:   11,
		StatementTypeDropIndex:      12,
		StatementTypeDropTable:      13,
		StatementTypeExplain:        14,
		StatementTypeInsert:         15,
		StatementTypeLoadData:       16,
		StatementTypeRollback:       17,
		StatementTypeSelect:         18,
		StatementTypeSet:            19,
		StatementTypeShow:           20,
		StatementTypeTruncateTable:  21,
		StatementTypeUpdate:         22,
		StatementTypeGrant:          23,
		StatementTypeRevoke:         24,
		StatementTypeDeallocate:     25,
		StatementTypeExecute:        26,
		StatementTypePrepare:        27,
		StatementTypeUse:            28,
	}
	seen := make(map[int]struct{}, len(want))
	for got, value := range want {
		if int(got) != value {
			t.Errorf("statement type %d: got %d, want %d", value, int(got), value)
		}
		if _, dup := seen[int(got)]; dup {
			t.Errorf("statement type %d is not unique", got)
		}
		seen[int(got)] = struct{}{}
	}
}

func TestStatementStatusValues(t *testing.T) {
	want := map[StatementStatus]int{
		StatementStatusPending:       1,
		StatementStatusVldFailure:    2,
		StatementStatusVldWarning:    3,
		StatementStatusWaitingForMrv: 4,
	}
	for got, value := range want {
		if int(got) != value {
			t.Errorf("statement status %d: got %d, want %d", value, int(got), value)
		}
	}
}

func TestRuleLevelValues(t *testing.T) {
	if RuleLevelMessage != 0 || RuleLevelError != 1 || RuleLevelWarning != 2 {
		t.Fatalf("rule levels drifted: %d %d %d", RuleLevelMessage, RuleLevelError, RuleLevelWarning)
	}
}

func TestUserStatusValues(t *testing.T) {
	if UserStatusNormal != 1 || UserStatusBlocked != 2 || UserStatusPending != 3 {
		t.Fatalf("user statuses drifted: %d %d %d", UserStatusNormal, UserStatusBlocked, UserStatusPending)
	}
}

func TestClusterStatusValues(t *testing.T) {
	if ClusterStatusNormal != 1 || ClusterStatusDisabled != 2 {
		t.Fatalf("cluster statuses drifted: %d %d", ClusterStatusNormal, ClusterStatusDisabled)
	}
}

func TestClusterTypeValues(t *testing.T) {
	if ClusterTypeMySQL != "mysql" || ClusterTypePostgres != "postgres" {
		t.Fatalf("cluster types drifted: %q %q", ClusterTypeMySQL, ClusterTypePostgres)
	}
}

func TestQueryTypeValues(t *testing.T) {
	if QueryTypeQuery != 1 || QueryTypeAnalyze != 2 || QueryTypeRewrite != 3 {
		t.Fatalf("query types drifted: %d %d %d", QueryTypeQuery, QueryTypeAnalyze, QueryTypeRewrite)
	}
}

func TestEdgeTypeValues(t *testing.T) {
	if EdgeTypeUserToReviewer != 1 || EdgeTypeUserToRole != 2 || EdgeTypeUserToCluster != 3 {
		t.Fatalf("edge types drifted: %d %d %d", EdgeTypeUserToReviewer, EdgeTypeUserToRole, EdgeTypeUserToCluster)
	}
}

func TestCronStatusValues(t *testing.T) {
	if CronStatusPending != "P" || CronStatusSuccess != "S" || CronStatusFailure != "F" {
		t.Fatalf("cron statuses drifted: %q %q %q", CronStatusPending, CronStatusSuccess, CronStatusFailure)
	}
}

// TestSentinelErrors proves the sentinels are distinct and wrap-preserving, so
// callers can use errors.Is through the repository layers.
func TestSentinelErrors(t *testing.T) {
	sentinels := []error{ErrNotFound, ErrInvalidInput, ErrConflict, ErrForbidden}
	for i, a := range sentinels {
		for j, b := range sentinels {
			if i != j && errors.Is(a, b) {
				t.Fatalf("sentinel %v unexpectedly matches %v", a, b)
			}
		}
		wrapped := fmt.Errorf("repo: %w", a)
		if !errors.Is(wrapped, a) {
			t.Fatalf("errors.Is lost the sentinel %v", a)
		}
	}
}
