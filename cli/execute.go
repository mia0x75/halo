package cli

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/mia0x75/halo/events"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
	"github.com/spf13/cobra"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute ticket",
	Long:  `This subcommand execute the specify ticket`,
	Run:   execute,
}

func init() {
	RootCmd.AddCommand(executeCmd)
	executeCmd.Flags().StringP("ticket", "T", "", "specify the ticket need to be executed")
}

func execute(cmd *cobra.Command, args []string) {
	var err error
	var buf bytes.Buffer
	var ticketUUID string
L:
	for {
		if cfg, ok := os.LookupEnv("HALO_CFG"); !ok {
			err = fmt.Errorf("Missing halo config")
			break
		} else {
			buf.WriteString(cfg)
			g.ParseConfig(cfg)
			g.InitDB()
		}

		if ticketUUID, err = cmd.Flags().GetString("ticket"); err != nil {
			break
		}
		if !isValidUUID(ticketUUID) {
			err = fmt.Errorf("String %s is not a valid UUID", ticketUUID)
			break
		}

		buf.WriteString("\n")
		buf.WriteString(ticketUUID)

		passwd := func(c *models.Cluster) []byte {
			bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
			return bs
		}

		ticket := &models.Ticket{}
		if _, err = g.Engine.Where("`uuid` = ?", ticketUUID).Get(ticket); err != nil {
			break
		}
		stmts := []*models.Statement{}
		if err = g.Engine.Where("`ticket_id` = ?", ticket.TicketID).Find(&stmts); err != nil {
			break
		}
		cluster := &models.Cluster{}
		if _, err = g.Engine.ID(ticket.ClusterID).Get(cluster); err != nil {
			break
		}
		var engine *xorm.Engine
		if engine, err = cluster.Connect(ticket.Database, passwd); err != nil {
			break
		}

		ticket.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumDone]
		for _, stmt := range stmts {
			var result sql.Result
			if result, err = engine.Exec(stmt.Content); err != nil {
				stmt.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumExecFailure]
				stmt.Results = err.Error()
				ticket.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumExecFailure]
				buf.WriteString("\nexecute statement error")
				buf.WriteString("\n")
				buf.WriteString(stmt.Content)
				break
			} else {
				stmt.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumDone]
				if ra, err := result.RowsAffected(); err == nil {
					stmt.RowsAffected = uint(ra)
				}
			}
		}
		if err != nil {
			defer events.Fire(events.EventTicketFailed, &events.TicketFailedArgs{
				Ticket:  *ticket,
				Cluster: *cluster,
			})
		} else {
			defer events.Fire(events.EventTicketExecuted, &events.TicketExecutedArgs{
				Ticket:  *ticket,
				Cluster: *cluster,
			})
		}
		session := g.Engine.NewSession()
		session.Begin()
		defer session.Close()
		if _, err = session.ID(ticket.TicketID).Update(ticket); err != nil {
			session.Rollback()
			break
		}
		for _, stmt := range stmts {
			if _, err := session.ID(core.PK{stmt.TicketID, stmt.Sequence}).Update(stmt); err != nil {
				session.Rollback()
				break L
			}
		}
		if err = session.Commit(); err != nil {
			break
		}

		break
	}
	buf.WriteString("\n")
	if err != nil {
		buf.WriteString(err.Error())
	} else {
		buf.WriteString("execute success")
	}
	name := fmt.Sprintf("/tmp/halocli-execute-%s.txt", ticketUUID)
	if ioutil.WriteFile(name, buf.Bytes(), 0644) == nil {
		os.Exit(0)
	}
}

// 大表结构变更自动使用gh-ost
// gh-ost -user=root \
//        -port=3306 \
//        -password=password \
//        -host=127.0.0.1 \
//        -allow-on-master \
//        -max-load='Threads_running=100,Threads_connected=500' \
//        -initially-drop-old-table \
//        -initially-drop-ghost-table \
//        -exact-rowcount=false \
//        -database=starwars \ /* 目标数据库 */
//        -table=t1 \          /* 目标表 */
//        -alter="add ghostc1 varchar(20) default '' comment 'ghost测试'" \ /* 使用Specs的Restore */
//        -execute \
//        -hooks-enabled \
//        -hooks-hint=${statementUUID}
