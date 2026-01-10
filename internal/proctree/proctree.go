package proctree

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Proc struct {
	PID     int32  `json:"pid"`
	PPID    int32  `json:"ppid,omitempty"`
	User    string `json:"user,omitempty"`
	Name    string `json:"name,omitempty"`
	Cmdline string `json:"cmdline,omitempty"`
}

type StartedBy struct {
	Kind    string `json:"kind"` // systemd|container|launchd|unknown
	Details string `json:"details,omitempty"`
}

func Build(pid int32, maxDepth int) ([]Proc, StartedBy) {
	if maxDepth <= 0 {
		maxDepth = 10
	}
	var chain []Proc
	cur := pid
	seen := map[int32]bool{}

	for i := 0; i < maxDepth && cur > 0; i++ {
		if seen[cur] {
			break
		}
		seen[cur] = true

		p := procInfo(cur)
		chain = append(chain, p)

		if p.PPID <= 0 || p.PPID == cur || p.PID == 1 {
			break
		}
		cur = p.PPID
	}

	return chain, whoStarted(pid)
}

func procInfo(pid int32) Proc {
	ppidStr := psField(pid, "ppid=")
	userStr := psField(pid, "user=")
	commStr := psField(pid, "comm=")
	cmdStr := psField(pid, "command=")

	return Proc{
		PID:     pid,
		PPID:    atoi32(ppidStr),
		User:    strings.TrimSpace(userStr),
		Name:    strings.TrimSpace(commStr),
		Cmdline: strings.TrimSpace(cmdStr),
	}
}

func psField(pid int32, format string) string {
	cmd := exec.Command("ps", "-p", itoa32(pid), "-o", format)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	_ = cmd.Run()
	return strings.TrimSpace(buf.String())
}

func whoStarted(pid int32) StartedBy {
	switch runtime.GOOS {
	case "linux":
		if unit := systemdUnitFromCgroup(pid); unit != "" {
			return StartedBy{Kind: "systemd", Details: unit}
		}
		if cid := containerIDFromCgroup(pid); cid != "" {
			return StartedBy{Kind: "container", Details: cid}
		}
		if s := systemctlStatusHint(pid); s != "" {
			return StartedBy{Kind: "systemd", Details: s}
		}
		return StartedBy{Kind: "unknown"}
	case "darwin":
		if parentLooksLikeLaunchd(pid) {
			return StartedBy{Kind: "launchd", Details: "parent chain includes launchd"}
		}
		return StartedBy{Kind: "unknown"}
	default:
		return StartedBy{Kind: "unknown"}
	}
}

func parentLooksLikeLaunchd(pid int32) bool {
	cur := pid
	for i := 0; i < 15 && cur > 0; i++ {
		name := strings.ToLower(psField(cur, "comm="))
		if strings.Contains(name, "launchd") {
			return true
		}
		ppid := atoi32(psField(cur, "ppid="))
		if ppid <= 0 || ppid == cur {
			break
		}
		cur = ppid
	}
	return false
}

func systemctlStatusHint(pid int32) string {
	if _, err := exec.LookPath("systemctl"); err != nil {
		return ""
	}
	out, err := exec.Command("systemctl", "status", itoa32(pid), "--no-pager").Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.Contains(l, ".service") {
			return l
		}
	}
	return ""
}

func systemdUnitFromCgroup(pid int32) string {
	b, err := os.ReadFile("/proc/" + itoa32(pid) + "/cgroup")
	if err != nil {
		return ""
	}
	txt := string(b)
	for _, line := range strings.Split(txt, "\n") {
		if strings.Contains(line, "system.slice/") && strings.Contains(line, ".service") {
			i := strings.Index(line, "system.slice/")
			if i < 0 {
				continue
			}
			sub := line[i+len("system.slice/"):]
			j := strings.Index(sub, ".service")
			if j < 0 {
				continue
			}
			return sub[:j+len(".service")]
		}
	}
	return ""
}

func containerIDFromCgroup(pid int32) string {
	b, err := os.ReadFile("/proc/" + itoa32(pid) + "/cgroup")
	if err != nil {
		return ""
	}
	txt := string(b)
	tokens := splitNonHex(txt)
	best := ""
	for _, t := range tokens {
		if len(t) >= 12 && len(t) <= 64 && isHex(t) {
			if len(t) > len(best) {
				best = t
			}
		}
	}
	return best
}

func splitNonHex(s string) []string {
	var out []string
	var cur strings.Builder
	for _, r := range s {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
			cur.WriteRune(r)
		} else {
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
		}
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}

func isHex(s string) bool {
	for _, r := range s {
		ok := (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
		if !ok {
			return false
		}
	}
	return true
}

func itoa32(n int32) string {
	if n == 0 {
		return "0"
	}
	x := int(n)
	var b [32]byte
	i := len(b)
	for x > 0 {
		i--
		b[i] = byte('0' + x%10)
		x /= 10
	}
	return string(b[i:])
}

func atoi32(s string) int32 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	sign := int32(1)
	if strings.HasPrefix(s, "-") {
		sign = -1
		s = strings.TrimPrefix(s, "-")
	}
	var n int32
	for _, r := range s {
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int32(r-'0')
	}
	return sign * n
}
