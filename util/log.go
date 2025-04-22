package util

type Logger interface {
	Info(eventID uint32, msg string) error
	Error(eventID uint32, msg string) error
	Close() error
}
