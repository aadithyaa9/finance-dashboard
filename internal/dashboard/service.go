package dashboard

import (
	"fmt"

	"github.com/aadithyaa9/finance-dashboard/internal/records"
	"github.com/jmoiron/sqlx"
)

type Summary struct {
	TotalIncome   float64 `db:"total_income"   json:"total_income"`
	TotalExpenses float64 `db:"total_expenses" json:"total_expenses"`
	NetBalance    float64 `db:"net_balance"    json:"net_balance"`
}

type CategoryTotal struct {
	Category string             `db:"category" json:"category"`
	Type     records.RecordType `db:"type"     json:"type"`
	Total    float64            `db:"total"    json:"total"`
}

type MonthlyTrend struct {
	Month string             `db:"month" json:"month"`
	Type  records.RecordType `db:"type"  json:"type"`
	Total float64            `db:"total" json:"total"`
}

type RecentRecord struct {
	ID       string             `db:"id"       json:"id"`
	Amount   float64            `db:"amount"   json:"amount"`
	Type     records.RecordType `db:"type"     json:"type"`
	Category string             `db:"category" json:"category"`
	Date     string             `db:"date"     json:"date"`
	Notes    string             `db:"notes"    json:"notes"`
}

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

// GetSummary returns total income, total expenses, and net balance.
func (s *Service) GetSummary() (*Summary, error) {
	summary := &Summary{}
	err := s.db.QueryRowx(fmt.Sprintf(`
		SELECT
			COALESCE(SUM(CASE WHEN type = '%s' THEN amount ELSE 0 END), 0)          AS total_income,
			COALESCE(SUM(CASE WHEN type = '%s' THEN amount ELSE 0 END), 0)          AS total_expenses,
			COALESCE(SUM(CASE WHEN type = '%s' THEN amount ELSE -amount END), 0)    AS net_balance
		FROM financial_records
		WHERE deleted_at IS NULL`,
		records.RecordTypeIncome,
		records.RecordTypeExpense,
		records.RecordTypeIncome,
	)).StructScan(summary)
	return summary, err
}

// GetByCategory returns totals grouped by category and type.
func (s *Service) GetByCategory() ([]CategoryTotal, error) {
	var rows []CategoryTotal
	err := s.db.Select(&rows, `
		SELECT category, type, COALESCE(SUM(amount), 0) AS total
		FROM financial_records
		WHERE deleted_at IS NULL
		GROUP BY category, type
		ORDER BY total DESC`,
	)
	return rows, err
}

// GetMonthlyTrends returns per-month income and expense totals for the last 12 months.
func (s *Service) GetMonthlyTrends() ([]MonthlyTrend, error) {
	var rows []MonthlyTrend
	err := s.db.Select(&rows, `
		SELECT
			TO_CHAR(date, 'YYYY-MM') AS month,
			type,
			COALESCE(SUM(amount), 0) AS total
		FROM financial_records
		WHERE deleted_at IS NULL
		  AND date >= NOW() - INTERVAL '12 months'
		GROUP BY month, type
		ORDER BY month DESC, type`,
	)
	return rows, err
}

// GetRecent returns the 10 most recent records.
func (s *Service) GetRecent() ([]RecentRecord, error) {
	var rows []RecentRecord
	err := s.db.Select(&rows, `
		SELECT id, amount, type, category, TO_CHAR(date, 'YYYY-MM-DD') AS date, COALESCE(notes, '') AS notes
		FROM financial_records
		WHERE deleted_at IS NULL
		ORDER BY date DESC, created_at DESC
		LIMIT 10`,
	)
	return rows, err
}


type Summary struct {
	TotalIncome   float64 `db:"total_income"   json:"total_income"`
	TotalExpenses float64 `db:"total_expenses" json:"total_expenses"`
	NetBalance    float64 `db:"net_balance"    json:"net_balance"`
}

type CategoryTotal struct {
	Category string  `db:"category" json:"category"`
	Type     string  `db:"type"     json:"type"`
	Total    float64 `db:"total"    json:"total"`
}

type MonthlyTrend struct {
	Month string  `db:"month" json:"month"`
	Type  string  `db:"type"  json:"type"`
	Total float64 `db:"total" json:"total"`
}

type RecentRecord struct {
	ID       string  `db:"id"       json:"id"`
	Amount   float64 `db:"amount"   json:"amount"`
	Type     string  `db:"type"     json:"type"`
	Category string  `db:"category" json:"category"`
	Date     string  `db:"date"     json:"date"`
	Notes    string  `db:"notes"    json:"notes"`
}

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

// GetSummary returns total income, total expenses, and net balance.
func (s *Service) GetSummary() (*Summary, error) {
	summary := &Summary{}
	err := s.db.QueryRowx(`
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expenses,
			COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE -amount END), 0) AS net_balance
		FROM financial_records
		WHERE deleted_at IS NULL`,
	).StructScan(summary)
	return summary, err
}

// GetByCategory returns totals grouped by category and type.
func (s *Service) GetByCategory() ([]CategoryTotal, error) {
	var rows []CategoryTotal
	err := s.db.Select(&rows, `
		SELECT category, type, COALESCE(SUM(amount), 0) AS total
		FROM financial_records
		WHERE deleted_at IS NULL
		GROUP BY category, type
		ORDER BY total DESC`,
	)
	return rows, err
}

// GetMonthlyTrends returns per-month income and expense totals for the last 12 months.
func (s *Service) GetMonthlyTrends() ([]MonthlyTrend, error) {
	var rows []MonthlyTrend
	err := s.db.Select(&rows, `
		SELECT
			TO_CHAR(date, 'YYYY-MM') AS month,
			type,
			COALESCE(SUM(amount), 0) AS total
		FROM financial_records
		WHERE deleted_at IS NULL
		  AND date >= NOW() - INTERVAL '12 months'
		GROUP BY month, type
		ORDER BY month DESC, type`,
	)
	return rows, err
}

// GetRecent returns the 10 most recent records.
func (s *Service) GetRecent() ([]RecentRecord, error) {
	var rows []RecentRecord
	err := s.db.Select(&rows, `
		SELECT id, amount, type, category, TO_CHAR(date, 'YYYY-MM-DD') AS date, COALESCE(notes, '') AS notes
		FROM financial_records
		WHERE deleted_at IS NULL
		ORDER BY date DESC, created_at DESC
		LIMIT 10`,
	)
	return rows, err
}
