//go:build windows

package listener

import (
	"net/http"
	"os/exec"

	"github.com/FoxDenHome/shutdownd/util"
	"golang.org/x/sys/windows/svc"
)

const shutdown_binary = "C:\\Windows\\System32\\shutdown.exe"

type Listener struct {
	Logger  util.Logger
	r       <-chan svc.ChangeRequest
	changes chan<- svc.Status
}

const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

func (h *Listener) doShutdown() error {
	return exec.Command(shutdown_binary, "-s", "-f", "-t", "60").Run()
}

func (h *Listener) doShutdownAbort() error {
	return exec.Command(shutdown_binary, "-a").Run()
}

func (h *Listener) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	h.changes = changes
	h.r = r

	changes <- svc.Status{State: svc.StartPending}

	return h.execute()
}

func (h *Listener) onReady(server *http.Server) {
	go func() {
		for {
			c := <-h.r
			switch c.Cmd {
			case svc.Interrogate:
				h.changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				server.Close()
				return
			}
		}
	}()

	h.changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
}
