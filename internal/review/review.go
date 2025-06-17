package review

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codegenius/cli/internal/interfaces"
)

// Reviewer implements the CodeReviewer interface
type Reviewer struct {
	config    interfaces.ConfigManager
	aiSession interfaces.AIProvider
}

// NewReviewer creates a new code reviewer
func NewReviewer(config interfaces.ConfigManager, ai interfaces.AIProvider) interfaces.CodeReviewer {
	return &Reviewer{
		config:    config,
		aiSession: ai,
	}
}

// HandleInteractive performs an interactive code review
func (r *Reviewer) HandleInteractive(diff string) error {
	if diff == "" {
		fmt.Println("No changes detected for review.")
		return nil
	}

	if err := r.validateDependencies(); err != nil {
		return fmt.Errorf("review setup error: %v", err)
	}

	supportedTypes := r.GetSupportedTypes()
	fmt.Println("🔍 Starting Code Review...")
	fmt.Println("Available review types:")
	for i, reviewType := range supportedTypes {
		fmt.Printf("  %d. %s\n", i+1, reviewType)
	}
	fmt.Print("Select review type (number or 'all'): ")

	var input string
	fmt.Scanln(&input)

	if input == "all" {
		return r.performAllReviews(diff)
	}

	// Handle single review type selection
	var selectedType string
	switch input {
	case "1":
		selectedType = "security"
	case "2":
		selectedType = "performance"
	case "3":
		selectedType = "style"
	case "4":
		selectedType = "structure"
	default:
		selectedType = input // Allow direct type input
	}

	// Validate the selected type
	if !r.isValidReviewType(selectedType) {
		return fmt.Errorf("invalid review type: %s", selectedType)
	}

	review, err := r.PerformReview(diff, selectedType)
	if err != nil {
		return fmt.Errorf("review failed: %v", err)
	}

	r.DisplayResults(review)
	return nil
}

// PerformReview performs a specific type of code review
func (r *Reviewer) PerformReview(diff, reviewType string) (*interfaces.ReviewResult, error) {
	if err := r.validateDependencies(); err != nil {
		return nil, fmt.Errorf("review setup error: %v", err)
	}

	if !r.isValidReviewType(reviewType) {
		return nil, fmt.Errorf("invalid review type: %s", reviewType)
	}

	if strings.TrimSpace(diff) == "" {
		return &interfaces.ReviewResult{
			Type:        reviewType,
			Issues:      []interfaces.ReviewItem{},
			Suggestions: []interfaces.ReviewItem{},
			Summary:     "No changes to review.",
		}, nil
	}

	response, err := r.aiSession.AnalyzeCode(diff, reviewType)
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %v", err)
	}

	review := r.parseReviewResponse(response, reviewType)
	return review, nil
}

// GetSupportedTypes returns the list of supported review types
func (r *Reviewer) GetSupportedTypes() []string {
	if r.config == nil {
		// Return default types if config is not available
		return []string{"security", "performance", "style", "structure"}
	}

	reviewConfig := r.config.GetReview()
	if len(reviewConfig.EnabledTypes) == 0 {
		// Return default types if none are configured
		return []string{"security", "performance", "style", "structure"}
	}

	return reviewConfig.EnabledTypes
}

// performAllReviews performs all enabled review types
func (r *Reviewer) performAllReviews(diff string) error {
	supportedTypes := r.GetSupportedTypes()

	for _, reviewType := range supportedTypes {
		fmt.Printf("\n🔍 Performing %s review...\n", reviewType)

		review, err := r.PerformReview(diff, reviewType)
		if err != nil {
			fmt.Printf("❌ %s review failed: %v\n", reviewType, err)
			continue
		}

		r.DisplayResults(review)
	}
	return nil
}

// isValidReviewType checks if the review type is supported
func (r *Reviewer) isValidReviewType(reviewType string) bool {
	supportedTypes := r.GetSupportedTypes()
	for _, supportedType := range supportedTypes {
		if supportedType == reviewType {
			return true
		}
	}
	return false
}

// validateDependencies ensures required dependencies are available
func (r *Reviewer) validateDependencies() error {
	if r.config == nil {
		return fmt.Errorf("configuration manager is not set")
	}
	if r.aiSession == nil {
		return fmt.Errorf("AI provider is not set")
	}
	return nil
}

