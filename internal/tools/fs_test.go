package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file with secure random name
	tempFile, err := os.CreateTemp("", "test_file_exists_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFileName := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFileName)

	// File should exist now
	if !FileExists(tempFileName) {
		t.Error("FileExists should return true for existing file")
	}

	// Remove file and test non-existence
	os.Remove(tempFileName)
	if FileExists(tempFileName) {
		t.Error("FileExists should return false for non-existent file")
	}
}

func TestSafeFileName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "hello-world"},
		{"file/with\\slash", "file-with-slash"},
		{"file:with:colon", "file-with-colon"},
		{"file*with*asterisk", "file-with-asterisk"},
		{"file?with?question", "file-with-question"},
		{"file\"with\"quote", "file-with-quote"},
		{"file<with>brackets", "file-with-brackets"},
		{"file|with|pipe", "file-with-pipe"},
		{"normal_file.txt", "normal-file.txt"},
	}

	for _, test := range tests {
		result := SafeFileName(test.input)
		if result != test.expected {
			t.Errorf("SafeFileName(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestGetHash(t *testing.T) {
	data := []byte("test data")
	expectedHash := sha256.Sum256(data)
	expected := hex.EncodeToString(expectedHash[:])

	result := GetHash(data)
	if result != expected {
		t.Errorf("GetHash() = %q, expected %q", result, expected)
	}

	// Test empty data
	emptyResult := GetHash([]byte{})
	emptyExpected := sha256.Sum256([]byte{})
	emptyExpectedStr := hex.EncodeToString(emptyExpected[:])

	if emptyResult != emptyExpectedStr {
		t.Errorf("GetHash(empty) = %q, expected %q", emptyResult, emptyExpectedStr)
	}
}

func TestCreateDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_create_dir_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Remove it first to test creation
	os.RemoveAll(tempDir)

	result, err := CreateDir(tempDir)
	if err != nil {
		t.Errorf("CreateDir should not return error: %v", err)
	}

	if result == "" {
		t.Error("CreateDir should return directory path")
	}

	// Check if directory was created
	if !FileExists(tempDir) {
		t.Error("Directory should have been created")
	}

	// Try creating the same directory again (should not error)
	_, err = CreateDir(tempDir)
	if err != nil {
		t.Errorf("CreateDir should not return error for existing directory: %v", err)
	}
}

func TestCreateDirFromFilename(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_subdir_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	tempFile := filepath.Join(tempDir, "subdir", "test_file.txt")
	expectedDir := filepath.Dir(tempFile)

	result, err := CreateDirFromFilename(tempFile, "test_suffix")
	if err != nil {
		t.Errorf("CreateDirFromFilename should not return error: %v", err)
	}

	if result == "" {
		t.Error("CreateDirFromFilename should return directory path")
	}

	// Check if directory was created
	if !FileExists(expectedDir) {
		t.Error("Directory should have been created")
	}
}

func TestTempFileName(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_temp_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	result := TempFileName(tempDir, "test", ".txt")

	if result == "" {
		t.Error("TempFileName should not return empty string")
	}

	if !filepath.IsAbs(result) {
		t.Error("TempFileName should return absolute path")
	}

	// Check that it includes the prefix and suffix
	base := filepath.Base(result)
	if len(base) < 4 { // At least "test" + ".txt"
		t.Error("TempFileName should include prefix and suffix")
	}
}

func TestHasBOM(t *testing.T) {
	// Create temp files with secure random names
	tempFileWithBOM, err := os.CreateTemp("", "test_with_bom_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file with BOM: %v", err)
	}
	tempFileWithBOMName := tempFileWithBOM.Name()
	defer os.Remove(tempFileWithBOMName)

	tempFileWithoutBOM, err := os.CreateTemp("", "test_without_bom_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file without BOM: %v", err)
	}
	tempFileWithoutBOMName := tempFileWithoutBOM.Name()
	defer os.Remove(tempFileWithoutBOMName)

	// Test with BOM
	dataWithBOM := []byte{0xEF, 0xBB, 0xBF, 'h', 'e', 'l', 'l', 'o'}
	err = os.WriteFile(tempFileWithBOMName, dataWithBOM, 0600)
	if err != nil {
		t.Fatalf("Failed to write test file with BOM: %v", err)
	}
	tempFileWithBOM.Close()

	if !HasBOM(tempFileWithBOMName) {
		t.Error("HasBOM should return true for file with BOM")
	}

	// Test without BOM
	dataWithoutBOM := []byte{'h', 'e', 'l', 'l', 'o'}
	err = os.WriteFile(tempFileWithoutBOMName, dataWithoutBOM, 0600)
	if err != nil {
		t.Fatalf("Failed to write test file without BOM: %v", err)
	}
	tempFileWithoutBOM.Close()

	if HasBOM(tempFileWithoutBOMName) {
		t.Error("HasBOM should return false for file without BOM")
	}

	// Test non-existent file
	if HasBOM("non-existent-file.txt") {
		t.Error("HasBOM should return false for non-existent file")
	}
}
