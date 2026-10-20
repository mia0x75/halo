package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mia0x75/halo/internal/agent"
	"github.com/mia0x75/halo/internal/domain"
)

type ClusterService struct {
	tx       TxManager
	clusters domain.ClusterRepository
	executor agent.Executor
}

func NewClusterService(tx TxManager, clusters domain.ClusterRepository, executor agent.Executor) *ClusterService {
	return &ClusterService{tx: tx, clusters: clusters, executor: executor}
}

// ConnectionResult 试连目标库的结果。连不上不是错误，Success 为 false 并带上原因。
type ConnectionResult struct {
	Success   bool
	Message   string
	Version   string
	Latency   time.Duration
	Databases []string
}

// TestConnection 用给定的凭据试连目标库。不落库，密码不出内存。
func (s *ClusterService) TestConnection(ctx context.Context, c *domain.Cluster, password string) (*ConnectionResult, error) {
	info, err := s.executor.TestConnection(ctx, c, password)
	if err != nil {
		return &ConnectionResult{Success: false, Message: err.Error()}, nil
	}
	return &ConnectionResult{
		Success:   true,
		Message:   "连接成功",
		Version:   info.Version,
		Latency:   info.Latency,
		Databases: info.Databases,
	}, nil
}

func (s *ClusterService) Get(ctx context.Context, id uint) (*domain.Cluster, error) {
	return s.clusters.GetByID(ctx, int64(id))
}

func (s *ClusterService) GetByUUID(ctx context.Context, uuid string) (*domain.Cluster, error) {
	return s.clusters.GetByUUID(ctx, uuid)
}

func (s *ClusterService) List(ctx context.Context, cursor int64, limit int) ([]*domain.Cluster, error) {
	return s.clusters.List(ctx, cursor, limit)
}

func (s *ClusterService) Create(ctx context.Context, c *domain.Cluster, operatorID uint) (uint, error) {
	var id uint
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		newID, err := repos.Cluster.Create(ctx, c)
		if err != nil {
			return err
		}
		id = uint(newID)
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "cluster.create", Target: "cluster:" + c.UUID})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ClusterService) Update(ctx context.Context, c *domain.Cluster, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Cluster.Update(ctx, c); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "cluster.update", Target: "cluster:" + c.UUID})
	})
}

func (s *ClusterService) Delete(ctx context.Context, id uint, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Cluster.Delete(ctx, int64(id)); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "cluster.delete", Target: fmt.Sprintf("cluster:%d", id)})
	})
}
