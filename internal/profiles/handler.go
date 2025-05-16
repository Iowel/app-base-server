package profiles

import (
	"github.com/Iowel/app-base-server/configs"
)

type ProfileHandlerDeps struct {
	*configs.Config
	*ProfileService
}

type profileHandler struct {
	*configs.Config
	*ProfileService
}

// func NewProfileHandler(router *http.ServeMux, deps ProfileHandlerDeps) {
// 	handler := &profileHandler{
// 		Config:         deps.Config,
// 		ProfileService: deps.ProfileService,
// 	}

// 	router.HandleFunc("POST /api/get-user-status", handler.GetUserStatus())

// }

// func (p *profileHandler) GetProfile() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		id, _ := strconv.Atoi(r.PathValue("id"))

// 		profile, _ := p.ProfileService.GetProfiles(id)

// 		response.Json(w, profile, http.StatusOK)
// 	}
// }
