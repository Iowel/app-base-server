package posts

type IPostCache interface {
	Set(key string, value *Post)
	Get(key string) *Post
	GetAll() []*Post
	Delete(key string)

	GetByUserID(userID string) []*Post 
}
