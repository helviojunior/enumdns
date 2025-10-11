package database

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/helviojunior/enumdns/pkg/models"
	"github.com/glebarez/sqlite"

	"github.com/helviojunior/enumdns/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func createGormConfig(debug bool) *gorm.Config {
	config := &gorm.Config{}
	if debug {
		// SECURITY: Debug logging enabled - should be disabled in production
		// This may log sensitive database queries and parameters
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Error)
	}
	return config
}

func validateSQLiteFile(db *url.URL, shouldExist bool) error {
	if !shouldExist {
		return nil
	}

	if runtime.GOOS == "windows" && db.Path[0:1] == "/" {
		db.Path = db.Path[1:]
	}
	dbpath := filepath.Join(db.Host, db.Path)
	dbpath = filepath.Clean(dbpath)

	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		return fmt.Errorf("sqlite database file does not exist: %s", dbpath)
	} else if err != nil {
		return fmt.Errorf("error checking sqlite database file: %w", err)
	}
	return nil
}

func openDatabaseConnection(uri string, db *url.URL, config *gorm.Config) (*gorm.DB, error) {
	var c *gorm.DB
	var err error

	switch db.Scheme {
	case "sqlite":
		c, err = gorm.Open(sqlite.Open(db.Host+db.Path+"?cache=shared"), config)
		if err != nil {
			return nil, err
		}
		c.Exec("PRAGMA foreign_keys = ON")
		c.Exec("PRAGMA cache_size = 10000")
	case "postgres":
		c, err = gorm.Open(postgres.Open(uri), config)
	case "mysql":
		c, err = gorm.Open(mysql.Open(uri), config)
	default:
		return nil, errors.New("invalid db uri scheme")
	}
	return c, err
}

func runMigrations(c *gorm.DB) error {
	return c.AutoMigrate(
		&Application{},
		&models.Result{},
		&models.FQDNData{},
		&models.ASN{},
		&models.ASNIpDelegate{},
	)
}

func initializeDefaultData(c *gorm.DB) error {
	var count int64

	// Initialize application info
	if err := c.Model(&Application{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		defaultApp := Application{
			Application: "enumdns",
			CreatedAt:   time.Now(),
		}
		if err := c.Create(&defaultApp).Error; err != nil {
			return err
		}
	}

	// Initialize ASN data
	if err := c.Model(&models.ASN{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		log.Warn("Populaing ASN table...")
		if err := c.CreateInBatches(models.ASNList, 50).Error; err != nil {
			return err
		}
		if err := c.CreateInBatches(models.ASNDelegated, 50).Error; err != nil {
			return err
		}
	}

	return nil
}

// Connection returns a Database connection based on a URI
func Connection(uri string, shouldExist, debug bool) (*gorm.DB, error) {
	db, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	config := createGormConfig(debug)

	if db.Scheme == "sqlite" {
		if err := validateSQLiteFile(db, shouldExist); err != nil {
			return nil, err
		}
	}

	c, err := openDatabaseConnection(uri, db, config)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(c); err != nil {
		return nil, err
	}

	if err := initializeDefaultData(c); err != nil {
		return nil, err
	}

	return c, nil
}

type Application struct {
	Application string    `json:"application"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Application) TableName() string {
	return "application_info"
}
