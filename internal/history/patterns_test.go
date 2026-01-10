package history

import (
	"testing"
	"time"
)

func TestDetectPatternsHour(t *testing.T) {
	var evs []OwnershipEvent
	now := time.Now()
	for i := 0; i < 8; i++ {
		evs = append(evs, OwnershipEvent{
			At:       time.Date(now.Year(), now.Month(), now.Day()-i, 9, 5, 0, 0, now.Location()),
			ProcName: "postgres",
			User:     "me",
		})
	}
	p := DetectPatterns(evs)
	if len(p) == 0 {
		t.Fatalf("expected some patterns")
	}
	found := false
	for _, x := range p {
		if x.Kind == "hour-of-day" || x.Kind == "owner-at-hour" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected hour patterns, got %#v", p)
	}
}
