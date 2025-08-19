# CLAUDE.md - Go TUI Expense Tracker

This document outlines the architecture, tech stack, and implementation plan for a command-line expense and income tracker application written in Go. The application will feature a Terminal User Interface (TUI).

## 🏛️ 1. Architecture

We will use **Clean Architecture**. This separates the project into distinct layers, ensuring that the core business logic is independent of the UI, database, and other external frameworks. The dependency rule is strict: dependencies only point inwards, from outer layers to inner layers.

The layers are:
1.  **Entities:** The core domain objects (`Expense`, `Income`, `Category`). These have no dependencies.
2.  **Use Cases:** The application's business logic (e.g., `AddExpense`, `ListIncomes`). Depends only on Entities.
3.  **Interface Adapters:** The bridge between the core logic and the outside world. This includes our TUI components, which will translate user input into calls to the Use Cases.
4.  **Frameworks & Drivers:** The outermost layer, containing implementation details like the TUI library (`Bubble Tea`) and the database driver (`go-sqlite3`).

## 🛠️ 2. Tech Stack

-   **TUI Library:** **Bubble Tea**. A functional, Elm-inspired framework that is excellent for managing state in complex TUI applications. We'll use its companion libraries, `Bubbles` for UI components and `Lip Gloss` for styling.
-   **Database:** **SQLite**. A serverless, file-based database perfect for a local desktop application.
-   **Database Driver:** **`mattn/go-sqlite3`**. The most popular and robust SQLite driver for Go.
-   **CLI Argument Parsing:** **Cobra**. For potential future command-line subcommands (e.g., quick-add functionality).

## 📂 3. Project Structure

The project will follow a standard Go project layout that aligns with Clean Architecture principles.

expense-tracker/
├── cmd/
│   └── app/
│       └── main.go         # Application entry point
├── internal/
│   ├── core/
│   │   ├── domain/         # Entities (Transaction, Category structs) + tests
│   │   └── usecase/        # Use Cases and Repository Interfaces + tests
│   ├── handler/
│   │   └── tui/            # Bubble Tea models, views, and components
│   └── repository/
│       └── sqlite/         # SQLite implementation of the repository interfaces
├── test/
│   ├── integration/        # Integration tests for repository layer
│   ├── mocks/              # Auto-generated mocks from mockery
│   └── fixtures/           # Test data and fixtures
├── .mockery.yaml           # Mock generation configuration
├── Makefile               # Comprehensive development commands
├── go.mod
└── go.sum


## 🗄️ 4. Database Schema

We use a simplified, unified schema with two tables that handle both expenses and income efficiently.

**`categories` table:**
- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `name` TEXT NOT NULL
- `type` TEXT NOT NULL CHECK (type IN ('income', 'expense'))
- `UNIQUE(name, type)` - allows same category name for different types

**`transactions` table:**
- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `description` TEXT NOT NULL
- `amount` REAL NOT NULL
- `date` TEXT NOT NULL (RFC3339 format)
- `type` TEXT NOT NULL CHECK (type IN ('income', 'expense'))
- `category_id` INTEGER
- `FOREIGN KEY(category_id) REFERENCES categories(id)`

**Pre-populated Categories:**
- **Expense Categories:** Food & Dining, Transportation, Shopping, Entertainment, Bills & Utilities, Healthcare, Other
- **Income Categories:** Salary, Freelance, Investment, Gift, Other

## 🎨 5. UX/UI Design Guidelines

### 5.1 TUI Design Principles

**Simplicity First**
- Clear, uncluttered interface with logical grouping
- Progressive disclosure for complex operations
- Consistent navigation patterns throughout the application

**Accessibility**
- High contrast color schemes for better readability
- Support for screen readers through semantic structure
- Keyboard-only navigation with clear focus indicators
- Configurable color themes for visual impairments

**Responsive Design**
- Graceful handling of different terminal sizes (min 80x24)
- Collapsible/expandable sections for small screens
- Horizontal scrolling for wide tables when necessary

### 5.2 Visual Hierarchy and Colors

