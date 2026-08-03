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

func defaultRunStatus() string {
	out, _ := exec.Command(binary, "status").CombinedOutput()
	return string(out)
}

func Query() Status {
	if _, err := lookPath(binary); err != nil {
		return Status{State: StateUnavailable}
	}
	return Parse(runStatus())
}

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
