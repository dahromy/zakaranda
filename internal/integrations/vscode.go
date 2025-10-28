package integrations

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"zakaranda/internal/theme"
)

type VSCodeIntegration struct {
	configPath string
	variant    VSCodeVariant
}

// VSCodeVariant represents different VS Code variants
type VSCodeVariant struct {
	Name       string
	ConfigDir  string
	CLICommand string
	AppPath    string
}

// VSCodeThemeExtension represents a VS Code theme extension
type VSCodeThemeExtension struct {
	ExtensionID      string
	ThemeName        string
	IconTheme        string
	ProductIconTheme string
}

// Map theme names to VS Code extensions
var vscodeThemeExtensions = map[string]VSCodeThemeExtension{
	"Nord": {
		ExtensionID:      "arcticicestudio.nord-visual-studio-code",
		ThemeName:        "Nord",
		IconTheme:        "charmed-light", // Requires: charmed-icons.charmed-icons
		ProductIconTheme: "fluent-icons",  // Requires: miguelsolorio.fluent-icons
	},
	"Catppuccin Latte": {
		ExtensionID:      "catppuccin.catppuccin-vsc",
		ThemeName:        "Catppuccin Latte",
		IconTheme:        "catppuccin-latte",
		ProductIconTheme: "catppuccin-latte",
	},
	"Catppuccin Frappe": {
		ExtensionID:      "catppuccin.catppuccin-vsc",
		ThemeName:        "Catppuccin Frappé",
		IconTheme:        "catppuccin-frappe",
		ProductIconTheme: "catppuccin-frappe",
	},
	"Catppuccin Macchiato": {
		ExtensionID:      "catppuccin.catppuccin-vsc",
		ThemeName:        "Catppuccin Macchiato",
		IconTheme:        "catppuccin-macchiato",
		ProductIconTheme: "catppuccin-macchiato",
	},
	"Catppuccin Mocha": {
		ExtensionID:      "catppuccin.catppuccin-vsc",
		ThemeName:        "Catppuccin Mocha",
		IconTheme:        "catppuccin-mocha",
		ProductIconTheme: "catppuccin-mocha",
	},
	"Rose Pine": {
		ExtensionID:      "mvllow.rose-pine",
		ThemeName:        "Rosé Pine",
		IconTheme:        "rose-pine-icons",
		ProductIconTheme: "fluent-icons",
	},
	"Rose Pine Moon": {
		ExtensionID:      "mvllow.rose-pine",
		ThemeName:        "Rosé Pine Moon",
		IconTheme:        "rose-pine-moon-icons",
		ProductIconTheme: "fluent-icons",
	},
	"Rose Pine Dawn": {
		ExtensionID:      "mvllow.rose-pine",
		ThemeName:        "Rosé Pine Dawn",
		IconTheme:        "rose-pine-dawn-icons",
		ProductIconTheme: "fluent-icons",
	},
}

// Additional extensions needed for icon themes
var iconThemeExtensions = map[string]string{
	"charmed-light":       "charmed-icons.charmed-icons",
	"material-icon-theme": "pkief.material-icon-theme",
	"fluent-icons":        "miguelsolorio.fluent-icons",
	"catppuccin-mocha":    "catppuccin.catppuccin-vsc", // Included in main theme
	"rose-pine-icons":     "mvllow.rose-pine",          // Included in main theme
}

// GetVSCodeVariants returns all available VS Code variants on the system
func GetVSCodeVariants() []VSCodeVariant {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	variants := []VSCodeVariant{
		{
			Name:       "VS Code",
			ConfigDir:  filepath.Join(home, "Library", "Application Support", "Code"),
			CLICommand: "code",
			AppPath:    "/Applications/Visual Studio Code.app",
		},
		{
			Name:       "VS Code Insiders",
			ConfigDir:  filepath.Join(home, "Library", "Application Support", "Code - Insiders"),
			CLICommand: "code-insiders",
			AppPath:    "/Applications/Visual Studio Code - Insiders.app",
		},
		{
			Name:       "Cursor",
			ConfigDir:  filepath.Join(home, "Library", "Application Support", "Cursor"),
			CLICommand: "cursor",
			AppPath:    "/Applications/Cursor.app",
		},
	}

	// Filter to only installed variants
	var installed []VSCodeVariant
	for _, variant := range variants {
		// Check if config directory exists OR app is installed
		configExists := false
		if _, err := os.Stat(variant.ConfigDir); err == nil {
			configExists = true
		}

		appExists := false
		if _, err := os.Stat(variant.AppPath); err == nil {
			appExists = true
		}

		if configExists || appExists {
			installed = append(installed, variant)
		}
	}

	return installed
}

