package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DropdownOption represents an option in the dropdown
type DropdownOption struct {
	Label string
	Value interface{}
	ID    int // For categories, this would be the category ID
}

// Dropdown represents a dropdown selection component
type Dropdown struct {
	Label       string
	Options     []DropdownOption
	Selected    int
	Open        bool
	Width       int
	Active      bool
	Disabled    bool
	Required    bool
	Placeholder string
	
	// Validation
	ValidationState ValidationState
	ErrorMessage   string
	
	// Internal state
	focused     bool
	highlighted int // Currently highlighted option when open
	maxVisible  int // Maximum visible options before scrolling
	scrollOffset int // Scroll offset for long lists
}

// NewDropdown creates a new dropdown component
func NewDropdown(label string) *Dropdown {
	return &Dropdown{
		Label:       label,
		Width:       25,
		Selected:    -1, // No selection initially
		highlighted: 0,
		maxVisible:  6, // Show max 6 options at once
		ValidationState: ValidationNone,
	}
}

// SetOptions sets the options for the dropdown
func (d *Dropdown) SetOptions(options []DropdownOption) *Dropdown {
	d.Options = options
	if d.Selected >= len(options) {
		d.Selected = -1
	}
	if d.highlighted >= len(options) {
		d.highlighted = 0
	}
	return d
}

// SetWidth sets the width of the dropdown
func (d *Dropdown) SetWidth(width int) *Dropdown {
	d.Width = width
	return d
}

// SetRequired marks the dropdown as required
func (d *Dropdown) SetRequired(required bool) *Dropdown {
	d.Required = required
	return d
}

// SetPlaceholder sets the placeholder text
func (d *Dropdown) SetPlaceholder(placeholder string) *Dropdown {
	d.Placeholder = placeholder
	return d
}

// Focus gives focus to the dropdown
func (d *Dropdown) Focus() {
	d.focused = true
	d.Active = true
}

// Blur removes focus from the dropdown
func (d *Dropdown) Blur() {
	d.focused = false
	d.Active = false
	d.Open = false
	d.validate()
}

// SetSelected sets the selected option by index
func (d *Dropdown) SetSelected(index int) {
	if index >= 0 && index < len(d.Options) {
		d.Selected = index
		d.validate()
	}
}

// SetSelectedByValue sets the selected option by value
func (d *Dropdown) SetSelectedByValue(value interface{}) {
	for i, option := range d.Options {
		if option.Value == value {
			d.Selected = i
			d.validate()
			return
		}
	}
}

// GetSelected returns the currently selected option
func (d *Dropdown) GetSelected() *DropdownOption {
	if d.Selected >= 0 && d.Selected < len(d.Options) {
		return &d.Options[d.Selected]
	}
	return nil
}

// GetSelectedValue returns the value of the currently selected option
func (d *Dropdown) GetSelectedValue() interface{} {
	if selected := d.GetSelected(); selected != nil {
		return selected.Value
	}
	return nil
}

// IsValid returns true if the dropdown has a valid selection
func (d *Dropdown) IsValid() bool {
	if d.Required {
		return d.Selected >= 0 && d.Selected < len(d.Options)
	}
	return true
}

// Update handles dropdown events
func (d *Dropdown) Update(msg tea.Msg) (*Dropdown, tea.Cmd) {
	if !d.focused || d.Disabled {
		return d, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			if d.Open {
				// Select the highlighted option
				d.Selected = d.highlighted
				d.Open = false
				d.validate()
			} else {
				// Open the dropdown
				d.Open = true
				if d.Selected >= 0 {
					d.highlighted = d.Selected
				}
			}
		case "esc":
			if d.Open {
				d.Open = false
			}
		case "up":
			if d.Open {
				if d.highlighted > 0 {
					d.highlighted--
					d.updateScrollOffset()
				}
			} else if !d.Open && len(d.Options) > 0 {
				// Quick selection without opening
				if d.Selected > 0 {
					d.Selected--
					d.validate()
				}
			}
		case "down":
			if d.Open {
				if d.highlighted < len(d.Options)-1 {
					d.highlighted++
					d.updateScrollOffset()
				}
			} else if !d.Open && len(d.Options) > 0 {
				// Quick selection without opening
				if d.Selected < len(d.Options)-1 {
					d.Selected++
					d.validate()
				} else if d.Selected == -1 && len(d.Options) > 0 {
					d.Selected = 0
					d.validate()
				}
			}
		case "home":
			if d.Open {
				d.highlighted = 0
				d.scrollOffset = 0
			}
		case "end":
			if d.Open {
				d.highlighted = len(d.Options) - 1
				d.updateScrollOffset()
			}
		}
	}

	return d, nil
}

