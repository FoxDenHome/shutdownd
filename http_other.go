//go:build !windows

package main

import (
	"net/http"
	"os/exec"
)

type shutdownHandler struct {
	logger Logger
}

func (h *shutdownHandler) doShutdown() error {
	return exec.Command("shutdown", "-P", "60").Run()
}

func (h *shutdownHandler) doShutdownAbort() error {
	return exec.Command("shutdown", "-c").Run()
}

func (h *shutdownHandler) onReady(*http.Server) {
}
