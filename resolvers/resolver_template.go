package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/models"
)

// UpdateTemplate 修改一个工单
func (r *mutationRootResolver) UpdateTemplate(ctx context.Context, input *models.UpdateTemplateInput) (template *models.Template, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		// credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		t, getErr := r.app.Services.Template.GetByUUID(ctx, input.TemplateUUID)
		if getErr != nil {
			if errors.Is(getErr, domain.ErrNotFound) {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 邮件模板(uuid=%s)不存在。", rc, input.TemplateUUID)
			} else {
				rc = gqlapi.ReturnCodeUnknowError
				err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, getErr.Error())
			}
			break
		}
		template = bridge.TemplateFromDomain(t)

		template.Subject = input.Subject
		template.Body = input.Body
		template.BeforeUpdate()

		if err = r.app.Services.Template.Update(ctx, bridge.TemplateToDomain(template), 0); err != nil {
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
func (r *queryRootResolver) Templates(ctx context.Context) (L []*models.Template, err error) {
	L = caches.TemplatesMap.All()
	return
}
