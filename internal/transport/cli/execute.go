package cli

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/event"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/platform/crypto"
	"github.com/mia0x75/halo/internal/service"
	"github.com/mia0x75/halo/internal/worker"
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
	ctx := context.Background()
	for {
		if cfg, ok := os.LookupEnv("HALO_CFG"); !ok {
			err = fmt.Errorf("Missing halo config")
			break
		} else {
			buf.WriteString(cfg)
			config.ParseConfig(cfg)
		}

		application, appErr := newApp()
		if appErr != nil {
			err = fmt.Errorf("错误代码: 1500, 错误信息: %s", appErr.Error())
			break
		}
		worker.Init(application.Repos)
		if ticketUUID, err = cmd.Flags().GetString("ticket"); err != nil {
			break
		}
		if !isValidUUID(ticketUUID) {
			err = fmt.Errorf("String %s is not a valid UUID", ticketUUID)
			break
		}

		buf.WriteString("\n")
		buf.WriteString(ticketUUID)

		ticket, tkErr := application.Services.Ticket.GetByUUID(ctx, ticketUUID)
		if tkErr != nil {
			err = fmt.Errorf("错误代码: 1500, 错误信息: %s", tkErr.Error())
			break
		}

		// 工单当前状态是LGTM才允许执行
		if ticket.Status != domain.TicketStatusLgtm {
			err = fmt.Errorf("错误代码: 1500, 错误信息: 只有审核通过等待上线执行的工单才可以执行。")
			break
		}

		stmts, stErr := application.Services.Statement.ListByTicket(ctx, ticket.TicketID)
		if stErr != nil {
			err = fmt.Errorf("错误代码: 1500, 错误信息: %s", stErr.Error())
			break
		}

		target, clErr := application.Services.Cluster.Get(ctx, ticket.ClusterID)
		if clErr != nil {
			err = fmt.Errorf("错误代码: 1500, 错误信息: %s", clErr.Error())
			break
		}
		target.Database = ticket.Database

		password, decErr := crypto.DecryptAES(target.Password, config.Config().Secret.Crypto)
		if decErr != nil {
			err = fmt.Errorf("错误代码: 1500, 错误信息: 群集(uuid=%s)的密码无法解密，请重新接入该群集。", target.UUID)
			break
		}

		ticket.Status = domain.TicketStatusDone
		for _, stmt := range stmts {
			result, execErr := application.Agent.Execute(ctx, target, string(password), stmt.Content)
			if execErr != nil {
				stmt.Status = domain.StatementStatusExecFailure
				message := execErr.Error()
				stmt.Results = &message
				ticket.Status = domain.TicketStatusExecFailure
				buf.WriteString(stmt.Content)
				err = fmt.Errorf("错误代码: 1500, 错误信息: %s", execErr.Error())
				break
			}
			stmt.Status = domain.StatementStatusDone
			rows := uint(result.RowsAffected)
			stmt.RowsAffected = &rows
		}
		// Execution is audited as user 0: the CLI acts on its own rather than
		// on behalf of a logged-in operator.
		actorID := uint(0)

		// The event bus speaks the legacy DTO, so the conversion happens only
		// at this boundary; the execution path itself stays on domain types.
		eventTicket := bridge.TicketFromDomain(ticket)
		eventCluster := bridge.ClusterFromDomain(target)
		if err != nil {
			event.FireSync(event.EventTicketFailed, &event.TicketFailed{
				UserID:  actorID,
				Ticket:  *eventTicket,
				Cluster: *eventCluster,
			})
		} else {
			event.FireSync(event.EventTicketExecuted, &event.TicketExecuted{
				UserID:  actorID,
				Ticket:  *eventTicket,
				Cluster: *eventCluster,
			})
		}
		err = application.Tx.WithTx(ctx, func(ctx context.Context, repos service.Repos) error {
			if err := repos.Ticket.Update(ctx, ticket); err != nil {
				return err
			}
			for _, stmt := range stmts {
				if err := repos.Statement.Update(ctx, stmt); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			err = fmt.Errorf("错误代码: 1500, 错误信息: %s", err.Error())
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
