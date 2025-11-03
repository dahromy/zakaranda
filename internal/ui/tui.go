package ui

import (
	"fmt"
	"os"
	"zakaranda/internal/integrations"
	"zakaranda/internal/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7aa2f7")).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bb9af7")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c0caf5"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ece6a")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89"))
)

type state int

const (
	selectingTheme state = iota
	selectingVariant
	previewingTheme
	selectingApps
	selectingVSCodeVariant
	applying
	complete
)

type model struct {
	baseThemes         []BaseTheme
	themes             []Theme
	selectedBaseTheme  int
	selectedVariant    int
	selectedTheme      int
	apps               []AppIntegration
	selectedApps       map[int]bool
	cursor             int
	state              state
	err                error
	results            []string
	vscodeVariants     []VSCodeVariant
	selectedVSCVariant int
}

// Type aliases for imported types
type Theme = theme.Theme
type AppIntegration = integrations.Integration
type BaseTheme = theme.BaseTheme
type ThemeVariant = theme.ThemeVariant
type VSCodeVariant = integrations.VSCodeVariant

func initialModel() model {
	baseThemes := theme.GetBuiltInBaseThemes()
	themes := theme.GetBuiltInThemes()
	apps := integrations.GetAllIntegrations()
	
	// Cache VS Code variants during initialization to avoid repeated calls
	vscodeVariants := integrations.GetVSCodeVariants()

	return model{
		baseThemes:        baseThemes,
		themes:            themes,
		selectedBaseTheme: 0,
		selectedVariant:   0,
		selectedTheme:     0,
		apps:              apps,
		selectedApps:      make(map[int]bool),
		cursor:            0,
		state:             selectingTheme,
		vscodeVariants:    vscodeVariants,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.state == selectingTheme {
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.state == selectingVariant {
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.state == selectingApps {
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.state == selectingVSCodeVariant {
				if m.cursor > 0 {
					m.cursor--
				}
			}

		case "down", "j":
			if m.state == selectingTheme {
				if m.cursor < len(m.baseThemes)-1 {
					m.cursor++
				}
			} else if m.state == selectingVariant {
				baseTheme := m.baseThemes[m.selectedBaseTheme]
				if m.cursor < len(baseTheme.Variants)-1 {
					m.cursor++
				}
			} else if m.state == selectingApps {
				if m.cursor < len(m.apps)-1 {
					m.cursor++
				}
			} else if m.state == selectingVSCodeVariant {
				if m.cursor < len(m.vscodeVariants)-1 {
					m.cursor++
				}
			}

		case "enter":
			if m.state == selectingTheme {
				m.selectedBaseTheme = m.cursor
				baseTheme := m.baseThemes[m.selectedBaseTheme]
				// If theme has only one variant, skip variant selection
				if len(baseTheme.Variants) == 1 {
					m.selectedVariant = 0
					// Calculate the theme index in the flattened themes list
					m.selectedTheme = m.calculateThemeIndex()
					m.state = previewingTheme
				} else {
					// Show variant selection
					m.cursor = 0
					m.state = selectingVariant
				}
			} else if m.state == selectingVariant {
				m.selectedVariant = m.cursor
				// Calculate the theme index in the flattened themes list
				m.selectedTheme = m.calculateThemeIndex()
				m.state = previewingTheme
			} else if m.state == previewingTheme {
				m.cursor = 0
				m.state = selectingApps
			} else if m.state == selectingApps {
				// Check if VS Code is selected
				vscodeSelected := false
				vscodeIdx := -1
				for idx, selected := range m.selectedApps {
					if selected && m.apps[idx].Name() == "VS Code" {
						vscodeSelected = true
						vscodeIdx = idx
						break
					}
				}

				if vscodeSelected {
					// Use cached VS Code variants
					if len(m.vscodeVariants) > 1 {
						// Multiple variants available, show selection
						m.cursor = 0
						m.state = selectingVSCodeVariant
					} else if len(m.vscodeVariants) == 1 {
						// Only one variant, use it directly
						m.selectedVSCVariant = 0
						m.updateVSCodeVariant(vscodeIdx)
						// After VS Code, proceed to apply themes
						return m, m.applyThemes()
					} else {
						// No variants found, proceed to apply themes
						return m, m.applyThemes()
					}
				} else {
					// VS Code not selected, proceed to apply themes
					return m, m.applyThemes()
				}
			} else if m.state == selectingVSCodeVariant {
				m.selectedVSCVariant = m.cursor
				// Find VS Code integration and update it
				for idx, app := range m.apps {
					if app.Name() == "VS Code" && m.selectedApps[idx] {
						m.updateVSCodeVariant(idx)
						break
					}
				}
				// After VS Code variant selection, proceed to apply themes
				return m, m.applyThemes()
			} else if m.state == complete {
				return m, tea.Quit
			}

		case " ":
			if m.state == selectingApps {
				m.selectedApps[m.cursor] = !m.selectedApps[m.cursor]
			}

		case "p":
			// Toggle preview from theme selection
			if m.state == selectingTheme {
				m.selectedTheme = m.cursor
				m.state = previewingTheme
			}

		case "esc":
			if m.state == selectingVariant {
				m.state = selectingTheme
				m.cursor = m.selectedBaseTheme
			} else if m.state == previewingTheme {
				// Go back to variant selection if theme has multiple variants
				baseTheme := m.baseThemes[m.selectedBaseTheme]
				if len(baseTheme.Variants) > 1 {
					m.state = selectingVariant
					m.cursor = m.selectedVariant
				} else {
					m.state = selectingTheme
					m.cursor = m.selectedBaseTheme
				}
			} else if m.state == selectingApps {
				m.state = previewingTheme
				m.cursor = 0
				m.selectedApps = make(map[int]bool)
			} else if m.state == selectingVSCodeVariant {
				m.cursor = 0
				m.state = selectingApps
			}
		}

	case applyCompleteMsg:
		m.state = complete
		m.results = msg.results
		m.err = msg.err
	}

	return m, nil
}

type applyCompleteMsg struct {
	results []string
	err     error
}

func (m *model) updateVSCodeVariant(appIdx int) {
	if vscode, ok := m.apps[appIdx].(*integrations.VSCodeIntegration); ok {
		variant := m.vscodeVariants[m.selectedVSCVariant]
		vscode.SetVariant(variant)
	}
}

// calculateThemeIndex calculates the index of the selected theme in the flattened themes list
func (m *model) calculateThemeIndex() int {
	index := 0
	for i := 0; i < m.selectedBaseTheme; i++ {
		index += len(m.baseThemes[i].Variants)
	}
	index += m.selectedVariant
	return index
}

func (m model) applyThemes() tea.Cmd {
	return func() tea.Msg {
		var results []string
		theme := m.themes[m.selectedTheme]

		for idx, selected := range m.selectedApps {
			if selected {
				app := m.apps[idx]
				if !app.IsInstalled() {
					results = append(results, fmt.Sprintf("⚠️  %s: Not installed or not found", app.Name()))
					continue
				}

				err := app.Apply(theme)
				if err != nil {
					results = append(results, fmt.Sprintf("❌ %s: %v", app.Name(), err))
				} else {
					results = append(results, fmt.Sprintf("✅ %s: Theme applied successfully", app.Name()))
				}
			}
		}

		return applyCompleteMsg{results: results}
	}
}

func (m model) View() string {
	s := titleStyle.Render("🎨 Theme Manager") + "\n\n"

	switch m.state {
	case selectingTheme:
		s += normalStyle.Render("Select a theme:") + "\n\n"
		for i, baseTheme := range m.baseThemes {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, baseTheme.Name)) + "\n"
				s += dimStyle.Render(fmt.Sprintf("  %s", baseTheme.Description)) + "\n"
			} else {
				s += normalStyle.Render(fmt.Sprintf("%s %s", cursor, baseTheme.Name)) + "\n"
			}
		}
		s += "\n" + dimStyle.Render("↑/↓: navigate • enter: select • q: quit")

	case selectingVariant:
		baseTheme := m.baseThemes[m.selectedBaseTheme]
		s += successStyle.Render(fmt.Sprintf("Theme: %s", baseTheme.Name)) + "\n\n"
		s += normalStyle.Render("Select variant:") + "\n\n"
		for i, variant := range baseTheme.Variants {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, variant.DisplayName)) + "\n"
			} else {
				s += normalStyle.Render(fmt.Sprintf("%s %s", cursor, variant.DisplayName)) + "\n"
			}
		}
		s += "\n" + dimStyle.Render("↑/↓: navigate • enter: preview • esc: back • q: quit")

	case previewingTheme:
		themeToPreview := m.themes[m.selectedTheme]
		preview := theme.NewThemePreview(themeToPreview)
		s += preview.Render()
		s += "\n" + dimStyle.Render("enter: continue to app selection • esc: back • q: quit")

	case selectingApps:
		s += successStyle.Render(fmt.Sprintf("Theme: %s", m.themes[m.selectedTheme].Name)) + "\n\n"
		s += normalStyle.Render("Select applications to theme:") + "\n\n"
		for i, app := range m.apps {
			cursor := " "
			checkbox := "[ ]"
			if m.selectedApps[i] {
				checkbox = "[✓]"
			}

			status := ""
			if !app.IsInstalled() {
				status = dimStyle.Render(" (not found)")
			}

			if m.cursor == i {
				cursor = ">"
				s += selectedStyle.Render(fmt.Sprintf("%s %s %s%s", cursor, checkbox, app.Name(), status)) + "\n"
				// Only show config path if it's not empty
				if configPath := app.ConfigPath(); configPath != "" {
					s += dimStyle.Render(fmt.Sprintf("   %s", configPath)) + "\n"
				}
			} else {
				s += normalStyle.Render(fmt.Sprintf("%s %s %s%s", cursor, checkbox, app.Name(), status)) + "\n"
			}
		}
		s += "\n" + dimStyle.Render("↑/↓: navigate • space: toggle • enter: apply • esc: back • q: quit")

	case selectingVSCodeVariant:
		s += successStyle.Render(fmt.Sprintf("Theme: %s", m.themes[m.selectedTheme].Name)) + "\n\n"
		s += normalStyle.Render("Multiple VS Code variants detected. Select one:") + "\n\n"
		for i, variant := range m.vscodeVariants {
			cursor := " "
			status := ""

			// Check if variant is actually installed
			if _, err := os.Stat(variant.AppPath); err == nil {
				status = successStyle.Render(" ✓")
			} else {
				status = dimStyle.Render(" (app not found)")
			}

			if m.cursor == i {
				cursor = ">"
				s += selectedStyle.Render(fmt.Sprintf("%s %s%s", cursor, variant.Name, status)) + "\n"
				s += dimStyle.Render(fmt.Sprintf("   Config: %s", variant.ConfigDir)) + "\n"
			} else {
				s += normalStyle.Render(fmt.Sprintf("%s %s%s", cursor, variant.Name, status)) + "\n"
			}
		}
		s += "\n" + dimStyle.Render("↑/↓: navigate • enter: select • esc: back • q: quit")

	case applying:
		s += normalStyle.Render("Applying themes...") + "\n"

	case complete:
		s += successStyle.Render("✨ Theme Application Complete!") + "\n\n"
		for _, result := range m.results {
			s += result + "\n"
		}
		s += "\n" + dimStyle.Render("Press enter or q to quit")
	}

	return s + "\n"
}

// Run starts the TUI application
func Run() error {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}
	return nil
}
