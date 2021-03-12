package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	// Sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// CreateDatabase creates the database
func CreateDatabase() error {
	dbFile := "backdropGo.db"

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	dbDirPath := filepath.Join(configDir, "BackdropGo")
	dbPath := filepath.Join(dbDirPath, dbFile)

	if _, err := os.Stat(dbDirPath); os.IsNotExist(err) {
		os.Mkdir(dbDirPath, os.ModePerm)
	}

	if err != nil {
		return err
	}

	file, err := os.Create(dbPath) // Create SQLite file
	if err != nil {
		return err
	}
	file.Close()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println(fmt.Sprintf("Created database at %s", dbPath))

	return nil
}
