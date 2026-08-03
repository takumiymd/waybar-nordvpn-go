// Package nordvpn queries the NordVPN command-line client and parses its human readable "nordvpn status" output into a structured value.
package nordvpn

import (
	"os/exec"
	"strings"
)

type State int

const (
	StateUnknown State = iota
	StateConnected
	StateDisconnected
	StateUnavailable
	StateError
)

// Status is the parsed result of a single "nordvpn status" invocation
type Status struct {
	State    State
	Country  string
	City     string
	IP       string
	Protocol string
	Raw      string
}

var binary = "nordvpn"

var (
	lookPath  = exec.LookPath
	runStatus = defaultRunStatus
)

// defaultRunStatus executes "nordvpn status" and returns its combined output
func defaultRunStatus() string {
	out, _ := exec.Command(binary, "status").CombinedOutput()
	return string(out)
}

// Query runs the NordVPN CLI status query and returns a parsed Status.
func Query() Status {
	if _, err := lookPath(binary); err != nil {
		return Status{State: StateUnavailable}
	}
	return Parse(runStatus())
}

// Parse extracts connection details from the raw NordVPN status output.
func Parse(raw string) Status {
	f := parseFields(raw)
	switch f["Status"] {
	case "Connected":
		return Status{
			State:    StateConnected,
			Country:  f["Country"],
			City:     f["City"],
			IP:       f["IP"],
			Protocol: f["Current protocol"],
			Raw:      raw,
		}
	case "Disconnected":
		return Status{State: StateDisconnected, Raw: raw}
	default:
		return Status{State: StateError, Raw: raw}
	}
}

// parseFields splits the raw status output into key-value pairs.
func parseFields(raw string) map[string]string {
	fields := make(map[string]string)
	for _, line := range strings.Split(raw, "\n") {
		key, val, ok := strings.Cut(line, ": ")
		if !ok {
			continue
		}
		fields[key] = strings.TrimSpace(val)
	}
	return fields
}
