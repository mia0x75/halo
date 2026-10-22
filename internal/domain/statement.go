package domain

// StatementStatus is the validation/execution state of a statement. Stored in
// mm_statements.status.
type StatementStatus int

const (
	StatementStatusPending       StatementStatus = 1
	StatementStatusVldFailure    StatementStatus = 2
	StatementStatusVldWarning    StatementStatus = 3
	StatementStatusWaitingForMrv StatementStatus = 4
	// 执行结果沿用工单状态的编号，便于两个层面直接对照。
	StatementStatusDone        StatementStatus = 7
	StatementStatusExecFailure StatementStatus = 8
)

// StatementType is the SQL operation of a statement. Stored in mm_statements.type.
type StatementType int

const (
	StatementTypeAlterTable     StatementType = 1
	StatementTypeAnalyzeTable   StatementType = 2
	StatementTypeBegin          StatementType = 3
	StatementTypeCommit         StatementType = 4
	StatementTypeCreateDatabase StatementType = 5
	StatementTypeCreateIndex    StatementType = 6
	StatementTypeCreateTable    StatementType = 7
	StatementTypeCreateView     StatementType = 8
	StatementTypeCreateUser     StatementType = 9
	StatementTypeDelete         StatementType = 10
	StatementTypeDropDatabase   StatementType = 11
	StatementTypeDropIndex      StatementType = 12
	StatementTypeDropTable      StatementType = 13
	StatementTypeExplain        StatementType = 14
	StatementTypeInsert         StatementType = 15
	StatementTypeLoadData       StatementType = 16
	StatementTypeRollback       StatementType = 17
	StatementTypeSelect         StatementType = 18
	StatementTypeSet            StatementType = 19
	StatementTypeShow           StatementType = 20
	StatementTypeTruncateTable  StatementType = 21
	StatementTypeUpdate         StatementType = 22
	StatementTypeGrant          StatementType = 23
	StatementTypeRevoke         StatementType = 24
	StatementTypeDeallocate     StatementType = 25
	StatementTypeExecute        StatementType = 26
	StatementTypePrepare        StatementType = 27
	StatementTypeUse            StatementType = 28
)

// Statement is a row of mm_statements. The primary key is (ticket_id, sequence);
// UUID is the stable lookup key used by the repository.
type Statement struct {
	TicketID     uint            `db:"ticket_id"`
	Sequence     uint16          `db:"sequence"`
	UUID         string          `db:"uuid"`
	Content      string          `db:"content"`
	Type         StatementType   `db:"type"`
	Status       StatementStatus `db:"status"`
	Report       string          `db:"report"`
	Plan         *string         `db:"plan"`
	Results      *string         `db:"results"`
	RowsAffected *uint           `db:"rows_affected"`
	Version      int             `db:"version"`
	UpdateAt     *uint           `db:"update_at"`
	CreateAt     uint            `db:"create_at"`
}