**Color Scheme (Lip Gloss styles):**
- **Primary**: Blue (#0066cc) for selected items and primary actions
- **Success**: Green (#22c55e) for income and positive balances
- **Warning**: Yellow (#eab308) for budget alerts and warnings  
- **Error**: Red (#ef4444) for expenses and error states
- **Neutral**: Gray (#6b7280) for secondary text and borders
- **Background**: Terminal default with subtle highlights

**Typography:**
- Bold for headers and important values
- Italic for labels and metadata
- Underline for focused interactive elements
- Consistent spacing and alignment

### 5.3 Navigation and Interaction Patterns

**Global Keybindings:**
- `q` or `Ctrl+C`: Quit application
- `?` or `h`: Show help/keybinding reference
- `Tab`/`Shift+Tab`: Navigate between sections
- `Enter`: Confirm/submit actions
- `Esc`: Cancel/go back

**Form Interactions:**
- Clear field labels and validation feedback
- Tab navigation between form fields
- Real-time validation with helpful error messages
- Autofocus on first field when forms open

**Table Navigation:**
- `j`/`k` or arrow keys for row navigation
- `Page Up`/`Page Down` for pagination
- `/` for search/filter mode
- `s` for sorting options

### 5.4 Error Handling and Feedback

**Error States:**
- Clear, actionable error messages
- Red color coding for critical errors
- Yellow for warnings and validation issues
- Inline validation for form fields

**Loading States:**
- Spinner animations for database operations
- Progress indicators for long-running tasks
- Graceful handling of slow operations

**Success Feedback:**
- Brief confirmation messages for completed actions
- Visual updates to reflect new data
- Smooth transitions between states

## 👥 6. User Stories

### 6.1 Core User Personas

**Sarah - Budget-Conscious Student**
- Needs to track every expense to stay within limited budget
- Values quick entry and simple categorization
- Uses application daily for small transactions

**Mike - Small Business Owner**
- Tracks both personal and business expenses
- Needs detailed reporting for tax purposes
- Values export capabilities and data accuracy

**Elena - Family Financial Manager**
- Manages household budget for family of four
- Needs to track multiple income sources
- Values analytics and spending trend insights

### 6.2 Detailed User Stories

#### Daily Usage Stories
```
As Sarah (student), I want to quickly add an expense while on-the-go
So that I can track my spending without interrupting my day
Given I have the application open
When I press 'a' and enter "Coffee $4.50"
Then the expense should be saved with today's date and default category
And I should return to the dashboard with updated balance
```

```
As Mike (business owner), I want to categorize my transactions properly
So that I can generate accurate reports for accounting
Given I'm adding a new expense
When I select the category field
Then I should see a list of existing categories plus option to create new
And I should be able to filter categories by typing
```

#### Analytics and Reporting Stories
```
As Elena (family manager), I want to see spending trends over time
So that I can identify areas where we're overspending
Given I have 3+ months of transaction data
When I navigate to the analytics view
Then I should see charts showing spending by category and month
And I should be able to compare current month to previous months
```

```
As Mike (business owner), I want to export my transaction data
So that I can import it into my accounting software
Given I have transactions for a specific date range
When I select export options
Then I should be able to choose CSV or JSON format
And specify date range and transaction types to include
```

#### Advanced User Stories
```
As a power user, I want to set up recurring transactions
So that I don't have to manually enter fixed monthly expenses
Given I have a monthly subscription expense
When I set it as recurring with frequency and amount
Then it should automatically appear each month for confirmation
```

```
As Elena (accessibility user), I want to use the app with screen reader
So that I can manage finances despite visual impairment
Given I'm using a screen reader
When I navigate through the application
Then all interface elements should be properly labeled and announced
And I should be able to complete all tasks using keyboard only
```

#### Error Handling Stories
```
As any user, I want clear feedback when something goes wrong
So that I can understand and fix issues quickly
Given the database becomes temporarily unavailable
When I try to add a transaction
Then I should see a clear error message explaining the problem
And suggestions for what to do next
```

## 🛠️ 7. Development Guidelines

### 7.1 Code Standards and Conventions

**Go Style Guidelines:**
- Follow official Go formatting (`gofmt`, `goimports`)
- Use meaningful variable and function names
- Keep functions small and focused (max 50 lines)
- Use interfaces for testability and flexibility
- Handle errors explicitly, never ignore them

**Clean Architecture Compliance:**
- Strict dependency direction (inward only)
- No direct database imports in use case layer
- All external dependencies injected via interfaces
- Domain entities should have no external dependencies

**Code Organization:**
```go
// Example use case structure
type AddExpenseUseCase struct {
    expenseRepo ExpenseRepository
    categoryRepo CategoryRepository
}

func (uc *AddExpenseUseCase) Execute(expense domain.Expense) error {
    // Validation
    if err := expense.Validate(); err != nil {
        return fmt.Errorf("invalid expense: %w", err)
    }
    
    // Business logic
    return uc.expenseRepo.Create(expense)
}
```

### 7.2 Git Workflow

**Branch Strategy:**
- `main`: Production-ready code
- `develop`: Integration branch for features
- `feature/*`: Individual feature development
- `hotfix/*`: Critical production fixes

**Commit Guidelines:**
- Use conventional commits format: `type(scope): description`
- Types: feat, fix, docs, test, refactor, style, chore
- Examples: `feat(tui): add expense filtering`, `fix(db): handle connection timeout`

**Pull Request Process:**
1. Create feature branch from `develop`
2. Implement feature with tests
3. Run full test suite and linting
4. Create PR with clear description and screenshots
5. Code review and approval required
6. Squash merge to maintain clean history

### 7.3 Performance Considerations

**Database Optimization:**
- Use prepared statements for repeated queries
- Implement connection pooling for concurrent operations
- Add database indexes for frequently queried columns
- Paginate large result sets

**TUI Performance:**
- Minimize re-renders by optimizing model updates
- Use efficient data structures for large transaction lists
- Implement virtual scrolling for very long lists
- Cache frequently accessed data

**Memory Management:**
- Close database connections properly
- Avoid memory leaks in long-running TUI sessions
- Profile memory usage during development
- Use appropriate data structures for scale

### 7.4 Security Best Practices

**Data Protection:**
- Store database file in user's home directory with restricted permissions
- Validate and sanitize all user input
- Use parameterized queries to prevent SQL injection
- Never log sensitive financial data

**Input Validation:**
```go
func (e *Expense) Validate() error {
    if e.Amount <= 0 {
        return errors.New("amount must be positive")
    }
    if len(e.Description) == 0 || len(e.Description) > 200 {
        return errors.New("description must be 1-200 characters")
    }
    return nil
}
```

### 7.5 Debugging and Troubleshooting

**Logging Strategy:**
- Use structured logging with levels (DEBUG, INFO, WARN, ERROR)
- Log to file in debug mode for troubleshooting
- Never log sensitive financial information
- Include correlation IDs for tracing user actions

**Debug Mode:**
```bash
# Enable debug logging
DEBUG=true ./expense-tracker

# Run with database logging
DB_DEBUG=true ./expense-tracker
```

**Common Issues:**
- **Database locked**: Handle SQLite concurrent access properly
- **Terminal size**: Test with various terminal dimensions
- **Unicode support**: Ensure proper UTF-8 handling for international users

## ✅ 8. Task Breakdown (Updated)

### Phase 1: Core Implementation ✅ COMPLETED
1.  **Project Setup:** ✅ DONE
    -   ✅ Initialize the Go module (`go mod init`)
    -   ✅ Create the directory structure following Clean Architecture
    -   ✅ Add dependencies (Bubble Tea, SQLite driver, etc.)

2.  **Core Logic (Inner Layers):** ✅ DONE
    -   ✅ Define the `Expense`, `Income`, and `Category` structs in `internal/core/domain/entities.go`
    -   ✅ Define repository interfaces in `internal/core/usecase/interfaces.go`
    -   ✅ Implement use cases (`AddExpenseUseCase`, `AddIncomeUseCase`, `SummaryUseCase`)

3.  **Database Implementation (Outer Layer):** ✅ DONE
    -   ✅ Implement repository interfaces in `internal/repository/sqlite/`
    -   ✅ Database initialization logic with all required tables
    -   ✅ CRUD operations for expenses, income, and categories

4.  **TUI Implementation (Outer Layer):** ✅ DONE
    -   ✅ Bubble Tea application structure in `internal/handler/tui/`
    -   ✅ Dashboard view with summary and recent transactions
    -   ✅ Add Expense and Add Income form views
    -   ✅ Transaction listing view with basic functionality

5.  **Integration:** ✅ DONE
    -   ✅ Wire everything together in `cmd/app/main.go`
    -   ✅ Database connection and repository initialization
    -   ✅ Use case dependency injection
    -   ✅ Main TUI model with complete user flows

### Phase 2: Testing & Quality ✅ COMPLETED
6.  **Test Infrastructure:** ✅ DONE
    -   ✅ Set up testing framework with testify and mockery
    -   ✅ Create test structure and mock generation with .mockery.yaml
    -   ✅ Add comprehensive Makefile with testing commands
    -   🔲 Add CI/CD pipeline with automated testing

7.  **Unit Testing:** ✅ DONE
    -   ✅ Domain layer tests (100% coverage achieved)
    -   ✅ Use case layer tests with mocked dependencies (67.9% coverage)
    -   ✅ Repository layer integration tests with real SQLite database
    -   🔲 TUI layer component tests (planned for Phase 3)

8.  **Code Quality:** ✅ DONE
    -   ✅ Add comprehensive input validation and error handling
    -   ✅ Implement proper error wrapping throughout application
    -   ✅ Code formatting, vetting, and quality checks in Makefile
    -   ✅ Security best practices implemented for data handling

### Phase 3: UX Enhancement 📋 PLANNED
9.  **User Experience:**
    -   🔲 Implement comprehensive keyboard shortcuts
    -   🔲 Add help system and onboarding
    -   🔲 Improve error messages and user feedback
    -   🔲 Accessibility improvements for screen readers

10. **Advanced TUI Features:**
    -   🔲 Search and filtering functionality
    -   🔲 Sorting options for transaction lists
    -   🔲 Pagination for large datasets
    -   🔲 Color themes and customization options

### Phase 4: Advanced Features 📋 PLANNED
11. **Analytics & Reporting:**
    -   🔲 Monthly/yearly spending trends
    -   🔲 Category-based analytics
    -   🔲 Budget tracking with alerts
    -   🔲 Spending pattern insights

12. **Data Management:**
    -   🔲 Export functionality (CSV, JSON)
    -   🔲 Data backup and restore
    -   🔲 Import from other applications
    -   🔲 Data migration tools

13. **Power User Features:**
    -   🔲 Recurring transaction support
    -   🔲 Multi-currency handling
    -   🔲 Advanced filtering and search
    -   🔲 Bulk operations for transactions

### Development Commands

Our comprehensive Makefile provides all necessary development commands:

```bash
# Basic Operations
make run              # Run the application
make run-dev          # Run with debug logging enabled
make build            # Build the application
make build-all        # Build for multiple platforms
make clean            # Clean generated files

# Testing
make test             # Run all tests
make test-unit        # Run only unit tests (domain & use case)
make test-integration # Run only integration tests
make test-coverage    # Run tests with coverage report
make test-coverage-html # Generate HTML coverage report
make test-race        # Run tests with race detection
make test-short       # Run tests with short flag

# Code Quality
make fmt              # Format code
make vet              # Vet code for issues
make lint             # Run linter (golangci-lint)
make check-all        # Run all quality checks (fmt, vet, test-coverage)

# Development Tools
make mocks            # Generate mocks from interfaces
make dev-deps         # Install development dependencies
make help             # Show all available commands
```

**Key Testing Infrastructure:**
- **Domain Layer:** 100% test coverage with comprehensive validation testing
- **Use Case Layer:** 67.9% coverage with mocked dependencies  
- **Integration Tests:** Full repository testing with real SQLite database
- **Mock Generation:** Automated with Mockery for clean, type-safe mocks
- **Quality Checks:** Integrated formatting, vetting, and linting

## 🧪 Testing Architecture Highlights

### Testing Strategy
Our testing follows the **Testing Pyramid** principle:
- **Unit Tests** (fast, numerous): Domain entities and use case business logic
- **Integration Tests** (moderate, focused): Repository layer with real database
- **End-to-End Tests** (slow, few): Planned for critical user workflows

### Coverage Achievements
```bash
$ make test-coverage
expense-tracker/internal/core/domain     100.0% coverage
expense-tracker/internal/core/usecase     67.9% coverage
expense-tracker/test/integration          [integration tests]
total:                                     6.1% overall
```

### Test Organization
- **Domain Tests** (`internal/core/domain/entities_test.go`): Comprehensive validation, edge cases, helper methods
- **Use Case Tests** (`internal/core/usecase/*_test.go`): Business logic with mocked repositories
- **Integration Tests** (`test/integration/*_test.go`): Repository operations with SQLite
- **Mocks** (`test/mocks/`): Auto-generated, type-safe mocks for interfaces

### Quality Assurance
- **Automated Formatting**: `go fmt` integration
- **Static Analysis**: `go vet` checks for common issues  
- **Mock Validation**: Testify mock expectations ensure correct repository usage
- **Error Handling**: Comprehensive error path testing

