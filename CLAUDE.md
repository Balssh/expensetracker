# CLAUDE.md - Expense Tracker TUI Project

## Project Overview
Building a Terminal User Interface (TUI) expense tracker application in Go to learn the language and good coding practices. The app will help users track their daily expenses, categorize spending, and view financial summaries.

## Primary Goals
1. **Learn Go fundamentals** - Structs, interfaces, error handling, concurrency
2. **Practice clean architecture** - Separation of concerns, dependency injection, testability
3. **Master TUI development** - Build an intuitive, keyboard-driven interface
4. **Implement good coding practices** - Testing, documentation, proper project structure

## Things to consider
- Always check TODO.md to see what is already implemented and what needs to be done

## Technical Decisions

### Architecture Pattern
- **Hexagonal Architecture** - Clear separation between domain, application, and infrastructure layers
- **Repository Pattern** - Abstract data access behind interfaces
- **Dependency Injection** - Explicit dependencies, no global state

### Tech Stack
- **Language**: Go 1.21+
- **TUI Framework**: Bubble Tea (Model-View-Update architecture)
- **Styling**: Lipgloss
- **Database**: SQLite with `database/sql`
- **Migrations**: golang-migrate/migrate
- **CLI**: Cobra (for future CLI commands)
- **Config**: Viper

### Project Structure
```
expense-tracker/
├── cmd/tracker/          # Application entry points
├── internal/             # Private application code
│   ├── model/           # Domain models (Expense, Category)
│   ├── storage/         # Data persistence layer
│   ├── service/         # Business logic
│   └── ui/              # TUI components and views
├── migrations/          # Database migration files
└── config/              # Configuration files
```

## Core Domain Model

### Transaction
- ID (UUID)
- Type (income or expense)
- Amount (float)
- Category (enum)
- Description (string, optional)
- Date (time.Time)
- Tags ([]string) - future enhancement

### Categories (MVP)
For expenses:
- Food & Dining
- Transportation  
- Shopping
- Entertainment
- Bills & Utilities
- Healthcare
- Other

For incomes:
- Salary
- Gifts
- Investments
- Other

## Key User Flows

1. **Quick Add**: User can quickly add incomes or expense with minimal friction
2. **Browse & Filter**: View transactions by type, date range, category, or search term
3. **Insights**: See spending patterns and category breakdowns
4. **Edit/Delete**: Modify or remove existing transactions
5. **Edit categories**: Modify or remove categories

## Development Principles

### Code Quality
- Write tests for service layer (minimum 70% coverage)
- Use `golangci-lint` for consistent code style
- Document public APIs with godoc comments
- Handle errors explicitly, wrap with context

### User Experience
- Keyboard shortcuts for common actions
- Responsive feedback for all actions
- Clear error messages
- Persistent state between sessions

### Data Integrity
- Use transactions for multi-step operations
- Validate inputs at UI and service layers
- Backup reminder for local database

## Future Enhancements (Post-MVP)
- Multi-currency support
- Recurring expenses
- Budget limits and alerts
- Export to CSV/JSON
- Charts and visualizations
- Cloud sync capability
- Split expenses (for shared costs)

## Testing Strategy
1. **Unit Tests**: Service layer business logic
2. **Integration Tests**: Repository layer with test database
3. **Manual Testing**: TUI interactions and user flows

## Performance Considerations
- Pagination for large expense lists
- Indexed database queries on date and category
- Lazy loading of expense details
- Efficient rendering of TUI components

## Security Considerations
- No sensitive data in logs
- Prepared statements for SQL queries (prevent injection)
- File permissions for database file
- Future: Optional encryption for database

## Development Workflow
1. Implement feature in domain/service layer
2. Write tests for new functionality
3. Add database migrations if needed
4. Build TUI interface
5. Manual testing of user flow

## Resources
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [SQLite Best Practices](https://www.sqlite.org/bestpractice.html)