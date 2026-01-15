package render

import (
	"fmt"
	"strings"
)

type TopRow struct {
	Port    int
	Proto   string
	Total   int
	Clients []TopClient
}

type TopClient struct {
	IP    string
	Count int
}

func TopTable(rows []TopRow, opt Options) string {
	var b strings.Builder
	if len(rows) == 0 {
		return "No active connections found.\n"
	}

	b.WriteString(label("TOP PORTS", opt))
	b.WriteString("\n")

	if opt.Summary {
		b.WriteString("  PORT   CONNS\n")
		b.WriteString("  -----  -----\n")
		for _, r := range rows {
			fmt.Fprintf(&b, "  %-5d  %-5d\n", r.Port, r.Total)
		}
		return b.String()
	}

	b.WriteString("  PORT   PROTO  CONNS  TOP CLIENTS\n")
	b.WriteString("  -----  -----  -----  ------------------------------\n")
	for _, r := range rows {
		b.WriteString(fmt.Sprintf("  %-5d  %-5s  %-5d  %s\n",
			r.Port, r.Proto, r.Total, formatClients(r.Clients)))
	}
	return b.String()
}

func formatClients(clients []TopClient) string {
	if len(clients) == 0 {
		return "-"
	}
	var parts []string
	for _, c := range clients {
		parts = append(parts, fmt.Sprintf("%s(%d)", c.IP, c.Count))
	}
	return strings.Join(parts, ", ")
}
