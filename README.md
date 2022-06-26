# Mabel

Deriving its name from the Hebrew word "מבול," meaning flood, deluge, or
(loosely) torrent, Mabel is a fancy BitTorrent client for the terminal.

Mabel builds on the excellent work of others in the open source
community, including:

- [charmbracelet/bubbletea]
- [anacrolix/torrent]

## Screenshots

<img width="50%" src="/default.png" alt="Mabel downloading several torrents, default theme">
<img width="50%" src="/desert.png" alt="Mabel looking at torrent details, desert theme">

## Prerequisites

- A [Nerd Font] installed and enabled in your terminal

## Install

Coming soon to a package manager near you!

```sh
go install github.com/smmr-software/mabel@latest
```

## Usage

```
mabel [OPTIONS] [TORRENT]...
```

When run without any arguments, Mabel starts a full TUI client.

When passed a single torrent (infohash, magnet link, `.torrent` file),
Mabel starts in "mini" mode.

When multiple torrents are provided as arguments, Mabel opens the full
TUI client with the corresponding torrents added.

### mini

The mini client downloads a single torrent, displaying its name and
download stats. Press `q` or `ctrl+c` to close the client.

### full

The full client manages a list of torrents downloading in parallel.

A brief breakdown of the controls:

- Scroll through the list with `hjkl` or the arrow keys
- Press `a` to add a torrent
- Get more details about a torrent with `return`
- Perhaps most importantly – press `?` to expand the help menu at the
  bottom of the view

Once again, press `q` or `ctrl+c` to close the client.

## Config and Flags

Mabel can be configured via a TOML file or via flags at runtime. Flags
take precedence over the config file.

### Config File

Located at `$XDG_CONFIG_HOME/mabel/config.toml`. A basic example:
```toml
# The default torrent download directory.
# Supports expansion of ~ characters
# Default: $XDG_DOWNLOAD_DIR
download = "~/movies/buffer"

# The port to which the client will bind
# Default: 42069
port = 126

# Toggle client logging (logs are written to $XDG_STATE_HOME/mabel)
log = false

# Mabel always prefers encrypted connections. If set to true,
# require_encryption will have Mabel ignore unencrypted peers.
require_encryption = false
```

#### Theme

The `theme` key is special, as it can be one of two types.

As a string, `theme` selects one of our default themes. The currently
available themes are:
- default
- desert
- purple
- 8-bit/ansi

As a table [a.k.a (hash)map or dictionary], the `theme` key can also
allow you to customize your colors in-depth. The `theme.base` key
provides a fallback for any unset values in the table, and follows the
same rules as the string `theme` key.

The `theme.gradient*` keys customize the gradients used in progress bars
throughout Mabel. `gradient-solid` takes precedence over
`gradient-start` and `gradient-end`. The latter two only work with
24-bit color.

Some examples:
```toml
# just the desert theme
theme = "desert"

# default, with a red primary color
[theme]
base = "default"
primary = "#FF0000"

# the same as previous
[theme]
primary = "#FF0000"

# 8-bit, with a blue error color
[theme]
base = "8-bit"
error = "12"

# default, with a gradient from Pink Lightly Toasted to SMMR Software
# Example Color 2 (https://colornames.org/color/990d35)
[theme]
gradient-start = "#D52941"
gradient-end = "#990D35"

# a completely custom theme where everything is green
[theme]
primary = "#00FF00"
light = "#00FF00"
dark = "#00FF00"
error = "#00FF00"
tooltip = "#00FF00"
gradient-solid = "#00FF00"
```

### Flags

Runtime flags are documented in Mabel's help message, which you can view
by passing `-h` on invocation. Flags encompass all the options
configurable in the TOML config, plus help and version information
messages.

## License

[GPLv3](COPYING).

***

A [SMMR Software] creation. 🏖

[charmbracelet/bubbletea]: https://github.com/charmbracelet/bubbletea
[anacrolix/torrent]: https://github.com/anacrolix/torrent
[Nerd Font]: https://www.nerdfonts.com
[SMMR Software]: https://smmr.software/
