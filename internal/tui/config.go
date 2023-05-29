package tui

import (
	"github.com/rytsh/yap/internal/tui/model"
	"github.com/rytsh/yap/internal/tui/view/login"
)

type Screen []View

type View struct {
	ID   string `cfg:"id"`
	Next string `cfg:"next"`
	Prev string `cfg:"prev"`
	Up   string `cfg:"up"`
	Down string `cfg:"down"`

	Selection Selection `cfg:"selection"`
}

type Selection struct {
	Login *login.Action `cfg:"login"`
}

func (s Selection) Action() model.Model {
	if s.Login != nil {
		return login.NewLoginModel(*s.Login)
	}

	return nil
}

func (s Screen) Models() []model.Model {
	var models []model.Model

	for _, v := range s {
		models = append(models, v.Selection.Action())
	}

	return models
}
