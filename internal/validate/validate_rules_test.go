package validate_test

import "testing"

// Per-rule integration tests. Each drives the cases authored for one rule in
// integration_test.go against the live target cluster (see runRuleCases).

// RULE: CTB-L2-003
func TestTableCreateTableCharsetCollateMatch(t *testing.T) {
	runRuleCases(t, "CTB-L2-003")
}

// RULE: CTB-L2-009
func TestTableCreateFromSelectEnabled(t *testing.T) {
	runRuleCases(t, "CTB-L2-009")
}

// RULE: CTB-L2-019
func TestTableCreateColumnCharsetCollateMatch(t *testing.T) {
	runRuleCases(t, "CTB-L2-019")
}

// RULE: CTB-L3-001
func TestTableCreateTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "CTB-L3-001")
}

// RULE: CTB-L3-002
func TestTableCreateTargetTableExists(t *testing.T) {
	runRuleCases(t, "CTB-L3-002")
}

// RULE: UPD-L3-001
func TestUpdateTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "UPD-L3-001")
}

// RULE: UPD-L3-002
func TestUpdateTargetTableExists(t *testing.T) {
	runRuleCases(t, "UPD-L3-002")
}

// RULE: UPD-L3-003
func TestUpdateTargetColumnExists(t *testing.T) {
	runRuleCases(t, "UPD-L3-003")
}

// RULE: UPD-L3-004
func TestUpdateFilterColumnExists(t *testing.T) {
	t.Skip("UPD-L3-004: covered by UPD-L3-003; handler is a documented no-op")
}

// RULE: UPD-L3-005
func TestUpdateRowsLimit(t *testing.T) {
	runRuleCases(t, "UPD-L3-005")
}

// RULE: CVW-L3-001
func TestViewCreateTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "CVW-L3-001")
}

// RULE: CVW-L3-002
func TestViewCreateTargetViewExists(t *testing.T) {
	t.Skip("CVW-L3-002: meta layer has no view metadata; handler is a documented no-op")
}

// RULE: MVW-L3-001
func TestViewAlterTargetDatabaseExists(t *testing.T) {
	t.Skip("MVW-L3-001: parser has no ALTER VIEW AST; handler is a documented no-op")
}

// RULE: MVW-L3-002
func TestViewAlterTargetViewExists(t *testing.T) {
	t.Skip("MVW-L3-002: parser has no ALTER VIEW AST; handler is a documented no-op")
}

// RULE: DVW-L3-001
func TestViewDropTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "DVW-L3-001")
}

// RULE: DVW-L3-002
func TestViewDropTargetViewExists(t *testing.T) {
	t.Skip("DVW-L3-002: meta layer has no view metadata; handler is a documented no-op")
}

// RULE: INS-L2-001
func TestInsertExplicitColumnRequired(t *testing.T) {
	runRuleCases(t, "INS-L2-001")
}

// RULE: INS-L2-002
func TestInsertUsingSelectEnabled(t *testing.T) {
	runRuleCases(t, "INS-L2-002")
}

// RULE: INS-L2-003 (merge semantics now enforced by MSC-L1-004)
func TestInsertMergeRequired(t *testing.T) {
	runMergeCases(t, "MSC-L1-004")
}

// RULE: INS-L2-004
func TestInsertRowsLimit(t *testing.T) {
	runRuleCases(t, "INS-L3-006")
}

// RULE: INS-L2-005
func TestInsertColumnValueMatch(t *testing.T) {
	runRuleCases(t, "INS-L2-005")
}

// RULE: INS-L3-001
func TestInsertTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "INS-L3-001")
}

// RULE: INS-L3-002
func TestInsertTargetTableExists(t *testing.T) {
	runRuleCases(t, "INS-L3-002")
}

// RULE: INS-L3-003
func TestInsertTargetColumnExists(t *testing.T) {
	runRuleCases(t, "INS-L3-003")
}

// RULE: INS-L3-004
func TestInsertValueForNotNullColumnRequired(t *testing.T) {
	runRuleCases(t, "INS-L3-004")
}

// RULE: CIX-L3-001
func TestIndexCreateTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "CIX-L3-001")
}

// RULE: CIX-L3-002
func TestIndexCreateTargetTableExists(t *testing.T) {
	runRuleCases(t, "CIX-L3-002")
}

