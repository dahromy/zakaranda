package integrations

import "zakaranda/internal/theme"

import (
	"fmt"
	"os"
	"path/filepath"
)

type StarshipIntegration struct {
	configPath string
}

func NewStarshipIntegration() *StarshipIntegration {
	home, err := os.UserHomeDir()
	if err != nil {
		return &StarshipIntegration{configPath: ""}
	}
	configPath := filepath.Join(home, ".config", "starship.toml")
	return &StarshipIntegration{configPath: configPath}
}

func (s *StarshipIntegration) Name() string {
	return "Starship"
}

func (s *StarshipIntegration) ConfigPath() string {
	return s.configPath
}

func (s *StarshipIntegration) IsInstalled() bool {
	// Check if starship config exists or if starship is in PATH
	if _, err := os.Stat(s.configPath); err == nil {
		return true
	}

	// Check if starship binary exists
	_, err := os.Stat("/usr/local/bin/starship")
	if err == nil {
		return true
	}

	// Check in common homebrew path
	_, err = os.Stat("/opt/homebrew/bin/starship")
	return err == nil
}

func (s *StarshipIntegration) Apply(t theme.Theme) error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(s.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create backup if config exists
	if data, err := os.ReadFile(s.configPath); err == nil && len(data) > 0 {
		backupPath := s.configPath + ".backup"
		if err := os.WriteFile(backupPath, data, 0644); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Get the full official configuration based on theme name
	var configContent string
	switch t.Name {
	case "Nord":
		configContent = s.getNordConfig()
	case "Catppuccin Latte":
		configContent = s.getCatppuccinConfig("latte")
	case "Catppuccin Frappe":
		configContent = s.getCatppuccinConfig("frappe")
	case "Catppuccin Macchiato":
		configContent = s.getCatppuccinConfig("macchiato")
	case "Catppuccin Mocha":
		configContent = s.getCatppuccinConfig("mocha")
	case "Rose Pine":
		configContent = s.getRosePineConfig("default")
	case "Rose Pine Moon":
		configContent = s.getRosePineConfig("moon")
	case "Rose Pine Dawn":
		configContent = s.getRosePineConfig("dawn")
	default:
		return fmt.Errorf("unsupported theme: %s", t.Name)
	}

	// Write the full configuration
	if err := os.WriteFile(s.configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (s *StarshipIntegration) getNordConfig() string {
	return `"$schema" = 'https://starship.rs/config-schema.json'

command_timeout = 3000

format = """
[](frost_blue)\
$os\
$username\
[](bg:frost_light fg:frost_blue)\
$directory\
[](bg:frost_cyan fg:frost_light)\
$git_branch\
$git_status\
[](fg:frost_cyan bg:frost_green)\
$c\
$rust\
$golang\
$nodejs\
$php\
$java\
$kotlin\
$haskell\
$python\
[](fg:frost_green bg:snow_2)\
$conda\
[](fg:snow_2 bg:snow_3)\
$time\
[ ](fg:snow_3)\
$line_break\
$character"""

palette = 'nord'

[os]
disabled = false
style = "bg:frost_blue fg:polar_0"

[os.symbols]
Windows = ""
Ubuntu = "󰕈"
SUSE = ""
Raspbian = "󰐿"
Mint = "󰣭"
Macos = "󰀵"
Manjaro = ""
Linux = "󰌽"
Gentoo = "󰣨"
Fedora = "󰣛"
Alpine = ""
Amazon = ""
Android = ""
Arch = "󰣇"
Artix = "󰣇"
CentOS = ""
Debian = "󰣚"
Redhat = "󱄛"
RedHatEnterprise = "󱄛"

[username]
show_always = true
style_user = "bg:frost_blue fg:polar_0"
style_root = "bg:frost_blue fg:polar_0"
format = '[ $user]($style)'

[directory]
style = "bg:frost_light fg:polar_0"
format = "[ $path ]($style)"
truncation_length = 3
truncation_symbol = "…/"

[directory.substitutions]
"Documents" = "󰈙 "
"Downloads" = " "
"Music" = "󰝚 "
"Pictures" = " "
"Developer" = "󰲋 "

[git_branch]
symbol = ""
style = "bg:frost_cyan"
format = '[[ $symbol $branch ](fg:polar_0 bg:frost_cyan)]($style)'

[git_status]
style = "bg:frost_cyan"
format = '[[($all_status$ahead_behind )](fg:polar_0 bg:frost_cyan)]($style)'

[nodejs]
symbol = ""
style = "bg:frost_green"
format = '[[ $symbol( $version) ](fg:polar_0 bg:frost_green)]($style)'

[c]
symbol = " "
style = "bg:frost_green"
format = '[[ $symbol( $version) ](fg:polar_0 bg:frost_green)]($style)'

[rust]
symbol = ""
style = "bg:frost_green"
format = '[[ $symbol( $version) ](fg:polar_0 bg:frost_green)]($style)'

[golang]
symbol = ""
style = "bg:frost_green"
format = '[[ $symbol( $version) ](fg:polar_0 bg:frost_green)]($style)'

[php]
symbol = ""
style = "bg:frost_green"
format = '[[ $symbol( $version) ](fg:polar_0 bg:frost_green)]($style)'
detect_extensions = ['php']
detect_files = ['composer.json', '.php-version']
detect_folders = ['vendor']

[java]
symbol = " "
style = "bg:frost_green"
format = '[[ $symbol( $version) ](fg:polar_0 bg:frost_green)]($style)'

[kotlin]
symbol = ""
style = "bg:frost_green"
format = '[[ $symbol( $version) ](fg:polar_0 bg:frost_green)]($style)'

[haskell]
symbol = ""
style = "bg:frost_green"
format = '[[ $symbol( $version) ](fg:polar_0 bg:frost_green)]($style)'

[python]
symbol = ""
style = "bg:frost_green"
format = '[[ $symbol( $version)(\(#$virtualenv\)) ](fg:polar_0 bg:frost_green)]($style)'

[docker_context]
symbol = ""
style = "bg:snow_2"
format = '[[ $symbol( $context) ](fg:polar_0 bg:snow_2)]($style)'

[conda]
symbol = "  "
style = "fg:polar_0 bg:snow_2"
format = '[$symbol$environment ]($style)'
ignore_base = false

[time]
disabled = false
time_format = "%R"
style = "bg:snow_3"
format = '[[  $time ](fg:polar_0 bg:snow_3)]($style)'

[line_break]
disabled = true

[character]
disabled = false
success_symbol = '[❯](bold fg:frost_green)'
error_symbol = '[❯](bold fg:frost_cyan)'
vimcmd_symbol = '[❮](bold fg:frost_green)'
vimcmd_replace_one_symbol = '[❮](bold fg:snow_2)'
vimcmd_replace_symbol = '[❮](bold fg:snow_2)'
vimcmd_visual_symbol = '[❮](bold fg:frost_blue)'

[cmd_duration]
show_milliseconds = true
format = " in $duration "
style = "bg:snow_3"
disabled = false
show_notifications = true
min_time_to_notify = 45000

[palettes.nord]
# Polar Night
polar_0 = "#2e3440"
polar_1 = "#3b4252"
polar_2 = "#434c5e"
polar_3 = "#4c566a"

# Snow Storm
snow_0 = "#d8dee9"
snow_1 = "#e5e9f0"
snow_2 = "#eceff4"

# Frost
frost_green = "#8fbcbb"
frost_cyan = "#88c0d0"
frost_light = "#81a1c1"
frost_blue = "#5e81ac"

# Extra convenience mapping
snow_3 = "#e5e9f0"
`
}

func (s *StarshipIntegration) getCatppuccinConfig(variant string) string {
	// Default to mocha if no variant specified
	if variant == "" {
		variant = "mocha"
	}

	paletteName := "catppuccin_" + variant

	return `"$schema" = 'https://starship.rs/config-schema.json'

# Global command timeout (in milliseconds)
# Increased to accommodate Lando PHP wrapper script
command_timeout = 3000

format = """
[](red)\
$os\
$username\
[](bg:peach fg:red)\
$directory\
[](bg:yellow fg:peach)\
$git_branch\
$git_status\
[](fg:yellow bg:green)\
$c\
$rust\
$golang\
$nodejs\
$php\
$java\
$kotlin\
$haskell\
$python\
[](fg:green bg:sapphire)\
$conda\
[](fg:sapphire bg:lavender)\
$time\
[ ](fg:lavender)\
$line_break\
$character"""

palette = '` + paletteName + `'

[os]
disabled = false
style = "bg:red fg:crust"

[os.symbols]
Windows = ""
Ubuntu = "󰕈"
SUSE = ""
Raspbian = "󰐿"
Mint = "󰣭"
Macos = "󰀵"
Manjaro = ""
Linux = "󰌽"
Gentoo = "󰣨"
Fedora = "󰣛"
Alpine = ""
Amazon = ""
Android = ""
Arch = "󰣇"
Artix = "󰣇"
CentOS = ""
Debian = "󰣚"
Redhat = "󱄛"
RedHatEnterprise = "󱄛"

[username]
show_always = true
style_user = "bg:red fg:crust"
style_root = "bg:red fg:crust"
format = '[ $user]($style)'

[directory]
style = "bg:peach fg:crust"
format = "[ $path ]($style)"
truncation_length = 3
truncation_symbol = "…/"

[directory.substitutions]
"Documents" = "󰈙 "
"Downloads" = " "
"Music" = "󰝚 "
"Pictures" = " "
"Developer" = "󰲋 "

[git_branch]
symbol = ""
style = "bg:yellow"
format = '[[ $symbol $branch ](fg:crust bg:yellow)]($style)'

[git_status]
style = "bg:yellow"
format = '[[($all_status$ahead_behind )](fg:crust bg:yellow)]($style)'

[nodejs]
symbol = ""
style = "bg:green"
format = '[[ $symbol( $version) ](fg:crust bg:green)]($style)'

[c]
symbol = " "
style = "bg:green"
format = '[[ $symbol( $version) ](fg:crust bg:green)]($style)'

[rust]
symbol = ""
style = "bg:green"
format = '[[ $symbol( $version) ](fg:crust bg:green)]($style)'

[golang]
symbol = ""
style = "bg:green"
format = '[[ $symbol( $version) ](fg:crust bg:green)]($style)'

[php]
symbol = ""
style = "bg:green"
format = '[[ $symbol( $version) ](fg:crust bg:green)]($style)'
# Only detect PHP in directories with PHP files to avoid unnecessary checks
detect_extensions = ['php']
detect_files = ['composer.json', '.php-version']
detect_folders = ['vendor']

[java]
symbol = " "
style = "bg:green"
format = '[[ $symbol( $version) ](fg:crust bg:green)]($style)'

[kotlin]
symbol = ""
style = "bg:green"
format = '[[ $symbol( $version) ](fg:crust bg:green)]($style)'

[haskell]
symbol = ""
style = "bg:green"
format = '[[ $symbol( $version) ](fg:crust bg:green)]($style)'

[python]
symbol = ""
style = "bg:green"
format = '[[ $symbol( $version)(\(#$virtualenv\)) ](fg:crust bg:green)]($style)'

[docker_context]
symbol = ""
style = "bg:sapphire"
format = '[[ $symbol( $context) ](fg:crust bg:sapphire)]($style)'

[conda]
symbol = "  "
style = "fg:crust bg:sapphire"
format = '[$symbol$environment ]($style)'
ignore_base = false

[time]
disabled = false
time_format = "%R"
style = "bg:lavender"
format = '[[  $time ](fg:crust bg:lavender)]($style)'

[line_break]
disabled = true

[character]
disabled = false
success_symbol = '[❯](bold fg:green)'
error_symbol = '[❯](bold fg:red)'
vimcmd_symbol = '[❮](bold fg:green)'
vimcmd_replace_one_symbol = '[❮](bold fg:lavender)'
vimcmd_replace_symbol = '[❮](bold fg:lavender)'
vimcmd_visual_symbol = '[❮](bold fg:yellow)'

[cmd_duration]
show_milliseconds = true
format = " in $duration "
style = "bg:lavender"
disabled = false
show_notifications = true
min_time_to_notify = 45000

[palettes.catppuccin_mocha]
rosewater = "#f5e0dc"
flamingo = "#f2cdcd"
pink = "#f5c2e7"
mauve = "#cba6f7"
red = "#f38ba8"
maroon = "#eba0ac"
peach = "#fab387"
yellow = "#f9e2af"
green = "#a6e3a1"
teal = "#94e2d5"
sky = "#89dceb"
sapphire = "#74c7ec"
blue = "#89b4fa"
lavender = "#b4befe"
text = "#cdd6f4"
subtext1 = "#bac2de"
subtext0 = "#a6adc8"
overlay2 = "#9399b2"
overlay1 = "#7f849c"
overlay0 = "#6c7086"
surface2 = "#585b70"
surface1 = "#45475a"
surface0 = "#313244"
base = "#1e1e2e"
mantle = "#181825"
crust = "#11111b"

[palettes.catppuccin_frappe]
rosewater = "#f2d5cf"
flamingo = "#eebebe"
pink = "#f4b8e4"
mauve = "#ca9ee6"
red = "#e78284"
maroon = "#ea999c"
peach = "#ef9f76"
yellow = "#e5c890"
green = "#a6d189"
teal = "#81c8be"
sky = "#99d1db"
sapphire = "#85c1dc"
blue = "#8caaee"
lavender = "#babbf1"
text = "#c6d0f5"
subtext1 = "#b5bfe2"
subtext0 = "#a5adce"
overlay2 = "#949cbb"
overlay1 = "#838ba7"
overlay0 = "#737994"
surface2 = "#626880"
surface1 = "#51576d"
surface0 = "#414559"
base = "#303446"
mantle = "#292c3c"
crust = "#232634"

[palettes.catppuccin_latte]
rosewater = "#dc8a78"
flamingo = "#dd7878"
pink = "#ea76cb"
mauve = "#8839ef"
red = "#d20f39"
maroon = "#e64553"
peach = "#fe640b"
yellow = "#df8e1d"
green = "#40a02b"
teal = "#179299"
sky = "#04a5e5"
sapphire = "#209fb5"
blue = "#1e66f5"
lavender = "#7287fd"
text = "#4c4f69"
subtext1 = "#5c5f77"
subtext0 = "#6c6f85"
overlay2 = "#7c7f93"
overlay1 = "#8c8fa1"
overlay0 = "#9ca0b0"
surface2 = "#acb0be"
surface1 = "#bcc0cc"
surface0 = "#ccd0da"
base = "#eff1f5"
mantle = "#e6e9ef"
crust = "#dce0e8"

[palettes.catppuccin_macchiato]
rosewater = "#f4dbd6"
flamingo = "#f0c6c6"
pink = "#f5bde6"
mauve = "#c6a0f6"
red = "#ed8796"
maroon = "#ee99a0"
peach = "#f5a97f"
yellow = "#eed49f"
green = "#a6da95"
teal = "#8bd5ca"
sky = "#91d7e3"
sapphire = "#7dc4e4"
blue = "#8aadf4"
lavender = "#b7bdf8"
text = "#cad3f5"
subtext1 = "#b8c0e0"
subtext0 = "#a5adcb"
overlay2 = "#939ab7"
overlay1 = "#8087a2"
overlay0 = "#6e738d"
surface2 = "#5b6078"
surface1 = "#494d64"
surface0 = "#363a4f"
base = "#24273a"
mantle = "#1e2030"
crust = "#181926"
`
}

func (s *StarshipIntegration) getRosePineConfig(variant string) string {
	// Default to rose_pine if no variant specified
	if variant == "" {
		variant = "default"
	}

	paletteName := "rose_pine"
	if variant != "default" {
		paletteName = "rose_pine_" + variant
	}

	return `# 🌹 Rosé Pine Starship Configuration (All Variants)
"$schema" = 'https://starship.rs/config-schema.json'

command_timeout = 3000

format = """
[](love)\
$os\
$username\
[](bg:gold fg:love)\
$directory\
[](bg:foam fg:gold)\
$git_branch\
$git_status\
[](fg:foam bg:pine)\
$c\
$rust\
$golang\
$nodejs\
$php\
$java\
$kotlin\
$haskell\
$python\
[](fg:pine bg:iris)\
$conda\
[](fg:iris bg:rose)\
$time\
[ ](fg:rose)\
$line_break\
$character"""

palette = '` + paletteName + `'

[os]
disabled = false
style = "bg:love fg:base"

[os.symbols]
Windows = ""
Ubuntu = "󰕈"
Macos = "󰀵"
Linux = "󰌽"
Debian = "󰣚"
Redhat = "󱄛"

[username]
show_always = true
style_user = "bg:love fg:base"
style_root = "bg:love fg:base"
format = '[ $user]($style)'

[directory]
style = "bg:gold fg:base"
format = "[ $path ]($style)"
truncation_length = 3
truncation_symbol = "…/"

[git_branch]
symbol = ""
style = "bg:foam"
format = '[[ $symbol $branch ](fg:base bg:foam)]($style)'

[git_status]
style = "bg:foam"
format = '[[($all_status$ahead_behind )](fg:base bg:foam)]($style)'

[nodejs]
symbol = ""
style = "bg:pine"
format = '[[ $symbol( $version) ](fg:base bg:pine)]($style)'

[c]
symbol = " "
style = "bg:pine"
format = '[[ $symbol( $version) ](fg:base bg:pine)]($style)'

[rust]
symbol = ""
style = "bg:pine"
format = '[[ $symbol( $version) ](fg:base bg:pine)]($style)'

[golang]
symbol = ""
style = "bg:pine"
format = '[[ $symbol( $version) ](fg:base bg:pine)]($style)'

[php]
symbol = ""
style = "bg:pine"
format = '[[ $symbol( $version) ](fg:base bg:pine)]($style)'
detect_extensions = ['php']
detect_files = ['composer.json', '.php-version']
detect_folders = ['vendor']

[java]
symbol = " "
style = "bg:pine"
format = '[[ $symbol( $version) ](fg:base bg:pine)]($style)'

[kotlin]
symbol = ""
style = "bg:pine"
format = '[[ $symbol( $version) ](fg:base bg:pine)]($style)'

[haskell]
symbol = ""
style = "bg:pine"
format = '[[ $symbol( $version) ](fg:base bg:pine)]($style)'

[python]
symbol = ""
style = "bg:pine"
format = '[[ $symbol( $version)(\(#$virtualenv\)) ](fg:base bg:pine)]($style)'

[docker_context]
symbol = ""
style = "bg:iris"
format = '[[ $symbol( $context) ](fg:base bg:iris)]($style)'

[conda]
symbol = "  "
style = "fg:base bg:iris"
format = '[$symbol$environment ]($style)'
ignore_base = false

[time]
disabled = false
time_format = "%R"
style = "bg:rose"
format = '[[  $time ](fg:base bg:rose)]($style)'

[line_break]
disabled = true

[character]
disabled = false
success_symbol = '[❯](bold fg:foam)'
error_symbol = '[❯](bold fg:love)'
vimcmd_symbol = '[❮](bold fg:foam)'
vimcmd_replace_one_symbol = '[❮](bold fg:rose)'
vimcmd_replace_symbol = '[❮](bold fg:rose)'
vimcmd_visual_symbol = '[❮](bold fg:gold)'

[cmd_duration]
show_milliseconds = true
format = " in $duration "
style = "bg:rose"
disabled = false
show_notifications = true
min_time_to_notify = 45000


# 🌑 Rosé Pine
[palettes.rose_pine]
base = "#191724"
surface = "#1f1d2e"
overlay = "#26233a"
muted = "#6e6a86"
subtle = "#908caa"
text = "#e0def4"
love = "#eb6f92"
gold = "#f6c177"
rose = "#ebbcba"
pine = "#31748f"
foam = "#9ccfd8"
iris = "#c4a7e7"
highlight_low = "#21202e"
highlight_med = "#403d52"
highlight_high = "#524f67"


# 🌙 Rosé Pine Moon
[palettes.rose_pine_moon]
base = "#232136"
surface = "#2a273f"
overlay = "#393552"
muted = "#6e6a86"
subtle = "#908caa"
text = "#e0def4"
love = "#eb6f92"
gold = "#f6c177"
rose = "#ea9a97"
pine = "#3e8fb0"
foam = "#9ccfd8"
iris = "#c4a7e7"
highlight_low = "#2a283e"
highlight_med = "#44415a"
highlight_high = "#56526e"


# 🌅 Rosé Pine Dawn
[palettes.rose_pine_dawn]
base = "#faf4ed"
surface = "#fffaf3"
overlay = "#f2e9e1"
muted = "#9893a5"
subtle = "#797593"
text = "#575279"
love = "#b4637a"
gold = "#ea9d34"
rose = "#d7827e"
pine = "#286983"
foam = "#56949f"
iris = "#907aa9"
highlight_low = "#f4ede8"
highlight_med = "#dfdad9"
highlight_high = "#cecacd"
`
}
