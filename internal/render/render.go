package render

import (
	"fmt"
	"strings"

	"portik/internal/model"
	"portik/internal/proctree"
	"portik/internal/sys"
)

func Who(rep model.Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Port %d/%s\n", rep.Port, rep.Proto)

	if len(rep.Listeners) == 0 {
		b.WriteString("  (no listeners found)\n")
	} else {
		for _, l := range rep.Listeners {
			fmt.Fprintf(&b, "%-7s %-24s pid=%d  user=%s  %-12s %s\n",
				l.State,
				fmt.Sprintf("%s:%d", fmtIP(l.LocalIP), l.LocalPort),
				l.PID,
				dash(l.User),
				dash(l.ProcName),
				dash(l.Cmdline),
			)
		}
	}

	if rep.Docker.Checked {
		if rep.Docker.Mapped {
			fmt.Fprintf(&b, "\nDocker: mapped from %s (%s) service=%s containerPort=%s\n",
				rep.Docker.ContainerID, rep.Docker.ContainerName, dash(rep.Docker.ComposeService), rep.Docker.ContainerPort)
		} else {
			b.WriteString("\nDocker: not mapped\n")
		}
	}

	return b.String()
}

func Explain(rep model.Report) string {
	var b strings.Builder
	b.WriteString(Who(rep))

	b.WriteString("\nSummary\n")
	if len(rep.Diagnostics) == 0 {
		b.WriteString("- No hints available\n")
	} else {
		for _, d := range rep.Diagnostics {
			fmt.Fprintf(&b, "- [%s] %s\n", strings.ToUpper(d.Severity), d.Summary)
		}
	}

	b.WriteString("\nDetails\n")
	for _, d := range rep.Diagnostics {
		fmt.Fprintf(&b, "• %s\n", d.Summary)
		if d.Details != "" {
			fmt.Fprintf(&b, "  %s\n", d.Details)
		}
		if d.Action != "" {
			fmt.Fprintf(&b, "  → %s\n", d.Action)
		}
	}
	return b.String()
}

func Blame(rep model.Report, chain []proctree.Proc, started proctree.StartedBy) string {
	var b strings.Builder
	b.WriteString(Who(rep))
	b.WriteString("\nProcess tree (child → parents)\n")
	for i, p := range chain {
		prefix := "└─"
		if i < len(chain)-1 {
			prefix = "├─"
		}
		fmt.Fprintf(&b, "%s pid=%d ppid=%d user=%s name=%s\n", prefix, p.PID, p.PPID, dash(p.User), dash(p.Name))
		if p.Cmdline != "" {
			fmt.Fprintf(&b, "   cmd: %s\n", p.Cmdline)
		}
	}
	b.WriteString("\nWho started this? (best-effort)\n")
	fmt.Fprintf(&b, "- %s", strings.ToUpper(started.Kind))
	if started.Details != "" {
		fmt.Fprintf(&b, ": %s", started.Details)
	}
	b.WriteString("\n")
	if rep.Docker.Mapped {
		b.WriteString("- Port is mapped from Docker; the owning process may be docker-proxy or inside the container.\n")
	}
	return b.String()
}

func ActionResult(r sys.ActionResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", r.Summary)
	if r.Details != "" {
		fmt.Fprintf(&b, "%s\n", r.Details)
	}
	return b.String()
}

func fmtIP(ip string) string {
	if ip == "" {
		return "*"
	}
	if strings.Contains(ip, ":") && !strings.HasPrefix(ip, "[") {
		return "[" + ip + "]"
	}
	return ip
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
