# Navigation & Keybinding Reference - Expense Tracker TUI

This document defines the complete navigation system and keybinding scheme for the expense tracker, ensuring consistent user experience across all screens.

## Navigation Philosophy

### Core Principles
1. **Vim-Inspired Movement**: `h,j,k,l` + arrow keys for navigation
2. **Clear Selection States**: Navigate → Select → Edit workflow
3. **Consistent Actions**: Same keys perform same actions across screens
4. **Escape Sequences**: `Esc` always cancels/goes back
5. **Enter Confirms**: `Enter` always confirms/selects current item

### Navigation Flow
```
Browse Mode → Selection Mode → Edit Mode
     ↑              ↓              ↓
     └──────── Esc ←────── Esc ←──┘
```

## Global Keybindings

### Application-Wide Controls

| Key | Action | Context |
|-----|--------|---------|
| `q` | Quit application | Dashboard only |
| `Ctrl+C` | Force quit | Any screen |
| `Esc` | Cancel/Go back | Any screen (context-dependent) |
| `?` or `h` | Show help | Any screen |

### Universal Navigation

| Key | Action | Notes |
|-----|--------|-------|
| `↑` or `k` | Move up | Lists, menus, forms |
| `↓` or `j` | Move down | Lists, menus, forms |
| `←` or `h` | Move left | Forms, panels |
| `→` or `l` | Move right | Forms, panels |
| `Enter` | Select/Confirm | Current highlighted item |
| `Esc` | Cancel/Back | Exit current mode/screen |

## Screen-Specific Navigation

### Dashboard (Main Screen)

#### Quick Actions
| Key | Action | Description |
|-----|--------|-------------|
| `a` | Add Expense | Open add expense form |
| `i` | Add Income | Open add income form |
| `l` | List Transactions | View all transactions |
| `s` | Summary View | Toggle extended summary |
| `r` | Refresh | Reload data from database |

#### Panel Navigation
| Key | Action | Description |
|-----|--------|-------------|
| `Tab` | Next Panel | Cycle through: Summary → Transactions → Help |
| `Shift+Tab` | Previous Panel | Reverse cycle through panels |
| `1` | Summary Panel | Jump directly to summary |
| `2` | Transactions Panel | Jump directly to transactions |
| `3` | Help Panel | Jump directly to help |

#### Summary Panel Actions
| Key | Action | Description |
|-----|--------|-------------|
| `m` | Month View | Toggle month/year view |
| `c` | Category Breakdown | Show expense by category |
| `Enter` | Expand Details | Show detailed breakdown |

#### Recent Transactions Panel
| Key | Action | Description |
|-----|--------|-------------|
| `↑/↓` | Navigate Transactions | Highlight transaction |
| `Enter` | View Details | Show transaction details |
| `e` | Edit Transaction | Edit selected transaction |
| `d` | Delete Transaction | Delete selected transaction |
| `f` | Filter by Category | Filter by transaction category |

### Add Expense Form

#### Form Navigation
| Key | Action | Description |
|-----|--------|-------------|
| `↑` or `k` | Previous Field | Move to previous form field |
| `↓` or `j` | Next Field | Move to next form field |
| `Enter` | Edit Field | Enter edit mode for current field |
| `Tab` | Next Field | Alternative next field navigation |
| `Shift+Tab` | Previous Field | Alternative previous field navigation |

#### Field-Specific Actions
| Key | Action | Field | Description |
|-----|--------|-------|-------------|
| `Enter` | Edit | Any | Start editing current field |
| `Esc` | Stop Editing | Any | Exit edit mode, keep changes |
| `Ctrl+U` | Clear Field | Text inputs | Clear current field content |
| `Space` | Open Dropdown | Category | Open category selection |

#### Category Selection
| Key | Action | Description |
|-----|--------|-------------|
| `↑` or `k` | Previous Category | Move up in category list |
| `↓` or `j` | Next Category | Move down in category list |
| `Enter` | Select Category | Choose highlighted category |
| `Esc` | Cancel Selection | Return to form without changing |
| `/` | Search Categories | Start typing to filter |

