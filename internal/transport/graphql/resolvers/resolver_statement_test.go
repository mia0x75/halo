package resolvers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeStatementCursorUsesUpperBoundOnFirstPage(t *testing.T) {
	cursor, err := decodeStatementCursor(nil)
	require.NoError(t, err)
	require.Equal(t, statementCursorUpperBound, cursor)

	empty := ""
	cursor, err = decodeStatementCursor(&empty)
	require.NoError(t, err)
	require.Equal(t, statementCursorUpperBound, cursor)
}

func TestDecodeStatementCursorStripsPrefix(t *testing.T) {
	after := encodeStatementCursor("11111111-2222-3333-4444-555555555555")
	cursor, err := decodeStatementCursor(&after)
	require.NoError(t, err)
	require.Equal(t, "11111111-2222-3333-4444-555555555555", cursor)
}

func TestDecodeStatementCursorRejectsGarbage(t *testing.T) {
	after := "not-base64!!"
	_, err := decodeStatementCursor(&after)
	require.Error(t, err)
}
