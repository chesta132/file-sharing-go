package filelib

import (
	"crypto/rand"
	"file-sharing/config"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

func CreateToken() string {
	const (
		length  = config.TOKEN_LENGTH
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	)

	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func GetPathname(size int64, id, filename string) string {
	return filepath.Join(GetPathBySize(size), id+GetExtension(filename))
}

func HashPassword(password string) string {
	pw := []byte(password)
	cost := bcrypt.DefaultCost

	hpw, _ := bcrypt.GenerateFromPassword(pw, cost)
	return string(hpw)
}
