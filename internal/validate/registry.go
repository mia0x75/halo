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
		"AvailableCharsets":             wrap((*DatabaseCreateVldr).AvailableCharsets),
		"AvailableCollates":             wrap((*DatabaseCreateVldr).AvailableCollates),
		"CharsetCollateMustMatch":       wrap((*DatabaseCreateVldr).CharsetCollateMustMatch),
		"DatabaseNameLowerCaseRequired": wrap((*DatabaseCreateVldr).DatabaseNameLowerCaseRequired),
		"DatabaseNameMaxLength":         wrap((*DatabaseCreateVldr).DatabaseNameMaxLength),
		"DatabaseNameQualified":         wrap((*DatabaseCreateVldr).DatabaseNameQualified),
		"TargetDatabaseDoesNotExist":    wrap((*DatabaseCreateVldr).TargetDatabaseDoesNotExist),
	},
	101: { // DatabaseAlterVldr
		"AvailableCharsets":          wrap((*DatabaseAlterVldr).AvailableCharsets),
		"AvailableCollates":          wrap((*DatabaseAlterVldr).AvailableCollates),
		"CharsetCollateMustMatch":    wrap((*DatabaseAlterVldr).CharsetCollateMustMatch),
		"TargetDatabaseDoesNotExist": wrap((*DatabaseAlterVldr).TargetDatabaseDoesNotExist),
	},
	102: { // DatabaseDropVldr
		"TargetDatabaseDoesNotExist": wrap((*DatabaseDropVldr).TargetDatabaseDoesNotExist),
	},
	110: { // TableCreateVldr
		"AutoIncColumnDuplicate":           wrap((*TableCreateVldr).AutoIncColumnDuplicate),
		"AvailableCharsets":                wrap((*TableCreateVldr).AvailableCharsets),
		"AvailableCollates":                wrap((*TableCreateVldr).AvailableCollates),
		"AvailableEngines":                 wrap((*TableCreateVldr).AvailableEngines),
		"ColumnAutoIncAvailableTypes":      wrap((*TableCreateVldr).ColumnAutoIncAvailableTypes),
		"ColumnAutoIncMustPrimaryKey":      wrap((*TableCreateVldr).ColumnAutoIncMustPrimaryKey),
		"ColumnAutoIncUnsignedRequired":    wrap((*TableCreateVldr).ColumnAutoIncUnsignedRequired),
		"ColumnAvailableCharsets":          wrap((*TableCreateVldr).ColumnAvailableCharsets),
		"ColumnAvailableCollates":          wrap((*TableCreateVldr).ColumnAvailableCollates),
		"ColumnCharsetCollateMustMatch":    wrap((*TableCreateVldr).ColumnCharsetCollateMustMatch),
		"ColumnCommentRequired":            wrap((*TableCreateVldr).ColumnCommentRequired),
		"ColumnNameDuplicate":              wrap((*TableCreateVldr).ColumnNameDuplicate),
		"ColumnNameLowerCaseRequired":      wrap((*TableCreateVldr).ColumnNameLowerCaseRequired),
		"ColumnNameMaxLength":              wrap((*TableCreateVldr).ColumnNameMaxLength),
		"ColumnNameQualified":              wrap((*TableCreateVldr).ColumnNameQualified),
		"ColumnNotNullWithDefaultRequired": wrap((*TableCreateVldr).ColumnNotNullWithDefaultRequired),
		"ColumnTypesDoesNotExpect":         wrap((*TableCreateVldr).ColumnTypesDoesNotExpect),
		"CreateTableFromSelectNotAllowed":  wrap((*TableCreateVldr).CreateTableFromSelectNotAllowed),
		"CreateTableUseLikeNotAllowed":     wrap((*TableCreateVldr).CreateTableUseLikeNotAllowed),
		"ForeignKeyNameExplicit":           wrap((*TableCreateVldr).ForeignKeyNameExplicit),
		"ForeignKeyNameLowerCaseRequired":  wrap((*TableCreateVldr).ForeignKeyNameLowerCaseRequired),
		"ForeignKeyNameMaxLength":          wrap((*TableCreateVldr).ForeignKeyNameMaxLength),
		"ForeignKeyNamePrefixRequired":     wrap((*TableCreateVldr).ForeignKeyNamePrefixRequired),
		"ForeignKeyNameQualified":          wrap((*TableCreateVldr).ForeignKeyNameQualified),
		"ForeignKeyNotAllowed":             wrap((*TableCreateVldr).ForeignKeyNotAllowed),
		"IndexNameExplicit":                wrap((*TableCreateVldr).IndexNameExplicit),
		"IndexNameLowerCaseRequired":       wrap((*TableCreateVldr).IndexNameLowerCaseRequired),
		"IndexNameMaxLength":               wrap((*TableCreateVldr).IndexNameMaxLength),
		"IndexNamePrefixRequired":          wrap((*TableCreateVldr).IndexNamePrefixRequired),
		"IndexNameQualified":               wrap((*TableCreateVldr).IndexNameQualified),
		"MaxAllowedColumnCount":            wrap((*TableCreateVldr).MaxAllowedColumnCount),
		"MaxAllowedIndexColumnCount":       wrap((*TableCreateVldr).MaxAllowedIndexColumnCount),
		"MaxAllowedIndexCount":             wrap((*TableCreateVldr).MaxAllowedIndexCount),
		"MaxAllowedTimestampCount":         wrap((*TableCreateVldr).MaxAllowedTimestampCount),
		"PrimaryKeyDuplicate":              wrap((*TableCreateVldr).PrimaryKeyDuplicate),
		"PrimaryKeyNameExplicit":           wrap((*TableCreateVldr).PrimaryKeyNameExplicit),
		"PrimaryKeyNameQualified":          wrap((*TableCreateVldr).PrimaryKeyNameQualified),
		"PrimaryKeyRequired":               wrap((*TableCreateVldr).PrimaryKeyRequired),
		"PrimryKeyLowerCaseRequired":       wrap((*TableCreateVldr).PrimryKeyLowerCaseRequired),
		"PrimryKeyMaxLength":               wrap((*TableCreateVldr).PrimryKeyMaxLength),
		"PrimryKeyPrefixRequired":          wrap((*TableCreateVldr).PrimryKeyPrefixRequired),
		"TableCharsetCollateMustMatch":     wrap((*TableCreateVldr).TableCharsetCollateMustMatch),
		"TableCommentRequired":             wrap((*TableCreateVldr).TableCommentRequired),
		"TableNameLowerCaseRequired":       wrap((*TableCreateVldr).TableNameLowerCaseRequired),
		"TableNameMaxLength":               wrap((*TableCreateVldr).TableNameMaxLength),
		"TableNameQualified":               wrap((*TableCreateVldr).TableNameQualified),
		"TargetDatabaseDoesNotExist":       wrap((*TableCreateVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":          wrap((*TableCreateVldr).TargetTableDoesNotExist),
		"UniqueNameExplicit":               wrap((*TableCreateVldr).UniqueNameExplicit),
		"UniqueNameLowerCaseRequired":      wrap((*TableCreateVldr).UniqueNameLowerCaseRequired),
		"UniqueNameMaxLength":              wrap((*TableCreateVldr).UniqueNameMaxLength),
		"UniqueNamePrefixRequired":         wrap((*TableCreateVldr).UniqueNamePrefixRequired),
		"UniqueNameQualified":              wrap((*TableCreateVldr).UniqueNameQualified),
	},
	120: { // TableAlterVldr
		"AvailableCharsets":                  wrap((*TableAlterVldr).AvailableCharsets),
		"AvailableCollates":                  wrap((*TableAlterVldr).AvailableCollates),
		"AvailableEngines":                   wrap((*TableAlterVldr).AvailableEngines),
		"ColumnAvailableCharsets":            wrap((*TableAlterVldr).ColumnAvailableCharsets),
		"ColumnAvailableCollates":            wrap((*TableAlterVldr).ColumnAvailableCollates),
		"ColumnCharsetCollateMustMatch":      wrap((*TableAlterVldr).ColumnCharsetCollateMustMatch),
		"ColumnCommentRequired":              wrap((*TableAlterVldr).ColumnCommentRequired),
		"CommentRequired":                    wrap((*TableAlterVldr).CommentRequired),
		"ColumnNameDoesNotExist":             wrap((*TableAlterVldr).ColumnNameDoesNotExist),
		"ColumnNameDuplicate":                wrap((*TableAlterVldr).ColumnNameDuplicate),
		"ColumnNameLowerCaseRequired":        wrap((*TableAlterVldr).ColumnNameLowerCaseRequired),
		"ColumnNameMaxLength":                wrap((*TableAlterVldr).ColumnNameMaxLength),
		"ColumnNameQualified":                wrap((*TableAlterVldr).ColumnNameQualified),
		"ColumnNotNullWithDefaultRequired":   wrap((*TableAlterVldr).ColumnNotNullWithDefaultRequired),
		"ColumnUnwantedTypes":                wrap((*TableAlterVldr).ColumnUnwantedTypes),
		"ForeignKeyDoesNotExist":             wrap((*TableAlterVldr).ForeignKeyDoesNotExist),
		"ForeignKeyNameExplicit":             wrap((*TableAlterVldr).ForeignKeyNameExplicit),
		"ForeignKeyNameLowerCaseRequired":    wrap((*TableAlterVldr).ForeignKeyNameLowerCaseRequired),
		"ForeignKeyNameMaxLength":            wrap((*TableAlterVldr).ForeignKeyNameMaxLength),
		"ForeignKeyNamePrefixRequired":       wrap((*TableAlterVldr).ForeignKeyNamePrefixRequired),
		"ForeignKeyNameQualified":            wrap((*TableAlterVldr).ForeignKeyNameQualified),
		"ForeignKeyNotAllowed":               wrap((*TableAlterVldr).ForeignKeyNotAllowed),
		"FullTextIndexNameDoesNotExist":      wrap((*TableAlterVldr).FullTextIndexNameDoesNotExist),
		"FullTextIndexNameDuplicate":         wrap((*TableAlterVldr).FullTextIndexNameDuplicate),
		"FullTextIndexNameExplicit":          wrap((*TableAlterVldr).FullTextIndexExplicit),
		"FullTextIndexNameLowerCaseRequired": wrap((*TableAlterVldr).FullTextIndexNameLowerCaseRequired),
		"FullTextIndexNameMaxLength":         wrap((*TableAlterVldr).FullTextIndexNameMaxLength),
		"FullTextIndexNamePrefixRequired":    wrap((*TableAlterVldr).FullTextIndexNamePrefixRequired),
		"FullTextIndexNameQualified":         wrap((*TableAlterVldr).FullTextIndexNameQualified),
		"FullTextIndexNotAllowed":            wrap((*TableAlterVldr).FullTextIndexNotAllowed),
		"IndexColumnDoesNotExist":            wrap((*TableAlterVldr).IndexColumnDoesNotExist),
		"IndexDoesNotExist":                  wrap((*TableAlterVldr).IndexDoesNotExist),
		"IndexNameDuplicate":                 wrap((*TableAlterVldr).IndexNameDuplicate),
		"IndexNameExplicit":                  wrap((*TableAlterVldr).IndexNameExplicit),
		"IndexNameLowerCaseRequired":         wrap((*TableAlterVldr).IndexNameLowerCaseRequired),
		"IndexNameMaxLength":                 wrap((*TableAlterVldr).IndexNameMaxLength),
		"IndexNamePrefixRequired":            wrap((*TableAlterVldr).IndexNamePrefixRequired),
		"IndexNameQualified":                 wrap((*TableAlterVldr).IndexNameQualified),
		"IndexOnBlobColumnNotAllowed":        wrap((*TableAlterVldr).IndexOnBlobColumnNotAllowed),
		"IndexOverlayNotAllowed":             wrap((*TableAlterVldr).IndexOverlayNotAllowed),
		"MaxAllowedColumnCount":              wrap((*TableAlterVldr).MaxAllowedColumnCount),
		"MaxAllowedIndexColumnCount":         wrap((*TableAlterVldr).MaxAllowedIndexColumnCount),
		"IndexCountLimit":                    wrap((*TableAlterVldr).IndexCountLimit),
		"MaxAllowedTimestampCount":           wrap((*TableAlterVldr).MaxAllowedTimestampCount),
		"NewTableNameLowerCaseRequired":      wrap((*TableAlterVldr).NewTableNameLowerCaseRequired),
		"NewTableNameMaxLength":              wrap((*TableAlterVldr).NewTableNameMaxLength),
		"NewTableNameQualified":              wrap((*TableAlterVldr).NewTableNameQualified),
		"PositionColumnDoesNotExist":         wrap((*TableAlterVldr).PositionColumnDoesNotExist),
		"PrimaryKeyDoesNotExist":             wrap((*TableAlterVldr).PrimaryKeyDoesNotExist),
		"TableCharsetCollateMustMatch":       wrap((*TableAlterVldr).TableCharsetCollateMustMatch),
		"TargetDatabaseDoesNotExist":         wrap((*TableAlterVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":            wrap((*TableAlterVldr).TargetTableDoesNotExist),
		"TargetTableDuplicate":               wrap((*TableAlterVldr).TargetTableDuplicate),
		"UniqueNameExplicit":                 wrap((*TableAlterVldr).UniqueNameExplicit),
		"UniqueNameLowerCaseRequired":        wrap((*TableAlterVldr).UniqueNameLowerCaseRequired),
		"UniqueNameMaxLength":                wrap((*TableAlterVldr).UniqueNameMaxLength),
		"UniqueNamePrefixRequired":           wrap((*TableAlterVldr).UniqueNamePrefixRequired),
		"UniqueNameQualified":                wrap((*TableAlterVldr).UniqueNameQualified),
	},
	130: { // TableRenameVldr
		"SourceDatabaseDoesNotExist":       wrap((*TableRenameVldr).SourceDatabaseDoesNotExist),
		"SourceTableDoesNotExist":          wrap((*TableRenameVldr).SourceTableDoesNotExist),
		"TablesIdentical":                  wrap((*TableRenameVldr).TablesIdentical),
		"TargetDatabaseDoesNotExist":       wrap((*TableRenameVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":          wrap((*TableRenameVldr).TargetTableDoesNotExist),
		"TargetTableNameLowerCaseRequired": wrap((*TableRenameVldr).TargetTableNameLowerCaseRequired),
		"TargetTableNameMaxLength":         wrap((*TableRenameVldr).TargetTableNameMaxLength),
		"TargetTableNameQualified":         wrap((*TableRenameVldr).TargetTableNameQualified),
	},
	140: { // TableDropVldr
		"TableDropTargetDatabaseDoesNotExist": wrap((*TableDropVldr).TargetDatabaseDoesNotExist),
		"TableDropTargetTableDoesNotExist":    wrap((*TableDropVldr).TargetTableDoesNotExist),
	},
	150: { // IndexCreateVldr
		"IndexColumnDuplicate":        wrap((*IndexCreateVldr).IndexColumnDuplicate),
		"IndexNameDuplicate":          wrap((*IndexCreateVldr).IndexNameDuplicate),
		"IndexNameLowerCaseRequired":  wrap((*IndexCreateVldr).IndexNameLowerCaseRequired),
		"IndexNameMaxLength":          wrap((*IndexCreateVldr).IndexNameMaxLength),
		"IndexNamePrefixRequired":     wrap((*IndexCreateVldr).IndexNamePrefixRequired),
		"IndexNameQualified":          wrap((*IndexCreateVldr).IndexNameQualified),
		"IndexOnBlobColumnNotAllowed": wrap((*IndexCreateVldr).IndexOnBlobColumnNotAllowed),
		"IndexOverlayNotAllowed":      wrap((*IndexCreateVldr).IndexOverlayNotAllowed),
		"MaxAllowedIndexColumnCount":  wrap((*IndexCreateVldr).MaxAllowedIndexColumnCount),
		"MaxAllowedIndexCount":        wrap((*IndexCreateVldr).MaxAllowedIndexCount),
		"TargetColumnDoesNotExist":    wrap((*IndexCreateVldr).TargetColumnDoesNotExist),
		"TargetDatabaseDoesNotExist":  wrap((*IndexCreateVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":     wrap((*IndexCreateVldr).TargetTableDoesNotExist),
	},
	151: { // IndexDropVldr
		"TargetDatabaseDoesNotExist": wrap((*IndexDropVldr).TargetDatabaseDoesNotExist),
		"TargetIndexDoesNotExist":    wrap((*IndexDropVldr).TargetIndexDoesNotExist),
		"TargetTableDoesNotExist":    wrap((*IndexDropVldr).TargetTableDoesNotExist),
	},
	160: { // InsertVldr
		"ColumnsValuesMustMatch":        wrap((*InsertVldr).ColumnsValuesMustMatch),
		"ExplicitColumnRequired":        wrap((*InsertVldr).ExplicitColumnRequired),
		"MergeRequired":                 wrap((*InsertVldr).MergeRequired),
		"RowsLimit":                     wrap((*InsertVldr).RowsLimit),
		"TargetColumnDoesNotExist":      wrap((*InsertVldr).TargetColumnDoesNotExist),
		"TargetDatabaseDoesNotExist":    wrap((*InsertVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":       wrap((*InsertVldr).TargetTableDoesNotExist),
		"UsingSelectNotAllowed":         wrap((*InsertVldr).UsingSelectNotAllowed),
		"ValueForNotNullColumnRequired": wrap((*InsertVldr).ValueForNotNullColumnRequired),
	},
	161: { // ReplaceVldr
		"ColumnsValuesMustMatch":        wrap((*ReplaceVldr).ColumnsValuesMustMatch),
		"ExplicitColumnRequired":        wrap((*ReplaceVldr).ExplicitColumnRequired),
		"MergeRequired":                 wrap((*ReplaceVldr).MergeRequired),
		"ReplaceRowsLimit":              wrap((*ReplaceVldr).RowsLimit),
		"TargetColumnDoesNotExist":      wrap((*ReplaceVldr).TargetColumnDoesNotExist),
		"TargetDatabaseDoesNotExist":    wrap((*ReplaceVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":       wrap((*ReplaceVldr).TargetTableDoesNotExist),
		"UsingSelectNotAllowed":         wrap((*ReplaceVldr).UsingSelectNotAllowed),
		"ValueForNotNullColumnRequired": wrap((*ReplaceVldr).ValueForNotNullColumnRequired),
	},
	162: { // UpdateVldr
		"FilterColumnExists":         wrap((*UpdateVldr).FilterColumnExists),
		"MaxAllowedUpdateRows":       wrap((*UpdateVldr).MaxAllowedUpdateRows),
		"TargetColumnDoesNotExist":   wrap((*UpdateVldr).TargetColumnDoesNotExist),
		"TargetDatabaseDoesNotExist": wrap((*UpdateVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":    wrap((*UpdateVldr).TargetTableDoesNotExist),
		"WithoutWhereNotAllowed":     wrap((*UpdateVldr).WithoutWhereNotAllowed),
	},
	163: { // DeleteVldr
		"MaxAllowedDeleteRows":            wrap((*DeleteVldr).MaxAllowedDeleteRows),
		"TargetColumnDoesNotDoesNotExist": wrap((*DeleteVldr).TargetColumnDoesNotDoesNotExist),
		"TargetDatabaseDoesNotExist":      wrap((*DeleteVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotDoesNotExist":  wrap((*DeleteVldr).TargetTableDoesNotDoesNotExist),
		"WithoutWhereNotAllowed":          wrap((*DeleteVldr).WithoutWhereNotAllowed),
	},
	170: { // SelectVldr
		"ReturnBlobOrTextNotAllowed":      wrap((*SelectVldr).ReturnBlobOrTextNotAllowed),
		"TargetColumnDoesNotDoesNotExist": wrap((*SelectVldr).TargetColumnDoesNotDoesNotExist),
		"TargetDatabaseDoesNotExist":      wrap((*SelectVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":         wrap((*SelectVldr).TargetTableDoesNotExist),
		"UseExplicitLockNotAllowed":       wrap((*SelectVldr).UseExplicitLockNotAllowed),
		"UseWildcardNotAllowed":           wrap((*SelectVldr).UseWildcardNotAllowed),
		"WithoutLimitNotAllowed":          wrap((*SelectVldr).WithoutLimitNotAllowed),
		"WithoutWhereNotAllowed":          wrap((*SelectVldr).WithoutWhereNotAllowed),
	},
	180: { // ViewCreateVldr
		"TargetDatabaseDoesNotExist": wrap((*ViewCreateVldr).TargetDatabaseDoesNotExist),
		"TargetViewDoesNotExist":     wrap((*ViewCreateVldr).TargetViewDoesNotExist),
		"ViewNameLowerCaseRequired":  wrap((*ViewCreateVldr).ViewNameLowerCaseRequired),
		"ViewNameMaxLength":          wrap((*ViewCreateVldr).ViewNameMaxLength),
		"ViewNamePrefixRequired":     wrap((*ViewCreateVldr).ViewNamePrefixRequired),
		"ViewNameQualified":          wrap((*ViewCreateVldr).ViewNameQualified),
	},
	181: { // ViewAlterVldr
		"TargetDatabaseDoesNotExist": wrap((*ViewAlterVldr).TargetDatabaseDoesNotExist),
		"TargetViewDoesNotExist":     wrap((*ViewAlterVldr).TargetViewDoesNotExist),
	},
	182: { // ViewDropVldr
		"TargetDatabaseDoesNotExist": wrap((*ViewDropVldr).TargetDatabaseDoesNotExist),
		"TargetViewDoesNotExist":     wrap((*ViewDropVldr).TargetViewDoesNotExist),
	},
	190: { // FuncCreateVldr
		"FuncNameLowerCaseRequired":  wrap((*FuncCreateVldr).FuncNameLowerCaseRequired),
		"FuncNameMaxLength":          wrap((*FuncCreateVldr).FuncNameMaxLength),
		"FuncNamePrefixRequired":     wrap((*FuncCreateVldr).FuncNamePrefixRequired),
		"FuncNameQuilified":          wrap((*FuncCreateVldr).FuncNameQuilified),
		"TargetDatabaseDoesNotExist": wrap((*FuncCreateVldr).TargetDatabaseDoesNotExist),
		"TargetFuncDoesNotExist":     wrap((*FuncCreateVldr).TargetFuncDoesNotExist),
	},
	191: { // FuncAlterVldr
		"TargetDatabaseDoesNotExist": wrap((*FuncAlterVldr).TargetDatabaseDoesNotExist),
		"TargetFuncDoesNotExist":     wrap((*FuncAlterVldr).TargetFuncDoesNotExist),
	},
	192: { // FuncDropVldr
		"TargetDatabaseDoesNotExist": wrap((*FuncDropVldr).TargetDatabaseDoesNotExist),
		"TargetFuncDoesNotExist":     wrap((*FuncDropVldr).TargetFuncDoesNotExist),
	},
	200: { // TriggerCreateVldr
		"TargetDatabaseDoesNotExist":   wrap((*TriggerCreateVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":      wrap((*TriggerCreateVldr).TargetTableDoesNotExist),
		"TargetTriggerDoesNotExist":    wrap((*TriggerCreateVldr).TargetTriggerDoesNotExist),
		"TriggerNameLowerCaseRequired": wrap((*TriggerCreateVldr).TriggerNameLowerCaseRequired),
		"TriggerNameMaxLength":         wrap((*TriggerCreateVldr).TriggerNameMaxLength),
		"TriggerNameQualified":         wrap((*TriggerCreateVldr).TriggerNameQualified),
		"TriggerPrefixRequired":        wrap((*TriggerCreateVldr).TriggerPrefixRequired),
	},
	201: { // TriggerAlterVldr
		"TargetDatabaseDoesNotExist": wrap((*TriggerAlterVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":    wrap((*TriggerAlterVldr).TargetTableDoesNotExist),
		"TargetTriggerDoesNotExist":  wrap((*TriggerAlterVldr).TargetTriggerDoesNotExist),
	},
	202: { // TriggerDropVldr
		"TargetDatabaseDoesNotExist": wrap((*TriggerDropVldr).TargetDatabaseDoesNotExist),
		"TargetTableDoesNotExist":    wrap((*TriggerDropVldr).TargetTableDoesNotExist),
		"TargetTriggerDoesNotExist":  wrap((*TriggerDropVldr).TargetTriggerDoesNotExist),
	},
	210: { // EventCreateVldr
		"EventNameLowerCaseRequired": wrap((*EventCreateVldr).EventNameLowerCaseRequired),
		"EventNameMaxLength":         wrap((*EventCreateVldr).EventNameMaxLength),
		"EventNamePrefixRequired":    wrap((*EventCreateVldr).EventNamePrefixRequired),
		"EventNameQualified":         wrap((*EventCreateVldr).EventNameQualified),
		"TargetDatabaseDoesNotExist": wrap((*EventCreateVldr).TargetDatabaseDoesNotExist),
		"TargetEventDoesNotExist":    wrap((*EventCreateVldr).TargetEventDoesNotExist),
	},
	211: { // EventAlterVldr
		"TargetDatabaseDoesNotExist": wrap((*EventAlterVldr).TargetDatabaseDoesNotExist),
		"TargetEventDoesNotExist":    wrap((*EventAlterVldr).TargetEventDoesNotExist),
	},
	212: { // EventDropVldr
		"TargetDatabaseDoesNotExist": wrap((*EventDropVldr).TargetDatabaseDoesNotExist),
		"TargetEventDoesNotExist":    wrap((*EventDropVldr).TargetEventDoesNotExist),
	},
	220: { // ProcCreateVldr
		"ProcNameLowerCaseRequired":  wrap((*ProcCreateVldr).ProcNameLowerCaseRequired),
		"ProcNameMaxLength":          wrap((*ProcCreateVldr).ProcNameMaxLength),
		"ProcNamePrefixRequired":     wrap((*ProcCreateVldr).ProcNamePrefixRequired),
		"ProcNameQualified":          wrap((*ProcCreateVldr).ProcNameQualified),
		"TargetDatabaseDoesNotExist": wrap((*ProcCreateVldr).TargetDatabaseDoesNotExist),
		"TargetProcDoesNotExist":     wrap((*ProcCreateVldr).TargetProcDoesNotExist),
	},
	221: { // ProcAlterVldr
		"TargetDatabaseDoesNotExist": wrap((*ProcAlterVldr).TargetDatabaseDoesNotExist),
		"TargetProcDoesNotExist":     wrap((*ProcAlterVldr).TargetProcDoesNotExist),
	},
	222: { // ProcDropVldr
		"TargetDatabaseDoesNotExist": wrap((*ProcDropVldr).TargetDatabaseDoesNotExist),
		"TargetProcDoesNotExist":     wrap((*ProcDropVldr).TargetProcDoesNotExist),
	},
	230: { // MiscVldr
		"FlushTableProhibited":    wrap((*MiscVldr).FlushTableProhibited),
		"KillProhibited":          wrap((*MiscVldr).KillProhibited),
		"LockTableProhibited":     wrap((*MiscVldr).LockTableProhibited),
		"MergeRequired":           wrap((*MiscVldr).MergeRequired),
		"PurgeLogsProhibited":     wrap((*MiscVldr).PurgeLogsProhibited),
		"SplitRequired":           wrap((*MiscVldr).SplitRequired),
		"SplitTicket":             wrap((*MiscVldr).SplitTicket),
		"TruncateTableProhibited": wrap((*MiscVldr).TruncateTableProhibited),
		"UnlockTableProhibited":   wrap((*MiscVldr).UnlockTableProhibited),
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
