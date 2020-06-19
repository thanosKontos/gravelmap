package routing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrLoop(t *testing.T) {
	expectedErr := fmt.Sprint(ErrLoopDetected.Error(), "From node '", 0, "' to node '", 1, "'")
	assert.Equal(t, expectedErr, newErrLoop(0, 1).Error())
}
