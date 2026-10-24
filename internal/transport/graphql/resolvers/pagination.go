package resolvers

import (
	"encoding/base64"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// probeLimit returns the row count to request from the store for a page of
// `first` rows. One extra row is fetched so the caller can tell whether a next
// page actually exists instead of guessing from a full page.
func probeLimit(first int) int {
	if first <= 0 {
		return first
	}
	return first + 1
}

// trimKeysetPage drops the probe row appended by a limit+1 fetch and reports
// whether it found one. hasNext is true only when the store returned more rows
// than the page size, i.e. at least one row remains past this page. A page that
// merely fills the limit is the last page, which is the case the old
// `len(edges) == limit` check got wrong.
func trimKeysetPage[T any](rows []T, limit int) (page []T, hasNextPage bool) {
	if limit > 0 && len(rows) > limit {
		return rows[:limit], true
	}
	return rows, false
}

// pageInMemory slices rows already loaded whole (cache hits, or a store call
// that returns the full list) into a cursor page keyed by id. Cursors are
// encoded ids and paging resumes below the cursor, so rows are reordered
// newest-first: stores return them oldest-first, which would otherwise hand
// back the same rows on every page.
func pageInMemory[T any](
	rows []T,
	idOf func(T) int,
	after *string,
	before *string,
	first *int,
	last *int,
) (page []T, hasNextPage bool, hasPreviousPage bool, err error) {
	if first != nil && last != nil {
		return nil, false, false, fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", gqlapi.ReturnCodeInvalidParams)
	}
	if after != nil && before != nil {
		return nil, false, false, fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", gqlapi.ReturnCodeInvalidParams)
	}

	from := math.MaxInt
	if after != nil && *after != "" {
		id, decodeErr := decodeIntCursor(*after)
		if decodeErr != nil {
			return nil, false, false, decodeErr
		}
		from = id
	}

	sorted := make([]T, len(rows))
	copy(sorted, rows)
	sort.SliceStable(sorted, func(i, j int) bool { return idOf(sorted[i]) > idOf(sorted[j]) })

	page = make([]T, 0, len(sorted))
	for _, row := range sorted {
		if idOf(row) >= from {
			continue
		}
		page = append(page, row)
	}

	limit := len(page)
	if first != nil && *first < limit {
		limit = *first
	}
	if last != nil && *last < limit {
		limit = *last
	}
	page, hasNextPage = trimKeysetPage(page, limit)

	return page, hasNextPage, after != nil && *after != "", nil
}

// decodeIntCursor decodes an id cursor produced by EncodeCursor.
func decodeIntCursor(cursor string) (int, error) {
	raw, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("错误代码: %s, 错误信息: 游标(cursor=%s)无效。", gqlapi.ReturnCodeInvalidParams, cursor)
	}
	id, err := strconv.Atoi(strings.TrimPrefix(string(raw), "cursor"))
	if err != nil {
		return 0, fmt.Errorf("错误代码: %s, 错误信息: 游标(cursor=%s)无效。", gqlapi.ReturnCodeInvalidParams, cursor)
	}
	return id, nil
}

// deletedUserPlaceholder stands in for a user reference that no longer exists.
// Audit rows and tickets must stay readable after the account is removed, and
// the schema marks those fields non-null, so returning an error would null out
// the whole row and blank the list that contains it.
//
// A zero id means the row predates any account (wizard seeding in old builds),
// so it is reported as an unowned entry rather than a deleted user.
func deletedUserPlaceholder(userID uint) *models.User {
	name := "已删除用户"
	if userID == 0 {
		name = "系统"
	}
	return &models.User{
		UserID: userID,
		UUID:   fmt.Sprintf("deleted-user-%d", userID),
		Email:  "-",
		Name:   name,
		Status: gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumBlocked],
	}
}

// deletedClusterPlaceholder mirrors deletedUserPlaceholder for removed
// clusters: the ticket/query row stays visible instead of vanishing.
func deletedClusterPlaceholder(clusterID uint) *models.Cluster {
	return &models.Cluster{
		ClusterID: clusterID,
		UUID:      fmt.Sprintf("deleted-cluster-%d", clusterID),
		Alias:     "已删除群集",
	}
}
