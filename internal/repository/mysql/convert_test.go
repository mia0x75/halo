package mysqlrepo

import (
	"database/sql"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestClampInt32(t *testing.T) {
	require.Equal(t, int32(0), clampInt32(-1))
	require.Equal(t, int32(0), clampInt32(0))
	require.Equal(t, int32(50), clampInt32(50))
	require.Equal(t, int32(math.MaxInt32), clampInt32(math.MaxInt32))
	require.Equal(t, int32(math.MaxInt32), clampInt32(math.MaxInt32+1))
}

func TestClampUint(t *testing.T) {
	require.Equal(t, uint32(0), clampUint(0))
	require.Equal(t, uint32(7), clampUint(7))
	require.Equal(t, uint32(math.MaxUint32), clampUint(math.MaxUint32))
	require.Equal(t, uint32(math.MaxUint32), clampUint(math.MaxUint32+1))
}

// TestWrapMapsSentinels locks the driver-to-domain error mapping every repo
// method relies on.
func TestWrapMapsSentinels(t *testing.T) {
	require.NoError(t, wrap("op", nil))

	err := wrap("op", sql.ErrNoRows)
	require.ErrorIs(t, err, domain.ErrNotFound)
	require.Contains(t, err.Error(), "op")

	err = wrap("op", &mysql.MySQLError{Number: mysqlErrDuplicateEntry, Message: "dup"})
	require.ErrorIs(t, err, domain.ErrConflict)

	boom := errors.New("dial")
	err = wrap("op", boom)
	require.ErrorIs(t, err, boom)
	require.NotErrorIs(t, err, domain.ErrNotFound)
}

func TestTimeConversions(t *testing.T) {
	require.Nil(t, timeFromNull(sql.NullTime{}))

	at := time.Unix(1700000000, 0)
	got := timeFromNull(sql.NullTime{Time: at, Valid: true})
	require.NotNil(t, got)
	require.Equal(t, uint(1700000000), *got)
	require.Equal(t, uint(1700000000), timeToUnix(at))
}

func TestUintConversions(t *testing.T) {
	require.Nil(t, uintFromNull(sql.NullInt32{}))
	u := uintFromNull(sql.NullInt32{Int32: 12, Valid: true})
	require.NotNil(t, u)
	require.Equal(t, uint(12), *u)

	require.Equal(t, sql.NullInt32{}, nullFromUint(nil))
	seven := uint(7)
	require.Equal(t, sql.NullInt32{Int32: 7, Valid: true}, nullFromUint(&seven))
}

func TestUint8Conversions(t *testing.T) {
	require.Nil(t, uint8FromNull(sql.NullInt16{}))
	u := uint8FromNull(sql.NullInt16{Int16: 3, Valid: true})
	require.NotNil(t, u)
	require.Equal(t, uint8(3), *u)

	require.Equal(t, sql.NullInt16{}, nullFromUint8(nil))
	three := uint8(3)
	require.Equal(t, sql.NullInt16{Int16: 3, Valid: true}, nullFromUint8(&three))
}

func TestStringConversions(t *testing.T) {
	require.Nil(t, stringFromNull(sql.NullString{}))
	s := stringFromNull(sql.NullString{String: "x", Valid: true})
	require.NotNil(t, s)
	require.Equal(t, "x", *s)

	require.Equal(t, sql.NullString{}, nullFromString(nil))
	y := "y"
	require.Equal(t, sql.NullString{String: "y", Valid: true}, nullFromString(&y))
}
