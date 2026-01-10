//go:build !linux && !darwin

package sockets

import (
	"fmt"

	"portik/internal/model"
)

func inspect(port int, proto string, includeConnections bool) ([]model.Listener, []model.Conn, error) {
	return nil, nil, fmt.Errorf("unsupported OS for socket inspection")
}
