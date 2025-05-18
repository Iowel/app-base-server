package auth

type ProfileResponse struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	About     string `json:"about"`
	Friends   []int  `json:"friends"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ResetPasswordResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type EditProfileRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	About  string `json:"about"`
	Role   string `json:"role"`
}
