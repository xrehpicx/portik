package render

import (
	"fmt"
	"sort"
	"strings"
)

type ConnAggRow struct {
	RemoteIP string
	Total    int
	ByState  map[string]int
	Samples  []string
}

// ConnTable renders top remote IPs and per-state counts.
func ConnTable(port int, proto string, rows any) string {
	// CLI passes []cli.connAggRow; we render via a tiny adapter to avoid import cycles.
	// So this function expects already-adapted []ConnAggRow in practice.
	switch v := rows.(type) {
	case []ConnAggRow:
		return connTableFrom(v, port, proto)
	default:
		return "conn: render: invalid input\n"
	}
}

func connTableFrom(rows []ConnAggRow, port int, proto string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Connections for %d/%s (top clients)\n", port, proto)
	b.WriteString("REMOTE IP              TOTAL   STATES                      SAMPLES\n")
	b.WriteString("────────────────────   ─────   ─────────────────────────   ─────────────────────────\n")

	for _, r := range rows {
		stateStr := formatStates(r.ByState, 3)
		samp := strings.Join(r.Samples, ", ")
		fmt.Fprintf(&b, "%-20s   %-5d   %-25s   %s\n",
			trunc(r.RemoteIP, 20),
			r.Total,
			trunc(stateStr, 25),
			trunc(samp, 25),
		)
	}
	return b.String()
}

func formatStates(m map[string]int, max int) string {
	if len(m) == 0 {
		return "-"
	}
	type kv struct {
		k string
		v int
	}
	var arr []kv
	for k, v := range m {
		arr = append(arr, kv{k: k, v: v})
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].v == arr[j].v {
			return arr[i].k < arr[j].k
		}
		return arr[i].v > arr[j].v
	})
	if max > 0 && len(arr) > max {
		arr = arr[:max]
	}
	var parts []string
	for _, x := range arr {
		parts = append(parts, fmt.Sprintf("%s:%d", x.k, x.v))
	}
	return strings.Join(parts, " ")
}
