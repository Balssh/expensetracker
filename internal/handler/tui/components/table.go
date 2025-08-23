package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TableColumn represents a column in the table
type TableColumn struct {
	Key      string            // Unique identifier for the column
	Header   string            // Display header text
	Width    int               // Column width
	MinWidth int               // Minimum width
	MaxWidth int               // Maximum width (0 = unlimited)
	Sortable bool              // Whether this column can be sorted
	Align    lipgloss.Position // Text alignment
	Color    string            // Default color for this column
}

// TableRow represents a row of data in the table
type TableRow struct {
	ID     string                 // Unique row identifier
	Data   map[string]interface{} // Column key -> value mapping
	Colors map[string]string      // Column key -> color override
	Style  lipgloss.Style         // Row-level style override
}

// SortDirection represents sort direction
type SortDirection int

const (
	SortNone SortDirection = iota
	SortAsc
	SortDesc
)

// Table represents a data table component
type Table struct {
	// Configuration
	Columns    []TableColumn
	Rows       []TableRow
	Width      int
	Height     int
	ShowHeader bool
	ShowBorder bool

	// Selection
	Selected     int          // Currently selected row index
	Selectable   bool         // Whether rows can be selected
	MultiSelect  bool         // Whether multiple rows can be selected
	SelectedRows map[int]bool // Set of selected row indices

	// Sorting
	SortBy        string // Column key to sort by
	SortDirection SortDirection

	// Pagination
	Page      int // Current page (0-based)
	PageSize  int // Rows per page
	TotalRows int // Total number of rows (for pagination display)

	// Scrolling
	ScrollOffset int // Vertical scroll offset

	// Filtering
	FilterFunc func(row TableRow) bool

	// State
	Active   bool
	Disabled bool
	focused  bool

	// Styling
	HeaderStyle   lipgloss.Style
	RowStyle      lipgloss.Style
	SelectedStyle lipgloss.Style
	BorderStyle   lipgloss.Style
	FooterStyle   lipgloss.Style
}

// NewTable creates a new table component
func NewTable() *Table {
	return &Table{
		Columns:       make([]TableColumn, 0),
		Rows:          make([]TableRow, 0),
		SelectedRows:  make(map[int]bool),
		ShowHeader:    true,
		ShowBorder:    true,
		Selectable:    true,
		PageSize:      10,
		Width:         80,
		Height:        15,
		SortDirection: SortNone,

		// Default styles
		HeaderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true).
			Underline(true),
		RowStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")),
		SelectedStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("39")).
			Foreground(lipgloss.Color("15")),
		BorderStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),
		FooterStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
	}
}

// SetColumns sets the table columns
func (t *Table) SetColumns(columns []TableColumn) *Table {
	t.Columns = columns
	return t
}

// SetRows sets the table rows
func (t *Table) SetRows(rows []TableRow) *Table {
	t.Rows = rows
	t.TotalRows = len(rows)

	// Reset selection if out of bounds
	if t.Selected >= len(rows) {
		t.Selected = len(rows) - 1
	}
	if t.Selected < 0 && len(rows) > 0 {
		t.Selected = 0
	}

	return t
}

// SetSize sets the table dimensions
func (t *Table) SetSize(width, height int) *Table {
	t.Width = width
	t.Height = height
	return t
}

// SetSelectable enables/disables row selection
func (t *Table) SetSelectable(selectable bool) *Table {
	t.Selectable = selectable
	return t
}

// SetMultiSelect enables/disables multi-row selection
func (t *Table) SetMultiSelect(multiSelect bool) *Table {
	t.MultiSelect = multiSelect
	return t
}

// SetPageSize sets the number of rows per page
func (t *Table) SetPageSize(pageSize int) *Table {
	t.PageSize = pageSize
	return t
}

// SetFilterFunc sets a function to filter rows
func (t *Table) SetFilterFunc(filterFunc func(row TableRow) bool) *Table {
	t.FilterFunc = filterFunc
	return t
}

// Focus gives focus to the table
func (t *Table) Focus() {
	t.focused = true
	t.Active = true
}

// Blur removes focus from the table
func (t *Table) Blur() {
	t.focused = false
	t.Active = false
}

// GetSelectedRow returns the currently selected row
func (t *Table) GetSelectedRow() *TableRow {
	if t.Selected >= 0 && t.Selected < len(t.Rows) {
		return &t.Rows[t.Selected]
	}
	return nil
}

