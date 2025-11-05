package main

import (
	"context"
	"encoding/json"
	"file-sharing/config"
	"file-sharing/internal/lib/filelib"
	"file-sharing/internal/services/db"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// docker compose run --rm clear
// docker compose run --rm clear go run .clear --clear-db

// while app running
// docker compose exec app go run ./cmd/clear
// docker compose exec app go run ./cmd/clear --clear-db

const dateFormat = "2006-01-02_15-04-05"

type fileInfo struct {
	Name         string `json:"fileName"`
	Types        string `json:"type"`
	Extension    string `json:"extension,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
	Size         string `json:"size,omitempty"`
}

type deleteLog struct {
	Large []*fileInfo `json:"largeSize"`
	Small []*fileInfo `json:"smallSize"`
}

func getType(f os.DirEntry) string {
	if f.IsDir() {
		return "directory"
	}
	return "file"
}

func createFileInfo(entries []os.DirEntry) []*fileInfo {
	result := []*fileInfo{}
	for _, f := range entries {
		i, err := f.Info()
		lastModif := ""
		size := ""
		if err == nil {
			lastModif = i.ModTime().Format(dateFormat)
			if !f.IsDir() {
				size = fmt.Sprintf("%.2fMB", float64(i.Size())/float64(config.MB))
			}
		}
		result = append(result, &fileInfo{
			Name:         f.Name(),
			Types:        getType(f),
			Extension:    filelib.GetExtension(f.Name()),
			LastModified: lastModif,
			Size:         size,
		})
	}
	return result
}

func main() {
	clearDB := flag.Bool("clear-db", false, "Clear database")
	flag.Parse()

	// prevent error on read dir
	if filelib.CreateDir() != nil {
		return
	}

	// read large/small directories
	ld, err := os.ReadDir(config.LARGE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	sd, err := os.ReadDir(config.SMALL_PATH)
	if err != nil {
		log.Fatal(err)
	}

	// create file info for delete log
	infoLarge := createFileInfo(ld)
	infoSmall := createFileInfo(sd)

	// creates delete log
	now := time.Now().Format(dateFormat)
	filename := now + ".json"

	os.MkdirAll(config.CLEAR_LOG_PATH, 0755)
	file, err := os.Create(filepath.Join(config.CLEAR_LOG_PATH, filename))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dl := deleteLog{infoLarge, infoSmall}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.Encode(dl)

	fmt.Printf("Delete logs created:  %v\n", filename)

	// clear files
	err = os.RemoveAll(config.LARGE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll(config.SMALL_PATH)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully clear all files")

	// clear db
	if *clearDB {
		client := db.Connect(filepath.Join(config.DB_PATH, "data.db"), true)
		defer client.Close()
		client.File.Delete().ExecX(context.Background())
		fmt.Println("Successfully clear files in database")
	}
}
