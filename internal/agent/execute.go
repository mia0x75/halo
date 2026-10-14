package agent

import (
	"context"
	"time"
)

// Execute 在目标库上执行单条语句。连接层错误直接返回；语句本身执行
// 失败时结果里带 Error，供上层落审计报告。
func (t *Target) Execute(ctx context.Context, query string) (*ExecResult, error) {
	start := time.Now()
	result := &ExecResult{}

	res, err := t.db.ExecContext(ctx, query)
	result.Duration = time.Since(start)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		result.Error = err.Error()
		return result, err
	}
	result.RowsAffected = affected

	return result, nil
}
