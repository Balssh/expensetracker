package components

// Breakpoint represents a responsive design breakpoint
type Breakpoint struct {
	Name     string
	MinWidth int
	MaxWidth int
}

// LayoutConfig defines layout configuration for different breakpoints
type LayoutConfig struct {
	Breakpoint      Breakpoint
	FilterLayout    FilterLayoutType
	ShowHelp        bool
	TablePadding    int
	ComponentWidths map[string]int
	TableColumns    TableColumnConfig
}

// FilterLayoutType represents different filter layout options
type FilterLayoutType string

const (
	FilterLayoutHorizontal FilterLayoutType = "horizontal"
	FilterLayoutVertical   FilterLayoutType = "vertical"
	FilterLayoutCollapsed  FilterLayoutType = "collapsed"
)

// TableColumnConfig defines column behavior for different screen sizes
type TableColumnConfig struct {
	ShowAllColumns bool
	HiddenColumns  []string
	MinWidths      map[string]int
	MaxWidths      map[string]int
	PreferredWidths map[string]int
}

// ResponsiveConfig manages responsive design breakpoints and configurations
type ResponsiveConfig struct {
	Breakpoints []Breakpoint
	Layouts     map[string]LayoutConfig
}

// NewResponsiveConfig creates a default responsive configuration
func NewResponsiveConfig() *ResponsiveConfig {
	return &ResponsiveConfig{
		Breakpoints: []Breakpoint{
			{Name: "mobile", MinWidth: 0, MaxWidth: 49},
			{Name: "narrow", MinWidth: 50, MaxWidth: 79},
			{Name: "standard", MinWidth: 80, MaxWidth: 119},
			{Name: "wide", MinWidth: 120, MaxWidth: 9999},
		},
		Layouts: map[string]LayoutConfig{
			"mobile": {
				Breakpoint:   Breakpoint{Name: "mobile", MinWidth: 0, MaxWidth: 49},
				FilterLayout: FilterLayoutCollapsed,
				ShowHelp:     false,
				TablePadding: 2,
				ComponentWidths: map[string]int{
					"dropdown": 12,
					"search":   15,
				},
				TableColumns: TableColumnConfig{
					ShowAllColumns: false,
					HiddenColumns:  []string{"description"},
					MinWidths: map[string]int{
						"date":     6,
						"type":     4,
						"category": 6,
						"amount":   7,
					},
					MaxWidths: map[string]int{
						"date":     8,
						"type":     6,
						"category": 10,
						"amount":   9,
					},
					PreferredWidths: map[string]int{
						"date":     8,
						"type":     6,
						"category": 8,
						"amount":   9,
					},
				},
			},
			"narrow": {
				Breakpoint:   Breakpoint{Name: "narrow", MinWidth: 50, MaxWidth: 79},
				FilterLayout: FilterLayoutVertical,
				ShowHelp:     false,
				TablePadding: 4,
				ComponentWidths: map[string]int{
					"dropdown": 15,
					"search":   20,
				},
				TableColumns: TableColumnConfig{
					ShowAllColumns: true,
					HiddenColumns:  []string{},
					MinWidths: map[string]int{
						"date":        6,
						"type":        6,
						"category":    8,
						"description": 12,
						"amount":      7,
					},
					MaxWidths: map[string]int{
						"date":        8,
						"type":        10,
						"category":    12,
						"description": 20,
						"amount":      9,
					},
					PreferredWidths: map[string]int{
						"date":        8,
						"type":        8,
						"category":    10,
						"description": 15,
						"amount":      9,
					},
				},
			},
			"standard": {
				Breakpoint:   Breakpoint{Name: "standard", MinWidth: 80, MaxWidth: 119},
				FilterLayout: FilterLayoutHorizontal,
				ShowHelp:     true,
				TablePadding: 6,
				ComponentWidths: map[string]int{
					"dropdown": 18,
					"search":   25,
				},
				TableColumns: TableColumnConfig{
					ShowAllColumns: true,
					HiddenColumns:  []string{},
					MinWidths: map[string]int{
						"date":        6,
						"type":        6,
						"category":    8,
						"description": 12,
						"amount":      7,
					},
					MaxWidths: map[string]int{
						"date":        12,
						"type":        12,
						"category":    15,
						"description": 30,
						"amount":      12,
					},
					PreferredWidths: map[string]int{
						"date":        8,
						"type":        10,
						"category":    12,
						"description": 20,
						"amount":      9,
					},
				},
			},
			"wide": {
				Breakpoint:   Breakpoint{Name: "wide", MinWidth: 120, MaxWidth: 9999},
				FilterLayout: FilterLayoutHorizontal,
				ShowHelp:     true,
				TablePadding: 6,
				ComponentWidths: map[string]int{
					"dropdown": 20,
					"search":   30,
				},
				TableColumns: TableColumnConfig{
					ShowAllColumns: true,
					HiddenColumns:  []string{},
					MinWidths: map[string]int{
						"date":        6,
						"type":        6,
						"category":    8,
						"description": 12,
						"amount":      7,
					},
					MaxWidths: map[string]int{
						"date":        12,
						"type":        12,
						"category":    20,
						"description": 40,
						"amount":      12,
					},
					PreferredWidths: map[string]int{
						"date":        8,
						"type":        12,
						"category":    15,
						"description": 25,
						"amount":      10,
					},
				},
			},
		},
	}
}

