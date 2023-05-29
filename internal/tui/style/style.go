package style

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	// In real life situations we'd adjust the document to fit the Width we've
	// detected. In the case of this example we're hardcoding the Width, and
	// later using the detected Width only to truncate in order to avoid jaggy
	// wrapping.
	Width = 96

	ColumnWidth = 30
)

// Style definitions.
var (

	// General.

	Subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	Highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	Special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	Divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(Subtle).
		String()

	Url = lipgloss.NewStyle().Foreground(Special).Render

	// Tabs

	BorderActiveTab = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	BorderActiveTabFirst = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "│",
		BottomRight: "└",
	}

	BorderActiveTabLast = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	BorderNormalTab = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	BorderNormalTabFirst = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "├",
		BottomRight: "┴",
	}

	BorderNormalTabLast = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	BorderTabGap = lipgloss.Border{
		Bottom:      "─",
		BottomRight: "┐",
	}

	Tab = lipgloss.NewStyle().
		Border(BorderNormalTab, true).
		BorderForeground(Highlight).
		Padding(0, 1)

	TabNormal      = Tab.Copy().Border(BorderNormalTab, true)
	TabNormalFirst = Tab.Copy().Border(BorderNormalTabFirst, true)
	TabNormalLast  = Tab.Copy().Border(BorderNormalTabLast, true)

	TabActive      = Tab.Copy().Border(BorderActiveTab, true)
	TabActiveFirst = Tab.Copy().Border(BorderActiveTabFirst, true)
	TabActiveLast  = Tab.Copy().Border(BorderActiveTabLast, true)

	TabGap = Tab.Copy().
		Border(BorderTabGap, true).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(true)

	// Title.

	TitleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(5).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			SetString("Lip Gloss")

	DescStyle = lipgloss.NewStyle().MarginTop(1)

	InfoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(Subtle)

	// Dialog.

	DialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0, 0, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	ActiveButtonStyle = ButtonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				Underline(true)

	// List.

	List = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(Subtle).
		MarginRight(2).
		Height(8).
		Width(ColumnWidth + 1)

	ListHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(Subtle).
			MarginRight(2).
			Render

	ListItem = lipgloss.NewStyle().PaddingLeft(2).Render

	CheckMark = lipgloss.NewStyle().SetString("✓").
			Foreground(Special).
			PaddingRight(1).
			String()

	ListDone = func(s string) string {
		return CheckMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}

	// Paragraphs/History.

	HistoryStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(Highlight).
			Margin(1, 3, 0, 0).
			Padding(1, 2).
			Height(19).
			Width(ColumnWidth)

	// Status Bar.

	StatusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	StatusStyle = lipgloss.NewStyle().
			Inherit(StatusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	EncodingStyle = StatusNugget.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	StatusText = lipgloss.NewStyle().Inherit(StatusBarStyle)

	FishCakeStyle = StatusNugget.Copy().Background(lipgloss.Color("#6124DF"))

	// Page.

	DocStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle  = FocusedStyle.Copy()
	NoStyle      = lipgloss.NewStyle()
	HelpStyle    = BlurredStyle.Copy()

	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(0)).
			Background(lipgloss.ANSIColor(3)).
			MarginTop(1)

	TabStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0, 0, 0)
)
