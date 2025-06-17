package container

import (
	"fmt"

	"github.com/codegenius/cli/internal/interfaces"
)

// Container implements dependency injection for the application
type Container struct {
	git     interfaces.GitRepository
	ai      interfaces.AIProvider
	config  interfaces.ConfigManager
	history interfaces.HistoryManager
	review  interfaces.CodeReviewer
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return &Container{}
}

// WithGit sets the Git repository implementation
func (c *Container) WithGit(git interfaces.GitRepository) interfaces.ServiceBuilder {
	c.git = git
	return c
}

// WithAI sets the AI provider implementation
func (c *Container) WithAI(ai interfaces.AIProvider) interfaces.ServiceBuilder {
	c.ai = ai
	return c
}

// WithConfig sets the configuration manager implementation
func (c *Container) WithConfig(config interfaces.ConfigManager) interfaces.ServiceBuilder {
	c.config = config
	return c
}

// WithHistory sets the history manager implementation
func (c *Container) WithHistory(history interfaces.HistoryManager) interfaces.ServiceBuilder {
	c.history = history
	return c
}

// WithReview sets the code reviewer implementation
func (c *Container) WithReview(review interfaces.CodeReviewer) interfaces.ServiceBuilder {
	c.review = review
	return c
}

// Build creates and validates the service with all dependencies
func (c *Container) Build() *interfaces.Service {
	if err := c.validate(); err != nil {
		panic(fmt.Sprintf("Invalid service configuration: %v", err))
	}

	return &interfaces.Service{
		Git:     c.git,
		AI:      c.ai,
		Config:  c.config,
		History: c.history,
		Review:  c.review,
	}
}

// validate ensures all required dependencies are provided
func (c *Container) validate() error {
	if c.git == nil {
		return fmt.Errorf("GitRepository is required")
	}
	if c.ai == nil {
		return fmt.Errorf("AIProvider is required")
	}
	if c.config == nil {
		return fmt.Errorf("ConfigManager is required")
	}
	if c.history == nil {
		return fmt.Errorf("HistoryManager is required")
	}
	if c.review == nil {
		return fmt.Errorf("CodeReviewer is required")
	}
	return nil
}

// BuildDefault creates a service with default implementations
func BuildDefault() (*interfaces.Service, error) {
	// This will import and create default implementations
	// We'll implement this after updating the individual packages
	return nil, fmt.Errorf("BuildDefault not yet implemented - use BuildWithConfig instead")
}

// ServiceFactory provides factory methods for common service configurations
type ServiceFactory struct{}

// NewServiceFactory creates a new service factory
func NewServiceFactory() *ServiceFactory {
	return &ServiceFactory{}
}

// CreateProductionService creates a service configured for production use
func (sf *ServiceFactory) CreateProductionService() (*interfaces.Service, error) {
	// Will be implemented after updating individual packages
	return nil, fmt.Errorf("not yet implemented")
}

// CreateTestService creates a service configured for testing with mocks
func (sf *ServiceFactory) CreateTestService() (*interfaces.Service, error) {
	// Will be implemented with mock implementations
	return nil, fmt.Errorf("not yet implemented")
}
