package theme

// SanitizeFileName removes spaces and special characters from a filename
func SanitizeFileName(name string) string {
	// Remove spaces and special characters
	result := ""
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_' {
			result += string(char)
		} else if char == ' ' {
			result += "_"
		}
	}
	return result
}

