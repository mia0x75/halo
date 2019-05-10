package resolvers

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/events"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// PatchRuleValues 修改规则值
func (r *mutationRootResolver) PatchRuleValues(ctx context.Context, input models.PatchRuleValuesInput) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		rule := caches.RulesMap.Any(func(elem *models.Rule) bool {
			if elem.UUID == input.RuleUUID {
				return true
			}
			return false
		})
		if rule == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 规则(uuid=%s)不存在。", rc, input.RuleUUID)
			break
		}
		if rule.Bitwise&2 != 2 {
			rc = gqlapi.ReturnCodeNotWritable
			err = fmt.Errorf("错误代码: %s, 错误信息: 规则(uuid=%s)不允许更新。", rc, input.RuleUUID)
			break
		}
		rule.Values = input.Values
		if _, err = g.Engine.Where("`uuid` = ?", rule.UUID).Update(rule); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
		}

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventRuleValuesPatched, &events.RuleValuesPatchedArgs{
			Manager: *credential.User,
			Rule:    *rule,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// PatchRuleBitwise 修改规则执行标志位
func (r *mutationRootResolver) PatchRuleBitwise(ctx context.Context, input models.PatchRuleBitwiseInput) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		rule := caches.RulesMap.Any(func(elem *models.Rule) bool {
			if elem.UUID == input.RuleUUID {
				return true
			}
			return false
		})
		if rule == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 规则(uuid=%s)不存在。", rc, input.RuleUUID)
			break
		}
		if rule.Bitwise&2 != 2 {
			rc = gqlapi.ReturnCodeNotWritable
			err = fmt.Errorf("错误代码: %s, 错误信息: 规则(uuid=%s)不允许更新。", rc, input.RuleUUID)
			break
		}
		enabled, _ := strconv.ParseBool(input.Enabled)
		value := uint8(0)
		if enabled {
			value = 1
		}
		rule.Bitwise = rule.Bitwise >> 1
		rule.Bitwise = rule.Bitwise << 1
		rule.Bitwise = rule.Bitwise | value

		// TODO: 需要单元测试
		if _, err = g.Engine.Where("`uuid` = ?", rule.UUID).Where("`bitwise` & 1 <> ?", value).Update(rule); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventRuleBitwisePatched, &events.RuleBitwisePatchedArgs{
			Manager: *credential.User,
			Rule:    *rule,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// TestRegexp 测试正则表达式的有效性
func (r *queryRootResolver) TestRegexp(ctx context.Context, input *models.ValidatePatternInput) (ok bool, err error) {
	rc := gqlapi.ReturnCodeOK
	if _, err = regexp.Compile(input.Pattern); err != nil {
		// TODO: 处理rc
		err = fmt.Errorf("错误代码: %s, 错误信息: `%s`不是一个有效的正则表达式。", rc, input.Pattern)
		return
	}
	ok = true

	return
}

// Rule 查看某一个具体的规则信息，好像目前用处不大
func (r *queryRootResolver) Rule(ctx context.Context, id string) (rule *models.Rule, err error) {
	rc := gqlapi.ReturnCodeOK
	rule = caches.RulesMap.Any(func(elem *models.Rule) bool {
		if elem.UUID == id {
			return true
		}
		return false
	})

	if rule == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 规则(uuid=%s)不存在。", rc, id)
	}
	return
}

// Rules 规则在页面上是一个页面处理所有内容，所以不需要考虑分页
func (r *queryRootResolver) Rules(ctx context.Context) (L []*models.Rule, err error) {
	L = caches.RulesMap.All()
	return
}

type ruleResolver struct{ *Resolver }
