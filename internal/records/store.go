package records

import (
	"fmt"
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

	// paginated results
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

func (s *Store) Update(id string, input UpdateInput) (*Record, error) {
	existing, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Apply partial updates — only overwrite fields that were provided
	if input.Amount != nil {
		existing.Amount = *input.Amount
	}
	if input.Type != nil {
		existing.Type = *input.Type
	}
	if input.Category != nil {
		existing.Category = *input.Category
	}
	if input.Date != nil {
		d, err := time.Parse("2006-01-02", *input.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
		}
		existing.Date = d
	}
	if input.Notes != nil {
		existing.Notes = *input.Notes
	}

	r := &Record{}
	err = s.db.QueryRowx(`
		UPDATE financial_records
		SET amount = $1, type = $2, category = $3, date = $4, notes = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING *`,
		existing.Amount, existing.Type, existing.Category, existing.Date, existing.Notes, id,
	).StructScan(r)
	return r, err
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
