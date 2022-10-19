package main

import (
	"encoding/json"
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

// Person Services

func GetPersonByDocument(conn *gorm.DB, body GetByDocumentBody) (response events.APIGatewayProxyResponse, err error) {
	person := Person{
		Name: `not_found`,
	}
	result := conn.Select(body.Fields).Where(&Person{Document: body.Document}).Find(&person)
	if result.Error != nil {
		fmt.Printf(`GetPersonByDocument(1): %s`, result.Error.Error())
		return ErrorMessage(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrorMessage(errors.New(`no person found`))
	}

	jsonBody, _ := json.Marshal(person)

	response = SetResponseHeaders()
	response.StatusCode = http.StatusOK
	response.Body = string(jsonBody)

	return
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

	result := conn.Where(&Person{Document: person.Document}).Updates(&person)
	if result.Error != nil {
		return ErrorMessage(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrorMessage(errors.New(`no data was updated`))
	}

	return SuccessMessage(`Person data updated successfully`)
}

func GetNameByPhone(conn *gorm.DB, body GetByPhoneBody) (response events.APIGatewayProxyResponse, err error) {
	person := Person{}
	result := conn.Model(&Person{}).Select("name, lastname").
		Joins("inner join user on person.document = user.document").
		Where(`user.phone = ?`, body.Phone).
		Scan(&person)
	if result.Error != nil {
		fmt.Printf(`GetNameByPhone(1): %s`, result.Error.Error())
		return ErrorMessage(result.Error)
	}

	if result.RowsAffected == 0 {
		fmt.Printf(`GetNameByPhone(2): No user found`)
		return ErrorMessage(errors.New(`no user found`))
	}

	response = SetResponseHeaders()
	response.StatusCode = 200
	response.Body = fmt.Sprintf(`{"name":"%s","lastname":"%s"}`, person.Name, person.Lastname)
	return
}

// User Services

func GetUserByDocument(conn *gorm.DB, body GetByDocumentBody) (response events.APIGatewayProxyResponse, err error) {
	user := User{}
	result := conn.Select(body.Fields).Where(&Person{Document: body.Document}).Find(&user)
	if result.Error != nil {
		fmt.Printf(`GetUserByDocument(1): %s`, result.Error.Error())
		return ErrorMessage(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrorMessage(errors.New(`no user found`))
	}

	jsonBody, _ := json.Marshal(user)

	response = SetResponseHeaders()
	response.StatusCode = http.StatusOK
	response.Body = string(jsonBody)

	return
}

func GetAuthorId(conn *gorm.DB, document string) (id int, err error) {
	user := User{}
	data, _ := GetUserByDocument(conn,
		GetByDocumentBody{
			Document: document,
			Fields:   []string{`user_id`},
		})
	if data.StatusCode != 200 {
		return -1, errors.New(`author not found`)
	}

	json.Unmarshal([]byte(data.Body), &user)
	id = user.UserId

	return
}

func GetUserByPhone(conn *gorm.DB, body GetByPhoneBody) (response events.APIGatewayProxyResponse, err error) {
	user := User{}
	result := conn.Select(body.Fields).Where(&User{Phone: body.Phone}).Find(&user)
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

func CheckUserByDocument(conn *gorm.DB, body GetByDocumentBody) (response events.APIGatewayProxyResponse, err error) {
	user := User{}
	data, _ := GetUserByDocument(conn,
		GetByDocumentBody{
			Document: body.Document,
			Fields:   []string{},
		})

	json.Unmarshal([]byte(data.Body), &user)

	if user.Document == "" {
		return SuccessMessage(`user does not exists`)
	}
	return SuccessMessage(`user already exists`)
}

func InsertUser(conn *gorm.DB, user User) (response events.APIGatewayProxyResponse, err error) {
	result := conn.Create([]User{user})
	if result.Error != nil {
		fmt.Printf("InsertUser(1) %s", result.Error.Error())
		return ErrorMessage(result.Error)
	}

	return SuccessMessage(`User inserted successfully`)
}

func GetAccountData(conn *gorm.DB, body GetByDocumentBody) (response events.APIGatewayProxyResponse, err error) {

	account := Account{}

	result := conn.Model(&User{}).Select(
		`CONCAT(name,' ',lastname) name`,
		`email`,
		`phone`,
		`photo`,
		`gender`,
	).Joins(`INNER JOIN person ON person.document = user.document`).
		Where(&User{Document: body.Document}).
		Scan(&account)

	if result.Error != nil {
		return ErrorMessage(result.Error)
	}

	response = SetResponseHeaders()
	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{"name":"%s","email":"%s","phone":"%s","photo":"%s","gender":"%s"}`,
		account.Name, account.Email, account.Phone, account.Photo, account.Gender)
	return response, nil
}

func GetDocumentByPhone(conn *gorm.DB, body GetByPhoneBody) (response events.APIGatewayProxyResponse, err error) {
	user := User{}
	result := conn.Model(&User{}).Select("document").
		Where(&User{Phone: body.Phone}).
		Scan(&user)
	if result.Error != nil {
		fmt.Printf(`GetDocumentByPhone(1): %s`, result.Error.Error())
		return ErrorMessage(result.Error)
	}

	if result.RowsAffected == 0 {
		fmt.Printf(`GetDocumentByPhone(2): No user found`)
		return ErrorMessage(errors.New(`no person found`))
	}

	response = SetResponseHeaders()
	response.StatusCode = 200
	response.Body = fmt.Sprintf(`{"document":"%s"}`, user.Document)
	return
}

//Score Services

func GetScoreByDocument(conn *gorm.DB, body GetByDocumentBody) (response events.APIGatewayProxyResponse, err error) {
	internalScore := InternalScore{}

	rows, err := conn.Model(&Score{}).
		Select(`avg(score)`,
			`sum(case when score > 50 then 1 else 0 end) positiveScores`,
			`sum(case when score < 50 then 1 else 0 end) negativeScores`,
			`avg(case when DATEDIFF(CURRENT_TIMESTAMP,creation_date) < 61 then score else null end) Average60Days`,
		).Where(&Score{
		Objective: body.Document}).Rows()

	if err != nil {
		fmt.Printf(`GetInternalScoreSummary(1): %s`, err.Error())
		return ErrorMessage(err)
	}

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
