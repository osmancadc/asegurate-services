package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
)

func ErrorMessage(functionError error) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()
	response.StatusCode = http.StatusInternalServerError
	response.Body = fmt.Sprintf(`{"message":"%s"}`, functionError.Error())

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

func GetPersonData(data RequestGetData) (events.APIGatewayProxyResponse, error) {
	authorizationToken := os.Getenv("AUTHORIZATION_TOKEN")
	baseUrl := os.Getenv("BASE_URL")
	client := &http.Client{}
	person := Person{}

	url := fmt.Sprintf(`%s/cedula/extra?documentType=CC&documentNumber=%s&date=%s`, baseUrl, data.Document, data.ExpeditionDate)
	bearerToken := fmt.Sprintf(`Bearer %s`, authorizationToken)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf(`GetPersonExtraData(1) %s`, err.Error())
		return ErrorMessage(err)
	}

	request.Header.Add(`Authorization`, bearerToken)
	result, err := client.Do(request)
	if err != nil {
		fmt.Printf(`GetPersonExtraData(2) %s`, err.Error())
		return ErrorMessage(err)
	}
	defer result.Body.Close()

	err = json.NewDecoder(result.Body).Decode(&person)
	if err != nil {
		fmt.Printf(`GetPersonExtraData(3) %s`, err.Error())
		return ErrorMessage(err)
	}

	response := SetResponseHeaders()
	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{ "full_name":"%s", "name":"%s", "lastname":"%s", "gender":"%s", "is_alive":%t}`,
		person.Data.FullName, person.Data.Name, person.Data.Lastname, person.Data.Gender, person.Data.IsAlive)

	return response, nil
}

func GetPersonName(data RequestGetName) (events.APIGatewayProxyResponse, error) {

	authorizationToken := os.Getenv("AUTHORIZATION_TOKEN")
	baseUrl := os.Getenv("BASE_URL")
	client := &http.Client{}
	person := &Person{}

	url := fmt.Sprintf(`%s/cedula?documentType=%s&documentNumber=%s`, baseUrl, data.DocumentType, data.Document)
	bearerToken := fmt.Sprintf(`Bearer %s`, authorizationToken)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf(`GetPersonName(1) %s`, err.Error())
		return ErrorMessage(err)
	}
	request.Header.Add(`Authorization`, bearerToken)

	result, err := client.Do(request)
	if err != nil {
		fmt.Printf(`GetPersonName(2) %s`, err.Error())
		return ErrorMessage(err)

	}
	defer result.Body.Close()

	err = json.NewDecoder(result.Body).Decode(person)
	if err != nil {
		fmt.Printf(`GetPersonName(3) %s`, err.Error())
		return ErrorMessage(err)
	}

	response := SetResponseHeaders()
	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{"name":"%s","lastname":"%s"}`, person.Data.Name, person.Data.Lastname)

	return response, nil
}

func GetProccedings(data RequestGetProccedings) (events.APIGatewayProxyResponse, error) {
	authorizationToken := os.Getenv("AUTHORIZATION_TOKEN")
	baseUrl := os.Getenv("BASE_URL")
	client := &http.Client{}
	proccedings := &ProccedingsResponse{}

	url := fmt.Sprintf(`%s/procesos?documentType=%s&documentNumber=%s`, baseUrl, data.DocumentType, data.Document)
	bearerToken := fmt.Sprintf(`Bearer %s`, authorizationToken)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf(`GetProccedings(1) %s`, err.Error())
		return ErrorMessage(err)
	}
	request.Header.Add(`Authorization`, bearerToken)

	result, err := client.Do(request)
	if err != nil {
		fmt.Printf(`GetProccedings(2) %s`, err.Error())
		return ErrorMessage(err)

	}
	defer result.Body.Close()

	err = json.NewDecoder(result.Body).Decode(proccedings)
	if err != nil {
		fmt.Printf(`GetPersonName(3) %s`, err.Error())
		return ErrorMessage(err)
	}

	response := SetResponseHeaders()
	response.StatusCode = http.StatusOK

	total := proccedings.Data.Record.Total
	response.Body = fmt.Sprintf(`{"formal_complaints":%d,"recent_complain_year":%d,"five_years_amount":%d}`,
		total, total, total)

	return response, nil
}
