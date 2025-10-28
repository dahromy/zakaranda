package integrations

import "zakaranda/internal/theme"

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ZedIntegration struct {
	configPath     string
	extensionsPath string
}

// ZedExtension represents a Zed theme extension
type ZedExtension struct {
	ExtensionID string
	ThemeName   string // Theme name to set in settings.json
	IsLight     bool   // Whether this is a light theme
}

// Map of our theme names to Zed extension information
var zedThemeExtensions = map[string]ZedExtension{
	"Nord": {
		ExtensionID: "nord",
		ThemeName:   "Nord Dark",
		IsLight:     false,
	},
	"Catppuccin Latte": {
		ExtensionID: "catppuccin",
		ThemeName:   "Catppuccin Latte",
		IsLight:     true,
	},
	"Catppuccin Frappe": {
		ExtensionID: "catppuccin",
		ThemeName:   "Catppuccin Frappé",
		IsLight:     false,
	},
	"Catppuccin Macchiato": {
		ExtensionID: "catppuccin",
		ThemeName:   "Catppuccin Macchiato",
		IsLight:     false,
	},
	"Catppuccin Mocha": {
		ExtensionID: "catppuccin",
		ThemeName:   "Catppuccin Mocha",
		IsLight:     false,
	},
	"Rose Pine": {
		ExtensionID: "rose-pine-theme",
		ThemeName:   "Rosé Pine",
		IsLight:     false,
	},
	"Rose Pine Moon": {
		ExtensionID: "rose-pine-theme",
		ThemeName:   "Rosé Pine Moon",
		IsLight:     false,
	},
	"Rose Pine Dawn": {
		ExtensionID: "rose-pine-theme",
		ThemeName:   "Rosé Pine Dawn",
		IsLight:     true,
	},
}

func NewZedIntegration() *ZedIntegration {
	home, err := os.UserHomeDir()
	if err != nil {
		return &ZedIntegration{configPath: "", extensionsPath: ""}
	}

	configPath := filepath.Join(home, ".config", "zed", "settings.json")

	// Extensions are installed in different locations based on OS
	// macOS: ~/Library/Application Support/Zed/extensions/installed
	// Linux: ~/.local/share/zed/extensions/installed (or $XDG_DATA_HOME)
	var extensionsPath string
	if os.Getenv("XDG_DATA_HOME") != "" {
		extensionsPath = filepath.Join(os.Getenv("XDG_DATA_HOME"), "zed", "extensions", "installed")
	} else {
		// Default paths for macOS and Linux
		extensionsPath = filepath.Join(home, "Library", "Application Support", "Zed", "extensions", "installed")
		// Check if macOS path exists, otherwise use Linux path
		if _, err := os.Stat(extensionsPath); os.IsNotExist(err) {
			extensionsPath = filepath.Join(home, ".local", "share", "zed", "extensions", "installed")
		}
	}

	return &ZedIntegration{
		configPath:     configPath,
		extensionsPath: extensionsPath,
	}
}

func (z *ZedIntegration) Name() string {
	return "Zed"
}

func (z *ZedIntegration) ConfigPath() string {
	return z.configPath
}

func (z *ZedIntegration) IsInstalled() bool {
	configDir := filepath.Dir(z.configPath)
	_, err := os.Stat(configDir)
	return err == nil
}

// IsExtensionInstalled checks if a Zed extension is installed
func (z *ZedIntegration) IsExtensionInstalled(extensionID string) bool {
	extensionPath := filepath.Join(z.extensionsPath, extensionID)
	_, err := os.Stat(extensionPath)
	return err == nil
}

// GetExtensionURL returns the zed:// URL for installing an extension
func (z *ZedIntegration) GetExtensionURL(extensionID string) string {
	return fmt.Sprintf("zed://extensions/%s", extensionID)
}

func (z *ZedIntegration) Apply(t theme.Theme) error {
	// Check if theme has official Zed extension
	themeExt, hasExtension := zedThemeExtensions[t.Name]
	if !hasExtension {
		return fmt.Errorf("no Zed extension available for theme: %s", t.Name)
	}

	// Check if extension is installed
	if !z.IsExtensionInstalled(themeExt.ExtensionID) {
		extensionURL := z.GetExtensionURL(themeExt.ExtensionID)
		return fmt.Errorf("extension not installed\n\nPlease install the %s extension first:\n%s\n\nAfter installation, press Enter to continue", themeExt.ExtensionID, extensionURL)
	}

	// Read existing settings
	var settings map[string]interface{}
	data, err := os.ReadFile(z.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			settings = make(map[string]interface{})
		} else {
			return fmt.Errorf("failed to read settings: %w", err)
		}
	} else {
		// Strip comments from JSONC
		cleanedData := z.stripJSONComments(string(data))
		if err := json.Unmarshal([]byte(cleanedData), &settings); err != nil {
			return fmt.Errorf("failed to parse settings: %w", err)
		}
	}

	// Create backup
	if len(data) > 0 {
		backupPath := z.configPath + ".backup"
		if err := os.WriteFile(backupPath, data, 0644); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Update theme configuration
	themeConfig := map[string]interface{}{
		"mode": "system",
	}

	// Set both light and dark to the same theme
	// User can manually adjust if they want different themes for light/dark mode
	themeConfig["light"] = themeExt.ThemeName
	themeConfig["dark"] = themeExt.ThemeName

	settings["theme"] = themeConfig

	// Write updated settings
	newData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(z.configPath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}

// stripJSONComments removes comments from JSONC (JSON with Comments)
// This is a character-by-character parser that preserves strings
func (z *ZedIntegration) stripJSONComments(jsonc string) string {
	var result strings.Builder
	inString := false
	escaped := false
	inSingleLineComment := false
	inMultiLineComment := false

	for i := 0; i < len(jsonc); i++ {
		ch := jsonc[i]

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

		// Handle multi-line comments
		if inMultiLineComment {
			if ch == '*' && i+1 < len(jsonc) && jsonc[i+1] == '/' {
				inMultiLineComment = false
				i++ // Skip the '/'
			}
			continue
		}

		// Handle single-line comments
		if inSingleLineComment {
			if ch == '\n' {
				inSingleLineComment = false
				result.WriteByte(ch) // Preserve newline
			}
			continue
		}

		// Check for comment start
		if ch == '/' && i+1 < len(jsonc) {
			next := jsonc[i+1]
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

		// Check for string start
		if ch == '"' {
			inString = true
		}

		result.WriteByte(ch)
	}

	return result.String()
}
