# 🎨 Zakaranda

A unified theme manager for macOS that applies consistent color schemes across multiple applications with a beautiful terminal UI.

![Zakaranda](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/Platform-macOS-000000?style=flat&logo=apple)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

## ✨ Features

- 🎭 **3 Beautiful Themes**: Nord, Catppuccin (4 variants), and Rose Pine (3 variants)
- 🖥️ **7 Application Integrations**: VS Code, Alacritty, Warp, iTerm2, Starship, Zed, and macOS Wallpaper
- 👁️ **Live Preview**: Preview themes with color palettes before applying
- 🎨 **Custom Themes**: Load your own themes from JSON/YAML/TOML files
- 💾 **Automatic Backups**: Creates `.backup` files before modifying configs
- ✨ **Interactive TUI**: Beautiful terminal interface built with Bubble Tea
- ⚡ **Powerline Support**: Perfect powerline rendering in Starship prompts

## 📸 Screenshots

### Theme Selection
```
🎨 Zakaranda


Select a theme:
❯ Nord
  Catppuccin
  Rose Pine
  Load Custom Theme
```

### Application Selection
```
Select applications to theme:
✓ VS Code
✓ Alacritty
✓ Warp
✓ iTerm2
✓ Starship
✓ Zed
✓ Wallpaper
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** - [Install Go](https://go.dev/doc/install)
- **macOS** - Currently supports macOS (can be adapted for Linux)
- One or more of the supported applications installed

### Installation

#### Option 1: Download Pre-built Binary (Recommended)

**macOS Intel (x86_64):**
```bash
curl -L https://github.com/YOUR_USERNAME/zakaranda/releases/latest/download/zakaranda-darwin-amd64.tar.gz | tar xz
cd zakaranda-darwin-amd64
sudo mv zakaranda /usr/local/bin/
sudo mkdir -p /usr/local/share/zakaranda
sudo cp -r assets /usr/local/share/zakaranda/
```

**macOS Apple Silicon (ARM64):**
```bash
curl -L https://github.com/YOUR_USERNAME/zakaranda/releases/latest/download/zakaranda-darwin-arm64.tar.gz | tar xz
cd zakaranda-darwin-arm64
sudo mv zakaranda /usr/local/bin/
sudo mkdir -p /usr/local/share/zakaranda
sudo cp -r assets /usr/local/share/zakaranda/
```

#### Option 2: Build from Source

1. **Clone the repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/zakaranda.git
   cd zakaranda
   ```

2. **Build the application**
   ```bash
   go build -o zakaranda ./cmd/zakaranda
   ```

3. **Run the theme manager**
   ```bash
   ./zakaranda
   ```

4. **(Optional) Install globally**
   ```bash
   sudo mv zakaranda /usr/local/bin/
   sudo mkdir -p /usr/local/share/zakaranda
   sudo cp -r assets /usr/local/share/zakaranda/
   ```

## 📖 Usage

### Basic Usage

1. **Launch the theme manager**
   ```bash
   ./zakaranda
   ```

2. **Select a theme**
   - Use arrow keys to navigate
   - Press Enter to select
   - Preview themes before applying

3. **Choose applications**
   - Use Space to toggle applications
   - Press Enter to confirm selection

4. **Apply the theme**
   - The theme will be applied to all selected applications
   - Backup files are created automatically

### Theme Variants

#### Catppuccin
- Mocha (dark, warm)
- Frappé (dark, cool)
- Macchiato (dark, balanced)
- Latte (light)

#### Rose Pine
- Default (dark)
- Moon (darker)
- Dawn (light)

### Custom Themes

Create a custom theme file in JSON, YAML, or TOML format:

```json
{
  "name": "My Custom Theme",
  "colors": {
    "background": "#1e1e2e",
    "foreground": "#cdd6f4",
    "black": "#45475a",
    "red": "#f38ba8",
    "green": "#a6e3a1",
    "yellow": "#f9e2af",
    "blue": "#89b4fa",
    "magenta": "#f5c2e7",
    "cyan": "#94e2d5",
    "white": "#bac2de",
    "brightBlack": "#585b70",
    "brightRed": "#f38ba8",
    "brightGreen": "#a6e3a1",
    "brightYellow": "#f9e2af",
    "brightBlue": "#89b4fa",
    "brightMagenta": "#f5c2e7",
    "brightCyan": "#94e2d5",
    "brightWhite": "#a6adc8"
  }
}
```

Load it by selecting "Load Custom Theme" from the menu.

