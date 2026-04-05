package records

import (
	"time"
)

type Record struct {
	ID        string     `db:"id"         json:"id"`
	UserID    string     `db:"user_id"    json:"user_id"`
	Amount    float64    `db:"amount"     json:"amount"`
	Type      string     `db:"type"       json:"type"`
	Category  string     `db:"category"   json:"category"`
	Date      time.Time  `db:"date"       json:"date"`
	Notes     string     `db:"notes"      json:"notes,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateInput struct {
	Amount   float64 `json:"amount"   validate:"required,gt=0"`
	Type     string  `json:"type"     validate:"required,oneof=income expense"`
	Category string  `json:"category" validate:"required,min=1,max=100"`
	Date     string  `json:"date"     validate:"required"`
	Notes    string  `json:"notes"`
}

type UpdateInput struct {
	Amount   *float64 `json:"amount"   validate:"omitempty,gt=0"`
	Type     *string  `json:"type"     validate:"omitempty,oneof=income expense"`
	Category *string  `json:"category" validate:"omitempty,min=1,max=100"`
	Date     *string  `json:"date"`
	Notes    *string  `json:"notes"`
}

type Filter struct {
	Type     string
	Category string
	From     string
	To       string
	Page     int
	Limit    int
}
