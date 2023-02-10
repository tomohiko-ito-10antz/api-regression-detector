package mock

import (
	"bytes"

	"github.com/Jumpaku/api-regression-detector/test"
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

type ErrNamedBuffer struct{}

func (b ErrNamedBuffer) Name() string {
	return "ErrMockNamedBuffer"
}

func (b ErrNamedBuffer) Read(p []byte) (int, error) {
	return 0, test.MockError
}

func (b ErrNamedBuffer) Write(p []byte) (int, error) {
	return 0, test.MockError
}
