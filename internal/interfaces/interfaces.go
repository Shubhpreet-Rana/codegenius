package interfaces

import (
	"time"
)

// GitRepository defines the contract for Git operations
type GitRepository interface {
	GetDiff() (string, error)
	GetChangedFiles() ([]string, error)
	GetCurrentBranch() (string, error)
	GetRecentCommits() ([]string, error)
	HasStagedChanges() (bool, error)
	CommitWithMessage(message string) error
	EditCommitMessage(message string) (string, error)
	AnalyzeDiffContext(diff string, ignorePatterns []string) (string, []string)
}

// AIProvider defines the contract for AI interactions
type AIProvider interface {
	GenerateCommitMessage(diff string, files []string, branchName, additionalContext string) (string, error)
	AnalyzeCode(code, analysisType string) (string, error)
	AddInteraction(interactionType, prompt, response, feedback string)
	GetContextualPrompt(basePrompt string) string
}

// ConfigManager defines the contract for configuration management
type ConfigManager interface {
	Load() error
	Save() error
	Initialize() error
	ShouldIgnoreFile(filename string) bool
	GetProject() ProjectConfig
	GetAI() AIConfig
	GetReview() ReviewConfig
}

// HistoryManager defines the contract for work history management
type HistoryManager interface {
	Load() error
	Save() error
	AddEntry(message string) error
	Display(monthYear string) error
	GetStats() map[string]interface{}
	FilterByMonthYear(monthYear string) []HistoryEntry
}

// CodeReviewer defines the contract for code review operations
type CodeReviewer interface {
	PerformReview(diff, reviewType string) (*ReviewResult, error)
	HandleInteractive(diff string) error
	DisplayResults(review *ReviewResult)
	GetSupportedTypes() []string
}

// Data structures (these will be shared across packages)
type ProjectConfig struct {
	Name        string   `yaml:"name"`
	Language    string   `yaml:"language"`
	Overview    string   `yaml:"overview"`
	Scopes      []string `yaml:"scopes"`
	Standards   string   `yaml:"standards"`
	IgnoreFiles []string `yaml:"ignore_files"`
}

type AIConfig struct {
	Model            string            `yaml:"model"`
	ContextTemplates map[string]string `yaml:"context_templates"`
	MaxTokens        int               `yaml:"max_tokens"`
}

type ReviewConfig struct {
	EnabledTypes     []string          `yaml:"enabled_types"`
	TextOnly         bool              `yaml:"text_only"`
	SecurityPatterns []string          `yaml:"security_patterns"`
	CustomRules      map[string]string `yaml:"custom_rules"`
}

type HistoryEntry struct {
	Date    string `json:"date"`
	Summary string `json:"summary"`
}

type ReviewResult struct {
	Type        string       `json:"type"`
	Issues      []ReviewItem `json:"issues"`
	Suggestions []ReviewItem `json:"suggestions"`
	Summary     string       `json:"summary"`
}

type ReviewItem struct {
	Line       int    `json:"line"`
	File       string `json:"file"`
	Category   string `json:"category"`
	Severity   string `json:"severity"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
}

type AIInteraction struct {
	Type      string    `json:"type"`
	Prompt    string    `json:"prompt"`
	Response  string    `json:"response"`
	Feedback  string    `json:"feedback"`
	Timestamp time.Time `json:"timestamp"`
}

// Service represents the main application service with all dependencies
type Service struct {
	Git     GitRepository
	AI      AIProvider
	Config  ConfigManager
	History HistoryManager
	Review  CodeReviewer
}

// ServiceBuilder provides a fluent interface for building services
type ServiceBuilder interface {
	WithGit(git GitRepository) ServiceBuilder
	WithAI(ai AIProvider) ServiceBuilder
	WithConfig(config ConfigManager) ServiceBuilder
	WithHistory(history HistoryManager) ServiceBuilder
	WithReview(review CodeReviewer) ServiceBuilder
	Build() *Service
}
