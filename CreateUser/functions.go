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

func InsertPerson(conn *sql.DB, document, name, lastname string) error {
	query, err := conn.Prepare(`INSERT INTO person (document, name, lastname, score, stars, reputation, last_update) VALUES(?, ?, ?, 50, 0, 50, CURRENT_TIMESTAMP)`)
	if err != nil {
		fmt.Printf("InsertPerson(1) %s", err.Error())
		return err
	}

	query.Exec(document, name, lastname)

	return nil
}

func InsertUser(conn *sql.DB, username, email, phone, password, document, role string) (int64, error) {

	query, err := conn.Prepare(`INSERT INTO user (username, email, phone, password, document, role) VALUES (?, ?, ?, ?, ?,?)`)
	if err != nil {
		fmt.Printf("InsertUser(1) %s", err.Error())
		return -1, err
	}

	data, err := query.Exec(username, email, phone, password, document, role)
	if err != nil {
		fmt.Printf("InsertUser(2) %s", err.Error())
		return -1, err
	}

	id, _ := data.LastInsertId()
	return id, nil
}

func InsertSeller(conn *sql.DB, id int64) error {
	query, err := conn.Prepare(`INSERT INTO seller (meli_id, user_id) VALUES('', ?)`)
	if err != nil {
		fmt.Printf("InsertSeller(1) %s", err.Error())
		return err
	}

	query.Exec(id)
	return nil
}
