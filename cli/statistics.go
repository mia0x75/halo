package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/internal/domain"
)

// statisticsCmd represents the statistics command
var statisticsCmd = &cobra.Command{
	Use:   "statistics",
	Short: "Add records into statistics table",
	Long:  `This subcommand will check the entries of statistics table for overall tickets and queries daily`,
	Run:   statistics,
}

func init() {
	RootCmd.AddCommand(statisticsCmd)
}

func statistics(cmd *cobra.Command, args []string) {
	var err error
	var buf bytes.Buffer

	today := time.Now().Format("2006-01-02")

	for {
		if cfg, ok := os.LookupEnv("HALO_CFG"); !ok {
			err = fmt.Errorf("Missing halo config")
			break
		} else {
			buf.WriteString(cfg)
			g.ParseConfig(cfg)
		}

		application, appErr := newApp()
		if appErr != nil {
			err = appErr
			break
		}
		ctx := context.Background()

		ensure := func(group string) error {
			if _, gerr := application.Services.Statistic.GetByKey(ctx, group, today); gerr == nil {
				return nil
			} else if !errors.Is(gerr, domain.ErrNotFound) {
				return gerr
			}
			return application.Services.Statistic.Create(ctx, &domain.Statistic{
				Group:    group,
				Key:      today,
				Value:    0,
				UUID:     uuid.NewString(),
				Version:  1,
				CreateAt: uint(time.Now().Unix()),
			}, 0)
		}
		if err = ensure("tickets-daily"); err != nil {
			break
		}
		if err = ensure("queries-daily"); err != nil {
			break
		}

		break
	}

	buf.WriteString("\n")
	if err != nil {
		buf.WriteString(err.Error())
	} else {
		buf.WriteString("statistics success")
	}
	name := fmt.Sprintf("/tmp/halocli-statistics-%s.txt", today)
	if ioutil.WriteFile(name, buf.Bytes(), 0644) == nil {
	}
}
