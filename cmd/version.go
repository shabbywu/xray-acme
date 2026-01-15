package cmd

import (
	"fmt"
	"github.com/xtls/xray-core/core"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current version of Xray",
	Long:  `Version prints the build information for Xray executables.`,
	Run:   runVersionCmd,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersionCmd(cmd *cobra.Command, args []string) {
	version := core.VersionStatement()
	for _, s := range version {
		fmt.Println(s)
	}
}
