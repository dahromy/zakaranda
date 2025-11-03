package integrations

import (
	"sync"
	"testing"
)

// TestGetVSCodeVariantsCaching verifies that VS Code variants are cached
func TestGetVSCodeVariantsCaching(t *testing.T) {
	// First call should populate the cache
	variants1 := GetVSCodeVariants()
	
	// Second call should return the same cached slice
	variants2 := GetVSCodeVariants()
	
	// Verify the slices have the same length
	if len(variants1) != len(variants2) {
		t.Error("Expected variants to have the same length")
	}
	
	// Verify caching is working by comparing slice headers
	// If caching works, both should point to the same underlying array
	if len(variants1) > 0 && len(variants2) > 0 {
		if &variants1[0] != &variants2[0] {
			t.Error("Expected variants to be cached and return the same slice")
		}
	}
}

// TestGetVSCodeVariantsContent verifies variant content is correct
func TestGetVSCodeVariantsContent(t *testing.T) {
	// Reset cache to test fresh scan
	vscodeVariantsCache = nil
	vscodeVariantsMutex = sync.Once{}
	
	variants := GetVSCodeVariants()
	
	// Verify all variants have required fields
	for _, variant := range variants {
		if variant.Name == "" {
			t.Error("Variant missing name")
		}
		if variant.ConfigDir == "" {
			t.Error("Variant missing config directory")
		}
		if variant.CLICommand == "" {
			t.Error("Variant missing CLI command")
		}
		if variant.AppPath == "" {
			t.Error("Variant missing app path")
		}
	}
}

// BenchmarkGetVSCodeVariants benchmarks the performance of getting VS Code variants
func BenchmarkGetVSCodeVariants(b *testing.B) {
	// Clear cache before benchmark
	vscodeVariantsCache = nil
	vscodeVariantsMutex = sync.Once{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetVSCodeVariants()
	}
}

// BenchmarkGetVSCodeVariantsCached benchmarks cached variant retrieval
func BenchmarkGetVSCodeVariantsCached(b *testing.B) {
	// Populate cache
	GetVSCodeVariants()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetVSCodeVariants()
	}
}

// BenchmarkStripJSONComments benchmarks JSON comment stripping
func BenchmarkStripJSONComments(b *testing.B) {
	testJSON := `{
		// This is a comment
		"key1": "value1",
		/* Multi-line
		   comment */
		"key2": "value2",
		"key3": "value3" // Inline comment
	}`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stripJSONComments(testJSON)
	}
}
