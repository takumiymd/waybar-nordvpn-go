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

type Output struct {
	Text    string `json:"text"`
	Class   string `json:"class"`
	Tooltip string `json:"tooltip"`
}

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

func (o Output) Emit(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(o)
}
