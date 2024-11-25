package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"github.com/timenglesf/personal-site/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
