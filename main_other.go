//go:build !windows

package main

import "log"

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
	runner := &shutdownHandler{}
	runner.logger = &genericLogger{}
	runner.execute()
}
