package heartbeat

import (
	"testing"
	"time"
)

func TestParseInterval(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
	}{
		{"5m", 5 * time.Minute},
		{"1h", 1 * time.Hour},
		{"30s", 30 * time.Second},
		{"1d", 24 * time.Hour},
		{"7d", 7 * 24 * time.Hour},
		{"", 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseInterval(tt.input)
			if got != tt.want {
				t.Errorf("ParseInterval(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsDue(t *testing.T) {
	// Never run
	check := Check{Interval: 5 * time.Minute}
	if !IsDue(check) {
		t.Error("should be due when never run")
	}

	// Recently run
	check = Check{
		Interval: 5 * time.Minute,
		LastRun:  time.Now().Add(-1 * time.Minute),
	}
	if IsDue(check) {
		t.Error("should not be due when recently run")
	}

	// Overdue
	check = Check{
		Interval: 5 * time.Minute,
		LastRun:  time.Now().Add(-10 * time.Minute),
	}
	if !IsDue(check) {
		t.Error("should be due when overdue")
	}

	// No interval
	check = Check{}
	if IsDue(check) {
		t.Error("should not be due with no interval")
	}
}

func TestFormatLogEntry(t *testing.T) {
	results := []CheckResult{
		{Name: "check1", OK: true},
		{Name: "check2", OK: false, Note: "timeout"},
	}
	log := FormatLogEntry(results)
	if log == "" {
		t.Error("expected non-empty log entry")
	}
}
