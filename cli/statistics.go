package cli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
	"github.com/spf13/cobra"
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
			g.InitDB()
		}

		tds := &models.Statistic{}
		qds := &models.Statistic{}
		if _, err = g.Engine.Where("`group` = ? AND `key` = ?", "tickets-daily", today).Get(tds); err != nil {
			break
		}
		if tds.UUID == "" {
			tds.Group = "tickets-daily"
			tds.Key = d
			tds.Value = 0
			if _, err = g.Engine.Insert(tds); err != nil {
				break
			}
		}
		if _, err = g.Engine.Where("`group` = ? AND `key` = ?", "queries-daily", today).Get(qds); err != nil {
			break
		}
		if qds.UUID == "" {
			qds.Group = "queries-daily"
			qds.Key = d
			qds.Value = 0
			if _, err = g.Engine.Insert(qds); err != nil {
				break
			}
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
