package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/codegenius/cli/internal/ai"
	"github.com/codegenius/cli/internal/config"
	"github.com/codegenius/cli/internal/container"
	"github.com/codegenius/cli/internal/git"
	"github.com/codegenius/cli/internal/history"
	"github.com/codegenius/cli/internal/interfaces"
	"github.com/codegenius/cli/internal/review"
	"github.com/codegenius/cli/internal/tui"
)

func main() {
	var (
		reviewFlag      = flag.Bool("review", false, "Perform code review")
		historyFlag     = flag.String("history", "", "Display work history for month-year (e.g., 'Dec 2024')")
		initFlag        = flag.Bool("init", false, "Initialize configuration")
		interactiveFlag = flag.Bool("interactive", false, "Run in interactive mode")
		tuiFlag         = flag.Bool("tui", false, "Run with beautiful terminal UI (recommended)")
		helpFlag        = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	// Show help if requested
	if *helpFlag {
		showHelp()
		return
	}

	// Build service with dependency injection
	service, err := buildService()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// If TUI mode is requested, use the beautiful terminal interface
	if *tuiFlag {
		handleTUIMode(service)
		return
	}

	// Handle different modes (legacy CLI)
	switch {
	case *initFlag:
		handleInit(service)
	case *reviewFlag:
		handleCodeReview(service)
	case *historyFlag != "":
		handleHistory(service, *historyFlag)
	case *interactiveFlag:
		handleInteractive(service)
	default:
		// If no flags specified, suggest TUI mode
		if !hasAnyFlags() {
			fmt.Println("🤖 Welcome to CodeGenius!")
			fmt.Println("For the best experience, try the beautiful TUI mode:")
			fmt.Println("  codegenius --tui")
			fmt.Println()
			fmt.Println("Or use the classic interactive mode:")
			fmt.Println("  codegenius --interactive")
			fmt.Println()
			fmt.Println("For help: codegenius --help")
			return
		}
		handleAutoCommit(service)
	}
}

// hasAnyFlags checks if any command line flags were provided
func hasAnyFlags() bool {
	flagSet := false
	flag.Visit(func(*flag.Flag) {
		flagSet = true
	})
	return flagSet
}

// showHelp displays help information
func showHelp() {
	fmt.Println(`🤖 CodeGenius - AI-Powered Git Assistant

USAGE:
    codegenius [FLAGS]

FLAGS:
    --tui              Launch beautiful terminal UI (recommended)
    --interactive      Run in interactive mode (legacy)
    --review           Perform code review on staged changes
    --history [month]  Display work history (e.g., "Dec 2024")
    --init             Initialize configuration
    --help             Show this help message

EXAMPLES:
    codegenius --tui                    # Launch beautiful terminal interface
    codegenius                          # Generate commit message for staged changes
    codegenius --review                 # Review staged changes
    codegenius --history "Dec 2024"     # Show December 2024 history
    codegenius --init                   # Setup configuration

SETUP:
    1. Get your Gemini API key: https://makersuite.google.com/app/apikey
    2. Set environment variable: export GEMINI_API_KEY="your-key-here"
    3. Initialize configuration: codegenius --init
    4. Stage your changes: git add .
    5. Run: codegenius --tui

For more information, visit: https://github.com/codegenius/cli`)
}

// handleTUIMode runs the beautiful terminal user interface
func handleTUIMode(service *interfaces.Service) {
	terminalUI := tui.NewTUI(service)

	for {
		err := terminalUI.MainMenu()
		if err != nil {
			// Check if it's an exit condition
			if strings.Contains(err.Error(), "exit") ||
				strings.Contains(err.Error(), "interrupt") ||
				strings.Contains(err.Error(), "EOF") {
				fmt.Println("👋 Goodbye! Happy coding!")
				return
			}
			log.Printf("TUI error: %v", err)
			return
		}

		// Ask if user wants to continue
		fmt.Println()
		fmt.Print("Press Enter to return to main menu or Ctrl+C to exit...")
		fmt.Scanln()
	}
}

// buildService creates the application service with all dependencies
func buildService() (*interfaces.Service, error) {
	// Create configuration manager
	configManager := config.NewManager()
	if err := configManager.Load(); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	// Create Git repository
	gitRepo := git.NewRepository(".")

	// Create AI provider
	aiProvider := ai.NewSessionManager(configManager)

	// Create history manager
	historyManager := history.NewManager("")

	// Create code reviewer
	codeReviewer := review.NewReviewer(configManager, aiProvider)

	// Build service using dependency injection container
	service := container.NewContainer().
		WithConfig(configManager).
		WithGit(gitRepo).
		WithAI(aiProvider).
		WithHistory(historyManager).
		WithReview(codeReviewer).
		Build()

	return service, nil
}

func handleInit(service *interfaces.Service) {
	if err := service.Config.Initialize(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}
	fmt.Println("✅ Configuration initialized successfully!")
}

func handleCodeReview(service *interfaces.Service) {
	diff, err := service.Git.GetDiff()
	if err != nil {
		log.Fatalf("Failed to get git diff: %v", err)
	}

	if diff == "" {
		fmt.Println("⚠️  No changes detected for review.")
		return
	}

	if err := service.Review.HandleInteractive(diff); err != nil {
		log.Fatalf("Code review failed: %v", err)
	}
}

func handleHistory(service *interfaces.Service, monthYear string) {
	if err := service.History.Load(); err != nil {
		log.Fatalf("Failed to load work history: %v", err)
	}

	if err := service.History.Display(monthYear); err != nil {
		log.Fatalf("Failed to display history: %v", err)
	}
}

func handleInteractive(service *interfaces.Service) {
	fmt.Println("=== CodeGenius Interactive Mode (Legacy) ===")
	fmt.Println("💡 Tip: Try the new TUI mode with: codegenius --tui")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  commit   - Generate and apply commit message")
	fmt.Println("  review   - Perform code review")
	fmt.Println("  history  - View work history")
	fmt.Println("  stats    - Show statistics")
	fmt.Println("  exit     - Exit interactive mode")
	fmt.Println()

	for {
		fmt.Print("codegenius> ")
		var command string
		fmt.Scanln(&command)

		switch command {
		case "commit":
			if err := handleAutoCommit(service); err != nil {
				fmt.Printf("❌ Commit failed: %v\n", err)
			}
		case "review":
			handleCodeReview(service)
		case "history":
			fmt.Print("Enter month-year (e.g., 'Dec 2024') or press Enter for all: ")
			var monthYear string
			fmt.Scanln(&monthYear)
			handleHistory(service, monthYear)
		case "stats":
			handleStats(service)
		case "exit":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Try: commit, review, history, stats, exit")
		}
	}
}

func handleAutoCommit(service *interfaces.Service) error {
	// Check for staged changes
	hasStaged, err := service.Git.HasStagedChanges()
	if err != nil {
		return fmt.Errorf("error checking staged changes: %v", err)
	}

	if !hasStaged {
		fmt.Println("⚠️  No staged changes detected. Please stage your changes first with 'git add'.")
		return nil
	}

	// Get git information
	diff, err := service.Git.GetDiff()
	if err != nil {
		return fmt.Errorf("error getting git diff: %v", err)
	}

	files, err := service.Git.GetChangedFiles()
	if err != nil {
		return fmt.Errorf("error getting changed files: %v", err)
	}

	branchName, err := service.Git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error getting current branch: %v", err)
	}

	// Generate commit message with AI
	message, err := service.AI.GenerateCommitMessage(diff, files, branchName, "")
	if err != nil {
		return fmt.Errorf("error generating commit message: %v", err)
	}

	fmt.Printf("📝 Generated commit message:\n%s\n\n", message)
	fmt.Print("Use this commit message? (y/n/e for edit): ")

	var response string
	fmt.Scanln(&response)

	switch response {
	case "y", "Y", "yes", "Yes":
		if err := service.Git.CommitWithMessage(message); err != nil {
			return fmt.Errorf("error committing: %v", err)
		}

		// Load and add to work history
		if err := service.History.Load(); err != nil {
			log.Printf("Warning: Failed to load work history: %v", err)
		} else if err := service.History.AddEntry(message); err != nil {
			log.Printf("Warning: Failed to update work history: %v", err)
		}

		fmt.Println("✅ Changes committed successfully!")
	case "e", "E", "edit", "Edit":
		editedMessage, err := service.Git.EditCommitMessage(message)
		if err != nil {
			return fmt.Errorf("error editing commit message: %v", err)
		}

		if err := service.Git.CommitWithMessage(editedMessage); err != nil {
			return fmt.Errorf("error committing: %v", err)
		}

		// Load and add to work history
		if err := service.History.Load(); err != nil {
			log.Printf("Warning: Failed to load work history: %v", err)
		} else if err := service.History.AddEntry(editedMessage); err != nil {
			log.Printf("Warning: Failed to update work history: %v", err)
		}

		fmt.Println("✅ Changes committed successfully!")
	default:
		fmt.Println("🚫 Commit cancelled.")
	}

	return nil
}

func handleStats(service *interfaces.Service) {
	if err := service.History.Load(); err != nil {
		fmt.Printf("❌ Failed to load work history: %v\n", err)
		return
	}

	stats := service.History.GetStats()

	fmt.Println("\n📊 CodeGenius Statistics")
	fmt.Println(strings.Repeat("═", 30))

	if totalCommits, ok := stats["total_commits"].(int); ok {
		fmt.Printf("📈 Total commits: %d\n", totalCommits)
	}

	if mostActive, ok := stats["most_active_month"].(string); ok && mostActive != "" {
		fmt.Printf("🏆 Most active month: %s\n", mostActive)
	}

	if breakdown, ok := stats["monthly_breakdown"].(map[string]int); ok && len(breakdown) > 0 {
		fmt.Println("\n📅 Monthly breakdown:")
		for month, count := range breakdown {
			fmt.Printf("  %s: %d commits\n", month, count)
		}
	}

	fmt.Println()
}
