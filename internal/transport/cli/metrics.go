package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/platform/config"
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "维护指标表",
	Long:  `This subcommand ensures the daily metric buckets exist and reports the counter totals`,
	Run:   metrics,
}

func init() {
	RootCmd.AddCommand(metricsCmd)
}

func metrics(cmd *cobra.Command, args []string) {
	var err error
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

	ctx := context.Background()
	today := time.Now().Format("2006-01-02")

	// Daily buckets must exist before the first event of the day bumps them,
	// otherwise the trend query has no row to start from.
	ensure := func(name string) error {
		key := domain.MetricKey{Name: name, Scope: domain.MetricScopeDaily, Bucket: today}
		if _, gerr := application.Services.Metric.GetByKey(ctx, key); gerr == nil {
			return nil
		} else if gerr != domain.ErrNotFound {
			return gerr
		}
		return application.Services.Metric.Create(ctx, &domain.Metric{
			Name:     name,
			Scope:    domain.MetricScopeDaily,
			Bucket:   today,
			Type:     domain.MetricTypeCounter,
			UUID:     uuid.NewString(),
			Version:  1,
			CreateAt: uint(time.Now().Unix()),
		}, 0)
	}

	for _, name := range []string{"total-tickets", "total-queries"} {
		if err = ensure(name); err != nil {
			break
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	metrics, err := application.Services.Metric.List(ctx, "", 100)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, metric := range metrics {
		fmt.Printf("%-16s scope=%d bucket=%-10s value=%f\n", metric.Name, metric.Scope, metric.Bucket, metric.Value)
	}
}
