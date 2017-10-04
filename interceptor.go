package fakenet

import "net/http"

// Interceptor is the interface matches (or not) incoming requests, and
// contains predetermined responses and errors.
type Interceptor struct {
	Matches  func(*http.Request) bool
	Response *http.Response
	Error    error
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
