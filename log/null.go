package log

type null struct {
}

// NewNullLog initialize and return an new DebugCLI object.
func NewNullLog() *null {
	return &null{}
}

func (null) Debug(log interface{}) {
}

func (null) Info(log interface{}) {
}

func (null) Warning(log interface{}) {
}

func (null) Error(log interface{}) {
}
