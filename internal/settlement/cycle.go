package settlement

import "time"

// CycleAndDate returns the settlement cycle (1–3) and date for a given UTC time.
//
// FPS settlement windows:
//   - 00:00–07:15 → cycle 1, today
//   - 07:15–13:00 → cycle 2, today
//   - 13:00–15:45 → cycle 3, today
//   - 15:45–24:00 → cycle 1, next calendar day
func CycleAndDate(t time.Time) (cycle int, date string) {
	t = t.UTC()
	minutes := t.Hour()*60 + t.Minute()

	switch {
	case minutes < 7*60+15:
		return 1, t.Format(time.DateOnly)
	case minutes < 13*60:
		return 2, t.Format(time.DateOnly)
	case minutes < 15*60+45:
		return 3, t.Format(time.DateOnly)
	default:
		return 1, t.AddDate(0, 0, 1).Format(time.DateOnly)
	}
}
