package styles

import gloss "github.com/charmbracelet/lipgloss"

type ColorTheme struct {
	Primary, Light, Dark, Error,
	Tooltip, GradientStart, GradientEnd gloss.AdaptiveColor
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
		Dark:  "#B2B2B2",
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

var DesertTheme = ColorTheme{
	Primary: gloss.AdaptiveColor{
		Light: "#D08B5D",
		Dark:  "#D08B5D",
	},
	Light: gloss.AdaptiveColor{
		Light: "#E3AF8C",
		Dark:  "#E3AF8C",
	},
	Dark: gloss.AdaptiveColor{
		Light: "#BC662E",
		Dark:  "#BC662E",
	},
	Error: gloss.AdaptiveColor{
		Light: "#5179C5",
		Dark:  "#5179C5",
	},
	Tooltip: gloss.AdaptiveColor{
		Light: "#564F45",
		Dark:  "#B6A589",
	},
	GradientStart: gloss.AdaptiveColor{
		Light: "#D37435",
		Dark:  "#D37435",
	},
	GradientEnd: gloss.AdaptiveColor{
		Light: "#E3AF8C",
		Dark:  "#E3AF8C",
	},
}

type CustomTheme struct {
	Base, Primary, Light, Dark, Error,
	Tooltip, GradientStart, GradientEnd string
}

func stringToTheme(s *string) *ColorTheme {
	switch *s {
	case "desert":
		return &DesertTheme
	default:
		return &DefaultTheme
	}
}

func customToTheme(c *CustomTheme) *ColorTheme {
	theme := stringToTheme(&c.Base)

	if c.Primary != "" {
		theme.Primary = gloss.AdaptiveColor{
			Light: c.Primary,
			Dark:  c.Primary,
		}
	}
	if c.Light != "" {
		theme.Light = gloss.AdaptiveColor{
			Light: c.Light,
			Dark:  c.Light,
		}
	}
	if c.Dark != "" {
		theme.Dark = gloss.AdaptiveColor{
			Light: c.Dark,
			Dark:  c.Dark,
		}
	}
	if c.Error != "" {
		theme.Error = gloss.AdaptiveColor{
			Light: c.Error,
			Dark:  c.Error,
		}
	}
	if c.Tooltip != "" {
		theme.Tooltip = gloss.AdaptiveColor{
			Light: c.Tooltip,
			Dark:  c.Tooltip,
		}
	}
	if c.GradientStart != "" {
		theme.GradientStart = gloss.AdaptiveColor{
			Light: c.GradientStart,
			Dark:  c.GradientStart,
		}
	}
	if c.GradientEnd != "" {
		theme.GradientEnd = gloss.AdaptiveColor{
			Light: c.GradientEnd,
			Dark:  c.GradientEnd,
		}
	}

	return theme
}
