package main

import (
	"database/sql"
	"fmt"
	"os"

	jwt "github.com/golang-jwt/jwt/v4"

	_ "github.com/go-sql-driver/mysql"
)

var ConnectDatabase = func() (connection *sql.DB, err error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")

	connection, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, database))
	if err != nil {
		fmt.Printf("Error conectando DB %v", err)
		return nil, err
	}

	return
}

func GetUserData(conn *sql.DB, data RequestBody) (found bool, user User, err error) {
	found = false
	results, err := conn.Query(`SELECT u.document id, CONCAT(p.name," ",p.lastname) name, u.role FROM user u
												INNER JOIN person p on u.document = p.document
												WHERE u.document = ? and u.password = ?`, data.Document, data.Password)
	if err != nil {
		fmt.Printf(`GetUserData(1): %s`, err.Error())
		return
	}

	if results.Next() {
		found = true
		err = results.Scan(&user.UserId, &user.Name, &user.Role)
		if err != nil {
			fmt.Printf(`GetUserData(2): %s`, err.Error())
			return
		}
	}

	return
}

func GenerateJWT(user User) (token string, err error) {
	tokenData := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"name":       user.Name,
		"id":         user.UserId,
		"role":       user.Role,
	})

	token, err = tokenData.SignedString([]byte("ASEGUR4TE"))
	if err != nil {
		fmt.Printf("GenerateJWT(1): %s", err.Error())
		return
	}

	return
}
