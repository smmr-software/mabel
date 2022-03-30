package main

import (
	"fmt"
	"strings"

	"github.com/smmr-software/mabel/internal/full"
	"github.com/smmr-software/mabel/internal/mini"

	"github.com/adrg/xdg"
	flag "github.com/spf13/pflag"

	gloss "github.com/charmbracelet/lipgloss"
)

func main() {
	var (
		download = flag.StringP("download", "d", xdg.UserDirs.Download, "Set the default directory for downloaded torrents.")
		port     = flag.UintP("port", "p", 42069, "Set the port number to which the client will bind.")
		help     = flag.BoolP("help", "h", false, "Print this help message.")
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
		menu.WriteString("\n    " + green.Render("-d") + ", " + green.Render("--download"))
		menu.WriteString("\n        Set the torrent download directory. Defaults to $XDG_DOWNLOAD_DIR.")
		menu.WriteString("\n    " + green.Render("-p") + ", " + green.Render("--port"))
		menu.WriteString("\n        Set the port number to which the client will bind. Defaults to 42069.")
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
		mini.Execute(&args[0], download, port)
	} else {
		full.Execute()
	}
}
