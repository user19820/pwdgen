package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	// Comma import is here because all of db functionality is encapsulated
	// within the database package.
	_ "github.com/mattn/go-sqlite3"
)

const (
	insertPwd = "INSERT INTO passwords(name, pwd) VALUES (?, ?)"
	//nolint:gosec // false positive
	retrievePwd = "SELECT pwd FROM passwords WHERE name = ?"

	setupTable = `CREATE TABLE IF NOT EXISTS passwords(
		name TEXT NOT NULL PRIMARY KEY,
		pwd TEXT NOT NULL
	);`
)

var ErrAlreadyInitialized error = errors.New("pwdgen has already been initialized")

func Init() error {
	dbp, dbpErr := dbPath()
	if dbpErr != nil {
		return dbpErr
	}

	_, statErr := os.Stat(dbp)
	if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
		return statErr
	}

	if statErr == nil {
		return ErrAlreadyInitialized
	}

	if mkdirErr := os.MkdirAll(dbp, 0o750); mkdirErr != nil {
		return mkdirErr
	}

	db, ghErr := GetHandle()
	if ghErr != nil {
		return ghErr
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Println("database: closing db connection ", closeErr.Error())
		}
	}()

	_, execErr := db.Exec(setupTable)
	if execErr != nil {
		return execErr
	}

	return nil
}

func GetHandle() (*sql.DB, error) {
	dbp, err := dbPath()
	if err != nil {
		return nil, fmt.Errorf("GetHandle: get dbPath: %w", err)
	}

	db, err := sql.Open("sqlite3", filepath.Join(dbp, "pwdgen.db"))
	if err != nil {
		return nil, fmt.Errorf("GetHandle: open db: %w", err)
	}

	return db, nil
}

func Store(db *sql.DB, name, pwd string) error {
	_, err := Retrieve(db, name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("hey!")
	}

	_, err = db.Exec(insertPwd, name, pwd)
	return err
}

func Retrieve(db *sql.DB, name string) (string, error) {
	row := db.QueryRow(retrievePwd, name)

	var pwd string
	if err := row.Scan(&pwd); err != nil {
		return "", err
	}

	return pwd, nil
}

func dbPath() (string, error) {
	root, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("check if user has initialized a database; %s", err.Error())
	}

	return path.Join(root, ".local", "share", "pwdgen"), nil
}
