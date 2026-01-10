package inspect

import (
	"testing"

	"portik/internal/model"
)

func TestDiagnoseIPv6Only(t *testing.T) {
	rep := model.Report{
		Port:  5432,
		Proto: "tcp",
		Listeners: []model.Listener{
			{Family: "ipv6", State: "LISTEN", PID: 10, ProcName: "x"},
		},
	}
	d := Diagnose(rep)
	found := false
	for _, x := range d {
		if x.Kind == "ipv6-only" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected ipv6-only diagnostic")
	}
}
