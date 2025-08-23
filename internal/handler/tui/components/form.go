package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormField represents a field in a form
type FormField interface {
	Focus()
	Blur()
	Update(tea.Msg) (FormField, tea.Cmd)
	View() string
	IsValid() bool
}

// FormComponent interface for components that can be used in forms
type FormComponent interface {
	Focus()
	Blur()
	View() string
	IsValid() bool
}

// FormFieldWrapper wraps any form component to implement FormField interface
type FormFieldWrapper struct {
	component FormComponent
	updateFunc func(tea.Msg) (FormComponent, tea.Cmd)
}

// NewInputFieldWrapper creates a wrapper for Input components
func NewInputFieldWrapper(input *Input) *FormFieldWrapper {
	return &FormFieldWrapper{
		component: input,
		updateFunc: func(msg tea.Msg) (FormComponent, tea.Cmd) {
			updatedInput, cmd := input.Update(msg)
			return updatedInput, cmd
		},
	}
}

// NewDropdownFieldWrapper creates a wrapper for Dropdown components
func NewDropdownFieldWrapper(dropdown *Dropdown) *FormFieldWrapper {
	return &FormFieldWrapper{
		component: dropdown,
		updateFunc: func(msg tea.Msg) (FormComponent, tea.Cmd) {
			updatedDropdown, cmd := dropdown.Update(msg)
			return updatedDropdown, cmd
		},
	}
}

// NewToggleFieldWrapper creates a wrapper for Toggle components
func NewToggleFieldWrapper(toggle *Toggle) *FormFieldWrapper {
	return &FormFieldWrapper{
		component: toggle,
		updateFunc: func(msg tea.Msg) (FormComponent, tea.Cmd) {
			updatedToggle, cmd := toggle.Update(msg)
			return updatedToggle, cmd
		},
	}
}

// Focus implements FormField interface
func (f *FormFieldWrapper) Focus() {
	if f.component != nil {
		f.component.Focus()
	}
}

// Blur implements FormField interface
func (f *FormFieldWrapper) Blur() {
	if f.component != nil {
		f.component.Blur()
	}
}

// Update implements FormField interface
func (f *FormFieldWrapper) Update(msg tea.Msg) (FormField, tea.Cmd) {
	if f.component != nil && f.updateFunc != nil {
		updatedComponent, cmd := f.updateFunc(msg)
		f.component = updatedComponent
		return f, cmd
	}
	return f, nil
}

// View implements FormField interface
func (f *FormFieldWrapper) View() string {
	if f.component != nil {
		return f.component.View()
	}
	return ""
}

// IsValid implements FormField interface
func (f *FormFieldWrapper) IsValid() bool {
	if f.component != nil {
		return f.component.IsValid()
	}
	return true
}

// GetComponent returns the underlying component (for type-specific operations)
func (f *FormFieldWrapper) GetComponent() FormComponent {
	return f.component
}

// Form represents a form with multiple fields
type Form struct {
	Title       string
	Fields      []FormField
	CurrentField int
	Submitted   bool
	
	// Styling
	Width       int
	Height      int
	
	// Callbacks
	OnSubmit    func() tea.Cmd
	OnCancel    func() tea.Cmd
	OnFieldChange func(fieldIndex int) tea.Cmd
	
	// State
	showHelp    bool
	message     string
	messageType MessageType
}

// MessageType represents different types of form messages
type MessageType int

const (
	MessageInfo MessageType = iota
	MessageSuccess
	MessageError
	MessageWarning
)

// NewForm creates a new form
func NewForm(title string) *Form {
	return &Form{
		Title:        title,
		Fields:       make([]FormField, 0),
		CurrentField: 0,
		Width:        80,
		Height:       24,
	}
}

// AddField adds a field to the form
func (f *Form) AddField(field FormField) *Form {
	f.Fields = append(f.Fields, field)
	
	// Focus the first field
	if len(f.Fields) == 1 {
		field.Focus()
	}
	
	return f
}

// AddInput adds an input field to the form
func (f *Form) AddInput(input *Input) *Form {
	wrapper := NewInputFieldWrapper(input)
	return f.AddField(wrapper)
}

// AddDropdown adds a dropdown field to the form
func (f *Form) AddDropdown(dropdown *Dropdown) *Form {
	wrapper := NewDropdownFieldWrapper(dropdown)
	return f.AddField(wrapper)
}

// AddToggle adds a toggle field to the form
func (f *Form) AddToggle(toggle *Toggle) *Form {
	wrapper := NewToggleFieldWrapper(toggle)
	return f.AddField(wrapper)
}

