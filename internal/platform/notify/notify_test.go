package notify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSanitizeHeaderStripsCRLF guards against SMTP header injection: a value
// containing CR/LF must not be able to introduce additional headers or end the
// header block.
func TestSanitizeHeaderStripsCRLF(t *testing.T) {
	require.Equal(t, "hello", sanitizeHeader("hello"))
	require.Equal(t, "aBcc: evil@example.com", sanitizeHeader("a\r\nBcc: evil@example.com"))
	require.Equal(t, "subject", sanitizeHeader("sub\nject"))
	require.Equal(t, "xInjected body", sanitizeHeader("x\r\n\r\nInjected body"))
}
