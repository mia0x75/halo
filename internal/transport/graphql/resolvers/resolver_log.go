package resolvers

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

type logResolver struct{ *Resolver }

// User 操作日志的发起人（委托人）。
// 用户被删除后日志仍要能翻出来，理由同 deletedUserPlaceholder。
func (r *logResolver) User(ctx context.Context, obj *models.Log) (*models.User, error) {
	user := cache.UsersMap.Any(func(elem *models.User) bool {
		return elem.UserID == obj.UserID
	})
	if user != nil {
		return user, nil
	}
	return deletedUserPlaceholder(obj.UserID), nil
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
		hasPreviousPage := from != math.MaxInt64
		limit := *first

		// 纯游标分页：多取一条探测是否有下一页，不统计全表总数（日志有
		// 保留期清理，总数既贵又会跳动）。
		domainLogs, listErr := r.app.Services.Log.List(ctx, int64(from), probeLimit(limit))
		if listErr != nil {
			return nil, listErr
		}
		pageLogs, hasNextPage := trimKeysetPage(domainLogs, limit)

		edges := []*gqlapi.LogEdge{}
		for _, l := range pageLogs {
			log := bridge.LogFromDomain(l)
			edges = append(edges, &gqlapi.LogEdge{
				Node:   log,
				Cursor: EncodeCursor(fmt.Sprintf("%d", log.LogID)),
			})
		}

		if len(edges) == 0 {
			return nil, nil
		}
		// 获取pageInfo
		startCursor := edges[0].Cursor
		endCursor := edges[len(edges)-1].Cursor
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
