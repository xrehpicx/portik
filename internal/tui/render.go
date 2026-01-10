//go:build tui

package tui

import (
	"fmt"
	"strings"
	"time"

	"portik/internal/history"
)

func renderUI(m modelTUI) string {
	var b strings.Builder

	header := fmt.Sprintf("portik TUI  proto=%s  interval=%s", m.opts.Proto, m.opts.Interval)
	if m.opts.Docker {
		header += "  docker=on"
	} else {
		header += "  docker=off"
	}
	if m.opts.Actions {
		header += "  actions=on (K=kill, R=restart)"
	} else {
		header += "  actions=off (run with --actions)"
	}

	b.WriteString(header)
	b.WriteString("\n")

	if m.filtering {
		b.WriteString(fmt.Sprintf("Filter: %s▌  (Enter to apply, Esc to clear)\n", m.filter))
	} else if strings.TrimSpace(m.filter) != "" {
		b.WriteString(fmt.Sprintf("Filter: %s  (press Esc to clear, / to edit)\n", m.filter))
	} else {
		b.WriteString("Filter: (press / to search)\n")
	}

	b.WriteString(strings.Repeat("─", clamp(m.width, 40, 140)))
	b.WriteString("\n")

	rows := m.filteredRows()
	if len(rows) == 0 {
		b.WriteString("(no rows; adjust filter)\n")
		return b.String()
	}

	sel := m.selected
	if sel < 0 {
		sel = 0
	}
	if sel >= len(rows) {
		sel = len(rows) - 1
	}

	b.WriteString("PORT   OWNER                    PID     ADDR                  DOCKER            Δ  SPARK\n")
	b.WriteString(strings.Repeat("─", clamp(m.width, 40, 140)))
	b.WriteString("\n")

	for i, r := range rows {
		prefix := " "
		if i == sel {
			prefix = ">"
		}
		changed := " "
		if r.Changed && !r.LastChange.IsZero() && time.Since(r.LastChange) < 5*time.Second {
			changed = "*"
		}
		owner := trunc(r.Owner, 23)
		addr := trunc(r.Addr, 20)
		dk := trunc(r.Docker, 16)
		spark := trunc(r.Spark, 12)
		if r.Err != "" {
			owner = trunc("ERR: "+r.Err, 23)
		}
		b.WriteString(fmt.Sprintf("%s%-5s %-23s %-7s %-20s %-16s %s  %-12s\n",
			prefix,
			fmt.Sprintf("%d", r.Port),
			owner,
			pidStr(r.PID),
			addr,
			dk,
			changed,
			spark,
		))
	}

	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", clamp(m.width, 40, 140)))
	b.WriteString("\n")

	row := rows[sel]

	title := "Details"
	if m.mode == viewExplain {
		title = "Explain"
	}
	b.WriteString(fmt.Sprintf("%s (Tab toggle, w=who, e=explain)\n", title))
	b.WriteString(truncLines(detailsForRow(row, m.mode), 18))
	b.WriteString("\n")

	b.WriteString(strings.Repeat("─", clamp(m.width, 40, 140)))
	b.WriteString("\n")
	b.WriteString("History (last 20)\n")
	b.WriteString(renderHistoryLast20(m, row.Port))
	b.WriteString("\n")

	if m.confirming {
		b.WriteString("\n")
		b.WriteString(strings.Repeat("─", clamp(m.width, 40, 140)))
		b.WriteString("\n")
		b.WriteString("CONFIRM: " + m.confirmMsg + "\n")
	} else {
		b.WriteString("\n")
		b.WriteString("Status: " + m.status + "\n")
		b.WriteString("Keys: ↑/↓ select, Tab toggle who/explain, r refresh, / filter, Esc clear, q quit")
		if m.opts.Actions {
			b.WriteString(", K kill, R restart")
		}
		b.WriteString("\n")
	}

	return b.String()
}

func renderHistoryLast20(m modelTUI, port int) string {
	if m.store == nil {
		return "  (no history loaded)\n"
	}
	key := fmt.Sprintf("%d/%s", port, m.opts.Proto)
	evs := m.store.Ports[key]
	if len(evs) == 0 {
		for _, k := range []string{fmt.Sprintf("%d/tcp", port), fmt.Sprintf("%d/udp", port)} {
			if len(m.store.Ports[k]) > 0 {
				key = k
				evs = m.store.Ports[k]
				break
			}
		}
	}
	if len(evs) == 0 {
		return "  (no events)\n"
	}
	start := 0
	if len(evs) > 20 {
		start = len(evs) - 20
	}
	var b strings.Builder
	for _, e := range evs[start:] {
		lbl := ownerLabelEvent(e)
		b.WriteString(fmt.Sprintf("  %s  %s\n", e.At.Format("01-02 15:04:05"), trunc(lbl, 48)))
	}
	return b.String()
}

func ownerLabelEvent(e history.OwnershipEvent) string {
	if e.DockerMapped {
		if e.ComposeService != "" {
			return fmt.Sprintf("docker:%s (svc=%s)", e.ContainerName, e.ComposeService)
		}
		return "docker:" + e.ContainerName
	}
	if e.ProcName != "" {
		if e.User != "" {
			return fmt.Sprintf("%s (%s)", e.ProcName, e.User)
		}
		return e.ProcName
	}
	if e.PID > 0 {
		return fmt.Sprintf("pid:%d", e.PID)
	}
	return "none"
}

func pidStr(pid int32) string {
	if pid <= 0 {
		return "-"
	}
	return fmt.Sprintf("%d", pid)
}

func trunc(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= n {
		return s
	}
	if n <= 1 {
		return s[:n]
	}
	return s[:n-1] + "…"
}

func truncLines(s string, maxLines int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s
	}
	return strings.Join(lines[:maxLines], "\n") + "\n…\n"
}

func clamp(v, lo, hi int) int {
	if v == 0 {
		return 100
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
