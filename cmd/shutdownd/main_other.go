//go:build !windows

package main

import (
	"log"

	"github.com/FoxDenHome/shutdownd/listener"
	"github.com/FoxDenHome/shutdownd/util"
)

type genericLogger struct{}

func (*genericLogger) Close() error {
	return nil
}

func (*genericLogger) Error(eventID uint32, msg string) error {
	log.Printf("Error %d: %s", eventID, msg)
	return nil
}

func (*genericLogger) Info(eventID uint32, msg string) error {
	log.Printf("Info %d: %s", eventID, msg)
	return nil
}

func main() {
	log.Printf("ShutdownD version %s", util.Commit())

	runner := &listener.Listener{}
	runner.Logger = &genericLogger{}
	runner.Execute()
}
