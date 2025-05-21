package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Iowel/app-base-server/configs"
	"github.com/Iowel/app-base-server/internal/auth"
	"github.com/Iowel/app-base-server/internal/posts"
	"github.com/Iowel/app-base-server/internal/profiles"
	"github.com/Iowel/app-base-server/internal/token"
	"github.com/Iowel/app-base-server/internal/user"
	"github.com/Iowel/app-base-server/pkg/cache"
	"github.com/Iowel/app-base-server/pkg/db"
	"github.com/Iowel/app-base-server/pkg/mail"
	"github.com/Iowel/app-base-server/pkg/mailer"
	"github.com/Iowel/app-base-server/pkg/middleware"
)

func main() {
	conf := configs.LoadConfig()

	db, err := db.NewDB(*conf)
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	// repository
	profileRepo := profiles.NewProfileRepository(db)
	userRepo := user.NewUserReposotory(db, profileRepo)
	tokenRepo := token.NewTokenRepository(db)
	mailer := mailer.NewMailer(conf)
	gmailer := mail.NewGmailSender(conf.SmtpGmail.SenderName, conf.SmtpGmail.SenderAddress, conf.SmtpGmail.SenderPassword)
	poseRepo := posts.NewPostRepository(db)

	cacheRepo := cache.NewRedisCache(conf.Redis.Port, conf.Redis.RedisDB, conf.Redis.Exp)
	postCacheRepo := posts.NewRedisCache(conf.Redis.Port, conf.Redis.RedisDB, conf.Redis.Exp)

	// Services
	authService := auth.NewAuthService(auth.AuthServiceDeps{UserRepo: userRepo, Token: tokenRepo, Profile: profileRepo, Mailer: mailer, Config: conf, Gmailer: gmailer, Cache: cacheRepo})
	postService := posts.NewPostService(poseRepo, postCacheRepo, userRepo, cacheRepo)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: conf, AuthService: authService})
	posts.NewPostHandler(router, conf, postService)

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
	)

	// Server
	server := http.Server{
		Addr:         conf.Web.Port,
		Handler:      stack(router),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting HTTP server on port %s\n", conf.Web.Port)
	server.ListenAndServe()
}
