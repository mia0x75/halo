package validate

import (
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
	log "github.com/sirupsen/logrus"
)

// RuleFunc runs a single rule against one validator and subject.
// The subject is nil for cross-statement (Misc) rules, which inspect
// the whole context rather than one statement.
type RuleFunc func(v Validator, s *Subject, r *Rule)

// wrap adapts a concrete validator method to RuleFunc. The type
// parameter keeps the assertion checked at registration time.
func wrap[V Validator](fn func(V, *Subject, *Rule)) RuleFunc {
	return func(v Validator, s *Subject, r *Rule) {
		fn(v.(V), s, r)
	}
}

// ruleRegistry maps (vldr_group, mm_rules.func) to the handler. Group is the
// outer key because the same func name recurs across groups.
var ruleRegistry = map[uint16]map[string]RuleFunc{
	100: { // DatabaseCreateVldr
		"AvailableCharsets":             wrap((*DatabaseCreateVldr).createDBAvailableCharsets),
		"AvailableCollates":             wrap((*DatabaseCreateVldr).createDBAvailableCollates),
		"CharsetCollateMustMatch":       wrap((*DatabaseCreateVldr).createDBCharsetCollateMustMatch),
		"DatabaseNameLowerCaseRequired": wrap((*DatabaseCreateVldr).createDBDatabaseNameLowerCaseRequired),
		"DatabaseNameMaxLength":         wrap((*DatabaseCreateVldr).createDBDatabaseNameMaxLength),
		"DatabaseNameQualified":         wrap((*DatabaseCreateVldr).createDBDatabaseNameQualified),
		"TargetDatabaseDoesNotExist":    wrap((*DatabaseCreateVldr).createDBTargetDatabaseDoesNotExist),
	},
	101: { // DatabaseAlterVldr
		"AvailableCharsets":          wrap((*DatabaseAlterVldr).alterDBAvailableCharsets),
		"AvailableCollates":          wrap((*DatabaseAlterVldr).alterDBAvailableCollates),
		"CharsetCollateMustMatch":    wrap((*DatabaseAlterVldr).alterDBCharsetCollateMustMatch),
		"TargetDatabaseDoesNotExist": wrap((*DatabaseAlterVldr).alterDBTargetDatabaseDoesNotExist),
	},
	102: { // DatabaseDropVldr
		"TargetDatabaseDoesNotExist": wrap((*DatabaseDropVldr).dropDBTargetDatabaseDoesNotExist),
	},
	110: { // TableCreateVldr
		"AutoIncColumnDuplicate":           wrap((*TableCreateVldr).createTableAutoIncColumnDuplicate),
		"AvailableCharsets":                wrap((*TableCreateVldr).createTableAvailableCharsets),
		"AvailableCollates":                wrap((*TableCreateVldr).createTableAvailableCollates),
		"AvailableEngines":                 wrap((*TableCreateVldr).createTableAvailableEngines),
		"ColumnAutoIncAvailableTypes":      wrap((*TableCreateVldr).createTableColumnAutoIncAvailableTypes),
		"ColumnAutoIncMustPrimaryKey":      wrap((*TableCreateVldr).createTableColumnAutoIncMustPrimaryKey),
		"ColumnAutoIncUnsignedRequired":    wrap((*TableCreateVldr).createTableColumnAutoIncUnsignedRequired),
		"ColumnAvailableCharsets":          wrap((*TableCreateVldr).createTableColumnAvailableCharsets),
		"ColumnAvailableCollates":          wrap((*TableCreateVldr).createTableColumnAvailableCollates),
		"ColumnCharsetCollateMustMatch":    wrap((*TableCreateVldr).createTableColumnCharsetCollateMustMatch),
		"ColumnCommentRequired":            wrap((*TableCreateVldr).createTableColumnCommentRequired),
		"ColumnNameDuplicate":              wrap((*TableCreateVldr).createTableColumnNameDuplicate),
		"ColumnNameLowerCaseRequired":      wrap((*TableCreateVldr).createTableColumnNameLowerCaseRequired),
		"ColumnNameMaxLength":              wrap((*TableCreateVldr).createTableColumnNameMaxLength),
		"ColumnNameQualified":              wrap((*TableCreateVldr).createTableColumnNameQualified),
		"ColumnNotNullWithDefaultRequired": wrap((*TableCreateVldr).createTableColumnNotNullWithDefaultRequired),
		"ColumnTypesDoesNotExpect":         wrap((*TableCreateVldr).createTableColumnTypesDoesNotExpect),
		"CreateTableFromSelectNotAllowed":  wrap((*TableCreateVldr).createTableCreateTableFromSelectNotAllowed),
		"CreateTableUseLikeNotAllowed":     wrap((*TableCreateVldr).createTableCreateTableUseLikeNotAllowed),
		"ForeignKeyNameExplicit":           wrap((*TableCreateVldr).createTableForeignKeyNameExplicit),
		"ForeignKeyNameLowerCaseRequired":  wrap((*TableCreateVldr).createTableForeignKeyNameLowerCaseRequired),
		"ForeignKeyNameMaxLength":          wrap((*TableCreateVldr).createTableForeignKeyNameMaxLength),
		"ForeignKeyNamePrefixRequired":     wrap((*TableCreateVldr).createTableForeignKeyNamePrefixRequired),
		"ForeignKeyNameQualified":          wrap((*TableCreateVldr).createTableForeignKeyNameQualified),
		"ForeignKeyNotAllowed":             wrap((*TableCreateVldr).createTableForeignKeyNotAllowed),
		"IndexNameExplicit":                wrap((*TableCreateVldr).createTableIndexNameExplicit),
		"IndexNameLowerCaseRequired":       wrap((*TableCreateVldr).createTableIndexNameLowerCaseRequired),
		"IndexNameMaxLength":               wrap((*TableCreateVldr).createTableIndexNameMaxLength),
		"IndexNamePrefixRequired":          wrap((*TableCreateVldr).createTableIndexNamePrefixRequired),
		"IndexNameQualified":               wrap((*TableCreateVldr).createTableIndexNameQualified),
		"MaxAllowedColumnCount":            wrap((*TableCreateVldr).createTableMaxAllowedColumnCount),
		"MaxAllowedIndexColumnCount":       wrap((*TableCreateVldr).createTableMaxAllowedIndexColumnCount),
		"MaxAllowedIndexCount":             wrap((*TableCreateVldr).createTableMaxAllowedIndexCount),
		"MaxAllowedTimestampCount":         wrap((*TableCreateVldr).createTableMaxAllowedTimestampCount),
		"PrimaryKeyDuplicate":              wrap((*TableCreateVldr).createTablePrimaryKeyDuplicate),
		"PrimaryKeyNameExplicit":           wrap((*TableCreateVldr).createTablePrimaryKeyNameExplicit),
		"PrimaryKeyNameQualified":          wrap((*TableCreateVldr).createTablePrimaryKeyNameQualified),
		"PrimaryKeyRequired":               wrap((*TableCreateVldr).createTablePrimaryKeyRequired),
		"PrimryKeyLowerCaseRequired":       wrap((*TableCreateVldr).createTablePrimryKeyLowerCaseRequired),
		"PrimryKeyMaxLength":               wrap((*TableCreateVldr).createTablePrimryKeyMaxLength),
		"PrimryKeyPrefixRequired":          wrap((*TableCreateVldr).createTablePrimryKeyPrefixRequired),
		"TableCharsetCollateMustMatch":     wrap((*TableCreateVldr).createTableTableCharsetCollateMustMatch),
		"TableCommentRequired":             wrap((*TableCreateVldr).createTableTableCommentRequired),
		"TableNameLowerCaseRequired":       wrap((*TableCreateVldr).createTableTableNameLowerCaseRequired),
		"TableNameMaxLength":               wrap((*TableCreateVldr).createTableTableNameMaxLength),
		"TableNameQualified":               wrap((*TableCreateVldr).createTableTableNameQualified),
		"TargetDatabaseDoesNotExist":       wrap((*TableCreateVldr).createTableTargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":          wrap((*TableCreateVldr).createTableTargetTableDoesNotExist),
		"UniqueNameExplicit":               wrap((*TableCreateVldr).createTableUniqueNameExplicit),
		"UniqueNameLowerCaseRequired":      wrap((*TableCreateVldr).createTableUniqueNameLowerCaseRequired),
		"UniqueNameMaxLength":              wrap((*TableCreateVldr).createTableUniqueNameMaxLength),
		"UniqueNamePrefixRequired":         wrap((*TableCreateVldr).createTableUniqueNamePrefixRequired),
		"UniqueNameQualified":              wrap((*TableCreateVldr).createTableUniqueNameQualified),
	},
	120: { // TableAlterVldr
		"AvailableCharsets":                  wrap((*TableAlterVldr).alterTableAvailableCharsets),
		"AvailableCollates":                  wrap((*TableAlterVldr).alterTableAvailableCollates),
		"AvailableEngines":                   wrap((*TableAlterVldr).alterTableAvailableEngines),
		"ColumnAvailableCharsets":            wrap((*TableAlterVldr).alterTableColumnAvailableCharsets),
		"ColumnAvailableCollates":            wrap((*TableAlterVldr).alterTableColumnAvailableCollates),
		"ColumnCharsetCollateMustMatch":      wrap((*TableAlterVldr).alterTableColumnCharsetCollateMustMatch),
		"ColumnCommentRequired":              wrap((*TableAlterVldr).alterTableColumnCommentRequired),
		"CommentRequired":                    wrap((*TableAlterVldr).alterTableCommentRequired),
		"ColumnNameDoesNotExist":             wrap((*TableAlterVldr).alterTableColumnNameDoesNotExist),
		"ColumnNameDuplicate":                wrap((*TableAlterVldr).alterTableColumnNameDuplicate),
		"ColumnNameLowerCaseRequired":        wrap((*TableAlterVldr).alterTableColumnNameLowerCaseRequired),
		"ColumnNameMaxLength":                wrap((*TableAlterVldr).alterTableColumnNameMaxLength),
		"ColumnNameQualified":                wrap((*TableAlterVldr).alterTableColumnNameQualified),
		"ColumnNotNullWithDefaultRequired":   wrap((*TableAlterVldr).alterTableColumnNotNullWithDefaultRequired),
		"ColumnUnwantedTypes":                wrap((*TableAlterVldr).alterTableColumnUnwantedTypes),
		"ForeignKeyDoesNotExist":             wrap((*TableAlterVldr).alterTableForeignKeyDoesNotExist),
		"ForeignKeyNameExplicit":             wrap((*TableAlterVldr).alterTableForeignKeyNameExplicit),
		"ForeignKeyNameLowerCaseRequired":    wrap((*TableAlterVldr).alterTableForeignKeyNameLowerCaseRequired),
		"ForeignKeyNameMaxLength":            wrap((*TableAlterVldr).alterTableForeignKeyNameMaxLength),
		"ForeignKeyNamePrefixRequired":       wrap((*TableAlterVldr).alterTableForeignKeyNamePrefixRequired),
		"ForeignKeyNameQualified":            wrap((*TableAlterVldr).alterTableForeignKeyNameQualified),
		"ForeignKeyNotAllowed":               wrap((*TableAlterVldr).alterTableForeignKeyNotAllowed),
		"FullTextIndexNameExplicit":          wrap((*TableAlterVldr).alterTableFullTextIndexExplicit),
		"FullTextIndexNameLowerCaseRequired": wrap((*TableAlterVldr).alterTableFullTextIndexNameLowerCaseRequired),
		"FullTextIndexNameMaxLength":         wrap((*TableAlterVldr).alterTableFullTextIndexNameMaxLength),
		"FullTextIndexNamePrefixRequired":    wrap((*TableAlterVldr).alterTableFullTextIndexNamePrefixRequired),
		"FullTextIndexNameQualified":         wrap((*TableAlterVldr).alterTableFullTextIndexNameQualified),
		"FullTextIndexNotAllowed":            wrap((*TableAlterVldr).alterTableFullTextIndexNotAllowed),
		"IndexColumnDoesNotExist":            wrap((*TableAlterVldr).alterTableIndexColumnDoesNotExist),
		"IndexDoesNotExist":                  wrap((*TableAlterVldr).alterTableIndexDoesNotExist),
		"IndexNameDuplicate":                 wrap((*TableAlterVldr).alterTableIndexNameDuplicate),
		"IndexNameExplicit":                  wrap((*TableAlterVldr).alterTableIndexNameExplicit),
		"IndexNameLowerCaseRequired":         wrap((*TableAlterVldr).alterTableIndexNameLowerCaseRequired),
		"IndexNameMaxLength":                 wrap((*TableAlterVldr).alterTableIndexNameMaxLength),
		"IndexNamePrefixRequired":            wrap((*TableAlterVldr).alterTableIndexNamePrefixRequired),
		"IndexNameQualified":                 wrap((*TableAlterVldr).alterTableIndexNameQualified),
		"IndexOnBlobColumnNotAllowed":        wrap((*TableAlterVldr).alterTableIndexOnBlobColumnNotAllowed),
		"IndexOverlayNotAllowed":             wrap((*TableAlterVldr).alterTableIndexOverlayNotAllowed),
		"MaxAllowedColumnCount":              wrap((*TableAlterVldr).alterTableMaxAllowedColumnCount),
		"MaxAllowedIndexColumnCount":         wrap((*TableAlterVldr).alterTableMaxAllowedIndexColumnCount),
		"IndexCountLimit":                    wrap((*TableAlterVldr).alterTableIndexCountLimit),
		"MaxAllowedTimestampCount":           wrap((*TableAlterVldr).alterTableMaxAllowedTimestampCount),
		"NewTableNameLowerCaseRequired":      wrap((*TableAlterVldr).alterTableNewTableNameLowerCaseRequired),
		"NewTableNameMaxLength":              wrap((*TableAlterVldr).alterTableNewTableNameMaxLength),
		"NewTableNameQualified":              wrap((*TableAlterVldr).alterTableNewTableNameQualified),
		"PositionColumnDoesNotExist":         wrap((*TableAlterVldr).alterTablePositionColumnDoesNotExist),
		"PrimaryKeyDoesNotExist":             wrap((*TableAlterVldr).alterTablePrimaryKeyDoesNotExist),
		"TableCharsetCollateMustMatch":       wrap((*TableAlterVldr).alterTableTableCharsetCollateMustMatch),
		"TargetDatabaseDoesNotExist":         wrap((*TableAlterVldr).alterTableTargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":            wrap((*TableAlterVldr).alterTableTargetTableDoesNotExist),
		"TargetTableDuplicate":               wrap((*TableAlterVldr).alterTableTargetTableDuplicate),
		"UniqueNameExplicit":                 wrap((*TableAlterVldr).alterTableUniqueNameExplicit),
		"UniqueNameLowerCaseRequired":        wrap((*TableAlterVldr).alterTableUniqueNameLowerCaseRequired),
		"UniqueNameMaxLength":                wrap((*TableAlterVldr).alterTableUniqueNameMaxLength),
		"UniqueNamePrefixRequired":           wrap((*TableAlterVldr).alterTableUniqueNamePrefixRequired),
		"UniqueNameQualified":                wrap((*TableAlterVldr).alterTableUniqueNameQualified),
	},
	130: { // TableRenameVldr
		"SourceDatabaseDoesNotExist":       wrap((*TableRenameVldr).renameTableSourceDatabaseDoesNotExist),
		"SourceTableDoesNotExist":          wrap((*TableRenameVldr).renameTableSourceTableDoesNotExist),
		"TablesIdentical":                  wrap((*TableRenameVldr).renameTableTablesIdentical),
		"TargetTableDoesNotExist":          wrap((*TableRenameVldr).renameTableTargetTableDoesNotExist),
		"TargetTableNameLowerCaseRequired": wrap((*TableRenameVldr).renameTableTargetTableNameLowerCaseRequired),
		"TargetTableNameMaxLength":         wrap((*TableRenameVldr).renameTableTargetTableNameMaxLength),
		"TargetTableNameQualified":         wrap((*TableRenameVldr).renameTableTargetTableNameQualified),
	},
	140: { // TableDropVldr
		"TableDropTargetDatabaseDoesNotExist": wrap((*TableDropVldr).dropTableTargetDatabaseDoesNotExist),
		"TableDropTargetTableDoesNotExist":    wrap((*TableDropVldr).dropTableTargetTableDoesNotExist),
	},
	150: { // IndexCreateVldr
		"IndexColumnDuplicate":        wrap((*IndexCreateVldr).createIndexIndexColumnDuplicate),
		"IndexNameDuplicate":          wrap((*IndexCreateVldr).createIndexIndexNameDuplicate),
		"IndexNameLowerCaseRequired":  wrap((*IndexCreateVldr).createIndexIndexNameLowerCaseRequired),
		"IndexNameMaxLength":          wrap((*IndexCreateVldr).createIndexIndexNameMaxLength),
		"IndexNamePrefixRequired":     wrap((*IndexCreateVldr).createIndexIndexNamePrefixRequired),
		"IndexNameQualified":          wrap((*IndexCreateVldr).createIndexIndexNameQualified),
		"IndexOnBlobColumnNotAllowed": wrap((*IndexCreateVldr).createIndexIndexOnBlobColumnNotAllowed),
		"IndexOverlayNotAllowed":      wrap((*IndexCreateVldr).createIndexIndexOverlayNotAllowed),
		"MaxAllowedIndexColumnCount":  wrap((*IndexCreateVldr).createIndexMaxAllowedIndexColumnCount),
		"MaxAllowedIndexCount":        wrap((*IndexCreateVldr).createIndexMaxAllowedIndexCount),
		"TargetColumnDoesNotExist":    wrap((*IndexCreateVldr).createIndexTargetColumnDoesNotExist),
		"TargetDatabaseDoesNotExist":  wrap((*IndexCreateVldr).createIndexTargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":     wrap((*IndexCreateVldr).createIndexTargetTableDoesNotExist),
	},
	151: { // IndexDropVldr
		"TargetDatabaseDoesNotExist": wrap((*IndexDropVldr).renameIndexTargetDatabaseDoesNotExist),
		"TargetIndexDoesNotExist":    wrap((*IndexDropVldr).renameIndexTargetIndexDoesNotExist),
		"TargetTableDoesNotExist":    wrap((*IndexDropVldr).renameIndexTargetTableDoesNotExist),
	},
	160: { // InsertVldr
		"ColumnsValuesMustMatch":        wrap((*InsertVldr).insertColumnsValuesMustMatch),
		"ExplicitColumnRequired":        wrap((*InsertVldr).insertExplicitColumnRequired),
		"RowsLimit":                     wrap((*InsertVldr).insertRowsLimit),
		"TargetColumnDoesNotExist":      wrap((*InsertVldr).insertTargetColumnDoesNotExist),
		"TargetDatabaseDoesNotExist":    wrap((*InsertVldr).insertTargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":       wrap((*InsertVldr).insertTargetTableDoesNotExist),
		"UsingSelectNotAllowed":         wrap((*InsertVldr).insertUsingSelectNotAllowed),
		"ValueForNotNullColumnRequired": wrap((*InsertVldr).insertValueForNotNullColumnRequired),
	},
	161: { // ReplaceVldr
	},
	162: { // UpdateVldr
		"MaxAllowedUpdateRows":       wrap((*UpdateVldr).updateMaxAllowedUpdateRows),
		"TargetColumnDoesNotExist":   wrap((*UpdateVldr).updateTargetColumnDoesNotExist),
		"TargetDatabaseDoesNotExist": wrap((*UpdateVldr).updateTargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":    wrap((*UpdateVldr).updateTargetTableDoesNotExist),
		"WithoutWhereNotAllowed":     wrap((*UpdateVldr).updateWithoutWhereNotAllowed),
	},
	163: { // DeleteVldr
		"MaxAllowedDeleteRows":            wrap((*DeleteVldr).deleteMaxAllowedDeleteRows),
		"TargetColumnDoesNotDoesNotExist": wrap((*DeleteVldr).deleteTargetColumnDoesNotDoesNotExist),
		"TargetDatabaseDoesNotExist":      wrap((*DeleteVldr).deleteTargetDatabaseDoesNotExist),
		"TargetTableDoesNotDoesNotExist":  wrap((*DeleteVldr).deleteTargetTableDoesNotDoesNotExist),
		"WithoutWhereNotAllowed":          wrap((*DeleteVldr).deleteWithoutWhereNotAllowed),
	},
	170: { // SelectVldr
		"ReturnBlobOrTextNotAllowed":      wrap((*SelectVldr).selectReturnBlobOrTextNotAllowed),
		"TargetColumnDoesNotDoesNotExist": wrap((*SelectVldr).selectTargetColumnDoesNotDoesNotExist),
		"TargetDatabaseDoesNotExist":      wrap((*SelectVldr).selectTargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":         wrap((*SelectVldr).selectTargetTableDoesNotExist),
		"UseExplicitLockNotAllowed":       wrap((*SelectVldr).selectUseExplicitLockNotAllowed),
		"UseWildcardNotAllowed":           wrap((*SelectVldr).selectUseWildcardNotAllowed),
		"WithoutLimitNotAllowed":          wrap((*SelectVldr).selectWithoutLimitNotAllowed),
		"WithoutWhereNotAllowed":          wrap((*SelectVldr).selectWithoutWhereNotAllowed),
	},
	180: { // ViewCreateVldr
		"TargetDatabaseDoesNotExist": wrap((*ViewCreateVldr).createViewTargetDatabaseDoesNotExist),
		"TargetViewDoesNotExist":     wrap((*ViewCreateVldr).createViewTargetViewDoesNotExist),
		"ViewNameLowerCaseRequired":  wrap((*ViewCreateVldr).createViewViewNameLowerCaseRequired),
		"ViewNameMaxLength":          wrap((*ViewCreateVldr).createViewViewNameMaxLength),
		"ViewNamePrefixRequired":     wrap((*ViewCreateVldr).createViewViewNamePrefixRequired),
		"ViewNameQualified":          wrap((*ViewCreateVldr).createViewViewNameQualified),
	},
	181: { // ViewAlterVldr
		"TargetDatabaseDoesNotExist": wrap((*ViewAlterVldr).alterViewTargetDatabaseDoesNotExist),
		"TargetViewDoesNotExist":     wrap((*ViewAlterVldr).alterViewTargetViewDoesNotExist),
	},
	182: { // ViewDropVldr
		"TargetDatabaseDoesNotExist": wrap((*ViewDropVldr).dropViewTargetDatabaseDoesNotExist),
		"TargetViewDoesNotExist":     wrap((*ViewDropVldr).dropViewTargetViewDoesNotExist),
	},
	190: { // FuncCreateVldr
		"FuncNameLowerCaseRequired":  wrap((*FuncCreateVldr).createFuncFuncNameLowerCaseRequired),
		"FuncNameMaxLength":          wrap((*FuncCreateVldr).createFuncFuncNameMaxLength),
		"FuncNamePrefixRequired":     wrap((*FuncCreateVldr).createFuncFuncNamePrefixRequired),
		"FuncNameQuilified":          wrap((*FuncCreateVldr).createFuncFuncNameQuilified),
		"TargetDatabaseDoesNotExist": wrap((*FuncCreateVldr).createFuncTargetDatabaseDoesNotExist),
		"TargetFuncDoesNotExist":     wrap((*FuncCreateVldr).createFuncTargetFuncDoesNotExist),
	},
	191: { // FuncAlterVldr
		"TargetDatabaseDoesNotExist": wrap((*FuncAlterVldr).alterFuncTargetDatabaseDoesNotExist),
		"TargetFuncDoesNotExist":     wrap((*FuncAlterVldr).alterFuncTargetFuncDoesNotExist),
	},
	192: { // FuncDropVldr
		"TargetDatabaseDoesNotExist": wrap((*FuncDropVldr).dropFuncTargetDatabaseDoesNotExist),
		"TargetFuncDoesNotExist":     wrap((*FuncDropVldr).dropFuncTargetFuncDoesNotExist),
	},
	200: { // TriggerCreateVldr
		"TargetDatabaseDoesNotExist":   wrap((*TriggerCreateVldr).createTriggerTargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":      wrap((*TriggerCreateVldr).createTriggerTargetTableDoesNotExist),
		"TargetTriggerDoesNotExist":    wrap((*TriggerCreateVldr).createTriggerTargetTriggerDoesNotExist),
		"TriggerNameLowerCaseRequired": wrap((*TriggerCreateVldr).createTriggerTriggerNameLowerCaseRequired),
		"TriggerNameMaxLength":         wrap((*TriggerCreateVldr).createTriggerTriggerNameMaxLength),
		"TriggerNameQualified":         wrap((*TriggerCreateVldr).createTriggerTriggerNameQualified),
		"TriggerPrefixRequired":        wrap((*TriggerCreateVldr).createTriggerTriggerPrefixRequired),
	},
	201: { // TriggerAlterVldr
		"TargetDatabaseDoesNotExist": wrap((*TriggerAlterVldr).alterTriggerTargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":    wrap((*TriggerAlterVldr).alterTriggerTargetTableDoesNotExist),
		"TargetTriggerDoesNotExist":  wrap((*TriggerAlterVldr).alterTriggerTargetTriggerDoesNotExist),
	},
	202: { // TriggerDropVldr
		"TargetDatabaseDoesNotExist": wrap((*TriggerDropVldr).dropTriggerTargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":    wrap((*TriggerDropVldr).dropTriggerTargetTableDoesNotExist),
		"TargetTriggerDoesNotExist":  wrap((*TriggerDropVldr).dropTriggerTargetTriggerDoesNotExist),
	},
	210: { // EventCreateVldr
		"EventNameLowerCaseRequired": wrap((*EventCreateVldr).createEventEventNameLowerCaseRequired),
		"EventNameMaxLength":         wrap((*EventCreateVldr).createEventEventNameMaxLength),
		"EventNamePrefixRequired":    wrap((*EventCreateVldr).createEventEventNamePrefixRequired),
		"EventNameQualified":         wrap((*EventCreateVldr).createEventEventNameQualified),
		"TargetDatabaseDoesNotExist": wrap((*EventCreateVldr).createEventTargetDatabaseDoesNotExist),
		"TargetEventDoesNotExist":    wrap((*EventCreateVldr).createEventTargetEventDoesNotExist),
	},
	211: { // EventAlterVldr
		"TargetDatabaseDoesNotExist": wrap((*EventAlterVldr).alterEventTargetDatabaseDoesNotExist),
		"TargetEventDoesNotExist":    wrap((*EventAlterVldr).alterEventTargetEventDoesNotExist),
	},
	212: { // EventDropVldr
		"TargetDatabaseDoesNotExist": wrap((*EventDropVldr).dropEventTargetDatabaseDoesNotExist),
		"TargetEventDoesNotExist":    wrap((*EventDropVldr).dropEventTargetEventDoesNotExist),
	},
	220: { // ProcCreateVldr
		"ProcNameLowerCaseRequired":  wrap((*ProcCreateVldr).createProcProcNameLowerCaseRequired),
		"ProcNameMaxLength":          wrap((*ProcCreateVldr).createProcProcNameMaxLength),
		"ProcNamePrefixRequired":     wrap((*ProcCreateVldr).createProcProcNamePrefixRequired),
		"ProcNameQualified":          wrap((*ProcCreateVldr).createProcProcNameQualified),
		"TargetDatabaseDoesNotExist": wrap((*ProcCreateVldr).createProcTargetDatabaseDoesNotExist),
		"TargetProcDoesNotExist":     wrap((*ProcCreateVldr).createProcTargetProcDoesNotExist),
	},
	221: { // ProcAlterVldr
		"TargetDatabaseDoesNotExist": wrap((*ProcAlterVldr).alterProcTargetDatabaseDoesNotExist),
		"TargetProcDoesNotExist":     wrap((*ProcAlterVldr).alterProcTargetProcDoesNotExist),
	},
	222: { // ProcDropVldr
		"TargetDatabaseDoesNotExist": wrap((*ProcDropVldr).dropProcTargetDatabaseDoesNotExist),
		"TargetProcDoesNotExist":     wrap((*ProcDropVldr).dropProcTargetProcDoesNotExist),
	},
	230: { // MiscVldr
		"FlushTableProhibited":    wrap((*MiscVldr).miscFlushTableProhibited),
		"KillProhibited":          wrap((*MiscVldr).miscKillProhibited),
		"LockTableProhibited":     wrap((*MiscVldr).miscLockTableProhibited),
		"MergeRequired":           wrap((*MiscVldr).miscMergeRequired),
		"PurgeLogsProhibited":     wrap((*MiscVldr).miscPurgeLogsProhibited),
		"SplitRequired":           wrap((*MiscVldr).miscSplitRequired),
		"TruncateTableProhibited": wrap((*MiscVldr).miscTruncateTableProhibited),
		"UnlockTableProhibited":   wrap((*MiscVldr).miscUnlockTableProhibited),
	},
}

// lookup resolves a rule handler by group and func name.
func lookup(group uint16, name string) (RuleFunc, bool) {
	g, ok := ruleRegistry[group]
	if !ok {
		return nil, false
	}
	fn, ok := g[name]
	return fn, ok
}

// ValidateRegistry reports rules whose (vldr_group, func) has no handler. It
// logs instead of failing so startup is not blocked by the known mismatches
// that the old reflection dispatch silently skipped anyway.
func ValidateRegistry(rules []*domain.Rule) {
	var missing []string
	seen := map[uint8]map[string]bool{}
	for _, r := range rules {
		if r == nil {
			continue
		}
		if seen[r.VldrGroup] == nil {
			seen[r.VldrGroup] = map[string]bool{}
		}
		if seen[r.VldrGroup][r.Func] {
			continue
		}
		seen[r.VldrGroup][r.Func] = true
		if _, ok := lookup(uint16(r.VldrGroup), r.Func); !ok {
			missing = append(missing, fmt.Sprintf("group=%d func=%s", r.VldrGroup, r.Func))
		}
	}
	if len(missing) > 0 {
		log.Errorf("validate: %d rules missing handlers (will be skipped): %v", len(missing), missing)
	}
}
