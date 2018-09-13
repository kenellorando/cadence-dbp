package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	_ "github.com/lib/pq"
	"gopkg.in/ini.v1"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s music_dir server_dir\n")
		fmt.Println("music_dir is the directory containing music to be parsed.")
		fmt.Println("server_dir is the path to a cadence-server install whose config\n  files to use for database connection.")
		return
	}

	MUSIC_DIR := os.Args[1]

	// Check if server default-config.ini exists. Return if err
	// If that file does not exist, we are not looking at a valid cadence-server instance
	if _, err := os.Stat(filepath.Join(os.Args[2], "default-config.ini")); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s is not a valid cadence-server install directory.\n", os.Args[2])
			return
		}
	}

	// Load the configuration from cadence-server
	// By loading the override file second, it overrides the defaults file automatically.
	cfg, err := ini.LooseLoad(filepath.Join(os.Args[2], "default-config.ini"),
		filepath.Join(os.Args[2], "config.ini"))
	if err != nil {
		fmt.Println("Error during config read.")
		return
	}

	sec, err := cfg.GetSection("DEFAULT")
	if err != nil {
		fmt.Println("Error during config parse.")
		return
	}

	SQLINSERT := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5, $6)",
		sec.Key("db_table").String(), sec.Key("db_column_title").String(),
		sec.Key("db_column_album").String(), sec.Key("db_column_artist").String(),
		sec.Key("db_column_genre").String(), sec.Key("db_column_year").String(),
		sec.Key("db_column_path").String())

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
	err = filepath.Walk(MUSIC_DIR, func(path string, info os.FileInfo, err error) error {
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
