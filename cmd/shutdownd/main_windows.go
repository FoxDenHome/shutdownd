//go:build windows

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/FoxDenHome/shutdownd/listener"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

const serviceName = "shutdownd"

func usage() {
	fmt.Fprintf(os.Stderr, "Must be run with a sub-command from: install, uninstall, run")
}

func main() {
	inService, err := svc.IsWindowsService()
	if err != nil {
		panic(err)
	}

	runner := &listener.Listener{}

	if inService {
		logger, err := eventlog.Open(serviceName)
		if err != nil {
			panic(err)
		}
		runner.Logger = logger
		err = svc.Run(serviceName, runner)
		if err != nil {
			panic(err)
		}
		return
	}

	if len(os.Args) < 2 {
		usage()
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "run":
		runner.Logger = debug.New(serviceName)
		err = debug.Run(serviceName, runner)
		if err != nil {
			panic(err)
		}
	case "install":
		err := installService(serviceName, "ShutdownD")
		if err != nil {
			panic(err)
		} else {
			fmt.Fprintf(os.Stderr, "Service installed")
		}
	case "uninstall":
		err := uninstallService(serviceName)
		if err != nil {
			panic(err)
		} else {
			fmt.Fprintf(os.Stderr, "Service uninstalled")
		}
	}
}
