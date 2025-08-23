package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ToggleOption represents an option in a toggle
type ToggleOption struct {
	Label string
	Value interface{}
	Color string // Color for styling this option
}

// Toggle represents a toggle/radio button group component
type Toggle struct {
	Label       string
	Options     []ToggleOption
	Selected    int
	Active      bool
	Disabled    bool
	Required    bool
	
	// Layout
	Horizontal  bool // If true, options are laid out horizontally
	Width       int
	
	// Validation
	ValidationState ValidationState
	ErrorMessage   string
	
	// Internal state
	focused     bool
}

// NewToggle creates a new toggle component
func NewToggle(label string) *Toggle {
	return &Toggle{
		Label:      label,
		Selected:   0, // First option selected by default
		Horizontal: true,
		Width:      40,
		ValidationState: ValidationNone,
	}
}

// SetOptions sets the options for the toggle
func (t *Toggle) SetOptions(options []ToggleOption) *Toggle {
	t.Options = options
	if t.Selected >= len(options) {
		t.Selected = 0
	}
	return t
}

// SetHorizontal sets the layout direction
func (t *Toggle) SetHorizontal(horizontal bool) *Toggle {
	t.Horizontal = horizontal
	return t
}

// SetWidth sets the width of the toggle
func (t *Toggle) SetWidth(width int) *Toggle {
	t.Width = width
	return t
}

// SetRequired marks the toggle as required
func (t *Toggle) SetRequired(required bool) *Toggle {
	t.Required = required
	return t
}

// Focus gives focus to the toggle
func (t *Toggle) Focus() {
	t.focused = true
	t.Active = true
}

// Blur removes focus from the toggle
func (t *Toggle) Blur() {
	t.focused = false
	t.Active = false
	t.validate()
}

// SetSelected sets the selected option by index
func (t *Toggle) SetSelected(index int) {
	if index >= 0 && index < len(t.Options) {
		t.Selected = index
		t.validate()
	}
}

// SetSelectedByValue sets the selected option by value
func (t *Toggle) SetSelectedByValue(value interface{}) {
	for i, option := range t.Options {
		if option.Value == value {
			t.Selected = i
			t.validate()
			return
		}
	}
}

// GetSelected returns the currently selected option
func (t *Toggle) GetSelected() *ToggleOption {
	if t.Selected >= 0 && t.Selected < len(t.Options) {
		return &t.Options[t.Selected]
	}
	return nil
}

// GetSelectedValue returns the value of the currently selected option
func (t *Toggle) GetSelectedValue() interface{} {
	if selected := t.GetSelected(); selected != nil {
		return selected.Value
	}
	return nil
}

// IsValid returns true if the toggle has a valid selection
func (t *Toggle) IsValid() bool {
	if t.Required {
		return t.Selected >= 0 && t.Selected < len(t.Options)
	}
	return true
}

// Update handles toggle events
func (t *Toggle) Update(msg tea.Msg) (*Toggle, tea.Cmd) {
	if !t.focused || t.Disabled {
		return t, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "up":
			if t.Selected > 0 {
				t.Selected--
				t.validate()
			}
		case "right", "down":
			if t.Selected < len(t.Options)-1 {
				t.Selected++
				t.validate()
			}
		case " ", "enter":
			// Space or Enter toggles between options
			t.Selected = (t.Selected + 1) % len(t.Options)
			t.validate()
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// Number keys for direct selection
			index := int(msg.String()[0] - '1')
			if index >= 0 && index < len(t.Options) {
				t.Selected = index
				t.validate()
			}
		}
	}

	return t, nil
}

// View renders the toggle component
func (t *Toggle) View() string {
	// Label with consistent styling
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)
	
	labelText := t.Label + ":"
	if t.Required {
		// Keep label color consistent, just add asterisk
		labelText = t.Label + " *:"
	}

	var optionViews []string
	
	// Render each option
	for i, option := range t.Options {
		// Determine styling
		optionStyle := lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 1, 0, 0)
		
		var indicator string
		var textColor string
		
		if i == t.Selected {
			// Selected option
			indicator = "●"
			textColor = option.Color
			if textColor == "" {
				textColor = "15" // Default white
			}
			
			if t.focused {
				// Add border when focused
				optionStyle = optionStyle.
					Border(lipgloss.NormalBorder()).
					BorderForeground(lipgloss.Color("39"))
			}
		} else {
			// Unselected option
			indicator = "○"
			textColor = "240" // Gray
			
			if t.focused {
				textColor = "15" // More visible when focused
			}
		}
		
		if t.Disabled {
			textColor = "240"
		}
		
		optionStyle = optionStyle.Foreground(lipgloss.Color(textColor))
		
		optionText := indicator + " " + option.Label
		optionViews = append(optionViews, optionStyle.Render(optionText))
	}
	
	// Layout options
	var optionsView string
	if t.Horizontal {
		optionsView = lipgloss.JoinHorizontal(lipgloss.Left, optionViews...)
	} else {
		optionsView = lipgloss.JoinVertical(lipgloss.Left, optionViews...)
	}
	
	// Validation indicator
	var validationPart string
	switch t.ValidationState {
	case ValidationValid:
		validationPart = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Render(" ✓")
	case ValidationError:
		validationPart = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render(" ✗")
	}
	
	// Combine label and options
	if t.Horizontal {
		// Horizontal layout: Label + Options + Validation on same line
		firstLine := lipgloss.JoinHorizontal(lipgloss.Left, 
			labelStyle.Render(labelText), " ", optionsView, validationPart)
		
		if t.ErrorMessage != "" && t.ValidationState == ValidationError {
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Margin(0, 0, 0, len(t.Label)+2)
			return lipgloss.JoinVertical(lipgloss.Left, firstLine, errorStyle.Render(t.ErrorMessage))
		}
		
		return firstLine
	} else {
		// Vertical layout: Label on top, options below
		result := lipgloss.JoinVertical(lipgloss.Left,
			labelStyle.Render(labelText)+validationPart,
			optionsView)
		
		if t.ErrorMessage != "" && t.ValidationState == ValidationError {
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("196"))
			result = lipgloss.JoinVertical(lipgloss.Left, result, errorStyle.Render(t.ErrorMessage))
		}
		
		return result
	}
}

// validate performs validation on the current selection
func (t *Toggle) validate() {
	if t.Required && (t.Selected < 0 || t.Selected >= len(t.Options)) {
		t.ValidationState = ValidationError
		t.ErrorMessage = "Please select an option"
	} else {
		t.ValidationState = ValidationValid
		t.ErrorMessage = ""
	}
}