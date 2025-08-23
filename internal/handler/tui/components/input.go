package components

import (
	"errors"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ValidationState represents the validation state of an input
type ValidationState int

const (
	ValidationNone ValidationState = iota
	ValidationValid
	ValidationError
	ValidationWarning
)

// InputType represents different types of inputs
type InputType int

const (
	InputText InputType = iota
	InputNumber
	InputEmail
	InputDate
)

// Input represents a form input component
type Input struct {
	Label        string
	Value        string
	Placeholder  string
	Width        int
	MaxLength    int
	Required     bool
	InputType    InputType
	Active       bool
	Disabled     bool
	
	// Validation
	ValidationFunc func(string) error
	ValidationState ValidationState
	ErrorMessage   string
	
	// Internal state
	cursor       int
	focused      bool
	blinkTick    int
}

// NewInput creates a new input component
func NewInput(label string) *Input {
	return &Input{
		Label:       label,
		Width:       20,
		MaxLength:   200,
		InputType:   InputText,
		cursor:      0,
		ValidationState: ValidationNone,
	}
}

// SetValidation sets the validation function for the input
func (i *Input) SetValidation(fn func(string) error) *Input {
	i.ValidationFunc = fn
	return i
}

// SetRequired marks the input as required
func (i *Input) SetRequired(required bool) *Input {
	i.Required = required
	return i
}

// SetWidth sets the width of the input
func (i *Input) SetWidth(width int) *Input {
	i.Width = width
	return i
}

// SetInputType sets the type of input (text, number, etc.)
func (i *Input) SetInputType(inputType InputType) *Input {
	i.InputType = inputType
	return i
}

// SetPlaceholder sets the placeholder text
func (i *Input) SetPlaceholder(placeholder string) *Input {
	i.Placeholder = placeholder
	return i
}

// Focus gives focus to the input
func (i *Input) Focus() {
	i.focused = true
	i.Active = true
}

// Blur removes focus from the input
func (i *Input) Blur() {
	i.focused = false
	i.Active = false
	i.validate()
}

// SetValue sets the value of the input
func (i *Input) SetValue(value string) {
	i.Value = value
	i.cursor = len(value)
	i.validate()
}

// GetValue returns the current value
func (i *Input) GetValue() string {
	return i.Value
}

// IsValid returns true if the input is valid
func (i *Input) IsValid() bool {
	return i.ValidationState == ValidationValid || (i.ValidationState == ValidationNone && !i.Required)
}

// Update handles input events
func (i *Input) Update(msg tea.Msg) (*Input, tea.Cmd) {
	if !i.focused || i.Disabled {
		return i, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyBackspace:
			if i.cursor > 0 {
				i.Value = i.Value[:i.cursor-1] + i.Value[i.cursor:]
				i.cursor--
				i.validate()
			}
		case tea.KeyDelete:
			if i.cursor < len(i.Value) {
				i.Value = i.Value[:i.cursor] + i.Value[i.cursor+1:]
				i.validate()
			}
		case tea.KeyLeft:
			if i.cursor > 0 {
				i.cursor--
			}
		case tea.KeyRight:
			if i.cursor < len(i.Value) {
				i.cursor++
			}
		case tea.KeyHome:
			i.cursor = 0
		case tea.KeyEnd:
			i.cursor = len(i.Value)
		case tea.KeyRunes:
			if len(i.Value) < i.MaxLength {
				char := msg.String()
				
				// Input type validation
				if i.InputType == InputNumber && !isNumericChar(char) {
					return i, nil
				}
				
				i.Value = i.Value[:i.cursor] + char + i.Value[i.cursor:]
				i.cursor++
				i.validate()
			}
		}
	}

	// Handle cursor blinking
	i.blinkTick++
	if i.blinkTick > 30 {
		i.blinkTick = 0
	}

	return i, nil
}