// GetLayoutForWidth returns the appropriate layout configuration for a given width
func (r *ResponsiveConfig) GetLayoutForWidth(width int) LayoutConfig {
	for _, breakpoint := range r.Breakpoints {
		if width >= breakpoint.MinWidth && width <= breakpoint.MaxWidth {
			if layout, exists := r.Layouts[breakpoint.Name]; exists {
				return layout
			}
		}
	}
	
	// Fallback to standard layout
	return r.Layouts["standard"]
}

// GetBreakpointName returns the breakpoint name for a given width
func (r *ResponsiveConfig) GetBreakpointName(width int) string {
	for _, breakpoint := range r.Breakpoints {
		if width >= breakpoint.MinWidth && width <= breakpoint.MaxWidth {
			return breakpoint.Name
		}
	}
	return "standard"
}

// AdjustTableColumns calculates column widths based on available space and layout config
func (r *ResponsiveConfig) AdjustTableColumns(columns []TableColumn, availableWidth int, layout LayoutConfig) []TableColumn {
	if len(columns) == 0 {
		return columns
	}
	
	// Filter out hidden columns
	visibleColumns := make([]TableColumn, 0)
	for _, col := range columns {
		hide := false
		for _, hiddenCol := range layout.TableColumns.HiddenColumns {
			if col.Key == hiddenCol {
				hide = true
				break
			}
		}
		if !hide {
			visibleColumns = append(visibleColumns, col)
		}
	}
	
	if len(visibleColumns) == 0 {
		return columns
	}
	
	// Account for padding
	contentWidth := availableWidth - layout.TablePadding
	if contentWidth < 30 {
		contentWidth = 30
	}
	
	// Calculate minimum width required
	minTotalWidth := 0
	for _, col := range visibleColumns {
		if minWidth, exists := layout.TableColumns.MinWidths[col.Key]; exists {
			minTotalWidth += minWidth
		} else {
			minTotalWidth += 5 // Fallback minimum
		}
	}
	
	// If screen too narrow, use minimum widths
	if contentWidth <= minTotalWidth {
		for i := range visibleColumns {
			if minWidth, exists := layout.TableColumns.MinWidths[visibleColumns[i].Key]; exists {
				visibleColumns[i].Width = minWidth
			} else {
				visibleColumns[i].Width = 5
			}
		}
		return visibleColumns
	}
	
	// Calculate proportional distribution
	totalPreferred := 0
	for _, col := range visibleColumns {
		if preferredWidth, exists := layout.TableColumns.PreferredWidths[col.Key]; exists {
			totalPreferred += preferredWidth
		} else {
			totalPreferred += col.Width
		}
	}
	
	// Distribute width proportionally
	remaining := contentWidth
	for i := range visibleColumns {
		col := &visibleColumns[i]
		
		// Get preferred width
		preferredWidth := col.Width
		if pref, exists := layout.TableColumns.PreferredWidths[col.Key]; exists {
			preferredWidth = pref
		}
		
		// Calculate proportional width
		proportionalWidth := (preferredWidth * contentWidth) / totalPreferred
		
		// Apply min/max constraints
		minWidth := 5
		if min, exists := layout.TableColumns.MinWidths[col.Key]; exists {
			minWidth = min
		}
		
		maxWidth := contentWidth / 2 // Reasonable maximum
		if max, exists := layout.TableColumns.MaxWidths[col.Key]; exists && max < maxWidth {
			maxWidth = max
		}
		
		if proportionalWidth < minWidth {
			proportionalWidth = minWidth
		} else if proportionalWidth > maxWidth {
			proportionalWidth = maxWidth
		}
		
		col.Width = proportionalWidth
		remaining -= proportionalWidth
	}
	
	// Distribute any remaining width to the description column if it exists
	if remaining > 0 {
		for i := range visibleColumns {
			if visibleColumns[i].Key == "description" {
				maxWidth := contentWidth / 2
				if max, exists := layout.TableColumns.MaxWidths["description"]; exists {
					maxWidth = max
				}
				additionalWidth := remaining
				if visibleColumns[i].Width+additionalWidth > maxWidth {
					additionalWidth = maxWidth - visibleColumns[i].Width
				}
				if additionalWidth > 0 {
					visibleColumns[i].Width += additionalWidth
				}
				break
			}
		}
	}
	
	return visibleColumns
}