// stripJSONComments removes single-line (//) and multi-line (/* */) comments from JSON
// while preserving strings that might contain // or /* */
func stripJSONComments(jsonStr string) string {
	var result strings.Builder
	inString := false
	inSingleLineComment := false
	inMultiLineComment := false
	escaped := false

	for i := 0; i < len(jsonStr); i++ {
		ch := jsonStr[i]

		// Handle escape sequences in strings
		if inString {
			result.WriteByte(ch)
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == '"' {
				inString = false
			}
			continue
		}

		// Check for string start
		if ch == '"' && !inSingleLineComment && !inMultiLineComment {
			inString = true
			result.WriteByte(ch)
			continue
		}

		// Handle multi-line comment end
		if inMultiLineComment {
			if ch == '*' && i+1 < len(jsonStr) && jsonStr[i+1] == '/' {
				inMultiLineComment = false
				i++ // Skip the '/'
			}
			continue
		}

		// Handle single-line comment end
		if inSingleLineComment {
			if ch == '\n' || ch == '\r' {
				inSingleLineComment = false
				result.WriteByte(ch) // Preserve newline
			}
			continue
		}

		// Check for comment start
		if ch == '/' && i+1 < len(jsonStr) {
			next := jsonStr[i+1]
			if next == '/' {
				inSingleLineComment = true
				i++ // Skip the second '/'
				continue
			} else if next == '*' {
				inMultiLineComment = true
				i++ // Skip the '*'
				continue
			}
		}

		// Regular character
		result.WriteByte(ch)
	}

	return result.String()
}

func NewVSCodeIntegration() *VSCodeIntegration {
	home, err := os.UserHomeDir()
	if err != nil {
		return &VSCodeIntegration{configPath: ""}
	}

	// Default to standard VS Code
	defaultVariant := VSCodeVariant{
		Name:       "VS Code",
		ConfigDir:  filepath.Join(home, "Library", "Application Support", "Code"),
		CLICommand: "code",
		AppPath:    "/Applications/Visual Studio Code.app",
	}

	configPath := filepath.Join(defaultVariant.ConfigDir, "User", "settings.json")
	return &VSCodeIntegration{
		configPath: configPath,
		variant:    defaultVariant,
	}
}

func (v *VSCodeIntegration) Name() string {
	return "VS Code"
}

func (v *VSCodeIntegration) ConfigPath() string {
	// Don't show config path since we have variant selection
	return ""
}

func (v *VSCodeIntegration) IsInstalled() bool {
	// Check if any VS Code variant is installed
	variants := GetVSCodeVariants()
	return len(variants) > 0
}

// SetVariant updates the VS Code variant and config path
func (v *VSCodeIntegration) SetVariant(variant VSCodeVariant) {
	v.variant = variant
	v.configPath = filepath.Join(variant.ConfigDir, "User", "settings.json")
}

