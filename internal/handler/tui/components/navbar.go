package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NavTab represents a navigation tab
type NavTab struct {
	Key        string // Keyboard shortcut key
	Label      string // Display label
	ViewID     int    // View identifier  
	Icon       string // Icon/emoji for the tab
	Shortcut   string // Keyboard shortcut display
}

// NavBar represents a navigation bar component
type NavBar struct {
	Tabs           []NavTab
	ActiveTab      int
	Width          int
	ShowShortcuts  bool
	
	// Styling
	TabStyle       lipgloss.Style
	ActiveTabStyle lipgloss.Style
	BorderStyle    lipgloss.Style
}

// NewNavBar creates a new navigation bar
func NewNavBar() *NavBar {
	return &NavBar{
		Tabs:          make([]NavTab, 0),
		ActiveTab:     0,
		Width:         80,
		ShowShortcuts: true,
		
		// Default styles
		TabStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Padding(0, 2).
			Margin(0, 1, 0, 0),
		ActiveTabStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("39")).
			Bold(true).
			Padding(0, 2).
			Margin(0, 1, 0, 0),
		BorderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Border(lipgloss.NormalBorder(), false, false, true, false),
	}
}

// AddTab adds a tab to the navigation bar
func (n *NavBar) AddTab(tab NavTab) *NavBar {
	n.Tabs = append(n.Tabs, tab)
	return n
}

// SetTabs sets all tabs at once
func (n *NavBar) SetTabs(tabs []NavTab) *NavBar {
	n.Tabs = tabs
	return n
}

// SetActiveTab sets the active tab by index
func (n *NavBar) SetActiveTab(index int) *NavBar {
	if index >= 0 && index < len(n.Tabs) {
		n.ActiveTab = index
	}
	return n
}

// SetActiveByViewID sets the active tab by view ID
func (n *NavBar) SetActiveByViewID(viewID int) *NavBar {
	for i, tab := range n.Tabs {
		if tab.ViewID == viewID {
			n.ActiveTab = i
			break
		}
	}
	return n
}

// SetWidth sets the width of the navigation bar
func (n *NavBar) SetWidth(width int) *NavBar {
	n.Width = width
	return n
}

// SetShowShortcuts controls whether to show keyboard shortcuts
func (n *NavBar) SetShowShortcuts(show bool) *NavBar {
	n.ShowShortcuts = show
	return n
}

// GetActiveTab returns the currently active tab
func (n *NavBar) GetActiveTab() *NavTab {
	if n.ActiveTab >= 0 && n.ActiveTab < len(n.Tabs) {
		return &n.Tabs[n.ActiveTab]
	}
	return nil
}

// Update handles navigation bar events
func (n *NavBar) Update(msg tea.Msg) (*NavBar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle number key navigation
		switch msg.String() {
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			index := int(msg.String()[0]) - int('1') // Convert '1' to 0, '2' to 1, etc.
			if index >= 0 && index < len(n.Tabs) {
				n.ActiveTab = index
				// Return a command to navigate to the selected view
				return n, func() tea.Msg {
					return NavigationRequestMsg{ViewID: n.Tabs[index].ViewID}
				}
			}
		}
	}
	
	return n, nil
}

// View renders the navigation bar
func (n *NavBar) View() string {
	if len(n.Tabs) == 0 {
		return ""
	}
	
	var tabViews []string
	
	for i, tab := range n.Tabs {
		var tabText string
		
		// Build tab text with icon and label
		if tab.Icon != "" {
			tabText = tab.Icon + " " + tab.Label
		} else {
			tabText = tab.Label
		}
		
		// Add keyboard shortcut if enabled and available
		if n.ShowShortcuts && tab.Key != "" {
			tabText += " (" + tab.Key + ")"
		}
		
		// Apply styling based on active state
		if i == n.ActiveTab {
			tabViews = append(tabViews, n.ActiveTabStyle.Render(tabText))
		} else {
			tabViews = append(tabViews, n.TabStyle.Render(tabText))
		}
	}
	
	// Join tabs horizontally
	tabsContent := lipgloss.JoinHorizontal(lipgloss.Left, tabViews...)
	
	// Apply border and full width
	navContent := n.BorderStyle.Width(n.Width).Render(tabsContent)
	
	return navContent
}

// NavigationRequestMsg is sent when a tab is selected
type NavigationRequestMsg struct {
	ViewID int
}