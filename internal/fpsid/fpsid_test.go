package fpsid

import (
	"regexp"
	"testing"
	"time"
)

func TestNumericReference(t *testing.T) {
	re := regexp.MustCompile(`^\d{6}$`)
	for range 20 {
		got := NumericReference()
		if !re.MatchString(got) {
			t.Errorf("NumericReference() = %q, want 6 digits", got)
		}
	}
}

func TestEndToEndReference(t *testing.T) {
	today := time.Now().UTC().Format("20060102")
	re := regexp.MustCompile(`^FPS` + today + `\d{6}$`)
	for range 20 {
		got := EndToEndReference()
		if !re.MatchString(got) {
			t.Errorf("EndToEndReference() = %q, want FPS+date+6digits", got)
		}
		if len(got) != 17 {
			t.Errorf("EndToEndReference() len = %d, want 17", len(got))
		}
	}
}

func TestSchemeTransactionID(t *testing.T) {
	re := regexp.MustCompile(`^\d{26}$`)
	for range 20 {
		got := SchemeTransactionID()
		if !re.MatchString(got) {
			t.Errorf("SchemeTransactionID() = %q, want 26 digits", got)
		}
	}
}

func TestProcessingDate(t *testing.T) {
	got := ProcessingDate()
	today := time.Now().UTC().Format(time.DateOnly)
	if got != today {
		t.Errorf("ProcessingDate() = %q, want %q", got, today)
	}
}
