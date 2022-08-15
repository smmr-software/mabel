// Package main creates and reads the config file, handles flags, and
// executes the mini and full clients.
package main

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/smmr-software/mabel/full"
	"github.com/smmr-software/mabel/internal/styles"
	"github.com/smmr-software/mabel/mini"

	"github.com/adrg/xdg"
	flag "github.com/spf13/pflag"

	gloss "github.com/charmbracelet/lipgloss"
)

var (
	version  = ""
	commit   = ""
	modified = ""
	builtBy  = ""
	date     = ""
)

func main() {
	var (
		download   = flag.StringP("download", "d", xdg.UserDirs.Download, "Set the default directory for downloaded torrents.")
		port       = flag.UintP("port", "p", 42069, "Set the port number to which the client will bind.")
		theme      = flag.StringP("theme", "t", "default", "Set the color theme that the client will use.")
		logging    = flag.BoolP("log", "l", false, "Enable client logging.")
		encrypt    = flag.BoolP("encrypt", "e", false, "Ignore unencrypted peers.")
		borderless = flag.BoolP("borderless", "b", false, "Do not render an outer border.")
		help       = flag.BoolP("help", "h", false, "Print this help message.")
		vrsn       = flag.BoolP("version", "v", false, "Print version information.")
	)
	flag.Parse()
	args := flag.Args()

	if *help { // print the flag help info
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
		menu.WriteString("\n    " + green.Render("-l") + ", " + green.Render("--log"))
		menu.WriteString("\n        Enable client logging. [dir: " + green.Render("$XDG_STATE_HOME/mabel") + "]")
		menu.WriteString("\n    " + green.Render("-e") + ", " + green.Render("--encrypt"))
		menu.WriteString("\n        Ignore unencrypted peers.")
		menu.WriteString("\n    " + green.Render("-b") + ", " + green.Render("--borderless"))
		menu.WriteString("\n        Do not render an outer border.")
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
	} else if *vrsn { // generate and print the version info
		info, ok := debug.ReadBuildInfo()
		if ok {
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					if commit == "" {
						commit = setting.Value
					}
				case "vcs.modified":
					if setting.Value == "true" {
						modified = " (modified)"
					}
				}
			}
		}

		if version == "" {
			version = info.Main.Version
		}
		if commit != "" {
			commit = fmt.Sprintf("Commit: %s%s\n", commit, modified)
		}
		if builtBy != "" {
			builtBy = fmt.Sprintf("Built by: %s\n", builtBy)
		}
		if date != "" {
			date = fmt.Sprintf("Built on: %s\n", date)
		}

		fmt.Printf(
			"Mabel %s\n%s%s%s",
			version, commit,
			builtBy, date,
		)
		return
	}

	// Check the config file and update download, port, theme, log, and
	// keys based on flags or user configuration. Otherwise, set them
	// to default
	conf := getConfig()
	downloadFlag := flag.Lookup("download")
	portFlag := flag.Lookup("port")
	themeFlag := flag.Lookup("theme")
	loggingFlag := flag.Lookup("log")
	encryptFlag := flag.Lookup("encrypt")
	borderlessFlag := flag.Lookup("borderless")
	thm := conf.getTheme()
	key := conf.Keys

	if !downloadFlag.Changed && conf.Download != "" {
		flag.Set("download", conf.Download)
	}
	if !portFlag.Changed && conf.Port != 0 {
		flag.Set("port", fmt.Sprint(conf.Port))
	}
	if !loggingFlag.Changed {
		flag.Set("log", fmt.Sprint(conf.Log))
	}
	if !encryptFlag.Changed {
		flag.Set("encrypt", fmt.Sprint(conf.RequireEncryption))
	}
	if !borderlessFlag.Changed {
		flag.Set("borderless", fmt.Sprint(conf.Borderless))
	}
	if themeFlag.Changed {
		thm = styles.StringToTheme(theme)
	}

	if *borderless {
		styles.Fullscreen = styles.Fullscreen.BorderStyle(gloss.HiddenBorder())
	}
	styles.Window = styles.Window.BorderForeground(thm.Primary)
	styles.Fullscreen = styles.Fullscreen.BorderForeground(thm.Primary)

	if flag.NArg() == 1 {
		mini.Execute(&args[0], download, port, logging, encrypt, thm)
	} else {
		full.Execute(&args, download, port, logging, encrypt, thm, key)
	}
}
