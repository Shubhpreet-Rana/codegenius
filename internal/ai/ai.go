package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codegenius/cli/internal/interfaces"
)

const geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="

// SessionManager implements the AIProvider interface
type SessionManager struct {
	currentSession *Session
	config         interfaces.ConfigManager
}

// Session represents an AI conversation session
type Session struct {
	History     []interfaces.AIInteraction `json:"history"`
	Context     string                     `json:"context"`
	LastMessage string                     `json:"last_message"`
}

// GeminiRequest represents the structure for Gemini API requests
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

// Content represents content in a Gemini request
type Content struct {
	Parts []Part `json:"parts"`
}

// Part represents a part of content in a Gemini request
type Part struct {
	Text string `json:"text"`
}

// GeminiResponse represents the structure for Gemini API responses
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

// Candidate represents a candidate response from Gemini
type Candidate struct {
	Content Content `json:"content"`
}

// NewSessionManager creates a new AI session manager
func NewSessionManager(config interfaces.ConfigManager) interfaces.AIProvider {
	return &SessionManager{
		currentSession: &Session{
			History: make([]interfaces.AIInteraction, 0),
		},
		config: config,
	}
}

// GenerateCommitMessage generates a commit message based on git diff and context
func (sm *SessionManager) GenerateCommitMessage(diff string, files []string, branchName, additionalContext string) (string, error) {
	if err := sm.validateConfig(); err != nil {
		return "", err
	}

	prompt := sm.buildCommitPrompt(diff, files, branchName, additionalContext)

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	response, err := sm.callGeminiAPI(prompt, apiKey)
	if err != nil {
		return "", fmt.Errorf("AI API call failed: %v", err)
	}

	// Clean up the response
	message := strings.TrimSpace(response)
	message = strings.Trim(message, "`\"'")

	// Validate the generated message
	if message == "" {
		return "", fmt.Errorf("AI generated an empty commit message")
	}

	// Add interaction to session
	sm.AddInteraction("commit", prompt, message, "")

	return message, nil
}

// buildCommitPrompt constructs the prompt for commit message generation
func (sm *SessionManager) buildCommitPrompt(diff string, files []string, branchName, additionalContext string) string {
	var prompt strings.Builder

	aiConfig := sm.config.GetAI()

	// Get context template
	template := aiConfig.ContextTemplates["default"]
	if branchName != "" {
		if strings.Contains(branchName, "bug") || strings.Contains(branchName, "fix") {
			if bugfixTemplate, exists := aiConfig.ContextTemplates["bugfix"]; exists {
				template = bugfixTemplate
			}
		} else if strings.Contains(branchName, "feature") || strings.Contains(branchName, "feat") {
			if featureTemplate, exists := aiConfig.ContextTemplates["feature"]; exists {
				template = featureTemplate
			}
		}
	}

	prompt.WriteString(fmt.Sprintf("Context: %s\n\n", template))
	prompt.WriteString("Generate a concise, meaningful Git commit message based on the following changes:\n\n")

	if branchName != "" {
		prompt.WriteString(fmt.Sprintf("Branch: %s\n", branchName))
	}

	if len(files) > 0 {
		prompt.WriteString(fmt.Sprintf("Modified files: %s\n", strings.Join(files, ", ")))
	}

	if additionalContext != "" {
		prompt.WriteString(fmt.Sprintf("Context: %s\n", additionalContext))
	}

	prompt.WriteString("\nGit diff:\n")
	prompt.WriteString(diff)
	prompt.WriteString("\n\nRequirements:")
	prompt.WriteString("\n- Use conventional commit format if applicable")
	prompt.WriteString("\n- Be specific about what changed")
	prompt.WriteString("\n- Keep it under 50 characters for the subject line")
	prompt.WriteString("\n- Focus on the 'why' and 'what', not the 'how'")
	prompt.WriteString("\n- Do not include file names unless essential")

	return prompt.String()
}

// AddInteraction adds an interaction to the current session
func (sm *SessionManager) AddInteraction(interactionType, prompt, response, feedback string) {
	interaction := interfaces.AIInteraction{
		Type:      interactionType,
		Prompt:    prompt,
		Response:  response,
		Feedback:  feedback,
		Timestamp: time.Now(),
	}
	sm.currentSession.History = append(sm.currentSession.History, interaction)
}

// GetContextualPrompt builds a prompt with session context
func (sm *SessionManager) GetContextualPrompt(basePrompt string) string {
	if len(sm.currentSession.History) == 0 {
		return basePrompt
	}

	var contextBuilder strings.Builder
	contextBuilder.WriteString("Previous interactions context:\n")

	// Include last few interactions for context
	start := 0
	if len(sm.currentSession.History) > 3 {
		start = len(sm.currentSession.History) - 3
	}

	for i := start; i < len(sm.currentSession.History); i++ {
		interaction := sm.currentSession.History[i]
		contextBuilder.WriteString(fmt.Sprintf("- %s: %s\n",
			interaction.Type,
			truncateString(interaction.Response, 100)))
	}

	contextBuilder.WriteString("\nCurrent request:\n")
	contextBuilder.WriteString(basePrompt)

	return contextBuilder.String()
}

