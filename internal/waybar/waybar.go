// Package waybar maps a parsed NordVPN status into the JSON document that Waybar's custom module protocol expects.
package waybar

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/takumiymd/waybar-nordvpn-go/internal/nordvpn"
)

const (
	prefix    = "VPN:"
	errorMark = "!!"
)

// Output is the JSON object Waybar reads from a custom module.
type Output struct {
	Text    string `json:"text"`
	Class   string `json:"class"`
	Tooltip string `json:"tooltip"`
}

// FromStatus maps a parsed NordVPN status to a Waybar output line.
func FromStatus(s nordvpn.Status) Output {
	switch s.State {
	case nordvpn.StateConnected:
		return Output{
			Text:  fmt.Sprintf("%s %s-%s up", prefix, s.Country, s.City),
			Class: "connected",
			Tooltip: fmt.Sprintf(
				"\n  NordVPN: Connected\n  Connection: %s, %s\n  IP: %s\n  Protocol: %s",
				s.Country, s.City, s.IP, s.Protocol),
		}
	case nordvpn.StateDisconnected:
		return Output{
			Text:    prefix + " down",
			Class:   "disconnected",
			Tooltip: "\n  NordVPN: Disconnected\n  Right click: Connect",
		}
	case nordvpn.StateUnavailable:
		return Output{
			Text:    prefix + " down",
			Class:   "error",
			Tooltip: "NordVPN CLI not found in PATH.",
		}
	default:
		return Output{
			Text:    prefix + " " + errorMark,
			Class:   "error",
			Tooltip: "\nNordVPN status unavailable:\n" + s.Raw,
		}
	}
}

// Emit encodes the output as a single JSON line to w. HTML escaping is disabled so characters such as <, > and & survive verbatim.
func (o Output) Emit(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(o)
}
