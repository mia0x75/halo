package app

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/platform/crypto"
	"github.com/mia0x75/halo/internal/platform/notify"
)

// Mail option names. They live in mm_options so an administrator can edit the
// SMTP settings from the settings page.
const (
	mailEnabled            = "mail.enabled"
	mailAddr               = "mail.addr"
	mailUser               = "mail.user"
	mailPassword           = "mail.password"
	mailInsecureSkipVerify = "mail.insecure-skip-verify"
)

// mailSettingsProvider resolves the SMTP settings from mm_options. The
// options are the only source: the config file no longer carries a mail
// section, so a missing option simply means "not configured".
func mailSettingsProvider(options domain.OptionRepository) notify.MailSettingsFunc {
	return func(ctx context.Context) (notify.MailSettings, error) {
		settings := notify.MailSettings{}

		read := func(name string) (string, bool, error) {
			option, err := options.GetByName(ctx, name)
			if err != nil {
				if err == domain.ErrNotFound {
					return "", false, nil
				}
				return "", false, fmt.Errorf("read %s: %w", name, err)
			}
			return option.Parameter, option.Parameter != "", nil
		}

		if value, ok, err := read(mailEnabled); err != nil {
			return settings, err
		} else if ok {
			settings.Enabled = value == "true"
		}
		if value, ok, err := read(mailAddr); err != nil {
			return settings, err
		} else if ok {
			settings.Addr = value
		}
		if value, ok, err := read(mailUser); err != nil {
			return settings, err
		} else if ok {
			settings.User = value
		}
		if value, ok, err := read(mailInsecureSkipVerify); err != nil {
			return settings, err
		} else if ok {
			settings.InsecureSkipVerify = value == "true"
		}
		if value, ok, err := read(mailPassword); err != nil {
			return settings, err
		} else if ok {
			// The password is stored base64-encoded because AES ciphertext is
			// not valid UTF-8 and the column is a text type.
			raw, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				return settings, fmt.Errorf("mail password is not valid base64: %w", err)
			}
			key := ""
			if cfg := config.Config(); cfg != nil && cfg.Secret != nil {
				key = cfg.Secret.Crypto
			}
			if len(key) < 16 {
				return settings, fmt.Errorf("mail password cannot be decrypted: secret.crypto 至少需要 16 字节")
			}
			plain, err := crypto.DecryptAES(raw, key)
			if err != nil {
				return settings, fmt.Errorf("mail password cannot be decrypted: %w", err)
			}
			settings.Password = string(plain)
		}

		return settings, nil
	}
}
