package fakenet

import (
	"net/http"
	"path/filepath"
)

// Interceptor is the interface matches (or not) incoming requests, and
// contains predetermined responses and errors.
type Interceptor struct {
	Matches  func(*http.Request) bool
	Response *http.Response
	Error    error
}

// Match returns whether or not the request is caught by that interceptor.
func (ic Interceptor) Match(req *http.Request) bool {
	if ic.Matches == nil {
		return false
	}
	return ic.Matches(req)
}

// GetResponse returns the response of that interceptor given the request that matched it.
func (ic Interceptor) GetResponse(_ *http.Request) (*http.Response, error) {
	return ic.Response, ic.Error
}

// WithBody returns a new interceptor that returns a body in the response.
func (ic Interceptor) WithBody(body string) Interceptor {
	if ic.Response == nil {
		ic.Response = &http.Response{}
	}
	ic.Response.Body = NewReadCloser(body)
	return ic
}

// WithURLMatcher returns a new interceptor to match requests with the url provided.
func (ic Interceptor) WithURLMatcher(url string) Interceptor {
	ic.Matches = func(req *http.Request) bool {
		match, err := filepath.Match(url, req.URL.String())
		return err == nil && match
	}
	return ic
}

// WithStatus returns a new interceptor with a status code in the response.
func (ic Interceptor) WithStatus(code int) Interceptor {
	if ic.Response == nil {
		ic.Response = &http.Response{}
	}
	ic.Response.StatusCode = code
	return ic
}

// WithHeader returns a new interceptor with a header value.
func (ic Interceptor) WithHeader(key string, value string, others ...string) Interceptor {
	if ic.Response == nil {
		ic.Response = &http.Response{}
	}
	if ic.Response.Header == nil {
		ic.Response.Header = map[string][]string{}
	}
	ic.Response.Header[key] = append([]string{value}, others...)
	return ic
}

// CatchAllInterceptor returns an interceptor that catches all requests and returns
// the response and error given as arguments.
func CatchAllInterceptor(response *http.Response, err error) Interceptor {
	return Interceptor{
		Matches:  func(_ *http.Request) bool { return true },
		Response: response,
		Error:    err,
	}
}

// CatchURLInterceptor returns an interceptor that catches all requests with the URL pattern specified,
// and returns the response and error given as arguments.
// See https://golang.org/pkg/path/filepath/#Match for information about what can be in patterns.
func CatchURLInterceptor(url string, response *http.Response, err error) Interceptor {
	matcher := func(req *http.Request) bool {
		match, err := filepath.Match(url, req.URL.String())
		return err == nil && match
	}
	return Interceptor{
		Matches:  matcher,
		Response: response,
		Error:    err,
	}
}
