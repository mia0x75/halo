// Package bridge converts between the legacy models.* structs (now tag-cleaned
// and still consumed by caches and resolvers) and the internal/domain entities.
// It is a temporary shim for the 2D migration: once models.* is deleted the
// converters go with it.

package bridge

import "database/sql"

// uintPtr boxes a uint. Legacy models model nullable columns as non-nullable
// values, so a domain *uint round-trips through a plain uint.
func uintPtr(u uint) *uint { return &u }

// deref unboxes a *uint, mapping nil to 0.
func deref(p *uint) uint {
	if p == nil {
		return 0
	}
	return *p
}

func strPtr(s string) *string { return &s }

func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func u8Ptr(v uint8) *uint8 { return &v }

func derefU8(p *uint8) uint8 {
	if p == nil {
		return 0
	}
	return *p
}

// nullInt64 boxes a nullable cron id for the legacy sql.NullInt64 column.
func nullInt64(u *uint) sql.NullInt64 {
	if u == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*u), Valid: true}
}

func uintFromNullInt64(n sql.NullInt64) *uint {
	if !n.Valid {
		return nil
	}
	v := uint(n.Int64)
	return &v
}
