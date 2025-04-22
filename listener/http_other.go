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
	return exec.Command("shutdown", "-P", "1").Run()
}

func (h *Listener) doShutdownAbort() error {
	return exec.Command("shutdown", "-c").Run()
}

func (h *Listener) onReady(*http.Server) {
}

func (h *Listener) Execute() (bool, uint32) {
	return h.execute()
}
