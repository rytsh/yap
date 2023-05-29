package login

import (
	"errors"
	"fmt"

	"github.com/rytsh/yap/internal/tui/view/login/auth"
)

var ErrLogin = errors.New("login failed")
var ErrNoSelection = errors.New("no selection")

type Action struct {
	Banner   string `cfg:"banner"`
	Tabs     []Auth `cfg:"tabs"`
	Selected string `cfg:"selected"`
}

func (a Action) GetTabNames() []string {
	v := make([]string, len(a.Tabs))
	for i, tab := range a.Tabs {
		v[i] = tab.Name
	}

	return v
}

func (a *Action) TabSelected() string {
	if a.Selected != "" {
		return a.Selected
	}

	if len(a.Tabs) > 0 {
		return a.Tabs[0].Name
	}

	return ""
}

func (a Action) Login(selected string, username, password string) error {
	for _, tab := range a.Tabs {
		if tab.Name == selected {
			return tab.Login(username, password)
		}
	}

	return ErrNoSelection
}

type Auth struct {
	Name string `cfg:"name"`

	BasicAuth *auth.BasicAuth `cfg:"basic_auth"`
}

func (a Auth) Login(username, password string) error {
	if a.BasicAuth != nil {
		if err := a.BasicAuth.Validator(username, password); err != nil {
			return fmt.Errorf("%w: %v", ErrLogin, err)
		}

		return nil
	}

	return ErrNoSelection
}
