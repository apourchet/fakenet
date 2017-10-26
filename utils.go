package fakenet

import (
	"io"
	"strings"
)

type readCloser struct{ io.Reader }

func (closer readCloser) Close() error { return nil }

// Returns a new io.ReadCloser from a string input. This is useful when creating http.Response objects
// with non-empty bodies.
func NewReadCloser(content string) io.ReadCloser {
	reader := strings.NewReader(content)
	return readCloser{reader}
}