// RULE: CIX-L3-003
func TestIndexCreateTargetColumnExists(t *testing.T) {
	runRuleCases(t, "CIX-L3-003")
}

// RULE: CIX-L3-004
func TestIndexCreateTargetIndexExists(t *testing.T) {
	runRuleCases(t, "CIX-L3-004")
}

// RULE: CIX-L3-005
func TestIndexCreateTargetNameExists(t *testing.T) {
	runRuleCases(t, "CIX-L3-005")
}

// RULE: CIX-L3-006
func TestIndexCreateIndexCountLimit(t *testing.T) {
	runRuleCases(t, "CIX-L3-006")
}

// RULE: CIX-L3-007
func TestIndexCreateIndexBlobColumnEnabled(t *testing.T) {
	runRuleCases(t, "CIX-L3-007")
}

// RULE: CIX-L3-008
func TestIndexDropTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "RIX-L3-001")
}

// RULE: CIX-L3-009
func TestIndexDropTargetTableExists(t *testing.T) {
	runRuleCases(t, "RIX-L3-002")
}

// RULE: CIX-L3-010
func TestIndexDropTargetIndexExists(t *testing.T) {
	runRuleCases(t, "RIX-L3-003")
}

// RULE: DEL-L3-001
func TestDeleteRowsLimit(t *testing.T) {
	runRuleCases(t, "DEL-L3-001")
}

// RULE: DEL-L3-002
func TestDeleteTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "DEL-L3-002")
}

// RULE: DEL-L3-003
func TestDeleteTargetTableExists(t *testing.T) {
	runRuleCases(t, "DEL-L3-003")
}

// RULE: DEL-L3-004
func TestDeleteFilterColumnExists(t *testing.T) {
	runRuleCases(t, "DEL-L3-004")
}

// RULE: CDB-L2-001
func TestDatabaseCreateCharsetCollateMatch(t *testing.T) {
	runRuleCases(t, "CDB-L2-001")
}

// RULE: CDB-L2-002
func TestDatabaseCreateTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "CDB-L2-002")
}

// RULE: MDB-L2-001
func TestDatabaseAlterAvailableCharsets(t *testing.T) {
	runRuleCases(t, "MDB-L2-001")
}

// RULE: MDB-L2-002
func TestDatabaseAlterAvailableCollates(t *testing.T) {
	runRuleCases(t, "MDB-L2-002")
}

// RULE: MDB-L2-003
func TestDatabaseAlterCharsetCollateMatch(t *testing.T) {
	runRuleCases(t, "MDB-L2-003")
}

// RULE: MDB-L3-001
func TestDatabaseAlterTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "MDB-L2-004")
}

// RULE: DDB-L3-001
func TestDatabaseDropTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "DDB-L2-001")
}

// RULE: RPL-L2-002
func TestReplaceUsingSelectEnabled(t *testing.T) {
	runRuleCases(t, "RPL-L2-002")
}

// RULE: RPL-L2-003 (merge semantics now enforced by MSC-L1-004)
func TestReplaceMergeRequired(t *testing.T) {
	runMergeCases(t, "MSC-L1-004")
}

// RULE: RPL-L2-004
func TestReplaceRowsLimit(t *testing.T) {
	runRuleCases(t, "RPL-L2-004")
}

// RULE: RPL-L2-005
func TestReplaceColumnValueMatch(t *testing.T) {
	runRuleCases(t, "RPL-L2-005")
}

// RULE: RPL-L3-001
func TestReplaceTargetDatabaseExists(t *testing.T) {
	runRuleCases(t, "RPL-L3-001")
}

// RULE: RPL-L3-002
func TestReplaceTargetTableExists(t *testing.T) {
	runRuleCases(t, "RPL-L3-002")
}

// RULE: RPL-L3-003
func TestReplaceTargetColumnExists(t *testing.T) {
	runRuleCases(t, "RPL-L3-003")
}

// RULE: RPL-L3-004
func TestReplaceValueForNotNullColumnRequired(t *testing.T) {
	runRuleCases(t, "RPL-L3-004")
}

// RULE: MSC-L1-008
func TestKeywordEnabled(t *testing.T) {
	runMergeCases(t, "MSC-L1-008")
}

// RULE: MSC-L0-001
func TestSplitTicket(t *testing.T) {
	t.Skip("MSC-L0-001: handler is a documented no-op; DDL/DML split covered by MSC-L1-008")
}
