package stats

import (
	"fmt"

	"github.com/anacrolix/torrent"
)

// Peers retrieves the peers information for a torrent and returns it
// as a string.
func Peers(t *torrent.Torrent) string {
	stats := t.Stats()

	return fmt.Sprintf(
		"%d/%d peers",
		stats.ActivePeers,
		stats.TotalPeers,
	)
}