// View renders the dropdown component
func (d *Dropdown) View() string {
	// Label with consistent styling
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)
	
	labelText := d.Label + ":"
	if d.Required {
		// Keep label color consistent, just add asterisk
		labelText = d.Label + " *:"
	}

	// Determine border color based on state
	borderColor := "240" // Default gray
	if d.focused {
		borderColor = "39" // Blue when focused
	}
	if d.ValidationState == ValidationError {
		borderColor = "196" // Red for errors
	} else if d.ValidationState == ValidationValid {
		borderColor = "46" // Green for valid
	}

	// Main dropdown field
	var displayText string
	if d.Selected >= 0 && d.Selected < len(d.Options) {
		displayText = d.Options[d.Selected].Label
	} else if d.Placeholder != "" {
		displayText = d.Placeholder
	} else {
		displayText = "Select an option..."
	}

	// Add dropdown arrow
	arrow := "▼"
	if d.Open {
		arrow = "▲"
	}

	fieldStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Width(d.Width - 3) // Account for arrow and padding

	if d.Disabled {
		fieldStyle = fieldStyle.Foreground(lipgloss.Color("240"))
	} else if d.Selected == -1 && d.Placeholder != "" {
		fieldStyle = fieldStyle.Foreground(lipgloss.Color("240"))
	}

	field := fieldStyle.Render(displayText) + " " + arrow

	// Validation indicator
	var validationPart string
	switch d.ValidationState {
	case ValidationValid:
		validationPart = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Render(" ✓")
	case ValidationError:
		validationPart = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render(" ✗")
	}

	// First line with label and field
	firstLine := lipgloss.JoinHorizontal(lipgloss.Left, 
		labelStyle.Render(labelText), " ", field, validationPart)

	// If not open, just return the first line
	if !d.Open {
		if d.ErrorMessage != "" && d.ValidationState == ValidationError {
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Margin(0, 0, 0, len(d.Label)+2)
			return lipgloss.JoinVertical(lipgloss.Left, firstLine, errorStyle.Render(d.ErrorMessage))
		}
		return firstLine
	}

	// Build dropdown options list
	var optionLines []string
	
	start := d.scrollOffset
	end := start + d.maxVisible
	if end > len(d.Options) {
		end = len(d.Options)
	}

	for i := start; i < end; i++ {
		option := d.Options[i]
		optionText := option.Label
		
		// Truncate if too long
		maxOptWidth := d.Width - 4
		if len(optionText) > maxOptWidth {
			optionText = optionText[:maxOptWidth-3] + "..."
		}

		optionStyle := lipgloss.NewStyle().
			Padding(0, 1).
			Width(d.Width - 2)

		if i == d.highlighted {
			// Highlighted option
			optionStyle = optionStyle.
				Background(lipgloss.Color("39")).
				Foreground(lipgloss.Color("15"))
		} else if i == d.Selected {
			// Currently selected option
			optionStyle = optionStyle.
				Background(lipgloss.Color("240")).
				Foreground(lipgloss.Color("15"))
		}

		optionLines = append(optionLines, optionStyle.Render(optionText))
	}

	// Add scroll indicators if needed
	if len(d.Options) > d.maxVisible {
		if d.scrollOffset > 0 {
			scrollUp := lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Width(d.Width - 2).
				Align(lipgloss.Center).
				Render("▲ More above")
			optionLines = append([]string{scrollUp}, optionLines...)
		}
		if d.scrollOffset + d.maxVisible < len(d.Options) {
			scrollDown := lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Width(d.Width - 2).
				Align(lipgloss.Center).
				Render("▼ More below")
			optionLines = append(optionLines, scrollDown)
		}
	}

	// Options container
	optionsStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Margin(0, 0, 0, len(d.Label)+2)

	optionsList := optionsStyle.Render(lipgloss.JoinVertical(lipgloss.Left, optionLines...))

	return lipgloss.JoinVertical(lipgloss.Left, firstLine, optionsList)
}

// updateScrollOffset updates the scroll offset to keep highlighted item visible
func (d *Dropdown) updateScrollOffset() {
	if d.highlighted < d.scrollOffset {
		d.scrollOffset = d.highlighted
	} else if d.highlighted >= d.scrollOffset + d.maxVisible {
		d.scrollOffset = d.highlighted - d.maxVisible + 1
	}
}

// validate performs validation on the current selection
func (d *Dropdown) validate() {
	if d.Required && d.Selected == -1 {
		d.ValidationState = ValidationError
		d.ErrorMessage = "Please select an option"
	} else {
		d.ValidationState = ValidationValid
		d.ErrorMessage = ""
	}
}