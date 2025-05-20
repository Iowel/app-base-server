package posts

import (
	"context"
	"fmt"
	"time"

	"github.com/Iowel/app-base-server/pkg/db"
)

type PostRepository struct {
	Db *db.Db
}

func NewPostRepository(db *db.Db) *PostRepository {
	return &PostRepository{
		Db: db,
	}
}

func (p *PostRepository) Create(post *Post) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO posts (user_id, title, content)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, created_at, updated_at
	`

	err := p.Db.Pool.QueryRow(ctx, query,
		post.UserID, post.Title, post.Content,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %v", err)
	}

	return post, nil
}

func (p *PostRepository) GetByUser(userID int) ([]*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM posts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := p.Db.Pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %v", err)
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostRepository) GetAllPosts() ([]*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT 
            p.id, p.title, p.content, p.user_id,
            u.name, u.avatar
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.id DESC
    `

	rows, err := p.Db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post

	for rows.Next() {

		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &p.UserName, &p.Avatar)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, rows.Err()
}
