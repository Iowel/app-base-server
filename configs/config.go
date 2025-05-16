package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB        Dbconfig
	Auth      AuthConfig
	Web       WebConfig
	Smtp      Smtp
	CryptLink CryptLink
	SmtpGmail SmtpGmail
}

type Dbconfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

type WebConfig struct {
	Port string
	Api  string
	Dsn  string
	Env  string
}

type Smtp struct {
	Host     string
	Port     string
	Username string
	Password string
}

type CryptLink struct {
	Secretkey string
	Frontend  string
}

type SmtpGmail struct {
	SenderName     string
	SenderAddress  string
	SenderPassword string
}



func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}

	return &Config{
		DB:        Dbconfig{Dsn: os.Getenv("DB_DSN")},
		Auth:      AuthConfig{Secret: os.Getenv("SECRET")},
		Web:       WebConfig{Port: os.Getenv("PORT"), Api: os.Getenv("API"), Dsn: os.Getenv("DSN"), Env: os.Getenv("ENV")},
		Smtp:      Smtp{Host: os.Getenv("SMTP_HOST"), Port: os.Getenv("SMTP_PORT"), Username: os.Getenv("SMTP_USERNAME"), Password: os.Getenv("SMTP_PASSWORD")},
		CryptLink: CryptLink{Secretkey: os.Getenv("SECRET_KEY"), Frontend: os.Getenv("FRONTEND_LINK")},
		SmtpGmail: SmtpGmail{SenderName: os.Getenv("EMAIL_SENDER_NAME"), SenderAddress: os.Getenv("EMAIL_SENDER_ADDRESS"), SenderPassword: os.Getenv("EMAIL_SENDER_PASSWORD")},
	}
}
