package resolvers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

func TestResolveCronJobUsesFixedCommands(t *testing.T) {
	job, err := resolveCronJob(gqlapi.CronJobEnumMetricsRollup, nil)
	require.NoError(t, err)
	require.Equal(t, "halocli", job.cmd)
	require.Equal(t, []string{"metrics", "rollup"}, job.params)

	job, err = resolveCronJob(gqlapi.CronJobEnumMetricsCleanup, nil)
	require.NoError(t, err)
	require.Equal(t, "halocli", job.cmd)
	require.Equal(t, []string{"metrics", "cleanup", "--days=90"}, job.params)
}

func TestResolveCronJobAppliesRetentionDays(t *testing.T) {
	days := 30
	job, err := resolveCronJob(gqlapi.CronJobEnumMetricsCleanup, &days)
	require.NoError(t, err)
	require.Equal(t, []string{"metrics", "cleanup", "--days=30"}, job.params)
}

func TestResolveCronJobRejectsUnknownJob(t *testing.T) {
	_, err := resolveCronJob(gqlapi.CronJobEnum("SOMETHING"), nil)
	require.Error(t, err)
}
