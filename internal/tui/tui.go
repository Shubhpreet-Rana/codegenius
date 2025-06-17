package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/codegenius/cli/internal/interfaces"
)

// Styles for the TUI
var (
	// Base styles
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true).
			Margin(1, 0)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF00")).
			Bold(true)

	// Layout styles
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Margin(1, 0)

	contentStyle = lipgloss.NewStyle().
			Margin(0, 2)
)

// TUI represents the terminal user interface
type TUI struct {
	service *interfaces.Service
}

// NewTUI creates a new terminal user interface
func NewTUI(service *interfaces.Service) *TUI {
	return &TUI{
		service: service,
	}
}

// MainMenu displays the main interactive menu
func (t *TUI) MainMenu() error {
	var action string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🤖 CodeGenius - AI-Powered Git Assistant").
				Description("Choose an action to perform").
				Options(
					huh.NewOption("Generate commit message", "commit"),
					huh.NewOption("Perform code review", "review"),
					huh.NewOption("View work history", "history"),
					huh.NewOption("Show statistics", "stats"),
					huh.NewOption("Initialize configuration", "init"),
					huh.NewOption("Exit", "exit"),
				).
				Value(&action),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return fmt.Errorf("menu selection failed: %v", err)
	}

	return t.handleAction(action)
}

// handleAction processes the selected action
func (t *TUI) handleAction(action string) error {
	switch action {
	case "commit":
		return t.handleCommit()
	case "review":
		return t.handleReview()
	case "history":
		return t.handleHistory()
	case "stats":
		return t.handleStats()
	case "init":
		return t.handleInit()
	case "exit":
		fmt.Println(successStyle.Render("👋 Goodbye! Happy coding!"))
		return nil
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

// handleCommit handles commit message generation with additional context
func (t *TUI) handleCommit() error {
	// Check for staged changes first
	hasStaged, err := t.service.Git.HasStagedChanges()
	if err != nil {
		return fmt.Errorf("error checking staged changes: %v", err)
	}

	if !hasStaged {
		fmt.Println(warningStyle.Render("⚠️  No staged changes detected"))
		fmt.Println(infoStyle.Render("Please stage your changes first with 'git add'"))
		return nil
	}

	// Get additional context from user
	var additionalContext string
	var includeContext bool

	contextForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("📝 Additional Context").
				Description("Would you like to provide additional context for the commit message?").
				Value(&includeContext),
		),
	).WithTheme(huh.ThemeCharm())

	if err := contextForm.Run(); err != nil {
		return fmt.Errorf("context form failed: %v", err)
	}

	if includeContext {
		contextInputForm := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("💭 Additional Context").
					Description("Provide any additional context or details about these changes").
					Placeholder("e.g., Fixed bug with user authentication, Added new feature for...").
					CharLimit(500).
					Value(&additionalContext),
			),
		).WithTheme(huh.ThemeCharm())

		if err := contextInputForm.Run(); err != nil {
			return fmt.Errorf("context input failed: %v", err)
		}
	}

	// Get git information
	diff, err := t.service.Git.GetDiff()
	if err != nil {
		return fmt.Errorf("error getting git diff: %v", err)
	}

	files, err := t.service.Git.GetChangedFiles()
	if err != nil {
		return fmt.Errorf("error getting changed files: %v", err)
	}

	branchName, err := t.service.Git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error getting current branch: %v", err)
	}

	// Display loading message
	fmt.Println(headerStyle.Render("🧠 Generating commit message..."))

	// Generate commit message
	message, err := t.service.AI.GenerateCommitMessage(diff, files, branchName, additionalContext)
	if err != nil {
		return fmt.Errorf("error generating commit message: %v", err)
	}

	// Display the generated message
	fmt.Println(containerStyle.Render(
		titleStyle.Render("Generated Commit Message") + "\n\n" +
			contentStyle.Render(message),
	))

	// Ask for user action
	var action string
	actionForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🚀 What would you like to do?").
				Options(
					huh.NewOption("✅ Commit with this message", "commit"),
					huh.NewOption("✏️  Edit the message", "edit"),
					huh.NewOption("🔄 Regenerate message", "regenerate"),
					huh.NewOption("❌ Cancel", "cancel"),
				).
				Value(&action),
		),
	).WithTheme(huh.ThemeCharm())

	if err := actionForm.Run(); err != nil {
		return fmt.Errorf("action form failed: %v", err)
	}

	switch action {
	case "commit":
		if err := t.service.Git.CommitWithMessage(message); err != nil {
			return fmt.Errorf("error committing: %v", err)
		}
		if err := t.service.History.Load(); err == nil {
			t.service.History.AddEntry(message)
		}
		fmt.Println(successStyle.Render("✅ Changes committed successfully!"))

	case "edit":
		editedMessage, err := t.service.Git.EditCommitMessage(message)
		if err != nil {
			return fmt.Errorf("error editing commit message: %v", err)
		}
		if err := t.service.Git.CommitWithMessage(editedMessage); err != nil {
			return fmt.Errorf("error committing: %v", err)
		}
		if err := t.service.History.Load(); err == nil {
			t.service.History.AddEntry(editedMessage)
		}
		fmt.Println(successStyle.Render("✅ Changes committed successfully!"))

	case "regenerate":
		return t.handleCommit() // Recursively regenerate

	case "cancel":
		fmt.Println(infoStyle.Render("🚫 Commit cancelled"))
	}

	return nil
}

