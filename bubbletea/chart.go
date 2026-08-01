package bubbletea

import "github.com/pol-cova/termkit-go/chart"

type ChartView struct {
	Data    chart.Chart
	Options chart.Options
}

func (v ChartView) View() string {
	out, err := chart.Render(v.Data, v.Options)
	if err != nil {
		return err.Error()
	}
	return out
}

func (v ChartView) Selected() chart.SelectedPoint {
	return v.Data.SelectedPoint(v.Options)
}

func ScrubSelection(key string, selected *int, count int) bool {
	if count <= 0 || selected == nil {
		return false
	}
	switch key {
	case "left", "h":
		if *selected > 0 {
			*selected--
			return true
		}
	case "right", "l":
		if *selected < count-1 {
			*selected++
			return true
		}
	case "home":
		if *selected != 0 {
			*selected = 0
			return true
		}
	case "end":
		last := count - 1
		if *selected != last {
			*selected = last
			return true
		}
	}
	return false
}

func CycleSeries(active *int, count int) bool {
	if count <= 0 || active == nil {
		return false
	}
	*active = (*active + 1) % count
	return true
}
