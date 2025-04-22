//go:build !windows

package listener

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/FoxDenHome/shutdownd/util"
)

type Listener struct {
	Logger util.Logger
}

func sudoIfNeeded(arg ...string) *exec.Cmd {
	if os.Getuid() != 0 {
		return exec.Command("sudo", arg...)
	}
	return exec.Command(arg[0], arg[1:]...)
}

func (h *Listener) doShutdown() error {
	return sudoIfNeeded("shutdown", "-P", "1").Run()
}

func (h *Listener) doShutdownAbort() error {
	return sudoIfNeeded("shutdown", "-c").Run()
}

func (h *Listener) onReady(*http.Server) {
}

func (h *Listener) Execute() (bool, uint32) {
	return h.execute()
}
