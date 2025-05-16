package posts

type PostService struct {
	post *PostRepository
}

func NewPostService(Post *PostRepository) *PostService {
	return &PostService{
		post: Post,
	}
}

func (p *PostService) CreatePost(post *Post) error {

	err := p.post.Create(post)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) GetPostsByUserID(id int) ([]*Post, error) {

	posts, err := p.post.GetByUser(id)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostService) GetPostsAllUsers() ([]Post, error) {

	posts, err := p.post.GetAllPosts()
	if err != nil {
		return nil, err
	}

	return posts, nil
}
