# TODO_MVP.md - Expense Tracker MVP Implementation

## 🎯 PROJECT STATUS

### ✅ PHASE 1 COMPLETE - Core Foundation Ready!
**Clean Architecture Implementation**: Domain layer, use cases, repositories, and basic TUI are fully functional.

**What's Working:**
- 🏗️ Complete Clean Architecture setup
- 💾 SQLite database with auto-initialization
- 📊 Dashboard with monthly summaries and recent transactions  
- 🔄 Full navigation system with keyboard shortcuts
- 📱 Responsive TUI interface with Bubble Tea
- 🏷️ Default categories seeded automatically

**Current State:** Application builds and runs successfully with a functional dashboard view.

**Next Priority:** Implement Add Transaction form in Phase 6 for full transaction management.

## Phase 1: Foundation (Core Setup) ✅ COMPLETED
- [x] Initialize Go module (`go mod init github.com/yourusername/expense-tracker`)
- [x] Create project directory structure (Clean Architecture layout)
- [ ] Set up `.gitignore` for Go projects
- [ ] Add README.md with basic project description
- [ ] Set up golangci-lint configuration
- [x] Install core dependencies:
  - [x] Bubble Tea (`go get github.com/charmbracelet/bubbletea`)
  - [x] Lipgloss (`go get github.com/charmbracelet/lipgloss`)
  - [x] SQLite driver (`go get github.com/mattn/go-sqlite3`)
  - [ ] Migration tool (`go get -u github.com/golang-migrate/migrate/v4/cmd/migrate`)
  - [x] UUID library (`go get github.com/google/uuid`)

## Phase 2: Domain Layer ✅ COMPLETED
- [x] Create `internal/core/domain/entities.go` (Clean Architecture):
  - [x] Define Transaction struct with ID, Type, Amount, Category, Description, Date
  - [x] Define TransactionType enum (Income, Expense)
  - [x] Add validation methods
  - [x] Define Category struct
  - [x] Define default expense categories (Food & Dining, Transportation, etc.)
  - [x] Define default income categories (Salary, Gifts, Investments, Other)
  - [x] Add methods to validate category for transaction type
  - [x] Add helper methods (FormatAmount, FormatDate, IsIncome, IsExpense)
- [x] Repository error types defined in `internal/core/usecase/interfaces.go`
- [ ] Write unit tests for models

## Phase 3: Storage Layer ✅ COMPLETED
- [x] Create `internal/core/usecase/interfaces.go`:
  - [x] Define TransactionRepository interface
  - [x] Define CategoryRepository interface
  - [x] Define common errors (ErrNotFound, ErrDuplicateCategory, etc.)
- [x] Create `internal/repository/sqlite/repository.go`:
  - [x] Implement SQLiteRepository struct
  - [x] Implement connection management
  - [x] Automatic schema initialization
- [x] Create `internal/repository/sqlite/transaction.go`:
  - [x] Implement transaction CRUD operations:
    - [x] Create(transaction *domain.Transaction) error
    - [x] GetByID(id string) (*domain.Transaction, error)
    - [x] List(limit, offset int) ([]*domain.Transaction, error)
    - [x] ListByType(transactionType domain.TransactionType, limit, offset int) ([]*domain.Transaction, error)
    - [x] ListByDateRange(start, end time.Time) ([]*domain.Transaction, error)
    - [x] ListByCategory(categoryID int) ([]*domain.Transaction, error)
    - [x] Update(transaction *domain.Transaction) error
    - [x] Delete(id string) error
    - [x] Count() (int, error)
    - [x] CountByType(transactionType domain.TransactionType) (int, error)
- [x] Create `internal/repository/sqlite/category.go`:
  - [x] Implement category CRUD operations:
    - [x] Create(category *domain.Category) error
    - [x] GetByID(id int) (*domain.Category, error)
    - [x] GetByName(name string, categoryType domain.TransactionType) (*domain.Category, error)
    - [x] ListByType(categoryType domain.TransactionType) ([]*domain.Category, error)
    - [x] List() ([]*domain.Category, error)
    - [x] Update(category *domain.Category) error
    - [x] Delete(id int) error
    - [x] GetUsageCount(id int) (int, error)
    - [x] Exists(name string, categoryType domain.TransactionType) (bool, error)
- [x] Database schema with automatic initialization:
  - [x] Categories table with constraints
  - [x] Transactions table with foreign keys
  - [x] Performance indexes
- [x] Default category seeding implemented in use cases
- [ ] Write integration tests with test database

## Phase 4: Use Case Layer (Business Logic) ✅ COMPLETED
- [x] Create `internal/core/usecase/transaction.go`:
  - [x] Implement TransactionUseCase struct
  - [x] Add business logic methods:
    - [x] AddTransaction(transactionType domain.TransactionType, amount float64, categoryID int, description string, date time.Time) (*domain.Transaction, error)
    - [x] GetTransaction(id string) (*domain.Transaction, error)
    - [x] ListTransactions(limit, offset int) ([]*domain.Transaction, error)
    - [x] ListTransactionsByType(transactionType domain.TransactionType, limit, offset int) ([]*domain.Transaction, error)
    - [x] GetRecentTransactions(days int) ([]*domain.Transaction, error)
    - [x] GetMonthlyBalance(year, month int) (income, expense, balance float64, error)
    - [x] GetCategoryBreakdown(transactionType domain.TransactionType, year, month int) (map[string]float64, error)
    - [x] UpdateTransaction(transaction *domain.Transaction) error
    - [x] DeleteTransaction(id string) error
