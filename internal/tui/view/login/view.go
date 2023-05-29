package login

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rytsh/yap/internal/tui/model"
	"github.com/rytsh/yap/internal/tui/style"
)

type LoginModel struct {
	width      int
	height     int
	keymap     keymapLogin
	innerWidth int
	help       help.Model
	inputs     []textinput.Model
	focusIndex int
	focusMax   int
	err        error
	time       time.Time

	index model.Index

	action      Action
	tabs        []string
	selectedTab string
}

type keymapLogin = struct {
	next, prev, nextTab, prevTab, login, quit, selection key.Binding
}

func NewLoginModel(action Action) *LoginModel {
	m := LoginModel{
		action:      action,
		tabs:        action.GetTabNames(),
		selectedTab: action.TabSelected(),
		inputs:      make([]textinput.Model, 2),
		innerWidth:  50,
		help:        help.New(),
		keymap: keymapLogin{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			login: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "login"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
			selection: key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "select"),
			),
			nextTab: key.NewBinding(
				key.WithKeys("ctrl+right"),
				key.WithHelp("ctrl+right", "next tab"),
			),
			prevTab: key.NewBinding(
				key.WithKeys("ctrl+left"),
				key.WithHelp("ctrl+left", "prev tab"),
			),
		},
	}

	for i := range m.inputs {
		t := textinput.New()

		switch i {
		case 0:
			// t.Placeholder = "username"
		case 1:
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		}

		m.inputs[i] = t
	}

	m.focusMax = len(m.inputs) + 2
	m.focusIndex = 0

	return &m
}

func (m *LoginModel) SetIndex(index model.Index) {
	m.index = index
}

func (m *LoginModel) Initialize(cfg model.Config) tea.Cmd {
	m.width = cfg.Width
	m.height = cfg.Height

	return m.Init()
}

func (m LoginModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.updateFocus())
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// fmt.Println("[" + msg.String() + "]")
		switch {
		case key.Matches(msg, m.keymap.quit):
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, tea.Quit
		case key.Matches(msg, m.keymap.next):
			m.focusIndex++
			if m.focusIndex > m.focusMax-1 {
				m.focusIndex = 0
			}
			return m, m.updateFocus()
		case key.Matches(msg, m.keymap.prev):
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = m.focusMax - 1
			}
			return m, m.updateFocus()
		case key.Matches(msg, m.keymap.login):
			return m.login()
		case key.Matches(msg, m.keymap.selection):
			// fmt.Printf("selection pressed %d\n", m.focusIndex)
			switch m.focusIndex {
			case len(m.inputs):
				// submit
				return m.login()
			case len(m.inputs) + 1:
				// quit
				return m, tea.Quit
			}
		case key.Matches(msg, m.keymap.nextTab):
			m.selectedTab = style.SwitchTab(m.tabs, m.selectedTab, true)
			return m, nil
		case key.Matches(msg, m.keymap.prevTab):
			m.selectedTab = style.SwitchTab(m.tabs, m.selectedTab, false)
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case model.TimeMsg:
		m.time = time.Time(msg)
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}
func (m *LoginModel) login() (tea.Model, tea.Cmd) {
	if err := m.action.Login(m.selectedTab, m.inputs[0].Value(), m.inputs[1].Value()); err != nil {
		m.err = model.TimeErr(err, m.time)

		return m, nil
	}

	return m.index.NextModel(model.Config{
		Width:  m.width,
		Height: m.height,
	})
}

func (m *LoginModel) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Set focused state for inputs, starts 0
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focusIndex {
			// Set focused state
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = style.FocusedStyle
			m.inputs[i].TextStyle = style.FocusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = style.NoStyle
		m.inputs[i].TextStyle = style.NoStyle
	}

	return tea.Batch(cmds...)
}

func (m *LoginModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m LoginModel) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		m.keymap.nextTab,
		m.keymap.prevTab,
		m.keymap.login,
		m.keymap.selection,
		m.keymap.quit,
	})

	var b strings.Builder

	for i := range m.inputs {
		switch i {
		case 0:
			b.WriteString("Username\n")
		case 1:
			b.WriteString("Password\n")
		}

		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	// box shape
	var submitButton string
	if m.focusIndex == len(m.inputs) {
		submitButton = style.ActiveButtonStyle.MarginRight(2).Render("Submit")
	} else {
		submitButton = style.ButtonStyle.MarginRight(2).Render("Submit")
	}

	var cancelButton string
	if m.focusIndex == len(m.inputs)+1 {
		cancelButton = style.ActiveButtonStyle.Render("Cancel")
	} else {
		cancelButton = style.ButtonStyle.Render("Cancel")
	}

	question := lipgloss.NewStyle().Width(m.innerWidth).Align(lipgloss.Left).Padding(0, 1).Render(b.String())
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, submitButton, cancelButton)

	uiVertical := []string{question, buttons}

	errStr := ""
	if m.err != nil {
		errStr = m.err.Error()
	}

	uiVertical = append(uiVertical, style.ErrorStyle.Width(m.innerWidth).Render(errStr))

	ui := lipgloss.JoinVertical(lipgloss.Center, uiVertical...)

	tabs := style.Tabs(m.tabs, m.selectedTab, m.innerWidth)

	banner := ""
	if m.action.Banner != "" {
		banner = bannerStyle.Render(m.action.Banner + strings.Repeat(" ", style.Max(0, m.innerWidth-lipgloss.Width(m.action.Banner)-1)))
	}

	dialog := lipgloss.Place(m.width, 0,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left, banner, tabs, boxStyle.Render(ui)),
		lipgloss.WithWhitespaceChars("  ##  "),
		lipgloss.WithWhitespaceForeground(style.Subtle),
	)

	return dialog + "\n\n" + help
}
