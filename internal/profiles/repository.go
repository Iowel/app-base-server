package profiles

import (
	"github.com/Iowel/app-base-server/pkg/db"
	"context"
	"fmt"
	"strconv"
	"time"
)

type ProfileRepository struct {
	Db *db.Db
}

func NewProfileRepository(db *db.Db) *ProfileRepository {
	return &ProfileRepository{
		Db: db,
	}
}

func (p *ProfileRepository) Create(profile *Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        INSERT INTO profiles (user_id, avatar, about, friends, wallet, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := p.Db.Pool.Exec(ctx, query, profile.UserID, profile.Avatar, profile.About, profile.Friends, profile.Wallet, profile.Status, profile.CreatedAt, profile.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepository) GetProfile(id int) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select user_id, avatar, about, friends, status, created_at, updated_at from profiles where user_id = $1`

	var profile Profile
	err := p.Db.Pool.QueryRow(ctx, query, id).Scan(&profile.UserID, &profile.Avatar, &profile.About, &profile.Friends, &profile.Status, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error find profile by id: %w", err)
	}

	return &profile, nil
}

func (p *ProfileRepository) GetStatusNameByID(statusID int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var name string
	err := p.Db.Pool.QueryRow(ctx, "SELECT name FROM statuses WHERE id = $1", statusID).Scan(&name)
	if err != nil {
		return name, err
	}

	return name, nil
}

// userID — id того, кто добавляет друга;
// friendID — id друга.
func (p *ProfileRepository) AddFriends(userID, friendID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := p.Db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		update profiles 
		set friends = array_append(friends, $1)
		where user_id = $2 and not ($1 = any(friends))
	`
	cmdTag, err := p.Db.Pool.Exec(ctx, query, strconv.Itoa(friendID), userID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("profile not found or friend already added")
	}

	return nil
}

func (p *ProfileRepository) DeletedFriends(userID, friendID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := p.Db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Запрос на удаление друга из списка
	query := `
        update profiles
        set friends = array_remove(friends, $1)
        where user_id = $2
    `
	cmdTag, err := tx.Exec(ctx, query, strconv.Itoa(friendID), userID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("profile not found or friend already removed")
	}

	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// func (p *ProfileRepository) GetUserStatus(ctx context.Context, userID int64) (string, error) {
// 	var status string

// 	query := `SELECT status FROM profiles WHERE user_id = $1`

// 	err := p.Db.Pool.QueryRow(ctx, query, userID).Scan(&status)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get user status for user %d: %w", userID, err)
// 	}
// 	return status, nil
// }