- [x] Create `internal/core/usecase/category.go`:
  - [x] Implement CategoryUseCase struct
  - [x] Add category management methods:
    - [x] GetCategories(categoryType domain.TransactionType) ([]*domain.Category, error)
    - [x] GetAllCategories() ([]*domain.Category, error)
    - [x] AddCategory(name string, categoryType domain.TransactionType) (*domain.Category, error)
    - [x] UpdateCategory(id int, name string) error
    - [x] DeleteCategory(id int) error
    - [x] CanDeleteCategory(id int) (bool, string, error)
    - [x] GetCategoryUsage(id int) (int, error)
    - [x] InitializeDefaultCategories() error
- [x] Add comprehensive input validation and business rules
- [ ] Write unit tests with mock repositories

## Phase 5: TUI Foundation ✅ COMPLETED
- [x] Create `internal/handler/tui/app.go`:
  - [x] Define main App model struct with state management
  - [x] Implement tea.Model interface (Init, Update, View)
  - [x] Set up navigation between views
  - [x] Add message/error handling system
  - [x] Responsive window size handling
- [x] Create `internal/handler/tui/dashboard.go`:
  - [x] Dashboard view with monthly summary
  - [x] Recent transactions display
  - [x] Navigation options
  - [x] Loading and error states
- [x] Basic styling integrated into components:
  - [x] Color scheme for income (green) and expense (red)
  - [x] Consistent styling with Lipgloss
  - [x] Header/footer layout
  - [x] Message display system
- [x] Keyboard shortcuts implemented:
  - [x] Global navigation (1-4, ?, q, ESC)
  - [x] Context-aware help text
  - [x] View switching

## Phase 6: TUI Views - Core 🚧 IN PROGRESS
- [x] Dashboard implemented (main menu functionality):
  - [x] Monthly summary with income/expense/balance
  - [x] Recent transactions display
  - [x] Navigation options and keyboard shortcuts
  - [x] Quick stats display (current month balance)
- [x] Create `internal/handler/tui/stubs.go`:
  - [x] Add Transaction view stub (ready for implementation)
  - [x] List Transactions view stub (ready for implementation)
  - [x] Categories view stub (ready for implementation)
  - [x] Help view with keyboard shortcuts
- [ ] Implement `internal/handler/tui/add_transaction.go`:
  - [ ] Toggle between Income/Expense mode
  - [ ] Form inputs for amount, category, description, date
  - [ ] Dynamic category list based on transaction type
  - [ ] Input validation with real-time feedback
  - [ ] Success confirmation with transaction details
  - [ ] Quick-add another option
- [ ] Implement `internal/handler/tui/list_transactions.go`:
  - [ ] Table view with color coding (green for income, red for expense)
  - [ ] Filters for type (All/Income/Expense)
  - [ ] Date range filter
  - [ ] Pagination controls
  - [ ] Sort options (date, amount, type)
  - [ ] Quick actions (edit, delete) with confirmation
  - [ ] Running balance column

## Phase 7: TUI Views - Management
- [ ] Create `internal/ui/views/edit_transaction.go`:
  - [ ] Pre-populate form with existing transaction data
  - [ ] Field-level editing
  - [ ] Validation before save
  - [ ] Cancel option with confirmation if changes made
- [ ] Create `internal/ui/views/category_manager.go`:
  - [ ] List categories by type (Income/Expense tabs)
  - [ ] Add new custom category
  - [ ] Edit category name
  - [ ] Delete category (with usage check)
  - [ ] Show usage count for each category
  - [ ] Highlight default vs custom categories
- [ ] Create `internal/ui/views/summary.go`:
  - [ ] Monthly income vs expense comparison
  - [ ] Net balance display
  - [ ] Category breakdown for income and expenses
  - [ ] Simple bar chart using box-drawing characters
  - [ ] Month/Year selector
  - [ ] YTD (Year-to-date) summary option

## Phase 8: TUI Components
- [ ] Create `internal/ui/components/input.go`:
  - [ ] Text input with validation
  - [ ] Number input with formatting
  - [ ] Date picker component
  - [ ] Currency formatting
- [ ] Create `internal/ui/components/select.go`:
  - [ ] Dropdown-like selection component
  - [ ] Search/filter capability
  - [ ] Custom vs default item styling
- [ ] Create `internal/ui/components/table.go`:
  - [ ] Configurable columns
  - [ ] Row selection and highlighting
  - [ ] Scrollable with fixed headers
  - [ ] Color coding support
- [ ] Create `internal/ui/components/tabs.go`:
  - [ ] Tab navigation component
  - [ ] Active tab highlighting
