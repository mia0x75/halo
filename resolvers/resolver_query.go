package resolvers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/mia0x75/parser"
	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/events"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// Query 查看某一查询的信息
func (r *queryRootResolver) Query(ctx context.Context, id string) (query *models.Query, err error) {
	rc := gqlapi.ReturnCodeOK
	found := false
	query = &models.Query{
		UUID: id,
	}
	if found, err = g.Engine.Get(query); err != nil {
		query = nil
		rc = gqlapi.ReturnCodeUnknowError
		err = fmt.Errorf("错误代码: %s, 错误信息: %s。", rc, err.Error())
	} else if !found {
		query = nil
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 查询(uuid=%s)不存在。", rc, id)
	}
	return
}

// Queries 分页查询数据查询
func (r *queryRootResolver) Queries(ctx context.Context, after *string, before *string, first *int, last *int) (*gqlapi.QueryConnection, error) {
	rc := gqlapi.ReturnCodeOK
	// 参数判断，只允许 first/before first/after last/before last/after 模式
	if first != nil && last != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", rc)
	}
	if after != nil && before != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", rc)
	}

	from := math.MaxInt64
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return nil, err
		}
		from = i
	}
	hasPreviousPage := true
	hasNextPage := true

	if from == math.MaxInt64 {
		hasPreviousPage = false
	}
	// 获取edges
	edges := []gqlapi.QueryEdge{}
	queries := []*models.Query{}
	var err error
	if err = g.Engine.Desc("user_id").Where("query_id < ?", from).Limit(*first).Find(&queries); err != nil {
		return nil, err
	}

	for _, query := range queries {
		edges = append(edges, gqlapi.QueryEdge{
			Node:   query,
			Cursor: EncodeCursor(fmt.Sprintf("%d", query.QueryID)),
		})
	}

	if len(edges) < *first {
		hasNextPage = false
	}

	if len(edges) == 0 {
		return nil, nil
	}
	// 获取pageInfo
	startCursor := EncodeCursor(fmt.Sprintf("%d", queries[0].QueryID))
	endCursor := EncodeCursor(fmt.Sprintf("%d", queries[len(queries)-1].QueryID))
	pageInfo := gqlapi.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     startCursor,
		EndCursor:       endCursor,
	}

	return &gqlapi.QueryConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}

// QuerySearch TODO: 下一个版本考虑实现
func (r *queryRootResolver) QuerySearch(ctx context.Context, search string, after *string, before *string, first *int, last *int) (*gqlapi.QueryConnection, error) {
	rc := gqlapi.ReturnCodeOK
	// 参数判断，只允许 first/before first/after last/before last/after 模式
	if first != nil && last != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", rc)
	}
	if after != nil && before != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", rc)
	}

	// TODO:
	panic("not implemented")
}

