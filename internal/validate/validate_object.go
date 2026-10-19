package validate

import (
	"strconv"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// The storage-program and view validators share the same name and existence
// checks, so they live here instead of being duplicated four times.

// nameRegexRule applies a regexp-based naming rule. withParams renders the rule
// parameters into the message (the "...需要满足正则\"%s\"。" family).
func nameRegexRule(s *Subject, r *Rule, name string, withParams bool) {
	var err error
	if withParams {
		err = Match(r, name, name, r.Parameters)
	} else {
		err = Match(r, name, name)
	}
	if err != nil {
		s.Violations.Append(&Clause{Description: err.Error(), Level: uint8(r.Level)})
	}
}

// namePrefixRule applies a prefix regexp rule; Match is used only as a boolean.
func namePrefixRule(s *Subject, r *Rule, name string) {
	if err := Match(r, name, r.Parameters); err != nil {
		appendViolation(s, r, name, r.Parameters)
	}
}

// nameMaxLengthRule flags names longer than the rule's numeric threshold.
func nameMaxLengthRule(s *Subject, r *Rule, name string) {
	threshold, _ := strconv.Atoi(r.Parameters)
	if len(name) > threshold {
		appendViolation(s, r, name, threshold)
	}
}

// requireDatabase flags when the target database is unknown to the metadata.
func (v *vldr) requireDatabase(s *Subject, r *Rule, schema string) {
	db := v.objectDB(schema)
	if v.DatabaseInfo(db) == nil {
		appendViolation(s, r, db)
	}
}

// requireObjectAbsent flags when an object already exists (CREATE rules).
func (v *vldr) requireObjectAbsent(s *Subject, r *Rule, schema string, kind meta.ObjectKind, name string) {
	if exists, known := v.lookupObject(v.objectDB(schema), kind, name); known && exists {
		appendViolation(s, r, name)
	}
}

// requireObjectPresent flags when an object is missing (ALTER/DROP rules).
func (v *vldr) requireObjectPresent(s *Subject, r *Rule, schema string, kind meta.ObjectKind, name string) {
	if exists, known := v.lookupObject(v.objectDB(schema), kind, name); known && !exists {
		appendViolation(s, r, name)
	}
}

// requireTable flags when the referenced base table does not exist.
func (v *vldr) requireTable(s *Subject, r *Rule, schema, table string) {
	db := v.objectDB(schema)
	if v.TableInfo(db, table) == nil {
		appendViolation(s, r, table)
	}
}
