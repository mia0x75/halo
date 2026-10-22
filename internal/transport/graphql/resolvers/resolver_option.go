package resolvers

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/event"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/platform/crypto"
	"github.com/mia0x75/halo/internal/platform/notify"
	transportauth "github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// Option 获取某一具体的选项信息
func (r *queryRootResolver) Option(ctx context.Context, id string) (option *models.Option, err error) {
	rc := gqlapi.ReturnCodeOK
	cached := cache.OptionsMap.Any(func(elem *models.Option) bool {
		return elem.UUID == id
	})
	if cached == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 系统选项(uuid=%s)不存在。", rc, id)
		return
	}
	option = maskOption(cached)
	return
}

// Options 查看全部系统选项
func (r *queryRootResolver) Options(ctx context.Context) (L []*models.Option, err error) {
	for _, cached := range cache.OptionsMap.All() {
		L = append(L, maskOption(cached))
	}
	return
}

// maskOption returns a copy with the secret value hidden. The control type and
// the writable flag come from the database, so the settings page renders what
// the row says.
func maskOption(option *models.Option) *models.Option {
	if option == nil {
		return nil
	}
	copied := *option
	if copied.Secret() {
		if copied.Parameter != "" {
			copied.Parameter = "********"
		}
	}
	return &copied
}

// SystemInfo 平台运行信息
func (r *queryRootResolver) SystemInfo(ctx context.Context) (*gqlapi.SystemInfo, error) {
	database := ""
	if cfg := config.Config(); cfg != nil && cfg.Database != nil {
		database = cfg.Database.Driver
		if database == "" {
			database = "mysql"
		}
	}
	// 邮件开关以 mm_options 为准，配置文件不再带 mail 段。
	mailEnabled := false
	if option, optErr := r.app.Services.Option.GetByName(ctx, "mail.enabled"); optErr == nil {
		mailEnabled = option.Parameter == "true"
	}
	return &gqlapi.SystemInfo{
		Version:      app.Version,
		Git:          app.Git,
		BuildAt:      app.BuildAt,
		GoVer:        app.GoVer,
		Branch:       app.Branch,
		Database:     database,
		MailEnabled:  mailEnabled,
		ClusterCount: cache.ClustersMap.Count(),
		UserCount:    cache.UsersMap.Count(),
	}, nil
}

// PatchOptionParameters 修改系统选项
func (r *mutationRootResolver) PatchOptionParameters(ctx context.Context, input models.PatchOptionParameterInput) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		option := cache.OptionsMap.Any(func(elem *models.Option) bool {
			return elem.UUID == input.OptionUUID
		})
		if option == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 系统选项(uuid=%s)不存在。", rc, input.OptionUUID)
			break
		}
		if option.Writable != 1 {
			rc = gqlapi.ReturnCodeNotWritable
			err = fmt.Errorf("错误代码: %s, 错误信息: 系统选项(uuid=%s)不允许更新。", rc, input.OptionUUID)
			break
		}

		if option.Secret() {
			// 前端回显的是掩码；原样提交表示不改密码。
			if input.Parameter == "********" {
				ok = true
				break
			}
			key := ""
			if cfg := config.Config(); cfg != nil && cfg.Secret != nil {
				key = cfg.Secret.Crypto
			}
			if len(key) < 16 {
				rc = gqlapi.ReturnCodeUnknowError
				err = fmt.Errorf("错误代码: %s, 错误信息: 保存密钥需要一个至少 16 字节的 secret.crypto。", rc)
				break
			}
			ciphertext, encErr := crypto.EncryptAES([]byte(input.Parameter), key)
			if encErr != nil {
				rc = gqlapi.ReturnCodeUnknowError
				err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, encErr.Error())
				break
			}
			option.Parameter = base64.StdEncoding.EncodeToString(ciphertext)
		} else {
			option.Parameter = input.Parameter
		}
		option.BeforeUpdate()
		credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
		if err = r.app.Services.Option.Update(ctx, bridge.OptionToDomain(option), credential.User.UserID); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		event.Fire(event.EventOptionParameterPatched, &event.OptionParameterPatched{
			Manager: *credential.User,
			Option:  *option,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// SendTestMail sends a test message through the configured notification
// channels so an administrator can verify the SMTP settings.
func (r *mutationRootResolver) SendTestMail(ctx context.Context, email string) (bool, error) {
	enabled, err := r.app.Services.Option.GetByName(ctx, "mail.enabled")
	if err != nil && err != domain.ErrNotFound {
		return false, err
	}
	if err == nil && enabled.Parameter != "true" {
		return false, fmt.Errorf("邮件通知未启用，请先在系统选项里开启 mail.enabled。")
	}
	if err == domain.ErrNotFound {
		return false, fmt.Errorf("邮件通知未启用，请先在系统选项里开启 mail.enabled。")
	}

	msg := notify.Message{
		ToEmail: email,
		Subject: "[Halo] 测试邮件",
		Body:    "这是一封来自 Halo 的测试邮件。收到它说明邮件通道配置可用。",
	}
	if err := r.app.Notifier.Send(ctx, msg); err != nil {
		return false, fmt.Errorf("发送测试邮件失败: %w", err)
	}
	return true, nil
}
