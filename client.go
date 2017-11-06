package fakenet

import (
	"net/http"
	"sync"
)

// HTTPClient is the struct that embeds *http.Client, and is returned by New.
// It is through this struct that most of the FakeNet operations will be performed.
// HTTPClient is thread-safe for all operations, just like the standard lib's client.
type HTTPClient struct {
	catchers []RequestCatcher
	lock     *sync.Mutex

	*http.Client
}

// RequestCatcher is the interface needed for a FakeNet to intercept and respond
// to requests.
type RequestCatcher interface {
	Match(req *http.Request) bool
	GetResponse(req *http.Request) (*http.Response, error)
}

// New creates a new FakeNet HTTP client.
func New() *HTTPClient {
	client := &HTTPClient{
		lock:   &sync.Mutex{},
		Client: &http.Client{},
	}
	client.Client.Transport = client
	return client
}

// Intercept adds a request interceptor to the http client. New interceptors will get priority over the old ones.
func (client *HTTPClient) Intercept(catcher RequestCatcher) {
	client.lock.Lock()
	defer client.lock.Unlock()
	client.catchers = append(client.catchers, catcher)
}

// RoundTrip is the implementation of the http.RoundTripper interface.
func (client *HTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	client.lock.Lock()
	for i := len(client.catchers) - 1; i >= 0; i-- {
		catcher := client.catchers[i]
		if catcher.Match(req) {
			client.lock.Unlock()
			return catcher.GetResponse(req)
		}
	}
	client.lock.Unlock()
	return http.DefaultTransport.RoundTrip(req)
}

// CatchAll creates a pre-determined response for all requests going through this client.
func (client *HTTPClient) CatchAll(code int, response string) {
	body := NewReadCloser(response)
	resp := &http.Response{StatusCode: code, Body: body}
	interceptor := CatchAllInterceptor(resp, nil)
	client.Intercept(interceptor)
}

// InterceptURL creates a pre-determined response for all requests that have the URL specified. The URL can
// also be a pattern that follows the patterns of the filepath.Match function: https://golang.org/pkg/path/filepath/#Match.
func (client *HTTPClient) InterceptURL(url string, code int, response string) {
	body := NewReadCloser(response)
	resp := &http.Response{StatusCode: code, Body: body}
	interceptor := CatchURLInterceptor(url, resp, nil)
	client.Intercept(interceptor)
}
