package view

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/rytsh/yap/internal/server/view/style"
)

// Just a generic tea.ModelGen to demo terminal information of ssh.
type ModelGen struct {
	Term   string
	Width  int
	Height int
	Time   time.Time

	// altScreen is true when the terminal is in the alternate screen buffer.
	altScreen bool
}

type TimeMsg time.Time

func (m ModelGen) Init() tea.Cmd {
	// return tea.Batch(
	// 	tea.ExitAltScreen,
	// )

	return nil
}

func (m ModelGen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TimeMsg:
		m.Time = time.Time(msg)
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "alt+m":
			m.altScreen = !m.altScreen
			if m.altScreen {
				return m, tea.EnterAltScreen
			} else {
				return m, tea.ExitAltScreen
			}
		}
	}
	return m, nil
}

func (m ModelGen) View() string {
	// style := lipgloss.NewStyle().Bold(true).SetString("Hello,")
	// s := fmt.Sprintf("%s\n%s\n", style.Render("kitty."), style.Render("puppy."))

	// return s
	// s := "Your term is %s\n"
	// s += "Your window size is x: %d y: %d\n"
	// s += "Time: " + m.Time.Format(time.RFC1123) + "\n\n"
	// s += "Press 'q' to quit\n"
	// return fmt.Sprintf(s, m.Term, m.Width, m.Height)

	physicalWidth := m.Width
	doc := strings.Builder{}
	// Tabs
	{
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			style.ActiveTab.Render("Lip Gloss"),
			style.Tab.Render("Blush"),
			style.Tab.Render("Eye Shadow"),
			style.Tab.Render("Mascara"),
			style.Tab.Render("Foundation"),
		)
		gap := style.TabGap.Render(strings.Repeat(" ", style.Max(0, style.Width-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row + "\n\n")
	}

	// Title
	{
		var (
			colors = style.ColorGrid(1, 5)
			title  strings.Builder
		)

		for i, v := range colors {
			const offset = 2
			c := lipgloss.Color(v[0])
			fmt.Fprint(&title, style.TitleStyle.Copy().MarginLeft(i*offset).Background(c))
			if i < len(colors)-1 {
				title.WriteRune('\n')
			}
		}

		desc := lipgloss.JoinVertical(lipgloss.Left,
			style.DescStyle.Render("Style Definitions for Nice Terminal Layouts"),
			style.InfoStyle.Render("From Charm"+style.Divider+style.Url("https://github.com/charmbracelet/lipgloss")),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, title.String(), desc)
		doc.WriteString(row + "\n\n")
	}

	// Dialog
	{
		okButton := style.ActiveButtonStyle.Render("Yes")
		cancelButton := style.ButtonStyle.Render("Maybe")

		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

		dialog := lipgloss.Place(style.Width, 9,
			lipgloss.Center, lipgloss.Center,
			style.DialogBoxStyle.Render(ui),
			lipgloss.WithWhitespaceChars("猫咪"),
			lipgloss.WithWhitespaceForeground(style.Subtle),
		)

		doc.WriteString(dialog + "\n\n")
	}

	// Color grid
	colors := func() string {
		colors := style.ColorGrid(14, 8)

		b := strings.Builder{}
		for _, x := range colors {
			for _, y := range x {
				s := lipgloss.NewStyle().SetString("  ").Background(lipgloss.Color(y))
				b.WriteString(s.String())
			}
			b.WriteRune('\n')
		}

		return b.String()
	}()

	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		style.List.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				style.ListHeader("Citrus Fruits to Try"),
				style.ListDone("Grapefruit"),
				style.ListDone("Yuzu"),
				style.ListItem("Citron"),
				style.ListItem("Kumquat"),
				style.ListItem("Pomelo"),
			),
		),
		style.List.Copy().Width(style.ColumnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				style.ListHeader("Actual Lip Gloss Vendors"),
				style.ListItem("Glossier"),
				style.ListItem("Claire‘s Boutique"),
				style.ListDone("Nyx"),
				style.ListItem("Mac"),
				style.ListDone("Milk"),
			),
		),
	)

	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists, colors))

	// Marmalade history
	{
		const (
			historyA = "The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum: Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos."
			historyB = "Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac."
			historyC = "In 1524, Henry VIII, King of England, received a “box of marmalade” from Mr. Hull of Exeter. This was probably marmelada, a solid quince paste from Portugal, still made and sold in southern Europe today. It became a favourite treat of Anne Boleyn and her ladies in waiting."
		)

		doc.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			style.HistoryStyle.Copy().Align(lipgloss.Right).Render(historyA),
			style.HistoryStyle.Copy().Align(lipgloss.Center).Render(historyB),
			style.HistoryStyle.Copy().MarginRight(0).Render(historyC),
		))

		doc.WriteString("\n\n")
	}

	// Status bar
	{
		w := lipgloss.Width

		statusKey := style.StatusStyle.Render("STATUS")
		encoding := style.EncodingStyle.Render("UTF-8")
		fishCake := style.FishCakeStyle.Render("🍥 Fish Cake")
		statusVal := style.StatusText.Copy().
			Width(style.Width - w(statusKey) - w(encoding) - w(fishCake)).
			Render("Ravishing")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			statusKey,
			statusVal,
			encoding,
			fishCake,
		)

		doc.WriteString(style.StatusBarStyle.Width(style.Width).Render(bar))
	}

	if physicalWidth > 0 {
		style.DocStyle = style.DocStyle.MaxWidth(physicalWidth).MaxHeight(m.Height)
	}

	// Okay, let's print it
	return fmt.Sprintln(style.DocStyle.Render(doc.String()))
}
