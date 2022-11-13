package file

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// function check if file exists
func IsFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// function to copy fiile from source to dest
func CopyFile(source string, dest string) error {
	// open files r and w
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	return nil
}

// func to find number filename
func GetUniqueFileName(sourceFilePath string, destDir string) (string, error) {
	_, fileName := filepath.Split(sourceFilePath)
	ext := filepath.Ext(fileName)
	fileName = strings.TrimSuffix(filepath.Base(fileName), ext)
	for {
		id, _ := gonanoid.New(5)
		uniqueFileName := filepath.Join(destDir, fmt.Sprintf("%s_%s%s", fileName, id, ext))

		exists, err := IsFileExists(uniqueFileName)

		if err != nil {
			return "", err
		}

		if !exists {
			return uniqueFileName, nil
		}

	}
}

func WalkDirChan(dir string, fileExtension string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			r, err := regexp.MatchString(fileExtension, info.Name())
			if err != nil {
				return err
			}

			if r {
				out <- path
			}

			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}()
	return out
}

func DeletePath(path string) error {
	return os.RemoveAll(path)
}

func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func RecreateDir(path string) error {
	DeletePath(path)
	return CreateDir(path)
}
