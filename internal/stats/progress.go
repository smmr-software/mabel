package stats

import (
	"github.com/smmr-software/mabel/internal/styles"

	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/progress"
)

// ProgressBar retrieves the progress information and defines the
// styles for the Bubbles progress bar.
func ProgressBar(t *torrent.Torrent, width *int, theme *styles.ColorTheme) string {
	var (
		done    = t.BytesCompleted()
		total   = t.Length()
		percent = float64(done) / float64(total)
	)

	var gradient progress.Option
	if theme.UseSolidGradient() {
		gradient = progress.WithSolidFill(
			styles.AdaptiveColorToString(&theme.GradientSolid))
	} else {
		gradient = progress.WithGradient(
			styles.AdaptiveColorToString(&theme.GradientStart),
			styles.AdaptiveColorToString(&theme.GradientEnd),
		)
	}

	prog := progress.New(gradient, progress.WithoutPercentage())

	if width != nil {
		prog.Width = *width
	}

	return prog.ViewAs(percent)
}
