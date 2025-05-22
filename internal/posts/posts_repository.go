package posts

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Iowel/app-base-server/pkg/db"
	"github.com/jackc/pgx/v5"
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
		INSERT INTO posts (user_id, title, content, image, likes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, created_at, updated_at
	`

	err := p.Db.Pool.QueryRow(ctx, query,
		post.UserID, post.Title, post.Content, post.Image, post.Likes,
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
		SELECT id, user_id, title, content, image, created_at, updated_at
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
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Image, &post.CreatedAt, &post.UpdatedAt)
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
            p.id, p.title, p.content, p.user_id, p.image,
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
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &p.Image, &p.UserName, &p.Avatar)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, rows.Err()
}

func (p *PostRepository) LikesUp(postID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE posts
		SET likes = likes + 1
		WHERE id = $1
		RETURNING likes;`

	var newLikes int

	err := p.Db.Pool.QueryRow(ctx, query, postID).Scan(&newLikes)
	if err != nil {
		return fmt.Errorf("failed to add count wallet: %v", err)
	}

	return nil
}

func (p *PostRepository) GetLikeToPost(postID int) (int, error) {
	query := `
		SELECT likes
		FROM posts
		WHERE id = $1
	`

	var newLike int

	err := p.Db.Pool.QueryRow(context.Background(), query, postID).Scan(&newLike)
	if err != nil {
		return 0, fmt.Errorf("failed to get like count: %w", err)
	}

	return newLike, nil
}

func (p *PostRepository) AddCountLikesForOneUser(post_id, user_id int) (*PostLikes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO post_likes (post_id, user_id)
		VALUES ($1, $2)
		RETURNING post_id, user_id
	`
	var postLikes PostLikes

	err := p.Db.Pool.QueryRow(ctx, query,
		post_id, user_id,
	).Scan(
		&postLikes.PostID,
		&postLikes.UserID,
	)
	if err != nil {
		return nil, err
	}

	return &postLikes, nil
}

func (p *PostRepository) DeleteCountLikesForOneUser(post_id, user_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE FROM post_likes 
		WHERE post_id = $1 AND user_id = $2
	`

	_, err := p.Db.Pool.Exec(ctx, query, post_id, user_id)
	if err != nil {
		return fmt.Errorf("failed to delete like from post_likes: %v", err)
	}

	return nil
}

func (p *PostRepository) LikesDown(postID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE posts
		SET likes = likes - 1
		WHERE id = $1 AND likes > 0
		RETURNING likes;`

	var newLikes int

	err := p.Db.Pool.QueryRow(ctx, query, postID).Scan(&newLikes)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("failed to decrease likes count: %v", err)
	}

	return nil
}
