package tui

import "github.com/charmbracelet/lipgloss"

// Color palette constants following the design system
const (
	colorPrimary   = lipgloss.Color("#0066cc")
	colorSuccess   = lipgloss.Color("#22c55e")
	colorWarning   = lipgloss.Color("#eab308")
	colorError     = lipgloss.Color("#ef4444")
	colorNeutral   = lipgloss.Color("#6b7280")
	
	// Text colors
	colorTextPrimary   = lipgloss.Color("#FAFAFA")
	colorTextSecondary = lipgloss.Color("#626262")
	colorTextMuted     = lipgloss.Color("#404040")
	
	// Background colors
	colorBackgroundInput    = lipgloss.Color("#333333")
	colorBackgroundSelected = lipgloss.Color("#0066cc")
)

// Application-level styles
var (
	// Main application container
	appContainerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(2, 3).
		Align(lipgloss.Center)

	// Application title
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorTextPrimary).
		Background(colorPrimary).
		Padding(0, 1)

	// Panel styles
	panelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorNeutral).
		Padding(1, 2)

	panelHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorTextPrimary).
		Background(colorNeutral).
		Padding(0, 1)
)

// Dashboard-specific styles
var (
	// Summary panel styles
	summaryBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorNeutral).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	summaryHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorTextPrimary).
		Align(lipgloss.Center)

	// Financial value styles
	incomeStyle = lipgloss.NewStyle().
		Foreground(colorSuccess).
		Bold(true)

	expenseStyle = lipgloss.NewStyle().
		Foreground(colorError).
		Bold(true)

	balancePositiveStyle = lipgloss.NewStyle().
		Foreground(colorSuccess).
		Bold(true)

	balanceNegativeStyle = lipgloss.NewStyle().
		Foreground(colorError).
		Bold(true)

	// Expense breakdown bar styles
	expenseBarStyle = lipgloss.NewStyle().
		Background(colorError).
		Foreground(colorTextPrimary)

	expenseBarCategoryStyle = lipgloss.NewStyle().
		Background(colorNeutral).
		Foreground(colorTextPrimary).
		Padding(0, 1)
)

// Table styles
var (
	tableHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorTextPrimary).
		Background(colorPrimary).
		Padding(0, 1)

	tableRowStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary)

	tableRowSelectedStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Background(colorBackgroundSelected).
		Padding(0, 1)

	tableRowAltStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Background(lipgloss.Color("#1a1a1a"))

	tableSeparatorStyle = lipgloss.NewStyle().
		Foreground(colorNeutral)
)

// Form styles
var (
	formStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2).
		MarginTop(1)

	formFieldLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorTextSecondary)

	formFieldStyle = lipgloss.NewStyle().
		Padding(0, 1).
		MarginBottom(1)
)

// Input styles
var (
	inputStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Background(colorBackgroundInput).
		Padding(0, 1)

	inputFocusedStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Background(colorPrimary).
		Padding(0, 1)

	inputErrorStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Background(colorError).
		Padding(0, 1)

	inputPlaceholderStyle = lipgloss.NewStyle().
		Foreground(colorTextMuted).
		Italic(true)
)

// Selection and dropdown styles
var (
	selectedItemStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Background(colorPrimary).
		Padding(0, 1)

	highlightedItemStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Background(colorBackgroundSelected).
		Padding(0, 1)

	dropdownStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2).
		MaxHeight(10)

	dropdownItemStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Padding(0, 1)

	dropdownItemSelectedStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Background(colorPrimary).
		Padding(0, 1)
)

// Status and message styles
var (
	successStyle = lipgloss.NewStyle().
		Foreground(colorSuccess).
		Bold(true)

	errorStyle = lipgloss.NewStyle().
		Foreground(colorError).
		Bold(true)

	warningStyle = lipgloss.NewStyle().
		Foreground(colorWarning).
		Bold(true)

	infoStyle = lipgloss.NewStyle().
		Foreground(colorNeutral).
		Bold(true)

	loadingStyle = lipgloss.NewStyle().
		Foreground(colorTextSecondary).
		Italic(true)
)

// Help and navigation styles
var (
	helpStyle = lipgloss.NewStyle().
		Foreground(colorTextSecondary).
		MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Bold(true)

	helpDescStyle = lipgloss.NewStyle().
		Foreground(colorTextSecondary)

	navigationStyle = lipgloss.NewStyle().
		Foreground(colorTextSecondary).
		Align(lipgloss.Center)
)

// Modal and dialog styles
var (
	modalOverlayStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")).
		Foreground(colorTextPrimary)

	modalStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(2, 3).
		Align(lipgloss.Center)

	modalHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorTextPrimary).
		Background(colorPrimary).
		Padding(0, 1).
		Align(lipgloss.Center)

	modalButtonStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Foreground(colorTextPrimary).
		Padding(0, 2)

	modalButtonFocusedStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Background(colorPrimary).
		Foreground(colorTextPrimary).
		Padding(0, 2)
)

// Search and filter styles
var (
	searchBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(0, 1)

	searchBoxFocusedStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Background(colorPrimary).
		Foreground(colorTextPrimary).
		Padding(0, 1)

	filterTagStyle = lipgloss.NewStyle().
		Background(colorNeutral).
		Foreground(colorTextPrimary).
		Padding(0, 1).
		MarginRight(1)

	activeFilterTagStyle = lipgloss.NewStyle().
		Background(colorPrimary).
		Foreground(colorTextPrimary).
		Padding(0, 1).
		MarginRight(1)
)

// Responsive utility styles
var (
	// For narrow terminals
	compactStyle = lipgloss.NewStyle().
		Padding(0, 1)

	// For wide terminals  
	expandedStyle = lipgloss.NewStyle().
		Padding(1, 3)
)
