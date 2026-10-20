package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"time"

	"github.com/jmoiron/sqlx"
)

// Engine identifies the metadata engine a migration set targets.
type Engine string

const (
	EngineMySQL    Engine = "mysql"
	EnginePostgres Engine = "postgres"
)

// Record is one row of the migration ledger.
type Record struct {
	Version    int
	Name       string
	Checksum   string
	DurationMS int
	Success    bool
	AppliedAt  time.Time
}

// Migrator applies schema migrations to the metadata database and keeps the
// mm_migrations ledger in sync.
type Migrator struct {
	db         *sqlx.DB
	engine     Engine
	migrations []Migration
}

// New builds a migrator for the given engine and migration set.
func New(db *sqlx.DB, engine Engine, fsys fs.FS) (*Migrator, error) {
	migrations, err := Load(fsys)
	if err != nil {
		return nil, err
	}
	return &Migrator{db: db, engine: engine, migrations: migrations}, nil
}

// Migrations returns the ordered migration set.
func (m *Migrator) Migrations() []Migration {
	return m.migrations
}

// Up applies every migration that is not yet recorded in the ledger.
//
// PostgreSQL runs each file inside a transaction, so a failure leaves neither
// the schema nor the ledger changed. MySQL commits implicitly on DDL: a failed
// file may have applied part of its statements, so the failure is recorded with
// success=false and later runs refuse to continue until an operator fixes the
// schema and clears that row.
func (m *Migrator) Up(ctx context.Context) ([]Migration, error) {
	if err := m.ensureLedger(ctx); err != nil {
		return nil, err
	}

	release, err := m.lock(ctx)
	if err != nil {
		return nil, err
	}
	defer release()

	applied, err := m.records(ctx)
	if err != nil {
		return nil, err
	}
	if err := m.verify(applied); err != nil {
		return nil, err
	}

	done := make([]Migration, 0, len(m.migrations))
	for _, mig := range m.migrations {
		record, ok := applied[mig.Version]
		if ok {
			if !record.Success {
				return done, fmt.Errorf(
					"迁移 %04d_%s 之前执行失败，需人工修复数据库后删除对应台账记录再重试",
					mig.Version, mig.Name)
			}
			continue
		}
		if err := m.apply(ctx, mig); err != nil {
			return done, err
		}
		done = append(done, mig)
	}

	return done, nil
}

// Status reports the state of every migration file against the ledger.
func (m *Migrator) Status(ctx context.Context) ([]Status, error) {
	if err := m.ensureLedger(ctx); err != nil {
		return nil, err
	}

	applied, err := m.records(ctx)
	if err != nil {
		return nil, err
	}

	status := make([]Status, 0, len(m.migrations))
	for _, mig := range m.migrations {
		record, ok := applied[mig.Version]
		switch {
		case !ok:
			status = append(status, Status{Migration: mig, State: StatePending})
		case !record.Success:
			status = append(status, Status{Migration: mig, State: StateFailed, Record: record})
		default:
			status = append(status, Status{Migration: mig, State: StateApplied, Record: record})
		}
	}

	return status, nil
}

// State describes a migration against the ledger.
type State string

const (
	StatePending State = "pending"
	StateApplied State = "applied"
	StateFailed  State = "failed"
)

// Status pairs a migration with its ledger state.
type Status struct {
	Migration Migration
	State     State
	Record    Record
}

// verify rejects a drifted ledger: a recorded file that no longer exists, a
// file whose content changed after it was applied, or a gap in the version
// sequence.
func (m *Migrator) verify(applied map[int]Record) error {
	byVersion := make(map[int]Migration, len(m.migrations))
	for _, mig := range m.migrations {
		byVersion[mig.Version] = mig
	}

	for version, record := range applied {
		mig, ok := byVersion[version]
		if !ok {
			return fmt.Errorf("台账中的迁移 %04d 在迁移目录中缺失", version)
		}
		if mig.Checksum != record.Checksum {
			return fmt.Errorf("迁移 %s 已被修改（checksum 不一致），已应用的迁移不可变更", mig.Filename)
		}
	}

	pending := make([]int, 0)
	for _, mig := range m.migrations {
		if _, ok := applied[mig.Version]; !ok {
			pending = append(pending, mig.Version)
		}
	}
	if len(pending) > 0 && len(applied) > 0 {
		maxApplied := 0
		for version := range applied {
			if version > maxApplied {
				maxApplied = version
			}
		}
		for _, version := range pending {
			if version < maxApplied {
				return fmt.Errorf("迁移版本不连续：%04d 未应用，但 %04d 已应用", version, maxApplied)
			}
		}
	}

	return nil
}

// records reads the ledger, indexed by version.
func (m *Migrator) records(ctx context.Context) (map[int]Record, error) {
	rows, err := m.db.QueryxContext(ctx, m.selectLedgerSQL())
	if err != nil {
		return nil, fmt.Errorf("read migration ledger: %w", err)
	}
	defer func() { _ = rows.Close() }()

	applied := make(map[int]Record)
	for rows.Next() {
		var (
			record  Record
			success sql.NullBool
		)
		if err := rows.Scan(&record.Version, &record.Name, &record.Checksum, &record.DurationMS, &success, &record.AppliedAt); err != nil {
			return nil, fmt.Errorf("scan migration ledger: %w", err)
		}
		record.Success = success.Valid && success.Bool
		applied[record.Version] = record
	}

	return applied, rows.Err()
}
