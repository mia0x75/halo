package domain

import "context"

// Group is a cluster group: the servers of one group belong together. Users
// are grouped by roles, servers are grouped by groups.
type Group struct {
	GroupID     uint   `db:"group_id"`
	UUID        string `db:"uuid"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Status      uint8  `db:"status"`
	Version     int    `db:"version"`
	UpdateAt    *uint  `db:"update_at"`
	CreateAt    uint   `db:"create_at"`
}

// GroupStatus is the state of a group.
const (
	GroupStatusNormal  uint8 = 1
	GroupStatusBlocked uint8 = 2
)

// GroupRepository stores the cluster groups.
type GroupRepository interface {
	Create(ctx context.Context, g *Group) error
	GetByUUID(ctx context.Context, uuid string) (*Group, error)
	GetByName(ctx context.Context, name string) (*Group, error)
	Update(ctx context.Context, g *Group) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor string, limit int) ([]*Group, error)
}
