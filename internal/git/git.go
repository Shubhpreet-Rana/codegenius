package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codegenius/cli/internal/interfaces"
)

// Repository implements the GitRepository interface
type Repository struct {
	workingDir string
}

// NewRepository creates a new Git repository instance
func NewRepository(workingDir string) interfaces.GitRepository {
	if workingDir == "" {
		workingDir = "."
	}
	return &Repository{workingDir: workingDir}
}

// New creates a new Git repository instance with current directory
func New() *Repository {
	return &Repository{workingDir: "."}
}

// GetDiff returns the git diff for staged changes
func (r *Repository) GetDiff() (string, error) {
	if err := r.validateGitRepo(); err != nil {
		return "", err
	}

	cmd := exec.Command("git", "diff", "--cached")
	cmd.Dir = r.workingDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %v", err)
	}
	return string(output), nil
}

// GetChangedFiles returns a list of files that have been changed
func (r *Repository) GetChangedFiles() ([]string, error) {
	if err := r.validateGitRepo(); err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	cmd.Dir = r.workingDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %v", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}
	return files, nil
}

// GetCurrentBranch returns the name of the current Git branch
func (r *Repository) GetCurrentBranch() (string, error) {
	if err := r.validateGitRepo(); err != nil {
		return "", err
	}

	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = r.workingDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetRecentCommits returns recent commit messages for context
func (r *Repository) GetRecentCommits() ([]string, error) {
	if err := r.validateGitRepo(); err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "log", "--oneline", "-10", "--pretty=format:%s")
	cmd.Dir = r.workingDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent commits: %v", err)
	}

	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(commits) == 1 && commits[0] == "" {
		return []string{}, nil
	}
	return commits, nil
}

// HasStagedChanges checks if there are any staged changes
func (r *Repository) HasStagedChanges() (bool, error) {
	if err := r.validateGitRepo(); err != nil {
		return false, err
	}

	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = r.workingDir
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return true, nil // Exit code 1 means there are differences
			}
		}
		return false, fmt.Errorf("error checking staged changes: %v", err)
	}
	return false, nil // Exit code 0 means no differences
}

// CommitWithMessage commits the staged changes with the provided message
func (r *Repository) CommitWithMessage(message string) error {
	if err := r.validateGitRepo(); err != nil {
		return err
	}

	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = r.workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit: %v", err)
	}
	return nil
}

// EditCommitMessage opens an editor for the user to edit the commit message
func (r *Repository) EditCommitMessage(message string) (string, error) {
	if err := r.validateGitRepo(); err != nil {
		return "", err
	}

	// Create a temporary file with the message
	tmpfile, err := os.CreateTemp("", "commit_message_*.txt")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write the message to the temporary file
	if _, err := tmpfile.WriteString(message); err != nil {
		tmpfile.Close()
		return "", fmt.Errorf("error writing to temporary file: %v", err)
	}
	tmpfile.Close()

	// Get the editor from environment or use nano as default
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	// Open the editor
	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Dir = r.workingDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running editor: %v", err)
	}

	// Read the edited message
	editedMessage, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return "", fmt.Errorf("error reading edited message: %v", err)
	}

	return strings.TrimSpace(string(editedMessage)), nil
}

// AnalyzeDiffContext analyzes git diff and extracts meaningful context
func (r *Repository) AnalyzeDiffContext(diff string, ignorePatterns []string) (string, []string) {
	lines := strings.Split(diff, "\n")
	var significantChanges []string
	var fileChanges []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip if should be ignored
		shouldIgnore := false
		for _, pattern := range ignorePatterns {
			if strings.Contains(line, pattern) {
				shouldIgnore = true
				break
			}
		}
		if shouldIgnore {
			continue
		}

		// Track file changes
		if strings.HasPrefix(line, "diff --git") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				file := strings.TrimPrefix(parts[3], "b/")
				fileChanges = append(fileChanges, file)
			}
		}

		// Analyze significant changes
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			if len(line) > 10 { // Ignore very short additions
				significantChanges = append(significantChanges, line[1:])
			}
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			if len(line) > 10 { // Ignore very short deletions
				significantChanges = append(significantChanges, "REMOVED: "+line[1:])
			}
		}
	}

	maxChanges := min(5, len(significantChanges))
	context := fmt.Sprintf("Modified files: %s\nKey changes: %s",
		strings.Join(fileChanges, ", "),
		strings.Join(significantChanges[:maxChanges], "; "))

	return context, fileChanges
}

// validateGitRepo checks if the current directory is a Git repository
func (r *Repository) validateGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = r.workingDir
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("not a git repository (or any of the parent directories)")
	}
	return nil
}

// GetWorkingDir returns the working directory
func (r *Repository) GetWorkingDir() string {
	return r.workingDir
}

// SetWorkingDir sets the working directory
func (r *Repository) SetWorkingDir(dir string) {
	r.workingDir = dir
}

// IsClean checks if the working directory is clean (no unstaged changes)
func (r *Repository) IsClean() (bool, error) {
	if err := r.validateGitRepo(); err != nil {
		return false, err
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = r.workingDir
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %v", err)
	}

	return strings.TrimSpace(string(output)) == "", nil
}

// GetStatus returns the git status output
func (r *Repository) GetStatus() (string, error) {
	if err := r.validateGitRepo(); err != nil {
		return "", err
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = r.workingDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git status: %v", err)
	}

	return string(output), nil
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
