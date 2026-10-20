package domain

// ClusterType is the target database engine. Stored in mm_clusters.type.
type ClusterType string

const (
	ClusterTypeMySQL    ClusterType = "mysql"
	ClusterTypePostgres ClusterType = "postgres"
)

// ClusterStatus is the lifecycle state of a cluster. Stored in mm_clusters.status.
type ClusterStatus int

const (
	ClusterStatusNormal   ClusterStatus = 1
	ClusterStatusDisabled ClusterStatus = 2
)

// Cluster is a row of mm_clusters. Database is not a column: it is the schema
// resolved at connect time and carried here so the executor takes one value.
type Cluster struct {
	ClusterID   uint          `db:"cluster_id"`
	UUID        string        `db:"uuid"`
	Host        string        `db:"host"`
	Alias       string        `db:"alias"`
	IP          string        `db:"ip"`
	Port        uint          `db:"port"`
	User        string        `db:"user"`
	Type        ClusterType   `db:"type"`
	GroupID     uint          `db:"group_id"`
	Password    []byte        `db:"password"`
	Fingerprint []byte        `db:"fingerprint"`
	Status      ClusterStatus `db:"status"`
	Version     int           `db:"version"`
	UpdateAt    *uint         `db:"update_at"`
	CreateAt    uint          `db:"create_at"`
	Database    string        `db:"-"`
}
