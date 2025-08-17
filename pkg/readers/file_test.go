package readers

import (
	"net/url"
	"os"
	"testing"
)

// Helper function to create temporary file with content for testing
func createTempFileWithContent(t *testing.T, pattern, content string) (string, func()) {
	tmpfile, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()

	if content != "" {
		if _, err := tmpfile.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}

	fileName := tmpfile.Name()
	cleanup := func() { os.Remove(fileName) }
	return fileName, cleanup
}

// Helper function to create default FileReaderOptions
func createDefaultFileReaderOptions() *FileReaderOptions {
	return &FileReaderOptions{
		DnsServer:         "8.8.8.8:53",
		IgnoreNonexistent: true,
	}
}

// Helper function to verify string slice matches expected values
func verifyStringSlice(t *testing.T, actual, expected []string, itemType string) {
	if len(actual) != len(expected) {
		t.Errorf("Expected %d %s, got %d", len(expected), itemType, len(actual))
	}

	for i, expectedItem := range expected {
		if i >= len(actual) || actual[i] != expectedItem {
			t.Errorf("Expected %s %d to be '%s', got '%s'", itemType, i, expectedItem, actual[i])
		}
	}
}

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
	content := "example\ntest\nSAMPLE\n\nempty\n"
	tmpFileName, cleanup := createTempFileWithContent(t, "wordlist_test", content)
	defer cleanup()

	opts := createDefaultFileReaderOptions()
	opts.HostFile = tmpFileName
	reader := NewFileReader(opts)

	var wordList []string
	err := reader.ReadWordList(&wordList)
	if err != nil {
		t.Fatalf("ReadWordList failed: %v", err)
	}

	expected := []string{"example", "test", "sample", "empty"}
	verifyStringSlice(t, wordList, expected, "words")
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
	content := "line1\nline2\nLINE3\n\nline4\n"
	tmpFileName, cleanup := createTempFileWithContent(t, "filelist_test", content)
	defer cleanup()

	opts := createDefaultFileReaderOptions()
	reader := NewFileReader(opts)

	var fileList []string
	err := reader.readFileList(tmpFileName, &fileList)
	if err != nil {
		t.Fatalf("readFileList failed: %v", err)
	}

	expected := []string{"line1", "line2", "line3", "line4"}
	verifyStringSlice(t, fileList, expected, "lines")
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
	tmpFileName, cleanup := createTempFileWithContent(t, "empty_test", "")
	defer cleanup()

	opts := createDefaultFileReaderOptions()
	reader := NewFileReader(opts)

	var fileList []string
	err := reader.readFileList(tmpFileName, &fileList)
	if err != nil {
		t.Fatalf("readFileList failed: %v", err)
	}

	if len(fileList) != 0 {
		t.Errorf("Expected empty list, got %d items", len(fileList))
	}
}

func TestReadDnsList(t *testing.T) {
	content := "com\norg\nnet\n.edu\n\ninfo\n"
	tmpFileName, cleanup := createTempFileWithContent(t, "dns_suffix_test", content)
	defer cleanup()

	opts := createDefaultFileReaderOptions()
	opts.DnsSuffixFile = tmpFileName
	reader := NewFileReader(opts)

	var dnsList []string
	err := reader.ReadDnsList(&dnsList)
	if err != nil {
		t.Fatalf("ReadDnsList failed: %v", err)
	}

	expected := []string{"com.", "org.", "net.", "edu.", "info."}
	verifyStringSlice(t, dnsList, expected, "DNS suffixes")
}

func TestReadDnsListFileNotFound(t *testing.T) {
	opts := &FileReaderOptions{
		DnsSuffixFile:     "nonexistent_dns_file.txt",
		DnsServer:         "8.8.8.8:53",
		IgnoreNonexistent: true,
	}
	reader := NewFileReader(opts)

	var dnsList []string
	err := reader.ReadDnsList(&dnsList)
	if err == nil {
		t.Error("Expected error for nonexistent DNS file")
	}
}
