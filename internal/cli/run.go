package cli

import (
	"fmt"
	"os"
)

func Run(args []string) int {
	if len(args) == 0 {
		printHelp()
		return 2
	}
	switch args[0] {
	case "help", "--help", "-h":
		printHelp()
		return 0
	case "version", "--version":
		fmt.Println("portik v0.2.0")
		return 0
	case "who":
		return runWho(args[1:])
	case "explain":
		return runExplain(args[1:])
	case "kill":
		return runKill(args[1:])
	case "restart":
		return runRestart(args[1:])
	case "watch":
		return runWatch(args[1:])
	case "history":
		return runHistory(args[1:])
	case "daemon":
		return runDaemon(args[1:])
	case "blame":
		return runBlame(args[1:])

	// NEW
	case "scan":
		return runScan(args[1:])
	case "free":
		return runFree(args[1:])
	case "reserve":
		return runReserve(args[1:])

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", args[0])
		printHelp()
		return 2
	}
}

func printHelp() {
	fmt.Print(`portik — port ownership inspector

Usage:
  portik <command> [args] [flags]

Commands:
  who <port>        Show who is listening on a port
  explain <port>    Explain likely reasons a port is stuck / bind fails
  kill <port>       Terminate the process owning a port (safe by default)
  restart <port>    Smart restart (kill + restart last command)
  watch <port>      Watch a port and record changes
  history <port>    Show port ownership history (+ pattern detection)
  daemon            Monitor multiple ports and record history (foreground)
  blame <port>      Process tree + who started this (best-effort)

  scan              Scan a set/range of ports and show a table
  free              Find a free port (optionally within a range)
  reserve           Reserve a port by binding it for a duration

  version           Show version

Common flags (per command):
  --proto tcp|udp
  --docker          Enable Docker mapping (shells out to docker)
  --json            JSON output (where supported)
  --yes             Skip confirmation prompts
`)
}
