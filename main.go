package main

import (
	"github.com/techtacles/sysmonitoring/internal/dashboard"
	"github.com/techtacles/sysmonitoring/internal/logging"
)

func main() {
	logging.Info("main", "starting application")
	if err := dashboard.Run(); err != nil {
		logging.Error("main", "dashboard server error", err)
	}
}