// parseReviewResponse parses the AI response into a structured review
func (r *Reviewer) parseReviewResponse(response, reviewType string) *interfaces.ReviewResult {
	review := &interfaces.ReviewResult{
		Type:        reviewType,
		Issues:      make([]interfaces.ReviewItem, 0),
		Suggestions: make([]interfaces.ReviewItem, 0),
		Summary:     r.cleanResponseText(response),
	}

	// Try to extract structured information from response
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip code snippets and code blocks
		if r.isCodeSnippet(line) {
			continue
		}

		// Look for issue patterns
		if r.isIssueLine(line) {
			item := r.parseReviewItem(line, "issue")
			if item != nil {
				review.Issues = append(review.Issues, *item)
			}
		}

		// Look for suggestion patterns
		if r.isSuggestionLine(line) {
			item := r.parseReviewItem(line, "suggestion")
			if item != nil {
				review.Suggestions = append(review.Suggestions, *item)
			}
		}
	}

	return review
}

// isCodeSnippet detects if a line contains code snippets
func (r *Reviewer) isCodeSnippet(line string) bool {
	// Common patterns that indicate code snippets
	codePatterns := []string{
		"```", "```go", "```javascript", "```python", "```java", "```rust", "```c++",
		"func ", "class ", "def ", "import ", "from ", "package ", "use ",
		"if (", "for (", "while (", "switch (", "try {", "catch (",
		"=>", "->", "::", "&&", "||", "!=", "===", "!==",
	}

	lowerLine := strings.ToLower(line)
	trimmedLine := strings.TrimSpace(line)

	// Skip code blocks
	if strings.HasPrefix(trimmedLine, "```") {
		return true
	}

	// Skip lines that look like code
	for _, pattern := range codePatterns {
		if strings.Contains(lowerLine, pattern) {
			return true
		}
	}

	// Skip lines with typical code syntax patterns
	if strings.Contains(line, "{") && strings.Contains(line, "}") {
		return true
	}

	// Skip indented lines that look like code (4+ spaces or tabs)
	if len(line) > 0 && (strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t")) {
		// Check if it contains code-like patterns
		if strings.Contains(line, "(") && strings.Contains(line, ")") {
			return true
		}
		if strings.Contains(line, "=") && (strings.Contains(line, ";") || strings.Contains(line, "{")) {
			return true
		}
	}

	return false
}

