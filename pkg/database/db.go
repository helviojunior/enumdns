package database

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"
	"runtime"

	"github.com/glebarez/sqlite"
	"github.com/helviojunior/enumdns/pkg/models"

    "github.com/helviojunior/enumdns/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connection returns a Database connection based on a URI
func Connection(uri string, shouldExist, debug bool) (*gorm.DB, error) {
	var err error
	var c *gorm.DB

	db, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	var config = &gorm.Config{}
	if debug {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Error)
	}

	switch db.Scheme {
	case "sqlite":
		if shouldExist {
			if runtime.GOOS == "windows" && db.Path[0:1] == "/" {
				db.Path = db.Path[1:]
			}
			dbpath := filepath.Join(db.Host, db.Path)
			dbpath = filepath.Clean(dbpath)

			if _, err := os.Stat(dbpath); os.IsNotExist(err) {
				return nil, fmt.Errorf("sqlite database file does not exist: %s", dbpath)
			} else if err != nil {
				return nil, fmt.Errorf("error checking sqlite database file: %w", err)
			}
		}

		//config.SkipDefaultTransaction = true

		c, err = gorm.Open(sqlite.Open(db.Host+db.Path+"?cache=shared"), config)
		if err != nil {
			return nil, err
		}
		c.Exec("PRAGMA foreign_keys = ON")
		c.Exec("PRAGMA cache_size = 10000")
	case "postgres":
		c, err = gorm.Open(postgres.Open(uri), config)
		if err != nil {
			return nil, err
		}
	case "mysql":
		c, err = gorm.Open(mysql.Open(uri), config)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid db uri scheme")
	}

	// run database migrations on the connection
	if err := c.AutoMigrate(
		&Application{},
		&models.Result{},
		&models.FQDNData{},
		&models.ASN{},
		&models.ASNIpDelegate{},
	); err != nil {
		return nil, err
	}

	//Check if app name was inserted at application info table
	var count int64
    if err := c.Model(&Application{}).Count(&count).Error; err != nil {
        return nil, err
    }

    if count == 0 {
        defaultApp := Application{
            Application:  "enumdns",
            CreatedAt: time.Now(),
        }
        if err := c.Create(&defaultApp).Error; err != nil {
            return nil, err
        }
    }

    // ASN List

    if err := c.Model(&models.ASN{}).Count(&count).Error; err != nil {
        return nil, err
    }

    if count == 0 {
    	log.Warn("Populaing ASN table...")
    	if err := c.CreateInBatches(models.AsnList, 50).Error; err != nil {
            return nil, err
        }

    	if err := c.CreateInBatches(models.AsnDelagated, 50).Error; err != nil {
            return nil, err
        }
    }

	return c, nil
}

type Application struct {
	Application           string    `json:"application"`
	CreatedAt             time.Time `json:"created_at"`
}

func (Application) TableName() string {
    return "application_info"
}
