package mysqlrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.ClusterRepository = (*ClusterRepo)(nil)

type ClusterRepo struct {
	q *Queries
}

func NewClusterRepo(q *Queries) *ClusterRepo {
	return &ClusterRepo{q: q}
}

func clusterFromRow(r MmCluster) *domain.Cluster {
	return &domain.Cluster{
		ClusterID:   uint(r.ClusterID),
		UUID:        r.Uuid,
		Host:        r.Host,
		Alias:       r.Alias,
		IP:          r.Ip,
		Port:        uint(r.Port),
		User:        r.User,
		Type:        domain.ClusterType(r.Type),
		Password:    r.Password,
		Fingerprint: r.Fingerprint,
		Status:      domain.ClusterStatus(r.Status),
		Version:     int(r.Version),
		UpdateAt:    timeFromNull(r.UpdateAt),
		CreateAt:    timeToUnix(r.CreateAt),
	}
}

func (r *ClusterRepo) Create(ctx context.Context, c *domain.Cluster) (int64, error) {
	id, err := r.q.CreateCluster(ctx, CreateClusterParams{
		Uuid:        c.UUID,
		Host:        c.Host,
		Alias:       c.Alias,
		Ip:          c.IP,
		Port:        uint32(c.Port),
		User:        c.User,
		Type:        string(c.Type),
		Password:    c.Password,
		Fingerprint: c.Fingerprint,
		Status:      uint8(c.Status),
		Version:     uint32(c.Version),
	})
	if err != nil {
		return 0, wrap("cluster.Create", err)
	}
	return id, nil
}

func (r *ClusterRepo) GetByID(ctx context.Context, id int64) (*domain.Cluster, error) {
	row, err := r.q.GetClusterByID(ctx, uint32(id))
	if err != nil {
		return nil, wrap("cluster.GetByID", err)
	}
	return clusterFromRow(row), nil
}

func (r *ClusterRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Cluster, error) {
	row, err := r.q.GetClusterByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("cluster.GetByUUID", err)
	}
	return clusterFromRow(row), nil
}

func (r *ClusterRepo) Update(ctx context.Context, c *domain.Cluster) error {
	err := r.q.UpdateCluster(ctx, UpdateClusterParams{
		Uuid:        c.UUID,
		Host:        c.Host,
		Alias:       c.Alias,
		Ip:          c.IP,
		Port:        uint32(c.Port),
		User:        c.User,
		Type:        string(c.Type),
		Password:    c.Password,
		Fingerprint: c.Fingerprint,
		Status:      uint8(c.Status),
		Version:     uint32(c.Version),
		ClusterID:   uint32(c.ClusterID),
	})
	return wrap("cluster.Update", err)
}

func (r *ClusterRepo) Delete(ctx context.Context, id int64) error {
	return wrap("cluster.Delete", r.q.DeleteClusterByID(ctx, uint32(id)))
}

func (r *ClusterRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("cluster.DeleteByUUID", r.q.DeleteClusterByUUID(ctx, uuid))
}

func (r *ClusterRepo) List(ctx context.Context, cursor int64, limit int) ([]*domain.Cluster, error) {
	rows, err := r.q.ListClusters(ctx, ListClustersParams{
		ClusterID: uint32(cursor),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, wrap("cluster.List", err)
	}
	out := make([]*domain.Cluster, 0, len(rows))
	for _, row := range rows {
		out = append(out, clusterFromRow(row))
	}
	return out, nil
}
