package posts

import "time"

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserName  string    `json:"user_name"`
	Avatar    string    `json:"avatar"`
	Image     string    `json:"image"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostWithUser struct {
	Post       `json:"post"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
	UserStatus string `json:"user_status"`
}

type PostLikes struct {
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
}
