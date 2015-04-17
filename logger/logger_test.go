package logger

import (
	"bytes"
	"testing"

	"github.com/bmizerany/assert"
)

func TestInfo(t *testing.T) {
	var b bytes.Buffer
	SetWriter(&b)

	Info("hello", "world")

	assert.Equal(t, b.String(), "Info> hello world\n")
}

func TestError(t *testing.T) {
	var b bytes.Buffer
	SetWriter(&b)

	Error("hello", "world")

	assert.Equal(t, b.String(), "Error> hello world\n")
}

func TestDebug(t *testing.T) {
	var b bytes.Buffer
	SetWriter(&b)

	Debug("hello")
	SetVerbose(true)
	Debug("world")

	assert.Equal(t, b.String(), "Debug> world\n")
}
