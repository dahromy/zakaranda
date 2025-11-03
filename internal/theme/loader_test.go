package theme

import (
	"testing"
)

// TestGetBuiltInThemesCaching verifies that built-in themes are cached
func TestGetBuiltInThemesCaching(t *testing.T) {
	// Clear cache for testing
	builtInThemesCache = nil
	
	// First call should populate the cache
	themes1 := GetBuiltInThemes()
	if len(themes1) == 0 {
		t.Fatal("Expected at least one built-in theme")
	}
	
	// Second call should return the same cached slice
	themes2 := GetBuiltInThemes()
	
	// Verify the slices are identical (same memory address)
	if &themes1[0] != &themes2[0] {
		t.Error("Expected themes to be cached and return the same slice")
	}
	
	// Verify we have the expected number of themes (Nord + 4 Catppuccin + 3 Rose Pine = 8)
	expectedCount := 8
	if len(themes1) != expectedCount {
		t.Errorf("Expected %d themes, got %d", expectedCount, len(themes1))
	}
}

// TestGetBuiltInThemesContent verifies theme content is correct
func TestGetBuiltInThemesContent(t *testing.T) {
	themes := GetBuiltInThemes()
	
	// Verify all themes have required fields
	for _, theme := range themes {
		if theme.Name == "" {
			t.Error("Theme missing name")
		}
		if theme.Description == "" {
			t.Error("Theme missing description")
		}
		if theme.Colors.Background == "" {
			t.Errorf("Theme %s missing background color", theme.Name)
		}
		if theme.Colors.Foreground == "" {
			t.Errorf("Theme %s missing foreground color", theme.Name)
		}
	}
}

// BenchmarkGetBuiltInThemes benchmarks the performance of getting built-in themes
func BenchmarkGetBuiltInThemes(b *testing.B) {
	// Clear cache before benchmark
	builtInThemesCache = nil
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetBuiltInThemes()
	}
}

// BenchmarkGetBuiltInThemesCached benchmarks cached theme retrieval
func BenchmarkGetBuiltInThemesCached(b *testing.B) {
	// Populate cache
	GetBuiltInThemes()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetBuiltInThemes()
	}
}
