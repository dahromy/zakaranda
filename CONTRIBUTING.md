# Contributing to Zakaranda

Thank you for your interest in contributing to Zakaranda! 🎨

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- macOS (for testing all integrations)
- Git

### Setup Development Environment

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/yourusername/zakaranda.git
   cd zakaranda
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Build the project**
   ```bash
   go build -o zakaranda ./cmd/zakaranda
   ```

4. **Run tests**
   ```bash
   cd cmd/zakaranda
   go test -v
   ```

## 📝 Development Guidelines

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Add comments for exported functions and types
- Keep functions focused and small

### Project Structure

```
zakaranda/
├── cmd/zakaranda/
│   ├── main.go           # Main application & TUI
│   ├── config.go         # Configuration management
│   ├── theme_loader.go   # Theme loading
│   ├── preview.go        # Preview UI
│   ├── *_integration.go  # Application integrations
│   └── theme_test.go     # Tests
```

### Adding a New Theme

1. **Define theme colors** in the respective integration files
2. **Add theme to the menu** in `main.go`
3. **Test with all integrations**
4. **Update README.md** with theme information

Example theme structure:
```go
type Theme struct {
    Name   string
    Colors struct {
        Background string
        Foreground string
        Black      string
        Red        string
        // ... other colors
    }
}
```

### Adding a New Integration

1. **Create a new file** (e.g., `myapp.go`)

2. **Implement the AppIntegration interface**:
   ```go
   type AppIntegration interface {
       Name() string
       ConfigPath() string
       IsInstalled() bool
       Apply(theme Theme) error
   }
   ```

3. **Add to integrations list** in `main.go`:
   ```go
   apps := []AppIntegration{
       // ... existing integrations
       NewMyAppIntegration(),
   }
   ```

4. **Create backup before modifying** config files:
   ```go
   backupPath := configPath + ".backup"
   os.Rename(configPath, backupPath)
   ```

5. **Write tests** in `theme_test.go`

6. **Update documentation** in README.md

### Testing

- Write unit tests for new features
- Test with all supported themes
- Test on a clean macOS installation if possible
- Verify backup files are created correctly

### Commit Messages

Use conventional commit format:

```
feat: add support for Kitty terminal
fix: correct Starship powerline rendering
docs: update installation instructions
refactor: simplify theme loading logic
test: add tests for VS Code integration
```

## 🐛 Reporting Bugs

When reporting bugs, please include:

- **OS version** (e.g., macOS 14.0)
- **Go version** (`go version`)
- **Steps to reproduce**
- **Expected behavior**
- **Actual behavior**
- **Error messages** (if any)
- **Screenshots** (if applicable)

## 💡 Feature Requests

We welcome feature requests! Please:

1. Check if the feature already exists
2. Describe the use case
3. Explain why it would be useful
4. Provide examples if possible

## 🔄 Pull Request Process

1. **Create a feature branch**
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes**
   - Write clean, documented code
   - Add tests if applicable
   - Update documentation

3. **Test your changes**
   ```bash
   go test ./...
   go build -o zakaranda ./cmd/zakaranda
   ./zakaranda
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add my new feature"
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/my-new-feature
   ```

6. **Create a Pull Request**
   - Provide a clear description
   - Reference any related issues
   - Include screenshots if UI changes

## 📚 Resources

- [Go Documentation](https://go.dev/doc/)
- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea)
- [Nord Theme Spec](https://www.nordtheme.com/)
- [Catppuccin Spec](https://github.com/catppuccin/catppuccin)
- [Rosé Pine Spec](https://rosepinetheme.com/)

## 🎨 Theme Color Guidelines

When adding themes, ensure:

- All 16 ANSI colors are defined
- Background and foreground colors have good contrast
- Colors are consistent across integrations
- Powerline characters render correctly in Starship

## ✅ Checklist for New Integrations

- [ ] Integration file created (`myapp.go`)
- [ ] Implements `AppIntegration` interface
- [ ] Creates backup before modifying configs
- [ ] Handles errors gracefully
- [ ] Added to integrations list in `main.go`
- [ ] Tests written
- [ ] Documentation updated
- [ ] Tested with all themes

## 🤝 Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Help others learn and grow
- Focus on what's best for the community

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Zakaranda! 🎉