// CreateQuery 创建一个查询
func (r *mutationRootResolver) CreateQuery(ctx context.Context, input models.CreateQueryInput) (rs string, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user := credential.User

		if strings.TrimSpace(user.Name) == "" {
			rc = gqlapi.ReturnCodeRegistrationIncomplete
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)信息不完整。", rc, user.UUID)
			break
		}

		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.UUID == input.ClusterUUID {
				return true
			}
			return false
		})

		if cluster == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, input.ClusterUUID)
			break
		}

		if cluster.Status != gqlapi.ClusterStatusEnumMap["NORMAL"] {
			rc = gqlapi.ReturnCodeClusterNotAvailable
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不可用。", rc, input.ClusterUUID)
			break
		}

		// 检查群集关联
		if !caches.EdgesMap.Include(func(elem *models.Edge) bool {
			if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster] &&
				elem.AncestorID == user.UserID &&
				elem.DescendantID == cluster.ClusterID {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeForbidden
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)没有关联群集(uuid=%s)。", rc, user.UUID, cluster.UUID)
			break
		}

		passwd := func(c *models.Cluster) []byte {
			bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
			return bs
		}

		if _, err = cluster.Stat(input.Database, passwd); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		p := parser.New()
		stmts := []ast.StmtNode{}
		stmts, _, err = p.Parse(input.Content, "", "")
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		// 记录用户的查询
		query := &models.Query{
			Type:      gqlapi.QueryTypeEnumMap[gqlapi.QueryTypeEnumQuery],
			UserID:    user.UserID,
			ClusterID: cluster.ClusterID,
			Database:  input.Database,
			Content:   input.Content,
		}
		if _, err = g.Engine.Insert(query); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		// TODO: 进行规则检查
		// 1. 类型断言
		if len(stmts) != 1 {
			// TODO: 处理rc
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: 每次只允许执行一条查询语句。", rc)
			break
		}
		stmt := stmts[0]
		switch stmt.(type) {
		case *ast.SelectStmt, *ast.ShowStmt:
		case *ast.KillStmt, *ast.AlterUserStmt, *ast.CreateUserStmt, *ast.ExplainStmt, *ast.GrantStmt, *ast.SetPwdStmt, *ast.SetStmt, *ast.FlushStmt:
			isAdmin := func(cred tools.Credential) bool {
				for _, role := range cred.Roles {
					if role.RoleID == gqlapi.RoleEnumMap[gqlapi.RoleEnumAdmin] {
						return true
					}
				}
				return false
			}(credential)
			if !isAdmin {
				rc = gqlapi.ReturnCodeUnknowError
				err = fmt.Errorf("错误代码: %s, 错误信息: 权限不足，无法执行管理类指令，请联系管理员。", rc)
			}
		default:
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: 不支持的查询语句，请参考使用手册。", rc)
		}
		if err != nil {
			break
		}

		// TODO: 2. 对SelectStmt
		if _, ok := stmt.(*ast.SelectStmt); ok {

		}

		engine := &xorm.Engine{}

		if engine, err = cluster.Connect(input.Database, func(c *models.Cluster) []byte {
			bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
			return bs
		}); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		rows := []map[string]string{}
		if rows, err = engine.QueryString(input.Content); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		// 退出for循环
		bs := []byte{}
		if bs, err = json.Marshal(rows); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}
		rs = string(bs)

		events.Fire(events.EventQueryCreated, &events.QueryCreatedArgs{
			User:  *user,
			Query: *query,
		})

		break
	}

	return
}

// AnalyzeQuery 调用SOAR的SQL分析功能
func (r *mutationRootResolver) AnalyzeQuery(ctx context.Context, input models.SoarQueryInput) (report string, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user := credential.User
		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.UUID == input.ClusterUUID {
				return true
			}
			return false
		})
		if cluster == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, input.ClusterUUID)
			break L
		}

		if cluster.Status != gqlapi.ClusterStatusEnumMap["NORMAL"] {
			rc = gqlapi.ReturnCodeClusterNotAvailable
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不可用。", rc, input.ClusterUUID)
			break L
		}

		// 检查群集关联
		if !caches.EdgesMap.Include(func(elem *models.Edge) bool {
			if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster] &&
				elem.AncestorID == user.UserID &&
				elem.DescendantID == cluster.ClusterID {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeForbidden
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)没有关联群集(uuid=%s)。", rc, user.UUID, cluster.UUID)
			break L
		}

		passwd := func(c *models.Cluster) []byte {
			bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
			return bs
		}

		if _, err = cluster.Stat(input.Database, passwd); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 记录用户的查询
		query := &models.Query{
			Type:      gqlapi.QueryTypeEnumMap[gqlapi.QueryTypeEnumAnalyze],
			UserID:    user.UserID,
			ClusterID: cluster.ClusterID,
			Database:  input.Database,
			Content:   input.Content,
		}
		if _, err = g.Engine.Insert(query); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		args := []string{
			// TODO: 如果input.Content中有特殊字符，会不会存在服务器安全问题
			//       密码中如果有特殊字符会不会有服务器安全问题
			fmt.Sprintf("-query=%s", input.Content),
			fmt.Sprintf("-online-dsn=%s:%s@%s:%d/%s", cluster.User, passwd(cluster), cluster.IP, cluster.Port, input.Database),
			fmt.Sprintf("-test-dsn=%s:%s@%s:%d/%s", cluster.User, passwd(cluster), cluster.IP, cluster.Port, input.Database),
			"-explain-format=json",
			"-report-type=json",
			"-log-output=/tmp/soar.log",
			"-allow-drop-index",
			"-sampling",
			"-allow-online-as-test",
		}
		cmd := "soar"
		report, err = tools.TimeoutedExec(5*time.Second, cmd, args...)
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		events.Fire(events.EventQueryAnalyzed, &events.QueryAnalyzedArgs{
			User:  *user,
			Query: *query,
		})

		break L
	}

	return
}

