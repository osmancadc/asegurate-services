package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDatabase() (connection *sql.DB) {

	os.Setenv("DB_USER", "administrator")
	os.Setenv("DB_PASSWORD", "35Yw!8uO5v5g")
	os.Setenv("DB_HOST", "dev-asegurate.cluster-cnaioe8hvyno.us-east-1.rds.amazonaws.com")
	os.Setenv("DB_NAME", "dev_asegurate")

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
