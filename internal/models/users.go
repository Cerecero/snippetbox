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
			if pqErr.Code == "23505" && pqErr.ConstraintName== "users_uc_email" {
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
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
