package api

import "github.com/spf13/cobra"

var (
	apiServerAddrPtr string
	apiTimeout       int
	apiJSON          bool

	pattern string
)

// Cmd the api command
var Cmd = &cobra.Command{
	Use: "api",
}

func init() {
	Cmd.PersistentFlags().StringVarP(&apiServerAddrPtr, "server", "s", "127.0.0.1:8080",
		"api server addr")
	Cmd.PersistentFlags().IntVarP(&apiTimeout, "timeout", "t", 3,
		"超时时间")
	Cmd.PersistentFlags().BoolVar(&apiJSON, "json", false,
		"是否以 json 格式输出")
}
