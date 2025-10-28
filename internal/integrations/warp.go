package integrations

import (
	"fmt"
	"os"
	"path/filepath"
	"zakaranda/internal/theme"

	"gopkg.in/yaml.v3"
)

type WarpVariant struct {
	Name        string
	DisplayName string
	ConfigPath  string
	AppPath     string
}

type WarpIntegration struct {
	configPath string
	themesPath string
}

func NewWarpIntegration() *WarpIntegration {
	home, err := os.UserHomeDir()
	if err != nil {
		return &WarpIntegration{configPath: "", themesPath: ""}
	}
	// Both Warp (Default) and Warp Preview share the same themes directory
	themesPath := filepath.Join(home, ".warp", "themes")
	return &WarpIntegration{
		configPath: themesPath,
		themesPath: themesPath,
	}
}

func (w *WarpIntegration) Name() string {
	// Check available variants
	variants := w.GetVariants()

	// If only one variant is installed, show its specific name
	if len(variants) == 1 {
		return variants[0].DisplayName
	}

	// Multiple variants installed, show that themes install to both
	if len(variants) > 1 {
		return "Warp (installs to both versions)"
	}

	// No variants found, show generic name
	return "Warp"
}

func (w *WarpIntegration) ConfigPath() string {
	// Always return the shared themes directory with a note
	// Both Warp (Default) and Warp Preview use ~/.warp/themes
	return w.configPath + " (shared by all Warp versions)"
}

func (w *WarpIntegration) IsInstalled() bool {
	variants := w.GetVariants()
	return len(variants) > 0
}

// GetVariants returns all installed Warp variants
func (w *WarpIntegration) GetVariants() []WarpVariant {
	var variants []WarpVariant
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return variants
	}

	// Both Warp (Default) and Warp Preview share the same themes directory
	sharedThemesPath := filepath.Join(homeDir, ".warp", "themes")

	// Check for Warp (default)
	warpDir := filepath.Join(homeDir, ".warp")
	if _, err := os.Stat(warpDir); err == nil {
		// Check if app exists
		appPaths := []string{
			"/Applications/Warp.app",
			filepath.Join(homeDir, "Applications", "Warp.app"),
		}
		for _, appPath := range appPaths {
			if _, err := os.Stat(appPath); err == nil {
				variants = append(variants, WarpVariant{
					Name:        "warp",
					DisplayName: "Warp (Default)",
					ConfigPath:  sharedThemesPath,
					AppPath:     appPath,
				})
				break
			}
		}
	}

	// Check for Warp Preview
	warpPreviewDir := filepath.Join(homeDir, ".warp-preview")
	if _, err := os.Stat(warpPreviewDir); err == nil {
		// Check if app exists
		appPaths := []string{
			"/Applications/WarpPreview.app",
			filepath.Join(homeDir, "Applications", "WarpPreview.app"),
		}
		for _, appPath := range appPaths {
			if _, err := os.Stat(appPath); err == nil {
				variants = append(variants, WarpVariant{
					Name:        "warp-preview",
					DisplayName: "Warp Preview",
					ConfigPath:  sharedThemesPath,
					AppPath:     appPath,
				})
				break
			}
		}
	}

	return variants
}

func (w *WarpIntegration) Apply(t theme.Theme) error {
	// Use the shared themes directory (both Warp variants use ~/.warp/themes)
	themesPath := w.themesPath

	// Create themes directory if it doesn't exist
	if err := os.MkdirAll(themesPath, 0755); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	// Write theme file
	themeFileName := fmt.Sprintf("%s.yaml", theme.SanitizeFileName(t.Name))
	themePath := filepath.Join(themesPath, themeFileName)

	// Create backup if file exists
	if existingData, err := os.ReadFile(themePath); err == nil && len(existingData) > 0 {
		backupPath := themePath + ".backup"
		if err := os.WriteFile(backupPath, existingData, 0644); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Generate Warp theme
	warpTheme := w.generateWarpTheme(t)

	data, err := yaml.Marshal(warpTheme)
	if err != nil {
		return fmt.Errorf("failed to marshal theme: %w", err)
	}

	if err := os.WriteFile(themePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write theme file: %w", err)
	}

	return nil
}

func (w *WarpIntegration) generateWarpTheme(t theme.Theme) map[string]interface{} {
	// Determine if theme is light or dark based on background color
	details := w.determineThemeDetails(t.Colors.Background)

	// According to official spec, the structure should be:
	// name, accent, cursor (optional), background, foreground, details, terminal_colors
	return map[string]interface{}{
		"name":       t.Name,
		"accent":     t.Colors.Blue,
		"cursor":     t.Colors.Green, // Use green as cursor color for visibility
		"background": t.Colors.Background,
		"foreground": t.Colors.Foreground,
		"details":    details,
		"terminal_colors": map[string]interface{}{
			"bright": map[string]string{
				"black":   t.Colors.BrightBlack,
				"blue":    t.Colors.BrightBlue,
				"cyan":    t.Colors.BrightCyan,
				"green":   t.Colors.BrightGreen,
				"magenta": t.Colors.BrightMagenta,
				"red":     t.Colors.BrightRed,
				"white":   t.Colors.BrightWhite,
				"yellow":  t.Colors.BrightYellow,
			},
			"normal": map[string]string{
				"black":   t.Colors.Black,
				"blue":    t.Colors.Blue,
				"cyan":    t.Colors.Cyan,
				"green":   t.Colors.Green,
				"magenta": t.Colors.Magenta,
				"red":     t.Colors.Red,
				"white":   t.Colors.White,
				"yellow":  t.Colors.Yellow,
			},
		},
	}
}

// determineThemeDetails determines if a theme is "darker" or "lighter" based on background color
func (w *WarpIntegration) determineThemeDetails(bgColor string) string {
	// Remove # if present
	hex := bgColor
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	// Parse hex to RGB
	var r, g, b uint8
	if len(hex) == 6 {
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	}

	// Calculate relative luminance using the formula from WCAG
	// https://www.w3.org/TR/WCAG20/#relativeluminancedef
	rLinear := float64(r) / 255.0
	gLinear := float64(g) / 255.0
	bLinear := float64(b) / 255.0

	// Apply gamma correction
	if rLinear <= 0.03928 {
		rLinear = rLinear / 12.92
	} else {
		rLinear = ((rLinear + 0.055) / 1.055)
		rLinear = rLinear * rLinear // Simplified power of 2.4
	}

	if gLinear <= 0.03928 {
		gLinear = gLinear / 12.92
	} else {
		gLinear = ((gLinear + 0.055) / 1.055)
		gLinear = gLinear * gLinear
	}

	if bLinear <= 0.03928 {
		bLinear = bLinear / 12.92
	} else {
		bLinear = ((bLinear + 0.055) / 1.055)
		bLinear = bLinear * bLinear
	}

	// Calculate luminance
	luminance := 0.2126*rLinear + 0.7152*gLinear + 0.0722*bLinear

	// If luminance is less than 0.5, it's a dark theme
	if luminance < 0.5 {
		return "darker"
	}
	return "lighter"
}
