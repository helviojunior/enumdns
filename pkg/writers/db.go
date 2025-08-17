package writers

import (
	"sync"

	"github.com/bob-reis/enumdns/pkg/database"
	"github.com/bob-reis/enumdns/pkg/models"
	"gorm.io/gorm"
)

var regThreshold = 200

// DbWriter is a Database writer
type DbWriter struct {
	URI       string
	conn      *gorm.DB
	mutex     sync.Mutex
	registers []models.Result
}

// NewDbWriter initialises a database writer
func NewDbWriter(uri string, debug bool) (*DbWriter, error) {
	c, err := database.Connection(uri, false, debug)
	if err != nil {
		return nil, err
	}
	/*
		if _, ok := c.Statement.Clauses["ON CONFLICT"]; !ok {
			c = c.Clauses(clause.OnConflict{UpdateAll: true})
		}*/

	return &DbWriter{
		URI:       uri,
		conn:      c,
		mutex:     sync.Mutex{},
		registers: []models.Result{},
	}, nil
}

// Write results to the database
func (dw *DbWriter) Write(result *models.Result) error {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()
	var err error

	if !result.Exists {
		dw.registers = append(dw.registers, *result)
		if len(dw.registers) >= regThreshold {
			err = dw.conn.CreateInBatches(dw.registers, 50).Error
			dw.registers = []models.Result{}
		}
	} else {
		err = dw.conn.CreateInBatches(result, 50).Error

		//err = dw.conn.Table("results").CreateInBatches( []models.Result{ *result }, 50).Error

		fqdn := result.ToFqdn()
		if fqdn != nil {
			// Not call WriteFqdn function because it will cause an deadlock at mutex
			err1 := dw.conn.CreateInBatches(fqdn, 50).Error
			if err1 != nil && err == nil {
				err = err1
			}
		}

	}

	return err
}

func (dw *DbWriter) WriteFqdn(fqdn *models.FQDNData) error {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	return dw.conn.Create(fqdn).Error
}

func (dw *DbWriter) Finish() error {
	return nil
}
