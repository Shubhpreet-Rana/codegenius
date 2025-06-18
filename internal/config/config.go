package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Shubhpreet-Rana/codegenius/internal/interfaces"
	"gopkg.in/yaml.v2"
)

const configFile = ".codegenius.yaml"

// Config represents the main configuration structure
type Config struct {
	Project interfaces.ProjectConfig `yaml:"project"`
	AI      interfaces.AIConfig      `yaml:"ai"`
	Review  interfaces.ReviewConfig  `yaml:"review"`
}

// Manager implements the ConfigManager interface
type Manager struct {
	config *Config
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{}
}

// Load reads and parses the configuration file
func (m *Manager) Load() error {
	config := getDefaultConfig()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		m.config = config
		return nil // Return default config
	}

	file, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	m.config = config
	return nil
}

// Save writes the configuration to file
func (m *Manager) Save() error {
	if m.config == nil {
		return fmt.Errorf("no configuration loaded")
	}

	data, err := yaml.Marshal(m.config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

// Initialize sets up the configuration with user input
func (m *Manager) Initialize() error {
	if m.config == nil {
		m.config = getDefaultConfig()
	}

	// Set basic project details
	m.config.Project.Language = detectProjectLanguage()

	// Save the configuration
	return m.Save()
}

// ShouldIgnoreFile checks if a file should be ignored based on patterns
func (m *Manager) ShouldIgnoreFile(filename string) bool {
	if m.config == nil {
		return false
	}

	for _, pattern := range m.config.Project.IgnoreFiles {
		if strings.Contains(filename, pattern) ||
			strings.HasSuffix(filename, strings.TrimPrefix(pattern, "*")) {
			return true
		}
	}
	return false
}

// GetProject returns the project configuration
func (m *Manager) GetProject() interfaces.ProjectConfig {
	if m.config == nil {
		return interfaces.ProjectConfig{}
	}
	return m.config.Project
}

// GetAI returns the AI configuration
func (m *Manager) GetAI() interfaces.AIConfig {
	if m.config == nil {
		return interfaces.AIConfig{}
	}
	return m.config.AI
}

// GetReview returns the review configuration
func (m *Manager) GetReview() interfaces.ReviewConfig {
	if m.config == nil {
		return interfaces.ReviewConfig{}
	}
	return m.config.Review
}

// GetConfig returns the full configuration (for backward compatibility)
func (m *Manager) GetConfig() *Config {
	return m.config
}

// SetProject updates the project configuration
func (m *Manager) SetProject(project interfaces.ProjectConfig) {
	if m.config == nil {
		m.config = getDefaultConfig()
	}
	m.config.Project = project
}

// SetAI updates the AI configuration
func (m *Manager) SetAI(ai interfaces.AIConfig) {
	if m.config == nil {
		m.config = getDefaultConfig()
	}
	m.config.AI = ai
}

// SetReview updates the review configuration
func (m *Manager) SetReview(review interfaces.ReviewConfig) {
	if m.config == nil {
		m.config = getDefaultConfig()
	}
	m.config.Review = review
}

// getDefaultConfig returns a configuration with sensible defaults
func getDefaultConfig() *Config {
	return &Config{
		Project: interfaces.ProjectConfig{
			Name:     "My Project",
			Language: "go",
			Overview: "A Go project using CodeGenius for intelligent commits and reviews",
			Scopes: []string{
				"core", "api", "docs", "deps", "scripts", "ci", "build",
			},
			Standards:   "https://golang.org/doc/effective_go.html",
			IgnoreFiles: []string{"go.mod", "go.sum", "*.lock", "node_modules/", ".git/"},
		},
		AI: interfaces.AIConfig{
			Model:     "gemini-2.0-flash",
			MaxTokens: 4000,
			ContextTemplates: map[string]string{
				"default": "This is a standard commit message generation request.",
				"bugfix":  "Focus on describing the bug that was fixed and its impact.",
				"feature": "Emphasize the new functionality and its benefits to users.",
			},
		},
		Review: interfaces.ReviewConfig{
			EnabledTypes: []string{"security", "performance", "style", "structure"},
			TextOnly:     true,
			SecurityPatterns: []string{
				`(?i)(password|secret|key|token)\s*[:=]\s*["'][^"']+["']`,
				`(?i)api[_-]?key\s*[:=]\s*["'][^"']+["']`,
				`(?i)(auth|bearer)\s*[:=]\s*["'][^"']+["']`,
			},
			CustomRules: map[string]string{},
		},
	}
}

// detectProjectLanguage attempts to detect the project language based on files
func detectProjectLanguage() string {
	files := []struct {
		pattern  string
		language string
	}{
		{"go.mod", "go"},
		{"package.json", "javascript"},
		{"requirements.txt", "python"},
		{"Cargo.toml", "rust"},
		{"pom.xml", "java"},
		{"Gemfile", "ruby"},
	}

	for _, file := range files {
		if _, err := os.Stat(file.pattern); err == nil {
			return file.language
		}
	}
	return "unknown"
}

// Backward compatibility functions

// Load reads and parses the configuration file (legacy function)
func Load() (*Config, error) {
	manager := NewManager()
	err := manager.Load()
	if err != nil {
		return nil, err
	}
	return manager.GetConfig(), nil
}
