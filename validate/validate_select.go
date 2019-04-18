package validate

import (
	"sync"

	"github.com/mia0x75/parser/ast"
)

// SelectVldr 查询语句审核
type SelectVldr struct {
	vldr

	ds *ast.SelectStmt
}

// Call 利用反射方法动态调用审核函数
func (v *SelectVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *SelectVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *SelectVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *SelectVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if ds, ok := node.(*ast.SelectStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.ds = ds
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}
