package user

import (
	"os/user"
	"runtime"

	"github.com/techtacles/sysmonitoring/internal/logging"
)

const logtag string = "user"

type UserInfo struct {
	Username string
	FullName string
	HomeDir  string
	Runtime  string
	Arch     string
}

func (u *UserInfo) Collect() error {
	if err := u.collectUser(); err != nil {
		return err
	}
	u.Runtime = runtime.GOOS
	u.Arch = runtime.GOARCH

	logging.Info(logtag, "successfully collected user info")
	return nil
}

func (u *UserInfo) collectUser() error {
	current, err := user.Current()
	if err != nil {
		logging.Error(logtag, "error getting current user", err)
		return err
	}

	u.Username = current.Username
	u.FullName = current.Name
	u.HomeDir = current.HomeDir

	return nil
}