func (v *VSCodeIntegration) Apply(t theme.Theme) error {
	// Check if theme has official VS Code extension
	themeExt, hasExtension := vscodeThemeExtensions[t.Name]

	// Install extensions if available
	if hasExtension {
		if err := v.installExtensions(themeExt); err != nil {
			// Don't fail if extension installation fails, just log and continue
			fmt.Printf("Warning: Failed to install extensions: %v\n", err)
		}
	}

	// Read existing settings
	var settings map[string]interface{}

	data, err := os.ReadFile(v.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			settings = make(map[string]interface{})
		} else {
			return fmt.Errorf("failed to read settings: %w", err)
		}
	} else {
		// Strip comments from JSON (VS Code allows comments in settings.json)
		cleanedData := stripJSONComments(string(data))
		if err := json.Unmarshal([]byte(cleanedData), &settings); err != nil {
			return fmt.Errorf("failed to parse settings: %w", err)
		}
	}

	// Create backup
	if len(data) > 0 {
		backupPath := v.configPath + ".backup"
		if err := os.WriteFile(backupPath, data, 0644); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// If official extension exists, set theme preferences
	if hasExtension {
		settings["workbench.colorTheme"] = themeExt.ThemeName
		settings["workbench.iconTheme"] = themeExt.IconTheme
		settings["workbench.productIconTheme"] = themeExt.ProductIconTheme
		// Don't clear color customizations - user may have custom overrides
		// Only remove if it's empty or doesn't exist
		if existingCustomizations, ok := settings["workbench.colorCustomizations"].(map[string]interface{}); ok {
			if len(existingCustomizations) == 0 {
				delete(settings, "workbench.colorCustomizations")
			}
			// Otherwise keep user's custom color overrides
		}
	} else {
		// Fallback to custom color customizations
		themeColors := v.generateVSCodeColors(t)
		terminalColors := v.generateTerminalColors(t)

		// Convert to map[string]interface{} for merging
		colors := make(map[string]interface{})

		// Add theme colors
		for k, v := range themeColors {
			colors[k] = v
		}

		// Add terminal colors
		for k, v := range terminalColors {
			colors[k] = v
		}

		// Preserve existing customizations and merge with theme colors
		// User customizations override theme colors
		if existingCustomizations, ok := settings["workbench.colorCustomizations"].(map[string]interface{}); ok {
			for k, v := range existingCustomizations {
				colors[k] = v
			}
		}

		settings["workbench.colorCustomizations"] = colors
	}

	// Write updated settings
	newData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(v.configPath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}

func (v *VSCodeIntegration) installExtensions(themeExt VSCodeThemeExtension) error {
	// Collect all extensions to install
	extensionsToInstall := []string{themeExt.ExtensionID}

	// Add icon theme extension if needed
	if iconExt, exists := iconThemeExtensions[themeExt.IconTheme]; exists {
		// Only add if it's not the same as the main theme extension
		if iconExt != themeExt.ExtensionID {
			extensionsToInstall = append(extensionsToInstall, iconExt)
		}
	}

	// Add product icon theme extension if needed
	if productIconExt, exists := iconThemeExtensions[themeExt.ProductIconTheme]; exists {
		// Only add if it's not already in the list
		alreadyAdded := false
		for _, ext := range extensionsToInstall {
			if ext == productIconExt {
				alreadyAdded = true
				break
			}
		}
		if !alreadyAdded {
			extensionsToInstall = append(extensionsToInstall, productIconExt)
		}
	}

	// Install each extension
	for _, extID := range extensionsToInstall {
		if err := v.installExtension(extID); err != nil {
			return fmt.Errorf("failed to install %s: %w", extID, err)
		}
	}

	return nil
}

func (v *VSCodeIntegration) installExtension(extensionID string) error {
	// Check if extension is already installed
	if v.isExtensionInstalled(extensionID) {
		fmt.Printf("Extension %s is already installed\n", extensionID)
		return nil
	}

	fmt.Printf("Installing VS Code extension: %s...\n", extensionID)

	// Try to find VS Code CLI
	codeCmd := v.findVSCodeCLI()
	if codeCmd == "" {
		return fmt.Errorf("VS Code CLI not found. Please install extensions manually")
	}

	// Install the extension
	cmd := exec.Command(codeCmd, "--install-extension", extensionID, "--force")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install extension: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Successfully installed %s\n", extensionID)
	return nil
}

func (v *VSCodeIntegration) isExtensionInstalled(extensionID string) bool {
	codeCmd := v.findVSCodeCLI()
	if codeCmd == "" {
		return false
	}

	cmd := exec.Command(codeCmd, "--list-extensions")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	installedExtensions := strings.Split(string(output), "\n")
	for _, installed := range installedExtensions {
		if strings.EqualFold(strings.TrimSpace(installed), extensionID) {
			return true
		}
	}

	return false
}

func (v *VSCodeIntegration) findVSCodeCLI() string {
	// First try the variant-specific CLI command
	if path, err := exec.LookPath(v.variant.CLICommand); err == nil {
		return path
	}

	// Try variant-specific app bundle path
	appBundleCLI := filepath.Join(v.variant.AppPath, "Contents", "Resources", "app", "bin", v.variant.CLICommand)
	if _, err := os.Stat(appBundleCLI); err == nil {
		return appBundleCLI
	}

	// Fallback: Common VS Code CLI locations on macOS
	possiblePaths := []string{
		"/usr/local/bin/" + v.variant.CLICommand,
		"/opt/homebrew/bin/" + v.variant.CLICommand,
		"/usr/local/bin/code",
		"/opt/homebrew/bin/code",
		"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code",
		"/Applications/Visual Studio Code - Insiders.app/Contents/Resources/app/bin/code-insiders",
		"/Applications/Cursor.app/Contents/Resources/app/bin/cursor",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Try to find in PATH
	if path, err := exec.LookPath("code"); err == nil {
		return path
	}

	return ""
}

func (v *VSCodeIntegration) generateVSCodeColors(t theme.Theme) map[string]string {
	return map[string]string{
		// Editor
		"editor.background":              t.Colors.Background,
		"editor.foreground":              t.Colors.Foreground,
		"editorCursor.foreground":        t.Colors.Blue,
		"editor.lineHighlightBackground": t.Colors.Black,
		"editor.selectionBackground":     t.Colors.BrightBlack,

		// Sidebar
		"sideBar.background":              t.Colors.Black,
		"sideBar.foreground":              t.Colors.Foreground,
		"sideBarSectionHeader.background": t.Colors.Background,

		// Activity Bar
		"activityBar.background":         t.Colors.Black,
		"activityBar.foreground":         t.Colors.Blue,
		"activityBar.inactiveForeground": t.Colors.BrightBlack,

		// Status Bar
		"statusBar.background":         t.Colors.Black,
		"statusBar.foreground":         t.Colors.Foreground,
		"statusBar.noFolderBackground": t.Colors.Background,

		// Title Bar
		"titleBar.activeBackground":   t.Colors.Black,
		"titleBar.activeForeground":   t.Colors.Foreground,
		"titleBar.inactiveBackground": t.Colors.Background,

		// Tabs
		"tab.activeBackground":   t.Colors.Background,
		"tab.inactiveBackground": t.Colors.Black,
		"tab.activeForeground":   t.Colors.Foreground,
		"tab.border":             t.Colors.Black,

		// Panel
		"panel.background":              t.Colors.Background,
		"panel.border":                  t.Colors.Black,
		"panelTitle.activeForeground":   t.Colors.Foreground,
		"panelTitle.inactiveForeground": t.Colors.BrightBlack,
	}
}

func (v *VSCodeIntegration) generateTerminalColors(t theme.Theme) map[string]string {
	return map[string]string{
		"terminal.background":        t.Colors.Background,
		"terminal.foreground":        t.Colors.Foreground,
		"terminal.ansiBlack":         t.Colors.Black,
		"terminal.ansiRed":           t.Colors.Red,
		"terminal.ansiGreen":         t.Colors.Green,
		"terminal.ansiYellow":        t.Colors.Yellow,
		"terminal.ansiBlue":          t.Colors.Blue,
		"terminal.ansiMagenta":       t.Colors.Magenta,
		"terminal.ansiCyan":          t.Colors.Cyan,
		"terminal.ansiWhite":         t.Colors.White,
		"terminal.ansiBrightBlack":   t.Colors.BrightBlack,
		"terminal.ansiBrightRed":     t.Colors.BrightRed,
		"terminal.ansiBrightGreen":   t.Colors.BrightGreen,
		"terminal.ansiBrightYellow":  t.Colors.BrightYellow,
		"terminal.ansiBrightBlue":    t.Colors.BrightBlue,
		"terminal.ansiBrightMagenta": t.Colors.BrightMagenta,
		"terminal.ansiBrightCyan":    t.Colors.BrightCyan,
		"terminal.ansiBrightWhite":   t.Colors.BrightWhite,
	}
}
