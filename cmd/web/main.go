package main

import (
	"encoding/gob"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Iowel/app-base-server/configs"
	"github.com/Iowel/app-base-server/pkg/db"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
)

const version = "1.0.0"
const cssVersion = "1"

var session *scs.SessionManager

type config struct {
	port        int
	env         string
	api         string
	auth_api    string
	saga_api    string
	prod_serv   string
	server_ip   string
	frontend_ip string
	db          struct {
		dsn string
	}
	secretkey string
	frontend  string
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	Db            *db.Db
	Session       *scs.SessionManager
}

func (app *application) serve() error {

	stack := Chain(
		SessionLoad,
		CORS,
	)

	srv := &http.Server{
		Addr:              ":8082",
		Handler:           stack(app.routes()),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting HTTP server in %s mode on port %d\n", app.config.env, app.config.port)

	return srv.ListenAndServe()
}

func main() {
	gob.Register(map[string]interface{}{})

	var cfg config

	flag.IntVar(&cfg.port, "port", 8082, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:8081", "URL to api")

	// для продакшна
	flag.StringVar(&cfg.server_ip, "server_ip", "http://84.201.150.225:8081", "URL to server_ip")
	flag.StringVar(&cfg.frontend_ip, "frontend_ip", "84.201.150.225", "URL to frontend_ip") // for ws
	flag.StringVar(&cfg.auth_api, "auth_api", "http://84.201.150.225:8083", "URL to auth_api")
	flag.StringVar(&cfg.saga_api, "saga_api", "http://84.201.150.225:8087", "URL to saga_api")
	flag.StringVar(&cfg.prod_serv, "prod_serv", "http://84.201.150.225:8088", "URL to auth_api")

	// для теста на локальной машине
	// flag.StringVar(&cfg.server_ip, "server_ip", "http://localhost:8081", "URL to server_ip")
	// flag.StringVar(&cfg.frontend_ip, "frontend_ip", "localhost", "URL to frontend_ip") // for ws
	// flag.StringVar(&cfg.auth_api, "auth_api", "http://localhost:8083", "URL to auth_api")
	// flag.StringVar(&cfg.saga_api, "saga_api", "http://localhost:8087", "URL to saga_api")
	// flag.StringVar(&cfg.prod_serv, "prod_serv", "http://localhost:8088", "URL to auth_api")

	flag.StringVar(&cfg.secretkey, "secret", "f92JkL5vZxP8tYqB7sNAeW4HrmCXg1dU", "secret key")
	flag.StringVar(&cfg.frontend, "interface", "http://84.201.150.225:8082", "url to front end")
	
	// flag.StringVar(&cfg.frontend, "interface", "http://localhost:8082", "url to front end") // for local test

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// init db
	conf := configs.LoadConfig()
	db, err := db.NewDB(*conf)
	if err != nil {
		log.Fatal(err)
	}

	session = scs.New()
	session.Store = pgxstore.New(db.Pool)

	session.Lifetime = 24 * time.Hour

	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
		Db:            db,
		Session:       session,
	}

	go app.ListenToWsChannel()

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}