// SetSize sets the form dimensions
func (f *Form) SetSize(width, height int) *Form {
	f.Width = width
	f.Height = height
	return f
}

// SetOnSubmit sets the submit callback
func (f *Form) SetOnSubmit(callback func() tea.Cmd) *Form {
	f.OnSubmit = callback
	return f
}

// SetOnCancel sets the cancel callback
func (f *Form) SetOnCancel(callback func() tea.Cmd) *Form {
	f.OnCancel = callback
	return f
}

// SetOnFieldChange sets the field change callback
func (f *Form) SetOnFieldChange(callback func(fieldIndex int) tea.Cmd) *Form {
	f.OnFieldChange = callback
	return f
}

// GetCurrentField returns the currently focused field
func (f *Form) GetCurrentField() FormField {
	if f.CurrentField >= 0 && f.CurrentField < len(f.Fields) {
		return f.Fields[f.CurrentField]
	}
	return nil
}

// GetField returns a field by index
func (f *Form) GetField(index int) FormField {
	if index >= 0 && index < len(f.Fields) {
		return f.Fields[index]
	}
	return nil
}

// GetInput returns an input field by index
func (f *Form) GetInput(index int) *Input {
	if field := f.GetField(index); field != nil {
		if wrapper, ok := field.(*FormFieldWrapper); ok {
			if input, ok := wrapper.GetComponent().(*Input); ok {
				return input
			}
		}
	}
	return nil
}

// GetDropdown returns a dropdown field by index
func (f *Form) GetDropdown(index int) *Dropdown {
	if field := f.GetField(index); field != nil {
		if wrapper, ok := field.(*FormFieldWrapper); ok {
			if dropdown, ok := wrapper.GetComponent().(*Dropdown); ok {
				return dropdown
			}
		}
	}
	return nil
}

// GetToggle returns a toggle field by index
func (f *Form) GetToggle(index int) *Toggle {
	if field := f.GetField(index); field != nil {
		if wrapper, ok := field.(*FormFieldWrapper); ok {
			if toggle, ok := wrapper.GetComponent().(*Toggle); ok {
				return toggle
			}
		}
	}
	return nil
}

// NextField moves to the next field
func (f *Form) NextField() tea.Cmd {
	if len(f.Fields) == 0 {
		return nil
	}
	
	// Blur current field
	if f.CurrentField >= 0 && f.CurrentField < len(f.Fields) {
		f.Fields[f.CurrentField].Blur()
	}
	
	// Move to next field
	f.CurrentField++
	if f.CurrentField >= len(f.Fields) {
		f.CurrentField = 0
	}
	
	// Focus new field
	f.Fields[f.CurrentField].Focus()
	
	// Call field change callback
	if f.OnFieldChange != nil {
		return f.OnFieldChange(f.CurrentField)
	}
	
	return nil
}

// PrevField moves to the previous field
func (f *Form) PrevField() tea.Cmd {
	if len(f.Fields) == 0 {
		return nil
	}
	
	// Blur current field
	if f.CurrentField >= 0 && f.CurrentField < len(f.Fields) {
		f.Fields[f.CurrentField].Blur()
	}
	
	// Move to previous field
	f.CurrentField--
	if f.CurrentField < 0 {
		f.CurrentField = len(f.Fields) - 1
	}
	
	// Focus new field
	f.Fields[f.CurrentField].Focus()
	
	// Call field change callback
	if f.OnFieldChange != nil {
		return f.OnFieldChange(f.CurrentField)
	}
	
	return nil
}

// IsValid returns true if all fields are valid
func (f *Form) IsValid() bool {
	for _, field := range f.Fields {
		if !field.IsValid() {
			return false
		}
	}
	return true
}

// GetInvalidFields returns a list of invalid field indices
func (f *Form) GetInvalidFields() []int {
	var invalid []int
	for i, field := range f.Fields {
		if !field.IsValid() {
			invalid = append(invalid, i)
		}
	}
	return invalid
}

// SetMessage sets a message to display in the form
func (f *Form) SetMessage(message string, messageType MessageType) {
	f.message = message
	f.messageType = messageType
}

// ClearMessage clears the current message
func (f *Form) ClearMessage() {
	f.message = ""
}

