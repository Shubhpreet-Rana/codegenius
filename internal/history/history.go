package history

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Shubhpreet-Rana/codegenius/internal/interfaces"
)

const (
	workHistoryFile   = ".git/work_history.json"
	dateFormat        = "02 Jan 2006"
	displayDateFormat = "02 Jan 2006"
)

// WorkHistory represents the work history data structure
type WorkHistory struct {
	Entries []interfaces.HistoryEntry `json:"entries"`
}

// Manager implements the HistoryManager interface
type Manager struct {
	history  *WorkHistory
	filePath string
}

// NewManager creates a new history manager
func NewManager(filePath string) interfaces.HistoryManager {
	if filePath == "" {
		filePath = workHistoryFile
	}
	return &Manager{
		filePath: filePath,
	}
}

// Load reads the work history from file or creates a new one
func (m *Manager) Load() error {
	history := &WorkHistory{
		Entries: make([]interfaces.HistoryEntry, 0),
	}

	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		m.history = history
		return nil // Return empty history if file doesn't exist
	}

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return fmt.Errorf("error reading work history file: %v", err)
	}

	err = json.Unmarshal(data, history)
	if err != nil {
		return fmt.Errorf("error parsing work history: %v", err)
	}

	m.history = history
	return nil
}

// Save writes the work history to file
func (m *Manager) Save() error {
	if m.history == nil {
		return fmt.Errorf("no history data to save")
	}

	// Ensure the .git directory exists
	dir := ".git"
	if strings.Contains(m.filePath, "/") {
		dir = m.filePath[:strings.LastIndex(m.filePath, "/")]
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %v", dir, err)
	}

	data, err := json.MarshalIndent(m.history, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling work history: %v", err)
	}

	err = os.WriteFile(m.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing work history file: %v", err)
	}

	return nil
}

// AddEntry adds a new entry to the work history
func (m *Manager) AddEntry(message string) error {
	if m.history == nil {
		if err := m.Load(); err != nil {
			return fmt.Errorf("failed to load history before adding entry: %v", err)
		}
	}

	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	now := time.Now()
	entry := interfaces.HistoryEntry{
		Date:    now.Format(dateFormat),
		Summary: message,
	}

	m.history.Entries = append(m.history.Entries, entry)
	return m.Save()
}

