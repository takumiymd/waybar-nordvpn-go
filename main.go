// Command waybar-nordvpn queries the NordVPN CLI once and prints a single waybar JSON status line describing the current VPN connection.
// waybar runs the binary on an interval and it emits one JSON object and exits.
package main

import (
	"fmt"
	"os"

	"github.com/takumiymd/waybar-nordvpn-go/internal/nordvpn"
	"github.com/takumiymd/waybar-nordvpn-go/internal/waybar"
)

// main runs the Waybar NordVPN status module query and writes it to stdout.
func main() {
	if err := waybar.FromStatus(nordvpn.Query()).Emit(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "waybar-nordvpn:", err)
		os.Exit(1)
	}
}
