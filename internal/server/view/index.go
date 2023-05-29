package view

import tea "github.com/charmbracelet/bubbletea"

type Model interface {
	tea.Model
	SetIndex(Index)
	Initialize(Config) tea.Cmd
}

type Index interface {
	InitModel(Config) (tea.Model, tea.Cmd)
	PrevModel(Config) (tea.Model, tea.Cmd)
	NextModel(Config) (tea.Model, tea.Cmd)
}

type Config struct {
	Width  int
	Height int
}

type IndexModel struct {
	Width  int
	Height int

	Models     []Model
	ModelIndex int
}

func (m *IndexModel) SetModels() tea.Model {
	for i := range m.Models {
		m.Models[i].SetIndex(m)
	}

	return m.Models[0]
}

func (m *IndexModel) getModel(index int, cfg Config) (tea.Model, tea.Cmd) {
	if index < 0 || index >= len(m.Models) {
		return nil, nil
	}

	return m.Models[index], m.Models[index].Initialize(cfg)
}

func (m *IndexModel) InitModel(cfg Config) (tea.Model, tea.Cmd) {
	m.ModelIndex = 0

	return m.getModel(m.ModelIndex, cfg)
}

func (m *IndexModel) PrevModel(cfg Config) (tea.Model, tea.Cmd) {
	if m.ModelIndex != 0 {
		m.ModelIndex--
	}

	return m.getModel(m.ModelIndex, cfg)
}

func (m *IndexModel) NextModel(cfg Config) (tea.Model, tea.Cmd) {
	m.ModelIndex++

	return m.getModel(m.ModelIndex, cfg)
}
