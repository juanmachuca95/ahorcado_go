package main

import (
	"flag"
	"os"

	"github.com/juanmachuca95/ahorcado_go/cmd/clients/cli/game"
	"github.com/pterm/pterm"
)

var (
	serverAddr = flag.String("host", "localhost:8080", "The server address in the format of host:port")
	username   = flag.String("user", "Guest", "Username to join the game")
)

func main() {
	flag.Parse()
	banner()

	g := game.NewGame(*serverAddr, *username)
	g.PanelOptions()
}

func banner() {
	f, err := os.ReadFile("./cmd/clients/cli/banner.txt")
	if err != nil {
		panic("Cannot read file banner.txt")
	}
	banner := string(f)
	style := pterm.NewStyle(pterm.FgYellow, pterm.Bold)
	style.Println(banner)
}
