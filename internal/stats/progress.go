package stats

import (
	"github.com/smmr-software/mabel/internal/styles"

	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/progress"
	gloss "github.com/charmbracelet/lipgloss"
)

func ProgressBar(t *torrent.Torrent, width *int, theme *styles.ColorTheme) string {
	var (
		done    = t.BytesCompleted()
		total   = t.Length()
		percent = float64(done) / float64(total)

		gradientStart = theme.GradientStart.Dark
		gradientEnd   = theme.GradientEnd.Dark
	)

	if !gloss.HasDarkBackground() {
		gradientStart = theme.GradientStart.Light
		gradientEnd = theme.GradientEnd.Light
	}

	prog := progress.New(
		progress.WithGradient(gradientStart, gradientEnd),
		progress.WithoutPercentage(),
	)

	if width != nil {
		prog.Width = *width
	}

	return prog.ViewAs(percent)
}
