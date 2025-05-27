package posts

import (
	"log"
	"strconv"

	"github.com/Iowel/app-base-server/internal/user"
	"github.com/Iowel/app-base-server/pkg/cache"
	"github.com/jackc/pgconn"
)

type PostService struct {
	post      *PostRepository
	cache     IPostCache
	user      *user.UserRepository
	userCache cache.IPostCache
}

func NewPostService(Post *PostRepository, cache IPostCache, user *user.UserRepository, userCache cache.IPostCache) *PostService {
	return &PostService{
		post:      Post,
		cache:     cache,
		user:      user,
		userCache: userCache,
	}
}

func (p *PostService) CreatePost(postInput *Post) error {

	post, err := p.post.Create(postInput)
	if err != nil {
		return err
	}

	var cacheUser = &Post{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		Image:     post.Image,
		Likes:     post.Likes,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	idStr := "post:" + strconv.Itoa(cacheUser.ID)
	p.cache.Set(idStr, cacheUser)

	return nil
}

func (p *PostService) GetPostsByUserID(id int) ([]*Post, error) {
	userID := strconv.Itoa(id)

	posts := p.cache.GetByUserID(userID)
	if len(posts) == 0 {
		posts, err := p.post.GetByUser(id)
		if err != nil {
			return nil, err
		}
		return posts, nil
	} else {
		return posts, nil
	}

}

func (p *PostService) GetPostsAllUsers() ([]*Post, error) {

	posts := p.cache.GetAll()

	// грузим из базы если кэш пуст
	if len(posts) == 0 {
		dbPosts, err := p.post.GetAllPosts()
		if err != nil {
			return nil, err
		}

		// кжшируем каждый пост
		for _, post := range dbPosts {
			key := "post:" + strconv.Itoa(post.ID)
			p.cache.Set(key, post)
		}

		posts = dbPosts
	}

	// дополняем имя с аватаркой
	for _, post := range posts {
		idStr := "user:" + strconv.Itoa(post.UserID)
		user := p.userCache.Get(idStr)

		like, _ := p.post.GetLikeToPost(post.ID)

		if like != 0 {
			post.Likes = like
		}

		if user != nil {
			post.UserName = user.Name
			post.Avatar = user.Avatar

		}

	}

	log.Printf("GETALLPOSTS: %v\n", posts[1].CreatedAt)

	return posts, nil

}

func (p *PostService) LikesUpper(postID int) error {

	if err := p.post.LikesUp(postID); err != nil {
		return err
	}

	return nil
}

func (p *PostService) GetLikeToPost(postID int) (int, error) {
	return p.post.GetLikeToPost(postID)
}

func (p *PostService) AddCountLikes(postID, userID int) (bool, error) {
	_, err := p.post.AddCountLikesForOneUser(postID, userID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (p *PostService) RemoveLikeToPost(postID, userID int) error {
	err := p.post.DeleteCountLikesForOneUser(postID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) LikesDowner(postID int) error {

	if err := p.post.LikesDown(postID); err != nil {
		return err
	}

	return nil
}
