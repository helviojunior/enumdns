package readers

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"
)

func TestNewFileReader(t *testing.T) {
	opts := &FileReaderOptions{
		DnsSuffixFile:     "test.txt",
		HostFile:          "hosts.txt",
		DnsServer:         "8.8.8.8:53",
		IgnoreNonexistent: true,
	}

	reader := NewFileReader(opts)

	if reader.Options != opts {
		t.Error("Expected Options to be set correctly")
	}

	if reader.Options.DnsSuffixFile != "test.txt" {
		t.Errorf("Expected DnsSuffixFile to be 'test.txt', got %s", reader.Options.DnsSuffixFile)
	}
}

func TestReadWordList(t *testing.T) {
	// Create temporary test file
	content := "example\ntest\nSAMPLE\n\nempty\n"
	tmpfile, err := ioutil.TempFile("", "wordlist_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test reading
	opts := &FileReaderOptions{
		HostFile: tmpfile.Name(),
	}
	reader := NewFileReader(opts)

	var wordList []string
	err = reader.ReadWordList(&wordList)
	if err != nil {
		t.Fatalf("ReadWordList failed: %v", err)
	}

	expected := []string{"example", "test", "sample", "empty"}
	if len(wordList) != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), len(wordList))
	}

	for i, word := range expected {
		if i >= len(wordList) || wordList[i] != word {
			t.Errorf("Expected word %d to be '%s', got '%s'", i, word, wordList[i])
		}
	}
}

func TestReadWordListFileNotFound(t *testing.T) {
	opts := &FileReaderOptions{
		HostFile: "nonexistent_file.txt",
	}
	reader := NewFileReader(opts)

	var wordList []string
	err := reader.ReadWordList(&wordList)
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestReadFileList(t *testing.T) {
	// Create temporary test file
	content := "line1\nline2\nLINE3\n\nline4\n"
	tmpfile, err := ioutil.TempFile("", "filelist_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test reading
	opts := &FileReaderOptions{}
	reader := NewFileReader(opts)

	var fileList []string
	err = reader.readFileList(tmpfile.Name(), &fileList)
	if err != nil {
		t.Fatalf("readFileList failed: %v", err)
	}

	expected := []string{"line1", "line2", "line3", "line4"}
	if len(fileList) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(fileList))
	}

	for i, line := range expected {
		if i >= len(fileList) || fileList[i] != line {
			t.Errorf("Expected line %d to be '%s', got '%s'", i, line, fileList[i])
		}
	}
}

func TestFileReaderOptionsWithProxy(t *testing.T) {
	proxyURL, _ := url.Parse("http://proxy.example.com:8080")
	opts := &FileReaderOptions{
		DnsSuffixFile:     "domains.txt",
		HostFile:          "hosts.txt",
		DnsServer:         "1.1.1.1:53",
		IgnoreNonexistent: false,
		ProxyUri:          proxyURL,
	}

	reader := NewFileReader(opts)

	if reader.Options.ProxyUri.String() != proxyURL.String() {
		t.Errorf("Expected ProxyUri to be %s, got %s", proxyURL.String(), reader.Options.ProxyUri.String())
	}
}

func TestReadFileListEmptyFile(t *testing.T) {
	// Create empty temporary file
	tmpfile, err := ioutil.TempFile("", "empty_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Test reading
	opts := &FileReaderOptions{}
	reader := NewFileReader(opts)

	var fileList []string
	err = reader.readFileList(tmpfile.Name(), &fileList)
	if err != nil {
		t.Fatalf("readFileList failed: %v", err)
	}

	if len(fileList) != 0 {
		t.Errorf("Expected empty list, got %d items", len(fileList))
	}
}
