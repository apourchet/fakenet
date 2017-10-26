package fakenet

import (
	"io"
	"strings"
)

type readCloser struct{ io.Reader }

func (closer readCloser) Close() error { return nil }

func NewReadCloser(content string) io.ReadCloser {
	reader := strings.NewReader(content)
	return readCloser{reader}
}
