package events

import "go.uber.org/zap"

type Bus interface {
	Publish(event Event) error
	Subscribe(event string, handler Handler) error
}

type Handler func(event Event, logger *zap.SugaredLogger) error
