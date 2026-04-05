package records

import (
	"time"
)

// RecordType represents the direction of a financial entry.
type RecordType string

const (
	RecordTypeIncome  RecordType = "income"
	RecordTypeExpense RecordType = "expense"
)

type Record struct {
	ID        string     `db:"id"         json:"id"`
	UserID    string     `db:"user_id"    json:"user_id"`
	Amount    float64    `db:"amount"     json:"amount"`
	Type      RecordType `db:"type"       json:"type"`
	Category  string     `db:"category"   json:"category"`
	Date      time.Time  `db:"date"       json:"date"`
	Notes     string     `db:"notes"      json:"notes,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateInput struct {
	Amount   float64    `json:"amount"   validate:"required,gt=0"`
	Type     RecordType `json:"type"     validate:"required,oneof=income expense"`
	Category string     `json:"category" validate:"required,min=1,max=100"`
	Date     string     `json:"date"     validate:"required"`
	Notes    string     `json:"notes"`
}

type UpdateInput struct {
	Amount   *float64    `json:"amount"   validate:"omitempty,gt=0"`
	Type     *RecordType `json:"type"     validate:"omitempty,oneof=income expense"`
	Category *string     `json:"category" validate:"omitempty,min=1,max=100"`
	Date     *string     `json:"date"`
	Notes    *string     `json:"notes"`
}

// Filter holds parameters for listing and paginating financial records.
//
// Pagination note: this implementation uses LIMIT/OFFSET which is simple and
// correct for small-to-medium datasets. For large datasets (100k+ rows),
// OFFSET pagination degrades because the database must scan and discard all
// preceding rows on every page. The production-grade alternative is keyset
// (cursor-based) pagination: the client sends the last seen (date, id) pair
// and the query uses "WHERE (date, id) < ($last_date, $last_id)" with an
// appropriate composite index. This avoids the full scan and stays fast at
// any scale.
type Filter struct {
	Type     string
	Category string
	From     string
	To       string
	Page     int
	Limit    int
}
