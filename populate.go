package main

import (
	"os"
	"fmt"
	//"time"
	"path/filepath"
	//"database/sql"
	_ "github.com/lib/pq"
)

const (
	DB_USER = "cadence"
	DB_NAME = "cadence"
	MUSIC_DIR = "/home/ken/cadence_testdir/"
)

func main() {
	err := filepath.Walk(MUSIC_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("Visited file: %q\n", path)
		return nil
	})

	if err != nil {
		fmt.Printf("Error in %q: %v\n", MUSIC_DIR, err)
	}
}