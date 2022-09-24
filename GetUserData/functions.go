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
		fmt.Printf(`Error conectando DB %s`, err.Error())
		panic(err.Error())
	}

	return
}

func GetUserData(document string, conn *sql.DB) (User, error) {
	user := User{
		Document: document,
	}

	results, err := conn.Query(`SELECT CONCAT(name,' ',lastname) name, email, phone, photo FROM person p 
								INNER JOIN user u ON p.document = u.document
								WHERE p.document = ?`, document)
	if err != nil {
		fmt.Printf(`GetStoredScore(1): %s`, err.Error())
		return user, err
	}

	if results.Next() {
		err = results.Scan(&user.Name, &user.Email, &user.Phone, &user.Photo)
		if err != nil {
			fmt.Printf(`GetStoredScore(2): %s`, err.Error())
			return user, err
		}
	}

	return user, nil
}