// GetSelectedRows returns all selected rows (for multi-select)
func (t *Table) GetSelectedRows() []TableRow {
	var selected []TableRow
	for i := range t.SelectedRows {
		if i >= 0 && i < len(t.Rows) {
			selected = append(selected, t.Rows[i])
		}
	}
	return selected
}

// ToggleRowSelection toggles selection for the current row
func (t *Table) ToggleRowSelection() {
	if !t.MultiSelect {
		return
	}

	if t.SelectedRows[t.Selected] {
		delete(t.SelectedRows, t.Selected)
	} else {
		t.SelectedRows[t.Selected] = true
	}
}

// ClearSelection clears all selections
func (t *Table) ClearSelection() {
	t.SelectedRows = make(map[int]bool)
}

// SortByColumn sorts the table by the specified column
func (t *Table) SortByColumn(columnKey string) {
	if t.SortBy == columnKey {
		// Toggle sort direction
		switch t.SortDirection {
		case SortNone:
			t.SortDirection = SortAsc
		case SortAsc:
			t.SortDirection = SortDesc
		case SortDesc:
			t.SortDirection = SortNone
			t.SortBy = ""
			return
		}
	} else {
		t.SortBy = columnKey
		t.SortDirection = SortAsc
	}

	t.sortRows()
}

// sortRows performs the actual sorting
func (t *Table) sortRows() {
	if t.SortBy == "" || t.SortDirection == SortNone {
		return
	}

	// Simple bubble sort for demonstration
	// In production, you'd want a more efficient algorithm
	for i := 0; i < len(t.Rows)-1; i++ {
		for j := 0; j < len(t.Rows)-i-1; j++ {
			val1 := t.Rows[j].Data[t.SortBy]
			val2 := t.Rows[j+1].Data[t.SortBy]

			shouldSwap := false
			if t.SortDirection == SortAsc {
				shouldSwap = compareValues(val1, val2) > 0
			} else {
				shouldSwap = compareValues(val1, val2) < 0
			}

			if shouldSwap {
				t.Rows[j], t.Rows[j+1] = t.Rows[j+1], t.Rows[j]
			}
		}
	}
}

// compareValues compares two interface{} values
func compareValues(a, b interface{}) int {
	// Convert to strings for comparison
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)

	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}
	return 0
}

// NextPage moves to the next page
func (t *Table) NextPage() {
	maxPage := (t.TotalRows - 1) / t.PageSize
	if t.Page < maxPage {
		t.Page++
		t.Selected = 0
	}
}

// PrevPage moves to the previous page
func (t *Table) PrevPage() {
	if t.Page > 0 {
		t.Page--
		t.Selected = 0
	}
}

// Update handles table events
func (t *Table) Update(msg tea.Msg) (*Table, tea.Cmd) {
	if !t.focused || t.Disabled {
		return t, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if t.Selected > 0 {
				t.Selected--
			} else if t.Page > 0 {
				t.PrevPage()
				t.Selected = t.PageSize - 1
				if t.Selected >= len(t.getVisibleRows()) {
					t.Selected = len(t.getVisibleRows()) - 1
				}
			}
		case "down", "j":
			visibleRows := t.getVisibleRows()
			if t.Selected < len(visibleRows)-1 {
				t.Selected++
			} else {
				maxPage := (t.TotalRows - 1) / t.PageSize
				if t.Page < maxPage {
					t.NextPage()
				}
			}
		case "left", "h":
			t.PrevPage()
		case "right", "l":
			t.NextPage()
		case "home":
			t.Selected = 0
		case "end":
			visibleRows := t.getVisibleRows()
			t.Selected = len(visibleRows) - 1
		case "pgup":
			t.PrevPage()
		case "pgdn":
			t.NextPage()
		case " ":
			if t.MultiSelect {
				t.ToggleRowSelection()
			}
		case "ctrl+a":
			if t.MultiSelect {
				// Select all visible rows
				visibleRows := t.getVisibleRows()
				for i := range visibleRows {
					t.SelectedRows[i] = true
				}
			}
		case "ctrl+d":
			if t.MultiSelect {
				t.ClearSelection()
			}
		}
	}

	return t, nil
}

