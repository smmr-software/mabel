# Mabel

Deriving its name from the Hebrew word "מבול," meaning flood, deluge, or
(loosely) torrent, Mabel is a fancy BitTorrent client for the terminal.

Mabel builds on the excellent work of others in the open source
community, including:

- [charmbracelet/bubbletea]
- [anacrolix/torrent]

## Prerequisites

- A [Nerd Font] installed and enabled in your terminal.

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

## License

[GPLv3](COPYING).

***

A [SMMR Software] creation. 🏖

[charmbracelet/bubbletea]: https://github.com/charmbracelet/bubbletea
[anacrolix/torrent]: https://github.com/anacrolix/torrent
[Nerd Font]: https://www.nerdfonts.com
[SMMR Software]: https://smmr.software/
