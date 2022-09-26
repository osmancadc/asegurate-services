package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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

func GetFromDatabase(conn *sql.DB, dataType, dataValue string) (bool, string, error) {
	name := ``

	query := ``

	if dataType == `CC` {
		query = `SELECT CONCAT(name,' ',lastname) name FROM person p WHERE p.document = ?`
	} else {
		query = `SELECT CONCAT(name,' ',lastname) name FROM person p INNER JOIN user u ON p.document =u.document WHERE u.phone = ?`
	}

	results, err := conn.Query(query, dataValue)
	if err != nil {
		fmt.Printf(`GetFromDatabase(1): %s`, err.Error())
		return false, "", err
	}

	if results.Next() {
		err = results.Scan(&name)
		if err != nil {
			fmt.Printf(`GetFromDatabase(2): %s`, err.Error())
			return false, "", err
		}
		return true, name, nil
	}

	return false, "", nil
}

func GetFromProvider(dataType, dataValue string) (bool, string, error) {
	if dataType == `CC` {
		url := fmt.Sprintf(`%s/cedula?documentType=CC&documentNumber=%s`, os.Getenv("DATA_URL"), dataValue)
		bearer := "Bearer " + os.Getenv("AUTHORIZATION_TOKEN")

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf(`GetFromProvider(1) %s`, err.Error())
			return false, "", err
		}
		request.Header.Add(`Authorization`, bearer)

		client := &http.Client{}
		result, err := client.Do(request)
		if err != nil {
			fmt.Printf(`GetFromProvider(2) %s`, err.Error())
			return false, "", err
		}
		defer result.Body.Close()

		data := &Person{}

		err = json.NewDecoder(result.Body).Decode(data)
		if err != nil {
			fmt.Printf(`GetFromProvider(3) %s`, err.Error())
			return false, "", err
		}

		return true, data.Data.FullName, nil

	}
	return false, ``, nil
}
