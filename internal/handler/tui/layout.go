package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Layout configuration constants
const (
	MinTerminalWidth  = 80
	MinTerminalHeight = 24
	MaxContentWidth   = 120
	DefaultPadding    = 2
)

// CenterConfig holds configuration for centering content
type CenterConfig struct {
	Width      int
	Height     int
	MinWidth   int
	MaxWidth   int
	HorizontalPadding int
	VerticalPadding   int
}

// NewCenterConfig creates a sensible default centering configuration
func NewCenterConfig(termWidth, termHeight int) CenterConfig {
	return CenterConfig{
		Width:             termWidth,
		Height:            termHeight,
		MinWidth:          MinTerminalWidth,
		MaxWidth:          MaxContentWidth,
		HorizontalPadding: DefaultPadding,
		VerticalPadding:   1,
	}
}

// CalculateContentWidth determines the optimal content width for centering
func (c CenterConfig) CalculateContentWidth() int {
	availableWidth := c.Width - (c.HorizontalPadding * 2)
	
	// Ensure minimum width
	if availableWidth < c.MinWidth {
		return c.MinWidth
	}
	
	// Enforce maximum width for better readability
	if availableWidth > c.MaxWidth {
		return c.MaxWidth
	}
	
	return availableWidth
}

// CalculateContentHeight determines available content height
func (c CenterConfig) CalculateContentHeight() int {
	return c.Height - (c.VerticalPadding * 2) - 2 // Account for borders
}

// CenterHorizontally centers content horizontally within the terminal
func CenterHorizontally(content string, config CenterConfig) string {
	contentWidth := config.CalculateContentWidth()
	
	// If content is already wider than available space, return as-is
	lines := strings.Split(content, "\n")
	maxLineWidth := 0
	for _, line := range lines {
		if len(line) > maxLineWidth {
			maxLineWidth = len(line)
		}
	}
	
	if maxLineWidth >= contentWidth {
		return content
	}
	
	// Calculate horizontal padding
	totalPadding := contentWidth - maxLineWidth
	leftPadding := totalPadding / 2
	
	// Apply padding to each line
	var centeredLines []string
	for _, line := range lines {
		padding := strings.Repeat(" ", leftPadding)
		centeredLines = append(centeredLines, padding+line)
	}
	
	return strings.Join(centeredLines, "\n")
}

// WrapWithBorder wraps content with a styled border
func WrapWithBorder(content string, style lipgloss.Style) string {
	return style.Render(content)
}

// CreateMainApplicationBorder creates the main application container
func CreateMainApplicationBorder(content string, config CenterConfig) string {
	contentWidth := config.CalculateContentWidth()
	
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Width(contentWidth).
		Padding(1, 2).
		Align(lipgloss.Center)
	
	return style.Render(content)
}

// CreatePanelBorder creates a border for individual panels
func CreatePanelBorder(content string, title string, width int) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorNeutral).
		Width(width-4). // Account for border and padding
		Padding(0, 1)
	
	if title != "" {
		style = style.BorderTop(true).
			BorderTopForeground(colorNeutral)
		// Add title to the border (Lipgloss limitation - title would need custom implementation)
	}
	
	return style.Render(content)
}

// CreateThreePanelLayout creates a 3-panel vertical layout
func CreateThreePanelLayout(topPanel, middlePanel, bottomPanel string, config CenterConfig) string {
	contentWidth := config.CalculateContentWidth()
	contentHeight := config.CalculateContentHeight()
	
	// Calculate panel heights - top panel gets less space, middle gets most
	topHeight := 8      // Fixed height for summary
	bottomHeight := 3   // Fixed height for help
	middleHeight := contentHeight - topHeight - bottomHeight - 6 // Account for borders and spacing
	
	if middleHeight < 5 {
		middleHeight = 5 // Minimum height for transaction list
	}
	
	// Create styled panels
	topPanelStyled := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorNeutral).
		Width(contentWidth-4).
		Height(topHeight).
		Padding(1, 2).
		Render(topPanel)
	
	middlePanelStyled := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorNeutral).
		Width(contentWidth-4).
		Height(middleHeight).
		Padding(1, 2).
		Render(middlePanel)
	
	bottomPanelStyled := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorNeutral).
		Width(contentWidth-4).
		Height(bottomHeight).
		Padding(0, 2).
		Render(bottomPanel)
	
	// Combine panels with spacing
	return lipgloss.JoinVertical(
		lipgloss.Left,
		topPanelStyled,
		"", // Spacing line
		middlePanelStyled,
		"", // Spacing line
		bottomPanelStyled,
	)
}

// CreateTableAlignment creates properly aligned table columns
type TableColumn struct {
	Header    string
	Width     int
	Alignment lipgloss.Position
}

// FormatTableRow formats a table row with proper alignment
func FormatTableRow(columns []TableColumn, values []string) string {
	if len(columns) != len(values) {
		return strings.Join(values, " | ") // Fallback
	}
	
	var parts []string
	for i, col := range columns {
		value := values[i]
		
		// Truncate if necessary
		if len(value) > col.Width {
			if col.Width > 3 {
				value = value[:col.Width-3] + "..."
			} else {
				value = value[:col.Width]
			}
		}
		
		// Apply alignment
		var formatted string
		switch col.Alignment {
		case lipgloss.Right:
			formatted = strings.Repeat(" ", col.Width-len(value)) + value
		case lipgloss.Center:
			totalPad := col.Width - len(value)
			leftPad := totalPad / 2
			rightPad := totalPad - leftPad
			formatted = strings.Repeat(" ", leftPad) + value + strings.Repeat(" ", rightPad)
		default: // Left alignment
			formatted = value + strings.Repeat(" ", col.Width-len(value))
		}
		
		parts = append(parts, formatted)
	}
	
	return strings.Join(parts, " ")
}

// CreateTableHeader creates a styled table header
func CreateTableHeader(columns []TableColumn) string {
	var headers []string
	for _, col := range columns {
		headers = append(headers, col.Header)
	}
	
	headerRow := FormatTableRow(columns, headers)
	
	// Apply style without padding that could cause wrapping
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorTextPrimary).
		Background(colorPrimary)
	
	return style.Render(headerRow)
}

// CreateTableSeparator creates a table separator line
func CreateTableSeparator(totalWidth int) string {
	return strings.Repeat("─", totalWidth)
}

// ResponsiveBreakpoint determines the layout breakpoint based on terminal width
type Breakpoint int

const (
	BreakpointNarrow Breakpoint = iota
	BreakpointStandard
	BreakpointWide
)

// GetBreakpoint returns the current responsive breakpoint
func GetBreakpoint(width int) Breakpoint {
	switch {
	case width < 100:
		return BreakpointNarrow
	case width < 120:
		return BreakpointStandard
	default:
		return BreakpointWide
	}
}

// TruncateWithEllipsis truncates a string to a maximum width with ellipsis
func TruncateWithEllipsis(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	
	if maxWidth <= 3 {
		return s[:maxWidth]
	}
	
	return s[:maxWidth-3] + "..."
}

// PadString pads a string to a specific width
func PadString(s string, width int, alignment lipgloss.Position) string {
	if len(s) >= width {
		return s
	}
	
	totalPad := width - len(s)
	
	switch alignment {
	case lipgloss.Right:
		return strings.Repeat(" ", totalPad) + s
	case lipgloss.Center:
		leftPad := totalPad / 2
		rightPad := totalPad - leftPad
		return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
	default: // Left
		return s + strings.Repeat(" ", totalPad)
	}
}