package auth

import (
	"errors"
	"time"

	"github.com/aadithyaa9/finance-dashboard/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInactiveUser       = errors.New("account is inactive")
	ErrEmailTaken         = errors.New("email already registered")
)

type Claims struct {
	UserID string     `json:"user_id"`
	Role   users.Role `json:"role"`
	jwt.RegisteredClaims
}

type RegisterInput struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role"     validate:"omitempty,oneof=viewer analyst admin"`
}

type LoginInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *users.User `json:"user"`
}

type Service struct {
	store       *users.Store
	jwtSecret   []byte
	expiryHours int
}

func NewService(store *users.Store, secret string, expiryHours int) *Service {
	return &Service{
		store:       store,
		jwtSecret:   []byte(secret),
		expiryHours: expiryHours,
	}
}

func (s *Service) Register(input RegisterInput) (*AuthResponse, error) {
	existing, err := s.store.FindByEmail(input.Email)
	if err == nil && existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := users.RoleViewer
	if input.Role != "" {
		role = users.Role(input.Role)
	}

	u := &users.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
		Role:     role,
	}

	if err := s.store.Create(u); err != nil {
		return nil, err
	}

	token, err := s.generateToken(u)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: u}, nil
}

func (s *Service) Login(input LoginInput) (*AuthResponse, error) {
	u, err := s.store.FindByEmail(input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !u.IsActive {
		return nil, ErrInactiveUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(u)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: u}, nil
}

func (s *Service) generateToken(u *users.User) (string, error) {
	claims := Claims{
		UserID: u.ID,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.expiryHours) * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func (s *Service) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
