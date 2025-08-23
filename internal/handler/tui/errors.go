package tui

import (
	"fmt"
	"time"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity int

const (
	ErrorSeverityInfo ErrorSeverity = iota
	ErrorSeverityWarning
	ErrorSeverityError
	ErrorSeverityCritical
)

// ErrorCategory represents the category/domain of an error
type ErrorCategory string

const (
	CategoryValidation    ErrorCategory = "validation"
	CategoryDataLoad      ErrorCategory = "data_load"
	CategoryDataSave      ErrorCategory = "data_save"
	CategoryNetwork       ErrorCategory = "network"
	CategoryPermission    ErrorCategory = "permission"
	CategoryConfiguration ErrorCategory = "configuration"
	CategoryGeneral       ErrorCategory = "general"
)

// UserError represents a user-friendly error with context and recovery options
type UserError struct {
	Code         string        // Unique error code for logging/debugging
	Title        string        // Brief title for the error
	Message      string        // User-friendly error message
	Details      string        // Technical details (optional, for debugging)
	Category     ErrorCategory // Error category
	Severity     ErrorSeverity // Severity level
	Timestamp    time.Time     // When the error occurred
	Recoverable  bool          // Whether the user can recover from this error
	RetryAction  string        // Suggested retry action (if recoverable)
	HelpText     string        // Additional help text for the user
	OriginalErr  error         // Original error (for logging/debugging)
}

// NewUserError creates a new user error with basic information
func NewUserError(code, title, message string, category ErrorCategory) *UserError {
	return &UserError{
		Code:        code,
		Title:       title,
		Message:     message,
		Category:    category,
		Severity:    ErrorSeverityError,
		Timestamp:   time.Now(),
		Recoverable: true,
	}
}

// NewCriticalError creates a new critical error that requires immediate attention
func NewCriticalError(code, title, message string, category ErrorCategory, originalErr error) *UserError {
	return &UserError{
		Code:        code,
		Title:       title,
		Message:     message,
		Category:    category,
		Severity:    ErrorSeverityCritical,
		Timestamp:   time.Now(),
		Recoverable: false,
		OriginalErr: originalErr,
	}
}

// WithDetails adds technical details to the error
func (e *UserError) WithDetails(details string) *UserError {
	e.Details = details
	return e
}

// WithRetry marks the error as recoverable with a specific retry action
func (e *UserError) WithRetry(action string) *UserError {
	e.Recoverable = true
	e.RetryAction = action
	return e
}

// WithHelp adds help text for the user
func (e *UserError) WithHelp(helpText string) *UserError {
	e.HelpText = helpText
	return e
}

// WithOriginalError attaches the original error for debugging
func (e *UserError) WithOriginalError(err error) *UserError {
	e.OriginalErr = err
	if e.Details == "" && err != nil {
		e.Details = err.Error()
	}
	return e
}

// Error implements the error interface
func (e *UserError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Title, e.Message)
}

// GetSeverityColor returns the color code for the error severity
func (e *UserError) GetSeverityColor() string {
	switch e.Severity {
	case ErrorSeverityInfo:
		return "39"  // Blue
	case ErrorSeverityWarning:
		return "226" // Yellow
	case ErrorSeverityError:
		return "196" // Red
	case ErrorSeverityCritical:
		return "160" // Dark red
	default:
		return "196"
	}
}

// GetSeveritySymbol returns the symbol for the error severity
func (e *UserError) GetSeveritySymbol() string {
	switch e.Severity {
	case ErrorSeverityInfo:
		return "ℹ️"
	case ErrorSeverityWarning:
		return "⚠️"
	case ErrorSeverityError:
		return "❌"
	case ErrorSeverityCritical:
		return "🔥"
	default:
		return "❌"
	}
}

// ErrorHandler manages error display and recovery
type ErrorHandler struct {
	errors       []*UserError
	maxErrors    int
	showDetails  bool
}

// NewErrorHandler creates a new error handler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		errors:      make([]*UserError, 0),
		maxErrors:   10, // Keep last 10 errors
		showDetails: false,
	}
}

// AddError adds a new error to the handler
func (h *ErrorHandler) AddError(err *UserError) {
	h.errors = append(h.errors, err)
	
	// Keep only the most recent errors
	if len(h.errors) > h.maxErrors {
		h.errors = h.errors[len(h.errors)-h.maxErrors:]
	}
}

