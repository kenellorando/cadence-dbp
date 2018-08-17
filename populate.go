package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	_ "github.com/lib/pq"
)

const (
	DB_USER   = "cadence"
	DB_NAME   = "cadence"
	MUSIC_DIR = "/home/ken/cadence_testdir/"
)

func main() {
	// Check if MUSIC_DIR exists. Return if err
	if _, err := os.Stat(MUSIC_DIR); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Music directory not found.\n")
			return
		}
	}

	// Recursive walk on MUSIC_DIR's contents
	err := filepath.Walk(MUSIC_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("Visited file: %q\n", path)

		if info.IsDir() {
			return nil
		}

		// Open a file for reading
		file, e := os.Open(path)
		if e != nil {
			return e
		}

		// Read metadata from the file
		tags, er := tag.ReadFrom(file)
		if er != nil {
			return er
		}

		fmt.Printf("title %q, album %q, artist %q, genre %q.\n",
			tags.Title(),
			tags.Album(),
			tags.Artist(),
			tags.Genre())

		// Close the file
		file.Close()
		return nil
	})

	if err != nil {
		fmt.Printf("Error in %q: %v\n", MUSIC_DIR, err)
	}
}