- [ ] Create `internal/ui/components/confirm.go`:
  - [ ] Confirmation dialog component
  - [ ] Yes/No options with clear defaults
- [ ] Create `internal/ui/components/message.go`:
  - [ ] Success/error/info message display
  - [ ] Auto-dismiss timer
  - [ ] Stack multiple messages

## Phase 9: Application Entry Point ✅ COMPLETED
- [x] Create `cmd/app/main.go`:
  - [x] Database path configuration (user home directory)
  - [x] Initialize database connection with SQLite
  - [x] Automatic schema initialization
  - [x] Seed default categories on first run
  - [x] Initialize repository and use case layers with dependency injection
  - [x] Start TUI application with Bubble Tea
  - [x] Handle graceful shutdown
  - [ ] Add panic recovery
- [ ] Add logging setup with levels
- [ ] Create database backup on startup (optional)

## Phase 10: Configuration
- [ ] Create `internal/config/config.go`:
  - [ ] Define Config struct
  - [ ] Load from environment variables
  - [ ] Load from config file (optional)
  - [ ] Set sensible defaults
- [ ] Add config for:
  - [ ] Database path
  - [ ] Currency symbol
  - [ ] Date format preference
  - [ ] Decimal places for amounts
  - [ ] Default transaction type (income/expense)
  - [ ] Color scheme preference
  - [ ] Backup frequency

## Phase 11: Testing & Polish
- [ ] Write comprehensive tests:
  - [ ] Model validation tests
  - [ ] Service layer business logic tests
  - [ ] Repository integration tests
  - [ ] Category management tests
  - [ ] Transaction filtering and sorting tests
  - [ ] Balance calculation tests
- [ ] Add data fixtures for testing
- [ ] Manual testing checklist:
  - [ ] Add income flow
  - [ ] Add expense flow
  - [ ] Edit transaction flow
  - [ ] Delete transaction with confirmation
  - [ ] Switch between income/expense in list view
  - [ ] Category management (add, edit, delete)
  - [ ] Monthly summary calculations
  - [ ] Date range filtering
  - [ ] Edge cases (empty states, invalid inputs, boundary values)
- [ ] Performance testing with large datasets (1000+ transactions)

## Phase 12: Documentation
- [ ] Write comprehensive README:
  - [ ] Installation instructions
  - [ ] Usage guide with examples
  - [ ] Screenshots/demo GIF
  - [ ] Architecture overview
  - [ ] Configuration options
- [ ] Add inline code documentation
- [ ] Create CONTRIBUTING.md
- [ ] Document keyboard shortcuts in-app and README
- [ ] Add example configuration file
- [ ] Create quick start guide

## Phase 13: Final MVP Polish
- [ ] Error handling review:
  - [ ] User-friendly error messages
  - [ ] Proper error propagation
  - [ ] Recovery from database errors
  - [ ] Network timeout handling (for future features)
- [ ] UX improvements:
  - [ ] Loading states for data fetching
  - [ ] Empty states with helpful messages
  - [ ] Confirmation dialogs for destructive actions
  - [ ] Success feedback for all actions
  - [ ] Smooth transitions between views
  - [ ] Consistent navigation patterns
- [ ] Data validation:
  - [ ] Prevent negative amounts (unless intentional)
  - [ ] Validate date ranges
  - [ ] Category existence checks
  - [ ] Description length limits
- [ ] Code cleanup:
  - [ ] Remove debug statements
  - [ ] Run linter and fix all issues
  - [ ] Refactor duplicated code
  - [ ] Ensure consistent naming
  - [ ] Optimize database queries

## Stretch Goals (Post-MVP)
- [ ] Export to CSV/JSON functionality
- [ ] Advanced search with multiple filters
- [ ] Transaction templates for recurring items
- [ ] Budget limits with visual warnings
- [ ] Backup/restore commands
- [ ] Multi-currency support
- [ ] Tags for transactions
- [ ] Monthly/yearly trends visualization
- [ ] Transaction attachments (receipt notes)
- [ ] Savings goals tracking
- [ ] Bill reminders

## Definition of Done for MVP
- [ ] Can add both income and expense transactions
- [ ] Can view list of transactions with type filtering
- [ ] Can edit existing transactions
- [ ] Can delete transactions with confirmation
- [ ] Can manage categories (add custom, edit, delete with usage check)
- [ ] Can view monthly summary with income vs expense
- [ ] Can see category breakdown for both types
- [ ] Balance calculations are accurate
- [ ] All core features have keyboard shortcuts
- [ ] Data persists between sessions
- [ ] Handles errors gracefully without crashes
- [ ] Has at least 70% test coverage on service layer
- [ ] Documentation is complete and helpful
- [ ] Default categories are properly seeded

## Notes
- Start with Phase 1-4 to get a working backend
- Test transaction type handling thoroughly before UI
- Keep income/expense UI clearly distinguished (colors, labels)
- Ensure category management doesn't break existing transactions
- Regular manual testing after each phase
- Consider using sample data for development
- Keep commits atomic and well-described