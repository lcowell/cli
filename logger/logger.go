package logger

import (
	"fmt"
	"io"
	"os"
)

var (
	verbose bool
	writer  io.Writer
)

// SetVerbose enables output when using Debug
func SetVerbose(v bool) {
	verbose = v
}

// SetWriter sets the output target
func SetWriter(w io.Writer) {
	writer = w
}

// Info writes output with an info tag
func Info(i ...interface{}) {
	write("Info>", i...)
}

// Error write output with an error tag
func Error(i ...interface{}) {
	write("Error>", i...)
}

// Debug writes output with a debug tag when verbose mode is enabled
func Debug(i ...interface{}) {
	if verbose {
		write("Debug>", i...)
	}
}

func write(tag string, i ...interface{}) {
	o := append([]interface{}{tag}, i...)
	fmt.Fprintln(writer, o...)
}

func init() {
	SetWriter(os.Stdout)
}
