package postgresrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/mia0x75/halo/internal/domain"
)

// pgErrUniqueViolation is the SQLSTATE for duplicate-key violations; the
// repository surfaces it as domain.ErrConflict so callers need not depend on
// the driver.
const pgErrUniqueViolation = "23505"

// wrap adds operation context and maps driver-level conditions to the domain
// sentinels. Every repository method funnels its error through it.
func wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s: %w", op, domain.ErrNotFound)
	}
	var pe *pgconn.PgError
	if errors.As(err, &pe) && pe.Code == pgErrUniqueViolation {
		return fmt.Errorf("%s: %w", op, domain.ErrConflict)
	}
	return fmt.Errorf("%s: %w", op, err)
}

// Timestamps are stored as TIMESTAMPTZ but the domain keeps the legacy
// Unix-second contract, so every repo converts at the row boundary.
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

// PG has no unsigned integers: BIGINT columns surface as int64 while the
// domain uses uint, so every repo converts at the row boundary.

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

func uintFromNullInt64(v sql.NullInt64) *uint {
	if !v.Valid {
		return nil
	}
	u := uint(v.Int64)
	return &u
}

func nullInt64FromUint(v *uint) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*v), Valid: true}
}

func uint8FromNullInt16(v sql.NullInt16) *uint8 {
	if !v.Valid {
		return nil
	}
	u := uint8(v.Int16)
	return &u
}

func nullInt16FromUint8(v *uint8) sql.NullInt16 {
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

// Phone numbers exceed int64 range in the worst case, so mm_users.phone is
// NUMERIC(20,0) read as text; the domain keeps uint64.

func parsePhone(s string) (uint64, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("phone %q: %w", s, err)
	}
	return n, nil
}

func formatPhone(n uint64) string {
	return strconv.FormatUint(n, 10)
}
