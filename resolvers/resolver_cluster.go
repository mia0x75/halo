package resolvers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fatih/structs"
	"github.com/mia0x75/yql"
	"xorm.io/core"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/events"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// TestCluster 测试群集的连接性
func (r queryRootResolver) TestCluster(ctx context.Context, input *models.ValidateConnectionInput) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		const pattern = "%s:%s@tcp(%s:%d)/mysql?loc=Local&parseTime=true"
		db := &sql.DB{}
		addr := fmt.Sprintf(pattern, input.User, input.Password, input.IP, input.Port)

		if db, err = sql.Open("mysql", addr); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: 连接参数不正确: %s", rc, err.Error())
			break
		}
		if err = db.Ping(); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: 无法连接到目标群集: %s", rc, err.Error())
			break
		}
		db.Close()

		// 退出for循环
		ok = true
		break
	}

	return
}

// Databases 获取某一个群集的所有用户数据库
func (r queryRootResolver) Databases(ctx context.Context, clusterUUID string) (L []*gqlapi.Database, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user := credential.User
		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.UUID == clusterUUID {
				return true
			}
			return false
		})
		if cluster == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, clusterUUID)
			break
		}

		if cluster.Status != gqlapi.ClusterStatusEnumMap["NORMAL"] {
			rc = gqlapi.ReturnCodeClusterNotAvailable
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不可用。", rc, cluster.UUID)
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

		databases := []models.Database{}
		if databases, err = cluster.Databases(func(c *models.Cluster) []byte {
			bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
			return bs
		}); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		for _, d := range databases {
			L = append(L, &gqlapi.Database{
				Name:    d.Name,
				Charset: d.Charset,
				Collate: d.Collate,
			})
		}

		break
	}

	return
}

// Cluster 获取某一指定群集的详细信息
func (r queryRootResolver) Cluster(ctx context.Context, id string) (cluster *models.Cluster, err error) {
	rc := gqlapi.ReturnCodeOK
	cluster = caches.ClustersMap.Any(func(elem *models.Cluster) bool {
		if elem.UUID == id {
			return true
		}
		return false
	})
	if cluster == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, id)
	}
	return
}

// Clusters 分页获取群集列表，TODO: 分页未完成
func (r queryRootResolver) Clusters(ctx context.Context, after *string, before *string, first *int, last *int) (*gqlapi.ClusterConnection, error) {
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

	clusters := caches.ClustersMap.All()
	edges := []*gqlapi.ClusterEdge{}
	for _, cluster := range clusters {
		edges = append(edges, &gqlapi.ClusterEdge{
			Node:   cluster,
			Cursor: EncodeCursor(fmt.Sprintf("%d", cluster.ClusterID)),
		})
	}
	if len(edges) == 0 {
		return nil, nil
	}
	// 获取pageInfo
	startCursor := EncodeCursor(fmt.Sprintf("%d", clusters[0].ClusterID))
	endCursor := EncodeCursor(fmt.Sprintf("%d", clusters[len(clusters)-1].ClusterID))
	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: false,
		HasNextPage:     false,
		StartCursor:     startCursor,
		EndCursor:       endCursor,
	}

	return &gqlapi.ClusterConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}

// ClusterSearch 群集查询，TODO: 分页未实现
func (r queryRootResolver) ClusterSearch(ctx context.Context, search string, after *string, before *string, first *int, last *int) (rs *gqlapi.ClusterConnection, err error) {
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

		var rule yql.Ruler
		if rule, err = yql.Rule(search); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: 语法错误，请参考帮助文档。", rc)
			break L
		}
		clusters := caches.ClustersMap.Filter(func(elem *models.Cluster) bool {
			data := structs.Map(elem)
			if matched, _ := rule.Match(data); matched {
				return true
			}
			return false
		})
		count := len(clusters)
		if count == 0 {
			break L
		}
		// 获取edges
		edges := []*gqlapi.ClusterEdge{}
		for _, cluster := range clusters {
			edges = append(edges, &gqlapi.ClusterEdge{
				Node:   cluster,
				Cursor: EncodeCursor(strconv.Itoa(0)),
			})
		}
		pageInfo := &gqlapi.PageInfo{
			StartCursor:     edges[0].Cursor,
			EndCursor:       edges[count-1].Cursor,
			HasPreviousPage: false,
			HasNextPage:     false,
		}

		rs = &gqlapi.ClusterConnection{
			PageInfo:   pageInfo,
			Edges:      edges,
			TotalCount: count,
		}

		break L
	}

	return
}

// Metadata 获取群集上某一个数据库的元数据信息
func (r queryRootResolver) Metadata(ctx context.Context, clusterUUID string, database string) (resp string, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user := credential.User
		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.UUID == clusterUUID {
				return true
			}
			return false
		})

		if cluster == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, clusterUUID)
			break L
		}

		if cluster.Status != gqlapi.ClusterStatusEnumMap["NORMAL"] {
			rc = gqlapi.ReturnCodeClusterNotAvailable
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不可用。", rc, clusterUUID)
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

		if _, err = cluster.Stat(database, passwd); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		tables := map[string][]*core.Table{}
		if tables, err = cluster.Metadata(database, passwd); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		var bs []byte
		bs, err = json.Marshal(cluster.Repack(tables[database]))
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		resp = string(bs)

		break L
	}

	return
}

