package reserve

import (
	"context"
	"fmt"
	"net"
	"time"
)

type Reservation struct {
	Port  int       `json:"port"`
	Proto string    `json:"proto"`
	Bind  string    `json:"bind"`
	Until time.Time `json:"until"`
}

type Handle interface {
	Close() error
}

type tcpHandle struct{ ln net.Listener }

func (h tcpHandle) Close() error { return h.ln.Close() }

type udpHandle struct{ pc net.PacketConn }

func (h udpHandle) Close() error { return h.pc.Close() }

// Reserve binds a port (or port=0 for ephemeral) and holds it until ctx is done.
// Returns reservation info + a handle you must Close (usually via defer).
func Reserve(ctx context.Context, proto, bind string, port int) (Reservation, Handle, error) {
	if bind == "" {
		bind = "127.0.0.1"
	}
	if proto != "tcp" && proto != "udp" {
		return Reservation{}, nil, fmt.Errorf("invalid proto: %s", proto)
	}

	addr := net.JoinHostPort(bind, fmt.Sprintf("%d", port))

	if proto == "tcp" {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return Reservation{}, nil, err
		}
		p := ln.Addr().(*net.TCPAddr).Port
		res := Reservation{Port: p, Proto: "tcp", Bind: bind, Until: time.Now().Add(0)}
		go func() {
			<-ctx.Done()
			_ = ln.Close()
		}()
		return res, tcpHandle{ln: ln}, nil
	}

	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		return Reservation{}, nil, err
	}
	p := pc.LocalAddr().(*net.UDPAddr).Port
	res := Reservation{Port: p, Proto: "udp", Bind: bind, Until: time.Now().Add(0)}
	go func() {
		<-ctx.Done()
		_ = pc.Close()
	}()
	return res, udpHandle{pc: pc}, nil
}
