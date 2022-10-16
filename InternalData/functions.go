package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
)

var ConnectDatabase = func() (db *gorm.DB, err error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s`, user, password, host, database)

	db, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			// Logger: logger.Discard,
		},
	)
	if err != nil {
		fmt.Printf(`Error conectando DB %s`, err.Error())
		return
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

func GetAuthorId(conn *gorm.DB, document string) (int, error) {
	user := User{}
	result := conn.Select(`user_id`).Where(`document = ?`, document).Find(&user)
	if result.Error != nil {
		fmt.Printf(`GetAuthorId(1): %s`, result.Error.Error())
		return -1, result.Error
	}

	if result.RowsAffected > 0 {
		return user.UserId, nil
	}

	return -1, errors.New("author not found")
}

func GetPersonByDocument(conn *gorm.DB, body GetByDocumentBody) (response events.APIGatewayProxyResponse, err error) {
	person := Person{
		Name: `not_found`,
	}
	result := conn.Select(`name`, `gender`).Where(`document = ?`, body.Document).Find(&person)
	if result.Error != nil {
		fmt.Printf(`GetPersonByDocument(1): %s`, result.Error.Error())
		return ErrorMessage(result.Error)
	}

	if result.RowsAffected > 0 {
		response = SetResponseHeaders()
		response.StatusCode = http.StatusOK
		response.Body = fmt.Sprintf(`{"name":"%s","gender":"%s"}`, person.Name, person.Gender)

		return
	}

	return ErrorMessage(errors.New(`no person found`))
}

func GetScoreByDocument(conn *gorm.DB, body GetByDocumentBody) (response events.APIGatewayProxyResponse, err error) {

	rows, err := conn.Raw(`SELECT avg(s.score),
							sum(case when s.score > 50 then 1 else 0 end) positiveScores,
							sum(case when s.score < 50 then 1 else 0 end) negativeScores,
							avg(case when DATEDIFF(CURRENT_TIMESTAMP,s.creation_date) < 61 then s.score else null end) Average60Days
							FROM score s 
							WHERE s.objective = ?`, body.Document).Rows()
	if err != nil {
		fmt.Printf(`GetInternalScoreSummary(1): %s`, err.Error())
		return ErrorMessage(err)
	}

	internalScore := InternalScore{}

	if rows.Next() {
		err = rows.Scan(&internalScore.Score, &internalScore.PositiveScores, &internalScore.NegativeScores, &internalScore.Average60Days)
		if err != nil {
			fmt.Printf(`GetInternalScoreSummary(2): %s`, err.Error())
			return ErrorMessage(err)
		}
	}

	response = SetResponseHeaders()
	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{ "score": %F, "positive_scores": %d, "negative_scores":%d, "average_60_days":%F }`,
		internalScore.Score, internalScore.PositiveScores, internalScore.NegativeScores, internalScore.Average60Days)
	return response, nil
}

func GetUserByPhone(conn *gorm.DB, body GetByPhoneBody) (response events.APIGatewayProxyResponse, err error) {

	user := User{}
	result := conn.Select(`document`).Where(`phone = ?`, body.Phone).Find(&user)
	if result.Error != nil {
		fmt.Printf(`GetUserByPhone(1): %s`, result.Error.Error())
		return ErrorMessage(result.Error)
	}

	if result.RowsAffected > 0 {
		response.StatusCode = http.StatusOK
		response.Body = fmt.Sprintf(`{"document":"%s"}`, user.Document)
		return response, nil
	}

	return ErrorMessage(errors.New(`user not found`))
}

func GetUserByDocument(conn *gorm.DB, body GetByDocumentBody) (response events.APIGatewayProxyResponse, err error) {
	user := User{}
	result := conn.Select(`user_id`).Where(`document = ?`, body.Document).Find(&user)
	if result.Error != nil {
		fmt.Printf(`GetUserByDocument(1): %s`, result.Error.Error())
		return ErrorMessage(result.Error)
	}

	if result.RowsAffected > 0 {
		return SuccessMessage(`user already exists`)
	}

	return SuccessMessage(`user does not exists`)
}

func InsertScore(conn *gorm.DB, body ScoreBody) (response events.APIGatewayProxyResponse, err error) {

	authorId, err := GetAuthorId(conn, body.Author)
	if err != nil {
		return ErrorMessage(err)
	}

	score := Score{
		Author:    authorId,
		Objective: body.Objective,
		Score:     body.Score,
		Comments:  body.Comments,
	}

	result := conn.Create([]Score{score})
	if result.Error != nil {
		return ErrorMessage(result.Error)
	}

	return SuccessMessage(`Score uploaded successfully`)
}

func InsertUser(conn *gorm.DB, user User) (response events.APIGatewayProxyResponse, err error) {
	result := conn.Create([]User{user})
	if result.Error != nil {
		fmt.Printf("InsertUser(1) %s", result.Error.Error())
		return ErrorMessage(result.Error)
	}

	return SuccessMessage(`User inserted successfully`)
}

func InsertPerson(conn *gorm.DB, person Person) (response events.APIGatewayProxyResponse, err error) {

	result := conn.Create([]Person{person})
	if result.Error != nil {
		fmt.Printf("InsertUser(1) %s", result.Error.Error())
		return ErrorMessage(result.Error)
	}

	return SuccessMessage(`Person inserted successfully`)
}

func UpdatePerson(conn *gorm.DB, person Person) (response events.APIGatewayProxyResponse, err error) {

	result := conn.Where(`document = ?`, person.Document).Updates(&person)
	if result.Error != nil {
		return ErrorMessage(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrorMessage(errors.New(`no data was updated`))
	}

	return SuccessMessage(`Person data updated successfully`)
}
