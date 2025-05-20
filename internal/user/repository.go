package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Iowel/app-base-server/internal/profiles"
	"github.com/Iowel/app-base-server/pkg/db"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	Db      *db.Db
	Profile *profiles.ProfileRepository
}

func NewUserReposotory(db *db.Db, profile *profiles.ProfileRepository) *UserRepository {
	return &UserRepository{
		Db:      db,
		Profile: profile,
	}
}

func (u *UserRepository) FindByEmail(email string) (*User, error) {
	query := `SELECT id, email, name, password FROM users WHERE email = $1`

	var user User

	err := u.Db.Pool.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &user, nil
}

func (u *UserRepository) GetBalance(id int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select wallet from profiles where user_id = $1
	`

	row := u.Db.Pool.QueryRow(ctx, query, id)

	var wallet int
	err := row.Scan(&wallet)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("user with id %d not found", id)
		}
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return wallet, nil
}

func (u *UserRepository) AddBalance(userID int64, amount int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        UPDATE
            profiles
        SET
            wallet = wallet + $1
        WHERE
            user_id = $2
        RETURNING wallet;`

	var newBalance int64
	err := u.Db.Pool.QueryRow(ctx, query, amount, userID).Scan(&newBalance)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user with ID %d not found", userID)
		}
		return fmt.Errorf("failed to update balance: %v", err)
	}

	return nil
}

func (u *UserRepository) GetAllUsers() ([]*UserCache, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT u.id, u.email, u.name, u.avatar, p.status
		FROM users u
		LEFT JOIN profiles p ON p.user_id = u.id
	`

	rows, err := u.Db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []*UserCache
	for rows.Next() {
		var user UserCache
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Avatar,
			&user.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during iteration: %w", err)
	}

	return users, nil
}

func (r *UserRepository) GetAllProfilesWithUser() ([]*ProfileWithUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT 
			p.id, p.user_id, u.name, u.email, 
			p.avatar, p.about, p.friends, p.wallet,
			p.created_at, p.updated_at
		FROM profiles p
		JOIN users u ON p.user_id = u.id
	`

	rows, err := r.Db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query profiles: %w", err)
	}
	defer rows.Close()

	var profiles []*ProfileWithUser
	for rows.Next() {
		var p ProfileWithUser
		err := rows.Scan(
			&p.ProfileID,
			&p.UserID,
			&p.UserName,
			&p.UserEmail,
			&p.Avatar,
			&p.About,
			&p.Friends,
			&p.Wallet,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profile: %w", err)
		}
		profiles = append(profiles, &p)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	return profiles, nil
}

func (u *UserRepository) GetUserByID(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select id, email, password, name, avatar, role from users where id = $1
	`

	row := u.Db.Pool.QueryRow(ctx, query, id)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Avatar, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (u *UserRepository) AddUser(user *User, hash string, wallet int) (*User, *profiles.Profile, error) {
	query := `
		INSERT INTO users (email, password, name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, password, created_at, updated_at
	`

	err := u.Db.Pool.QueryRow(context.Background(), query, user.Email, hash, user.Name, user.Role).Scan(
		&user.ID,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %v", err)
	}

	updateQuery := `
		UPDATE
			users
		SET
			is_email_verified = true
		WHERE id = $1
	`

	_, err = u.Db.Pool.Exec(context.Background(), updateQuery, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update is_email_verified: %v", err)
	}

	// Создаем профиль
	profile := &profiles.Profile{
		UserID:    user.ID,
		Avatar:    "static/fox-icon.png",
		About:     "default description",
		Friends:   []int{},
		Wallet:    wallet,
		Status:    "Серебренный",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.Profile.Create(profile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create profile: %v", err)
	}

	return user, profile, nil
}

func (u *UserRepository) UpdateUser(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE
			users
		SET
			email = $2,
			name = $3,
			password = $4,
			role = $5,
			avatar = $6
		WHERE
			id = $1
		RETURNING id, email, name, password, role, avatar 
	`

	err := u.Db.Pool.QueryRow(ctx, query, user.ID, user.Email, user.Name, user.Password, user.Role, user.Avatar).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.Role,
		&user.Avatar,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (u *UserRepository) UpdateUserOne(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update users
		set
			name = $2,
			avatar = $3
		where
			id = $1
	`

	_, err := u.Db.Pool.Exec(ctx, query, user.ID, user.Name, user.Avatar)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (u *UserRepository) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmdTag, err := u.Db.Pool.Exec(ctx, `delete from users where id = $1`, id)

	if cmdTag.RowsAffected() == 0 {
		return err
	}

	stmt := `delete from tokens where user_id = $1`
	cmdTag, err = u.Db.Pool.Exec(ctx, stmt, id)

	if cmdTag.RowsAffected() == 0 {
		return err
	}

	return nil
}

func (u *UserRepository) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	err := u.Db.Pool.QueryRow(ctx, "select id, password from users where email = $1", email).Scan(&id, &hashedPassword)
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

func (u *UserRepository) UpdatePasswordForUser(user *User, hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update users set password = $1 where id = $2`

	_, err := u.Db.Pool.Exec(ctx, stmt, hash, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	return nil
}

func (u *UserRepository) UpdateUserRole(role string, user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update users set role = $1 where id = $2`

	_, err := u.Db.Pool.Exec(ctx, stmt, role, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	return nil
}

func (u *UserRepository) GetProfileByID(id int) (*profiles.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select id, avatar, status, wallet, about from profiles where user_id = $1
	`

	row := u.Db.Pool.QueryRow(ctx, query, id)

	var profile profiles.Profile

	err := row.Scan(&profile.ID, &profile.Avatar, &profile.Status, &profile.Wallet, &profile.About)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("profile with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &profile, nil
}

func (u *UserRepository) UpdatePrifile(id int, profile *profiles.Profile) (*profiles.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update profiles
		set
			avatar = $2,
			about = $3
		where
			id = $1
	`

	_, err := u.Db.Pool.Exec(ctx, query, id, profile.Avatar, profile.About)

	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return profile, nil
}

func (u *UserRepository) GetUserStatus(userID int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var status string

	query := `SELECT status FROM profiles WHERE user_id = $1`

	err := u.Db.Pool.QueryRow(ctx, query, userID).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("failed to get user status for user %d: %w", userID, err)
	}
	return status, nil
}

func (u *UserRepository) UpdateWalletProfile(id int, wallet int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	update profiles
	set
		wallet = $2
	where
		user_id = $1
`

	_, err := u.Db.Pool.Exec(ctx, query, id, wallet)

	if err != nil {
		return fmt.Errorf("failed to update profiles: %w", err)
	}

	return nil
}
