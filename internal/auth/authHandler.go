package auth

import (
	"github.com/Iowel/app-base-server/configs"
	"github.com/Iowel/app-base-server/internal/user"
	"log"
	"strconv"

	"github.com/Iowel/app-base-server/pkg/response"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type authHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &authHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}

	router.HandleFunc("GET /api/get-profile", handler.GetProfile())
	router.HandleFunc("POST /api/update-profile/{id}", handler.UpdateProfile())



	router.HandleFunc("GET /api/get_balance/{id}", handler.GetBalance())
	router.HandleFunc("POST /api/add-balance/{id}", handler.AddBalance())



	router.Handle("GET /api/get-all-users", Auth(handler.AllUsers(), *deps.AuthService))
	router.HandleFunc("POST /api/get-all-users", handler.AllUsers())
	router.HandleFunc("POST /api/get-all-users/{id}", handler.OneUser())
	router.HandleFunc("POST /api/update-user/{id}", handler.UpdateUser())
	router.HandleFunc("DELETE /api/delete-user/{id}", handler.DeleteUser())



	router.HandleFunc("POST /api/forgot-password", handler.SendPasswordResetEmail())
	router.HandleFunc("POST /api/reset-password", handler.ResetPassword())


	router.Handle("GET /add-friends/{id}", Auth(handler.AddFriends(), *deps.AuthService))
	router.HandleFunc("GET /friends/{id}", handler.GetFriends())
	router.HandleFunc("GET /delete-friends/{id}", handler.DeleteFriends())

	router.Handle("POST /api/admin/requestform", Auth(handler.RequestForm(), *deps.AuthService))
}

func (h *authHandler) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _ := h.AuthService.GetProfiles(r)

		response.Json(w, u, http.StatusOK)
	}
}


func (h *authHandler) UpdateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var reqUser EditProfileRequest

		if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("UPDATE PROFILE %+v\n", reqUser)

		// Получаем текущие данные пользователя и профиля
		existUser, err := h.AuthService.GetOneUser(userID)
		if err != nil {
			http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
			return
		}

		log.Printf("EXIST USER %+v\n", existUser)

		existProfile, err := h.AuthService.GetProfileByID(userID)
		if err != nil {
			http.Error(w, "Profile not found: "+err.Error(), http.StatusNotFound)
			return
		}

		// Обновляем только те поля, которые были переданы
		if reqUser.Name != "" {
			existUser.Name = reqUser.Name
		}

		if reqUser.Avatar != "" {
			existProfile.Avatar = reqUser.Avatar
		}

		if reqUser.About != "" {
			existProfile.About = reqUser.About
		}

		if reqUser.Role != "" {
			existUser.Role = reqUser.Role
		}

		if reqUser.Avatar != "" {
			existUser.Avatar = reqUser.Avatar
		}

		// Обновляем данные в базе
		if err := h.AuthService.UpdateUserOne(existUser); err != nil {
			http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h.AuthService.UpdateProfile(existProfile.ID, existProfile); err != nil {
			http.Error(w, "Failed to update profile: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		payload.Error = false
		payload.Message = "success"

		existUsers, err := h.AuthService.GetOneUser(userID)
		if err != nil {
			http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("AFTER UPDATE PROFILE %+v\n", existUsers)

		response.Json(w, payload, http.StatusOK)
	}
}

func (h *authHandler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		idd, _ := strconv.Atoi(id)

		log.Println(idd)

		u, _ := h.AuthService.GetUserBalance(idd)
		log.Println(u)
		response.Json(w, u, http.StatusOK)
	}
}

type BalanceRequest struct {
	Amount int64 `json:"amount"`
}

func (h *authHandler) AddBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		idd, _ := strconv.Atoi(id)

		var req BalanceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("idd", idd)
		log.Println("amount", req.Amount)

		err = h.AuthService.AddUserBalance(int64(idd), req.Amount)

		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		if err != nil {
			payload.Error = true
			payload.Message = "Error while adding balance: " + err.Error()
			http.Error(w, payload.Message, http.StatusInternalServerError)
			return
		}

		payload.Error = false
		payload.Message = "success"
		response.Json(w, payload, http.StatusOK)

	}
}


func (h *authHandler) AllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.AuthService.GetAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, users, http.StatusOK)
	}
}