// CreateCluster 创建一个群集
func (r mutationRootResolver) CreateCluster(ctx context.Context, input models.CreateClusterInput) (cluster *models.Cluster, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		// 检查Alias唯一性
		if caches.ClustersMap.Include(func(elem *models.Cluster) bool {
			if elem.Alias == input.Alias {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeDuplicateAlias
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(alias=%s)已经存在。", rc, input.Alias)
			break L
		}
		// 主机+端口唯一性,IP+端口唯一性
		if caches.ClustersMap.Include(func(elem *models.Cluster) bool {
			if (elem.Host == input.Host || elem.IP == input.IP) && elem.Port == input.Port {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeDuplicateHost
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(cluster=%s:%d 或 cluster=%s:%d)已经存在。", rc, input.Host, input.Port, input.IP, input.Port)
			break L
		}
		var passwd []byte
		passwd, err = tools.EncryptAES([]byte(input.Password), g.Config().Secret.Crypto)
		if err != nil {
			break L
		}
		cluster = &models.Cluster{}
		cluster.Host = input.Host
		cluster.Port = input.Port
		cluster.IP = input.IP
		cluster.User = input.User
		cluster.Status = input.Status
		cluster.Password = passwd
		cluster.Alias = input.Alias

		if _, err = g.Engine.Insert(cluster); err != nil {
			cluster = nil
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 主动同步缓存
		caches.ClustersMap.Append(cluster)

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventClusterCreated, &events.ClusterCreatedArgs{
			Manager: *credential.User,
			Cluster: *cluster,
		})

		// 退出for循环
		break L
	}

	return
}

// UpdateCluster 更新一个群集
func (r mutationRootResolver) UpdateCluster(ctx context.Context, input models.UpdateClusterInput) (cluster *models.Cluster, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		cluster = caches.ClustersMap.Any(func(elem *models.Cluster) bool {
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
		// 检查Alias唯一性
		exists := func(clusterID uint, alias string) bool {
			c := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
				if elem.Alias == alias {
					return true
				}
				return false
			})
			if c != nil && c.ClusterID != clusterID {
				return true
			}
			return false
		}(cluster.ClusterID, input.Alias)
		if exists {
			rc = gqlapi.ReturnCodeDuplicateAlias
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(alias=%s)已经存在。", rc, input.Alias)
			break L
		}

		// 避免修改同一实例未更换IP Port Host报错
		if cluster.Host != input.Host || cluster.IP != input.IP || cluster.Port != input.Port {
			// 主机+端口唯一性
			// IP+端口唯一性
			if caches.ClustersMap.Include(func(elem *models.Cluster) bool {
				if (elem.Host == input.Host || elem.IP == input.IP) && elem.Port == input.Port {
					return true
				}
				return false
			}) {
				rc = gqlapi.ReturnCodeDuplicateHost
				err = fmt.Errorf("错误代码: %s, 错误信息: 群集(cluster=%s:%d 或 cluster=%s:%d)已经存在。", rc, input.Host, input.Port, input.IP, input.Port)
				break L
			}
		}

		var passwd []byte

		passwd, err = tools.EncryptAES([]byte(input.Password), g.Config().Secret.Crypto)
		if err != nil {
			break L
		}

		cluster.Host = input.Host
		cluster.Password = passwd
		cluster.IP = input.IP
		cluster.Port = input.Port
		cluster.User = input.User
		cluster.Alias = input.Alias
		cluster.Status = input.Status

		if _, err = g.Engine.ID(cluster.ClusterID).AllCols().Update(cluster); err != nil {
			cluster = nil
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventClusterUpdated, &events.ClusterUpdatedArgs{
			Manager: *credential.User,
			Cluster: *cluster,
		})

		// 退出for循环
		break L
	}

	return
}

// RemoveCluster 有限的条件下，删除群集
func (r mutationRootResolver) RemoveCluster(ctx context.Context, id string) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.UUID == id {
				return true
			}
			return false
		})
		if cluster == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, id)
			break
		}
		if _, err = g.Engine.ID(cluster.ClusterID).Delete(cluster); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		// 主动刷新缓存
		caches.ClustersMap.Remove(func(elem *models.Cluster) bool {
			if elem.UUID == id {
				return true
			}
			return false
		})

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventClusterRemoved, &events.ClusterRemovedArgs{
			Manager: *credential.User,
			Cluster: *cluster,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// PatchClusterStatus 修改群集状态
func (r mutationRootResolver) PatchClusterStatus(ctx context.Context, input models.PatchClusterStatusInput) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		if input.Status < 1 || input.Status > 2 {
			rc = gqlapi.ReturnCodeInvalidParams
			err = fmt.Errorf("错误代码: %s, 错误信息: 参数(status=%d)无效。", rc, input.Status)
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

		cluster.Status = input.Status
		// TODO: 需要单元测试
		if _, err = g.Engine.ID(cluster.ClusterID).Where("`status` <> ?", input.Status).Update(cluster); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventClusterStatusPatched, &events.ClusterStatusPatchedArgs{
			Manager: *credential.User,
			Cluster: *cluster,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

type clusterResolver struct{ *Resolver }
