package integrations

import "zakaranda/internal/theme"

// Integration defines the interface that all application integrations must implement
type Integration interface {
	// Name returns the display name of the integration
	Name() string

	// IsInstalled checks if the application is installed on the system
	IsInstalled() bool

	// Apply applies the given theme to the application
	Apply(theme theme.Theme) error

	// ConfigPath returns the path to the application's configuration file
	ConfigPath() string
}
