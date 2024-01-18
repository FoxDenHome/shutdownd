//go:build windows

package main

import (
	"net/http"

	"golang.org/x/sys/windows/svc"
)

type shutdownHandler struct {
	logger  Logger
	r       <-chan svc.ChangeRequest
	changes chan<- svc.Status
}

const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

func (h *shutdownHandler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	h.changes = changes
	h.r = r

	changes <- svc.Status{State: svc.StartPending}

	return h.execute(args)
}

func (h *shutdownHandler) onReady(server *http.Server) {
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
