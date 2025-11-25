package db

import (
    "database/sql"
    "fmt"
    "os"

    _ "github.com/lib/pq"
)

func GetDSNFromEnv() string {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    pass := os.Getenv("DB_PASSWORD")
    name := os.Getenv("DB_NAME")
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)
}

func Connect(dsn string) (*sql.DB, error) {
    dbConn, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    if err := dbConn.Ping(); err != nil {
        return nil, err
    }
    return dbConn, nil
}