// handleReview handles code review with multi-select options
func (t *TUI) handleReview() error {
	// Check for changes
	diff, err := t.service.Git.GetDiff()
	if err != nil {
		return fmt.Errorf("failed to get git diff: %v", err)
	}

	if diff == "" {
		fmt.Println(warningStyle.Render("⚠️  No changes detected for review"))
		return nil
	}

	// Get additional context for review
	var additionalContext string
	var includeContext bool

	contextForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("📝 Review Context").
				Description("Would you like to provide additional context for the code review?").
				Value(&includeContext),
		),
	).WithTheme(huh.ThemeCharm())

	if err := contextForm.Run(); err != nil {
		return fmt.Errorf("context form failed: %v", err)
	}

	if includeContext {
		contextInputForm := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("🎯 Review Focus").
					Description("Specify what aspects you'd like the review to focus on").
					Placeholder("e.g., Performance optimization, Security concerns, Code style...").
					CharLimit(300).
					Value(&additionalContext),
			),
		).WithTheme(huh.ThemeCharm())

		if err := contextInputForm.Run(); err != nil {
			return fmt.Errorf("context input failed: %v", err)
		}
	}

	// Multi-select review types
	supportedTypes := t.service.Review.GetSupportedTypes()
	var selectedTypes []string

	// Create options for multi-select
	options := make([]huh.Option[string], len(supportedTypes))
	for i, reviewType := range supportedTypes {
		emoji := getReviewTypeEmoji(reviewType)
		options[i] = huh.NewOption(fmt.Sprintf("%s %s", emoji, strings.Title(reviewType)), reviewType)
	}

	reviewForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("🔍 Select Review Types").
				Description("Choose which types of review to perform (multiple selections allowed)").
				Options(options...).
				Value(&selectedTypes),
		),
	).WithTheme(huh.ThemeCharm())

	if err := reviewForm.Run(); err != nil {
		return fmt.Errorf("review type selection failed: %v", err)
	}

	if len(selectedTypes) == 0 {
		fmt.Println(warningStyle.Render("⚠️  No review types selected"))
		return nil
	}

	// Perform reviews
	fmt.Println(headerStyle.Render("🔍 Performing code review..."))

	for _, reviewType := range selectedTypes {
		fmt.Printf("\n%s Analyzing %s...\n", getReviewTypeEmoji(reviewType), reviewType)

		// Add additional context to the diff if provided
		reviewDiff := diff
		if additionalContext != "" {
			reviewDiff = fmt.Sprintf("Additional Context: %s\n\n%s", additionalContext, diff)
		}

		review, err := t.service.Review.PerformReview(reviewDiff, reviewType)
		if err != nil {
			fmt.Printf("%s %s review failed: %v\n", errorStyle.Render("❌"), reviewType, err)
			continue
		}

		t.displayReviewResults(review)
	}

	return nil
}

