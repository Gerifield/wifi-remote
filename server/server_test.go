package server

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer_RootHandler(t *testing.T) {
	s := &Server{}

	srv := httptest.NewServer(s.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL+"/")
	require.NoError(t, err)
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, []byte("Hello"), b)
}

func TestServer_NewLoadConfig(t *testing.T) {
	s, err := New("../config/config_test.json")
	require.NoError(t, err)

	assert.Len(t, s.keyMap, 2)
	assert.Equal(t, 0, s.keyMap["1"])
	assert.Equal(t, 1, s.keyMap["2"])
}

func TestServer_LoadConfig(t *testing.T) {
	s := &Server{
		configFile: "../config/config_test.json",
	}

	assert.Len(t, s.keyMap, 0)

	assert.NoError(t, s.LoadConfig())

	assert.Len(t, s.keyMap, 2)
	assert.Equal(t, 0, s.keyMap["1"])
	assert.Equal(t, 1, s.keyMap["2"])
}
