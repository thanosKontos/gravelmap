package cli

type nullCLI struct {
}

// NewDebugCLI initialize and return an new DebugCLI object.
func NewNullCLI() *nullCLI {
	return &nullCLI{}
}

func (nullCLI) Debug(log interface{}) {
}

func (nullCLI) Info(log interface{}) {
}

func (nullCLI) Warning(log interface{}) {
}

func (nullCLI) Error(log interface{}) {
}
