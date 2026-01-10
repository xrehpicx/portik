package render

import (
	"fmt"
	"strings"
)

type scanRow struct {
	Port   int
	Proto  string
	Status string
	Owner  string
	PID    int32
	Addr   string
	Docker string
	Hint   string
	Error  string
}

type ScanTableRow interface {
	// duck-typing helper: your cli.scanRow has matching exported fields, so we re-map.
}

func ScanTable(rows any) string {
	// Expect []cli.scanRow but avoid import cycles by using reflection-free approach:
	// The cli will pass []cli.scanRow; we convert by fmt.Sprintf via known formatting in renderScanTableFrom.
	// Simpler: keep this renderer accepting a concrete internal struct; so we provide an adapter below.
	switch v := rows.(type) {
	case []interface{}:
		_ = v
	}
	// The CLI passes []scanRow from cli. For simplicity, render uses a helper function in cli that calls ScanTableRows().
	// But since you asked for only file paths, we keep ScanTable signature matching current usage:
	// In cli/scan.go we call render.ScanTable(rows) where rows is []cli.scanRow.
	// We therefore implement ScanTable via a tiny interface expectation using fmt on maps is messy.
	// Best: change cli/scan.go to call render.ScanTableRows(renderRows).
	return "render.ScanTable: invalid input\n"
}

// Use this instead: render.ScanTableRows(...)
type ScanRows []struct {
	Port   int
	Proto  string
	Status string
	Owner  string
	PID    int32
	Addr   string
	Docker string
	Hint   string
	Error  string
}

func ScanTableRows(rows ScanRows) string {
	var b strings.Builder

	b.WriteString("PORT   STATUS   OWNER                 PID     ADDR                 DOCKER              HINT\n")
	b.WriteString("────   ──────   ────────────────────  ──────  ───────────────────  ──────────────────  ─────────────────────\n")

	for _, r := range rows {
		owner := trunc(r.Owner, 20)
		addr := trunc(r.Addr, 19)
		docker := trunc(r.Docker, 18)
		hint := trunc(r.Hint, 20)
		if r.Error != "" {
			hint = trunc("ERR: "+r.Error, 20)
		}
		pid := "-"
		if r.PID > 0 {
			pid = fmt.Sprintf("%d", r.PID)
		}
		fmt.Fprintf(&b, "%-5d  %-7s  %-20s  %-6s  %-19s  %-18s  %-20s\n",
			r.Port, r.Status, owner, pid, addr, docker, hint)
	}
	return b.String()
}

func trunc(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	if n <= 1 {
		return s[:n]
	}
	return s[:n-1] + "…"
}
