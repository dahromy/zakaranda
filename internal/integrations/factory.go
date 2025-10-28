package integrations

// GetAllIntegrations returns all available integrations
func GetAllIntegrations() []Integration {
	return []Integration{
		NewVSCodeIntegration(),
		NewAlacrittyIntegration(),
		NewWarpIntegration(),
		NewITerm2Integration(),
		NewStarshipIntegration(),
		NewZedIntegration(),
		NewWallpaperIntegration(),
		NewSlackIntegration(),
	}
}
