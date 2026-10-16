package resolvers

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/event"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
	transportauth "github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// Option 获取某一具体的选项信息
func (r *queryRootResolver) Option(ctx context.Context, id string) (option *models.Option, err error) {
	rc := gqlapi.ReturnCodeOK
	option = cache.OptionsMap.Any(func(elem *models.Option) bool {
		return elem.UUID == id
	})
	if option == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 系统选项(uuid=%s)不存在。", rc, id)
	}
	return
}

// Options 查看全部系统选项
func (r *queryRootResolver) Options(ctx context.Context) (L []*models.Option, err error) {
	L = cache.OptionsMap.All()
	return
}

// PatchOptionValues 修改系统选项
func (r *mutationRootResolver) PatchOptionValues(ctx context.Context, input models.PatchOptionValueInput) (ok bool, err error) {
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
		option.Value = input.Value
		option.BeforeUpdate()
		credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
		if err = r.app.Services.Option.Update(ctx, bridge.OptionToDomain(option), credential.User.UserID); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		event.Fire(event.EventOptionValuePatched, &event.OptionValuePatchedArgs{
			Manager: *credential.User,
			Option:  *option,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}
