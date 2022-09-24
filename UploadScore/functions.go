package main

import (
	"database/sql"
	"fmt"
	"os"
)

func ConnectDatabase() (connection *sql.DB) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")

	connection, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, database))
	if err != nil {
		fmt.Printf(`Error conectando DB %s`, err.Error())
		panic(err.Error())
	}

	return
}

func UploadScore(conn *sql.DB, author, stars int, objective, comments string) error {
	query, err := conn.Prepare(`INSERT INTO dev_asegurate.score (author, objective, stars, coments) VALUES(?, ?, ?, ?); `)
	if err != nil {
		fmt.Printf("UploadScore(1) %s", err.Error())
		return err
	}

	_, err = query.Exec(author, objective, stars, comments)
	if err != nil {
		fmt.Printf("UploadScore(2) %s", err.Error())
		return err
	}

	return nil
}
