package platform

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Summary struct {
	OS       string
	Arch     string
	Hostname string
	Kernel   string
}

func HostSummary() (s Summary) {
	s.OS = runtime.GOOS
	s.Arch = runtime.GOARCH
	if h, err := os.Hostname(); err == nil {
		s.Hostname = h
	}
	s.Kernel = unameR()
	return
}

func unameR() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(bytes.TrimSpace(out)))
}

func InContainer() bool {
	b, err := os.ReadFile("/proc/1/cgroup")
	if err == nil {
		txt := string(b)
		if strings.Contains(txt, "docker") || strings.Contains(txt, "kubepods") || strings.Contains(txt, "containerd") {
			return true
		}
	}
	if os.Getenv("container") != "" {
		return true
	}
	return false
}

func InWSL() bool {
	if b, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		return strings.Contains(strings.ToLower(string(b)), "microsoft")
	}
	return false
}

func InVM() bool {
	if b, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		txt := strings.ToLower(string(b))
		if strings.Contains(txt, "hypervisor") {
			return true
		}
	}
	return false
}