## 🎯 Supported Applications

### VS Code
- **Config**: `~/Library/Application Support/Code/User/settings.json`
- **Features**: Workbench colors, terminal colors, editor theme
- **Extensions**: Automatically installs theme extensions if needed

### Alacritty
- **Config**: `~/.config/alacritty/alacritty.toml` or `alacritty.yml`
- **Features**: Full terminal color scheme, cursor, selection colors

### Warp
- **Config**: `~/.warp/themes/`
- **Features**: Custom theme files with accent colors

### iTerm2
- **Config**: `~/Library/Application Support/iTerm2/DynamicProfiles/`
- **Features**: Dynamic profiles with full color schemes

### Starship
- **Config**: `~/.config/starship.toml`
- **Features**: Powerline prompts with theme-matched colors
- **Special**: Perfect powerline character rendering (U+E0B0, U+E0B4, U+E0B6)

### Zed
- **Config**: `~/.config/zed/settings.json`
- **Features**: Theme selection from installed extensions
- **Note**: Requires Catppuccin/Nord/Rose Pine extensions installed

### Wallpaper
- **Features**: Sets macOS desktop wallpaper to match theme
- **Wallpapers**: Stored in `~/.config/zakaranda/wallpapers/`

## 🛠️ Development

### Project Structure

```
zakaranda/
├── README.md
├── go.mod
├── go.sum
├── zakaranda           # Compiled binary (gitignored)
├── assets/             # Static assets
│   └── wallpapers/     # Theme wallpapers
│       ├── nord.jpg                    # Nord theme wallpaper
│       └── catppuccin-rosepine.jpg     # Catppuccin/Rose Pine wallpaper
├── cmd/
│   └── zakaranda/
│       └── main.go         # Entry point
└── internal/               # Private application code
    ├── integrations/       # Application integrations
    │   ├── integration.go  # Interface definition
    │   ├── factory.go      # Integration factory
    │   ├── vscode.go       # VS Code integration
    │   ├── alacritty.go    # Alacritty integration
    │   ├── warp.go         # Warp integration
    │   ├── iterm2.go       # iTerm2 integration
    │   ├── starship.go     # Starship integration
    │   ├── zed.go          # Zed integration
    │   └── wallpaper.go    # Wallpaper integration
    ├── theme/              # Theme logic
    │   ├── theme.go        # Theme types
    │   ├── loader.go       # Theme loading
    │   ├── preview.go      # Theme preview
    │   └── utils.go        # Utilities
    ├── config/             # Configuration
    │   └── config.go       # Config management
    └── ui/                 # Terminal UI
        └── tui.go          # Bubble Tea TUI
```

### Building from Source

```bash
# Build
go build -o zakaranda ./cmd/zakaranda

# Run tests
go test ./...

# Install globally (optional)
sudo mv zakaranda /usr/local/bin/
```

### Adding a New Theme

1. Add theme colors to the respective integration files
2. Update the theme selection menu in `main.go`
3. Test with all integrations

### Adding a New Integration

1. Create a new file (e.g., `myapp.go`)
2. Implement the `AppIntegration` interface:
   ```go
   type AppIntegration interface {
       Name() string
       ConfigPath() string
       IsInstalled() bool
       Apply(theme Theme) error
   }
   ```
3. Add to the integrations list in `main.go`

## 🐛 Troubleshooting

### Theme not applying to VS Code
- Ensure VS Code is closed when applying themes
- Check that the config path exists: `~/Library/Application Support/Code/User/`
- Verify the settings.json is valid JSON

### Starship powerline characters broken
- Install a Nerd Font (e.g., JetBrains Mono Nerd Font)
- Set your terminal to use the Nerd Font
- Restart your terminal

### Zed theme not changing
- Install the required theme extension first:
  - Catppuccin: `zed://extensions/catppuccin`
  - Nord: `zed://extensions/nord`
  - Rose Pine: `zed://extensions/rose-pine-theme`

### Backup files accumulating
- Backups are created as `.backup` files
- Safe to delete old backups manually
- Only the most recent backup is needed for rollback

## 📝 License

MIT License - feel free to use and modify!

## 🙏 Acknowledgments

- [Nord Theme](https://www.nordtheme.com/)
- [Catppuccin](https://github.com/catppuccin/catppuccin)
- [Rosé Pine](https://rosepinetheme.com/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Starship](https://starship.rs/) - Cross-shell prompt

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

Made with ❤️ for terminal enthusiasts

