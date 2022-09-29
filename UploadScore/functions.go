package main

import (
	"database/sql"
	"errors"
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

func UploadScorePhone(conn *sql.DB, author, score int, phone, comments string) error {
	objective := ""

	results, err := conn.Query(`SELECT document FROM user u WHERE u.phone = ?`, phone)
	if err != nil {
		fmt.Printf(`UploadScorePhone(1): %s`, err.Error())
		return err
	}

	if results.Next() {
		err = results.Scan(&objective)
		if err != nil {
			fmt.Printf(`UploadScorePhone(2): %s`, err.Error())
			return err
		}
		return UploadScoreDocument(conn, author, score, objective, comments)
	}

	return errors.New("no se encontró ningún usuario asociado")
}

func UploadScoreDocument(conn *sql.DB, author, score int, objective, comments string) error {
	query, err := conn.Prepare(`INSERT INTO score (author, objective, score, coments) VALUES(?, ?, ?, ?)`)
	if err != nil {
		fmt.Printf("UploadScore(1) %s", err.Error())
		return err
	}

	fmt.Println(query)
	fmt.Println(author)
	fmt.Println(objective)
	fmt.Println(score)
	fmt.Println(comments)

	_, err = query.Exec(author, objective, score, comments)
	if err != nil {
		fmt.Printf("UploadScore(2) %s", err.Error())
		return err
	}

	return nil
}
