// Package remote .
package remote

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type eventer interface {
	SetKeys(keys ...int)
	Press() error
	Release() error
}

type config struct {
	KeyPressDelay Duration `json:"key_press_delay"`
	Buttons       []button `json:"buttons"`
}

type button struct {
	ID        string `json:"id"`
	ColorCode string `json:"color_code"`
	IconClass string `json:"icon_class"`

	Keycode int `json:"keycode"`
}

var (
	// ErrInvalidButton .
	ErrInvalidButton = errors.New("invalid button")
)

type Remote struct {
	configFile string
	keyboard   eventer

	keyMapLock    sync.Mutex
	keyMap        map[string]int
	keyPressDelay Duration
}

// New .
func New(configFile string, keyboard eventer) (*Remote, error) {
	s := &Remote{
		configFile: configFile,
		keyboard:   keyboard,
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

	events := make(map[string]int)
	for _, b := range conf.Buttons {
		events[b.ID] = b.Keycode
	}

	r.keyMapLock.Lock()
	r.keyMap = events
	if conf.KeyPressDelay.Duration == 0 {
		conf.KeyPressDelay = Duration{10 * time.Millisecond}
	}
	r.keyPressDelay = conf.KeyPressDelay
	r.keyMapLock.Unlock()

	log.Println(r.keyPressDelay)

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

	err := r.keyboard.Press()
	if err != nil {
		return err
	}
	time.Sleep(r.keyPressDelay.Duration)
	return r.keyboard.Release()
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
