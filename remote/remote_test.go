package remote

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_NewLoadConfig(t *testing.T) {
	s, err := New("../config/config_test.json", nil)
	require.NoError(t, err)

	assert.Len(t, s.keyMap, 2)
	assert.Equal(t, 0, s.keyMap["1"])
	assert.Equal(t, 1, s.keyMap["2"])
}

func TestServer_LoadConfig(t *testing.T) {
	s := &Remote{
		configFile: "../config/config_test.json",
	}

	assert.Len(t, s.keyMap, 0)

	assert.NoError(t, s.LoadConfig())

	assert.Len(t, s.keyMap, 2)
	assert.Equal(t, 0, s.keyMap["1"])
	assert.Equal(t, 1, s.keyMap["2"])
}

func TestServer_LoadConfigInvalidFile(t *testing.T) {
	s := &Remote{
		configFile: "fake_news.json",
	}
	assert.Error(t, s.LoadConfig())
}

func TestKeyPress_InvalidButton(t *testing.T) {
	s := &Remote{keyMap: map[string]int{
		"1": 42,
	}}

	assert.Equal(t, ErrInvalidButton, s.KeyPress("2"))
}

func TestKeyPress_Success(t *testing.T) {
	testError := errors.New("testError")
	te := &testEvent{
		retPressError: testError,
	}

	s := &Remote{
		keyMap: map[string]int{
			"1": 42,
		},
		keyboard: te,
	}

	assert.Equal(t, testError, s.KeyPress("1"))
	require.Len(t, te.paramKeys, 1)
	assert.Equal(t, 42, te.paramKeys[0])
}

type testEvent struct {
	paramKeys       []int
	retPressError   error
	retReleaseError error
}

func (te *testEvent) SetKeys(keys ...int) {
	te.paramKeys = keys
}

func (te *testEvent) Press() error {
	return te.retPressError
}

func (te *testEvent) Release() error {
	return te.retReleaseError
}