// getVisibleRows returns the rows visible on the current page
func (t *Table) getVisibleRows() []TableRow {
	filteredRows := t.Rows

	// Apply filter
	if t.FilterFunc != nil {
		var filtered []TableRow
		for _, row := range t.Rows {
			if t.FilterFunc(row) {
				filtered = append(filtered, row)
			}
		}
		filteredRows = filtered
	}

	// Apply pagination
	start := t.Page * t.PageSize
	end := start + t.PageSize
	if end > len(filteredRows) {
		end = len(filteredRows)
	}
	if start > len(filteredRows) {
		start = len(filteredRows)
	}

	return filteredRows[start:end]
}

// View renders the table
func (t *Table) View() string {
	if len(t.Columns) == 0 {
		return "No columns defined"
	}

	visibleRows := t.getVisibleRows()

	var content []string

	// Header
	if t.ShowHeader {
		header := t.renderHeader()
		content = append(content, header)
	}

	// Rows
	if len(visibleRows) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Align(lipgloss.Center).
			Width(t.Width - 4).
			Render("No data to display")
		content = append(content, emptyMsg)
	} else {
		for i, row := range visibleRows {
			rowContent := t.renderRow(row, i == t.Selected, t.SelectedRows[i])
			content = append(content, rowContent)
		}
	}

	// Footer with pagination info
	if t.TotalRows > t.PageSize {
		footer := t.renderFooter()
		content = append(content, footer)
	}

	tableContent := lipgloss.JoinVertical(lipgloss.Left, content...)

	if t.ShowBorder {
		return t.BorderStyle.Width(t.Width).Height(t.Height).Render(tableContent)
	}

	return tableContent
}

// renderHeader renders the table header
func (t *Table) renderHeader() string {
	var headerCells []string

	for _, col := range t.Columns {
		headerText := col.Header

		// Add sort indicator
		if t.SortBy == col.Key {
			switch t.SortDirection {
			case SortAsc:
				headerText += " ▲"
			case SortDesc:
				headerText += " ▼"
			}
		}

		cellStyle := t.HeaderStyle.Width(col.Width)
		if col.Align != 0 {
			cellStyle = cellStyle.Align(col.Align)
		}

		cell := cellStyle.Render(headerText)
		headerCells = append(headerCells, cell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, headerCells...)
}

// renderRow renders a single table row
func (t *Table) renderRow(row TableRow, isSelected, isMultiSelected bool) string {
	var rowCells []string

	style := t.RowStyle
	if isSelected && t.Selectable {
		style = t.SelectedStyle
	}

	for _, col := range t.Columns {
		value := row.Data[col.Key]
		text := fmt.Sprintf("%v", value)

		// Truncate if too long
		if len(text) > col.Width-2 {
			text = text[:col.Width-5] + "..."
		}

		cellStyle := style.Width(col.Width)
		if col.Align != 0 {
			cellStyle = cellStyle.Align(col.Align)
		}

		// Apply column color
		if col.Color != "" {
			cellStyle = cellStyle.Foreground(lipgloss.Color(col.Color))
		}

		// Apply row-specific color override
		if rowColor, exists := row.Colors[col.Key]; exists {
			cellStyle = cellStyle.Foreground(lipgloss.Color(rowColor))
		}

		// Multi-select indicator
		if t.MultiSelect && isMultiSelected {
			text = "✓ " + text
		}

		cell := cellStyle.Render(text)
		rowCells = append(rowCells, cell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, rowCells...)
}

// renderFooter renders the pagination footer
func (t *Table) renderFooter() string {
	totalPages := (t.TotalRows + t.PageSize - 1) / t.PageSize
	currentPage := t.Page + 1

	paginationInfo := fmt.Sprintf("Page %d of %d • %d total rows",
		currentPage, totalPages, t.TotalRows)

	if t.MultiSelect && len(t.SelectedRows) > 0 {
		paginationInfo += fmt.Sprintf(" • %d selected", len(t.SelectedRows))
	}

	controls := "←/→: Pages • ↑/↓: Select"
	if t.MultiSelect {
		controls += " • Space: Toggle • Ctrl+A: All"
	}

	footer := lipgloss.JoinHorizontal(lipgloss.Left,
		t.FooterStyle.Render(paginationInfo),
		t.FooterStyle.Align(lipgloss.Right).Width(t.Width-len(paginationInfo)).Render(controls),
	)

	return footer
}
