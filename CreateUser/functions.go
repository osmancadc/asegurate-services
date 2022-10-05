package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

func CheckExistingUser(conn *sql.DB, document string) (exists bool, err error) {
	results, err := conn.Query(`SELECT user_id FROM user WHERE document =  ?`, document)
	if err != nil {
		fmt.Printf(`CheckExistingUser(1): %s`, err.Error())
		return
	}

	if results.Next() {
		fmt.Printf("CheckExistingUser(2) el usuario ya existe")
		exists = true
		return
	}

	return
}

func CheckExistingPerson(conn *sql.DB, document string) (name string, exists bool, err error) {
	results, err := conn.Query(`SELECT name FROM person p WHERE p.document = ?`, document)
	if err != nil {
		fmt.Printf(`CheckExistingPerson(1): %s`, err.Error())
		return
	}

	if results.Next() {
		fmt.Printf("CheckExistingPerson(2) la persona ya existe")
		err = results.Scan(&name)
		if err != nil {
			fmt.Printf(`CheckExistingPerson(3): %s`, err.Error())
			return
		}
		exists = true
		return
	}

	return
}

func GetPersonData(document, expirationDate string) (PersonData, error) {
	person := &Person{
		Data: PersonData{
			IsAlive: false,
		},
	}

	url := fmt.Sprintf(`%s/cedula/extra?documentType=CC&documentNumber=%s&date=%s`, os.Getenv("DATA_URL"), document, expirationDate)
	bearer := "Bearer " + os.Getenv("AUTHORIZATION_TOKEN")

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf(`GetPersonName(1) %s`, err.Error())
		return person.Data, err
	}
	request.Header.Add(`Authorization`, bearer)

	client := &http.Client{}
	result, err := client.Do(request)
	if err != nil {
		fmt.Printf(`GetPersonName(2) %s`, err.Error())
		return person.Data, err
	}
	defer result.Body.Close()

	err = json.NewDecoder(result.Body).Decode(person)
	if err != nil {
		fmt.Printf(`GetPersonName(3) %s`, err.Error())
		return person.Data, err
	}

	return person.Data, nil
}

func InsertPerson(conn *sql.DB, document, expeditionDate string) (name string, err error) {

	name, exists, err := CheckExistingPerson(conn, document)
	if err != nil {
		fmt.Printf(`InsertPerson(1): %s`, err.Error())
		return
	}

	if exists {
		fmt.Println(`InsertPerson(2): Person already exists`)
		return
	}

	personData, err := GetPersonData(document, expeditionDate)
	if err != nil {
		fmt.Printf(`InsertPerson(3): %s`, err.Error())
		return
	}
	query, err := conn.Prepare(`INSERT INTO person (document, name, lastname, score, stars, reputation, last_update) VALUES(?, ?, ?, 50, 0, 50, CURRENT_TIMESTAMP)`)
	if err != nil {
		fmt.Printf("InsertPerson(5) %s", err.Error())
		return
	}

	query.Exec(document, personData.Name, personData.Lastname)
	name = personData.Name

	return
}

func InsertUser(conn *sql.DB, email, phone, password, document, role string) (err error) {

	exists, err := CheckExistingUser(conn, document)
	if err != nil {
		fmt.Printf(`InsertUser(1): %s`, err.Error())
		return
	}

	if exists {
		fmt.Println(`InsertUser(2): User already exists`)
		err = errors.New(`user already exists`)
		return
	}

	query, err := conn.Prepare(`INSERT INTO user (email, phone, password, document, role) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Printf("InsertUser(3) %s", err.Error())
		return
	}

	query.Exec(email, phone, password, document, role)

	return
}
