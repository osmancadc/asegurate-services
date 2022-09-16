package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

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

func GetStoredScore(conn *sql.DB, document string) (Score, bool, error) {
	score := Score{}

	results, err := conn.Query(`select name, lastname, score, reputation, stars, last_update  from person p  where p.document = ?`, document)
	if err != nil {
		fmt.Println("ERROR 1 DATABASE")
		return score, false, err
	}

	if results.Next() {
		err = results.Scan(&score.Name, &score.Lastname, &score.Score, &score.Reputation, &score.Stars, &score.Updated)
		if err != nil {
			fmt.Println("error 2")
			return score, false, err
		}

		return score, true, nil
	}

	return score, false, nil
}

func DaysSinceLastUpdate(lastUpdate string) (int, error) {
	lastUpdated, err := time.Parse("2006-01-02 15:04:05", lastUpdate)
	if err != nil {
		return -1, err
	}

	return int(time.Since(lastUpdated).Hours() / 24), nil

}

func GetAssociatedName(document, documentType string) (string, string, error) {

	url := fmt.Sprintf(`%s/cedula?documentType=%s&documentNumber=%s`, os.Getenv("DATA_URL"), documentType, document)
	bearer := "Bearer " + os.Getenv("AUTHORIZATION_TOKEN")

	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add(`Authorization`, bearer)

	client := &http.Client{}
	result, err := client.Do(request)
	if err != nil {
		fmt.Println("ERROR 1 GETNAME")
		return "", "", err
	}
	defer result.Body.Close()

	data := &Person{}

	err = json.NewDecoder(result.Body).Decode(data)
	if err != nil {
		fmt.Println("ERROR 2")
		return "", "", err
	}

	return data.Data.FirstName, data.Data.Lastname, nil
}

func CalculateScore(document, documentType string) (Score, error) {

	score := Score{}

	name, lastname, _ := GetAssociatedName(document, documentType)

	min := 50
	max := 100
	score.Name = name
	score.Lastname = lastname
	score.Reputation = rand.Intn(max-min) + min
	score.Score = rand.Intn(max-min) + min
	score.Stars = rand.Intn(4) + 1

	return score, nil
}

func GetResponseBody(score Score, document string) string {
	certified := (rand.Intn(1) == 1)
	fullname := fmt.Sprintf(`%s %s`, score.Name, score.Lastname)
	profile_picture := "https://i.postimg.cc/yxNwV2Cm/user-01.png"

	return fmt.Sprintf(`{
		"name": "%s",
		"document": "%s",
		"stars": %d,
		"reputation": %d,
		"score": %d,
		"certified": %t,
		"photo": "%s"
	}`, fullname, document, score.Stars, score.Reputation, score.Score, certified, profile_picture)
}
