package integrations

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"zakaranda/internal/theme"
)

type WallpaperIntegration struct {
	wallpaperPath string
	assetsPath    string
}

func NewWallpaperIntegration() *WallpaperIntegration {
	home, err := os.UserHomeDir()
	if err != nil {
		return &WallpaperIntegration{wallpaperPath: "", assetsPath: ""}
	}
	wallpaperPath := filepath.Join(home, ".config", "zakaranda", "wallpapers")

	// Get the executable path to find assets
	execPath, err := os.Executable()
	if err != nil {
		return &WallpaperIntegration{wallpaperPath: wallpaperPath, assetsPath: ""}
	}
	execDir := filepath.Dir(execPath)
	assetsPath := filepath.Join(execDir, "assets", "wallpapers")

	return &WallpaperIntegration{
		wallpaperPath: wallpaperPath,
		assetsPath:    assetsPath,
	}
}

func (w *WallpaperIntegration) Name() string {
	return "macOS Wallpaper"
}

func (w *WallpaperIntegration) ConfigPath() string {
	return w.wallpaperPath
}

func (w *WallpaperIntegration) IsInstalled() bool {
	// Always available on macOS
	return true
}

func (w *WallpaperIntegration) Apply(t theme.Theme) error {
	// Create wallpapers directory if it doesn't exist
	if err := os.MkdirAll(w.wallpaperPath, 0755); err != nil {
		return fmt.Errorf("failed to create wallpapers directory: %w", err)
	}

	// Map theme name to wallpaper file
	sourceWallpaper := w.getWallpaperForTheme(t.Name)
	if sourceWallpaper == "" {
		return fmt.Errorf("no wallpaper found for theme: %s", t.Name)
	}

	// Copy wallpaper to user's config directory
	destWallpaper := filepath.Join(w.wallpaperPath, filepath.Base(sourceWallpaper))
	if err := w.copyFile(sourceWallpaper, destWallpaper); err != nil {
		return fmt.Errorf("failed to copy wallpaper: %w", err)
	}

	// Set as desktop wallpaper using AppleScript
	if err := w.setWallpaper(destWallpaper); err != nil {
		return fmt.Errorf("failed to set wallpaper: %w", err)
	}

	return nil
}

// getWallpaperForTheme maps theme names to wallpaper files
func (w *WallpaperIntegration) getWallpaperForTheme(themeName string) string {
	// Normalize theme name for comparison
	normalizedName := strings.ToLower(themeName)

	// Map theme names to wallpaper files
	var wallpaperFile string
	if strings.Contains(normalizedName, "nord") {
		wallpaperFile = "nord.png"
	} else if strings.Contains(normalizedName, "catppuccin") || strings.Contains(normalizedName, "rose pine") {
		wallpaperFile = "catppuccin-rosepine.jpg"
	} else {
		// Default to catppuccin-rosepine for unknown themes
		wallpaperFile = "catppuccin-rosepine.jpg"
	}

	// Try to find the wallpaper in assets directory
	wallpaperPath := filepath.Join(w.assetsPath, wallpaperFile)
	if _, err := os.Stat(wallpaperPath); err == nil {
		return wallpaperPath
	}

	// If not found in assets, return empty string
	return ""
}

// copyFile copies a file from src to dst
func (w *WallpaperIntegration) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func (w *WallpaperIntegration) setWallpaper(imagePath string) error {
	// Use osascript to set wallpaper
	script := fmt.Sprintf(`
tell application "System Events"
	tell every desktop
		set picture to "%s"
	end tell
end tell
`, imagePath)

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set wallpaper: %w, output: %s", err, string(output))
	}

	return nil
}
