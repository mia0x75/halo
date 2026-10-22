package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mia0x75/halo/internal/platform/config"
)

var retentionCmd = &cobra.Command{
	Use:   "retention",
	Short: "维护数据保留期",
	Long:  `Clean up the data types that grow without bound, using the retention windows configured in the platform options`,
}

var retentionCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "按保留期清理指标、命中、日志、站内信与任务记录",
	Long:  `Remove metric buckets, rule hits, audit logs, messages and task records older than the windows stored in mm_options`,
	Run:   retentionCleanup,
}

func init() {
	RootCmd.AddCommand(retentionCmd)
	retentionCmd.AddCommand(retentionCleanupCmd)
}

func retentionCleanup(cmd *cobra.Command, args []string) {
	if cfg, ok := os.LookupEnv("HALO_CFG"); ok {
		config.ParseConfig(cfg)
	} else {
		config.ParseConfig("")
	}

	application, appErr := newApp()
	if appErr != nil {
		fmt.Println(appErr)
		os.Exit(1)
	}

	result, err := application.Services.Retention.Cleanup(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("已清理 指标桶 %d、规则命中 %d、操作日志 %d、站内信 %d、任务记录 %d\n",
		result.Metrics, result.Hits, result.Logs, result.Messages, result.Tasks)
}
