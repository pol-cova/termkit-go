package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	ditherkit "github.com/pol-cova/ditherkit-go"
)

type model struct {
	kind, selected, series int
	width, height          int
}

var kinds = []ditherkit.Kind{ditherkit.Area, ditherkit.Line, ditherkit.Bar, ditherkit.Pie, ditherkit.Radar}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--static" {
		chart, _ := ditherkit.Render(sample(kinds[0]), ditherkit.Options{Width: 58, Height: 13, Selected: 2, Color: true})
		fmt.Println(chart)
		return
	}
	if _, err := tea.NewProgram(model{width: 68, height: 15}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1", "2", "3", "4", "5":
			m.kind = int(key.String()[0] - '1')
		case "left", "h":
			m.selected = max(0, m.selected-1)
		case "right", "l":
			m.selected = min(5, m.selected+1)
		case "tab":
			m.series = (m.series + 1) % 2
		}
	}
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = max(44, size.Width-4)
		m.height = max(8, size.Height-9)
	}
	return m, nil
}
func (m model) View() string {
	c := sample(kinds[m.kind])
	c.Title = "ditherkit-go  /  " + string(c.Kind) + " chart"
	chart, err := ditherkit.Render(c, ditherkit.Options{Width: m.width, Height: m.height, Selected: m.selected, ActiveSeries: m.series, Color: true})
	if err != nil {
		return err.Error()
	}
	return chart + "\n\n  1–5 chart type  •  ←/→ select point  •  tab series  •  q quit\n"
}
func sample(kind ditherkit.Kind) ditherkit.Chart {
	return ditherkit.Chart{Kind: kind, Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}, Series: []ditherkit.Series{{Name: "desktop", Values: []float64{18, 42, 29, 68, 53, 76}}, {Name: "mobile", Values: []float64{25, 17, 56, 39, 72, 46}}}}
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
