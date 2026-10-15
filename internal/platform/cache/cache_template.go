package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/models"
)

// TemplatesMap 线程安全的数据缓存对象，TemplatesMap 规则缓存对象
var TemplatesMap = SafeMap[models.Template]{}

// Init 缓存初始化
func TemplatesInit() (err error) {
	_, err = TemplatesMap.WithInit(func() ([]*models.Template, error) {
		return loadUUIDKeyed(repos.Template.List, bridge.TemplateFromDomain)
	})
	return
}
