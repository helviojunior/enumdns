package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"encoding/base64"
	"io/ioutil"
	"crypto/sha1"
	"encoding/hex"
	"os/user"
	"errors"
)

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

	if err := os.MkdirAll(dir, 0755); err != nil {
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
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Codifica em Base64
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil

}

func GetHash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func ResolveFullPath(file_path string) (string, error) {

	if file_path == "" {
		return "", errors.New("File path cannot be empty!") 
	}

	fc := file_path[0:1]
	if fc == "~" {
		usr, err := user.Current()
	    if err != nil {
	       return "", err
	    }

		file_path = strings.Replace(file_path, "~", usr.HomeDir, 1)
		if !IsValid(file_path) {
			return "", errors.New("File path '"+ file_path + "' is not a valid path") 
		}

		return file_path, nil
	}

	currentPath, err := os.Getwd()
    if err != nil {
       return "", err
    }

	if fc == "." {
		file_path = strings.Replace(file_path, ".", currentPath, 1)
		if !IsValid(file_path) {
			return "", errors.New("File path '"+ file_path + "' is not a valid path") 
		}

		return file_path, nil
	}

	file_path = filepath.Join(currentPath, file_path)
	if !IsValid(file_path) {
		return "", errors.New("File path '"+ file_path + "' is not a valid path") 
	}

	return file_path, nil
}

func IsValid(fp string) bool {
  // Check if file already exists
  if _, err := os.Stat(fp); err == nil {
    return true
  }

  // Attempt to create it
  var d []byte
  if err := ioutil.WriteFile(fp, d, 0644); err == nil {
    os.Remove(fp) // And delete it
    return true
  }

  return false
}