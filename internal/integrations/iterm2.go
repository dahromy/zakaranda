package integrations

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"zakaranda/internal/theme"
)

type ITerm2Integration struct {
	themesPath string
}

func NewITerm2Integration() *ITerm2Integration {
	home, err := os.UserHomeDir()
	if err != nil {
		return &ITerm2Integration{themesPath: ""}
	}
	themesPath := filepath.Join(home, ".config", "theme-manager", "iterm2")
	return &ITerm2Integration{themesPath: themesPath}
}

func (i *ITerm2Integration) Name() string {
	return "iTerm2"
}

func (i *ITerm2Integration) ConfigPath() string {
	return i.themesPath
}

func (i *ITerm2Integration) IsInstalled() bool {
	// Check multiple possible iTerm2 installation locations
	possiblePaths := []string{
		"/Applications/iTerm.app",
		"/Applications/iTerm2.app",
	}

	// Also check user's Applications folder
	if home, err := os.UserHomeDir(); err == nil {
		possiblePaths = append(possiblePaths,
			filepath.Join(home, "Applications", "iTerm.app"),
			filepath.Join(home, "Applications", "iTerm2.app"),
		)
	}

	// Check if any of the paths exist
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (i *ITerm2Integration) Apply(t theme.Theme) error {
	// Create themes directory if it doesn't exist
	if err := os.MkdirAll(i.themesPath, 0755); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	// Generate iTerm2 color preset
	preset := i.generateITerm2Preset(t)
	presetFileName := fmt.Sprintf("%s.itermcolors", theme.SanitizeFileName(t.Name))
	presetPath := filepath.Join(i.themesPath, presetFileName)

	// Create backup if file exists
	if existingData, err := os.ReadFile(presetPath); err == nil && len(existingData) > 0 {
		backupPath := presetPath + ".backup"
		if err := os.WriteFile(backupPath, existingData, 0644); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	if err := os.WriteFile(presetPath, []byte(preset), 0644); err != nil {
		return fmt.Errorf("failed to write preset file: %w", err)
	}

	// Try to import the theme using PlistBuddy (official method)
	if err := i.importTheme(presetPath); err != nil {
		return fmt.Errorf("theme saved to %s, but auto-import failed: %w\n\nTo import manually:\n1. Open iTerm2 → Preferences → Profiles → Colors\n2. Click 'Color Presets' → 'Import'\n3. Select: %s\n4. Restart iTerm2 to see the theme", presetPath, err, presetPath)
	}

	return nil
}

func (i *ITerm2Integration) generateITerm2Preset(t theme.Theme) string {
	// iTerm2 uses XML plist format for color schemes
	// Format follows the official iTerm2 Color Schemes specification
	template := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Ansi 0 Color</key>%s
	<key>Ansi 1 Color</key>%s
	<key>Ansi 2 Color</key>%s
	<key>Ansi 3 Color</key>%s
	<key>Ansi 4 Color</key>%s
	<key>Ansi 5 Color</key>%s
	<key>Ansi 6 Color</key>%s
	<key>Ansi 7 Color</key>%s
	<key>Ansi 8 Color</key>%s
	<key>Ansi 9 Color</key>%s
	<key>Ansi 10 Color</key>%s
	<key>Ansi 11 Color</key>%s
	<key>Ansi 12 Color</key>%s
	<key>Ansi 13 Color</key>%s
	<key>Ansi 14 Color</key>%s
	<key>Ansi 15 Color</key>%s
	<key>Background Color</key>%s
	<key>Badge Color</key>%s
	<key>Bold Color</key>%s
	<key>Cursor Color</key>%s
	<key>Cursor Guide Color</key>%s
	<key>Cursor Text Color</key>%s
	<key>Foreground Color</key>%s
	<key>Link Color</key>%s
	<key>Selected Text Color</key>%s
	<key>Selection Color</key>%s
</dict>
</plist>`

	return fmt.Sprintf(template,
		i.hexToITermColor(t.Colors.Black),               // Ansi 0
		i.hexToITermColor(t.Colors.Red),                 // Ansi 1
		i.hexToITermColor(t.Colors.Green),               // Ansi 2
		i.hexToITermColor(t.Colors.Yellow),              // Ansi 3
		i.hexToITermColor(t.Colors.Blue),                // Ansi 4
		i.hexToITermColor(t.Colors.Magenta),             // Ansi 5
		i.hexToITermColor(t.Colors.Cyan),                // Ansi 6
		i.hexToITermColor(t.Colors.White),               // Ansi 7
		i.hexToITermColor(t.Colors.BrightBlack),         // Ansi 8
		i.hexToITermColor(t.Colors.BrightRed),           // Ansi 9
		i.hexToITermColor(t.Colors.BrightGreen),         // Ansi 10
		i.hexToITermColor(t.Colors.BrightYellow),        // Ansi 11
		i.hexToITermColor(t.Colors.BrightBlue),          // Ansi 12
		i.hexToITermColor(t.Colors.BrightMagenta),       // Ansi 13
		i.hexToITermColor(t.Colors.BrightCyan),          // Ansi 14
		i.hexToITermColor(t.Colors.BrightWhite),         // Ansi 15
		i.hexToITermColor(t.Colors.Background),          // Background
		i.hexToITermColorWithAlpha(t.Colors.Black, 0.5), // Badge (semi-transparent)
		i.hexToITermColor(t.Colors.BrightWhite),         // Bold
		i.hexToITermColor(t.Colors.Foreground),          // Cursor
		i.hexToITermColor(t.Colors.BrightBlack),         // Cursor Guide (subtle)
		i.hexToITermColor(t.Colors.Background),          // Cursor Text
		i.hexToITermColor(t.Colors.Foreground),          // Foreground
		i.hexToITermColor(t.Colors.BrightCyan),          // Link (bright cyan for visibility)
		i.hexToITermColor(t.Colors.BrightBlack),         // Selected Text
		i.hexToITermColor(t.Colors.Foreground),          // Selection
	)
}

func (i *ITerm2Integration) hexToITermColor(hex string) string {
	return i.hexToITermColorWithAlpha(hex, 1.0)
}

func (i *ITerm2Integration) hexToITermColorWithAlpha(hex string, alpha float64) string {
	// Remove # if present
	hex = strings.TrimPrefix(hex, "#")

	// Parse hex to RGB
	r, g, b := i.hexToRGB(hex)

	// Convert to iTerm2 format (0.0 - 1.0)
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	return fmt.Sprintf(`
	<dict>
		<key>Alpha Component</key>
		<real>%.6f</real>
		<key>Blue Component</key>
		<real>%.6f</real>
		<key>Color Space</key>
		<string>sRGB</string>
		<key>Green Component</key>
		<real>%.6f</real>
		<key>Red Component</key>
		<real>%.6f</real>
	</dict>`, alpha, bf, gf, rf)
}

func (i *ITerm2Integration) hexToRGB(hex string) (uint8, uint8, uint8) {
	var r, g, b uint8
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

func (i *ITerm2Integration) importTheme(presetPath string) error {
	// Use PlistBuddy to import the theme directly into iTerm2's preferences
	// This method is based on the official iTerm2-Color-Schemes import script
	// Reference: https://github.com/mbadolato/iTerm2-Color-Schemes/blob/master/tools/import-scheme.sh

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	plistPath := filepath.Join(home, "Library", "Preferences", "com.googlecode.iterm2.plist")

	// Check if iTerm2 preferences file exists
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		return fmt.Errorf("iTerm2 preferences file not found at %s. Please launch iTerm2 at least once", plistPath)
	}

	// Extract theme name from file path
	themeName := i.extractThemeName(presetPath)

	// Create 'Custom Color Presets' entry if it doesn't exist
	if err := i.ensureCustomColorPresetsExists(plistPath); err != nil {
		return fmt.Errorf("failed to ensure Custom Color Presets exists: %w", err)
	}

	// Check if theme already exists
	themeExists, err := i.themeExists(plistPath, themeName)
	if err != nil {
		return fmt.Errorf("failed to check if theme exists: %w", err)
	}

	// Import the theme
	if themeExists {
		// Delete existing theme first, then reinstall
		if err := i.deleteTheme(plistPath, themeName); err != nil {
			return fmt.Errorf("failed to delete existing theme: %w", err)
		}
	}

	// Add and merge the theme
	if err := i.addAndMergeTheme(plistPath, themeName, presetPath); err != nil {
		return fmt.Errorf("failed to import theme: %w", err)
	}

	return nil
}

// extractThemeName extracts and formats the theme name from the file path
func (i *ITerm2Integration) extractThemeName(filePath string) string {
	// Get base name without extension
	baseName := filepath.Base(filePath)
	name := strings.TrimSuffix(baseName, ".itermcolors")

	// Replace underscores and hyphens with spaces
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	return name
}

// ensureCustomColorPresetsExists creates the 'Custom Color Presets' entry if it doesn't exist
func (i *ITerm2Integration) ensureCustomColorPresetsExists(plistPath string) error {
	// Check if 'Custom Color Presets' exists
	cmd := exec.Command("/usr/libexec/PlistBuddy",
		"-c", "Print \"Custom Color Presets\"",
		plistPath)

	if err := cmd.Run(); err != nil {
		// Entry doesn't exist, create it
		cmd = exec.Command("/usr/libexec/PlistBuddy",
			"-c", "Add \"Custom Color Presets\" dict",
			plistPath)

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create Custom Color Presets entry: %w", err)
		}
	}

	return nil
}

// themeExists checks if a theme with the given name already exists
func (i *ITerm2Integration) themeExists(plistPath, themeName string) (bool, error) {
	cmd := exec.Command("/usr/libexec/PlistBuddy",
		"-c", fmt.Sprintf("Print \"Custom Color Presets:%s\"", themeName),
		plistPath)

	err := cmd.Run()
	if err != nil {
		// Theme doesn't exist
		return false, nil
	}

	return true, nil
}

// deleteTheme deletes an existing theme
func (i *ITerm2Integration) deleteTheme(plistPath, themeName string) error {
	cmd := exec.Command("/usr/libexec/PlistBuddy",
		"-c", fmt.Sprintf("Delete \"Custom Color Presets:%s\"", themeName),
		plistPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete theme: %w", err)
	}

	return nil
}

// addAndMergeTheme adds a new theme entry and merges the color scheme
func (i *ITerm2Integration) addAndMergeTheme(plistPath, themeName, presetPath string) error {
	cmd := exec.Command("/usr/libexec/PlistBuddy",
		"-c", fmt.Sprintf("Add \"Custom Color Presets:%s\" dict", themeName),
		"-c", fmt.Sprintf("Merge \"%s\" \"Custom Color Presets:%s\"", presetPath, themeName),
		plistPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add and merge theme: %w\nOutput: %s", err, string(output))
	}

	return nil
}
