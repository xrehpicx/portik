package reserve

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type FreeOptions struct {
	Proto string // tcp|udp
	Bind  string // ip/host, default 127.0.0.1
	// If RangeStart/RangeEnd are 0, caller should request ephemeral port (0) via FindFreeEphemeral.
	RangeStart int
	RangeEnd   int
	// Attempts for random sampling in range; if 0 uses full scan
	Attempts int
}

func FindFreeInRange(ctx context.Context, opt FreeOptions) (int, error) {
	if opt.Bind == "" {
		opt.Bind = "127.0.0.1"
	}
	if opt.Proto != "tcp" && opt.Proto != "udp" {
		return 0, fmt.Errorf("invalid proto: %s", opt.Proto)
	}
	if opt.RangeStart <= 0 || opt.RangeEnd <= 0 || opt.RangeStart > 65535 || opt.RangeEnd > 65535 {
		return 0, fmt.Errorf("invalid range")
	}
	if opt.RangeStart > opt.RangeEnd {
		opt.RangeStart, opt.RangeEnd = opt.RangeEnd, opt.RangeStart
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := opt.RangeEnd - opt.RangeStart + 1

	try := func(p int) bool {
		ok, _ := isBindable(opt.Proto, opt.Bind, p)
		return ok
	}

	// random attempts first (fast)
	attempts := opt.Attempts
	if attempts <= 0 {
		attempts = min(64, total)
	}
	for i := 0; i < attempts; i++ {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		p := opt.RangeStart + r.Intn(total)
		if try(p) {
			return p, nil
		}
	}

	// fallback: linear scan
	for p := opt.RangeStart; p <= opt.RangeEnd; p++ {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		if try(p) {
			return p, nil
		}
	}

	return 0, fmt.Errorf("no free port found in range %d-%d", opt.RangeStart, opt.RangeEnd)
}

func FindFreeEphemeral(opt FreeOptions) (int, error) {
	if opt.Bind == "" {
		opt.Bind = "127.0.0.1"
	}
	if opt.Proto != "tcp" && opt.Proto != "udp" {
		return 0, fmt.Errorf("invalid proto: %s", opt.Proto)
	}

	if opt.Proto == "tcp" {
		ln, err := net.Listen("tcp", net.JoinHostPort(opt.Bind, "0"))
		if err != nil {
			return 0, err
		}
		defer ln.Close()
		return ln.Addr().(*net.TCPAddr).Port, nil
	}

	pc, err := net.ListenPacket("udp", net.JoinHostPort(opt.Bind, "0"))
	if err != nil {
		return 0, err
	}
	defer pc.Close()
	return pc.LocalAddr().(*net.UDPAddr).Port, nil
}

func isBindable(proto, bind string, port int) (bool, error) {
	addr := net.JoinHostPort(bind, fmt.Sprintf("%d", port))
	if proto == "tcp" {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return false, err
		}
		_ = ln.Close()
		return true, nil
	}
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		return false, err
	}
	_ = pc.Close()
	return true, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
