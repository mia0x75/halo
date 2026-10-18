package meta_test

import (
	"testing"

	"github.com/mia0x75/halo/internal/agent/meta"
)

func TestColumnTypeClassification(t *testing.T) {
	text := []meta.ColumnType{"ENUM", "SET", "JSON", "JSONB", "CHAR", "VARCHAR", "NCHAR", "NVARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT"}
	blob := []meta.ColumnType{"BINARY", "VARBINARY", "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB", "BYTEA"}
	large := map[meta.ColumnType]bool{
		"TINYTEXT": true, "TEXT": true, "MEDIUMTEXT": true, "LONGTEXT": true,
		"TINYBLOB": true, "BLOB": true, "MEDIUMBLOB": true, "LONGBLOB": true,
	}

	for _, typ := range text {
		if !typ.IsText() {
			t.Errorf("%s: IsText() = false, want true", typ)
		}
		if typ.IsBlob() {
			t.Errorf("%s: IsBlob() = true, want false", typ)
		}
	}
	for _, typ := range blob {
		if !typ.IsBlob() {
			t.Errorf("%s: IsBlob() = false, want true", typ)
		}
		if typ.IsText() {
			t.Errorf("%s: IsText() = true, want false", typ)
		}
	}

	// IsLargeObject is the narrow subset rules use to reject payload columns; it
	// must exclude ordinary string types such as VARCHAR.
	all := append(append([]meta.ColumnType{}, text...), blob...)
	for _, typ := range all {
		if got := typ.IsLargeObject(); got != large[typ] {
			t.Errorf("%s: IsLargeObject() = %v, want %v", typ, got, large[typ])
		}
	}
	for _, typ := range []meta.ColumnType{"", "INT", "BIGINT", "DECIMAL", "DATETIME"} {
		if typ.IsText() || typ.IsBlob() || typ.IsLargeObject() {
			t.Errorf("%s: unexpected classification", typ)
		}
	}
}

func TestTableGetColumn(t *testing.T) {
	tbl := &meta.Table{Columns: map[string]*meta.Column{"id": {Name: "id"}}}
	if got := tbl.GetColumn("id"); got == nil || got.Name != "id" {
		t.Fatalf("GetColumn(id) = %+v, want id", got)
	}
	if got := tbl.GetColumn("missing"); got != nil {
		t.Errorf("GetColumn(missing) = %+v, want nil", got)
	}
}
