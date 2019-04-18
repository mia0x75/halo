package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// ghostCmd represents the execute command
var ghostCmd = &cobra.Command{
	Use:   "ghost",
	Short: "Gh-ost hooks",
	Long:  `Github online schema transmogrifier hooks`,
	Run:   ghost,
}

func init() {
	RootCmd.AddCommand(ghostCmd)
	ghostCmd.PersistentFlags().BoolP("on-startup", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-validated", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-rowcount-complete", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-before-row-copy", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-row-copy-complete", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-begin-postponed", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-before-cut-over", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-interactive-command", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-success", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-failure", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-status", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-stop-replication", "", false, "")
	ghostCmd.PersistentFlags().BoolP("on-start-replication", "", false, "")
}

// onStartup -> onValidated -> onBeforeRowCopy -> onRowCopyComplete -> onBeforeCutOver -> onSuccess
func ghost(cmd *cobra.Command, args []string) {
	onStartup, _ := cmd.Flags().GetBool("on-startup")
	onValidated, _ := cmd.Flags().GetBool("on-validated")
	onRowcountComplete, _ := cmd.Flags().GetBool("on-rowcount-complete")
	onBeforeRowcopy, _ := cmd.Flags().GetBool("on-before-row-copy")
	onRowcopyComplete, _ := cmd.Flags().GetBool("on-row-copy-complete")
	onBeginPostponed, _ := cmd.Flags().GetBool("on-begin-postponed")
	onBeforeCutover, _ := cmd.Flags().GetBool("on-before-cut-over")
	onInteractiveCommand, _ := cmd.Flags().GetBool("on-interactive-command")
	onSuccess, _ := cmd.Flags().GetBool("on-success")
	onFailure, _ := cmd.Flags().GetBool("on-failure")
	onStatus, _ := cmd.Flags().GetBool("on-status")
	onStopReplication, _ := cmd.Flags().GetBool("on-stop-replication")
	onStartReplication, _ := cmd.Flags().GetBool("on-start-replication")

	var err error

	for {
		statementUUID, ok := os.LookupEnv("GH_OST_HOOKS_HINT")
		if !ok {
			break
		}
		if !isValidUUID(statementUUID) {
			fmt.Printf("String %s is not a valid UUID.", statementUUID)
			break
		}

		name := fmt.Sprintf("/tmp/halocli-ghost-%s.txt", statementUUID)
		if ioutil.WriteFile(name, []byte(statementUUID), 0644) == nil {
		}

		if cfg, ok := os.LookupEnv("HALO_CFG"); !ok {
			err = fmt.Errorf("Missing halo config")
			break
		} else {
			g.ParseConfig(cfg)
			g.InitDB()
		}

		statement := &models.Statement{}
		if _, err = g.Engine.Where("`uuid` = ?", statementUUID).Get(statement); err != nil {
			break
		}

		if onStartup {
			fmt.Println("halocli ghost --onStartup invoked")
			break
		}
		if onValidated {
			fmt.Println("halocli ghost --onValidated invoked")
			break
		}
		if onRowcountComplete {
			fmt.Println("halocli ghost --onRowcountComplete invoked")
			break
		}
		if onBeforeRowcopy {
			fmt.Println("halocli ghost --onBeforeRowcopy invoked")
			break
		}
		if onRowcopyComplete {
			fmt.Println("halocli ghost --onRowcopyComplete invoked")
			break
		}
		if onBeginPostponed {
			fmt.Println("halocli ghost --onBeginPostponed invoked")
			break
		}
		if onBeforeCutover {
			fmt.Println("halocli ghost --onBeforeCutover invoked")
			break
		}
		if onInteractiveCommand {
			fmt.Println("halocli ghost --onInteractiveCommand invoked")
			break
		}
		if onSuccess {
			fmt.Println("halocli ghost --onSuccess invoked")
			break
		}
		if onFailure {
			fmt.Println("halocli ghost --onFailure invoked")
			break
		}
		if onStatus {
			if _, ok := os.LookupEnv("GH_OST_STATUS"); !ok {
				break
			}
			fmt.Println("halocli ghost --onStatus invoked")
			break
		}
		if onStopReplication {
			fmt.Println("halocli ghost --onStopReplication invoked")
			break
		}
		if onStartReplication {
			fmt.Println("halocli ghost --onStartReplication invoked")
			break
		}
		break
	}
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
