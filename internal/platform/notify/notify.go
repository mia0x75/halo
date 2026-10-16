// Package notify delivers notifications to external channels. Delivery is
// best-effort: callers (the event bus) log and drop failures.
package notify

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/platform/config"
)

// Message is a channel-agnostic notification.
type Message struct {
	ToName  string
	ToEmail string
	Subject string
	Body    string
}

// Notifier delivers a message over one channel.
type Notifier interface {
	Name() string
	Send(ctx context.Context, msg Message) error
}

// Multi fans a message out to every configured notifier. It returns the
// joined errors but always attempts every channel.
type Multi struct {
	notifiers []Notifier
}

// NewMulti builds a fan-out notifier. Nil notifiers are ignored.
func NewMulti(notifiers ...Notifier) *Multi {
	ns := make([]Notifier, 0, len(notifiers))
	for _, n := range notifiers {
		if n != nil {
			ns = append(ns, n)
		}
	}
	return &Multi{notifiers: ns}
}

// Name identifies the fan-out notifier.
func (m *Multi) Name() string { return "multi" }

// Send delivers to every notifier, collecting errors.
func (m *Multi) Send(ctx context.Context, msg Message) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Send(ctx, msg); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", n.Name(), err))
		}
	}
	return errors.Join(errs...)
}

// EmailNotifier sends mail over SMTP with implicit TLS (port 465 style).
type EmailNotifier struct{}

// Name identifies the email channel.
func (EmailNotifier) Name() string { return "email" }

// Send delivers an email. It is a no-op when mail is disabled in config.
func (EmailNotifier) Send(_ context.Context, msg Message) error {
	cfg := config.Config().Mail
	if cfg == nil || !cfg.Enabled {
		return nil
	}

	// Header values are attacker-influenced (ticket subject, user name/email),
	// so strip CR/LF to prevent SMTP header injection, and reject a recipient
	// that is not a valid address.
	from := mail.Address{Name: sanitizeHeader("系统用户"), Address: cfg.User}
	to, err := mail.ParseAddress(msg.ToEmail)
	if err != nil {
		return fmt.Errorf("invalid recipient %q: %w", msg.ToEmail, err)
	}
	to.Name = sanitizeHeader(msg.ToName)

	header := textproto.MIMEHeader{}
	header.Set(textproto.CanonicalMIMEHeaderKey("from"), fmt.Sprintf("%s <%s>", from.Name, from.Address))
	header.Set(textproto.CanonicalMIMEHeaderKey("to"), fmt.Sprintf("%s <%s>", to.Name, to.Address))
	header.Set(textproto.CanonicalMIMEHeaderKey("content-type"), "text/plain; charset=UTF-8")
	header.Set(textproto.CanonicalMIMEHeaderKey("mime-version"), "1.0")
	header.Set(textproto.CanonicalMIMEHeaderKey("subject"), sanitizeHeader(msg.Subject))

	var buffer bytes.Buffer
	for key, value := range header {
		fmt.Fprintf(&buffer, "%s: %s\r\n", key, value[0])
	}
	fmt.Fprintf(&buffer, "\r\n%s", msg.Body)

	servername := cfg.Addr
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", cfg.User, cfg.Password, host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	if err = c.Auth(auth); err != nil {
		return err
	}
	if err = c.Mail(from.Address); err != nil {
		return err
	}
	if err = c.Rcpt(to.Address); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(buffer.Bytes()); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

// sanitizeHeader strips CR and LF so an attacker-controlled value cannot inject
// additional SMTP headers or terminate the header block.
func sanitizeHeader(s string) string {
	if !strings.ContainsAny(s, "\r\n") {
		return s
	}
	return strings.NewReplacer("\r", "", "\n", "").Replace(s)
}

// DingTalkNotifier is a placeholder for the DingTalk webhook channel.
type DingTalkNotifier struct {
	Webhook string
	Secret  string
}

// Name identifies the DingTalk channel.
func (DingTalkNotifier) Name() string { return "dingtalk" }

// Send is not implemented yet. A missing webhook is a silent no-op so the
// channel can be enabled purely by configuration.
func (n DingTalkNotifier) Send(_ context.Context, _ Message) error {
	if n.Webhook == "" {
		return nil
	}
	// TODO: POST markdown payload to n.Webhook, signing with n.Secret.
	log.Warnf("[W] dingtalk notifier not implemented, dropping message")
	return nil
}

// WeChatNotifier is a placeholder for the WeCom webhook channel.
type WeChatNotifier struct {
	Webhook string
}

// Name identifies the WeChat channel.
func (WeChatNotifier) Name() string { return "wechat" }

// Send is not implemented yet. A missing webhook is a silent no-op.
func (n WeChatNotifier) Send(_ context.Context, _ Message) error {
	if n.Webhook == "" {
		return nil
	}
	// TODO: POST markdown payload to n.Webhook.
	log.Warnf("[W] wechat notifier not implemented, dropping message")
	return nil
}
