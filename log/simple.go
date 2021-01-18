package log

import (
	"log"
	"time"
)

type simple struct {
	level  string
	logger *log.Logger
}

// NewSimple initializes and returns an new simple logging object.
func NewSimple(level string, logger *log.Logger) *simple {
	return &simple{
		level:  level,
		logger: logger,
	}
}

func (l simple) Debug(log interface{}) {
	if l.level != "debug" {
		return
	}

	t := time.Now().Format("2006/01/02 15:04:05")
	l.logger.Println(t, log)
}

func (l simple) Info(log interface{}) {
	if l.level != "info" && l.level != "debug" {
		return
	}

	t := time.Now().Format("2006/01/02 15:04:05")
	l.logger.Println(t, log)
}

func (l simple) Warning(log interface{}) {
	if l.level != "info" && l.level != "debug" && l.level != "warning" {
		return
	}

	t := time.Now().Format("2006/01/02 15:04:05")
	l.logger.Println(t, log)
}

func (l simple) Error(log interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05")
	l.logger.Println(t, log)
}
