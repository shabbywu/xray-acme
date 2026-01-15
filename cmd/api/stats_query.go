package api

import (
	"fmt"
	"github.com/spf13/cobra"
	statsService "github.com/xtls/xray-core/app/stats/command"
)

// statsQueryCmd represents the statsquery command
var statsQueryCmd = &cobra.Command{
	Use:   "statsquery [--server=127.0.0.1:8080] [--pattern '']",
	Short: "Query statistics",
	Long:  `Query statistics from Xray.`,
	RunE:  runQueryStats,
}

func init() {
	statsQueryCmd.Flags().StringVar(&pattern, "pattern", "", "Filter pattern for the statistics query")
	Cmd.AddCommand(statsQueryCmd)
}

func runQueryStats(cmd *cobra.Command, args []string) error {
	conn, ctx, close := dialAPIServer()
	defer close()

	client := statsService.NewStatsServiceClient(conn)
	r := &statsService.QueryStatsRequest{
		Pattern: pattern,
		Reset_:  false,
	}
	resp, err := client.QueryStats(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to query stats: %s", err)
	}
	showJSONResponse(resp)
	return err
}
