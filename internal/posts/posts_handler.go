package posts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Iowel/app-base-server/configs"
	"github.com/Iowel/app-base-server/pkg/response"
)

type PostHandler struct {
	conf        *configs.Config
	postService *PostService
}

func NewPostHandler(router *http.ServeMux, conf *configs.Config, postService *PostService) {
	handler := &PostHandler{
		conf:        conf,
		postService: postService,
	}

	router.HandleFunc("POST /api/send-posts", handler.SendPost())
	router.HandleFunc("POST /api/get-posts", handler.GetPosts())
	router.HandleFunc("GET /api/all-posts", handler.GetAllPosts())

	router.HandleFunc("POST /api/like-post", handler.LikePosts())

}

func (p *PostHandler) SendPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var post *Post

		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := p.postService.CreatePost(post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		payload.Error = false
		payload.Message = "success"

		response.Json(w, payload, http.StatusOK)
	}
}

func (p *PostHandler) GetPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var post *Post

		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		posts, err := p.postService.GetPostsByUserID(post.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response.Json(w, posts, http.StatusOK)
	}
}

func (p *PostHandler) GetAllPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		posts, err := p.postService.GetPostsAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response.Json(w, posts, http.StatusOK)
	}
}

type LikeRequest struct {
	PostID string `json:"post_id"`
	UserID int    `json:"user_id"`
}

func (p *PostHandler) LikePosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var like LikeRequest

		if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		postIdInt, err := strconv.Atoi(like.PostID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// заносим лайк юзера в таблицу для уникальнсти лайков
		_, err = p.postService.AddCountLikes(postIdInt, like.UserID)
		if err != nil {
			p.postService.RemoveLikeToPost(postIdInt, like.UserID)
			p.postService.LikesDowner(postIdInt)
		} else {
			p.postService.LikesUpper(postIdInt)
		}

		// // полуачем обновлённое количество лайков
		updatedLikes, err := p.postService.post.GetLikeToPost(postIdInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]int{"likes": updatedLikes})
	}
}
