package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/models"
)

// CronFromDomain converts a domain entity to its legacy model. domain.Cron.Hash
// has no counterpart on the legacy model, so it is dropped.
func CronFromDomain(c *domain.Cron) *models.Cron {
	if c == nil {
		return nil
	}
	return &models.Cron{
		CronID:    c.CronID,
		UUID:      c.UUID,
		Status:    string(c.Status),
		Name:      c.Name,
		Cmd:       c.Cmd,
		Params:    c.Params,
		Interval:  derefStr(c.Interval),
		Duration:  derefStr(c.Duration),
		LastRun:   derefStr(c.LastRun),
		NextRun:   derefStr(c.NextRun),
		Recurrent: derefU8(c.Recurrent),
		Version:   c.Version,
		UpdateAt:  deref(c.UpdateAt),
		CreateAt:  c.CreateAt,
	}
}

// CronToDomain converts a legacy model to its domain entity. Hash is left nil
// because the legacy model never carried it.
func CronToDomain(c *models.Cron) *domain.Cron {
	if c == nil {
		return nil
	}
	return &domain.Cron{
		CronID:    c.CronID,
		UUID:      c.UUID,
		Status:    domain.CronStatus(c.Status),
		Name:      c.Name,
		Cmd:       c.Cmd,
		Params:    c.Params,
		Interval:  strPtr(c.Interval),
		Duration:  strPtr(c.Duration),
		LastRun:   strPtr(c.LastRun),
		NextRun:   strPtr(c.NextRun),
		Recurrent: u8Ptr(c.Recurrent),
		Version:   c.Version,
		UpdateAt:  uintPtr(c.UpdateAt),
		CreateAt:  c.CreateAt,
	}
}
