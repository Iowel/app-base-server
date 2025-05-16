package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	err := app.Db.Pool.QueryRow(ctx, "select id, password from users where email = $1", email).Scan(&id, &hashedPassword)
	if err != nil {
		return id, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Role     string
}

func (app *application) GetOneUser(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, name, email, password, role
		FROM users
		WHERE id = $1
	`

	row := app.Db.Pool.QueryRow(ctx, query, id)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}

	return &u, nil
}
