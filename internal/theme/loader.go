package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type ThemeLoader struct {
	builtInThemes []Theme
	customPath    string
}

func NewThemeLoader(customPath string) *ThemeLoader {
	return &ThemeLoader{
		builtInThemes: GetBuiltInThemes(),
		customPath:    customPath,
	}
}

// LoadAllThemes loads built-in and custom themes
func (tl *ThemeLoader) LoadAllThemes() ([]Theme, error) {
	themes := make([]Theme, len(tl.builtInThemes))
	copy(themes, tl.builtInThemes)

	// Load custom themes
	customThemes, err := tl.loadCustomThemes()
	if err != nil {
		// Don't fail if custom themes can't be loaded
		fmt.Printf("Warning: Could not load custom themes: %v\n", err)
	} else {
		themes = append(themes, customThemes...)
	}

	return themes, nil
}

func (tl *ThemeLoader) loadCustomThemes() ([]Theme, error) {
	var themes []Theme

	// Create custom themes directory if it doesn't exist
	if err := os.MkdirAll(tl.customPath, 0755); err != nil {
		return themes, err
	}

	// Read all files in custom themes directory
	entries, err := os.ReadDir(tl.customPath)
	if err != nil {
		return themes, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(tl.customPath, entry.Name())
		ext := strings.ToLower(filepath.Ext(entry.Name()))

		var theme Theme
		var loadErr error

		switch ext {
		case ".json":
			theme, loadErr = tl.loadJSONTheme(filePath)
		case ".yaml", ".yml":
			theme, loadErr = tl.loadYAMLTheme(filePath)
		case ".toml":
			theme, loadErr = tl.loadTOMLTheme(filePath)
		default:
			continue // Skip unsupported formats
		}

		if loadErr != nil {
			fmt.Printf("Warning: Failed to load theme %s: %v\n", entry.Name(), loadErr)
			continue
		}

		themes = append(themes, theme)
	}

	return themes, nil
}

func (tl *ThemeLoader) loadJSONTheme(path string) (Theme, error) {
	var theme Theme

	data, err := os.ReadFile(path)
	if err != nil {
		return theme, err
	}

	if err := json.Unmarshal(data, &theme); err != nil {
		return theme, err
	}

	return theme, nil
}

func (tl *ThemeLoader) loadYAMLTheme(path string) (Theme, error) {
	var theme Theme

	data, err := os.ReadFile(path)
	if err != nil {
		return theme, err
	}

	if err := yaml.Unmarshal(data, &theme); err != nil {
		return theme, err
	}

	return theme, nil
}

func (tl *ThemeLoader) loadTOMLTheme(path string) (Theme, error) {
	var theme Theme

	data, err := os.ReadFile(path)
	if err != nil {
		return theme, err
	}

	// Use the BurntSushi/toml package
	if err := toml.Unmarshal(data, &theme); err != nil {
		return theme, err
	}

	return theme, nil
}

// SaveCustomTheme saves a theme to the custom themes directory
func (tl *ThemeLoader) SaveCustomTheme(theme Theme, format string) error {
	// Ensure directory exists
	if err := os.MkdirAll(tl.customPath, 0755); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	var data []byte
	var err error
	var ext string

	switch strings.ToLower(format) {
	case "json":
		data, err = json.MarshalIndent(theme, "", "  ")
		ext = ".json"
	case "yaml", "yml":
		data, err = yaml.Marshal(theme)
		ext = ".yaml"
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal theme: %w", err)
	}

	filename := SanitizeFileName(theme.Name) + ext
	filePath := filepath.Join(tl.customPath, filename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write theme file: %w", err)
	}

	return nil
}

