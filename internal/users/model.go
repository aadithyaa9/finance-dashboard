package users

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Role string

const (
	RoleViewer  Role = "viewer"
	RoleAnalyst Role = "analyst"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID        string    `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	Email     string    `db:"email"      json:"email"`
	Password  string    `db:"password"   json:"-"`
	Role      Role      `db:"role"       json:"role"`
	IsActive  bool      `db:"is_active"  json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(u *User) error {
	query := `
		INSERT INTO users (name, email, password, role)
		VALUES (:name, :email, :password, :role)
		RETURNING id, created_at, updated_at`
	rows, err := s.db.NamedQuery(query, u)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	}
	return nil
}

func (s *Store) FindByEmail(email string) (*User, error) {
	u := &User{}
	err := s.db.Get(u, `SELECT * FROM users WHERE email = $1`, email)
	return u, err
}

func (s *Store) FindByID(id string) (*User, error) {
	u := &User{}
	err := s.db.Get(u, `SELECT * FROM users WHERE id = $1`, id)
	return u, err
}

func (s *Store) List() ([]User, error) {
	var list []User
	err := s.db.Select(&list, `SELECT * FROM users ORDER BY created_at DESC`)
	return list, err
}

func (s *Store) UpdateRole(id string, role Role) error {
	_, err := s.db.Exec(
		`UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`,
		role, id,
	)
	return err
}

func (s *Store) UpdateStatus(id string, isActive bool) error {
	_, err := s.db.Exec(
		`UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2`,
		isActive, id,
	)
	return err
}
