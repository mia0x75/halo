package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/domain"
)

// LogService records operations for audit. It deliberately does not open its own
// transaction: callers run it inside TxManager.WithTx so the audit row commits
// atomically with the change it describes.
type LogService struct {
	log domain.LogRepository
}

func NewLogService(log domain.LogRepository) *LogService {
	return &LogService{log: log}
}

// LogEntry is a structured audit event. Operation is a stable verb such as
// "ticket.create"; Target identifies the affected object; Detail is free text
// (usually a JSON fragment or a human reason).
type LogEntry struct {
	UserID    uint
	Operation string
	Target    string
	Detail    string
}

// Log persists the entry. mm_logs.operation is a single TINYTEXT column, so the
// parts are joined with '|': "ticket.create|ticket:<uuid>|<detail>".
func (s *LogService) Log(ctx context.Context, e LogEntry) error {
	entry := &domain.Log{
		UUID:      uuid.NewString(),
		UserID:    e.UserID,
		Operation: fmt.Sprintf("%s|%s|%s", e.Operation, e.Target, e.Detail),
		Version:   1,
		CreateAt:  uint(time.Now().Unix()),
	}
	if _, err := s.log.Create(ctx, entry); err != nil {
		return fmt.Errorf("log %s: %w", e.Operation, err)
	}
	return nil
}

// List returns a keyset page of audit rows, newest first.
func (s *LogService) List(ctx context.Context, cursor int64, limit int) ([]*domain.Log, error) {
	return s.log.List(ctx, cursor, limit)
}
