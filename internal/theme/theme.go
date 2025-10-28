package theme

// Theme represents a color theme with a name, description, and color palette
type Theme struct {
	Name        string
	Description string
	Colors      ColorPalette
}

// ColorPalette defines the color scheme for a theme
type ColorPalette struct {
	Background    string
	Foreground    string
	Black         string
	Red           string
	Green         string
	Yellow        string
	Blue          string
	Magenta       string
	Cyan          string
	White         string
	BrightBlack   string
	BrightRed     string
	BrightGreen   string
	BrightYellow  string
	BrightBlue    string
	BrightMagenta string
	BrightCyan    string
	BrightWhite   string
}

