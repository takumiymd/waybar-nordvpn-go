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
