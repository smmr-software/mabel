package stats

import (
	"fmt"

	"github.com/anacrolix/torrent"
)

func Peers(t *torrent.Torrent) string {
	stats := t.Stats()

	return fmt.Sprintf(
		"%d/%d peers",
		stats.ActivePeers,
		stats.TotalPeers,
	)
}
