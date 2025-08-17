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
│   │   ├── domain/         # Entities (Expense, Category structs)
│   │   └── usecase/        # Use Cases and Repository Interfaces
│   ├── handler/
│   │   └── tui/            # Bubble Tea models, views, and components
│   └── repository/
│       └── sqlite/         # SQLite implementation of the repository interfaces
├── go.mod
└── go.sum


## 🗄️ 4. Database Schema

We'll start with four simple tables in our SQLite database to track both expenses and income.

**`expense_categories` table:**
- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `name` TEXT NOT NULL UNIQUE

**`expenses` table:**
- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `description` TEXT NOT NULL
- `amount` REAL NOT NULL
- `date` TEXT NOT NULL
- `category_id` INTEGER
- `FOREIGN KEY(category_id) REFERENCES expense_categories(id)`

**`income_categories` table:**
- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `name` TEXT NOT NULL UNIQUE

**`income` table:**
- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `description` TEXT NOT NULL
- `amount` REAL NOT NULL
- `date` TEXT NOT NULL
- `category_id` INTEGER
- `FOREIGN KEY(category_id) REFERENCES income_categories(id)`

## 🌊 5. User Flows

1.  **Main View (Dashboard):**
    -   On launch, the user sees a summary for the current month: total income, total expenses, and the net balance (Income - Expenses).
    -   Also shows a list of recent transactions (both income and expenses).
    -   Keybindings are displayed for navigation (e.g., `a` to add expense, `i` to add income, `l` to list all, `q` to quit).

2.  **Add Expense / Add Income:**
    -   User presses `a` (for expense) or `i` (for income).
    -   An input form appears with fields for Description, Amount, Category, and Date.
    -   The 'Category' field is a dropdown/selectable list populated from the relevant category table.
    -   User submits the form, the data is saved, and the view returns to the dashboard, which is now updated.

3.  **List & Filter Transactions:**
    -   User presses `l`.
    -   A new view shows a paginated table of all transactions, sorted by date. Income and expenses are clearly distinguished (e.g., by color or a symbol).
    -   The user can press `/` to enter a filter mode, allowing them to search by description, filter by category, or show only income/expenses.

## ✅ 6. Task Breakdown

Here is a step-by-step plan to build the application:

1.  **Project Setup:**
    -   Initialize the Go module (`go mod init`).
    -   Create the directory structure outlined above.
    -   Add dependencies (`go get`).

2.  **Core Logic (Inner Layers):**
    -   Define the `Expense`, `Income`, and `Category` structs in `internal/core/domain/`.
    -   Define the repository interfaces (e.g., `ExpenseRepository`, `IncomeRepository`) in `internal/core/usecase/`.
    -   Implement the use cases (e.g., `AddExpenseUseCase`, `AddIncomeUseCase`) that use these interfaces.

3.  **Database Implementation (Outer Layer):**
    -   Implement the repository interfaces in `internal/repository/sqlite/`. This package will contain the actual SQL queries to interact with the SQLite database.
    -   Write the database initialization logic (creating all four tables if they don't exist).

4.  **TUI Implementation (Outer Layer):**
    -   Start with a simple "Hello, World" Bubble Tea application in `internal/handler/tui/` to ensure it runs.
    -   Build the main dashboard model/view to display income, expenses, and net balance.
    -   Build the "Add Expense" and "Add Income" form models/views.
    -   Build the unified "List Transactions" table model/view.

5.  **Integration:**
    -   In `cmd/app/main.go`, wire everything together:
        -   Initialize the database connection.
        -   Initialize the repository implementations.
        -   Initialize the use cases, injecting the repositories.
        -   Initialize the main TUI model, injecting the use cases.
        -   Run the Bubble Tea program.

