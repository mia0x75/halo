package domain

// QueryType is the purpose of a logged query. Stored in mm_queries.type.
type QueryType int

const (
	QueryTypeQuery   QueryType = 1
	QueryTypeAnalyze QueryType = 2
	QueryTypeRewrite QueryType = 3
)

// Query is a row of mm_queries.
type Query struct {
	QueryID   uint      `db:"query_id"`
	UUID      string    `db:"uuid"`
	Type      QueryType `db:"type"`
	ClusterID uint      `db:"cluster_id"`
	Database  string    `db:"database"`
	Content   string    `db:"content"`
	Plan      string    `db:"plan"`
	UserID    uint      `db:"user_id"`
	Version   int       `db:"version"`
	UpdateAt  *uint     `db:"update_at"`
	CreateAt  uint      `db:"create_at"`
}
