package styles

// ColorScheme defines all colors used throughout the application
type ColorScheme struct {
	Primary   string `yaml:"primary"`   // Titles, headers, emphasis
	Success   string `yaml:"success"`   // Success messages
	Error     string `yaml:"error"`     // Error messages
	Normal    string `yaml:"normal"`    // Normal text, table data cells
	Muted     string `yaml:"muted"`     // Borders, separators, help text
	Highlight string `yaml:"highlight"` // Selected backgrounds, search match bg
	Accent    string `yaml:"accent"`    // Keywords, strings, relationships, accents
}
