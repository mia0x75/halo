package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"

	"github.com/go-xorm/core"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// 接口确认，如果有报错表示没有实现接口的全部方法
var (
	_ Validator = &DatabaseCreateVldr{}
	_ Validator = &DatabaseAlterVldr{}
	_ Validator = &DatabaseDropVldr{}
	_ Validator = &TableCreateVldr{}
	_ Validator = &TableAlterVldr{}
	_ Validator = &TableRenameVldr{}
	_ Validator = &TableDropVldr{}
	_ Validator = &IndexCreateVldr{}
	_ Validator = &IndexDropVldr{}
	_ Validator = &InsertVldr{}
	_ Validator = &UpdateVldr{}
	_ Validator = &DeleteVldr{}
	_ Validator = &ViewCreateVldr{}
	_ Validator = &ViewCreateVldr{}
	_ Validator = &ViewCreateVldr{}
	_ Validator = &FuncCreateVldr{}
	_ Validator = &FuncAlterVldr{}
	_ Validator = &FuncDropVldr{}
	_ Validator = &TriggerCreateVldr{}
	_ Validator = &TriggerAlterVldr{}
	_ Validator = &TriggerDropVldr{}
	_ Validator = &EventCreateVldr{}
	_ Validator = &EventAlterVldr{}
	_ Validator = &EventDropVldr{}
	_ Validator = &ProcCreateVldr{}
	_ Validator = &ProcCreateVldr{}
	_ Validator = &ProcCreateVldr{}
	_ Validator = &MiscVldr{}
)

// Context 在不同的验证组中共享内容
type Context struct {
}

// 相当于注册表
var registry = map[uint16]Validator{
	100: &DatabaseCreateVldr{},
	101: &DatabaseAlterVldr{},
	102: &DatabaseDropVldr{},
	110: &TableCreateVldr{},
	120: &TableAlterVldr{},
	130: &TableRenameVldr{},
	140: &TableDropVldr{},
	150: &IndexCreateVldr{},
	151: &IndexDropVldr{},
	160: &InsertVldr{},
	161: &ReplaceVldr{},
	162: &UpdateVldr{},
	163: &DeleteVldr{},
	170: &SelectVldr{},
	180: &ViewCreateVldr{},
	181: &ViewCreateVldr{},
	182: &ViewCreateVldr{},
	190: &FuncCreateVldr{},
	191: &FuncAlterVldr{},
	192: &FuncDropVldr{},
	200: &TriggerCreateVldr{},
	201: &TriggerAlterVldr{},
	202: &TriggerDropVldr{},
	210: &EventCreateVldr{},
	211: &EventAlterVldr{},
	212: &EventDropVldr{},
	220: &ProcCreateVldr{},
	221: &ProcCreateVldr{},
	222: &ProcCreateVldr{},
	230: &MiscVldr{},
}

// Run 调用入口
func Run(stmts []*models.Statement, cluster *models.Cluster, ticket *models.Ticket) {
	wg := &sync.WaitGroup{}
	for g, v := range registry {
		if v.Enabled() {
			wg.Add(1)
			v.SetGroup(g) //
			v.Preprocess(cluster, ticket, stmts)
			go v.Validate(wg)
		}
	}
	wg.Wait()
}

// Validator 接口定义
type Validator interface {
	SetGroup(gid uint16)                                                                  // 设定分组，获得规则列表
	Preprocess(cluster *models.Cluster, ticket *models.Ticket, stmts []*models.Statement) // 预处理
	Enabled() bool                                                                        // 分组规则是否启用
	SetContext(ctx Context)                                                               // 设定上下文
	Validate(wg *sync.WaitGroup)                                                          // 验证
}

type vldr struct {
	rules   []*models.Rule      // 全部适用的规则
	stmts   []*models.Statement // 全部等待审核的数据
	ticket  *models.Ticket      // 等待审核的工单
	cluster *models.Cluster     // 目标群集

	databases []models.Database // 目标群集所有的用户数据库
	tables    []*core.Table     // 目标群集上目标数据库的元数据
	ctx       Context           // 上下文，保存检测报告？
}

// SetGroup 设置验证分组，并初始化验证规则
func (v *vldr) SetGroup(gid uint16) {
	rules := caches.RulesMap.Filter(func(elem *models.Rule) bool {
		if elem.VldrGroup == gid {
			return true
		}
		return true
	})
	v.rules = rules
}

// Preprocess 验证预处理
func (v *vldr) Preprocess(cluster *models.Cluster, ticket *models.Ticket, stmts []*models.Statement) {
	v.stmts = stmts
	v.ticket = ticket
	v.cluster = cluster

	v.tables, _ = cluster.Metadata(ticket.Database, func(c *models.Cluster) []byte {
		bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
		return bs
	})
	v.databases, _ = cluster.Databases(func(c *models.Cluster) []byte {
		bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
		return bs
	})
}

// Call 方法反射
func Call(object interface{}, method string, params ...interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("错误代码: 1500, 错误信息: %s\n", err)
			return
		}
	}()

	v := reflect.ValueOf(object)
	// Check if the passed interface is a pointer
	if v.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface, so we have a pointer to work with
		v = reflect.New(reflect.TypeOf(object))
	}
	// Get the method by name
	f := v.MethodByName(method)
	if !f.IsValid() {
		return fmt.Errorf("Couldn't find method `%s` in interface `%s`, is it exported?", method, v.Type())
	}
	// 'dereference' with Elem() and get the field by name
	// f := ValueIface.Elem().FieldByName(FieldName)
	// if !f.IsValid() {
	// 	return fmt.Errorf("Interface `%s` does not have the field `%s`", v.Type(), f)
	// }
	// 判断函数参数和传入的参数是否相等
	if len(params) != f.Type().NumIn() {
		return fmt.Errorf("")
	}
	if f.Type().NumIn() == 0 {
		f.Call(nil)
	} else {
		// 然后将传入参数转为反射类型切片
		in := make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
		// 利用函数反射对象的call方法调用函数.
		f.Call(in)
	}
	return nil
}

// Match 正则匹配
func Match(r *models.Rule, src string, params ...interface{}) (err error) {
	match := false
	match, err = regexp.MatchString(r.Values, src)
	if err != nil {
		err = fmt.Errorf("`%s`不是一个有的效正则表达式。", r.Values)
		return
	}

	if match {
		return nil
	}
	err = fmt.Errorf(r.Message, params...)
	return
}
