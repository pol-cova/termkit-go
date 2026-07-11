package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/pol-cova/termkit-go/animate"
	"github.com/pol-cova/termkit-go/chart"
	"github.com/pol-cova/termkit-go/component"
)

type model struct {
	kind, selected, series int
	width, height          int
	frame                  int
}

type frame struct{}

var kinds = []chart.Kind{chart.Area, chart.Line, chart.Bar, chart.Pie, chart.Radar}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--static" {
		view, _ := chart.Render(sample(kinds[0]), chart.Options{Width: 58, Height: 13, Selected: 2, Color: true})
		fmt.Println(view)
		return
	}
	if _, err := tea.NewProgram(model{width: 68, height: 15}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd { return nextFrame() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1", "2", "3", "4", "5":
			m.kind = int(msg.String()[0] - '1')
		case "left", "h":
			m.selected = max(0, m.selected-1)
		case "right", "l":
			m.selected = min(5, m.selected+1)
		case "tab":
			m.series = (m.series + 1) % 2
		}
	case tea.WindowSizeMsg:
		m.width = max(44, msg.Width-4)
		m.height = max(7, msg.Height-13)
	case frame:
		m.frame++
		return m, nextFrame()
	}
	return m, nil
}

func (m model) View() string {
	c := sample(kinds[m.kind])
	c.Title = "termkit-go  /  " + string(c.Kind) + " chart"
	view, err := chart.Render(c, chart.Options{Width: m.width, Height: m.height, Selected: m.selected, ActiveSeries: m.series, Color: true})
	if err != nil {
		return err.Error()
	}
	motion := animate.Tween(float64(m.frame%40)/39, animate.EaseInOut)
	widgets := component.Badge("LIVE", component.Success) + "  " + component.Progress("CPU", motion, 14, component.Accent) + "   " + component.Gauge("memory", 1-motion*.45, 12, component.Success) + "   " + component.SpinnerFrame(m.frame, "sampling", component.Warning)
	footer := component.StatusBar(component.StatusBarOptions{
		Left:  []component.Segment{{Text: "1–5 chart", Tone: component.Accent}, {Text: "←/→ scrub"}, {Text: "tab series"}},
		Right: fmt.Sprintf("motion %.0f%%", motion*100),
	})
	return view + "\n\n" + widgets + "\n" + footer + "\n"
}

func sample(kind chart.Kind) chart.Chart {
	stack := chart.Default
	if kind == chart.Area || kind == chart.Bar {
		stack = chart.Stacked
	}
	return chart.Chart{
		Kind: kind, Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}, StackType: stack,
		Series: []chart.Series{
			{Name: "desktop", Values: []float64{18, 42, 29, 68, 53, 76}, Variant: chart.Dotted},
			{Name: "mobile", Values: []float64{25, 17, 56, 39, 72, 46}, Variant: chart.Hatched},
		},
	}
}
func nextFrame() tea.Cmd {
	return tea.Tick(75*time.Millisecond, func(time.Time) tea.Msg { return frame{} })
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
