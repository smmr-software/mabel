package styles

import gloss "github.com/charmbracelet/lipgloss"

type ColorTheme struct {
	Primary, Light, Dark, Error,
	Tooltip, GradientStart,
	GradientEnd, GradientSolid gloss.AdaptiveColor
}

type CustomTheme struct {
	Base, Primary, Light,
	Dark, Error, Tooltip string

	GradientStart string `toml:"gradient-start"`
	GradientEnd   string `toml:"gradient-end"`
	GradientSolid string `toml:"gradient-solid"`
}

// StringToTheme converts strings from the config file into ColorTheme
// pointers.
func StringToTheme(s *string) *ColorTheme {
	switch *s {
	case "desert":
		return &DesertTheme
	case "purple", "lean", "drank":
		return &PurpleTheme
	case "8-bit", "ansi":
		return &ANSITheme
	default:
		return &DefaultTheme
	}
}

// ToTheme translates a CustomTheme of strings to a ColorTheme of Lip
// Gloss AdaptiveColors.
func (c *CustomTheme) ToTheme() *ColorTheme {
	theme := StringToTheme(&c.Base)

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
	if c.GradientSolid != "" {
		theme.GradientSolid = gloss.AdaptiveColor{
			Light: c.GradientSolid,
			Dark:  c.GradientSolid,
		}
	}

	return theme
}

// UseSolidGradient returns a boolean based on whether the gradient for
// the progress bar is solid or not.
func (t *ColorTheme) UseSolidGradient() bool {
	return t.GradientSolid != gloss.AdaptiveColor{}
}

// AdaptiveColorToString returns the Light or Dark alternative of an
// AdaptiveColor as a string, based on the background color of the
// terminal.
func AdaptiveColorToString(c *gloss.AdaptiveColor) string {
	if gloss.HasDarkBackground() {
		return c.Dark
	}
	return c.Light
}

// Define the default theme
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

// Define the desert theme
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

// Define the purple theme
var PurpleTheme = ColorTheme{
	Primary: gloss.AdaptiveColor{
		Light: "#8C00FF",
		Dark:  "#8C00FF",
	},
	Light: gloss.AdaptiveColor{
		Light: "#8C00FF",
		Dark:  "#8C00FF",
	},
	Dark: gloss.AdaptiveColor{
		Light: "#8C00FF",
		Dark:  "#8C00FF",
	},
	Error: gloss.AdaptiveColor{
		Light: "#8C00FF",
		Dark:  "#8C00FF",
	},
	Tooltip: gloss.AdaptiveColor{
		Light: "#8C00FF",
		Dark:  "#8C00FF",
	},
	GradientSolid: gloss.AdaptiveColor{
		Light: "#8C00FF",
		Dark:  "#8C00FF",
	},
}

// Define the ANSI/8-bit theme
var ANSITheme = ColorTheme{
	Primary:       gloss.AdaptiveColor{Light: "4", Dark: "4"},
	Light:         gloss.AdaptiveColor{Light: "14", Dark: "14"},
	Dark:          gloss.AdaptiveColor{Light: "12", Dark: "12"},
	Error:         gloss.AdaptiveColor{Light: "9", Dark: "9"},
	Tooltip:       gloss.AdaptiveColor{Light: "8", Dark: "8"},
	GradientSolid: gloss.AdaptiveColor{Light: "5", Dark: "5"},
}