// handleHistory handles work history display
func (t *TUI) handleHistory() error {
	var monthYear string
	var showAll bool

	historyForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("📊 Work History").
				Description("Show all history or filter by month/year?").
				Affirmative("Show All").
				Negative("Filter by Month").
				Value(&showAll),
		),
	).WithTheme(huh.ThemeCharm())

	if err := historyForm.Run(); err != nil {
		return fmt.Errorf("history form failed: %v", err)
	}

	if !showAll {
		filterForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("📅 Month/Year Filter").
					Description("Enter month and year (e.g., 'Dec 2024', 'Jan 2025')").
					Placeholder("Dec 2024").
					Value(&monthYear),
			),
		).WithTheme(huh.ThemeCharm())

		if err := filterForm.Run(); err != nil {
			return fmt.Errorf("filter form failed: %v", err)
		}
	}

	if err := t.service.History.Load(); err != nil {
		return fmt.Errorf("failed to load work history: %v", err)
	}

	return t.service.History.Display(monthYear)
}

// handleStats displays statistics
func (t *TUI) handleStats() error {
	if err := t.service.History.Load(); err != nil {
		return fmt.Errorf("failed to load work history: %v", err)
	}

	stats := t.service.History.GetStats()

	fmt.Println(headerStyle.Render("📊 CodeGenius Statistics"))
	fmt.Println(strings.Repeat("═", 50))

	if totalCommits, ok := stats["total_commits"].(int); ok {
		fmt.Printf("📈 Total commits: %s\n", successStyle.Render(fmt.Sprintf("%d", totalCommits)))
	}

	if mostActive, ok := stats["most_active_month"].(string); ok && mostActive != "" {
		fmt.Printf("🏆 Most active month: %s\n", successStyle.Render(mostActive))
	}

	if breakdown, ok := stats["monthly_breakdown"].(map[string]int); ok && len(breakdown) > 0 {
		fmt.Println(headerStyle.Render("\n📅 Monthly Breakdown:"))
		for month, count := range breakdown {
			fmt.Printf("  %s: %s commits\n", month, successStyle.Render(fmt.Sprintf("%d", count)))
		}
	}

	fmt.Println()
	return nil
}

// handleInit handles configuration initialization
func (t *TUI) handleInit() error {
	var confirm bool

	initForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("⚙️  Initialize Configuration").
				Description("This will create a new configuration file. Continue?").
				Value(&confirm),
		),
	).WithTheme(huh.ThemeCharm())

	if err := initForm.Run(); err != nil {
		return fmt.Errorf("init form failed: %v", err)
	}

	if !confirm {
		fmt.Println(infoStyle.Render("🚫 Configuration initialization cancelled"))
		return nil
	}

	if err := t.service.Config.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize configuration: %v", err)
	}

	fmt.Println(successStyle.Render("✅ Configuration initialized successfully!"))
	return nil
}

