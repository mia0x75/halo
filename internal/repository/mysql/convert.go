package mysqlrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/mia0x75/halo/internal/domain"
)

// MySQL error 1062 is a duplicate-key violation; the repository surfaces it as
// domain.ErrConflict so callers need not depend on the driver.
const mysqlErrDuplicateEntry = 1062

// clampUint32 narrows an int64 id/cursor to the schema's UNSIGNED INT range.
// Out-of-range values are clamped rather than silently wrapped, so a malformed
// cursor (negative, or beyond 2^32) cannot address an unrelated row.
func clampUint32(v int64) uint32 {
	if v <= 0 {
		return 0
	}
	if v > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(v)
}

// clampInt32 narrows a page size to int32, clamping to [0, MaxInt32] so an
// oversized limit cannot wrap negative and break the LIMIT clause.
func clampInt32(v int) int32 {
	if v <= 0 {
		return 0
	}
	if int64(v) > math.MaxInt32 {
		return math.MaxInt32
	}
	return int32(v)
}

// clampUint narrows a uint to the schema's UNSIGNED INT range, clamping instead
// of truncating on 64-bit platforms.
func clampUint(v uint) uint32 {
	if v > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(v)
}

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

// intFromNull converts a nullable signed int column.
func intFromNull(v sql.NullInt32) *int {
	if !v.Valid {
		return nil
	}
	i := int(v.Int32)
	return &i
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

// nullFromTime converts a unix second pointer into a nullable timestamp.
func nullFromTime(v *uint) sql.NullTime {
	if v == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: time.Unix(int64(*v), 0).UTC(), Valid: true}
}
