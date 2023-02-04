package mock

import (
	"bytes"
)

type MockNamedBuffer struct {
	Buffer *bytes.Buffer
}

func (b MockNamedBuffer) Name() string {
	return "MockNamedBuffer"
}

func (b MockNamedBuffer) Read(p []byte) (n int, err error) {
	return b.Buffer.Read(p)
}

func (b MockNamedBuffer) Write(p []byte) (n int, err error) {
	return b.Buffer.Write(p)
}
