package domain

// EdgeType is the relation kind carried by an edge. Stored in mm_edges.type.
type EdgeType int

const (
	EdgeTypeUserToReviewer EdgeType = 1
	EdgeTypeUserToRole     EdgeType = 2
	EdgeTypeUserToCluster  EdgeType = 3
)

// Edge is a row of mm_edges; it links two entities by kind.
type Edge struct {
	EdgeID       uint     `db:"edge_id"`
	UUID         string   `db:"uuid"`
	Type         EdgeType `db:"type"`
	AncestorID   uint     `db:"ancestor_id"`
	DescendantID uint     `db:"descendant_id"`
	Version      int      `db:"version"`
	UpdateAt     *uint    `db:"update_at"`
	CreateAt     uint     `db:"create_at"`
}
