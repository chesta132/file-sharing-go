package main

import (
	"context"
	"file-sharing/config"
	"file-sharing/internal/lib/filelib"
	"file-sharing/internal/services/db"
	"flag"
	"fmt"
	"log"
	"os"
)

// docker compose run --rm clear
// docker compose run --rm clear go run .clear --clear-db

// while app running
// docker compose exec app go run ./cmd/clear
// docker compose exec app go run ./cmd/clear --clear-db

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
	infoLarge := filelib.CreateFileInfo(ld)
	infoSmall := filelib.CreateFileInfo(sd)

	// creates delete log
	file, err := filelib.CreateDeleteLogFile(infoLarge, infoSmall, config.CLEAR_LOG_PATH)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Delete logs created:  %v\n", file.Name())

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
		client := db.Connect(config.DB_PATH, true)
		defer client.Close()
		client.File.Delete().ExecX(context.Background())
		fmt.Println("Successfully clear files in database")
	}
}
