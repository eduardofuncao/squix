package styles

// Built-in color schemes
var (
	// DefaultScheme - Original magenta/cyan theme
	DefaultScheme = ColorScheme{
		Primary:   "205", // Magenta
		Success:   "171", // Purple
		Error:     "196", // Red
		Normal:    "252", // Light gray/white
		Muted:     "238", // Gray
		Highlight: "62",  // Dark Cyan
		Accent:    "86",  // Cyan
	}

	// DraculaScheme - Popular dark theme
	DraculaScheme = ColorScheme{
		Primary:   "141", // Purple
		Success:   "48",  // Green
		Error:     "203", // Orange/Red
		Normal:    "15",  // White
		Muted:     "59",  // Comment Gray
		Highlight: "237", // Dark Blue
		Accent:    "117", // Sky Blue
	}

	// GruvboxDarkScheme - Warm retro colors
	GruvboxDarkScheme = ColorScheme{
		Primary:   "214", // Orange
		Success:   "142", // Green
		Error:     "203", // Red/Orange
		Normal:    "223", // Yellow/white
		Muted:     "244", // Gray
		Highlight: "237", // Dark Gray
		Accent:    "223", // Yellow
	}

	// SolarizedDarkScheme - Precision contrast
	SolarizedDarkScheme = ColorScheme{
		Primary:   "33",  // Blue
		Success:   "106", // Green
		Error:     "203", // Red/Orange
		Normal:    "15",  // White
		Muted:     "244", // Gray
		Highlight: "235", // Dark Gray
		Accent:    "220", // Yellow
	}

	// NordScheme - Arctic, bluish
	NordScheme = ColorScheme{
		Primary:   "68",  // Blue
		Success:   "150", // Green
		Error:     "203", // Red
		Normal:    "15",  // White
		Muted:     "244", // Gray
		Highlight: "236", // Dark Gray
		Accent:    "109", // Aqua
	}

	// MonokaiScheme - Classic vibrant
	MonokaiScheme = ColorScheme{
		Primary:   "141", // Purple
		Success:   "77",  // Green
		Error:     "203", // Red/Orange
		Normal:    "188", // Light gray
		Muted:     "59",  // Gray
		Highlight: "237", // Dark Gray
		Accent:    "81",  // Purple
	}

	// Base16 Schemes

	// BlackMetalScheme - Dark, minimal
	BlackMetalScheme = ColorScheme{
		Primary:   "210", // Pink/red #dd9999
		Success:   "210", // Pink #dd9999
		Error:     "73",  // Cyan #5f8787
		Normal:    "189", // Light gray #c1c1c1
		Muted:     "240", // Dark gray #333333
		Highlight: "210", // Pink/red #dd9999
		Accent:    "144", // Gray #aaaaaa
	}

	// BlackMetalGorgorothScheme - Dark metal variant
	BlackMetalGorgorothScheme = ColorScheme{
		Primary:   "180", // Beige #9b8d7f
		Success:   "180", // Beige #9b8d7f
		Error:     "73",  // Cyan #5f8787
		Normal:    "189", // Light gray #c1c1c1
		Muted:     "240", // Dark gray #333333
		Highlight: "180", // Beige #9b8d7f
		Accent:    "144", // Gray #aaaaaa
	}

	// VesperScheme - Clean dark theme
	VesperScheme = ColorScheme{
		Primary:   "109", // Cyan #8eaaaa
		Success:   "73",  // Cyan #5f8787
		Error:     "167", // Red #de6e6e
		Normal:    "145", // Light gray #b7b7b7
		Muted:     "240", // Dark gray #333333
		Highlight: "235", // Black #222222
		Accent:    "73",  // Cyan #60a592
	}

	// CatppuccinMochaScheme - Modern pastel theme
	CatppuccinMochaScheme = ColorScheme{
		Primary:   "117", // Blue #89b4fa
		Success:   "158", // Green #a6e3a1
		Error:     "210", // Red #f38ba8
		Normal:    "189", // White #cdd6f4
		Muted:     "146", // Subtext #6c7086
		Highlight: "59",  // Surface #45475a
		Accent:    "151", // Teal #94e2d5
	}

	// TokyoNightScheme - Professional dark theme
	TokyoNightScheme = ColorScheme{
		Primary:   "74",  // Cyan #2ac3de
		Success:   "149", // Green #9ece6a
		Error:     "210", // Red #f7768e
		Normal:    "146", // Light gray #a9b1d6
		Muted:     "147", // Muted #787c99
		Highlight: "23",  // Dark blue #2f3549
		Accent:    "153", // Cyan #b4f9f8
	}

	// RosePineScheme - Soft purple/pink
	RosePineScheme = ColorScheme{
		Primary:   "182", // Purple #c4a7e7
		Success:   "31",  // Blue #31748f
		Error:     "168", // Pink #eb6f92
		Normal:    "188", // White #e0def4
		Muted:     "97",  // Muted purple #6e6a86
		Highlight: "54",  // Dark purple #26233a
		Accent:    "152", // Cyan #9ccfd8
	}

	// TerracottaScheme - Light earth tones
	TerracottaScheme = ColorScheme{
		Primary:   "97",  // Purple #625574
		Success:   "107", // Green #7a894a
		Error:     "131", // Red #a75045
		Normal:    "59",  // Brown #473731
		Muted:     "138", // Light brown #c0aca4
		Highlight: "181", // Light beige #d0c1bb
		Accent:    "103", // Purple #847f9e
	}
)

// GetScheme returns a color scheme by name
func GetScheme(name string) ColorScheme {
	switch name {
	case "dracula":
		return DraculaScheme
	case "gruvbox":
		return GruvboxDarkScheme
	case "solarized":
		return SolarizedDarkScheme
	case "nord":
		return NordScheme
	case "monokai":
		return MonokaiScheme
	case "black-metal":
		return BlackMetalScheme
	case "black-metal-gorgoroth":
		return BlackMetalGorgorothScheme
	case "vesper":
		return VesperScheme
	case "catppuccin-mocha":
		return CatppuccinMochaScheme
	case "tokyo-night":
		return TokyoNightScheme
	case "rose-pine":
		return RosePineScheme
	case "terracotta":
		return TerracottaScheme
	default:
		return DefaultScheme
	}
}
