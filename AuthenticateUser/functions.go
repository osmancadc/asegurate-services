package main

import (
	"database/sql"
	"fmt"
	"os"

	jwt "github.com/golang-jwt/jwt/v4"

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

func GenerateJWT(user User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"name":       user.Name,
		"id":         user.UserId,
		"role":       user.Role,
	})

	tokenString, err := token.SignedString([]byte("ASEGUR4TE"))
	if err != nil {
		fmt.Println("ERROR TOKEN 1")
	}

	return tokenString, nil
}
