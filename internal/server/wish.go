package server

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"

	"github.com/rytsh/yap/internal/tui"
	"github.com/rytsh/yap/internal/tui/model"
)

// You can write your own custom bubbletea middleware that wraps tea.Program.
// Make sure you set the program input and output to ssh.Session.
func screenMiddleware(screen tui.Screen) wish.Middleware {
	newProg := func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		p := tea.NewProgram(m, opts...)
		go func() {
			for {
				<-time.After(1 * time.Second)
				p.Send(model.TimeMsg(time.Now()))
			}
		}()
		return p
	}

	teaHandler := func(s ssh.Session) *tea.Program {
		pty, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}

		m := model.IndexModel{
			Width:  pty.Window.Width,
			Height: pty.Window.Height,

			Models: screen.Models(),
		}

		return newProg(m.SetModels(), tea.WithInput(s), tea.WithOutput(s), tea.WithAltScreen())
	}

	return bm.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}
