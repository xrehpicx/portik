package inspect

import (
	"fmt"
	"os/user"
	"time"

	"portik/internal/docker"
	"portik/internal/model"
	"portik/internal/platform"
	"portik/internal/proc"
	"portik/internal/sockets"
)

type Options struct {
	EnableDocker       bool
	IncludeConnections bool
}

func InspectPort(port int, proto string, opt Options) (model.Report, error) {
	if proto != "tcp" && proto != "udp" {
		return model.Report{}, fmt.Errorf("unsupported proto: %s", proto)
	}

	u, _ := user.Current()
	hs := platform.HostSummary()

	rep := model.Report{
		Port:      port,
		Proto:     proto,
		Generated: time.Now(),
		Host: model.HostSummary{
			OS:       hs.OS,
			Arch:     hs.Arch,
			Hostname: hs.Hostname,
			Kernel:   hs.Kernel,
		},
		User: model.UserSummary{Username: safeUsername(u)},
	}

	listeners, conns, err := sockets.Inspect(port, proto, opt.IncludeConnections)
	if err != nil {
		return model.Report{}, err
	}
	for i := range listeners {
		proc.Enrich(&listeners[i])
	}
	for i := range conns {
		proc.EnrichConn(&conns[i])
	}

	rep.Listeners = listeners
	rep.Connections = conns

	if opt.EnableDocker {
		rep.Docker = docker.MapPort(port, proto)
	}

	rep.Diagnostics = Diagnose(rep)
	return rep, nil
}

func safeUsername(u *user.User) string {
	if u == nil {
		return ""
	}
	return u.Username
}
