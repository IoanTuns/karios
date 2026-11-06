package calcul

import (
	"testing"
	"time"
)

func TestCalendarDistingeSarbatoareWeekendSiZiLucratoare(t *testing.T) {
	c := Calendar{Sarbatori: map[string]string{"2026-10-30": "Test vineri", "2026-10-31": "Test sambata"}}
	for _, v := range []struct {
		data string
		tip  TipZi
	}{{"2026-10-29", Lucratoare}, {"2026-10-30", Sarbatoare}, {"2026-10-31", Sarbatoare}, {"2026-11-01", Duminica}, {"2026-11-02", Lucratoare}} {
		d, _ := time.Parse(time.DateOnly, v.data)
		if got := c.TipZi(d); got != v.tip {
			t.Fatal(v.data, got, v.tip)
		}
	}
}
