package mock

import (
	"bytes"
)

type NamedBuffer struct {
	Buffer *bytes.Buffer
}

func (b NamedBuffer) Name() string {
	return "MockNamedBuffer"
}

func (b NamedBuffer) Read(p []byte) (int, error) {
	return b.Buffer.Read(p)
}

func (b NamedBuffer) Write(p []byte) (int, error) {
	return b.Buffer.Write(p)
}
