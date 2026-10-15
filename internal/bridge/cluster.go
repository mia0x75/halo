package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// ClusterFromDomain converts a domain entity to its legacy model. domain.Cluster
// carries a non-persisted Database, which the legacy model has no slot for.
func ClusterFromDomain(c *domain.Cluster) *models.Cluster {
	if c == nil {
		return nil
	}
	return &models.Cluster{
		ClusterID:   c.ClusterID,
		UUID:        c.UUID,
		Host:        c.Host,
		IP:          c.IP,
		Port:        uint16(c.Port),
		Alias:       c.Alias,
		User:        c.User,
		Type:        string(c.Type),
		Password:    c.Password,
		FingerPrint: c.Fingerprint,
		Status:      uint8(c.Status),
		Version:     c.Version,
		UpdateAt:    deref(c.UpdateAt),
		CreateAt:    c.CreateAt,
	}
}

// ClusterToDomain converts a legacy model to its domain entity.
func ClusterToDomain(c *models.Cluster) *domain.Cluster {
	if c == nil {
		return nil
	}
	return &domain.Cluster{
		ClusterID:   c.ClusterID,
		UUID:        c.UUID,
		Host:        c.Host,
		IP:          c.IP,
		Port:        uint(c.Port),
		Alias:       c.Alias,
		User:        c.User,
		Type:        domain.ClusterType(c.Type),
		Password:    c.Password,
		Fingerprint: c.FingerPrint,
		Status:      domain.ClusterStatus(c.Status),
		Version:     c.Version,
		UpdateAt:    uintPtr(c.UpdateAt),
		CreateAt:    c.CreateAt,
	}
}