// View renders the input component
func (i *Input) View() string {
	var parts []string

	// Label with consistent styling
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)
	
	if i.Required {
		// Keep label color consistent, just add asterisk
		parts = append(parts, labelStyle.Render(i.Label+" *:"))
	} else {
		parts = append(parts, labelStyle.Render(i.Label+":"))
	}

	// Input field
	inputValue := i.Value
	if inputValue == "" && !i.focused {
		inputValue = i.Placeholder
	}

	// Add cursor if focused
	if i.focused && i.blinkTick < 15 {
		if i.cursor >= len(i.Value) {
			inputValue = i.Value + "│"
		} else {
			inputValue = i.Value[:i.cursor] + "│" + i.Value[i.cursor:]
		}
	}

	// Determine border color based on state
	borderColor := "240" // Default gray
	if i.focused {
		borderColor = "39" // Blue when focused
	}
	if i.ValidationState == ValidationError {
		borderColor = "196" // Red for errors
	} else if i.ValidationState == ValidationValid {
		borderColor = "46" // Green for valid
	} else if i.ValidationState == ValidationWarning {
		borderColor = "226" // Yellow for warnings
	}

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Width(i.Width)

	if i.Disabled {
		inputStyle = inputStyle.Foreground(lipgloss.Color("240"))
	} else if inputValue == i.Placeholder {
		inputStyle = inputStyle.Foreground(lipgloss.Color("240"))
	}

	inputField := inputStyle.Render(inputValue)

	// Validation indicator and message
	var validationPart string
	switch i.ValidationState {
	case ValidationValid:
		validationPart = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Render(" ✓")
	case ValidationError:
		validationPart = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render(" ✗")
		if i.ErrorMessage != "" {
			validationPart += lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Render(" " + i.ErrorMessage)
		}
	case ValidationWarning:
		validationPart = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Render(" ⚠")
		if i.ErrorMessage != "" {
			validationPart += lipgloss.NewStyle().
				Foreground(lipgloss.Color("226")).
				Render(" " + i.ErrorMessage)
		}
	}

	// Combine all parts
	firstLine := lipgloss.JoinHorizontal(lipgloss.Left, parts[0], " ", inputField, validationPart)
	
	if i.ErrorMessage != "" && i.ValidationState == ValidationError {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Margin(0, 0, 0, len(i.Label)+2)
		
		return lipgloss.JoinVertical(lipgloss.Left, firstLine, errorStyle.Render(i.ErrorMessage))
	}

	return firstLine
}

// validate performs validation on the current value
func (i *Input) validate() {
	// Don't show validation errors while focused (live typing)
	if i.focused {
		i.ValidationState = ValidationNone
		i.ErrorMessage = ""
		return
	}

	if i.ValidationFunc == nil {
		if i.Required && strings.TrimSpace(i.Value) == "" {
			i.ValidationState = ValidationError
			i.ErrorMessage = "This field is required"
		} else if i.Value != "" {
			i.ValidationState = ValidationValid
			i.ErrorMessage = ""
		} else {
			i.ValidationState = ValidationNone
			i.ErrorMessage = ""
		}
		return
	}

	// If value is empty and field is required, show required error
	if strings.TrimSpace(i.Value) == "" {
		if i.Required {
			i.ValidationState = ValidationError
			i.ErrorMessage = "This field is required"
		} else {
			i.ValidationState = ValidationNone
			i.ErrorMessage = ""
		}
		return
	}

	// Run custom validation for non-empty values
	if err := i.ValidationFunc(i.Value); err != nil {
		i.ValidationState = ValidationError
		i.ErrorMessage = err.Error()
	} else {
		i.ValidationState = ValidationValid
		i.ErrorMessage = ""
	}
}

// isNumericChar checks if a character is numeric (including decimal point)
func isNumericChar(s string) bool {
	if len(s) != 1 {
		return false
	}
	char := s[0]
	return (char >= '0' && char <= '9') || char == '.'
}

// Common validation functions

// ValidateAmount validates monetary amounts
func ValidateAmount(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil // Allow empty for optional fields
	}
	
	// Remove any currency symbols
	value = strings.ReplaceAll(value, "$", "")
	value = strings.TrimSpace(value)
	
	if value == "" {
		return nil
	}
	
	// Check for valid decimal format
	if !isValidDecimal(value) {
		return errors.New("invalid amount format")
	}
	
	// Parse and validate amount
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return errors.New("invalid amount")
	}
	
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	
	if amount > 999999.99 {
		return errors.New("amount too large")
	}
	
	return nil
}

// ValidateDescription validates transaction descriptions
func ValidateDescription(value string) error {
	value = strings.TrimSpace(value)
	if len(value) > 200 {
		return errors.New("description too long (max 200 characters)")
	}
	return nil
}

// ValidateRequired validates required fields
func ValidateRequired(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("this field is required")
	}
	return nil
}

// Helper function to validate decimal format
func isValidDecimal(s string) bool {
	dotCount := 0
	for _, char := range s {
		if char == '.' {
			dotCount++
			if dotCount > 1 {
				return false
			}
		} else if char < '0' || char > '9' {
			return false
		}
	}
	return true
}