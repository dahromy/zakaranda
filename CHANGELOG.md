# Changelog

All notable changes to Zakaranda will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-10-28

### Added
- Slack integration with automatic clipboard copy
  - Generates 4-color theme string for Slack's custom theme
  - Automatically copies theme to clipboard
  - Provides clear instructions for manual application in Slack
  - Maps terminal colors to Slack UI colors (System navigation, Selected items, Presence indication, Notifications)
- New dependency: `github.com/atotto/clipboard` for cross-platform clipboard support

## [1.0.0] - 2025-10-28

### Added
- Initial release of Zakaranda theme manager
- Support for 3 theme families: Nord, Catppuccin (4 variants), Rose Pine (3 variants)
- 7 application integrations:
  - VS Code (with variant selection)
  - Alacritty
  - Warp
  - iTerm2
  - Starship (with powerline support)
  - Zed
  - macOS Wallpaper
- Interactive TUI built with Bubble Tea
- Live theme preview with color palettes
- Automatic backup creation before modifications
- Custom theme loading from JSON/YAML/TOML files
- High-quality wallpapers for Nord and Catppuccin/Rose Pine themes
- GitHub Actions CI/CD workflows for automated builds and releases
- Support for both macOS Intel (x86_64) and Apple Silicon (ARM64)

