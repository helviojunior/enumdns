package writers

import (
	"os"
	"testing"
	"time"

	"github.com/bob-reis/enumdns/pkg/models"
)

// Helper function to create temporary file for testing
func createTempFile(t *testing.T, pattern string) (string, func()) {
	tempFile, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFileName := tempFile.Name()
	tempFile.Close()
	cleanup := func() { os.Remove(tempFileName) }
	return tempFileName, cleanup
}

// Helper function to create sample Result for testing
func createSampleResult() *models.Result {
	return &models.Result{
		FQDN:     "example.com",
		RType:    "A",
		IPv4:     "1.2.3.4",
		ProbedAt: time.Now(),
	}
}

// Helper function to create sample FQDNData for testing
func createSampleFQDN() *models.FQDNData {
	return &models.FQDNData{
		FQDN:   "example.com",
		Source: "test",
	}
}

func TestNewStdoutWriter(t *testing.T) {
	writer, err := NewStdoutWriter()
	if err != nil {
		t.Fatalf("Failed to create stdout writer: %v", err)
	}
	if writer == nil {
		t.Error("Stdout writer should not be nil")
	}
}

func TestStdoutWriterWrite(t *testing.T) {
	writer, err := NewStdoutWriter()
	if err != nil {
		t.Fatalf("Failed to create stdout writer: %v", err)
	}

	result := createSampleResult()

	err = writer.Write(result)
	if err != nil {
		t.Errorf("Write should not return error: %v", err)
	}
}

func TestStdoutWriterFinish(t *testing.T) {
	writer, err := NewStdoutWriter()
	if err != nil {
		t.Fatalf("Failed to create stdout writer: %v", err)
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Finish should not return error: %v", err)
	}
}

func TestNewTextWriter(t *testing.T) {
	tempFileName, cleanup := createTempFile(t, "test_output_*.txt")
	defer cleanup()

	writer, err := NewTextWriter(tempFileName)
	if err != nil {
		t.Fatalf("Failed to create text writer: %v", err)
	}
	if writer == nil {
		t.Error("Text writer should not be nil")
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Finish should not return error: %v", err)
	}
}

