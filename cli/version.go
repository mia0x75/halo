package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mia0x75/halo/g"
)

// versionCmd represents the bye command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  `This subcommand shows version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n",
			"Version", g.Version,
			"Git commit", g.Git,
			"Compile", g.Compile,
			"Distro", g.Distro,
			"Kernel", g.Kernel,
			"Branch", g.Branch,
		)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// byeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// byeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
