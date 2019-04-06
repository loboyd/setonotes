package postgres

import (
    "log"
    "database/sql"

    "github.com/setonotes/pkg/config"

    _ "github.com/lib/pq" // postgres drivers
)

// Repository defines a wrapper for a postgres database pool
type Repository struct {
    // this should eventually be changed back to `db` to avoid exporting
    DB *sql.DB
}

func New(c *config.Config) (*Repository, error) {
    log.Println("creating new Postgres repository...")

    // set postgres connection string
    psqlStr :=
        "host="     + c.DBHost +
        " port="     + c.DBPort +
        " user="     + c.DBUser +
        " password=" + c.DBPass +
        " dbname="   + c.DBName

    // open postgres connection
    db, err := sql.Open("postgres", psqlStr)
    if err != nil {
        log.Println("failed to open postgres connection")
        return nil, err
    }
    log.Println("successfully created postgres database")

    r := &Repository{DB: db}

    return r, nil
}
