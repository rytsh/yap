package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

func Tabs(tabs []string, selected string, width int) string {
	var renderTabs []string
	for i, tab := range tabs {
		if tab == selected {
			if i == 0 {
				renderTabs = append(renderTabs, TabActiveFirst.Render(tab))
			} else if i == len(tabs)-1 {
				renderTabs = append(renderTabs, TabActiveLast.Render(tab))
			} else {
				renderTabs = append(renderTabs, TabActive.Render(tab))
			}
		} else {
			if i == 0 {
				renderTabs = append(renderTabs, TabNormalFirst.Render(tab))
			} else if i == len(tabs)-1 {
				renderTabs = append(renderTabs, TabNormalLast.Render(tab))
			} else {
				renderTabs = append(renderTabs, TabNormal.Render(tab))
			}
		}
	}

	if len(renderTabs) == 0 {
		renderTabs = append(renderTabs, Tab.Render("nothing"))
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderTabs...,
	)

	gap := TabGap.Render(strings.Repeat(" ", Max(0, width-lipgloss.Width(row)-1)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
}

func SwitchTab(tabs []string, current string, IsNext bool) string {
	if len(tabs) < 2 {
		return current
	}

	var index int
	for i, tab := range tabs {
		if tab == current {
			index = i
			break
		}
	}

	if IsNext {
		index++
		if index >= len(tabs) {
			index = 0
		}
	} else {
		index--
		if index < 0 {
			index = len(tabs) - 1
		}
	}

	return tabs[index]
}

func ColorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := 0; y < xSteps; y++ {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
