package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
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

	//////////////////////////////////////////////
	// Initialize Config and Set ENV Variables //
	// //////////////////////////////////////////
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

	/////////////////////////////////////
	// Initialize Database and Session //
	/////////////////////////////////////

	// Configure database connection
	cfg.configDBConnection()

	// Connect and migrate database
	db, err := cfg.connectAndMigrateDB(logger)
	if err != nil {
		panic("failed to connect or migrate database")
	}

	// Initialize session manager
	sessionManager, err := initializeSessionManager(db)
	if err != nil {
		logger.Error("unable to initialize session manager", "error", err)
		panic("failed to initialize session manager")
	}

	/////////////////////////////////////
	// Initialize Application Struct ///
	////////////////////////////////////

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

///////////////////////////////
//// main function helpers ////
///////////////////////////////

// setObjectStorageURL sets the object storage URL on the config struct
func (cfg *config) setObjectStorageURL() {
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

// configDBConnection sets the database connection on the config struct
// and initializes the DSN
func (cfg *config) configDBConnection() {
	cfg.db.local = os.Getenv(DBLOCAL) == "true"
	cfg.db.host = os.Getenv(DBHOST)
	cfg.db.port = os.Getenv(DBPORT)
	cfg.db.name = os.Getenv(DBNAME)
	cfg.db.user = os.Getenv(DBUSER)
	cfg.db.password = os.Getenv(DBPASSWORD)

	if cfg.db.local {
		cfg.db.setLocalDSN()
	} else {
		cfg.db.setDSN()
	}
}

// connectToDB establishes a connection to the database using the configuration provided.
// It returns a pointer to the gorm.DB instance and an error if the connection fails.
func (cfg *config) connectToDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.db.getDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// migrateDB performs the database migration for the specified models.
// It returns an error if the migration fails.
func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Post{}, &models.Tag{}, &models.Meta{}, &models.Category{})
}

// connectAndMigrateDB connects to the database and performs the migration.
// It logs the success or failure of these operations and returns a pointer to the gorm.DB instance and an error if any operation fails.
func (cfg *config) connectAndMigrateDB(logger *slog.Logger) (*gorm.DB, error) {
	db, err := cfg.connectToDB()
	if err != nil {
		logger.Error("unable to connect to database", "error", err)
		return nil, err
	}
	err = migrateDB(db)
	if err != nil {
		logger.Error("unable to migrate database", "error", err)
		return nil, err
	}

	logger.Info(
		"successfully connected to the database",
		"name", cfg.db.name,
		"host", cfg.db.host,
		"port", cfg.db.port,
	)

	return db, nil
}

// setEnviornmentDependentVariables sets the secureCookies configuration based on the environment.
func (cfg *config) setEnviornmentDependentVariables() {
	if cfg.environment == DEVENV {
		cfg.secureCookies = false
	} else {
		cfg.secureCookies = true
	}
}

// initializeSessionManager initializes the session manager with the given database connection.
// It returns a pointer to the scs.SessionManager instance and an error if the initialization fails.
func initializeSessionManager(db *gorm.DB) (*scs.SessionManager, error) {
	var err error
	sessionManager := scs.New()
	sessionManager.Store, err = gormstore.New(db)
	if err != nil {
		return nil, err
	}
	sessionManager.Lifetime = 24 * 7 * time.Hour
	return sessionManager, nil
}

// fetchOrInsertMetaData fetches the most recent meta data from the database or inserts new meta data if necessary.
// It returns a pointer to the fetched or inserted models.Meta instance.
func (app *application) fetchOrInsertMetaData() *models.Meta {
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
			app.logger.Error("Unable to insert meta", "error", err)
		} else {
			app.logger.Info("Successfully inserted meta")
		}
	}

	return fetchedMeta
}

// ///////////////////////
// ENV Variables Helpers//
// ///////////////////////

// getEnvOrDefault returns the value of the environment variable if it exists,
// otherwise it returns the default value provided.
func getEnvOrDefault(env, defaultValue string) string {
	if v, exists := os.LookupEnv(env); exists && v != "" {
		return v
	}
	return defaultValue
}

// setHostPortEnv sets the host port configuration value. It first checks for an
// environment variable (HOSTPORT). If the environment variable is not set, it uses
// the provided default value. It logs the final port value being set.
func (cfg *config) setHostPortEnv(defaultVal string, logger *slog.Logger) {
	v := getEnvOrDefault(HOSTPORT, defaultVal)
	logger.Info(fmt.Sprintf("setting port to %s", v))
	cfg.port = v
}

// setGoEnv sets the Go environment configuration value. It first checks for an
// environment variable (GOENV). If the environment variable is not set, it uses
// the provided default value. If the value is not "production" or "development",
// it defaults to "development". It logs the final environment value being set.
func (cfg *config) setGoEnv(defaultVal string, logger *slog.Logger) {
	v := getEnvOrDefault(GOENV, defaultVal)
	// if not production or development, set to development
	if v != DEVENV && v != PRODENV {
		cfg.environment = DEVENV
	} else {
		cfg.environment = v
	}
	logger.Info(fmt.Sprintf("setting environment to %s", v))
}

func (cfg *config) setUseObjStorage(defaultVal bool, logger *slog.Logger) {
	if v, exists := os.LookupEnv(SHOULD_USE_OBJ_STORAGE_URL); exists {
		cfg.objectStorage.serveStaticObjectStorage = strings.ToLower(v) == "true"
	} else {
		cfg.objectStorage.serveStaticObjectStorage = defaultVal
	}
}
