package auth

import (
	"github.com/Iowel/app-base-server/configs"
	"github.com/Iowel/app-base-server/internal/profiles"
	"github.com/Iowel/app-base-server/internal/token"
	"github.com/Iowel/app-base-server/internal/user"
	"github.com/Iowel/app-base-server/pkg/encryption"
	"github.com/Iowel/app-base-server/pkg/mail"
	"github.com/Iowel/app-base-server/pkg/mailer"
	"github.com/Iowel/app-base-server/pkg/urlsigner"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceDeps struct {
	UserRepo *user.UserRepository
	Token    *token.TokenRepository
	Profile  *profiles.ProfileRepository
	Mailer   *mailer.Mailer
	Config   *configs.Config
	Gmailer  mail.EmailSender
}

type AuthService struct {
	UserRepo *user.UserRepository
	Token    *token.TokenRepository
	Profile  *profiles.ProfileRepository
	*mailer.Mailer
	Gmailer mail.EmailSender
	*configs.Config
}

func NewAuthService(deps AuthServiceDeps) *AuthService {
	return &AuthService{UserRepo: deps.UserRepo, Token: deps.Token, Profile: deps.Profile, Mailer: deps.Mailer, Config: deps.Config, Gmailer: deps.Gmailer}
}

func (s *AuthService) GetAllUsers() ([]*user.User, error) {
	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		return nil, errors.New(ErrGetAllUsers)
	}

	return users, nil
}

func (s *AuthService) GetUserBalance(id int) (int, error) {
	wallet, err := s.UserRepo.GetBalance(id)
	if err != nil {
		return wallet, errors.New(ErrGetAllUsers)
	}

	return wallet, nil
}

func (s *AuthService) AddUserBalance(id int64, amount int64) error {
	err := s.UserRepo.AddBalance(id, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) AddUser(user *user.User, password string, wallet int) error {

	// create password
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.UserRepo.AddUser(user, string(newHash), wallet)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) AuthenticateToken(r *http.Request) (*user.User, error) {
	// check header request
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("Authorization header is missing")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("Invalid token format")
	}

	// get token
	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("Authentication token size mismatch")
	}

	// get the user from tokens table
	user, err := s.Token.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("No matching user found for the given token")
	}

	log.Printf("AuthenticateTokenUser %+v\n", user)

	return user, nil
}

// PROFILE

func (s *AuthService) GetProfiles(r *http.Request) (*ProfileResponse, error) {

	user, _ := s.AuthenticateToken(r)

	// получаем профиль
	profile, err := s.Profile.GetProfile(user.ID)
	if err != nil {
		return nil, err
	}

	resp := ProfileResponse{
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    profile.Avatar,
		About:     profile.About,
		Friends:   profile.Friends,
		Status:    profile.Status,
		CreatedAt: profile.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: profile.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return &resp, nil
}


func (s *AuthService) ForgotPassword(w http.ResponseWriter, email string) error {

	// verify user exist
	_, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	// create url address for signer
	link := fmt.Sprintf("%s/reset-password?email=%s", s.Config.CryptLink.Frontend, email)

	sign := urlsigner.Signer{
		Secret: []byte(s.Config.CryptLink.Secretkey),
	}

	signedLink := sign.GenerateTokenFromString(link)

	var data struct {
		Link string
	}
	data.Link = signedLink

	// send mail
	err = s.Mailer.SendMail("foxy@mail.com", email, "Запрос на восстановление пароля", "password-reset", data)
	if err != nil {
		return err
	}

	// send gmail
	subject := "Some organization"

	err = s.Gmailer.Sendmail("password-reset", data, subject, email, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	return nil
}

func (s *AuthService) ResetPassword(w http.ResponseWriter, email, password string) error {

	encryptor := encryption.Encryption{
		Key: []byte(s.Config.CryptLink.Secretkey),
	}

	realEmail, err := encryptor.Decrypt(email)
	if err != nil {
		return err
	}

	user, err := s.UserRepo.FindByEmail(realEmail)
	if err != nil {
		return err
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	err = s.UserRepo.UpdatePasswordForUser(user, string(newHash))
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UpdatePassword(user *user.User, password string) error {

	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	err = s.UserRepo.UpdatePasswordForUser(user, string(newHash))
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UpdateRole(role string, user *user.User) error {
	err := s.UserRepo.UpdateUserRole(role, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetOneUser(id int) (*user.User, error) {
	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	fmt.Println(user)

	return user, nil
}

func (s *AuthService) UpdateUser(user *user.User) error {

	// обновляем
	err := s.UserRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UpdateUserOne(user *user.User) error {

	// обновляем
	err := s.UserRepo.UpdateUserOne(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) DeleteUser(id int) error {
	err := s.UserRepo.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}

// userID — id того, кто добавляет друга
// friendID — id друга
func (s *AuthService) AddFriends(userID, friendID int) error {
	err := s.Profile.AddFriends(userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

// userID — id того, кто удаляет друга
// friendID — id друга
func (s *AuthService) DeleteFriends(userID, friendID int) error {
	err := s.Profile.DeletedFriends(userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetProfileByID(id int) (*profiles.Profile, error) {
	profile, err := s.UserRepo.GetProfileByID(id)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *AuthService) UpdateProfile(id int, profile *profiles.Profile) error {

	// обновляем
	err := s.UserRepo.UpdatePrifile(id, profile)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UpdateWalletProfiles(id int, wallet int) error {
	err := s.UserRepo.UpdateWalletProfile(id, wallet)
	if err != nil {
		return err
	}

	return nil
}
