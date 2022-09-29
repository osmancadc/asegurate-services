package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDatabase() (connection *sql.DB) {

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")

	connection, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, database))
	if err != nil {
		fmt.Printf("ConnectDatabase(1) %s", err.Error())
		panic(err.Error())
	}

	return
}
func UploadScore(conn *sql.DB, author, score int, objective, comments string) error {

	fmt.Printf("New score: \nauthor: %d \nobjective: %s \nscore: %d comments:%s", author, objective, score, comments)

	query, err := conn.Prepare(`INSERT INTO score (author, objective, score, coments) VALUES(?, ?, ?, ?)`)
	if err != nil {
		fmt.Printf("UploadScore(1) %s", err.Error())
		return err
	}

	_, err = query.Exec(author, objective, score, comments)
	if err != nil {
		fmt.Printf("UploadScore(2) %s", err.Error())
		return err
	}

	return nil
}
