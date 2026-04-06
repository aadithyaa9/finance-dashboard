package records

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(userID string, input CreateInput) (*Record, error) {
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}

	r := &Record{}
	err = s.db.QueryRowx(`
		INSERT INTO financial_records (user_id, amount, type, category, date, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *`,
		userID, input.Amount, input.Type, input.Category, date, input.Notes,
	).StructScan(r)
	return r, err
}

func (s *Store) List(f Filter) ([]Record, int, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}

	where := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	i := 1

	if f.Type != "" {
		where += fmt.Sprintf(" AND type = $%d", i)
		args = append(args, f.Type)
		i++
	}
	if f.Category != "" {
		where += fmt.Sprintf(" AND category ILIKE $%d", i)
		args = append(args, "%"+f.Category+"%")
		i++
	}
	if f.From != "" {
		where += fmt.Sprintf(" AND date >= $%d", i)
		args = append(args, f.From)
		i++
	}
	if f.To != "" {
		where += fmt.Sprintf(" AND date <= $%d", i)
		args = append(args, f.To)
		i++
	}

	// total count
	var total int
	err := s.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM financial_records %s", where), args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// LIMIT/OFFSET is used here for simplicity. See Filter for a note on
	// keyset pagination as a future scaling improvement.
	offset := (f.Page - 1) * f.Limit
	query := fmt.Sprintf(
		"SELECT * FROM financial_records %s ORDER BY date DESC, created_at DESC LIMIT $%d OFFSET $%d",
		where, i, i+1,
	)
	args = append(args, f.Limit, offset)

	var list []Record
	err = s.db.Select(&list, query, args...)
	return list, total, err
}

func (s *Store) GetByID(id string) (*Record, error) {
	r := &Record{}
	err := s.db.Get(r, `SELECT * FROM financial_records WHERE id = $1 AND deleted_at IS NULL`, id)
	return r, err
}

// Update builds a dynamic UPDATE statement containing only the fields present
// in the payload. This avoids the read-modify-write pattern that causes lost
// updates when two requests modify the same record concurrently.
func (s *Store) Update(id string, input UpdateInput) (*Record, error) {
	setClauses := []string{}
	args := []interface{}{}
	i := 1

	if input.Amount != nil {
		setClauses = append(setClauses, fmt.Sprintf("amount = $%d", i))
		args = append(args, *input.Amount)
		i++
	}
	if input.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", i))
		args = append(args, *input.Type)
		i++
	}
	if input.Category != nil {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", i))
		args = append(args, *input.Category)
		i++
	}
	if input.Date != nil {
		d, err := time.Parse("2006-01-02", *input.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
		}
		setClauses = append(setClauses, fmt.Sprintf("date = $%d", i))
		args = append(args, d)
		i++
	}
	if input.Notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", i))
		args = append(args, *input.Notes)
		i++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields provided to update")
	}

	// updated_at is always refreshed
	setClauses = append(setClauses, "updated_at = NOW()")

	// id is the final argument
	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE financial_records SET %s WHERE id = $%d AND deleted_at IS NULL RETURNING *",
		strings.Join(setClauses, ", "),
		i,
	)

	r := &Record{}
	err := s.db.QueryRowx(query, args...).StructScan(r)
	if err != nil {
		return nil, fmt.Errorf("record not found or already deleted")
	}
	return r, nil
}

func (s *Store) SoftDelete(id string) error {
	res, err := s.db.Exec(
		`UPDATE financial_records SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}
