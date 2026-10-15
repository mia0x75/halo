package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/internal/platform/db"
)

// newApp builds the dependency container from the configuration already loaded
// by the command (g.ParseConfig must have run first).
func newApp() (*app.App, error) {
	return app.New(context.Background(), db.Config{
		Addr:           g.Config().Database.Addr,
		MaxIdle:        g.Config().Database.MaxIdle,
		MaxConnections: g.Config().Database.MaxConnections,
		WaitTimeout:    g.Config().Database.WaitTimeout,
		Driver:         g.Config().Database.Driver,
	})
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "halocli",
	Short: "halo cli program",
	Long: `This application shows how to create modern CLI 
applications in go using Cobra CLI library`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
