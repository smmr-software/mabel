package stats

import (
	"fmt"

	"github.com/anacrolix/torrent"
	"github.com/dustin/go-humanize"
)

func Upload(t *torrent.Torrent) string {
	var (
		stats  = t.Stats()
		upload = stats.BytesWritten.Int64()
	)

	return fmt.Sprintf(
		"%s ↑",
		humanize.Bytes(uint64(upload)),
	)
}
