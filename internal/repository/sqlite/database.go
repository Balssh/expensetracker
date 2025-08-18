package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const schema = `
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
    UNIQUE(name, type)
);

CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    amount REAL NOT NULL,
    date TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
    category_id INTEGER,
    FOREIGN KEY(category_id) REFERENCES categories(id)
);

INSERT OR IGNORE INTO categories (name, type) VALUES 
    ('Food & Dining', 'expense'),
    ('Transportation', 'expense'),
    ('Shopping', 'expense'),
    ('Entertainment', 'expense'),
    ('Bills & Utilities', 'expense'),
    ('Healthcare', 'expense'),
    ('Other', 'expense'),
    ('Salary', 'income'),
    ('Freelance', 'income'),
    ('Investment', 'income'),
    ('Gift', 'income'),
    ('Other', 'income');
`

type Database struct {
	db *sql.DB
}

func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{db: db}

	if err := database.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return database, nil
}

func (d *Database) initialize() error {
	_, err := d.db.Exec(schema)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) DB() *sql.DB {
	return d.db
}