func TestTextWriterWrite(t *testing.T) {
	tempFileName, cleanup := createTempFile(t, "test_output_*.txt")
	defer cleanup()

	writer, err := NewTextWriter(tempFileName)
	if err != nil {
		t.Fatalf("Failed to create text writer: %v", err)
	}

	result := createSampleResult()

	err = writer.Write(result)
	if err != nil {
		t.Errorf("Write should not return error: %v", err)
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Finish should not return error: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(tempFileName); os.IsNotExist(err) {
		t.Error("Output file should have been created")
	}
}

func TestTextWriterWriteFqdn(t *testing.T) {
	tempFileName, cleanup := createTempFile(t, "test_fqdn_output_*.txt")
	defer cleanup()

	writer, err := NewTextWriter(tempFileName)
	if err != nil {
		t.Fatalf("Failed to create text writer: %v", err)
	}

	fqdn := createSampleFQDN()

	err = writer.WriteFqdn(fqdn)
	if err != nil {
		t.Errorf("WriteFqdn should not return error: %v", err)
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Finish should not return error: %v", err)
	}
}

func TestNewJsonWriter(t *testing.T) {
	tempFileName, cleanup := createTempFile(t, "test_output_*.jsonl")
	defer cleanup()

	writer, err := NewJsonWriter(tempFileName)
	if err != nil {
		t.Fatalf("Failed to create JSON writer: %v", err)
	}
	if writer == nil {
		t.Error("JSON writer should not be nil")
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Finish should not return error: %v", err)
	}
}

func TestJsonWriterWrite(t *testing.T) {
	tempFileName, cleanup := createTempFile(t, "test_output_*.jsonl")
	defer cleanup()

	writer, err := NewJsonWriter(tempFileName)
	if err != nil {
		t.Fatalf("Failed to create JSON writer: %v", err)
	}

	result := createSampleResult()

	err = writer.Write(result)
	if err != nil {
		t.Errorf("Write should not return error: %v", err)
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Finish should not return error: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(tempFileName); os.IsNotExist(err) {
		t.Error("Output file should have been created")
	}
}

func TestNewCsvWriter(t *testing.T) {
	tempFileName, cleanup := createTempFile(t, "test_output_*.csv")
	defer cleanup()

	writer, err := NewCsvWriter(tempFileName)
	if err != nil {
		t.Fatalf("Failed to create CSV writer: %v", err)
	}
	if writer == nil {
		t.Error("CSV writer should not be nil")
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Finish should not return error: %v", err)
	}
}

func TestCsvWriterWrite(t *testing.T) {
	tempFileName, cleanup := createTempFile(t, "test_output_*.csv")
	defer cleanup()

	writer, err := NewCsvWriter(tempFileName)
	if err != nil {
		t.Fatalf("Failed to create CSV writer: %v", err)
	}

	result := createSampleResult()

	err = writer.Write(result)
	if err != nil {
		t.Errorf("Write should not return error: %v", err)
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Finish should not return error: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(tempFileName); os.IsNotExist(err) {
		t.Error("Output file should have been created")
	}
}

func TestNewNoneWriter(t *testing.T) {
	writer, err := NewNoneWriter()
	if err != nil {
		t.Fatalf("Failed to create none writer: %v", err)
	}
	if writer == nil {
		t.Error("None writer should not be nil")
	}
}

func TestNoneWriterOperations(t *testing.T) {
	writer, err := NewNoneWriter()
	if err != nil {
		t.Fatalf("Failed to create none writer: %v", err)
	}

	result := createSampleResult()

	err = writer.Write(result)
	if err != nil {
		t.Errorf("Write should not return error: %v", err)
	}

	fqdn := createSampleFQDN()

	err = writer.WriteFqdn(fqdn)
	if err != nil {
		t.Errorf("WriteFqdn should not return error: %v", err)
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Finish should not return error: %v", err)
	}
}

func TestNewMemoryWriter(t *testing.T) {
	// Test valid slots
	writer, err := NewMemoryWriter(10)
	if err != nil {
		t.Fatalf("Failed to create memory writer: %v", err)
	}
	if writer == nil {
		t.Error("Memory writer should not be nil")
	}

	// Test invalid slots
	_, err = NewMemoryWriter(0)
	if err == nil {
		t.Error("Expected error for zero slots")
	}

	_, err = NewMemoryWriter(-1)
	if err == nil {
		t.Error("Expected error for negative slots")
	}
}

func TestMemoryWriterOperations(t *testing.T) {
	writer, err := NewMemoryWriter(3)
	if err != nil {
		t.Fatalf("Failed to create memory writer: %v", err)
	}

	// Test empty state
	if writer.GetLatest() != nil {
		t.Error("GetLatest should return nil for empty writer")
	}
	if writer.GetFirst() != nil {
		t.Error("GetFirst should return nil for empty writer")
	}
	if len(writer.GetAllResults()) != 0 {
		t.Error("GetAllResults should return empty slice for empty writer")
	}

	// Add first result
	result1 := &models.Result{
		FQDN:  "example1.com",
		RType: "A",
		IPv4:  "1.2.3.4",
	}
	err = writer.Write(result1)
	if err != nil {
		t.Errorf("Write should not return error: %v", err)
	}

	if writer.GetLatest().FQDN != "example1.com" {
		t.Error("GetLatest should return first result")
	}
	if writer.GetFirst().FQDN != "example1.com" {
		t.Error("GetFirst should return first result")
	}

	// Add second result
	result2 := &models.Result{
		FQDN:  "example2.com",
		RType: "A",
		IPv4:  "2.3.4.5",
	}
	err = writer.Write(result2)
	if err != nil {
		t.Errorf("Write should not return error: %v", err)
	}

	if writer.GetLatest().FQDN != "example2.com" {
		t.Error("GetLatest should return second result")
	}
	if writer.GetFirst().FQDN != "example1.com" {
		t.Error("GetFirst should still return first result")
	}

	// Add third result
	result3 := &models.Result{
		FQDN:  "example3.com",
		RType: "A",
		IPv4:  "3.4.5.6",
	}
	err = writer.Write(result3)
	if err != nil {
		t.Errorf("Write should not return error: %v", err)
	}

	if len(writer.GetAllResults()) != 3 {
		t.Errorf("Expected 3 results, got %d", len(writer.GetAllResults()))
	}

	// Add fourth result (should overflow and remove first)
	result4 := &models.Result{
		FQDN:  "example4.com",
		RType: "A",
		IPv4:  "4.5.6.7",
	}
	err = writer.Write(result4)
	if err != nil {
		t.Errorf("Write should not return error: %v", err)
	}

	if len(writer.GetAllResults()) != 3 {
		t.Errorf("Expected 3 results after overflow, got %d", len(writer.GetAllResults()))
	}
	if writer.GetFirst().FQDN != "example2.com" {
		t.Error("GetFirst should return second result after overflow")
	}
	if writer.GetLatest().FQDN != "example4.com" {
		t.Error("GetLatest should return fourth result")
	}

	// Test WriteFqdn and Finish
	fqdn := &models.FQDNData{
		FQDN:   "test.com",
		Source: "test",
	}
	err = writer.WriteFqdn(fqdn)
	if err != nil {
		t.Errorf("WriteFqdn should not return error: %v", err)
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Finish should not return error: %v", err)
	}
}
