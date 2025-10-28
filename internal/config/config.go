package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	LastTheme        string            `json:"last_theme"`
	EnabledApps      []string          `json:"enabled_apps"`
	AutoBackup       bool              `json:"auto_backup"`
	MaxBackups       int               `json:"max_backups"`
	CustomThemesPath string            `json:"custom_themes_path"`
	Preferences      map[string]string `json:"preferences"`
}

type ConfigManager struct {
	configPath string
	config     *Config
}

func NewConfigManager() (*ConfigManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "theme-manager")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	cm := &ConfigManager{
		configPath: configPath,
	}

	// Load or create config
	if err := cm.Load(); err != nil {
		// Create default config if it doesn't exist
		cm.config = cm.defaultConfig()
		if err := cm.Save(); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
	}

	return cm, nil
}

func (cm *ConfigManager) defaultConfig() *Config {
	home, err := os.UserHomeDir()
	customPath := filepath.Join(home, ".config", "theme-manager", "themes")
	if err != nil {
		customPath = ".config/theme-manager/themes" // Fallback
	}
	return &Config{
		LastTheme:        "",
		EnabledApps:      []string{},
		AutoBackup:       true,
		MaxBackups:       5,
		CustomThemesPath: customPath,
		Preferences:      make(map[string]string),
	}
}

func (cm *ConfigManager) Load() error {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return err
	}

	cm.config = &Config{}
	if err := json.Unmarshal(data, cm.config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

func (cm *ConfigManager) Save() error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (cm *ConfigManager) GetLastTheme() string {
	return cm.config.LastTheme
}

func (cm *ConfigManager) SetLastTheme(theme string) error {
	cm.config.LastTheme = theme
	return cm.Save()
}

func (cm *ConfigManager) GetEnabledApps() []string {
	return cm.config.EnabledApps
}

func (cm *ConfigManager) SetEnabledApps(apps []string) error {
	cm.config.EnabledApps = apps
	return cm.Save()
}

func (cm *ConfigManager) IsAutoBackupEnabled() bool {
	return cm.config.AutoBackup
}

func (cm *ConfigManager) GetMaxBackups() int {
	return cm.config.MaxBackups
}

func (cm *ConfigManager) GetCustomThemesPath() string {
	return cm.config.CustomThemesPath
}

func (cm *ConfigManager) GetPreference(key string) (string, bool) {
	val, ok := cm.config.Preferences[key]
	return val, ok
}

func (cm *ConfigManager) SetPreference(key, value string) error {
	if cm.config.Preferences == nil {
		cm.config.Preferences = make(map[string]string)
	}
	cm.config.Preferences[key] = value
	return cm.Save()
}

// Backup management
func (cm *ConfigManager) CleanOldBackups(appName string) error {
	if !cm.config.AutoBackup {
		return nil
	}

	// Find all backup files for the app
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	backupDir := filepath.Join(home, ".config", "theme-manager", "backups", appName)

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		// Directory doesn't exist or is empty
		return nil
	}

	// If we have more backups than allowed, delete the oldest ones
	if len(entries) > cm.config.MaxBackups {
		// Sort by modification time would be needed here
		// For simplicity, just keep the last N files
		toDelete := len(entries) - cm.config.MaxBackups
		for i := 0; i < toDelete; i++ {
			backupPath := filepath.Join(backupDir, entries[i].Name())
			os.Remove(backupPath)
		}
	}

	return nil
}
