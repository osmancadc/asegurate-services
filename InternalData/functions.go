package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
)

var ConnectDatabase = func() (connection *sql.DB, err error) {
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

func ErrorMessage(functionError error) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()

	response.StatusCode = http.StatusInternalServerError
	response.Body = fmt.Sprintf(`{"message":"%s"}`, functionError.Error())

	return
}

func SuccessMessage(message string) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()
	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{"message":"%s"}`, message)

	return
}

func SetResponseHeaders() (response events.APIGatewayProxyResponse) {
	response = events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}
	return
}

func GetAuthorId(conn *sql.DB, document string) (int, error) {

	id := 0

	results, err := conn.Query(`SELECT user_id FROM user u WHERE u.document = ?`, document)
	if err != nil {
		fmt.Printf(`GetAuthorId(1): %s`, err.Error())
		return -1, err
	}

	if results.Next() {
		err = results.Scan(&id)
		if err != nil {
			fmt.Printf(`GetAuthorId(2): %s`, err.Error())
			return -1, err
		}
		return id, nil
	}

	return -1, errors.New("no user found")
}

func GetInternalScoreSummary(conn *sql.DB, body GetScoreBody) (response events.APIGatewayProxyResponse, err error) {
	results, err := conn.Query(`SELECT avg(s.score),
							sum(case when s.score > 50 then 1 else 0 end) positiveScores,
							sum(case when s.score < 50 then 1 else 0 end) negativeScores,
							avg(case 
								when DATEDIFF(CURRENT_TIMESTAMP,s.creation_date) < 61 
									then s.score 
									else null 
								end) Average60Days
							FROM score s 
							WHERE s.objective = ?`, body.Document)
	if err != nil {
		fmt.Printf(`GetInternalScore(1): %s`, err.Error())
		return ErrorMessage(err)
	}

	internalScore := InternalScore{}

	if results.Next() {
		err = results.Scan(&internalScore.Score, &internalScore.PositiveScores, &internalScore.NegativeScores, &internalScore.Average60Days)
		if err != nil {
			fmt.Printf(`GetInternalScore(2): %s`, err.Error())
			return ErrorMessage(err)
		}
	}

	response = SetResponseHeaders()
	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{ "score": %F, "positive_scores": %d, "negative_scores":%d, "average_60_days":%F }`,
		internalScore.Score, internalScore.PositiveScores, internalScore.NegativeScores, internalScore.Average60Days)
	return response, nil
}

func UpdateInternalScore(conn *sql.DB, body UpdateScoreBody) (response events.APIGatewayProxyResponse, err error) {
	query, err := conn.Prepare(`UPDATE person
								SET score=?, reputation=?, last_update=CURRENT_TIMESTAMP
								WHERE document= ?`)
	if err != nil {
		fmt.Printf("UpdateInternalScore(1) %s", err.Error())
		return ErrorMessage(err)
	}

	_, err = query.Exec(body.Score, body.Reputation, body.Document)
	if err != nil {
		fmt.Printf("UpdateInternalScore(2) %s", err.Error())
		return ErrorMessage(err)
	}

	return SuccessMessage(`User score updated successfully`)
}

func InsertInternalScore(conn *sql.DB, body InsertScoreBody) (response events.APIGatewayProxyResponse, err error) {

	authorId, err := GetAuthorId(conn, body.Author)
	if err != nil {
		return ErrorMessage(err)
	}

	query, err := conn.Prepare(`INSERT INTO score (author, objective, score, comments) VALUES(?, ?, ?, ?)`)
	if err != nil {
		fmt.Printf("InsertInternalScore(1) %s", err.Error())
		return ErrorMessage(err)
	}

	_, err = query.Exec(authorId, body.Objective, body.Score, body.Comments)
	if err != nil {
		fmt.Printf("InsertInternalScore(2) %s", err.Error())
		return ErrorMessage(err)
	}

	return SuccessMessage(`Score uploaded successfully`)
}

func GetUserByPhone(conn *sql.DB, body GetUserByPhoneBody) (response events.APIGatewayProxyResponse, err error) {
	objective := ""

	results, err := conn.Query(`SELECT document FROM user u WHERE u.phone = ?`, body.Phone)
	if err != nil {
		fmt.Printf(`UploadScorePhone(1): %s`, err.Error())
		return ErrorMessage(err)
	}

	if results.Next() {
		err = results.Scan(&objective)
		if err != nil {
			fmt.Printf(`UploadScorePhone(2): %s`, err.Error())
			return ErrorMessage(err)
		}
	}

	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{"document":"%s"}`, objective)
	return response, nil
}

func GetPersonByDocument(conn *sql.DB, body GetByDocumentBody) (response events.APIGatewayProxyResponse, err error) {
	name := "not_found"
	gender := "not_found"

	results, err := conn.Query(`SELECT name, gender FROM person where document = ?`, body.Document)
	if err != nil {
		fmt.Printf(`GetPersonByDocument(1): %s`, err.Error())
		return ErrorMessage(err)
	}

	if results.Next() {
		err = results.Scan(&name, &gender)
		if err != nil {
			fmt.Printf(`GetPersonByDocument(2): %s`, err.Error())
			return ErrorMessage(err)
		}
	}

	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{"name":"%s","gender":"%s"}`, name, gender)
	return response, nil
}

func GetUserByDocument(conn *sql.DB, body GetByDocumentBody) (response events.APIGatewayProxyResponse, err error) {
	results, err := conn.Query(`SELECT user_id FROM user WHERE document =  ?`, body.Document)
	if err != nil {
		fmt.Printf(`CheckExistingUser(1): %s`, err.Error())
		return ErrorMessage(err)
	}

	if results.Next() {
		fmt.Printf("CheckExistingUser(2) el usuario ya existe")
		return SuccessMessage(`user already exists`)
	}

	return SuccessMessage(`user does not exists`)
}

func InsertUser(conn *sql.DB, body InsertUserBody) (response events.APIGatewayProxyResponse, err error) {
	query, err := conn.Prepare(`INSERT INTO user (email, phone, password, document, role) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Printf("InsertUser(1) %s", err.Error())
		return ErrorMessage(err)
	}

	_, err = query.Exec(body.Email, body.Phone, body.Password, body.Document, body.Role)
	if err != nil {
		fmt.Printf("InsertUser(2) %s", err.Error())
		return ErrorMessage(err)
	}

	return SuccessMessage(`User inserted successfully`)
}

func InsertPerson(conn *sql.DB, body InsertPersonBody) (response events.APIGatewayProxyResponse, err error) {

	query, err := conn.Prepare(`INSERT INTO person  (document, name, lastname, gender, score, reputation, photo, last_update) 
								VALUES(?, ?, ?, ?, 50, 50, '', CURRENT_TIMESTAMP)`)
	if err != nil {
		fmt.Printf("InsertPerson(1) %s", err.Error())
		return ErrorMessage(err)
	}

	_, err = query.Exec(body.Document, body.Name, body.Lastname, body.Gender)
	if err != nil {
		fmt.Printf("InsertPerson(2) %s", err.Error())
		return ErrorMessage(err)
	}
	return SuccessMessage(`Person inserted successfully`)
}
