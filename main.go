// Command waybar-nordvpn queries the NordVPN CLI once and prints a single waybar JSON status line describing the current VPN connection.
// waybar runs the binary on an interval and it emits one JSON object and exits.
package main

import (
	"fmt"
	"os"

	"github.com/takumiymd/waybar-nordvpn-go/internal/nordvpn"
	"github.com/takumiymd/waybar-nordvpn-go/internal/waybar"
)

// version is the build version, overridden at link time by the Makefile via
// -ldflags '-X main.version=...'. It stays "dev" for plain `go build`.
var version = "dev"

// main runs the Waybar NordVPN status module query and writes it to stdout.
// It prints the build version and exits when invoked with -version.
func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-version" || os.Args[1] == "--version") {
		fmt.Println("waybar-nordvpn", version)
		return
	}
	if err := waybar.FromStatus(nordvpn.Query()).Emit(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "waybar-nordvpn:", err)
		os.Exit(1)
	}
}
