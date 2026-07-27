package main

import (
	"fmt"
	"math"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pol-cova/termkit-go/chart"
	"github.com/pol-cova/termkit-go/component"
)

type monitorModel struct {
	width, height, frame, selected int
	cpuHist                        *chart.History
	memHist                        *chart.History
	netHist                        *chart.History
	diskHist                       *chart.History
}

func runMonitor() {
	cpu, _ := chart.NewHistory(120, []string{"All"}, []string{"monitor-primary"})
	mem, _ := chart.NewHistory(120, []string{"RAM", "SWP"}, []string{"monitor-primary", "monitor-secondary"})
	net, _ := chart.NewHistory(120, []string{"RX", "TX"}, []string{"monitor-primary", "monitor-secondary"})
	disk, _ := chart.NewHistory(120, []string{"read", "write"}, []string{"monitor-primary", "monitor-secondary"})

	m := monitorModel{
		width: 100, height: 32, cpuHist: cpu, memHist: mem, netHist: net, diskHist: disk,
	}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (m monitorModel) Init() tea.Cmd { return monitorTick() }

func (m monitorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left", "h":
			m.selected = max(0, m.selected-1)
		case "right", "l":
			m.selected = min(119, m.selected+1)
		}
	case tea.WindowSizeMsg:
		m.width = max(60, msg.Width)
		m.height = max(20, msg.Height)
	case monitorFrame:
		m.frame++
		t := float64(m.frame) * 0.08
		m.cpuHist.Push(15 + 25*math.Sin(t) + 10*math.Sin(t*2.3))
		m.memHist.Push(65+8*math.Sin(t*0.7), 20+5*math.Sin(t*1.1))
		m.netHist.Push(200+150*math.Abs(math.Sin(t*1.4)), 400+200*math.Abs(math.Cos(t*0.9)))
		m.diskHist.Push(50+40*math.Abs(math.Sin(t*0.5)), 30+25*math.Abs(math.Cos(t*0.6)))
		return m, monitorTick()
	}
	return m, nil
}

func (m monitorModel) View() string {
	rows, cols := 3, 2
	gridH := m.height - 2
	gridW := m.width

	cpuWidget := brailleWidget("CPU", m.cpuHist, 60*time.Second, chart.AxisSpec{Min: 0, Max: 100, Formatter: chart.PercentFormatter}, m.selected)
	memWidget := brailleWidget("Memory", m.memHist, 60*time.Second, chart.AxisSpec{Min: 0, Max: 100, Formatter: chart.PercentFormatter}, m.selected, memLegend...)
	netWidget := brailleWidget("Network", m.netHist, 60*time.Second, chart.AxisSpec{Formatter: chart.SIFormatter}, m.selected)
	diskWidget := brailleWidget("Disk", m.diskHist, 60*time.Second, chart.AxisSpec{Formatter: chart.SIFormatter}, m.selected)
	procWidget := component.FuncWidget{
		Render: func(w, h int) string {
			return component.BorderTitle("Processes", processTable(w-2, h-2, m.frame), w, h, component.Accent)
		},
	}
	coreWidget := component.FuncWidget{
		Render: func(w, h int) string {
			return component.BorderTitle("CPU", coreTable(w-2, h-2, m.frame), w, h, component.Accent)
		},
	}

	return component.Grid(gridW, gridH, rows, cols,
		component.Slot{Widget: cpuWidget, Row: 0, Col: 0, RowSpan: 2, WeightW: 3, WeightH: 2},
		component.Slot{Widget: coreWidget, Row: 0, Col: 1, WeightW: 1},
		component.Slot{Widget: memWidget, Row: 1, Col: 0, WeightW: 3},
		component.Slot{Widget: procWidget, Row: 1, Col: 1, RowSpan: 2, WeightW: 1, WeightH: 2},
		component.Slot{Widget: netWidget, Row: 2, Col: 0, WeightW: 3},
		component.Slot{Widget: diskWidget, Row: 2, Col: 1, WeightW: 1},
	) + "\n" + component.StatusBar(component.StatusBarOptions{
		Left: []component.Segment{
			{Text: "monitor demo", Tone: component.Accent},
			{Text: "←/→ scrub"},
			{Text: "q quit"},
		},
		Right: "synthetic data",
	})
}

type monitorFrame struct{}

func monitorTick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg { return monitorFrame{} })
}

var memLegend = []chart.LegendEntry{
	{Name: "RAM", Current: "69%", Extra: "25.0GiB/36.0GiB"},
	{Name: "SWP", Current: "25%", Extra: "0.5GiB/2.0GiB"},
}

func brailleWidget(title string, hist *chart.History, window time.Duration, y chart.AxisSpec, selected int, legendOverride ...chart.LegendEntry) component.FuncWidget {
	return component.FuncWidget{
		MinW: 24, MinH: 6,
		Render: func(w, h int) string {
			series := hist.Series()
			legend := chart.LegendSpec{}
			if len(legendOverride) > 0 {
				legend.Entries = updateLegendValues(legendOverride, series, y.Formatter)
			}
			ts := chart.TimeSeries{Title: title, Window: window, Series: series, Y: y, Legend: legend}
			body, err := chart.RenderBraille(ts, chart.BrailleOptions{
				Width: w - 2, Height: h - 2, Color: true, ShowAxes: true, ShowLegend: true, Selected: selected,
			})
			if err != nil {
				return component.BorderTitle(title, err.Error(), w, h, component.Danger)
			}
			return component.BorderTitle(title, body, w, h, component.Accent)
		},
	}
}

func updateLegendValues(entries []chart.LegendEntry, series []chart.Series, formatter chart.ValueFormatter) []chart.LegendEntry {
	if formatter == nil {
		formatter = chart.DefaultFormatter
	}
	out := make([]chart.LegendEntry, len(entries))
	copy(out, entries)
	for i := range out {
		if i < len(series) && len(series[i].Values) > 0 {
			out[i].Current = formatter(series[i].Values[len(series[i].Values)-1])
		}
	}
	return out
}

func processTable(w, h int, frame int) string {
	headers := []string{"PID", "Name", "CPU%", "Mem%"}
	rows := [][]string{
		{"384", "WindowServer", fmt.Sprintf("%d", 12+frame%5), "2.1"},
		{"991", "Cursor", fmt.Sprintf("%d", 8+frame%3), "4.5"},
		{"1204", "chrome", fmt.Sprintf("%d", 5+frame%4), "3.2"},
		{"2201", "node", fmt.Sprintf("%d", 3+frame%2), "1.8"},
		{"3300", "go", fmt.Sprintf("%d", 2+frame%3), "0.9"},
	}
	return component.Table(headers, rows, w)
}

func coreTable(w, h int, frame int) string {
	headers := []string{"Core", "Use"}
	rows := [][]string{
		{"All", fmt.Sprintf("%d%%", 15+frame%20)},
		{"CPU0", fmt.Sprintf("%d%%", frame%30)},
		{"CPU1", fmt.Sprintf("%d%%", (frame+3)%25)},
		{"CPU2", fmt.Sprintf("%d%%", (frame+7)%18)},
		{"CPU3", fmt.Sprintf("%d%%", (frame+2)%22)},
	}
	_ = h
	return component.Table(headers, rows, w)
}
