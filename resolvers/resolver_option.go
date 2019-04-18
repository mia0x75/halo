package resolvers

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/events"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// Option 获取某一具体的选项信息
func (r *queryRootResolver) Option(ctx context.Context, id string) (option *models.Option, err error) {
	rc := gqlapi.ReturnCodeOK
	option = caches.OptionsMap.Any(func(elem *models.Option) bool {
		if elem.UUID == id {
			return true
		}
		return false
	})
	if option == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 系统选项(uuid=%s)不存在。", rc, id)
	}
	return
}

// Options 查看全部系统选项
func (r *queryRootResolver) Options(ctx context.Context) (L []*models.Option, err error) {
	L = caches.OptionsMap.All()
	return
}

type optionResolver struct{ *Resolver }

// PatchOptionValues 修改系统选项
func (r *mutationRootResolver) PatchOptionValues(ctx context.Context, input models.PatchOptionValueInput) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		option := caches.OptionsMap.Any(func(elem *models.Option) bool {
			if elem.UUID == input.OptionUUID {
				return true
			}
			return false
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
		// TODO: 需要单元测试
		if _, err = g.Engine.Where("`uuid` = ?", option.UUID).Update(option); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventOptionValuePatched, &events.OptionValuePatchedArgs{
			Manager: *credential.User,
			Option:  *option,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}
