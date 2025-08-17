package domain

import "time"

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Expense struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	CategoryID  int       `json:"category_id"`
	Category    *Category `json:"category,omitempty"`
}

type Income struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	CategoryID  int       `json:"category_id"`
	Category    *Category `json:"category,omitempty"`
}

type Summary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetBalance   float64 `json:"net_balance"`
}

type Transaction struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Category    string    `json:"category"`
	Type        string    `json:"type"` // "income" or "expense"
}