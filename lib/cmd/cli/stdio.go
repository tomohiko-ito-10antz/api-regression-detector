package cli

import (
	"io"
)

type Stdio struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}
