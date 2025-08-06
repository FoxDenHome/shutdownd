//go:build !windows

package listener

import (
	"net/http"
	"os/exec"

	"github.com/FoxDenHome/shutdownd/util"
)

type Listener struct {
	Logger util.Logger
}

func (h *Listener) doShutdown() error {
	return exec.Command("/usr/bin/systemctl", "start", "shutdownd-run.service").Run()
}

func (h *Listener) doShutdownAbort() error {
	return exec.Command("/usr/bin/systemctl", "stop", "shutdownd-run.service").Run()
}

func (h *Listener) onReady(*http.Server) {
}

func (h *Listener) Execute() (bool, uint32) {
	return h.execute()
}
