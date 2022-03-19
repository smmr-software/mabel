package shared

import (
	"fmt"

	"github.com/anacrolix/torrent"
)

func PeerStats(t *torrent.Torrent) string {
	stats := t.Stats()

	return fmt.Sprintf(
		"%d/%d peers",
		stats.ActivePeers,
		stats.TotalPeers,
	)
}
