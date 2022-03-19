package shared

import (
	"fmt"

	"github.com/anacrolix/torrent"
	"github.com/dustin/go-humanize"
)

func DownloadStats(t *torrent.Torrent, showPercent bool) string {
	var (
		done    = t.BytesCompleted()
		total   = t.Length()
		percent = float64(done) / float64(total) * 100

		tail string
	)

	if showPercent {
		tail = fmt.Sprintf(" (%d%%)", uint64(percent))
	}

	return fmt.Sprintf(
		"%s/%s%s ↓",
		humanize.Bytes(uint64(done)),
		humanize.Bytes(uint64(total)),
		tail,
	)
}
