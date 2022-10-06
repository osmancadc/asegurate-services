package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDatabase() (connection *sql.DB, err error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")

	connection, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, database))
	if err != nil {
		fmt.Printf(`Error conectando DB %s`, err.Error())
		return nil, err
	}

	return

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
		return response, err
	}

	if !isStored {
		fmt.Println("No previous score was found")

		score, err := CalculateScore(conn, reqBody.Value, reqBody.Type, score, isStored)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, err
		}

		err = SaveNewPerson(conn, score, reqBody.Value)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, err
		}

		photo, err := GetPersonPhoto(conn, reqBody.Value, reqBody.Type)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, err
		}

		response.Body = GetResponseBody(score, reqBody.Value, photo)
		response.StatusCode = http.StatusOK
		return response, nil
	}

	fmt.Println("Found previous score")

	score, err = CalculateScore(conn, reqBody.Value, reqBody.Type, score, isStored)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	photo, err := GetPersonPhoto(conn, reqBody.Value, reqBody.Type)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	response.Body = GetResponseBody(score, reqBody.Value, photo)
	response.StatusCode = http.StatusOK
	return response, nil
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
		return response, err
	}

	score, err = CalculateScore(conn, document, "CC", score, true)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	photo, err := GetPersonPhoto(conn, reqBody.Value, reqBody.Type)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	response.Body = GetResponseBody(score, reqBody.Value, photo)
	response.StatusCode = http.StatusOK
	return response, nil
}

func CalculateScore(conn *sql.DB, document, documentType string, data Score, existing bool) (score Score, err error) {

	reputation := data.Reputation
	scoreAverage, err := CalculateInternalScore(conn, document)
	if err != nil {
		return score, err
	}

	elapsed, err := DaysSinceLastUpdate(score.Updated)
	if err != nil {
		return
	}

	if elapsed > 7 || !existing {
		fmt.Println("Calculating reputation")
		reputation, err = CalculateReputation(document, documentType)
		if err != nil {
			return score, err
		}
	}

	if existing {

		score = Score{
			Name:       data.Name,
			Lastname:   data.Lastname,
			Gender:     data.Gender,
			Score:      scoreAverage,
			Reputation: reputation,
			Stars:      int((scoreAverage + reputation) / 40),
		}

		return
	} else {
		name, lastname, err := GetAssociatedName(document, documentType)
		if err != nil {
			fmt.Printf(`CalculateScore(1): %s`, err.Error())
			return score, err
		}

		score = Score{
			Name:       name,
			Lastname:   lastname,
			Gender:     `male`,
			Score:      scoreAverage,
			Reputation: reputation,
			Stars:      int((scoreAverage + reputation) / 40),
		}
	}

	return score, nil
}

func ValidatePhone(conn *sql.DB, phone string) (score Score, document string, err error) {

	results, err := conn.Query(`SELECT p.document, name, lastname, gender, score, reputation , stars, last_update FROM user u 
								INNER JOIN person p ON u.document =p.document 
								WHERE u.phone = ?`, phone)
	if err != nil {
		fmt.Printf(`ValidatePhone(1): %s`, err.Error())
		return
	}

	if results.Next() {
		err = results.Scan(&document, &score.Name, &score.Lastname, &score.Gender, &score.Score, &score.Reputation, &score.Stars, &score.Updated)
		if err != nil {
			fmt.Printf(`ValidatePhone(2): %s`, err.Error())
			return
		}
		return
	}

	return score, "", errors.New("número de celular no encontrado")
}

func GetStoredScore(conn *sql.DB, document string) (score Score, exists bool, err error) {
	results, err := conn.Query(`SELECT name, lastname, gender, score, reputation, stars, last_update FROM person p where p.document = ?`, document)
	if err != nil {
		fmt.Printf(`GetStoredScore(1): %s`, err.Error())
		return
	}

	if results.Next() {
		err = results.Scan(&score.Name, &score.Lastname, &score.Gender, &score.Score, &score.Reputation, &score.Stars, &score.Updated)
		if err != nil {
			fmt.Printf(`GetStoredScore(2): %s`, err.Error())
			return
		}

		exists = true
		return
	}

	return
}

func DaysSinceLastUpdate(lastUpdate string) (int, error) {

	if lastUpdate == "" {
		return 1, nil
	}

	lastUpdated, err := time.Parse(`2006-01-02 15:04:05`, lastUpdate)
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
	if err != nil {
		fmt.Printf(`GetAssociatedName(1) %s`, err.Error())
		return "", "", err
	}
	request.Header.Add(`Authorization`, bearer)

	client := &http.Client{}
	result, err := client.Do(request)
	if err != nil {
		fmt.Printf(`GetAssociatedName(2) %s`, err.Error())
		return "", "", err
	}
	defer result.Body.Close()

	data := &Person{}
	fmt.Println(result.StatusCode)

	err = json.NewDecoder(result.Body).Decode(data)
	if err != nil {
		fmt.Printf(`GetAssociatedName(3) %s`, err.Error())
		return "", "", err
	}

	return data.Data.FirstName, data.Data.Lastname, nil
}

func SaveNewPerson(conn *sql.DB, score Score, document string) error {

	query, err := conn.Prepare(`INSERT INTO person (document, name, lastname, gender, score, stars, reputation, last_update) VALUES(?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		fmt.Printf(`SaveNewPerson(1) %s`, err.Error())
		return err
	}

	data, err := query.Exec(document, score.Name, score.Lastname, ``, score.Score, score.Stars, score.Reputation)
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

func GetPersonPhoto(conn *sql.DB, dataValue, dataType string) (string, error) {
	photo := ``
	query := ``

	if dataType == `CC` {
		query = `SELECT photo FROM person p WHERE p.document = ?`
	} else {
		query = `SELECT p.photo FROM person p INNER JOIN user u ON p.document =u.document WHERE u.phone = ?`
	}

	results, err := conn.Query(query, dataValue)
	if err != nil {
		fmt.Printf(`GetPersonPhoto(1): %s`, err.Error())
		return "", err
	}

	if results.Next() {
		err = results.Scan(&photo)
		if err != nil {
			fmt.Printf(`GetPersonPhoto(2): %s`, err.Error())
			return "", err
		}
		return photo, nil
	}

	return "", nil
}

func GetResponseBody(score Score, document, photo string) string {
	certified := true
	fullname := fmt.Sprintf(`%s %s`, score.Name, score.Lastname)

	return fmt.Sprintf(`{
		"name": "%s",
		"document": "%s",
		"gender": "%s",
		"stars": %d,
		"reputation": %d,
		"score": %d,
		"certified": %t,
		"photo": "%s"
	}`, fullname, document, score.Gender, score.Stars, score.Reputation, score.Score, certified, photo)
}
