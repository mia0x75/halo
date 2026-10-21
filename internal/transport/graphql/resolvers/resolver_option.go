package resolvers

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/event"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
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

// maskOption returns a copy with the secret value hidden and the control type
// filled in from the option spec.
func maskOption(option *models.Option) *models.Option {
	if option == nil {
		return nil
	}
	copied := *option
	copied.Parameter = domain.MaskOptionValue(option.Name, option.Parameter)
	spec := domain.OptionSpecFor(option.Name)
	copied.Element = spec.Element
	copied.Writable = boolToUint8(spec.Writable)
	return &copied
}

func boolToUint8(value bool) uint8 {
	if value {
		return 1
	}
	return 0
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
	mailEnabled := false
	if cfg := config.Config(); cfg != nil && cfg.Mail != nil {
		mailEnabled = cfg.Mail.Enabled
	}
	return &gqlapi.SystemInfo{
		Version:      app.Version,
		Git:          app.Git,
		Compile:      app.Compile,
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

		option.Parameter = input.Parameter
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
