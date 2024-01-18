//go:build !windows

package main

import "net/http"

type shutdownHandler struct {
	logger Logger
}

func (h *shutdownHandler) onReady(*http.Server) {
}
