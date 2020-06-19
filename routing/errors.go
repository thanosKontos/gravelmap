package routing

import (
	"errors"
	"fmt"
)

//ErrNoPath is thrown when there is no path from src to dest
var ErrNoPath = errors.New("no path found")

//ErrLoopDetected is thrown when a loop is detected, causing the cost to go
// to inf (or -inf), or just generally loop forever
var ErrLoopDetected = errors.New("infinite loop detected")

//NewErrLoop generates a new error with details for loop error
func newErrLoop(a, b int) error {
	return errors.New(fmt.Sprint(ErrLoopDetected.Error(), "From node '", a, "' to node '", b, "'"))
}
