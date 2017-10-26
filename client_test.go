package fakenet_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/apourchet/fakenet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatchAll(t *testing.T) {
	client := fakenet.New()
	client.CatchAll(http.StatusOK, "OK")

	resp, err := client.Get("http://example.org/")
	require.Nil(t, err)
	require.NotNil(t, resp)

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	assert.Equal(t, "OK", string(content))
}

func TestPathInterceptor(t *testing.T) {
	client := fakenet.New()
	catchall := fakenet.CatchAllInterceptor(nil, errors.New("Fell through to the catch all"))
	client.Intercept(catchall)
	client.InterceptURL("http://example.org", http.StatusOK, "OK")

	resp, err := client.Get("http://example.org")
	require.Nil(t, err)
	require.NotNil(t, resp)

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	assert.Equal(t, "OK", string(content))

	resp, err = client.Get("http://unknown.org")
	require.NotNil(t, err)
}

func TestResponseBody(t *testing.T) {
	client := fakenet.New()
	catchall := fakenet.CatchAllInterceptor(nil, errors.New("Fell through to the catch all"))
	client.Intercept(catchall)

	mycatcher := fakenet.Interceptor{}.WithBody("hello").WithStatus(201).WithHeader("TEST", "true")
	client.Intercept(mycatcher.WithURLMatcher("http://example.org"))

	resp, err := client.Get("http://example.org")
	require.Nil(t, err)
	require.NotNil(t, resp)

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)

	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, "hello", string(content))
	assert.Contains(t, resp.Header["TEST"], "true")
}
