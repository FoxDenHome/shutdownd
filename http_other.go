//go:build !windows

package main

import "net/http"

type shutdownHandler struct {
	logger   Logger
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *shutdownHandler) onReady(*http.Server) {
}
