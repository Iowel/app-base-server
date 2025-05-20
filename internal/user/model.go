package user

import (
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Token     string
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar"`
	Status    string    `json:"status"`
	Wallet    *int      `json:"wallet"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProfileWithUser struct {
	ProfileID int
	UserID    int
	Status    string
	UserName  string
	UserEmail string
	Avatar    string
	About     string
	Friends   []int32
	Wallet    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCache struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Token     string
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar"`
	Status    string    `json:"status"`
	Wallet    int       `json:"wallet"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
