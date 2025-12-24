package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/techtacles/sysmonitoring/internal/dashboard"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

var logtag string = "cmd"

var kubeconfigpath string
var collectDocker bool
var isDetached bool
var collectKubernetes bool

var RunCmd = &cobra.Command{
	Use:   "start",
	Short: "Run the dashboard server",
	Long:  `Run the dashboard server to display system metrics in a web interface.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if isDetached {
			// Find and remove the detached flag from arguments
			newArgs := []string{}
			for _, arg := range os.Args[1:] {
				if arg != "-D" && arg != "--detached" {
					newArgs = append(newArgs, arg)
				}
			}

			// Spawn a new process
			execPath, err := os.Executable()
			if err != nil {
				return err
			}

			child := exec.Command(execPath, newArgs...)
			child.Stdout = nil
			child.Stderr = nil
			child.Stdin = nil

			if err := child.Start(); err != nil {
				logging.Error(logtag, "failed to start background process", err)
				return err
			}

			logging.Info(logtag, fmt.Sprintf("dashboard server started in background on http://localhost:%s (PID: %d)", dashboard.Port, child.Process.Pid))
			os.Exit(0)
		}

		logging.Info(logtag, "running dashboard server")
		if err := dashboard.Run(collectDocker, collectKubernetes, kubeconfigpath); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(RunCmd)
	RunCmd.Flags().StringVarP(&dashboard.Port, "port", "p", "8080", "Port to run the dashboard server on")
	RunCmd.Flags().BoolVarP(&collectDocker, "docker", "d", false, "Whether to collect docker metrics. Make sure docker is running when passing this flag")
	RunCmd.Flags().BoolVarP(&collectKubernetes, "kubernetes", "k", false, "Whether to collect kubernetes metrics.")
	RunCmd.Flags().StringVarP(&kubeconfigpath, "kubeconfig", "", "", "absolute path to the kubeconfig file (optional)")
	RunCmd.Flags().BoolVarP(&isDetached, "detached", "D", false, "Run the dashboard server in the background")
}
