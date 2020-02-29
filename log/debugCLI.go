package log

import (
	"fmt"
	"time"
)

type debugCLI struct {
}

// NewDebugCLI initialize and return an new DebugCLI object.
func NewDebugCLI() *debugCLI {
	return &debugCLI{}
}

func (debugCLI) Debug(log interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(t, log)
}

func (debugCLI) Info(log interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(t, log)
}

func (debugCLI) Warning(log interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(t, log)
}

func (debugCLI) Error(log interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Println(t, log)
}
