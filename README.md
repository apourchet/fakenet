# FakeNet
A Golang HTTP client that can dropped-in for `*http.Client` or any interface that `*http.Client` satisfies. It gives testing flexibility through the use of request interceptors, bypassing the network stack when certain user-described conditions are met. This package is useful in the context of unit tests, where we do not want to depend on the network and need consistent pre-determined responses to requests.

### Example
The following example illustrates how one could implement a website status checker package. In a production/live environment, we want the requests to be sent through to the network; but in unit tests we want to setup pre-determined responses to the GET requests that it issues.
```go
// website_checker.go
package checker

type Checker struct {
	NetworkClient HTTPClient
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

func New(client HTTPClient) Checker {
	return Checker{
    	NetworkClient: client,
    }
}

func (checker Checker) IsWebsiteUp() bool {
	resp, err := checker.NetworkClient.Get("http://example.org")
    return err == nil && resp.StatusCode == 200
}
```
And now for the unit tests:
```go
// website_checker_test.go
package checker_test

import "testing"

func TestErrorResponse(t *testing.T) {
	client := fakenet.New()
    catchall := fakenet.CatchAllInterceptor(nil, errors.New("Fell through to the catch all"))
	client.Intercept(catchall)
    
    checker := checker.New(client)
    if checker.IsWebsiteUp() {
		t.Fail("Checker should say website is down if we receive an error in the request")
    }
}

func TestBadStatusCode(t *testing.T) {
	client := fakenet.New()
	client.CatchAll(http.StatusInternalServerError, "Failed to load website.")
    
    checker := checker.New(client)
    if checker.IsWebsiteUp() {
		t.Fail("Checker should say website is down if we receive non-200 status code")
    }
}

func TestWebsiteIsLive(t *testing.T) {
	client := fakenet.New()	
    client.CatchAll(http.StatusOK, "Website is live!")
    
    checker := checker.New(client)
    if !checker.IsWebsiteUp() {
    	t.Fail("Checker should say website is up if the status code is 200")
    }
}
```
