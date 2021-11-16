// Package remote .
package remote

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
)

type eventer interface {
	SetKeys(keys ...int)
	Launching() error
}

type config struct {
	Events map[string]int `json:"events"`
}

var (
	// ErrInvalidButton .
	ErrInvalidButton = errors.New("invalid button")
)

type Remote struct {
	configFile string
	keyboard eventer

	keyMapLock sync.Mutex
	keyMap map[string]int
}


// New .
func New(configFile string, keyboard eventer) (*Remote, error) {
	s := &Remote{
		configFile: configFile,
		keyboard: keyboard,
		keyMap: make(map[string]int),
	}

	return s, s.LoadConfig()
}

// LoadConfig .
func (r *Remote) LoadConfig() error {
	b, err := ioutil.ReadFile(r.configFile)
	if err != nil {
		return err
	}

	var conf config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		return err
	}

	r.keyMapLock.Lock()
	r.keyMap = conf.Events
	r.keyMapLock.Unlock()

	return nil
}

// KeyPress .
func (r *Remote) KeyPress(button string) error {
	r.keyMapLock.Lock()
	key, ok := r.keyMap[button]
	r.keyMapLock.Unlock()

	if !ok {
		return ErrInvalidButton
	}

	r.keyboard.SetKeys(key)

	return r.keyboard.Launching()
}
