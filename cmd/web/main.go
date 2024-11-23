package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"gorm.io/gorm"

	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/ui/template"
)

// ENUMS
const (
	sessionUserId  = "authenticatedUserID"
	sessionIsAdmin = "isAdmin"
	// DB env variable keys
	DBLOCAL    = "DBLOCAL"
	DBHOST     = "DBHOST"
	DBPORT     = "DBPORT"
	DBNAME     = "DBNAME"
	DBUSER     = "DBUSER"
	DBPASSWORD = "DBPASSWORD"
	// Server Environment variable keys
	HOSTPORT = "PORT"
	GOENV    = "GOENV"
	DEVENV   = "development"
	PRODENV  = "production"

	SHOULD_USE_OBJ_STORAGE_URL = "USE_OBJ_STORAGE"
)

var version = "1.0.3"

// application struct
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

// config struct
type config struct {
	port          string
	environment   string
	db            psqlConfig
	objectStorage objectStorageConfig
	secureCookies bool
}

type objectStorageConfig struct {
	objectStorageURL         string
	serveStaticObjectStorage bool
}

// database config struct
type psqlConfig struct {
	local    bool
	host     string
	port     string
	user     string
	password string
	name     string
	dsn      string
}

// database config methods
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

////////////////////////////////////////
////////// MAIN FUNCTION ///////////////
////////////////////////////////////////

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	///////// Initialize Config and Set ENV Variables /////////
	var cfg config
	portFlag := flag.String("port", "8080", "HTTP port")
	envFlag := flag.String("env", "development", "Environment")
	// Parse command-line flags
	objStorageFlag := flag.Bool(
		"object-storage",
		false,
		"Serve static files from object storage",
	)

	flag.Parse()

	// Check if ENV variables exist, if so override the flag values
	// set port variable
	cfg.setHostPortEnv(*portFlag, logger)
	// set development enviornment variable
	cfg.setGoEnv(*envFlag, logger)
	// set object storage flag
	cfg.setUseObjStorage(*objStorageFlag, logger)

	// Set secure cookies based on environment
	cfg.setEnviornmentDependentVariables()

	// Configure object storage if enabled
	if cfg.objectStorage.serveStaticObjectStorage {
		cfg.setObjectStorageURL()
	}

	////// Initialize Database and Session //////

	// Configure database connection
	cfg.configDBConnection()

	// Connect and migrate database
	db, err := cfg.connectAndMigrateDB(logger)
	if err != nil {
		logger.Error("failed to connect or migrate database", "error", err)
		panic("failed to connect or migrate database")
	}

	// Initialize session manager
	sessionManager, err := initializeSessionManager(db)
	if err != nil {
		logger.Error("unable to initialize session manager", "error", err)
		panic("failed to initialize session manager")
	}

	/////// Initialize Application Struct ////////

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

	fetchedMeta := app.fetchOrInsertMetaData()

	/////// Configure and start the HTTP server ///////
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
