package style

// Monitor palette tokens for btm-style time-series graphs.
const (
	MonitorPrimary   Color = "monitor-primary"
	MonitorSecondary Color = "monitor-secondary"
	MonitorGrid      Color = "monitor-grid"
	MonitorBorder    Color = "monitor-border"
)

func init() {
	rgb[MonitorPrimary] = [3]int{215, 153, 33}   // ≈ #D79921
	rgb[MonitorSecondary] = [3]int{69, 133, 136} // ≈ #458588
	rgb[MonitorGrid] = [3]int{80, 80, 80}
	rgb[MonitorBorder] = [3]int{180, 180, 180}
	ansi256[MonitorPrimary] = 208
	ansi256[MonitorSecondary] = 66
	ansi256[MonitorGrid] = 240
	ansi256[MonitorBorder] = 245
}

// MonitorSeriesColor returns the ANSI-256 colour for a monitor graph series index.
func MonitorSeriesColor(index int, bright bool) int {
	order := []Color{MonitorPrimary, MonitorSecondary, Green, Blue, Pink, Orange}
	return Resolve(string(order[index%len(order)]), index).ANSI256(bright)
}
