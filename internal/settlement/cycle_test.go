package settlement

import (
	"testing"
	"time"
)

func TestCycleAndDate(t *testing.T) {
	tests := []struct {
		name      string
		utcTime   string // RFC3339
		wantCycle int
		wantDate  string
	}{
		{"midnight", "2026-03-18T00:00:00Z", 1, "2026-03-18"},
		{"early morning", "2026-03-18T06:30:00Z", 1, "2026-03-18"},
		{"just before cycle 2", "2026-03-18T07:14:59Z", 1, "2026-03-18"},
		{"cycle 2 start", "2026-03-18T07:15:00Z", 2, "2026-03-18"},
		{"midday", "2026-03-18T12:00:00Z", 2, "2026-03-18"},
		{"just before cycle 3", "2026-03-18T12:59:59Z", 2, "2026-03-18"},
		{"cycle 3 start", "2026-03-18T13:00:00Z", 3, "2026-03-18"},
		{"afternoon", "2026-03-18T15:00:00Z", 3, "2026-03-18"},
		{"just before next-day rollover", "2026-03-18T15:44:59Z", 3, "2026-03-18"},
		{"next-day rollover", "2026-03-18T15:45:00Z", 1, "2026-03-19"},
		{"evening", "2026-03-18T20:00:00Z", 1, "2026-03-19"},
		{"just before midnight", "2026-03-18T23:59:59Z", 1, "2026-03-19"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm, err := time.Parse(time.RFC3339, tt.utcTime)
			if err != nil {
				t.Fatalf("parse time: %v", err)
			}
			cycle, date := CycleAndDate(tm)
			if cycle != tt.wantCycle {
				t.Errorf("cycle = %d, want %d", cycle, tt.wantCycle)
			}
			if date != tt.wantDate {
				t.Errorf("date = %s, want %s", date, tt.wantDate)
			}
		})
	}
}
