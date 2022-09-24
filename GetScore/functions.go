package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
)

func ValidatePhone(conn *sql.DB, phone string) (Score, string, error) {
	score := Score{}

	results, err := conn.Query(`SELECT p.document, name, lastname, score, reputation , stars, last_update FROM user u 
								INNER JOIN person p ON u.document =p.document 
								WHERE u.phone = ?`, phone)
	if err != nil {
		fmt.Printf(`ValidatePhone(1): %s`, err.Error())
		return score, "", err
	}

	if results.Next() {
		document := ""
		err = results.Scan(&document, &score.Name, &score.Lastname, &score.Score, &score.Reputation, &score.Stars, &score.Updated)
		if err != nil {
			fmt.Printf(`ValidatePhone(2): %s`, err.Error())

			return score, "", err
		}

		return score, document, nil
	}

	return score, "", errors.New("Número de celular no encontrado")
}

func ConnectDatabase() (connection *sql.DB) {
	os.Setenv("DB_USER", "administrator")
	os.Setenv("DB_PASSWORD", "35Yw!8uO5v5g")
	os.Setenv("DB_HOST", "dev-asegurate.cluster-cnaioe8hvyno.us-east-1.rds.amazonaws.com")
	os.Setenv("DB_NAME", "dev_asegurate")

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

func GetStoredScore(conn *sql.DB, document string) (Score, bool, error) {
	score := Score{}

	results, err := conn.Query(`select name, lastname, score, reputation, stars, last_update  from person p  where p.document = ?`, document)
	if err != nil {
		fmt.Printf(`GetStoredScore(1): %s`, err.Error())
		return score, false, err
	}

	if results.Next() {
		err = results.Scan(&score.Name, &score.Lastname, &score.Score, &score.Reputation, &score.Stars, &score.Updated)
		if err != nil {
			fmt.Printf(`GetStoredScore(2): %s`, err.Error())
			return score, false, err
		}

		return score, true, nil
	}

	return score, false, nil
}

func DaysSinceLastUpdate(lastUpdate string) (int, error) {
	lastUpdated, err := time.Parse("2006-01-02 15:04:05", lastUpdate)
	if err != nil {
		fmt.Printf(`DaysSinceLastUpdate(1)  %s`, err.Error())
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
		fmt.Printf(`GetAssociatedName(1) %s`, err.Error())
		return "", "", err
	}
	defer result.Body.Close()

	data := &Person{}

	err = json.NewDecoder(result.Body).Decode(data)
	if err != nil {
		fmt.Printf(`GetAssociatedName(2) %s`, err.Error())
		return "", "", err
	}

	return data.Data.FirstName, data.Data.Lastname, nil
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

func SaveNewPerson(conn *sql.DB, score Score, document string) error {

	query, err := conn.Prepare(`INSERT INTO person (document, name, lastname, score, stars, reputation, last_update) VALUES(?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		fmt.Printf(`SaveNewPerson(1) %s`, err.Error())
		return err
	}

	data, err := query.Exec(document, score.Name, score.Lastname, score.Score, score.Stars, score.Reputation)
	if err != nil {
		fmt.Printf(`SaveNewPerson(2) %s`, err.Error())
		return err
	}

	_, err = data.LastInsertId()
	if err != nil {
		fmt.Printf(`SaveNewPerson(3) %s`, err.Error())
		return err
	}
	return nil
}

func CalculateReputation(document, documentType string) (int, error) {
	return 50, nil
}

func CalculateInternalScore(conn *sql.DB, document string) (int, error) {
	var scoreArray []int

	results, err := conn.Query(`SELECT score FROM score WHERE objective = ?`, document)
	if err != nil {
		fmt.Printf(`CalculateScore(2): %s`, err.Error())
		return -1, err
	}

	for results.Next() {
		auxScore := 0
		err = results.Scan(&auxScore)
		if err != nil {
			fmt.Printf(`CalculateScore(3): %s`, err.Error())
			return -1, err
		}
		scoreArray = append(scoreArray, auxScore)
	}

	if len(scoreArray) > 0 {
		sum := 0
		for _, i := range scoreArray {
			sum += i
		}
		return int(sum / len(scoreArray)), nil
	}

	return 50, nil

}

func CalculateScore(conn *sql.DB, document, documentType string, score Score) (Score, error) {

	elapsed, err := DaysSinceLastUpdate(score.Updated)
	if err != nil {
		return score, err
	}

	reputation := 50
	if elapsed > 7 {
		fmt.Println("The score was updated over a week ago")
		reputation, err = CalculateReputation(document, documentType)
		if err != nil {
			return score, err
		}
	}

	name, lastname, err := GetAssociatedName(document, documentType)
	if err != nil {
		fmt.Printf(`CalculateScore(1): %s`, err.Error())
		return score, err
	}

	scoreAverage, err := CalculateInternalScore(conn, document)
	if err != nil {
		return score, nil
	}

	stars := int((scoreAverage + reputation) / 40)

	score = Score{
		Name:       name,
		Lastname:   lastname,
		Score:      scoreAverage,
		Reputation: reputation,
		Stars:      stars,
	}

	return score, nil
}

func CalculateScorePhone(reqBody RequestBody, conn *sql.DB) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}

	score, document, err := ValidatePhone(conn, reqBody.Value)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	score, err = CalculateScore(conn, document, "CC", score)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.Body = GetResponseBody(score, reqBody.Value)
	response.StatusCode = http.StatusOK
	return response, nil
}

func CalculateScoreDocument(reqBody RequestBody, conn *sql.DB) (events.APIGatewayProxyResponse, error) {

	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}

	score, isStored, err := GetStoredScore(conn, reqBody.Value)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	if !isStored {
		fmt.Println("No previous score was found")

		score, err := CalculateScore(conn, reqBody.Value, reqBody.Type, score)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}

		err = SaveNewPerson(conn, score, reqBody.Value)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}

		response.Body = GetResponseBody(score, reqBody.Value)
		response.StatusCode = http.StatusOK
		return response, nil
	}

	fmt.Println("Found previous score")

	score, err = CalculateScore(conn, reqBody.Value, reqBody.Type, score)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.Body = GetResponseBody(score, reqBody.Value)
	response.StatusCode = http.StatusOK
	return response, nil
}
