package nordvpn

import (
	"errors"
	"testing"
)

const updateBanner = "A new version of NordVPN is available!\nPlease update the app.\n"

func TestParseConnected(t *testing.T) {
	raw := updateBanner + `Status: Connected
Hostname: jp123.nordvpn.com
IP: 203.0.113.7
Country: Japan
City: Tokyo
Current protocol: NORDLYNX
`
	got := Parse(raw)
	if got.State != StateConnected {
		t.Fatalf("State = %v, want StateConnected", got.State)
	}
	for _, c := range []struct{ name, got, want string }{
		{"Country", got.Country, "Japan"},
		{"City", got.City, "Tokyo"},
		{"IP", got.IP, "203.0.113.7"},
		{"Protocol", got.Protocol, "NORDLYNX"},
	} {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestParseDisconnected(t *testing.T) {
	got := Parse(updateBanner + "Status: Disconnected\n")
	if got.State != StateDisconnected {
		t.Fatalf("State = %v, want StateDisconnected", got.State)
	}
}

func TestParseUnparseable(t *testing.T) {
	got := Parse("some unexpected error text\n")
	if got.State != StateError {
		t.Fatalf("State = %v, want StateError", got.State)
	}
	if got.Raw == "" {
		t.Error("Raw output should be preserved for the error tooltip")
	}
}

func TestQueryUnavailable(t *testing.T) {
	orig := lookPath
	defer func() { lookPath = orig }()
	lookPath = func(string) (string, error) { return "", errors.New("not found") }

	if got := Query(); got.State != StateUnavailable {
		t.Fatalf("State = %v, want StateUnavailable", got.State)
	}
}

func TestQueryParsesRunner(t *testing.T) {
	origLook, origRun := lookPath, runStatus
	defer func() { lookPath, runStatus = origLook, origRun }()
	lookPath = func(string) (string, error) { return "/usr/bin/nordvpn", nil }
	runStatus = func() string { return "Status: Disconnected\n" }

	if got := Query(); got.State != StateDisconnected {
		t.Fatalf("State = %v, want StateDisconnected", got.State)
	}
}
