package writers

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"

	"github.com/bob-reis/enumdns/internal/tools"
	"github.com/bob-reis/enumdns/pkg/models"
)

// fields in the main model to ignore
var csvExludedFields = []string{"HTML"}

// CsvWriter writes CSV files
type CsvWriter struct {
	FilePath  string
	finalPath string
}

// NewCsvWriter gets a new CsvWriter
func NewCsvWriter(destination string) (*CsvWriter, error) {
	p, err := tools.CreateFileWithDir(destination)
	if err != nil {
		return nil, err
	}

	// open the file and write the CSV headers to it
	file, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(csvHeaders()); err != nil {
		return nil, err
	}

	return &CsvWriter{
		FilePath:  destination,
		finalPath: p,
	}, nil
}

func (cw *CsvWriter) Finish() error {
	return nil
}

// Write a CSV line
func (cw *CsvWriter) Write(result *models.Result) error {

	if !result.Exists {
		return nil
	}

	file, err := os.OpenFile(cw.finalPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// get values from the result
	val := reflect.ValueOf(*result)
	numField := val.NumField()

	var values []string
	for i := 0; i < numField; i++ {
		// skip excluded fields
		if tools.SliceHasStr(csvExludedFields, val.Type().Field(i).Name) {
			continue
		}

		// skip slices
		if val.Field(i).Kind() == reflect.Slice {
			continue // Optionally skip slice fields, or handle them differently
		}

		values = append(values, fmt.Sprintf("%v", val.Field(i).Interface()))
	}

	return writer.Write(values)
}

func (cw *CsvWriter) WriteFqdn(result *models.FQDNData) error {
	return nil
}

// headers returns the headers a CSV file should have.
func csvHeaders() []string {
	val := reflect.ValueOf(models.Result{})
	numField := val.NumField()

	var fieldNames []string
	for i := 0; i < numField; i++ {
		// skip excluded fields
		if tools.SliceHasStr(csvExludedFields, val.Type().Field(i).Name) {
			continue
		}

		// skip slices
		if val.Field(i).Kind() == reflect.Slice {
			continue // Optionally skip slice fields, or handle them differently
		}

		fieldNames = append(fieldNames, val.Type().Field(i).Name)
	}

	return fieldNames
}
