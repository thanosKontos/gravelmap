package log

import (
	"fmt"
	"time"
)

type stdout struct {
	level string
}

// NewStdout initializes and returns an new DebugCLI object.
func NewStdout(level string) *stdout {
	return &stdout{
		level: level,
	}
}

func (l stdout) Debug(log interface{}) {
	if l.level != "debug" {
		return
	}

	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(t, log)
}

func (l stdout) Info(log interface{}) {
	if l.level != "info" && l.level != "debug" {
		return
	}

	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(t, log)
}

func (l stdout) Warning(log interface{}) {
	if l.level != "info" && l.level != "debug" && l.level != "warning" {
		return
	}

	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(t, log)
}

func (l stdout) Error(log interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(t, log)
}
