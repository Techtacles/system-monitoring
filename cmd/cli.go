package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

var rootCmd = &cobra.Command{
	Use:   "sysmonitoring",
	Short: "System Monitoring Tool",
	Long:  `A simple system monitoring tool that provides real-time insights into system performance, resource usage, and system health.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logging.Error(logtag, "error executing command", err)
		os.Exit(1)
	}
}
