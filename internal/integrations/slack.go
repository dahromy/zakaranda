package integrations

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"zakaranda/internal/theme"

	"github.com/atotto/clipboard"
)

type SlackIntegration struct{}

func NewSlackIntegration() *SlackIntegration {
	return &SlackIntegration{}
}

func (s *SlackIntegration) Name() string {
	return "Slack"
}

func (s *SlackIntegration) ConfigPath() string {
	return "Manual (copy to clipboard)"
}

func (s *SlackIntegration) IsInstalled() bool {
	// Check if Slack is installed by looking for the application
	switch runtime.GOOS {
	case "darwin":
		// macOS: Check if Slack.app exists
		_, err := exec.Command("mdfind", "kMDItemKind == 'Application' && kMDItemFSName == 'Slack.app'").Output()
		return err == nil
	case "linux":
		// Linux: Check if slack command exists
		_, err := exec.LookPath("slack")
		return err == nil
	case "windows":
		// Windows: Check if slack.exe exists in common locations
		_, err := exec.LookPath("slack.exe")
		return err == nil
	default:
		// For other platforms, assume it might be installed
		return true
	}
}

func (s *SlackIntegration) Apply(t theme.Theme) error {
	// Generate Slack theme string (4 colors)
	themeString := s.generateSlackTheme(t.Colors)

	// Copy to clipboard
	if err := clipboard.WriteAll(themeString); err != nil {
		return fmt.Errorf("failed to copy theme to clipboard: %w\n\nTheme colors: %s", err, themeString)
	}

	// Display success message with instructions
	fmt.Println("\n✓ Slack theme copied to clipboard!")
	fmt.Println("\n📋 Theme colors:", themeString)
	fmt.Println("\n🎨 To apply in Slack:")
	fmt.Println("1. Open Slack")
	fmt.Println("2. Go to Preferences → Appearance → Custom theme")
	fmt.Println("3. Paste the theme colors (Cmd+V / Ctrl+V)")
	fmt.Println("4. The theme will be applied automatically")

	// Try to open Slack preferences (optional, may not work on all systems)
	s.tryOpenSlackPreferences()

	return nil
}

// generateSlackTheme creates a Slack theme string from the color palette
// Slack uses 4 colors in this order:
// 1. System navigation (sidebar/navigation background)
// 2. Selected items (background for selected channels/items)
// 3. Presence indication (online/active status indicators)
// 4. Notifications (notification badges and alerts)
func (s *SlackIntegration) generateSlackTheme(colors theme.ColorPalette) string {
	// Map the 18-color terminal palette to Slack's 4 UI colors
	slackColors := []string{
		colors.Background, // 1. System navigation - use background color
		colors.Blue,       // 2. Selected items - use blue for active selection
		colors.Green,      // 3. Presence indication - use green for online status
		colors.Red,        // 4. Notifications - use red for alerts/mentions
	}

	// Join with commas (Slack expects comma-separated hex codes with # prefix)
	return strings.Join(slackColors, ",")
}

// tryOpenSlackPreferences attempts to open Slack's preferences
// This is a best-effort attempt and may not work on all systems
func (s *SlackIntegration) tryOpenSlackPreferences() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// macOS: Try to open Slack with deep link
		cmd = exec.Command("open", "slack://preferences")
	case "linux":
		// Linux: Try to open Slack with xdg-open
		cmd = exec.Command("xdg-open", "slack://preferences")
	case "windows":
		// Windows: Try to open Slack with start
		cmd = exec.Command("cmd", "/c", "start", "slack://preferences")
	default:
		// Unsupported platform
		return
	}

	// Run the command, but don't fail if it doesn't work
	// This is just a convenience feature
	_ = cmd.Run()
}
