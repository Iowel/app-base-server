package posts

import (
	"github.com/Iowel/app-base-server/configs"
	"github.com/Iowel/app-base-server/pkg/response"
	"encoding/json"
	"net/http"
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

}

func (p *PostHandler) SendPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var post *Post

		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		err := p.postService.CreatePost(post)
		if err != nil {
			http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
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
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		posts, err := p.postService.GetPostsByUserID(post.UserID)
		if err != nil {
			http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
			return
		}

		response.Json(w, posts, http.StatusOK)
	}
}

func (p *PostHandler) GetAllPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		posts, err := p.postService.GetPostsAllUsers()
		if err != nil {
			http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
			return
		}


		response.Json(w, posts, http.StatusOK)
	}
}
