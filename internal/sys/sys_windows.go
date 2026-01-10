//go:build windows

package sys

import "time"

type ActionResult struct {
	ExitCode int    `json:"exit_code"`
	Summary  string `json:"summary"`
	Details  string `json:"details,omitempty"`
}

func EnsureSameUser(pid int32) error { return nil }

func TerminateProcess(pid int32, timeout time.Duration) ActionResult {
	return ActionResult{ExitCode: 1, Summary: "Not implemented on Windows yet"}
}

func SmartRestart(pid int32, cmdline string, timeout time.Duration) ActionResult {
	return ActionResult{ExitCode: 1, Summary: "Not implemented on Windows yet"}
}

func RestartDockerContainer(containerID string, timeout time.Duration) ActionResult {
	return ActionResult{ExitCode: 1, Summary: "Not implemented on Windows yet"}
}
