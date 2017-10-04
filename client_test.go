package fakenet_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/apourchet/fakenet"
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
	require.Equal(t, "OK", string(content))
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
	require.Equal(t, "OK", string(content))

	resp, err = client.Get("http://unknown.org")
	require.NotNil(t, err)
}
