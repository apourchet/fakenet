package fakenet

import "net/http"

// Interceptor is the interface matches (or not) incoming requests, and
// contains predetermined responses and errors.
type Interceptor struct {
	Matches  func(*http.Request) bool
	Response *http.Response
	Error    error
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
		return req.URL.String() == url
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

// CatchURLInterceptor returns an interceptor that catches all requests with the path specified,
// and returns the response and error given as arguments.
func CatchURLInterceptor(url string, response *http.Response, err error) Interceptor {
	matcher := func(req *http.Request) bool {
		return req.URL.String() == url
	}
	return Interceptor{
		Matches:  matcher,
		Response: response,
		Error:    err,
	}
}