// ExportTheme exports a built-in theme to a file
func (tl *ThemeLoader) ExportTheme(theme Theme, outputPath string) error {
	ext := strings.ToLower(filepath.Ext(outputPath))
	format := strings.TrimPrefix(ext, ".")

	var data []byte
	var err error

	switch format {
	case "json":
		data, err = json.MarshalIndent(theme, "", "  ")
	case "yaml", "yml":
		data, err = yaml.Marshal(theme)
	default:
		return fmt.Errorf("unsupported format: %s (use .json or .yaml)", format)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal theme: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// BaseTheme represents a theme family with multiple variants
type BaseTheme struct {
	Name        string
	Description string
	Variants    []ThemeVariant
}

// ThemeVariant represents a specific variant of a theme
type ThemeVariant struct {
	Name        string // e.g., "Latte", "Moon", "Nord"
	DisplayName string // e.g., "Latte (Light)", "Moon (Dark)"
	FullName    string // e.g., "Catppuccin Latte", "Rose Pine Moon", "Nord"
	Colors      ColorPalette
}

// GetBuiltInBaseThemes returns all built-in base themes with their variants
func GetBuiltInBaseThemes() []BaseTheme {
	return []BaseTheme{
		{
			Name:        "Nord",
			Description: "An arctic, north-bluish color palette",
			Variants: []ThemeVariant{
				{
					Name:        "Nord",
					DisplayName: "Nord",
					FullName:    "Nord",
					Colors: ColorPalette{
						Background:    "#2e3440", // nord0
						Foreground:    "#d8dee9", // nord4
						Black:         "#3b4252", // nord1
						Red:           "#bf616a", // nord11
						Green:         "#a3be8c", // nord14
						Yellow:        "#ebcb8b", // nord13
						Blue:          "#81a1c1", // nord9
						Magenta:       "#b48ead", // nord15
						Cyan:          "#88c0d0", // nord8
						White:         "#e5e9f0", // nord5
						BrightBlack:   "#4c566a", // nord3
						BrightRed:     "#bf616a", // nord11
						BrightGreen:   "#a3be8c", // nord14
						BrightYellow:  "#ebcb8b", // nord13
						BrightBlue:    "#81a1c1", // nord9
						BrightMagenta: "#b48ead", // nord15
						BrightCyan:    "#8fbcbb", // nord7
						BrightWhite:   "#eceff4", // nord6
					},
				},
			},
		},
		{
			Name:        "Catppuccin",
			Description: "Soothing pastel theme for the high-spirited!",
			Variants: []ThemeVariant{
				{
					Name:        "Latte",
					DisplayName: "Latte (Light)",
					FullName:    "Catppuccin Latte",
					Colors: ColorPalette{
						Background:    "#eff1f5", // Base
						Foreground:    "#4c4f69", // Text
						Black:         "#5c5f77", // Subtext 1
						Red:           "#d20f39", // Red
						Green:         "#40a02b", // Green
						Yellow:        "#df8e1d", // Yellow
						Blue:          "#1e66f5", // Blue
						Magenta:       "#ea76cb", // Pink
						Cyan:          "#179299", // Teal
						White:         "#4c4f69", // Text
						BrightBlack:   "#6c6f85", // Subtext 0
						BrightRed:     "#d20f39", // Red
						BrightGreen:   "#40a02b", // Green
						BrightYellow:  "#df8e1d", // Yellow
						BrightBlue:    "#1e66f5", // Blue
						BrightMagenta: "#ea76cb", // Pink
						BrightCyan:    "#179299", // Teal
						BrightWhite:   "#4c4f69", // Text
					},
				},
				{
					Name:        "Frappe",
					DisplayName: "Frappé (Dark)",
					FullName:    "Catppuccin Frappe",
					Colors: ColorPalette{
						Background:    "#303446", // Base
						Foreground:    "#c6d0f5", // Text
						Black:         "#51576d", // Surface 1
						Red:           "#e78284", // Red
						Green:         "#a6d189", // Green
						Yellow:        "#e5c890", // Yellow
						Blue:          "#8caaee", // Blue
						Magenta:       "#f4b8e4", // Pink
						Cyan:          "#81c8be", // Teal
						White:         "#b5bfe2", // Subtext 1
						BrightBlack:   "#626880", // Surface 2
						BrightRed:     "#e78284", // Red
						BrightGreen:   "#a6d189", // Green
						BrightYellow:  "#e5c890", // Yellow
						BrightBlue:    "#8caaee", // Blue
						BrightMagenta: "#f4b8e4", // Pink
						BrightCyan:    "#81c8be", // Teal
						BrightWhite:   "#a5adce", // Subtext 0
					},
				},
				{
					Name:        "Macchiato",
					DisplayName: "Macchiato (Dark)",
					FullName:    "Catppuccin Macchiato",
					Colors: ColorPalette{
						Background:    "#24273a", // Base
						Foreground:    "#cad3f5", // Text
						Black:         "#494d64", // Surface 1
						Red:           "#ed8796", // Red
						Green:         "#a6da95", // Green
						Yellow:        "#eed49f", // Yellow
						Blue:          "#8aadf4", // Blue
						Magenta:       "#f5bde6", // Pink
						Cyan:          "#8bd5ca", // Teal
						White:         "#b8c0e0", // Subtext 1
						BrightBlack:   "#5b6078", // Surface 2
						BrightRed:     "#ed8796", // Red
						BrightGreen:   "#a6da95", // Green
						BrightYellow:  "#eed49f", // Yellow
						BrightBlue:    "#8aadf4", // Blue
						BrightMagenta: "#f5bde6", // Pink
						BrightCyan:    "#8bd5ca", // Teal
						BrightWhite:   "#a5adcb", // Subtext 0
					},
				},
				{
					Name:        "Mocha",
					DisplayName: "Mocha (Dark)",
					FullName:    "Catppuccin Mocha",
					Colors: ColorPalette{
						Background:    "#1e1e2e", // Base
						Foreground:    "#cdd6f4", // Text
						Black:         "#45475a", // Surface 1
						Red:           "#f38ba8", // Red
						Green:         "#a6e3a1", // Green
						Yellow:        "#f9e2af", // Yellow
						Blue:          "#89b4fa", // Blue
						Magenta:       "#f5c2e7", // Pink
						Cyan:          "#94e2d5", // Teal
						White:         "#bac2de", // Subtext 1
						BrightBlack:   "#585b70", // Surface 2
						BrightRed:     "#f38ba8", // Red
						BrightGreen:   "#a6e3a1", // Green
						BrightYellow:  "#f9e2af", // Yellow
						BrightBlue:    "#89b4fa", // Blue
						BrightMagenta: "#f5c2e7", // Pink
						BrightCyan:    "#94e2d5", // Teal
						BrightWhite:   "#a6adc8", // Subtext 0
					},
				},
			},
		},
		{
			Name:        "Rose Pine",
			Description: "All natural pine, faux fur and a bit of soho vibes",
			Variants: []ThemeVariant{
				{
					Name:        "Main",
					DisplayName: "Main (Dark)",
					FullName:    "Rose Pine",
					Colors: ColorPalette{
						Background:    "#191724", // Base
						Foreground:    "#e0def4", // Text
						Black:         "#26233a", // Overlay
						Red:           "#eb6f92", // Love
						Green:         "#9ccfd8", // Foam
						Yellow:        "#f6c177", // Gold
						Blue:          "#31748f", // Pine
						Magenta:       "#c4a7e7", // Iris
						Cyan:          "#ebbcba", // Rose
						White:         "#e0def4", // Text
						BrightBlack:   "#6e6a86", // Muted
						BrightRed:     "#eb6f92", // Love
						BrightGreen:   "#9ccfd8", // Foam
						BrightYellow:  "#f6c177", // Gold
						BrightBlue:    "#31748f", // Pine
						BrightMagenta: "#c4a7e7", // Iris
						BrightCyan:    "#ebbcba", // Rose
						BrightWhite:   "#e0def4", // Text
					},
				},
				{
					Name:        "Moon",
					DisplayName: "Moon (Dark)",
					FullName:    "Rose Pine Moon",
					Colors: ColorPalette{
						Background:    "#232136", // Base
						Foreground:    "#e0def4", // Text
						Black:         "#393552", // Overlay
						Red:           "#eb6f92", // Love
						Green:         "#9ccfd8", // Foam
						Yellow:        "#f6c177", // Gold
						Blue:          "#3e8fb0", // Pine
						Magenta:       "#c4a7e7", // Iris
						Cyan:          "#ea9a97", // Rose
						White:         "#e0def4", // Text
						BrightBlack:   "#6e6a86", // Muted
						BrightRed:     "#eb6f92", // Love
						BrightGreen:   "#9ccfd8", // Foam
						BrightYellow:  "#f6c177", // Gold
						BrightBlue:    "#3e8fb0", // Pine
						BrightMagenta: "#c4a7e7", // Iris
						BrightCyan:    "#ea9a97", // Rose
						BrightWhite:   "#e0def4", // Text
					},
				},
				{
					Name:        "Dawn",
					DisplayName: "Dawn (Light)",
					FullName:    "Rose Pine Dawn",
					Colors: ColorPalette{
						Background:    "#faf4ed", // Base
						Foreground:    "#575279", // Text
						Black:         "#f2e9e1", // Overlay
						Red:           "#b4637a", // Love
						Green:         "#56949f", // Foam
						Yellow:        "#ea9d34", // Gold
						Blue:          "#286983", // Pine
						Magenta:       "#907aa9", // Iris
						Cyan:          "#d7827e", // Rose
						White:         "#575279", // Text
						BrightBlack:   "#9893a5", // Muted
						BrightRed:     "#b4637a", // Love
						BrightGreen:   "#56949f", // Foam
						BrightYellow:  "#ea9d34", // Gold
						BrightBlue:    "#286983", // Pine
						BrightMagenta: "#907aa9", // Iris
						BrightCyan:    "#d7827e", // Rose
						BrightWhite:   "#575279", // Text
					},
				},
			},
		},
	}
}

var (
	// Cache for built-in themes to avoid rebuilding on every call
	builtInThemesCache []Theme
)

// GetBuiltInThemes returns all built-in themes (flattened from base themes)
func GetBuiltInThemes() []Theme {
	// Return cached themes if available
	if builtInThemesCache != nil {
		return builtInThemesCache
	}
	
	baseThemes := GetBuiltInBaseThemes()
	// Pre-calculate total number of themes to avoid slice reallocations
	totalThemes := 0
	for _, baseTheme := range baseThemes {
		totalThemes += len(baseTheme.Variants)
	}
	
	themes := make([]Theme, 0, totalThemes)
	for _, baseTheme := range baseThemes {
		for _, variant := range baseTheme.Variants {
			themes = append(themes, Theme{
				Name:        variant.FullName,
				Description: baseTheme.Description,
				Colors:      variant.Colors,
			})
		}
	}
	
	// Cache the result
	builtInThemesCache = themes
	return themes
}
