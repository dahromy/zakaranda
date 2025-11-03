package integrations

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"zakaranda/internal/theme"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type AlacrittyIntegration struct {
	configPath string
	themesPath string
}

// Map theme names to official Alacritty theme repository file names
var alacrittyThemeMap = map[string]string{
	"Nord":                 "nord.toml",
	"Catppuccin Latte":     "catppuccin_latte.toml",
	"Catppuccin Frappe":    "catppuccin_frappe.toml",
	"Catppuccin Macchiato": "catppuccin_macchiato.toml",
	"Catppuccin Mocha":     "catppuccin_mocha.toml",
	"Rose Pine":            "rose_pine.toml",
	"Rose Pine Moon":       "rose_pine_moon.toml",
	"Rose Pine Dawn":       "rose_pine_dawn.toml",
}

func NewAlacrittyIntegration() *AlacrittyIntegration {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to empty string, will be caught by IsInstalled
		return &AlacrittyIntegration{configPath: "", themesPath: ""}
	}

	configPath := filepath.Join(home, ".config", "alacritty", "alacritty.toml")
	// Try .yml as fallback
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(home, ".config", "alacritty", "alacritty.yml")
	}

	themesPath := filepath.Join(home, ".config", "alacritty", "themes")

	return &AlacrittyIntegration{
		configPath: configPath,
		themesPath: themesPath,
	}
}

func (a *AlacrittyIntegration) Name() string {
	return "Alacritty"
}

func (a *AlacrittyIntegration) ConfigPath() string {
	return a.configPath
}

func (a *AlacrittyIntegration) IsInstalled() bool {
	// Check if config directory exists
	_, err := os.Stat(filepath.Dir(a.configPath))
	return err == nil
}

func (a *AlacrittyIntegration) Apply(t theme.Theme) error {
	// Ensure theme repository is cloned
	if err := a.ensureThemeRepo(); err != nil {
		return fmt.Errorf("failed to setup theme repository: %w", err)
	}

	// Check if official theme exists
	officialTheme, hasOfficial := alacrittyThemeMap[t.Name]
	officialThemePath := filepath.Join(a.themesPath, "alacritty", "themes", officialTheme)

	// Read existing config
	var config map[string]any
	data, err := os.ReadFile(a.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			config = make(map[string]any)
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	} else {
		// Determine file format based on extension
		isTOML := strings.HasSuffix(strings.ToLower(a.configPath), ".toml")

		if isTOML {
			if err := toml.Unmarshal(data, &config); err != nil {
				return fmt.Errorf("failed to parse TOML config: %w", err)
			}
		} else {
			if err := yaml.Unmarshal(data, &config); err != nil {
				return fmt.Errorf("failed to parse YAML config: %w", err)
			}
		}
	}

	// Create backup
	if len(data) > 0 {
		backupPath := a.configPath + ".backup"
		if err := os.WriteFile(backupPath, data, 0644); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Determine file format
	isTOML := strings.HasSuffix(strings.ToLower(a.configPath), ".toml")

	// Use import method if official theme exists and config is TOML
	if hasOfficial && isTOML && a.themeFileExists(officialThemePath) {
		// Use import directive
		general, ok := config["general"].(map[string]any)
		if !ok {
			general = make(map[string]any)
			config["general"] = general
		}

		// Set import path
		importPath := fmt.Sprintf("~/.config/alacritty/themes/alacritty/themes/%s", officialTheme)
		general["import"] = []string{importPath}

		// Remove colors section if it exists (let import handle it)
		delete(config, "colors")
	} else {
		// Fallback to manual color palette
		colors := a.generateAlacrittyColors(t)
		config["colors"] = colors

		// Remove import if it exists
		if general, ok := config["general"].(map[string]any); ok {
			delete(general, "import")
		}
	}

	// Write updated config in the same format
	var newData []byte

	if isTOML {
		// For TOML, we need to encode properly
		buf := new(strings.Builder)
		encoder := toml.NewEncoder(buf)
		if err := encoder.Encode(config); err != nil {
			return fmt.Errorf("failed to marshal TOML config: %w", err)
		}
		newData = []byte(buf.String())
	} else {
		// For YAML, use manual colors (import not supported in YAML)
		colors := a.generateAlacrittyColors(t)
		config["colors"] = colors

		newData, err = yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML config: %w", err)
		}
	}

	if err := os.WriteFile(a.configPath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (a *AlacrittyIntegration) ensureThemeRepo() error {
	// Check if theme repository already exists
	repoPath := filepath.Join(a.themesPath, "alacritty")
	if _, err := os.Stat(repoPath); err == nil {
		// Repository exists, skip update to avoid slow git operations
		// User can manually update if needed
		return nil
	}

	// Repository doesn't exist, clone it
	return a.cloneThemeRepo()
}

func (a *AlacrittyIntegration) cloneThemeRepo() error {
	// Create themes directory
	if err := os.MkdirAll(a.themesPath, 0755); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	// Clone the repository with minimal depth for faster cloning
	cmd := exec.Command("git", "clone",
		"--depth", "1",
		"--single-branch",
		"https://github.com/alacritty/alacritty-theme",
		filepath.Join(a.themesPath, "alacritty"))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone theme repository: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func (a *AlacrittyIntegration) themeFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (a *AlacrittyIntegration) generateAlacrittyColors(t theme.Theme) map[string]any {
	return map[string]any{
		"primary": map[string]string{
			"background": t.Colors.Background,
			"foreground": t.Colors.Foreground,
		},
		"cursor": map[string]string{
			"text":   t.Colors.Background,
			"cursor": t.Colors.Foreground,
		},
		"normal": map[string]string{
			"black":   t.Colors.Black,
			"red":     t.Colors.Red,
			"green":   t.Colors.Green,
			"yellow":  t.Colors.Yellow,
			"blue":    t.Colors.Blue,
			"magenta": t.Colors.Magenta,
			"cyan":    t.Colors.Cyan,
			"white":   t.Colors.White,
		},
		"bright": map[string]string{
			"black":   t.Colors.BrightBlack,
			"red":     t.Colors.BrightRed,
			"green":   t.Colors.BrightGreen,
			"yellow":  t.Colors.BrightYellow,
			"blue":    t.Colors.BrightBlue,
			"magenta": t.Colors.BrightMagenta,
			"cyan":    t.Colors.BrightCyan,
			"white":   t.Colors.BrightWhite,
		},
		"selection": map[string]string{
			"text":       t.Colors.Foreground,
			"background": t.Colors.BrightBlack,
		},
	}
}
