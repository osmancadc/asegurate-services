package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var ConnectDatabase = func() (connection *sql.DB, err error) {

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")

	connection, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, database))
	if err != nil {
		fmt.Printf("ConnectDatabase(1) %s", err.Error())
		return nil, err
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

	_, err = query.Exec(author, objective, score, comments)
	if err != nil {
		fmt.Printf("UploadScore(2) %s", err.Error())
		return err
	}

	return nil
}

func GetAuthorId(conn *sql.DB, document string) (int, error) {

	id := 0

	results, err := conn.Query(`SELECT user_id FROM user u WHERE u.document = ?`, document)
	if err != nil {
		fmt.Printf(`GetFromDatabase(1): %s`, err.Error())
		return -1, err
	}

	if results.Next() {
		err = results.Scan(&id)
		if err != nil {
			fmt.Printf(`GetFromDatabase(2): %s`, err.Error())
			return -1, err
		}
		return id, nil
	}

	return -1, errors.New("no user found")

}
