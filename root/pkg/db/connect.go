package db

import (
    "database/sql"
    "fmt"
    "os"

    _ "github.com/lib/pq"
)

type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
}

func LoadConfigFromEnv() Config {
    return Config{
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
    }
}

func BuildDSN(cfg Config) string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.Host,
        cfg.Port,
        cfg.User,
        cfg.Password,
        cfg.DBName,
    )
}

func Connect(driver string) func(string) (*sql.DB, error) {
    return func(dsn string) (*sql.DB, error) {
        dbConn, err := sql.Open(driver, dsn)
        if err != nil {
            return nil, err
        }
        if err := dbConn.Ping(); err != nil {
            return nil, err
        }
        return dbConn, nil
    }
}