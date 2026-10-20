package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// lockName is the advisory lock key. It is identical on both engines so a
// cluster of Halo instances sharing one metadata database serializes migrations.
const lockName = "halo_migrate"

// ledgerDDL creates the ledger table. It is owned by the migrator, therefore
// every statement is idempotent.
func ledgerDDL(engine Engine) string {
	if engine == EnginePostgres {
		return `CREATE TABLE IF NOT EXISTS mm_migrations (
    version    INTEGER NOT NULL,
    name       VARCHAR(75) NOT NULL,
    checksum   CHAR(64) NOT NULL,
    duration_ms INTEGER NOT NULL DEFAULT 0,
    success    BOOLEAN NOT NULL DEFAULT TRUE,
    apply_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (version),
    CONSTRAINT uq_mm_migrations_name UNIQUE (name)
)`
	}

	return "CREATE TABLE IF NOT EXISTS `mm_migrations` (" +
		"`version` INT UNSIGNED NOT NULL COMMENT '迁移版本号'," +
		"`name` VARCHAR(75) NOT NULL COMMENT '迁移名称'," +
		"`checksum` CHAR(64) NOT NULL COMMENT '文件内容 sha256'," +
		"`duration_ms` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '执行耗时（毫秒）'," +
		"`success` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否成功'," +
		"`apply_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '应用时间'," +
		"PRIMARY KEY (`version`)," +
		"UNIQUE KEY `unique_1` (`name`)" +
		") ENGINE = InnoDB CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '迁移版本台账'"
}

func (m *Migrator) selectLedgerSQL() string {
	if m.engine == EnginePostgres {
		return `SELECT version, name, checksum, duration_ms, success, apply_at FROM mm_migrations ORDER BY version`
	}
	return "SELECT `version`, `name`, `checksum`, `duration_ms`, `success`, `apply_at` FROM `mm_migrations` ORDER BY `version`"
}

func (m *Migrator) insertLedgerSQL() string {
	if m.engine == EnginePostgres {
		return `INSERT INTO mm_migrations (version, name, checksum, duration_ms, success, apply_at)
VALUES ($1, $2, $3, $4, $5, NOW())`
	}
	return "INSERT INTO `mm_migrations` (`version`, `name`, `checksum`, `duration_ms`, `success`) " +
		"VALUES (?, ?, ?, ?, ?)"
}

// ensureLedger creates the ledger table when it is missing.
func (m *Migrator) ensureLedger(ctx context.Context) error {
	if _, err := m.db.ExecContext(ctx, ledgerDDL(m.engine)); err != nil {
		return fmt.Errorf("create migration ledger: %w", err)
	}
	return nil
}

// lock takes the engine specific advisory lock so that concurrent instances do
// not apply the same migration twice.
func (m *Migrator) lock(ctx context.Context) (func(), error) {
	if m.engine == EnginePostgres {
		if _, err := m.db.ExecContext(ctx, "SELECT pg_advisory_lock(hashtext($1))", lockName); err != nil {
			return nil, fmt.Errorf("acquire migration lock: %w", err)
		}
		return func() {
			_, _ = m.db.ExecContext(context.WithoutCancel(ctx), "SELECT pg_advisory_unlock(hashtext($1))", lockName)
		}, nil
	}

	var acquired sql.NullInt64
	if err := m.db.QueryRowxContext(ctx, "SELECT GET_LOCK(?, 30)", lockName).Scan(&acquired); err != nil {
		return nil, fmt.Errorf("acquire migration lock: %w", err)
	}
	if !acquired.Valid || acquired.Int64 != 1 {
		return nil, fmt.Errorf("acquire migration lock: timeout waiting for %q", lockName)
	}

	return func() {
		_, _ = m.db.ExecContext(context.WithoutCancel(ctx), "SELECT RELEASE_LOCK(?)", lockName)
	}, nil
}

// apply executes one migration file and records the outcome in the ledger.
func (m *Migrator) apply(ctx context.Context, mig Migration) error {
	start := time.Now()

	// A migration folded into the baseline has already been applied by the
	// baseline snapshot; it is only recorded so the ledger keeps one row per
	// migration file.
	if mig.Folded {
		if err := m.record(ctx, mig, start, true); err != nil {
			return fmt.Errorf("迁移 %s 写入台账失败: %w", mig.Filename, err)
		}
		return nil
	}

	if m.engine == EnginePostgres {
		return m.applyTransactional(ctx, mig, start)
	}

	// MySQL commits implicitly on DDL, so the statements cannot be rolled back.
	for _, stmt := range mig.Statements {
		if _, err := m.db.ExecContext(ctx, stmt); err != nil {
			if recordErr := m.record(ctx, mig, start, false); recordErr != nil {
				return fmt.Errorf("迁移 %s 执行失败: %w（写入失败台账时出错: %v）", mig.Filename, err, recordErr)
			}
			return fmt.Errorf("迁移 %s 执行失败: %w", mig.Filename, err)
		}
	}
	if err := m.record(ctx, mig, start, true); err != nil {
		return fmt.Errorf("迁移 %s 已执行但写入台账失败: %w", mig.Filename, err)
	}
	return nil
}

// applyTransactional runs the statements and the ledger insert in one
// transaction. PostgreSQL supports transactional DDL, so a failure leaves the
// database untouched and the migration stays pending.
func (m *Migrator) applyTransactional(ctx context.Context, mig Migration, start time.Time) error {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("迁移 %s 开启事务失败: %w", mig.Filename, err)
	}
	defer func() { _ = tx.Rollback() }()

	for _, stmt := range mig.Statements {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("迁移 %s 执行失败: %w", mig.Filename, err)
		}
	}

	duration := int(time.Since(start).Milliseconds())
	if _, err := tx.ExecContext(ctx, m.insertLedgerSQL(), mig.Version, mig.Name, mig.Checksum, duration, true); err != nil {
		return fmt.Errorf("迁移 %s 写入台账失败: %w", mig.Filename, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("迁移 %s 提交失败: %w", mig.Filename, err)
	}
	return nil
}

// record writes the ledger row of a migration.
func (m *Migrator) record(ctx context.Context, mig Migration, start time.Time, success bool) error {
	duration := int(time.Since(start).Milliseconds())
	if _, err := m.db.ExecContext(ctx, m.insertLedgerSQL(), mig.Version, mig.Name, mig.Checksum, duration, success); err != nil {
		return err
	}
	return nil
}
