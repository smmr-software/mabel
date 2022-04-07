package main

import (
	"fmt"
	"strings"

	"github.com/smmr-software/mabel/full"
	"github.com/smmr-software/mabel/internal/styles"
	"github.com/smmr-software/mabel/mini"

	"github.com/adrg/xdg"
	flag "github.com/spf13/pflag"

	gloss "github.com/charmbracelet/lipgloss"
)

var (
	version = "v0.0.0"
	commit  = "unknown"
	builtBy = "unknown"
)

func main() {
	var (
		download = flag.StringP("download", "d", xdg.UserDirs.Download, "Set the default directory for downloaded torrents.")
		port     = flag.UintP("port", "p", 42069, "Set the port number to which the client will bind.")
		theme    = flag.StringP("theme", "t", "default", "Set the color theme that the client will use.")
		help     = flag.BoolP("help", "h", false, "Print this help message.")
		vrsn     = flag.BoolP("version", "v", false, "Print version information.")
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
		menu.WriteString("\n        Set the torrent download directory. [default: " + green.Render("$XDG_DOWNLOAD_DIR") + "]")
		menu.WriteString("\n    " + green.Render("-p") + ", " + green.Render("--port"))
		menu.WriteString("\n        Set the port number to which the client will bind. [default: " + green.Render("42069") + "]")
		menu.WriteString("\n    " + green.Render("-t") + ", " + green.Render("--theme"))
		menu.WriteString("\n        Set the color theme that the client will use. [default: " + green.Render("default") + "]")
		menu.WriteString("\n    " + green.Render("-h") + ", " + green.Render("--help"))
		menu.WriteString("\n        Print this help message.")
		menu.WriteString("\n    " + green.Render("-v") + ", " + green.Render("--version"))
		menu.WriteString("\n        Print version information.\n\n")
		menu.WriteString(yellow.Render("ARGS:"))
		menu.WriteString("\n    " + green.Render("<TORRENT>..."))
		menu.WriteString("\n        An optional list of infohashes, magnet links, and/or torrent")
		menu.WriteString("\n        files. If one torrent is provided, Mabel starts in \"mini\" mode.")
		menu.WriteString("\n        Otherwise, the full TUI client opens with the corresponding")
		menu.WriteString("\n        torrents added.")

		fmt.Println(menu.String())
		return
	} else if *vrsn {
		fmt.Printf(
			"Mabel %s\nCommit: %s\nBuilt by: %s\n",
			version, commit, builtBy,
		)
		return
	}

	conf := getConfig()
	downloadFlag := flag.Lookup("download")
	portFlag := flag.Lookup("port")
	themeFlag := flag.Lookup("theme")

	if !downloadFlag.Changed && conf.Download != "" {
		flag.Set("download", conf.Download)
	}
	if !portFlag.Changed && conf.Port != 0 {
		flag.Set("port", fmt.Sprint(conf.Port))
	}
	if !themeFlag.Changed && conf.Theme != "" {
		flag.Set("theme", conf.Theme)
	}

	realTheme := styles.DefaultTheme
	/*if *theme == "placeholder" {
		realTheme = styles.PlaceholderTheme
	}*/

	if flag.NArg() == 1 {
		mini.Execute(&args[0], download, port, &realTheme)
	} else {
		full.Execute(&args, download, port, &realTheme)
	}
}