// displayReviewResults shows the review results with beautiful formatting
func (t *TUI) displayReviewResults(review *interfaces.ReviewResult) {
	fmt.Print(containerStyle.Render(
		titleStyle.Render("📋 Code Review Results") + "\n\n",
	))

	// Clean the review content to remove any code snippets
	cleanedReview := t.cleanReviewContent(review)

	// Display review type
	emoji := getReviewTypeEmoji(review.Type)
	title := fmt.Sprintf("%s %s Review", emoji, strings.Title(review.Type))

	fmt.Print(containerStyle.Render(
		titleStyle.Render(title) + "\n" +
			strings.Repeat("─", 50) + "\n\n",
	))

	// Display issues if any
	if len(cleanedReview.Issues) > 0 {
		fmt.Print(containerStyle.Render(
			errorStyle.Render(fmt.Sprintf("❌ Issues Found: %d", len(cleanedReview.Issues))) + "\n",
		))
		for i, issue := range cleanedReview.Issues {
			cleanMessage := t.cleanText(issue.Message)
			fmt.Printf("  %d. %s\n", i+1, cleanMessage)
			if issue.Line > 0 {
				fmt.Printf("     📍 Line: %d\n", issue.Line)
			}
			if issue.File != "" {
				fmt.Printf("     📁 File: %s\n", issue.File)
			}
		}
		fmt.Println()
	}

	// Display suggestions if any
	if len(cleanedReview.Suggestions) > 0 {
		fmt.Print(containerStyle.Render(
			warningStyle.Render(fmt.Sprintf("💡 Suggestions: %d", len(cleanedReview.Suggestions))) + "\n",
		))
		for i, suggestion := range cleanedReview.Suggestions {
			cleanMessage := t.cleanText(suggestion.Message)
			fmt.Printf("  %d. %s\n", i+1, cleanMessage)
			if suggestion.Line > 0 {
				fmt.Printf("     📍 Line: %d\n", suggestion.Line)
			}
			if suggestion.File != "" {
				fmt.Printf("     📁 File: %s\n", suggestion.File)
			}
		}
		fmt.Println()
	}

	// Display cleaned summary
	if cleanedReview.Summary != "" {
		cleanSummary := t.cleanText(cleanedReview.Summary)
		fmt.Print(containerStyle.Render(
			infoStyle.Render("📝 Summary:") + "\n" +
				cleanSummary + "\n\n",
		))
	}

	if len(cleanedReview.Issues) == 0 && len(cleanedReview.Suggestions) == 0 {
		fmt.Print(containerStyle.Render(
			successStyle.Render("✅ No specific issues found") + "\n\n",
		))
	}

	// Display final summary
	summary := fmt.Sprintf("📊 Review Summary: %d issues, %d suggestions",
		len(cleanedReview.Issues), len(cleanedReview.Suggestions))
	fmt.Print(containerStyle.Render(summary) + "\n")
}

// cleanReviewContent removes code snippets from review content
func (t *TUI) cleanReviewContent(review *interfaces.ReviewResult) *interfaces.ReviewResult {
	cleaned := &interfaces.ReviewResult{
		Type:        review.Type,
		Issues:      make([]interfaces.ReviewItem, 0),
		Suggestions: make([]interfaces.ReviewItem, 0),
		Summary:     t.cleanText(review.Summary),
	}

	// Clean issues
	for _, issue := range review.Issues {
		cleanedIssue := interfaces.ReviewItem{
			Message:    t.cleanText(issue.Message),
			Severity:   issue.Severity,
			File:       issue.File,
			Line:       issue.Line,
			Category:   issue.Category,
			Suggestion: t.cleanText(issue.Suggestion),
		}
		cleaned.Issues = append(cleaned.Issues, cleanedIssue)
	}

	// Clean suggestions
	for _, suggestion := range review.Suggestions {
		cleanedSuggestion := interfaces.ReviewItem{
			Message:    t.cleanText(suggestion.Message),
			Severity:   suggestion.Severity,
			File:       suggestion.File,
			Line:       suggestion.Line,
			Category:   suggestion.Category,
			Suggestion: t.cleanText(suggestion.Suggestion),
		}
		cleaned.Suggestions = append(cleaned.Suggestions, cleanedSuggestion)
	}

	return cleaned
}

// cleanText removes code snippets and code blocks from text
func (t *TUI) cleanText(text string) string {
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	inCodeBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip code blocks entirely
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			continue
		}

		// Skip lines that look like code
		if !t.isCodeLine(line) {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}

// isCodeLine detects if a line contains code
func (t *TUI) isCodeLine(line string) bool {
	// Common patterns that indicate code snippets
	codePatterns := []string{
		"func ", "class ", "def ", "import ", "from ", "package ", "use ",
		"if (", "for (", "while (", "switch (", "try {", "catch (",
		"=>", "->", "::", "&&", "||", "!=", "===", "!==",
	}

	lowerLine := strings.ToLower(line)

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

// getReviewTypeEmoji returns emoji for review types
func getReviewTypeEmoji(reviewType string) string {
	switch reviewType {
	case "security":
		return "🔒"
	case "performance":
		return "⚡"
	case "style":
		return "🎨"
	case "structure":
		return "🏗️"
	case "test":
		return "🧪"
	default:
		return "🔍"
	}
}