func (h *authHandler) GetFriends() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		idd, _ := strconv.Atoi(id)

		users, _ := h.AuthService.GetOneUser(idd)

		friends := []*user.User{}

		friends = append(friends, users)

		response.Json(w, friends, http.StatusOK)
	}
}

func (h *authHandler) DeleteFriends() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		idd, _ := strconv.Atoi(id)

		user, _ := h.AuthService.AuthenticateToken(r)

		err := h.AuthService.DeleteFriends(user.ID, idd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("UserID", user.ID)
		fmt.Println("pathID", idd)

		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		payload.Error = false
		payload.Message = "Success add friend!"

		response.Json(w, payload, http.StatusOK)
	}
}



func (h *authHandler) AddFriends() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		idd, _ := strconv.Atoi(id)

		user, _ := h.AuthService.AuthenticateToken(r)

		err := h.AuthService.AddFriends(user.ID, idd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		payload.Error = false
		payload.Message = "Success add friend!"

		response.Json(w, payload, http.StatusOK)
	}
}

func (h *authHandler) OneUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		idd, _ := strconv.Atoi(id)

		user, _ := h.AuthService.GetOneUser(idd)

		response.Json(w, user, http.StatusOK)
	}
}

func (h *authHandler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := strconv.Atoi(r.PathValue("id"))

		var user user.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if userID > 0 {
			// достаем существующего юзера
			existUser, err := h.AuthService.GetOneUser(userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			profile, err := h.AuthService.GetProfileByID(existUser.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if user.Avatar == "" {
				user.Avatar = existUser.Avatar
			}

			// сохраняем старый баланс если он пришел пустой
			// если поле Wallet не передано (null), сохраняем старое значение
			if user.Wallet == nil {
				user.Wallet = &profile.Wallet
			}

			// обновляем только если значение отличается
			if user.Wallet != nil && *user.Wallet != profile.Wallet {
				err = h.AuthService.UpdateWalletProfiles(user.ID, *user.Wallet)
				if err != nil {
					http.Error(w, "ошибка при обновлении баланса", http.StatusInternalServerError)
					return
				}
			}

			if user.Password == "" {
				user.Password = existUser.Password
			}

			// обновляем юзера
			err = h.AuthService.UpdateUser(&user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// если пароль был введён — хэшируем и обновляем

			// обновляем пароль (хэшируем и сохраняем)
			if user.Password != existUser.Password {
				err := h.AuthService.UpdatePassword(&user, user.Password)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			// обновляем роль, если она задана
			if user.Role != "" {
				err := h.AuthService.UpdateRole(user.Role, &user)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

		} else {

			var wallet int
			if user.Wallet != nil {
				wallet = *user.Wallet
			}

			// логика для создания нового юзера
			err := h.AuthService.AddUser(&user, user.Password, wallet)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		var payload struct {
			Error bool `json:"error"`
		}
		payload.Error = false

		response.Json(w, payload, http.StatusOK)
	}
}

func (h *authHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))

		err := h.AuthService.DeleteUser(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		payload.Error = false

		response.Json(w, payload, http.StatusOK)
	}
}

type FormRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (h *authHandler) RequestForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form FormRequest

		err := json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
			Form    FormRequest
		}
		payload.Error = false
		payload.Message = "Success!"
		payload.Form = form

		response.Json(w, payload, http.StatusOK)
	}
}

func (h *authHandler) SendPasswordResetEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Email string `json:"email"`
		}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.AuthService.ForgotPassword(w, payload.Email)
		if err != nil {
			var resp struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}
			resp.Error = true
			resp.Message = "No matching user found for the email"

			response.Json(w, resp, http.StatusOK)
			return
		}

		var resp struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		resp.Error = false

		response.Json(w, resp, http.StatusOK)
	}
}

func (h *authHandler) ResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.AuthService.ResetPassword(w, payload.Email, payload.Password)
		if err != nil {
			var resp struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}
			resp.Error = true
			resp.Message = "Error to change password"

			response.Json(w, resp, http.StatusOK)
			return
		}

		var resp struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		resp.Error = false
		resp.Message = "Password successfully changed"

		response.Json(w, resp, http.StatusOK)
	}
}
