package waybar

import (
	"strings"
	"testing"

	"github.com/takumiymd/waybar-nordvpn-go/internal/nordvpn"
)

func TestFromStatusConnected(t *testing.T) {
	out := FromStatus(nordvpn.Status{
		State:    nordvpn.StateConnected,
		Country:  "Japan",
		City:     "Tokyo",
		IP:       "203.0.113.7",
		Protocol: "NORDLYNX",
	})
	if out.Text != "VPN: Japan-Tokyo up" {
		t.Errorf("Text = %q", out.Text)
	}
	if out.Class != "connected" {
		t.Errorf("Class = %q", out.Class)
	}
	if !strings.Contains(out.Tooltip, "IP: 203.0.113.7") {
		t.Errorf("Tooltip missing IP: %q", out.Tooltip)
	}
}

func TestFromStatusDisconnected(t *testing.T) {
	out := FromStatus(nordvpn.Status{State: nordvpn.StateDisconnected})
	if out.Text != "VPN: down" || out.Class != "disconnected" {
		t.Errorf("got %+v", out)
	}
	if strings.Contains(out.Tooltip, "NordVPN: Connected") {
		t.Errorf("disconnected tooltip wrongly says Connected: %q", out.Tooltip)
	}
}

func TestFromStatusUnavailable(t *testing.T) {
	out := FromStatus(nordvpn.Status{State: nordvpn.StateUnavailable})
	if out.Class != "error" || !strings.Contains(out.Tooltip, "not found") {
		t.Errorf("got %+v", out)
	}
}

func TestFromStatusError(t *testing.T) {
	out := FromStatus(nordvpn.Status{State: nordvpn.StateError, Raw: "boom"})
	if out.Text != "VPN: !!" || out.Class != "error" {
		t.Errorf("got %+v", out)
	}
	if !strings.Contains(out.Tooltip, "boom") {
		t.Errorf("error tooltip should include raw output: %q", out.Tooltip)
	}
}

func TestEmitIsSingleJSONLineFieldOrder(t *testing.T) {
	var b strings.Builder
	if err := (Output{Text: "t", Class: "c", Tooltip: "tt"}).Emit(&b); err != nil {
		t.Fatal(err)
	}
	got := b.String()
	want := `{"text":"t","class":"c","tooltip":"tt"}` + "\n"
	if got != want {
		t.Errorf("Emit = %q, want %q", got, want)
	}
}
