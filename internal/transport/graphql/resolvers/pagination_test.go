package resolvers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

func TestProbeLimitAddsOneRow(t *testing.T) {
	require.Equal(t, 11, probeLimit(10))
	// A non-positive page size must stay as-is so the store keeps its own
	// default instead of receiving a meaningless limit.
	require.Equal(t, 0, probeLimit(0))
	require.Equal(t, -1, probeLimit(-1))
}

func TestTrimKeysetPageDropsProbeRow(t *testing.T) {
	rows := []int{1, 2, 3, 4}

	page, hasNext := trimKeysetPage(rows, 3)
	require.Equal(t, []int{1, 2, 3}, page)
	require.True(t, hasNext)
}

func TestTrimKeysetPageFullPageIsLastPage(t *testing.T) {
	// A page that exactly fills the limit is the last page. The old
	// `len(edges) == limit` check reported a phantom next page here.
	rows := []int{1, 2, 3}

	page, hasNext := trimKeysetPage(rows, 3)
	require.Equal(t, []int{1, 2, 3}, page)
	require.False(t, hasNext)
}

func TestTrimKeysetPageEmptyPage(t *testing.T) {
	page, hasNext := trimKeysetPage([]int{}, 3)
	require.Empty(t, page)
	require.False(t, hasNext)
}

func TestTrimKeysetPageKeepsAllRowsWithoutLimit(t *testing.T) {
	rows := []int{1, 2, 3}

	page, hasNext := trimKeysetPage(rows, 0)
	require.Equal(t, rows, page)
	require.False(t, hasNext)
}

func TestTicketConnectionReportsTotalAcrossPages(t *testing.T) {
	tickets := []*domain.Ticket{
		{TicketID: 5},
		{TicketID: 4},
		{TicketID: 3},
		{TicketID: 2},
		{TicketID: 1},
	}
	first := 2

	conn, err := ticketConnection(tickets, nil, nil, &first, nil)
	require.NoError(t, err)
	require.Len(t, conn.Edges, 2)
	require.True(t, conn.PageInfo.HasNextPage)
	// The badge shows the real amount, not how many rows this page carries.
	require.Equal(t, 5, conn.TotalCount)
}

func TestTicketConnectionLastPageHasNoNext(t *testing.T) {
	tickets := []*domain.Ticket{
		{TicketID: 2},
		{TicketID: 1},
	}
	first := 2

	conn, err := ticketConnection(tickets, nil, nil, &first, nil)
	require.NoError(t, err)
	require.Len(t, conn.Edges, 2)
	require.False(t, conn.PageInfo.HasNextPage)
	require.Equal(t, 2, conn.TotalCount)
}

func TestPageInMemoryResumesAfterCursor(t *testing.T) {
	ids := []int{5, 4, 3, 2, 1}
	cursor := EncodeCursor("4")
	first := 2

	page, hasNext, hasPrevious, err := pageInMemory(ids, func(id int) int { return id }, &cursor, nil, &first, nil)
	require.NoError(t, err)
	require.Equal(t, []int{3, 2}, page)
	require.True(t, hasNext)
	require.True(t, hasPrevious)
}

func TestPageInMemoryPagingAscendingRowsDoesNotRepeat(t *testing.T) {
	// Stores return rows oldest-first; paging must still walk newer to older.
	ids := []int{1, 2, 3, 4, 5}
	first := 2

	page1, hasNext, _, err := pageInMemory(ids, func(id int) int { return id }, nil, nil, &first, nil)
	require.NoError(t, err)
	require.Equal(t, []int{5, 4}, page1)
	require.True(t, hasNext)

	cursor := EncodeCursor("4")
	page2, hasNext, _, err := pageInMemory(ids, func(id int) int { return id }, &cursor, nil, &first, nil)
	require.NoError(t, err)
	require.Equal(t, []int{3, 2}, page2)
	require.True(t, hasNext)

	cursor = EncodeCursor("2")
	page3, hasNext, _, err := pageInMemory(ids, func(id int) int { return id }, &cursor, nil, &first, nil)
	require.NoError(t, err)
	require.Equal(t, []int{1}, page3)
	require.False(t, hasNext)
}

func TestPageInMemoryLastPageAfterCursor(t *testing.T) {
	ids := []int{5, 4, 3, 2, 1}
	cursor := EncodeCursor("2")
	first := 2

	page, hasNext, _, err := pageInMemory(ids, func(id int) int { return id }, &cursor, nil, &first, nil)
	require.NoError(t, err)
	require.Equal(t, []int{1}, page)
	require.False(t, hasNext)
}

func TestPageInMemoryRejectsConflictingArgs(t *testing.T) {
	first, last := 1, 1
	_, _, _, err := pageInMemory([]int{1}, func(id int) int { return id }, nil, nil, &first, &last)
	require.Error(t, err)
}

func TestPageInMemoryRejectsInvalidCursor(t *testing.T) {
	bad := "not-a-cursor"
	_, _, _, err := pageInMemory([]int{1}, func(id int) int { return id }, &bad, nil, nil, nil)
	require.Error(t, err)
}

func TestDeletedUserPlaceholderKeepsRowReadable(t *testing.T) {
	// Non-null schema fields must never come back as an error: that nulls the
	// whole edge and blanks the list the user is paging through.
	user := deletedUserPlaceholder(42)
	require.NotNil(t, user)
	require.Equal(t, uint(42), user.UserID)
	require.NotEmpty(t, user.UUID)
	require.NotEmpty(t, user.Email)
}

func TestDeletedClusterPlaceholderKeepsRowReadable(t *testing.T) {
	cluster := deletedClusterPlaceholder(7)
	require.NotNil(t, cluster)
	require.Equal(t, uint(7), cluster.ClusterID)
	require.NotEmpty(t, cluster.UUID)
}

func TestLogUserResolverFallsBackForDeletedAccount(t *testing.T) {
	resolver := &logResolver{}
	user, err := resolver.User(context.Background(), &models.Log{UserID: 999999})
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, uint(999999), user.UserID)
}

func TestTicketResolversFallBackForDeletedRefs(t *testing.T) {
	resolver := &ticketResolver{}
	ticket := &models.Ticket{UserID: 999998, ReviewerID: 999999, ClusterID: 999997}

	user, err := resolver.User(context.Background(), ticket)
	require.NoError(t, err)
	require.NotNil(t, user)

	reviewer, err := resolver.Reviewer(context.Background(), ticket)
	require.NoError(t, err)
	require.NotNil(t, reviewer)

	cluster, err := resolver.Cluster(context.Background(), ticket)
	require.NoError(t, err)
	require.NotNil(t, cluster)
}
