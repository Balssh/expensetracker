# Styling Guide - Expense Tracker TUI

This document defines the comprehensive styling system for the expense tracker terminal user interface, ensuring consistency across all components and screens.

## Global Design System

### Color Palette

#### Primary Colors
- **Primary Blue**: `#0066cc` - Selected items, primary actions, focused elements
- **Success Green**: `#22c55e` - Income, positive balances, success messages
- **Warning Yellow**: `#eab308` - Alerts, warnings, pending states
- **Error Red**: `#ef4444` - Expenses, negative balances, error states
- **Neutral Gray**: `#6b7280` - Secondary text, borders, inactive elements

#### Background & Text
- **Background**: Terminal default with subtle highlights
- **Text Primary**: `#FAFAFA` - Main content, headers
- **Text Secondary**: `#626262` - Help text, metadata
- **Text Muted**: `#404040` - Placeholder text, disabled states

### Typography Hierarchy

#### Headers
- **Application Title**: Bold, Primary Blue background, white text, padding 0,1
- **Section Headers**: Bold, larger text, appropriate color for context
- **Subsection Headers**: Bold, normal size, context-appropriate color

#### Content Text
- **Body Text**: Normal weight, primary text color
- **Labels**: Bold, slightly smaller, secondary text color
- **Values**: Bold for important numbers (amounts, balances)
- **Metadata**: Italic, secondary text color

### Border System

#### Border Styles
- **Main Application**: Rounded border, Primary Blue, padding 2,3
- **Panel Sections**: Rounded border, Neutral Gray, padding 1,2
- **Form Elements**: Rounded border, Primary Blue, padding 1,2
- **Focus Indicators**: Thick border, Primary Blue, enhanced padding

#### Spacing Standards
- **Inter-panel spacing**: 1 line
- **Form field spacing**: 1 line between fields
- **Content margins**: 2 spaces horizontal, 1 line vertical
- **Button padding**: 0,1 (vertical, horizontal)

## Component-Specific Styling

### Input Elements

#### Text Inputs
- **Default State**: Gray background `#333333`, white text, padding 0,1
- **Focused State**: Primary Blue background, white text, padding 0,1
- **Error State**: Red border, white text on dark background
- **Placeholder**: Muted gray text

#### Selection Lists
- **Unselected Items**: Normal text, no background
- **Selected Item**: Primary Blue background, white text, padding 0,1
- **Highlighted Item**: Primary Blue background with arrow indicator "→ "

### Tables

#### Table Headers
- **Style**: Bold, white text, Primary Blue background, padding 0,1
- **Alignment**: Left-aligned text, consistent column widths
- **Separator**: Horizontal line using "─" character

#### Table Rows
- **Alignment**: 
  - Date: Left-aligned, 12 characters
  - Category: Left-aligned, 15 characters  
  - Description: Left-aligned, 20 characters
  - Type: Left-aligned, 10 characters
  - Amount: Right-aligned, 10 characters
- **Color Coding**: 
  - Income rows: Success Green for amount
  - Expense rows: Error Red for amount

### Status Indicators

#### Messages
- **Success**: Success Green, bold text
- **Error**: Error Red, bold text  
- **Warning**: Warning Yellow, bold text
- **Info**: Neutral Gray, normal text

#### Loading States
- **Spinner**: Rotating character animation
- **Progress**: Horizontal bar with percentage
- **Placeholder**: "Loading..." in secondary text color

## Screen-Specific Styling

### Dashboard (Main Screen)

#### Overall Layout
- **Container**: Full-screen with application border
- **Centering**: Horizontally centered, minimum 80 characters wide
- **Vertical Sections**: 3 distinct panels with borders

#### Top Panel - Summary
- **Header**: Month/Year in bold, larger text
- **Metrics Layout**: 
  ```
  Income:   $X,XXX.XX  (Success Green)
  Expenses: $X,XXX.XX  (Error Red)  
  Balance:  $X,XXX.XX  (Green/Red based on value)
  ```
- **Visual Bar**: Horizontal expense breakdown (dummy data initially)
- **Border**: Rounded, Neutral Gray

#### Middle Panel - Recent Transactions
- **Table Style**: Clean alignment, alternating subtle backgrounds
- **Headers**: Primary Blue background, white text
- **Date Format**: "Jan 02" for recent, "Jan 02, 2006" for older
- **Amount Alignment**: Right-aligned with $ symbol
- **Border**: Rounded, Neutral Gray

#### Bottom Panel - Help
- **Content**: Context-aware keybindings
- **Format**: Key combinations in bold, descriptions in normal text
- **Layout**: Horizontal flow with " • " separators
- **Style**: Secondary text color, no background

### Add Expense/Income Forms

#### Form Container
- **Border**: Rounded, Primary Blue when active
- **Layout**: Vertical field stack with consistent spacing
- **Width**: Centered, fixed width for consistency

#### Field Layout
```
Label: [Input Field                    ]

Description: [Enter expense description  ]

Amount:      [0.00                      ]

Date:        [YYYY-MM-DD or today       ]

Category:    [Selected Category ▼       ]
```

#### Selection States
- **Navigation Mode**: Arrow indicators, no editing
- **Edit Mode**: Cursor visible, input accepting
- **Category Dropdown**: Vertical list with selection indicator

### Transaction List View

#### Search Section
- **Search Bar**: Full-width, prominent positioning
- **Active State**: Primary Blue border and background
- **Placeholder**: "Search transactions..." in muted text

#### Results Table
- **Full Headers**: Date, Description, Category, Type, Amount
- **Pagination**: Bottom status line with page indicators
- **Empty State**: Centered message with helpful text

#### Controls Section
- **Keybinding Help**: Bottom panel with current options
- **Status Information**: Page numbers, result counts

## Responsive Design

### Minimum Dimensions
- **Width**: 80 characters minimum
- **Height**: 24 lines minimum
- **Graceful Degradation**: Horizontal scroll for tables if needed

### Adaptive Elements
- **Panel Widths**: Scale with terminal width
- **Table Columns**: Truncate with "..." if necessary
- **Help Text**: Collapse to essential keys on narrow terminals

### Breakpoints
- **Narrow** (80-100 chars): Compact layout, abbreviated help
- **Standard** (100-120 chars): Full layout, complete help
- **Wide** (120+ chars): Expanded layout, additional spacing

## Implementation Notes

### Lipgloss Style Objects
- Create reusable style objects for each component type
- Use consistent naming: `componentStateStyle` (e.g., `inputFocusedStyle`)
- Group related styles in logical sections

### Color Consistency
- Define colors as constants at package level
- Use semantic names: `colorPrimary`, `colorSuccess`, etc.
- Apply colors through style objects, not inline

### Border and Layout Helpers
- Create utility functions for common layouts
- Implement centering calculations for different screen sizes
- Use consistent padding/margin calculations