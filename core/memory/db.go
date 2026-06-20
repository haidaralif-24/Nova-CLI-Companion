package memory

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() error {
    db, err := sql.Open("sqlite3", "./nova.db")
    if err != nil {
        return err
    }

    DB = db
    return createTables()
}

func createTables() error {
    _, err := DB.Exec(`
    CREATE TABLE IF NOT EXISTS commands (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        command TEXT,
        cwd TEXT,
        exit_code INTEGER,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `)
    return err
}
