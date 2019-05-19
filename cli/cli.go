package cli

import (
	"fmt"
	"time"
)

type CLI struct {
}

// NewCLI initialize and return an new CLI object.
func NewCLI() *CLI {
	return &CLI{}
}

func (CLI) Debug(log interface{}) {
	t := time.Now().String()
	fmt.Println(t, log)
}

func (CLI) Info(log interface{}) {
	t := time.Now().String()
	fmt.Println(t, log)
}

func (CLI) Warning(log interface{}) {
	t := time.Now().String()
	fmt.Println(t, log)
}

func (CLI) Error(log interface{}) {
	t := time.Now().String()
	fmt.Println(t, log)
}
