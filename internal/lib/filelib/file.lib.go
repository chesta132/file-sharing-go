package filelib

import (
	"errors"
	"file-sharing/config"
	"file-sharing/ent"
	"file-sharing/internal/lib/crypto"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetPathBySize(size int64) string {
	threshold := int64(config.SAVE_SPLIT) * config.MB
	if size > threshold {
		return config.LARGE_PATH
	}
	return config.SMALL_PATH
}

func CreateDir() error {
	if err := os.MkdirAll(config.SMALL_PATH, 0755); err != nil {
		log.Printf("Error creating %v directories:\n%v", config.SMALL_PATH, err)
		return err
	}

	if err := os.MkdirAll(config.LARGE_PATH, 0755); err != nil {
		log.Printf("Error creating %v directories:\n%v", config.LARGE_PATH, err)
		return err
	}
	return nil
}

func GetExtension(filename string) string {
	i := strings.LastIndex(filename, ".")
	return filename[i:]
}

func GetPathname(file *ent.File) string {
	return filepath.Join(GetPathBySize(file.FileSize), fmt.Sprintf("%v##%v", file.Token, file.FileName))
}

func ExtractFileName(path string) (string, error) {
	i := strings.LastIndex(path, "\\")
	if strings.LastIndex(path[i:], "/") != i {
		i = strings.LastIndex(path[i:], "/")
	}
	s := path[i:]
	if !strings.Contains(s, ".") {
		return "", errors.New("filename: Invalid path, file name does not have extension")
	}
	return s, nil
}

func IsDownloadable(file *ent.File, password string) (downloadable bool, cause string) {
	if file.MaxDownloads != nil && file.DownloadCount > *file.MaxDownloads {
		return false, "MAX_DOWNLOADS"
	}
	if file.Password != nil && !crypto.ComparePassword(*file.Password, password) {
		return false, "PASSWORD"
	}
	return true, ""
}
