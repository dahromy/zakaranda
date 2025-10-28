package theme

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ThemePreview struct {
	theme Theme
}

func NewThemePreview(theme Theme) *ThemePreview {
	return &ThemePreview{theme: theme}
}

// Render generates a visual preview of the theme
func (tp *ThemePreview) Render() string {
	var preview strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(tp.theme.Colors.Blue)).
		MarginBottom(1)

	preview.WriteString(titleStyle.Render(fmt.Sprintf("Preview: %s", tp.theme.Name)))
	preview.WriteString("\n\n")

	// Color palette
	preview.WriteString(tp.renderColorPalette())
	preview.WriteString("\n\n")

	// Code example
	preview.WriteString(tp.renderCodeExample())
	preview.WriteString("\n\n")

	// Terminal example
	preview.WriteString(tp.renderTerminalExample())

	return preview.String()
}

func (tp *ThemePreview) renderColorPalette() string {
	var palette strings.Builder

	palette.WriteString("Color Palette:\n")

	colors := []struct {
		name  string
		value string
	}{
		{"Background", tp.theme.Colors.Background},
		{"Foreground", tp.theme.Colors.Foreground},
		{"Black", tp.theme.Colors.Black},
		{"Red", tp.theme.Colors.Red},
		{"Green", tp.theme.Colors.Green},
		{"Yellow", tp.theme.Colors.Yellow},
		{"Blue", tp.theme.Colors.Blue},
		{"Magenta", tp.theme.Colors.Magenta},
		{"Cyan", tp.theme.Colors.Cyan},
		{"White", tp.theme.Colors.White},
	}

	for _, color := range colors {
		colorBox := lipgloss.NewStyle().
			Background(lipgloss.Color(color.value)).
			Foreground(lipgloss.Color(tp.theme.Colors.Foreground)).
			Padding(0, 2).
			Render("  ")

		label := lipgloss.NewStyle().
			Foreground(lipgloss.Color(tp.theme.Colors.Foreground)).
			Render(fmt.Sprintf("%-12s %s", color.name, color.value))

		palette.WriteString(fmt.Sprintf("  %s %s\n", colorBox, label))
	}

	return palette.String()
}

func (tp *ThemePreview) renderCodeExample() string {
	var code strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(tp.theme.Colors.Blue)).
		Bold(true)

	code.WriteString(headerStyle.Render("Code Example:"))
	code.WriteString("\n")

	// Background box
	codeBoxStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(tp.theme.Colors.Background)).
		Foreground(lipgloss.Color(tp.theme.Colors.Foreground)).
		Padding(1, 2).
		MarginTop(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(tp.theme.Colors.Black))

	// Syntax-highlighted code
	keywordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Magenta))
	functionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Blue))
	stringStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Green))
	commentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.BrightBlack))
	numberStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Magenta))

	codeContent := fmt.Sprintf(`%s main() {
    %s name %s %s
    %s count %s %s
    %s(%s)
    %s
}`,
		keywordStyle.Render("func"),
		keywordStyle.Render("var"),
		lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Foreground)).Render("="),
		stringStyle.Render(`"Theme Manager"`),
		keywordStyle.Render("var"),
		lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Foreground)).Render("="),
		numberStyle.Render("42"),
		functionStyle.Render("println"),
		stringStyle.Render(`"Hello, World!"`),
		commentStyle.Render("// This is a comment"),
	)

	code.WriteString(codeBoxStyle.Render(codeContent))

	return code.String()
}

func (tp *ThemePreview) renderTerminalExample() string {
	var terminal strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(tp.theme.Colors.Blue)).
		Bold(true)

	terminal.WriteString(headerStyle.Render("Terminal Example:"))
	terminal.WriteString("\n")

	termBoxStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(tp.theme.Colors.Background)).
		Foreground(lipgloss.Color(tp.theme.Colors.Foreground)).
		Padding(1, 2).
		MarginTop(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(tp.theme.Colors.Black))

	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Green))
	pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Cyan))
	commandStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Yellow))

	termContent := fmt.Sprintf(`%s %s %s
%s
%s
%s %s %s
%s`,
		promptStyle.Render("user@host"),
		pathStyle.Render("~/projects"),
		lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Blue)).Render("❯"),
		commandStyle.Render("ls -la"),
		lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Foreground)).Render("total 42"),
		promptStyle.Render("user@host"),
		pathStyle.Render("~/projects"),
		lipgloss.NewStyle().Foreground(lipgloss.Color(tp.theme.Colors.Blue)).Render("❯"),
		commandStyle.Render("git status"),
	)

	terminal.WriteString(termBoxStyle.Render(termContent))

	return terminal.String()
}

// RenderCompact creates a compact single-line preview
func (tp *ThemePreview) RenderCompact() string {
	colors := []string{
		tp.theme.Colors.Red,
		tp.theme.Colors.Green,
		tp.theme.Colors.Yellow,
		tp.theme.Colors.Blue,
		tp.theme.Colors.Magenta,
		tp.theme.Colors.Cyan,
	}

	var boxes strings.Builder
	for _, color := range colors {
		box := lipgloss.NewStyle().
			Background(lipgloss.Color(color)).
			Render("  ")
		boxes.WriteString(box)
	}

	return boxes.String()
}

// RenderColorGrid creates a grid of all colors
func (tp *ThemePreview) RenderColorGrid() string {
	var grid strings.Builder

	grid.WriteString("Color Grid:\n\n")

	allColors := []struct {
		name  string
		value string
	}{
		{"Bg", tp.theme.Colors.Background},
		{"Fg", tp.theme.Colors.Foreground},
		{"Blk", tp.theme.Colors.Black},
		{"Red", tp.theme.Colors.Red},
		{"Grn", tp.theme.Colors.Green},
		{"Ylw", tp.theme.Colors.Yellow},
		{"Blu", tp.theme.Colors.Blue},
		{"Mag", tp.theme.Colors.Magenta},
		{"Cyn", tp.theme.Colors.Cyan},
		{"Wht", tp.theme.Colors.White},
		{"BBlk", tp.theme.Colors.BrightBlack},
		{"BRed", tp.theme.Colors.BrightRed},
		{"BGrn", tp.theme.Colors.BrightGreen},
		{"BYlw", tp.theme.Colors.BrightYellow},
		{"BBlu", tp.theme.Colors.BrightBlue},
		{"BMag", tp.theme.Colors.BrightMagenta},
		{"BCyn", tp.theme.Colors.BrightCyan},
		{"BWht", tp.theme.Colors.BrightWhite},
	}

	// Create grid rows
	row := 0
	for i, color := range allColors {
		box := lipgloss.NewStyle().
			Background(lipgloss.Color(color.value)).
			Width(8).
			Align(lipgloss.Center).
			Render(color.name)

		grid.WriteString(box)

		if (i+1)%6 == 0 {
			grid.WriteString("\n")
			row++
		}
	}

	return grid.String()
}