// RewriteQuery 调用SOAR的重写SQL功能
func (r *mutationRootResolver) RewriteQuery(ctx context.Context, input models.SoarQueryInput) (sql string, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user := credential.User
		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.UUID == input.ClusterUUID {
				return true
			}
			return false
		})
		if cluster == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, input.ClusterUUID)
			break L
		}

		if cluster.Status != gqlapi.ClusterStatusEnumMap["NORMAL"] {
			rc = gqlapi.ReturnCodeClusterNotAvailable
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不可用。", rc, input.ClusterUUID)
			break L
		}

		// 检查群集关联
		if !caches.EdgesMap.Include(func(elem *models.Edge) bool {
			if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster] &&
				elem.AncestorID == user.UserID &&
				elem.DescendantID == cluster.ClusterID {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeForbidden
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)没有关联群集(uuid=%s)。", rc, user.UUID, cluster.UUID)
			break L
		}

		passwd := func(c *models.Cluster) []byte {
			bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
			return bs
		}

		if _, err = cluster.Stat(input.Database, passwd); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 记录用户的查询
		query := &models.Query{
			Type:      gqlapi.QueryTypeEnumMap[gqlapi.QueryTypeEnumRewrite],
			UserID:    user.UserID,
			ClusterID: cluster.ClusterID,
			Database:  input.Database,
			Content:   input.Content,
		}
		if _, err = g.Engine.Insert(query); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		args := []string{
			// TODO: 如果input.Content中有特殊字符，会不会存在服务器安全问题
			//       密码中如果有特殊字符会不会有服务器安全问题
			fmt.Sprintf("-query=%s", input.Content),
			fmt.Sprintf("-online-dsn=%s:%s@%s:%d/%s", cluster.User, passwd(cluster), cluster.IP, cluster.Port, input.Database),
			fmt.Sprintf("-test-dsn=%s:%s@%s:%d/%s", cluster.User, passwd(cluster), cluster.IP, cluster.Port, input.Database),
			// TODO: 为什么这个rules是硬编码
			fmt.Sprintf("-rewrite-rules=%s", "star2columns,delimiter"),
			"-report-type=rewrite",
			"-log-output=/tmp/soar.log",
			"-allow-online-as-test",
		}
		cmd := "soar"
		sql, err = tools.TimeoutedExec(5*time.Second, cmd, args...)
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		events.Fire(events.EventQueryRewrited, &events.QueryRewritedArgs{
			User:  *user,
			Query: *query,
		})

		break L
	}

	return
}

type queryResolver struct{ *Resolver }

// Cluster 查询发起的目标群集信息
func (r *queryResolver) Cluster(ctx context.Context, obj *models.Query) (cluster *models.Cluster, err error) {
	rc := gqlapi.ReturnCodeOK
	cluster = caches.ClustersMap.Any(func(elem *models.Cluster) bool {
		if elem.ClusterID == obj.ClusterID {
			return true
		}
		return false
	})
	if cluster == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 查询(uuid=%s)依赖的群集不存在。", rc, obj.UUID)
	}
	return
}

// User 查询发起的用户
func (r *queryResolver) User(ctx context.Context, obj *models.Query) (user *models.User, err error) {
	rc := gqlapi.ReturnCodeOK
	user = caches.UsersMap.Any(func(elem *models.User) bool {
		if elem.UserID == obj.UserID {
			return true
		}
		return false
	})

	if user == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 查询(uuid=%s)的发起人不存在。", rc, obj.UUID)
	}
	return
}
