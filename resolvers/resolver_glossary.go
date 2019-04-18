package resolvers

import (
	"context"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// Glossaries 根据分组信息获取字典信息
func (r *queryRootResolver) Glossaries(ctx context.Context, groups []string) (L []*models.Glossary, err error) {
	L = caches.GlossariesMap.Filter(func(elem *models.Glossary) bool {
		return tools.Contains(groups, elem.Group)
	})
	return
}

type glossaryResolver struct{ *Resolver }
