package styles

import gloss "github.com/charmbracelet/lipgloss"

type ColorTheme struct {
	Primary       gloss.AdaptiveColor
	Light         gloss.AdaptiveColor
	Dark          gloss.AdaptiveColor
	Error         gloss.AdaptiveColor
	Tooltip       gloss.AdaptiveColor
	GradientStart gloss.AdaptiveColor
	GradientEnd   gloss.AdaptiveColor
}

var DefaultTheme = ColorTheme{
	Primary: gloss.AdaptiveColor{
		Light: "#464FB6",
		Dark:  "#5FFFD7",
	},
	Light: gloss.AdaptiveColor{
		Light: "#9EA4D0",
		Dark:  "#DFFFF7",
	},
	Dark: gloss.AdaptiveColor{
		Light: "#353B88",
		Dark:  "#00AC81",
	},
	Error: gloss.AdaptiveColor{
		Light: "#FF5E87",
		Dark:  "#FF5E87",
	},
	Tooltip: gloss.AdaptiveColor{
		Light: "#4A4A4A",
		Dark:  "B2B2B2",
	},
	GradientStart: gloss.AdaptiveColor{
		Light: "#5A56E0",
		Dark:  "#5A56E0",
	},
	GradientEnd: gloss.AdaptiveColor{
		Light: "#EE6FF8",
		Dark:  "#EE6FF8",
	},
}

var PlaceholderTheme = ColorTheme{
	Primary: gloss.AdaptiveColor{
		Light: "#464FB6",
		Dark:  "#F456CA",
	},
	Light: gloss.AdaptiveColor{
		Light: "#9EA4D0",
		Dark:  "#DFFFF7",
	},
	Dark: gloss.AdaptiveColor{
		Light: "#353B88",
		Dark:  "#00AC81",
	},
	Error: gloss.AdaptiveColor{
		Light: "#FF5E87",
		Dark:  "#FF5E87",
	},
	Tooltip: gloss.AdaptiveColor{
		Light: "#4A4A4A",
		Dark:  "B2B2B2",
	},
	GradientStart: gloss.AdaptiveColor{
		Light: "#5A56E0",
		Dark:  "#CC9254",
	},
	GradientEnd: gloss.AdaptiveColor{
		Light: "#EE6FF8",
		Dark:  "#B854CC",
	},
}
