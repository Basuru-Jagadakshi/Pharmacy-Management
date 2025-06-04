package main

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func ConnectDB() error {
    dbURL := "postgres://user:password@db:5432/mydb"
    var err error
    dbPool, err = pgxpool.New(context.Background(), dbURL)
    if err != nil {
        return err
    }

    err = dbPool.Ping(context.Background())
    if err != nil {
        return err
    }

    fmt.Println("Connected to PostgreSQL!")
    return nil
}

func CloseDB() {
    dbPool.Close()
}