// Submit attempts to submit the form
func (f *Form) Submit() tea.Cmd {
	if !f.IsValid() {
		f.SetMessage("Please correct the errors below", MessageError)
		
		// Focus first invalid field
		invalidFields := f.GetInvalidFields()
		if len(invalidFields) > 0 {
			f.Fields[f.CurrentField].Blur()
			f.CurrentField = invalidFields[0]
			f.Fields[f.CurrentField].Focus()
		}
		
		return nil
	}
	
	f.Submitted = true
	
	if f.OnSubmit != nil {
		return f.OnSubmit()
	}
	
	return nil
}

// Cancel cancels the form
func (f *Form) Cancel() tea.Cmd {
	if f.OnCancel != nil {
		return f.OnCancel()
	}
	return nil
}

// Reset resets the form to its initial state
func (f *Form) Reset() {
	f.CurrentField = 0
	f.Submitted = false
	f.ClearMessage()
	
	// Reset all fields
	for i, field := range f.Fields {
		field.Blur()
		if i == 0 {
			field.Focus()
		}
	}
}

// Update handles form events
func (f *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			return f, f.NextField()
		case "shift+tab":
			return f, f.PrevField()
		case "ctrl+s":
			return f, f.Submit()
		case "esc":
			return f, f.Cancel()
		case "f1", "ctrl+h":
			f.showHelp = !f.showHelp
			return f, nil
		}
	}
	
	// Update current field
	if f.CurrentField >= 0 && f.CurrentField < len(f.Fields) {
		field, fieldCmd := f.Fields[f.CurrentField].Update(msg)
		f.Fields[f.CurrentField] = field
		if fieldCmd != nil {
			cmd = fieldCmd
		}
	}
	
	return f, cmd
}

// View renders the form
func (f *Form) View() string {
	var sections []string
	
	// Title with consistent emoji and styling
	if f.Title != "" {
		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true).
			Margin(0, 0, 1, 0)
		sections = append(sections, titleStyle.Render("💰 "+f.Title))
	}
	
	// Message with consistent styling
	if f.message != "" {
		messageStyle := lipgloss.NewStyle().
			Margin(0, 0, 1, 0).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder())
		
		switch f.messageType {
		case MessageSuccess:
			messageStyle = messageStyle.
				BorderForeground(lipgloss.Color("46")).
				Foreground(lipgloss.Color("46"))
			sections = append(sections, messageStyle.Render("✓ "+f.message))
		case MessageError:
			messageStyle = messageStyle.
				BorderForeground(lipgloss.Color("196")).
				Foreground(lipgloss.Color("196"))
			sections = append(sections, messageStyle.Render("✗ "+f.message))
		case MessageWarning:
			messageStyle = messageStyle.
				BorderForeground(lipgloss.Color("226")).
				Foreground(lipgloss.Color("226"))
			sections = append(sections, messageStyle.Render("⚠ "+f.message))
		default:
			messageStyle = messageStyle.
				BorderForeground(lipgloss.Color("39")).
				Foreground(lipgloss.Color("39"))
			sections = append(sections, messageStyle.Render("ℹ "+f.message))
		}
	}
	
	// Fields with spacing
	for i, field := range f.Fields {
		fieldView := field.View()
		if i < len(f.Fields)-1 {
			// Add spacing between fields except the last one
			fieldView += "\n"
		}
		sections = append(sections, fieldView)
	}
	
	// Buttons/Actions with improved styling
	buttonStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 2).
		Margin(1, 1, 0, 0).
		Bold(true)
	
	saveButton := buttonStyle.Copy().
		BorderForeground(lipgloss.Color("46")).
		Foreground(lipgloss.Color("46")).
		Render("💾 Save (Ctrl+S)")
	
	cancelButton := buttonStyle.Copy().
		BorderForeground(lipgloss.Color("240")).
		Foreground(lipgloss.Color("240")).
		Render("❌ Cancel (ESC)")
	
	buttons := lipgloss.JoinHorizontal(lipgloss.Left, saveButton, cancelButton)
	sections = append(sections, buttons)
	
	// Help with improved styling
	if f.showHelp {
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2).
			Margin(1, 0, 0, 0)
		
		helpText := []string{
			"📋 Navigation Help:",
			"",
			"• Tab/Shift+Tab - Navigate fields",
			"• Enter - Activate/Select options",
			"• Ctrl+S - Submit form",
			"• ESC - Cancel and return",
			"• F1 - Toggle this help panel",
		}
		
		sections = append(sections, helpStyle.Render(lipgloss.JoinVertical(lipgloss.Left, helpText...)))
	} else {
		helpHint := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Margin(1, 0, 0, 0).
			Render("💡 Press F1 for navigation help")
		sections = append(sections, helpHint)
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}