package filelib

import (
	"file-sharing/config"
	"file-sharing/ent"
	"file-sharing/internal/lib/crypto"
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

func GetPathname(size int64, id, filename string) string {
	return filepath.Join(GetPathBySize(size), id+GetExtension(filename))
}

func IsDownloadable(file *ent.File, password string) bool {
	if file.MaxDownloads != nil && file.DownloadCount > *file.MaxDownloads {
		return false
	}
	if file.Password != nil && !crypto.ComparePassword(*file.Password, password) {
		return false
	}
	return true
}