// GetLatestError returns the most recent error
func (h *ErrorHandler) GetLatestError() *UserError {
	if len(h.errors) == 0 {
		return nil
	}
	return h.errors[len(h.errors)-1]
}

// GetErrorsByCategory returns errors filtered by category
func (h *ErrorHandler) GetErrorsByCategory(category ErrorCategory) []*UserError {
	var filtered []*UserError
	for _, err := range h.errors {
		if err.Category == category {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

// ClearErrors clears all stored errors
func (h *ErrorHandler) ClearErrors() {
	h.errors = h.errors[:0]
}

// HasCriticalErrors checks if there are any critical errors
func (h *ErrorHandler) HasCriticalErrors() bool {
	for _, err := range h.errors {
		if err.Severity == ErrorSeverityCritical {
			return true
		}
	}
	return false
}

// SetShowDetails controls whether technical details are shown
func (h *ErrorHandler) SetShowDetails(show bool) {
	h.showDetails = show
}

// UnifiedErrorMsg is the unified error message type for Bubble Tea
type UnifiedErrorMsg struct {
	Error *UserError
}

// UnifiedMessageMsg is the unified success/info message type for Bubble Tea  
type UnifiedMessageMsg struct {
	Title       string
	Message     string
	Severity    ErrorSeverity
	Dismissible bool
	AutoDismiss time.Duration
	Timestamp   time.Time
}

// ToastNotification represents a toast notification
type ToastNotification struct {
	ID          string
	Title       string
	Message     string
	Severity    ErrorSeverity
	Icon        string
	Timestamp   time.Time
	AutoDismiss time.Duration
	Dismissed   bool
	Progress    float64 // Progress towards auto-dismiss (0.0 to 1.0)
}

// NewToastNotification creates a new toast notification
func NewToastNotification(title, message string, severity ErrorSeverity, autoDismiss time.Duration) *ToastNotification {
	return &ToastNotification{
		ID:          fmt.Sprintf("toast_%d", time.Now().UnixNano()),
		Title:       title,
		Message:     message,
		Severity:    severity,
		Icon:        getSeverityIcon(severity),
		Timestamp:   time.Now(),
		AutoDismiss: autoDismiss,
		Dismissed:   false,
		Progress:    0.0,
	}
}

// getSeverityIcon returns the appropriate icon for a severity level
func getSeverityIcon(severity ErrorSeverity) string {
	switch severity {
	case ErrorSeverityInfo:
		return "ℹ️"
	case ErrorSeverityWarning:
		return "⚠️"
	case ErrorSeverityError:
		return "❌"
	case ErrorSeverityCritical:
		return "🔥"
	default:
		return "✅" // Success/default
	}
}

// IsExpired checks if the toast should be auto-dismissed
func (t *ToastNotification) IsExpired() bool {
	if t.AutoDismiss <= 0 {
		return false // Never expires
	}
	return time.Since(t.Timestamp) >= t.AutoDismiss
}

// UpdateProgress updates the dismissal progress
func (t *ToastNotification) UpdateProgress() {
	if t.AutoDismiss <= 0 {
		t.Progress = 0.0
		return
	}
	
	elapsed := time.Since(t.Timestamp)
	t.Progress = float64(elapsed) / float64(t.AutoDismiss)
	if t.Progress > 1.0 {
		t.Progress = 1.0
		t.Dismissed = true
	}
}

// GetSeverityColor returns the color for the toast severity
func (t *ToastNotification) GetSeverityColor() string {
	switch t.Severity {
	case ErrorSeverityInfo:
		return "39"  // Blue
	case ErrorSeverityWarning:
		return "226" // Yellow
	case ErrorSeverityError:
		return "196" // Red
	case ErrorSeverityCritical:
		return "160" // Dark red
	default:
		return "46" // Green (success)
	}
}

// ToastManager manages toast notifications
type ToastManager struct {
	toasts    []*ToastNotification
	maxToasts int
}

// NewToastManager creates a new toast manager
func NewToastManager() *ToastManager {
	return &ToastManager{
		toasts:    make([]*ToastNotification, 0),
		maxToasts: 3, // Show max 3 toasts at once
	}
}

// AddToast adds a new toast notification
func (tm *ToastManager) AddToast(toast *ToastNotification) {
	tm.toasts = append(tm.toasts, toast)
	
	// Keep only the most recent toasts
	if len(tm.toasts) > tm.maxToasts {
		tm.toasts = tm.toasts[len(tm.toasts)-tm.maxToasts:]
	}
}

// AddSuccessToast adds a success toast with auto-dismiss
func (tm *ToastManager) AddSuccessToast(title, message string) {
	toast := NewToastNotification(title, message, ErrorSeverityInfo, 3*time.Second)
	toast.Icon = "✅"
	tm.AddToast(toast)
}

// AddWarningToast adds a warning toast with longer auto-dismiss
func (tm *ToastManager) AddWarningToast(title, message string) {
	toast := NewToastNotification(title, message, ErrorSeverityWarning, 5*time.Second)
	tm.AddToast(toast)
}

// AddErrorToast adds an error toast that requires manual dismissal
func (tm *ToastManager) AddErrorToast(title, message string) {
	toast := NewToastNotification(title, message, ErrorSeverityError, 0) // No auto-dismiss
	tm.AddToast(toast)
}

// UpdateToasts updates all toast notifications (call this periodically)
func (tm *ToastManager) UpdateToasts() {
	// Update progress and remove expired toasts
	activeToasts := make([]*ToastNotification, 0)
	
	for _, toast := range tm.toasts {
		toast.UpdateProgress()
		if !toast.Dismissed && !toast.IsExpired() {
			activeToasts = append(activeToasts, toast)
		}
	}
	
	tm.toasts = activeToasts
}

// GetActiveToasts returns all active (non-dismissed) toasts
func (tm *ToastManager) GetActiveToasts() []*ToastNotification {
	activeToasts := make([]*ToastNotification, 0)
	for _, toast := range tm.toasts {
		if !toast.Dismissed {
			activeToasts = append(activeToasts, toast)
		}
	}
	return activeToasts
}

// DismissToast dismisses a toast by ID
func (tm *ToastManager) DismissToast(id string) {
	for _, toast := range tm.toasts {
		if toast.ID == id {
			toast.Dismissed = true
			break
		}
	}
}

// DismissAll dismisses all active toasts
func (tm *ToastManager) DismissAll() {
	for _, toast := range tm.toasts {
		toast.Dismissed = true
	}
}

// HasActiveToasts returns true if there are any active toasts
func (tm *ToastManager) HasActiveToasts() bool {
	for _, toast := range tm.toasts {
		if !toast.Dismissed {
			return true
		}
	}
	return false
}

// Common error constructors for frequent scenarios

// ValidationError creates a validation error
func ValidationError(field, message string) *UserError {
	return NewUserError(
		fmt.Sprintf("VALIDATION_%s", field),
		"Validation Error",
		message,
		CategoryValidation,
	).WithRetry("Please correct the input and try again")
}

// DataLoadError creates a data loading error
func DataLoadError(resource, message string) *UserError {
	return NewUserError(
		fmt.Sprintf("DATA_LOAD_%s", resource),
		"Failed to Load Data",
		message,
		CategoryDataLoad,
	).WithRetry("Press 'r' to retry loading")
}

// DataSaveError creates a data saving error
func DataSaveError(resource, message string) *UserError {
	return NewUserError(
		fmt.Sprintf("DATA_SAVE_%s", resource),
		"Failed to Save Data",
		message,
		CategoryDataSave,
	).WithRetry("Please try saving again or check your input")
}

// NetworkError creates a network-related error
func NetworkError(operation, message string) *UserError {
	return NewUserError(
		fmt.Sprintf("NETWORK_%s", operation),
		"Connection Error",
		message,
		CategoryNetwork,
	).WithRetry("Please check your connection and try again").
		WithHelp("Ensure you have a stable internet connection")
}

// PermissionError creates a permission-related error
func PermissionError(resource, message string) *UserError {
	return NewUserError(
		fmt.Sprintf("PERMISSION_%s", resource),
		"Permission Denied",
		message,
		CategoryPermission,
	).WithHelp("Please check file permissions or contact an administrator")
}

// DatabaseError creates a database-related error
func DatabaseError(operation, message string, originalErr error) *UserError {
	return NewCriticalError(
		fmt.Sprintf("DB_%s", operation),
		"Database Error",
		message,
		CategoryDataSave,
		originalErr,
	).WithHelp("The application database may be corrupted or inaccessible")
}