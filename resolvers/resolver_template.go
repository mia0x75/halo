package resolvers

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
)

// UpdateTemplate 修改一个工单
func (r mutationRootResolver) UpdateTemplate(ctx context.Context, input *models.UpdateTemplateInput) (template *models.Template, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		// credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		template = &models.Template{}
		if _, err = g.Engine.Where("`uuid` = ?", input.TemplateUUID).Get(template); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		if template.UUID == "" {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 邮件模板(uuid=%s)不存在。", rc, input.TemplateUUID)
			break
		}

		template.Subject = input.Subject
		template.Body = input.Body

		if _, err = g.Engine.Where("`uuid` = ?", input.TemplateUUID).Update(template); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		break
	}

	if err != nil {
		template = nil
	}

	return
}

// Templates 邮件模板在页面上是一个页面处理所有内容，所以不需要考虑分页
func (r queryRootResolver) Templates(ctx context.Context) (L []*models.Template, err error) {
	L = caches.TemplatesMap.All()
	return
}

type templateResolver struct{ *Resolver }
