package token

import (
	"context"
	"crypto/sha256"
	"errors"
	"log"
	"time"

	"github.com/Iowel/app-base-server/internal/user"
	"github.com/Iowel/app-base-server/pkg/db"

	"github.com/jackc/pgx/v4"
)

type TokenRepository struct {
	Db *db.Db
}

func NewTokenRepository(Db *db.Db) *TokenRepository {
	return &TokenRepository{Db: Db}
}

func (m *TokenRepository) GetUserForToken(token string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Хешируем токен
	tokenHash := sha256.Sum256([]byte(token))

	var user user.User

	query := `
		SELECT
			u.id, u.name, u.email, u.role
		FROM
			users u
		INNER JOIN tokens t ON u.id = t.user_id
		WHERE
			t.token_hash = $1 AND t.expiry > $2
	`

	err := m.Db.Pool.QueryRow(ctx, query, tokenHash[:], time.Now()).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("Token not found or expired")
			return nil, nil
		}
		log.Println("DB error:", err)
		return nil, err
	}

	return &user, nil
}

func (m *TokenRepository) GetTokenFromUser(user *user.User) (*Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()


	var token Token

	query := `
		SELECT
			expiry
		FROM
			tokens
		WHERE
			user_id = $1
	`

	err := m.Db.Pool.QueryRow(ctx, query, user.ID).Scan(
		&token.Expiry,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("Token not found or expired")
			return nil, nil
		}
		log.Println("DB error:", err)
		return nil, err
	}

	return &token, nil
}
