package resolvers

import (
	"context"

	"github.com/mia0x75/halo/internal/platform/cache"
	"github.com/mia0x75/halo/models"
)

// Avatars 全部的头像，头像数量不多，这个不进行分页处理
func (r *queryRootResolver) Avatars(ctx context.Context) (L []*models.Avatar, err error) {
	L = caches.AvatarsMap.All()
	return
}
