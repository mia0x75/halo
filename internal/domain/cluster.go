package domain

// ClusterType 目标数据库的类型。存储值对齐 mm_clusters.type 列。
type ClusterType string

const (
	ClusterTypeMySQL    ClusterType = "mysql"
	ClusterTypePostgres ClusterType = "postgres"
)

// Cluster 目标数据库群集。仅承载执行器必需的最小字段集，
// 业务字段由 models.Cluster 持有，二者在阶段 2 收敛。
type Cluster struct {
	ClusterID uint
	UUID      string
	Host      string
	IP        string
	Port      uint16
	Alias     string
	User      string
	Type      ClusterType
	Database  string
	Status    uint8
}
