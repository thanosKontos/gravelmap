package log

import (
	"log"
	"time"
)

type levelized struct {
	level  string
	logger *log.Logger
}

// NewLevelized initializes and returns an new levelized logging object.
func NewLevelized(level string, logger *log.Logger) *levelized {
	return &levelized{
		level:  level,
		logger: logger,
	}
}

func (l levelized) Debug(log interface{}) {
	if l.level != "debug" {
		return
	}

	t := time.Now().Format("2006/01/02 15:04:05")
	l.logger.Println(t, log)
}

func (l levelized) Info(log interface{}) {
	if l.level != "info" && l.level != "debug" {
		return
	}

	t := time.Now().Format("2006/01/02 15:04:05")
	l.logger.Println(t, log)
}

func (l levelized) Warning(log interface{}) {
	if l.level != "info" && l.level != "debug" && l.level != "warning" {
		return
	}

	t := time.Now().Format("2006/01/02 15:04:05")
	l.logger.Println(t, log)
}

func (l levelized) Error(log interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05")
	l.logger.Println(t, log)
}
