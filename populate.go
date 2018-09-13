package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	_ "github.com/lib/pq"
)

const (
	DB_USER   = "postgres"
	DB_NAME   = "cadence"
	MUSIC_DIR = "/home/ken/cadence_testdir/"
	SQLINSERT = `INSERT INTO cadence (title, album, artist, genre, year, path) VALUES ($1, $2, $3, $4, $5, $6)`
)

func main() {
	db, err = sql.Open()
	var extensions = [...]string{
		".mp3",
		".m4a",
		".ogg",
		".flac"}

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
		//fmt.Printf("Visited file: %q\n", path)

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-music files
		music := false
		for _, ext := range extensions {
			if strings.HasSuffix(path, ext) {
				music = true
				break
			}
		}
		if !music {
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

		fmt.Printf("title %q, album %q, artist %q, genre %q, year %d.\n",
			tags.Title(),
			tags.Album(),
			tags.Artist(),
			tags.Genre(),
			tags.Year())

		// Todo: connect to database

		// Insert into database
		_, err = db.Exec(SQLINSERT, tags.Title(), tags.Album(), tags.Artist(), tags.Genre(), tags.Year(), path)
		if err != nil {
			panic(err)
		}

		// Close the file
		file.Close()
		return nil
	})

	if err != nil {
		fmt.Printf("Error in %q: %v\n", MUSIC_DIR, err)
	}
}
