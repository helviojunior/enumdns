package tools

import (
	"archive/zip"
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/helviojunior/enumdns/internal/disk"
	"github.com/helviojunior/enumdns/pkg/log"
)

func GetMimeType(s string) (string, error) {
	file, err := os.Open(s)

	if err != nil {
		return "", err
	}

	defer file.Close()

	buff := make([]byte, 512)

	// why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = file.Read(buff)

	if err != nil {
		return "", err
	}

	filetype := http.DetectContentType(buff)
	if strings.Contains(filetype, ";") {
		s1 := strings.SplitN(filetype, ";", 2)
		if s1[0] != "" && strings.Contains(s1[0], "/") {
			filetype = s1[0]
		}
	}

	return filetype, nil
}

// CreateDir creates a directory if it does not exist, returning the final
// normalized directory as a result.
func CreateDir(dir string) (string, error) {
	var err error

	if strings.HasPrefix(dir, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(homeDir, dir[1:])
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}

	return dir, nil
}

// CreateFileWithDir creates a file, relative to a directory, returning the
// final normalized path as a result.
func CreateFileWithDir(destination string) (string, error) {
	dir := filepath.Dir(destination)
	file := filepath.Base(destination)

	if file == "." || file == "/" {
		return "", fmt.Errorf("destination does not appear to be a valid file path: %s", destination)
	}

	absDir, err := CreateDir(dir)
	if err != nil {
		return "", err
	}

	absPath := filepath.Join(absDir, file)
	fileHandle, err := os.Create(absPath)
	if err != nil {
		return "", err
	}
	defer fileHandle.Close()

	return absPath, nil
}

func CreateDirFromFilename(destination string, s string) (string, error) {
	fn := SafeFileName(strings.TrimSuffix(filepath.Base(s), filepath.Ext(s)))
	if fn == "" {
		fn = "temp"
	}

	return CreateDir(filepath.Join(destination, fn))
}

// SafeFileName takes a string and returns a string safe to use as
// a file name.
func SafeFileName(s string) string {
	var builder strings.Builder

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' {
			builder.WriteRune(r)
		} else {
			builder.WriteRune('-')
		}
	}

	return builder.String()
}

func selectTempDirectory(basePath string) string {
	if basePath != "" {
		return basePath
	}

	tempDir := os.TempDir()
	di, err := disk.GetInfo(tempDir, false)
	if err != nil {
		log.Debug("Error getting disk stats", "path", tempDir, "err", err)
		return tempDir
	}

	log.Debug("Free disk space", "path", tempDir, "free", di.Free)
	if di.Free <= (5 * 1024 * 1024 * 1024) { // Less than 5GB
		currentPath, err := os.Getwd()
		if err != nil {
			log.Debug("Error getting working directory", "err", err)
			return tempDir
		}
		log.Debug("Free disk <= 5Gb, changing temp path location", "temp_path", currentPath)
		return currentPath
	}

	return tempDir
}

func TempFileName(basePath, prefix, suffix string) string {
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		log.Error("failed to generate random bytes", "err", err)
	}

	finalPath := selectTempDirectory(basePath)
	return filepath.Join(finalPath, prefix+hex.EncodeToString(randBytes)+suffix)
}

// FileExists returns true if a path exists
func FileExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

// MoveFile moves a file from a to b
func MoveFile(sourcePath, destPath string) error {
	if err := os.Rename(sourcePath, destPath); err == nil {
		return nil
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = os.Remove(sourcePath)
	if err != nil {
		return err
	}

	return nil
}

func EncodeFileToBase64(filename string) (string, error) {

	var file *os.File
	var err error

	file, err = os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Lê o conteúdo do arquivo
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Codifica em Base64
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil

}

func GetHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func RemoveFolder(path string) error {
	if path == "" {
		return nil
	}

	fi, err := os.Stat(path)

	if err != nil {
		return err
	}

	if fi.Mode().IsDir() {
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}

	} else {
		return errors.New("path is not a directory")
	}

	return nil
}

func validateZipPath(fpath, dest string) error {
	if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filepath.Base(fpath))
	}
	return nil
}

func extractZipDirectory(fpath string, mode os.FileMode) error {
	return os.MkdirAll(fpath, mode)
}

func extractZipFile(fpath string, mode os.FileMode, rc io.ReadCloser) error {
	fdir := filepath.Dir(fpath)
	if err := os.MkdirAll(fdir, mode); err != nil {
		return err
	}

	file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, rc)
	return err
}

func extractZipEntry(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	fpath := filepath.Join(dest, f.Name)
	if err := validateZipPath(fpath, dest); err != nil {
		return err
	}

	if f.FileInfo().IsDir() {
		return extractZipDirectory(fpath, f.Mode())
	}

	return extractZipFile(fpath, f.Mode(), rc)
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err := extractZipEntry(f, dest); err != nil {
			return err
		}
	}
	return nil
}

func HasBOM(fileName string) bool {
	f, err := os.Open(fileName)
	if err != nil {
		return false
	}
	defer f.Close()

	br := bufio.NewReader(f)
	r, _, err := br.ReadRune()
	if err != nil {
		return false
	}
	if r != '\uFEFF' {
		//br.UnreadRune() // Not a BOM -- put the rune back
		return false
	}

	return true
}
