package style

import "testing"

func TestMonitorSeriesColor(t *testing.T) {
	if MonitorSeriesColor(0, false) != 208 {
		t.Fatalf("primary = %d", MonitorSeriesColor(0, false))
	}
	if MonitorSeriesColor(1, false) != 66 {
		t.Fatalf("secondary = %d", MonitorSeriesColor(1, false))
	}
}

func TestMonitorResolve(t *testing.T) {
	c := Resolve("monitor-primary", 0)
	if c != MonitorPrimary {
		t.Fatal(c)
	}
	if c.ANSI256(false) != 208 {
		t.Fatal(c.ANSI256(false))
	}
}
