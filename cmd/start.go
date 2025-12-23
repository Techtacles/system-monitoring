package cmd

import (
	"github.com/spf13/cobra"
	"github.com/techtacles/sysmonitoring/internal/dashboard"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

var logtag string = "cmd"

var collectDocker bool

var RunCmd = &cobra.Command{
	Use:   "start",
	Short: "Run the dashboard server",
	Long:  `Run the dashboard server to display system metrics in a web interface.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logging.Info(logtag, "running dashboard server")
		return dashboard.Run(collectDocker)
	},
}

func init() {
	rootCmd.AddCommand(RunCmd)
	RunCmd.Flags().StringVarP(&dashboard.Port, "port", "p", "8080", "Port to run the dashboard server on")
	RunCmd.Flags().BoolVarP(&collectDocker, "docker", "d", false, "Whether to collect docker metrics. Make sure docker is running when passing this flag")

}
