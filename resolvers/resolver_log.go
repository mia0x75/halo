package resolvers

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
)

type logResolver struct{ *Resolver }

// User 操作日志的发起人（委托人）
func (r *logResolver) User(ctx context.Context, obj *models.Log) (user *models.User, err error) {
	rc := gqlapi.ReturnCodeOK
	user = caches.UsersMap.Any(func(elem *models.User) bool {
		if elem.UserID == obj.UserID {
			return true
		}
		return false
	})
	if user == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 用户不存在。", rc)
	}
	return
}

// Logs 操作日志分页查看
func (r *queryRootResolver) Logs(ctx context.Context, after *string, before *string, first *int, last *int) (data *gqlapi.LogConnection, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		// 参数判断，只允许 first/before first/after last/before last/after 模式
		if first != nil && last != nil {
			rc = gqlapi.ReturnCodeInvalidParams
			err = fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", rc)
			break L
		}
		if after != nil && before != nil {
			rc = gqlapi.ReturnCodeInvalidParams
			err = fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", rc)
			break L
		}

		from := math.MaxInt64
		if after != nil {
			var bs []byte
			bs, err = base64.StdEncoding.DecodeString(*after)
			if err != nil {
				rc = gqlapi.ReturnCodeInvalidParams
				err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
				break L
			}
			var i int
			i, err = strconv.Atoi(strings.TrimPrefix(string(bs), "cursor"))
			if err != nil {
				rc = gqlapi.ReturnCodeInvalidParams
				err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
				break L
			}
			from = i
		}
		hasPreviousPage := true
		hasNextPage := true

		if from == math.MaxInt64 {
			hasPreviousPage = false
		}
		// 获取edges
		edges := []*gqlapi.LogEdge{}
		logs := []*models.Log{}
		if err = g.Engine.Desc("log_id").Where("log_id < ?", from).Limit(*first).Find(&logs); err != nil {
			return nil, err
		}

		for _, log := range logs {
			edges = append(edges, &gqlapi.LogEdge{
				Node:   log,
				Cursor: EncodeCursor(fmt.Sprintf("%d", log.LogID)),
			})
		}

		if len(edges) < *first {
			hasNextPage = false
		}

		if len(edges) == 0 {
			return nil, nil
		}
		// 获取pageInfo
		startCursor := EncodeCursor(fmt.Sprintf("%d", logs[0].LogID))
		endCursor := EncodeCursor(fmt.Sprintf("%d", logs[len(logs)-1].LogID))
		pageInfo := &gqlapi.PageInfo{
			HasPreviousPage: hasPreviousPage,
			HasNextPage:     hasNextPage,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		}

		data = &gqlapi.LogConnection{
			PageInfo:   pageInfo,
			Edges:      edges,
			TotalCount: len(edges),
		}

		break
	}

	return
}
