package mysqlrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/mia0x75/halo/internal/domain"
)

// MySQL error 1062 is a duplicate-key violation; the repository surfaces it as
// domain.ErrConflict so callers need not depend on the driver.
const mysqlErrDuplicateEntry = 1062

// wrap adds operation context and maps driver-level conditions to the domain
// sentinels. Every repository method funnels its error through it.
func wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s: %w", op, domain.ErrNotFound)
	}
	var me *mysql.MySQLError
	if errors.As(err, &me) && me.Number == mysqlErrDuplicateEntry {
		return fmt.Errorf("%s: %w", op, domain.ErrConflict)
	}
	return fmt.Errorf("%s: %w", op, err)
}

// Timestamps are stored as DATETIME but the domain keeps the legacy Unix-second
// contract, so every repo converts at the row boundary. Reads return the instant
// in UTC; writes treat the domain value as an absolute Unix second.
func timeFromNull(t sql.NullTime) *uint {
	if !t.Valid {
		return nil
	}
	s := uint(t.Time.Unix())
	return &s
}

func timeToUnix(t time.Time) uint {
	return uint(t.Unix())
}

// The domain models nullable unsigned ints as *uint while sqlc uses
// sql.NullInt32 (MySQL UNSIGNED INT fits in int32 for the values this metadata
// schema stores: versions, foreign keys).

func uintFromNull(v sql.NullInt32) *uint {
	if !v.Valid {
		return nil
	}
	u := uint(v.Int32)
	return &u
}

func nullFromUint(v *uint) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: int32(*v), Valid: true}
}

func uint8FromNull(v sql.NullInt16) *uint8 {
	if !v.Valid {
		return nil
	}
	u := uint8(v.Int16)
	return &u
}

func nullFromUint8(v *uint8) sql.NullInt16 {
	if v == nil {
		return sql.NullInt16{}
	}
	return sql.NullInt16{Int16: int16(*v), Valid: true}
}

func stringFromNull(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func nullFromString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}