// Display shows work history for a specific month/year
func (m *Manager) Display(monthYear string) error {
	if m.history == nil {
		if err := m.Load(); err != nil {
			return fmt.Errorf("failed to load history: %v", err)
		}
	}

	if monthYear == "" {
		// Show all entries
		return m.displayAll()
	}

	// Filter entries by month/year
	filtered := m.FilterByMonthYear(monthYear)
	if len(filtered) == 0 {
		fmt.Printf("No work history found for %s\n", monthYear)
		return nil
	}

	fmt.Printf("\n🗓️  Work History for %s\n", monthYear)
	fmt.Println(strings.Repeat("=", 50))

	// Group by date
	dateGroups := make(map[string][]interfaces.HistoryEntry)
	for _, entry := range filtered {
		dateGroups[entry.Date] = append(dateGroups[entry.Date], entry)
	}

	// Sort dates
	dates := make([]string, 0, len(dateGroups))
	for date := range dateGroups {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Display entries
	for _, date := range dates {
		entries := dateGroups[date]
		fmt.Printf("\n📅 %s\n", date)
		for i, entry := range entries {
			fmt.Printf("   %d. %s\n", i+1, entry.Summary)
		}
	}

	fmt.Printf("\nTotal commits: %d\n", len(filtered))
	return nil
}

// displayAll shows all work history entries
func (m *Manager) displayAll() error {
	if len(m.history.Entries) == 0 {
		fmt.Println("No work history found.")
		return nil
	}

	fmt.Println("\n🗓️  Complete Work History")
	fmt.Println(strings.Repeat("=", 50))

	// Group by month/year
	monthGroups := make(map[string][]interfaces.HistoryEntry)
	for _, entry := range m.history.Entries {
		monthYear := extractMonthYear(entry.Date)
		monthGroups[monthYear] = append(monthGroups[monthYear], entry)
	}

	// Sort month/years
	monthYears := make([]string, 0, len(monthGroups))
	for monthYear := range monthGroups {
		monthYears = append(monthYears, monthYear)
	}
	sort.Strings(monthYears)

	// Display by month
	for _, monthYear := range monthYears {
		entries := monthGroups[monthYear]
		fmt.Printf("\n📅 %s (%d commits)\n", monthYear, len(entries))

		// Show first few entries as preview
		maxShow := 5
		for i, entry := range entries {
			if i >= maxShow {
				fmt.Printf("   ... and %d more\n", len(entries)-maxShow)
				break
			}
			fmt.Printf("   • %s\n", truncateSummary(entry.Summary, 60))
		}
	}

	fmt.Printf("\nTotal commits: %d\n", len(m.history.Entries))
	return nil
}

// FilterByMonthYear filters entries by month/year string
func (m *Manager) FilterByMonthYear(monthYear string) []interfaces.HistoryEntry {
	if m.history == nil {
		return []interfaces.HistoryEntry{}
	}

	var filtered []interfaces.HistoryEntry
	for _, entry := range m.history.Entries {
		if strings.Contains(entry.Date, monthYear) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// GetStats returns statistics about the work history
func (m *Manager) GetStats() map[string]interface{} {
	if m.history == nil {
		return map[string]interface{}{
			"total_commits":     0,
			"monthly_breakdown": map[string]int{},
			"most_active_month": "",
		}
	}

	stats := make(map[string]interface{})

	stats["total_commits"] = len(m.history.Entries)

	// Group by month for monthly stats
	monthGroups := make(map[string]int)
	for _, entry := range m.history.Entries {
		monthYear := extractMonthYear(entry.Date)
		monthGroups[monthYear]++
	}

	stats["monthly_breakdown"] = monthGroups
	stats["most_active_month"] = findMostActiveMonth(monthGroups)

	return stats
}

// GetHistory returns the full history (for advanced usage)
func (m *Manager) GetHistory() *WorkHistory {
	return m.history
}

// GetEntries returns all history entries
func (m *Manager) GetEntries() []interfaces.HistoryEntry {
	if m.history == nil {
		return []interfaces.HistoryEntry{}
	}
	return m.history.Entries
}

// GetEntriesForDate returns entries for a specific date
func (m *Manager) GetEntriesForDate(date string) []interfaces.HistoryEntry {
	if m.history == nil {
		return []interfaces.HistoryEntry{}
	}

	var entries []interfaces.HistoryEntry
	for _, entry := range m.history.Entries {
		if entry.Date == date {
			entries = append(entries, entry)
		}
	}
	return entries
}

// GetEntriesInDateRange returns entries within a date range
func (m *Manager) GetEntriesInDateRange(startDate, endDate time.Time) []interfaces.HistoryEntry {
	if m.history == nil {
		return []interfaces.HistoryEntry{}
	}

	var entries []interfaces.HistoryEntry
	for _, entry := range m.history.Entries {
		entryDate, err := time.Parse(dateFormat, entry.Date)
		if err != nil {
			continue // Skip entries with invalid dates
		}

		if (entryDate.Equal(startDate) || entryDate.After(startDate)) &&
			(entryDate.Equal(endDate) || entryDate.Before(endDate)) {
			entries = append(entries, entry)
		}
	}
	return entries
}

// Clear removes all history entries
func (m *Manager) Clear() error {
	m.history = &WorkHistory{
		Entries: make([]interfaces.HistoryEntry, 0),
	}
	return m.Save()
}

// extractMonthYear extracts month/year from a date string
func extractMonthYear(dateStr string) string {
	// Parse the date and extract month/year
	if date, err := time.Parse(dateFormat, dateStr); err == nil {
		return date.Format("Jan 2006")
	}
	// Fallback: try to extract manually
	parts := strings.Fields(dateStr)
	if len(parts) >= 3 {
		return fmt.Sprintf("%s %s", parts[1], parts[2])
	}
	return dateStr
}

// truncateSummary truncates a summary to fit display
func truncateSummary(summary string, maxLen int) string {
	if len(summary) <= maxLen {
		return summary
	}
	return summary[:maxLen-3] + "..."
}

// findMostActiveMonth finds the month with the most commits
func findMostActiveMonth(monthGroups map[string]int) string {
	maxCommits := 0
	mostActive := ""

	for month, commits := range monthGroups {
		if commits > maxCommits {
			maxCommits = commits
			mostActive = month
		}
	}

	return mostActive
}

// Backward compatibility functions

// Load reads the work history from file or creates a new one (legacy function)
func Load() (*WorkHistory, error) {
	manager := NewManager(workHistoryFile)
	err := manager.Load()
	if err != nil {
		return nil, err
	}

	// Convert to legacy format
	historyManager := manager.(*Manager)
	return historyManager.GetHistory(), nil
}
