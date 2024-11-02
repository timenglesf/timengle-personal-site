package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/ui/template"
)

const (
	sessionUserId  = "authenticatedUserID"
	sessionIsAdmin = "isAdmin"
)

var version = "1.0.2"

type application struct {
	logger            *slog.Logger
	cfg               *config
	meta              *models.MetaModel
	user              models.UserModelInterface
	post              models.PostModelInterface
	db                *gorm.DB
	sessionManager    *scs.SessionManager
	formDecoder       *form.Decoder
	pageTemplates     *template.Pages
	partialTemplates  *template.Partials
	mostRecentPost    *models.Post
	latestPublicPosts *[]models.Post
}

type config struct {
	port          string
	environment   string
	db            psqlConfig
	objectStorage objectStorageConfig
	secureCookies bool
}

type psqlConfig struct {
	local    bool
	host     string
	port     string
	user     string
	password string
	name     string
	dsn      string
}

func (c *psqlConfig) setDSN() {
	c.dsn = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.host,
		c.port,
		c.user,
		c.password,
		c.name,
	)
}

func (c *psqlConfig) setLocalDSN() {
	c.dsn = fmt.Sprintf("host=%s port=%s dbname=%s sslmode=disable", c.host, c.port, c.name)
}

func (c *psqlConfig) getDSN() string {
	return c.dsn
}

type objectStorageConfig struct {
	objectStorageURL         string
	serveStaticObjectStorage bool
}

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var cfg config

	// Set port to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	cfg.port = port

	// Set environment to development if not set to production
	env := os.Getenv("GOENV")
	if env != "production" {
		logger.Info("setting environment to development")
		env = "development"
	} else {
		logger.Info("setting environment to production")
	}
	cfg.environment = env

	// Parse command-line flags
	flag.BoolVar(&cfg.objectStorage.serveStaticObjectStorage, "object-storage", false, "Serve static files from object storage")
	flag.Parse()

	// Configure object storage if enabled
	if cfg.objectStorage.serveStaticObjectStorage {
		osURL := os.Getenv("OBJECT_STORAGE_URL")
		if osURL == "" {
			log.Fatal("OBJECT_STORAGE_URL must be set when object storage is enabled")
		}
		targetFile := fmt.Sprintf("%s/static/dist/js/form-prevent.js", osURL)
		// #nosec G107
		resp, err := http.Get(targetFile)
		if err != nil {
			log.Fatal("Unable to connect to object storage")
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatal("Unable to connect to object storage")
		}
		cfg.objectStorage.objectStorageURL = osURL
	}

	// Configure database connection
	cfg.db.local = os.Getenv("DBLOCAL") == "true"
	cfg.db.host = os.Getenv("DBHOST")
	cfg.db.port = os.Getenv("DBPORT")
	cfg.db.name = os.Getenv("DBNAME")
	cfg.db.user = os.Getenv("DBUSER")
	cfg.db.password = os.Getenv("DBPASSWORD")

	if cfg.db.local {
		cfg.db.setLocalDSN()
	} else {
		cfg.db.setDSN()
	}

	// Set secure cookies based on environment
	if cfg.environment == "development" {
		cfg.secureCookies = false
	} else {
		cfg.secureCookies = true
	}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(cfg.db.getDSN()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate database schema
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Tag{}, &models.Meta{}, &models.Category{})
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(
		"successfully connected to the database",
		"name", cfg.db.name,
		"host", cfg.db.host,
		"port", cfg.db.port,
	)

	// Initialize session manager
	sessionManager := scs.New()
	sessionManager.Store, err = gormstore.New(db)
	if err != nil {
		log.Fatal(err)
	}
	sessionManager.Lifetime = 24 * 7 * time.Hour

	// Initialize form decoder
	formDecoder := form.NewDecoder()

	// Initialize page and partial templates
	pageTemplates := template.CreatePageTemplates()
	partialTemplates := template.CreatePartialTemplates()

	// Create application struct
	app := &application{
		logger:           logger,
		cfg:              &cfg,
		meta:             &models.MetaModel{DB: db},
		user:             &models.UserModel{DB: db},
		post:             &models.PostModel{DB: db},
		db:               db,
		sessionManager:   sessionManager,
		formDecoder:      formDecoder,
		pageTemplates:    pageTemplates,
		partialTemplates: partialTemplates,
	}

	// Reset mostRecentPublicPost & latestPublicPosts app field
	if err := app.UpdatePostsOnAppStruct(); err != nil {
		app.logger.Error("Unable to update posts on app struct", "error", err)
	}

	// Insert meta information if necessary
	meta := models.Meta{
		Version:     version,
		Name:        "personal-site",
		LastUpdated: "2024-11-02",
		Description: "tim engle's blog",
		Author:      "Tim Engle",
		Environment: "Development",
		BuildNumber: "1",
		License:     "MIT",
	}

	fetchedMeta, err := app.meta.GetMostRecentMeta()
	if err != nil {
		app.logger.Error("Unable to fetch most recent meta", "error", err)
	}

	if fetchedMeta == nil || fetchedMeta.Version != meta.Version || fetchedMeta.LastUpdated != meta.LastUpdated {
		err = app.meta.InsertMeta(meta)
		if err != nil {
			logger.Error("Unable to insert meta", "error", err)
		} else {
			logger.Info("Successfully inserted meta")
		}
	}

	// Configure and start the HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("Successfully fetched meta", "meta", fetchedMeta)

	logger.Info("Starting the server", "port", cfg.port)
	err = srv.ListenAndServe()
	logger.Error("Server error", "error", err.Error())
	os.Exit(1)
}