#### Form Actions
| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+S` | Save | Submit the form |
| `Ctrl+Enter` | Save and New | Save and start new expense |
| `Esc` | Cancel | Return to dashboard without saving |
| `Ctrl+R` | Reset Form | Clear all fields |

### Add Income Form

*Uses identical navigation pattern as Add Expense Form*

#### Form Navigation
| Key | Action | Description |
|-----|--------|-------------|
| `↑/↓` or `k/j` | Navigate Fields | Move between form fields |
| `Enter` | Edit Field | Enter edit mode |
| `Space` | Category Dropdown | Open income category selection |
| `Ctrl+S` | Save Income | Submit form |
| `Esc` | Cancel | Return to dashboard |

### Transaction List View

#### List Navigation
| Key | Action | Description |
|-----|--------|-------------|
| `↑` or `k` | Previous Transaction | Move up in transaction list |
| `↓` or `j` | Next Transaction | Move down in transaction list |
| `Page Up` | Previous Page | Move up one page |
| `Page Down` | Next Page | Move down one page |
| `Home` | First Transaction | Jump to top of list |
| `End` | Last Transaction | Jump to bottom of list |

#### Search and Filter
| Key | Action | Description |
|-----|--------|-------------|
| `/` | Search Mode | Enter search/filter mode |
| `Enter` | Execute Search | Apply search filter |
| `Esc` | Clear Search | Exit search, clear filter |
| `c` | Clear All Filters | Remove all active filters |
| `f` | Filter Menu | Open filter options |

#### Transaction Actions
| Key | Action | Description |
|-----|--------|-------------|
| `Enter` | View Details | Show detailed transaction view |
| `e` | Edit Transaction | Edit selected transaction |
| `d` | Delete Transaction | Delete with confirmation |
| `x` | Toggle Selection | Multi-select for bulk actions |
| `Ctrl+A` | Select All | Select all visible transactions |

#### View Options
| Key | Action | Description |
|-----|--------|-------------|
| `s` | Sort Options | Open sort menu |
| `g` | Group by Category | Toggle category grouping |
| `t` | Toggle Type Filter | Show only income or expenses |
| `m` | Month Filter | Filter by specific month |

#### Bulk Actions (Multi-select Mode)
| Key | Action | Description |
|-----|--------|-------------|
| `Space` | Toggle Selection | Add/remove from selection |
| `Delete` | Bulk Delete | Delete all selected |
| `e` | Bulk Edit | Edit common fields |
| `Esc` | Clear Selection | Exit multi-select mode |

## Advanced Navigation Patterns

### Quick Jump Navigation
| Key Combination | Action | Description |
|------------------|--------|-------------|
| `g` + `g` | Go to Top | Jump to first item in any list |
| `G` | Go to Bottom | Jump to last item in any list |
| `Ctrl+Home` | First Page | Jump to first page |
| `Ctrl+End` | Last Page | Jump to last page |

### Modal and Dialog Navigation
| Key | Action | Description |
|-----|--------|-------------|
| `Tab` | Next Element | Navigate dialog elements |
| `Shift+Tab` | Previous Element | Reverse navigate elements |
| `Enter` | Confirm | Accept dialog action |
| `Esc` | Cancel | Close dialog without action |
| `Space` | Toggle | Toggle checkboxes/options |

### Search Mode Navigation
| Key | Action | Description |
|-----|--------|-------------|
| `↑` or `Ctrl+P` | Previous Match | Navigate search results |
| `↓` or `Ctrl+N` | Next Match | Navigate search results |
| `Enter` | Select Match | Choose current search result |
| `Esc` | Exit Search | Return to normal mode |
| `Ctrl+U` | Clear Search | Clear search term |

## Context-Sensitive Help

### Dynamic Help Display
Each screen shows relevant keybindings in the bottom help panel:

#### Dashboard Help
```
(a) Add Expense • (i) Add Income • (l) List All • (?) Help • (q) Quit
```

#### Form Help (Edit Mode)
```
↑/↓ Navigate • Enter Edit • Space Category • Ctrl+S Save • Esc Cancel
```

#### List View Help
```
↑/↓ Navigate • (/) Search • (e) Edit • (d) Delete • (q) Back
```

#### Search Mode Help
```
Type to search • Enter Apply • Esc Cancel • ↑/↓ Navigate results
```

## Accessibility Features

### Screen Reader Support
- All navigation states announced
- Form field labels clearly defined
- Table headers properly associated
- Status changes verbally indicated

### Keyboard-Only Operation
- No mouse required for any operation
- All functionality accessible via keyboard
- Clear visual focus indicators
- Logical tab order throughout

### Visual Indicators
- Clear highlight for current selection
- Distinct styles for different modes
- Color coding with text alternatives
- Progress indicators for operations

## Error Handling in Navigation

### Invalid Key Presses
- Ignored silently in most contexts
- Brief error flash for destructive actions
- Context help reminder for complex screens

### Navigation Boundaries
- List navigation stops at first/last item
- Form navigation cycles through fields
- Panel navigation wraps around

### Mode Conflicts
- Clear mode indicators always visible
- Escape always returns to safe state
- Consistent behavior across contexts