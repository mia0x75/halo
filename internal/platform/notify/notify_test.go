package notify

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/platform/config"
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

func TestSanitizeBodyKeepsNewlineDropsControl(t *testing.T) {
	require.Equal(t, "line1\nline2", sanitizeBody("line1\nline2"))
	// Bare CR and other control characters are dropped; \n survives.
	require.Equal(t, "a\nb", sanitizeBody("a\r\nb"))
	require.Equal(t, "ab", sanitizeBody("a\x00\x1bb"))
	require.Equal(t, "tab", sanitizeBody("ta\tb"))
	require.Equal(t, "", sanitizeBody("\x00\x01\x02"))
}

type recordingNotifier struct {
	name string
	sent []Message
	err  error
}

func (r *recordingNotifier) Name() string { return r.name }

func (r *recordingNotifier) Send(_ context.Context, msg Message) error {
	r.sent = append(r.sent, msg)
	return r.err
}

func TestMultiFansOutAndJoinsErrors(t *testing.T) {
	ok := &recordingNotifier{name: "ok"}
	bad := &recordingNotifier{name: "bad", err: errors.New("boom")}
	m := NewMulti(nil, ok, bad, nil)

	require.Equal(t, "multi", m.Name())
	err := m.Send(context.Background(), Message{Subject: "s"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "bad: boom")
	// Both channels are attempted even when one fails.
	require.Len(t, ok.sent, 1)
	require.Len(t, bad.sent, 1)
}

func TestMultiWithoutNotifiers(t *testing.T) {
	m := NewMulti()
	require.NoError(t, m.Send(context.Background(), Message{}))
}

// loadMailConfig points the global config at a temp file with the given mail
// block so EmailNotifier can be exercised without a live SMTP server.
func loadMailConfig(t *testing.T, mail string) {
	t.Helper()
	body := `{"mail":` + mail + `}`
	path := filepath.Join(t.TempDir(), "cfg.json")
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))
	config.ParseConfig(path)
}

func TestEmailNotifierDisabledIsNoOp(t *testing.T) {
	loadMailConfig(t, `{"enabled":false,"addr":"smtp.test:465"}`)
	require.NoError(t, EmailNotifier{}.Send(context.Background(), Message{ToEmail: "a@b.test"}))
	require.Equal(t, "email", EmailNotifier{}.Name())
}

func TestEmailNotifierRejectsInvalidRecipient(t *testing.T) {
	loadMailConfig(t, `{"enabled":true,"addr":"127.0.0.1:1","user":"halo"}`)
	err := EmailNotifier{}.Send(context.Background(), Message{ToEmail: "not an address"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid recipient")
}

// TestEmailNotifierInjectionSurvivesDialFailure documents that CRLF in the
// recipient/body cannot abort Send before the transport stage: the call reaches
// the network and fails there, not on malformed input.
func TestEmailNotifierDialFailure(t *testing.T) {
	loadMailConfig(t, `{"enabled":true,"addr":"127.0.0.1:1","user":"halo"}`)
	err := EmailNotifier{}.Send(context.Background(), Message{
		ToName:  "Eve\r\nBcc: evil@example.com",
		ToEmail: "victim@example.com",
		Subject: "hi\r\nX-Evil: 1",
		Body:    "body",
	})
	require.Error(t, err)
	require.False(t, strings.Contains(err.Error(), "invalid recipient"))
}

func TestPlaceholderChannelsAreNoOps(t *testing.T) {
	require.Equal(t, "dingtalk", DingTalkNotifier{}.Name())
	require.NoError(t, DingTalkNotifier{}.Send(context.Background(), Message{}))
	require.NoError(t, DingTalkNotifier{Webhook: "https://example.test"}.Send(context.Background(), Message{}))

	require.Equal(t, "wechat", WeChatNotifier{}.Name())
	require.NoError(t, WeChatNotifier{}.Send(context.Background(), Message{}))
	require.NoError(t, WeChatNotifier{Webhook: "https://example.test"}.Send(context.Background(), Message{}))
}
