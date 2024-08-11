package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

func New() (db *sql.DB, connector *libsql.Connector, dir string, err error) {
	dbName := os.Getenv("TURSO_LOCAL_DB")
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	encryptionKey := os.Getenv("TURSO_ENCRYPTION_KEY")

	dir, err = os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}

	dbPath := filepath.Join(dir, dbName)

	connector, err = libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
		libsql.WithEncryption(encryptionKey),
	)

	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}

	db = sql.OpenDB(connector)

	return db, connector, dir, nil
}
