package main

import (
	"fmt"
	"strings"

	"github.com/smmr-software/mabel/internal/full"
	"github.com/smmr-software/mabel/internal/mini"

	flag "github.com/spf13/pflag"

	gloss "github.com/charmbracelet/lipgloss"
)

func main() {
	var (
		help = flag.BoolP("help", "h", false, "Print this help message.")
	)
	flag.Parse()
	args := flag.Args()

	if *help {
		green := gloss.NewStyle().Foreground(gloss.Color("2"))
		yellow := gloss.NewStyle().Foreground(gloss.Color("3"))

		var menu strings.Builder
		menu.WriteString(green.Render("Mabel"))
		menu.WriteString("\nSMMR Software <hello@smmr.software>\n")
		menu.WriteString("A fancy BitTorrent client for the terminal.\n\n")
		menu.WriteString(yellow.Render("USAGE:"))
		menu.WriteString("\n    mabel [OPTIONS] [TORRENT]...\n\n")
		menu.WriteString(yellow.Render("OPTIONS:"))
		menu.WriteString("\n    " + green.Render("-h") + ", " + green.Render("--help"))
		menu.WriteString("\n        Print this help message.\n\n")
		menu.WriteString(yellow.Render("ARGS:"))
		menu.WriteString("\n    " + green.Render("<TORRENT>..."))
		menu.WriteString("\n        An optional list of infohashes, magnet links, and/or torrent")
		menu.WriteString("\n        files. If one torrent is provided, Mabel starts in \"mini\" mode.")
		menu.WriteString("\n        Otherwise, the full TUI client opens with the corresponding")
		menu.WriteString("\n        torrents added.")

		fmt.Println(menu.String())
	} else if len(args) == 1 {
		mini.Execute(&args[0])
	} else {
		full.Execute()
	}
}
