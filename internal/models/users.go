package models

import (
	"database/sql"
	"errors"
	"fmt"

	"time"

	"github.com/jackc/pgx/v5/pgconn"
	// "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*User, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, NOW() AT TIME ZONE 'UTC')"

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var pqErr *pgconn.PgError
		// unwrappedErr := errors.Unwrap(err)
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" && pqErr.ConstraintName == "users_uc_email" {
				return ErrDuplicateEmail
			}
		} else {
			fmt.Println("Errors is not pq.Error, type", fmt.Sprintf("%T", err))
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = $1"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) Get(id int) (*User, error) {
	var user User

	stmt := "SELECT id, name, email, created FROM users WHERE id = $1"

	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string)  error {
	var currentHashedPassword []byte

	stmt := "SELECT hashed_password FROM users WHERE id = $1"

	err := m.DB.QueryRow(stmt, id).Scan(&currentHashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(currentHashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	stmt = "UPDATE users SET hashed_password = $1 WHERE id = $2"
	_, err = m.DB.Exec(stmt, string(newHashedPassword), id)
	return err
}
