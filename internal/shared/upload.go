package shared

import (
	"fmt"

	"github.com/anacrolix/torrent"
	"github.com/dustin/go-humanize"
)

func UploadStats(t *torrent.Torrent) string {
	var (
		stats  = t.Stats()
		upload = stats.BytesWritten.Int64()
	)

	return fmt.Sprintf(
		"%s ↑",
		humanize.Bytes(uint64(upload)),
	)
}
