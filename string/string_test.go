package string

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExistsHappyPath(t *testing.T) {
	assert.True(t, String("foo").Exists([]string{"foo", "bar"}))
	assert.True(t, String("foo").Exists([]string{"bar", "foo"}))
	assert.True(t, String("foo").Exists([]string{"foo"}))
	assert.True(t, String("").Exists([]string{"bar", "foo", ""}))
}

func TestExistsUnhappyPath(t *testing.T) {
	assert.False(t, String("foo").Exists([]string{"blah", "bar"}))
	assert.False(t, String("").Exists([]string{"bar", "foo"}))
	assert.False(t, String("foo").Exists([]string{}))
	assert.False(t, String("").Exists([]string{}))
}
