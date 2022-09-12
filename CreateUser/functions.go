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
		fmt.Printf("Error conectando DB %v", err)
		panic(err.Error())
	}

	return
}

func InsertPerson(conn *sql.DB, document, name, lastname string) error {
	query, err := conn.Prepare(`INSERT INTO person (document, name, lastname) VALUES(?, ?, ?)`)
	if err != nil {
		fmt.Println("ERROR 1.1")
		return err
	}

	query.Exec(document, name, lastname)

	return nil
}

func InsertUser(conn *sql.DB, username, email, phone, password, document, role string) (int64, error) {

	query, err := conn.Prepare(`INSERT INTO user (username, email, phone, password, document, role) VALUES (?, ?, ?, ?, ?,?)`)
	if err != nil {
		fmt.Println("ERROR 2.1")
		return -1, err
	}

	data, err := query.Exec(username, email, phone, password, document, role)
	if err != nil {
		fmt.Println("ERROR 2.2")
		fmt.Println(err.Error())
		return -1, err
	}

	id, _ := data.LastInsertId()
	return id, nil
}

func InsertSeller(conn *sql.DB, id int64) error {
	query, err := conn.Prepare(`INSERT INTO seller (score, reputation, stars, meli_id, user_id) VALUES(50, 50, 0, '', ?)`)
	if err != nil {
		fmt.Println("ERROR 4")
		return err
	}

	query.Exec(id)
	return nil
}