// cleanResponseText removes code snippets and cleans up the response text
func (r *Reviewer) cleanResponseText(text string) string {
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	inCodeBlock := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Skip code blocks entirely
		if strings.HasPrefix(trimmedLine, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			continue
		}

		// Skip lines that look like code
		if !r.isCodeSnippet(line) {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}

// isIssueLine checks if a line represents an issue
func (r *Reviewer) isIssueLine(line string) bool {
	issuePatterns := []string{
		"issue", "problem", "error", "bug", "vulnerability", "risk", "concern",
		"dangerous", "unsafe", "insecure", "critical", "warning",
	}

	lowerLine := strings.ToLower(line)
	for _, pattern := range issuePatterns {
		if strings.Contains(lowerLine, pattern) {
			return true
		}
	}
	return false
}

// isSuggestionLine checks if a line represents a suggestion
func (r *Reviewer) isSuggestionLine(line string) bool {
	suggestionPatterns := []string{
		"suggest", "recommend", "consider", "improve", "optimize", "better",
		"should", "could", "might", "enhancement", "refactor",
	}

	lowerLine := strings.ToLower(line)
	for _, pattern := range suggestionPatterns {
		if strings.Contains(lowerLine, pattern) {
			return true
		}
	}
	return false
}

// parseReviewItem extracts a review item from a line
func (r *Reviewer) parseReviewItem(line, itemType string) *interfaces.ReviewItem {
	// Basic parsing - can be enhanced with more sophisticated parsing
	item := &interfaces.ReviewItem{
		Category: itemType,
		Message:  line,
		Severity: r.extractSeverity(line),
	}

	// Try to extract line number
	lineNumRegex := regexp.MustCompile(`line\s+(\d+)`)
	if matches := lineNumRegex.FindStringSubmatch(line); len(matches) > 1 {
		fmt.Sscanf(matches[1], "%d", &item.Line)
	}

	// Try to extract file name
	fileRegex := regexp.MustCompile(`(?:in|file)\s+([a-zA-Z0-9_\-./]+\.[a-zA-Z0-9]+)`)
	if matches := fileRegex.FindStringSubmatch(line); len(matches) > 1 {
		item.File = matches[1]
	}

	return item
}

// extractSeverity attempts to extract severity from the text
func (r *Reviewer) extractSeverity(text string) string {
	lowerText := strings.ToLower(text)

	if strings.Contains(lowerText, "critical") || strings.Contains(lowerText, "severe") {
		return "critical"
	}
	if strings.Contains(lowerText, "high") || strings.Contains(lowerText, "major") {
		return "high"
	}
	if strings.Contains(lowerText, "medium") || strings.Contains(lowerText, "moderate") {
		return "medium"
	}
	if strings.Contains(lowerText, "low") || strings.Contains(lowerText, "minor") {
		return "low"
	}

	return "medium" // default
}

// DisplayResults shows the review results in a formatted way
func (r *Reviewer) DisplayResults(review *interfaces.ReviewResult) {
	if review == nil {
		fmt.Println("No review results to display.")
		return
	}

	fmt.Printf("\n📋 %s Review Results\n", strings.Title(review.Type))
	fmt.Println(strings.Repeat("=", 50))

	if len(review.Issues) > 0 {
		fmt.Printf("\n❌ Issues Found (%d):\n", len(review.Issues))
		for i, issue := range review.Issues {
			cleanMessage := r.cleanResponseText(issue.Message)
			fmt.Printf("  %d. [%s] %s\n", i+1, strings.ToUpper(issue.Severity), cleanMessage)
			if issue.Line > 0 {
				fmt.Printf("     📍 Line: %d\n", issue.Line)
			}
			if issue.File != "" {
				fmt.Printf("     📁 File: %s\n", issue.File)
			}
		}
	}

	if len(review.Suggestions) > 0 {
		fmt.Printf("\n💡 Suggestions (%d):\n", len(review.Suggestions))
		for i, suggestion := range review.Suggestions {
			cleanMessage := r.cleanResponseText(suggestion.Message)
			fmt.Printf("  %d. %s\n", i+1, cleanMessage)
			if suggestion.Line > 0 {
				fmt.Printf("     📍 Line: %d\n", suggestion.Line)
			}
			if suggestion.File != "" {
				fmt.Printf("     📁 File: %s\n", suggestion.File)
			}
		}
	}

	if len(review.Issues) == 0 && len(review.Suggestions) == 0 {
		fmt.Println("\n✅ No specific issues or suggestions found.")
	}

	// Display cleaned summary
	cleanSummary := r.cleanResponseText(review.Summary)
	fmt.Printf("\n📝 Summary:\n%s\n", cleanSummary)
}

// GetReviewStats returns statistics about the review
func (r *Reviewer) GetReviewStats(review *interfaces.ReviewResult) map[string]interface{} {
	if review == nil {
		return map[string]interface{}{}
	}

	stats := make(map[string]interface{})
	stats["type"] = review.Type
	stats["total_issues"] = len(review.Issues)
	stats["total_suggestions"] = len(review.Suggestions)

	// Count by severity
	severityCounts := make(map[string]int)
	for _, issue := range review.Issues {
		severityCounts[issue.Severity]++
	}
	stats["severity_breakdown"] = severityCounts

	// Count by category
	categoryCounts := make(map[string]int)
	for _, issue := range review.Issues {
		categoryCounts[issue.Category]++
	}
	for _, suggestion := range review.Suggestions {
		categoryCounts[suggestion.Category]++
	}
	stats["category_breakdown"] = categoryCounts

	return stats
}

// BatchReview performs multiple review types and returns all results
func (r *Reviewer) BatchReview(diff string, reviewTypes []string) (map[string]*interfaces.ReviewResult, error) {
	if err := r.validateDependencies(); err != nil {
		return nil, fmt.Errorf("review setup error: %v", err)
	}

	results := make(map[string]*interfaces.ReviewResult)

	for _, reviewType := range reviewTypes {
		if !r.isValidReviewType(reviewType) {
			continue // Skip invalid types
		}

		review, err := r.PerformReview(diff, reviewType)
		if err != nil {
			// Log the error but continue with other reviews
			fmt.Printf("Warning: %s review failed: %v\n", reviewType, err)
			continue
		}

		results[reviewType] = review
	}

	return results, nil
}
