package models

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func TestModels(t *testing.T) {
	engine, err := xorm.NewEngine("mysql", "root:Admin?!##12@tcp(127.0.0.1:3306)/mcms?charset=utf8&loc=Local&parseTime=true")
	if err != nil {
		t.Error(err)
	}
	cluster := Cluster{
		ClusterID: 2,
	}
	_, err = engine.Get(&cluster)
	cluster.Status = 2
	affected, err := engine.Update(&cluster)
	if affected != 1 {
		t.Log("数据被其他人改写或者删除。")
	}
	engine.ShowSQL(true)
	if err != nil {
		t.Error(err)
	}
}