// callGeminiAPI makes a request to the Gemini API
func (sm *SessionManager) callGeminiAPI(prompt, apiKey string) (string, error) {
	request := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(geminiAPIURL+apiKey, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// AnalyzeCode analyzes code for various purposes (review, optimization, etc.)
func (sm *SessionManager) AnalyzeCode(code, analysisType string) (string, error) {
	if err := sm.validateConfig(); err != nil {
		return "", err
	}

	prompt := sm.buildAnalysisPrompt(code, analysisType)

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	response, err := sm.callGeminiAPI(prompt, apiKey)
	if err != nil {
		return "", fmt.Errorf("AI analysis failed: %v", err)
	}

	// Add interaction to session
	sm.AddInteraction("analysis", prompt, response, "")

	return response, nil
}

// buildAnalysisPrompt constructs the prompt for code analysis
func (sm *SessionManager) buildAnalysisPrompt(code, analysisType string) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Perform %s analysis on the following code changes.\n\n", analysisType))

	// Add specific instruction for text-only review
	prompt.WriteString("IMPORTANT: Provide ONLY text-based analysis and recommendations. ")
	prompt.WriteString("Do NOT include any code snippets, code blocks, or code examples in your response. ")
	prompt.WriteString("Focus on descriptive explanations, recommendations, and actionable insights.\n\n")

	prompt.WriteString("Code changes to analyze:\n")
	prompt.WriteString(code)
	prompt.WriteString("\n\nPlease provide a comprehensive text-based review covering:")

	switch analysisType {
	case "security":
		prompt.WriteString("\n- Security vulnerabilities and potential risks identified")
		prompt.WriteString("\n- Authentication and authorization concerns")
		prompt.WriteString("\n- Data protection and privacy considerations")
		prompt.WriteString("\n- Input validation and sanitization recommendations")
		prompt.WriteString("\n- Injection attack prevention measures")
		prompt.WriteString("\n- Cryptographic and secrets management improvements")
		prompt.WriteString("\n- Specific actionable security recommendations")
	case "performance":
		prompt.WriteString("\n- Performance bottlenecks and inefficiencies detected")
		prompt.WriteString("\n- Algorithm complexity and optimization opportunities")
		prompt.WriteString("\n- Memory usage patterns and potential improvements")
		prompt.WriteString("\n- Database query optimization suggestions")
		prompt.WriteString("\n- Network and I/O operation improvements")
		prompt.WriteString("\n- Caching strategies and resource utilization")
		prompt.WriteString("\n- Scalability considerations and recommendations")
	case "style":
		prompt.WriteString("\n- Code style and formatting consistency issues")
		prompt.WriteString("\n- Naming convention improvements")
		prompt.WriteString("\n- Code organization and structure suggestions")
		prompt.WriteString("\n- Documentation and comment quality assessment")
		prompt.WriteString("\n- Language-specific best practices recommendations")
		prompt.WriteString("\n- Readability and maintainability improvements")
	case "structure":
		prompt.WriteString("\n- Architectural design and modularity assessment")
		prompt.WriteString("\n- Design pattern usage and recommendations")
		prompt.WriteString("\n- Separation of concerns evaluation")
		prompt.WriteString("\n- Dependencies and coupling analysis")
		prompt.WriteString("\n- Error handling and exception management")
		prompt.WriteString("\n- Code maintainability and extensibility suggestions")
		prompt.WriteString("\n- Overall structural improvements")
	default:
		prompt.WriteString("\n- Overall code quality assessment")
		prompt.WriteString("\n- Potential bugs and issues identified")
		prompt.WriteString("\n- Best practices and improvement recommendations")
		prompt.WriteString("\n- Maintainability and reliability considerations")
	}

	prompt.WriteString("\n\nFormat your response as:")
	prompt.WriteString("\n1. Summary: Brief overview of findings")
	prompt.WriteString("\n2. Issues: List specific problems found (if any)")
	prompt.WriteString("\n3. Recommendations: Actionable improvement suggestions")
	prompt.WriteString("\n4. Priority: Indicate which items should be addressed first")
	prompt.WriteString("\n\nRemember: Use descriptive text only, no code snippets or examples.")

	return prompt.String()
}

// validateConfig ensures the configuration is available and valid
func (sm *SessionManager) validateConfig() error {
	if sm.config == nil {
		return fmt.Errorf("configuration manager is not set")
	}

	aiConfig := sm.config.GetAI()
	if aiConfig.Model == "" {
		return fmt.Errorf("AI model is not configured")
	}

	return nil
}

// GetSession returns the current session (for advanced usage)
func (sm *SessionManager) GetSession() *Session {
	return sm.currentSession
}

// ResetSession creates a new session, clearing history
func (sm *SessionManager) ResetSession() {
	sm.currentSession = &Session{
		History: make([]interfaces.AIInteraction, 0),
	}
}

// GetInteractionHistory returns the interaction history
func (sm *SessionManager) GetInteractionHistory() []interfaces.AIInteraction {
	return sm.currentSession.History
}

// SetContextualFeedback adds feedback to the last interaction
func (sm *SessionManager) SetContextualFeedback(feedback string) error {
	if len(sm.currentSession.History) == 0 {
		return fmt.Errorf("no interactions to provide feedback for")
	}

	lastIndex := len(sm.currentSession.History) - 1
	sm.currentSession.History[lastIndex].Feedback = feedback
	return nil
}

// truncateString truncates a string to a specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
