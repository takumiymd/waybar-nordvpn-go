package main

import (
	"fmt"
	"os"

	"github.com/takumiymd/waybar-nordvpn-go/internal/nordvpn"
	"github.com/takumiymd/waybar-nordvpn-go/internal/waybar"
)

func main() {
	if err := waybar.FromStatus(nordvpn.Query()).Emit(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "waybar-nordvpn:", err)
		os.Exit(1)
	}
}
